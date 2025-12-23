package rabbit

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumeWithDLQ consume mensajes con soporte para Dead Letter Queue.
// Implementa el mismo control de ciclo de vida que Consume.
func (c *RabbitMQConsumer) ConsumeWithDLQ(ctx context.Context, queueName string, handler MessageHandler) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("consumer already running")
	}
	c.running = true
	c.mu.Unlock()

	// Configurar prefetch si está configurado
	ch := c.conn.GetChannel()
	prefetchCount := c.config.PrefetchCount
	if prefetchCount == 0 {
		prefetchCount = DefaultPrefetchCount
	}
	if err := ch.Qos(
		prefetchCount,
		0,
		false,
	); err != nil {
		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Declarar DLX y DLQ si está habilitado
	if c.config.DLQ.Enabled {
		if err := c.setupDLQ(ch); err != nil {
			c.mu.Lock()
			c.running = false
			c.mu.Unlock()
			return fmt.Errorf("failed to setup DLQ: %w", err)
		}
	}

	// Consumir mensajes
	msgs, err := ch.Consume(
		queueName,
		c.config.Name,
		c.config.AutoAck,
		c.config.Exclusive,
		c.config.NoLocal,
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	workerLimit := prefetchCount
	if workerLimit <= 0 {
		workerLimit = DefaultPrefetchCount
	}
	semaphore := make(chan struct{}, workerLimit)

	// Registrar goroutine principal en WaitGroup
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer func() {
			c.mu.Lock()
			c.running = false
			c.mu.Unlock()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.stopCh:
				return
			case msg, ok := <-msgs:
				if !ok {
					// Canal cerrado inesperadamente
					select {
					case c.errChan <- fmt.Errorf("message channel closed unexpectedly"):
					default:
					}
					return
				}

				semaphore <- struct{}{}
				go func(delivery amqp.Delivery) {
					defer func() { <-semaphore }()
					c.processMessage(ctx, ch, queueName, handler, delivery)
				}(msg)
			}
		}
	}()

	return nil
}

// setupDLQ configura el Dead Letter Exchange y Queue
func (c *RabbitMQConsumer) setupDLQ(ch *amqp.Channel) error {
	// Declarar DLX (exchange para mensajes fallidos)
	if err := ch.ExchangeDeclare(
		c.config.DLQ.DLXExchange, // name
		"direct",                 // type
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	); err != nil {
		return fmt.Errorf("failed to declare DLX: %w", err)
	}

	// Declarar DLQ (queue para mensajes fallidos)
	_, err := ch.QueueDeclare(
		c.config.DLQ.DLXRoutingKey, // name (usa routing key como nombre)
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Bindear DLQ al DLX
	if err := ch.QueueBind(
		c.config.DLQ.DLXRoutingKey, // queue name
		c.config.DLQ.DLXRoutingKey, // routing key
		c.config.DLQ.DLXExchange,   // exchange
		false,                      // no-wait
		nil,                        // arguments
	); err != nil {
		return fmt.Errorf("failed to bind DLQ: %w", err)
	}

	return nil
}

// sendToDLQ envía un mensaje al Dead Letter Queue
func (c *RabbitMQConsumer) sendToDLQ(ch *amqp.Channel, msg amqp.Delivery) error {
	// Agregar metadata al mensaje
	headers := msg.Headers
	if headers == nil {
		headers = amqp.Table{}
	}
	headers["x-original-exchange"] = msg.Exchange
	headers["x-original-routing-key"] = msg.RoutingKey
	headers["x-failed-at"] = time.Now().Unix()
	headers["x-retry-count"] = getRetryCount(msg.Headers)

	// Publicar a DLX
	return ch.Publish(
		c.config.DLQ.DLXExchange,   // exchange
		c.config.DLQ.DLXRoutingKey, // routing key
		false,                      // mandatory
		false,                      // immediate
		amqp.Publishing{
			ContentType: msg.ContentType,
			Body:        msg.Body,
			Headers:     headers,
		},
	)
}

func (c *RabbitMQConsumer) processMessage(
	ctx context.Context,
	ch *amqp.Channel,
	queueName string,
	handler MessageHandler,
	msg amqp.Delivery,
) {
	retries := getRetryCount(msg.Headers)
	err := handler(ctx, msg.Body)

	if c.config.AutoAck {
		return
	}

	if err != nil {
		c.handleProcessingError(ctx, ch, queueName, msg, retries)
		return
	}

	_ = msg.Ack(false) //nolint:errcheck // Error de Ack en procesamiento exitoso se ignora
}

func (c *RabbitMQConsumer) handleProcessingError(
	ctx context.Context,
	ch *amqp.Channel,
	queueName string,
	msg amqp.Delivery,
	retries int,
) {
	nextRetry := retries + 1

	if c.config.DLQ.Enabled {
		if nextRetry > c.config.DLQ.MaxRetries {
			if err := c.sendToDLQ(ch, msg); err != nil {
				_ = msg.Nack(false, true) //nolint:errcheck // Error de Nack en fallback se ignora
			} else {
				_ = msg.Ack(false) //nolint:errcheck // Error de Ack en DLQ exitoso se ignora
			}
			return
		}

		backoff := c.config.DLQ.CalculateBackoff(retries)
		timer := time.NewTimer(backoff)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			_ = msg.Nack(false, true) //nolint:errcheck // Error de Nack en cancelación se ignora
			return
		case <-timer.C:
		}

		retryHeaders := cloneHeaders(msg.Headers)
		retryHeaders["x-retry-count"] = nextRetry

		publishErr := ch.PublishWithContext(
			ctx,
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				ContentType:  msg.ContentType,
				Body:         msg.Body,
				Headers:      retryHeaders,
				DeliveryMode: amqp.Persistent,
				Priority:     msg.Priority,
			},
		)
		if publishErr != nil {
			_ = msg.Nack(false, true) //nolint:errcheck // Error de Nack en fallo de publish se ignora
			return
		}

		_ = msg.Ack(false) //nolint:errcheck // Error de Ack en retry exitoso se ignora
		return
	}

	_ = msg.Nack(false, true) //nolint:errcheck // Error de Nack sin DLQ se ignora
}

// getRetryCount extrae el número de reintentos del header
func getRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}
	if count, ok := headers["x-retry-count"].(int); ok {
		return count
	}
	// Intentar con int32 (tipo que RabbitMQ puede usar)
	if count, ok := headers["x-retry-count"].(int32); ok {
		return int(count)
	}
	// Intentar con int64 (algunos clientes usan este tipo)
	if count, ok := headers["x-retry-count"].(int64); ok {
		return int(count)
	}
	return 0
}

// cloneHeaders crea una copia superficial de los headers para evitar mutar el original
func cloneHeaders(headers amqp.Table) amqp.Table {
	if headers == nil {
		return amqp.Table{}
	}
	copyHeaders := make(amqp.Table, len(headers))
	for k, v := range headers {
		copyHeaders[k] = v
	}
	return copyHeaders
}
