package rabbit

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumeWithDLQ consume mensajes con soporte para Dead Letter Queue
func (c *RabbitMQConsumer) ConsumeWithDLQ(ctx context.Context, queueName string, handler MessageHandler) error {
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
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Declarar DLX y DLQ si está habilitado
	if c.config.DLQ.Enabled {
		if err := c.setupDLQ(ch); err != nil {
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
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	// Procesar mensajes en un loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				// Obtener número de reintentos del header
				retries := getRetryCount(msg.Headers)

				// Procesar mensaje
				err := handler(ctx, msg.Body)

				// Manejar acknowledgment si no es auto-ack
				if !c.config.AutoAck {
					if err != nil {
						// Verificar si excedió reintentos
						if c.config.DLQ.Enabled && retries >= c.config.DLQ.MaxRetries {
							// Enviar a DLQ
							if err := c.sendToDLQ(ch, msg); err != nil {
								// Si falla envío a DLQ, reencolar como fallback
								_ = msg.Nack(false, true)
							} else {
								// Acknowledge porque ya está en DLQ
								_ = msg.Ack(false)
							}
						} else {
							// Reencolar con delay (si DLQ está habilitado)
							if c.config.DLQ.Enabled {
								backoff := c.config.DLQ.CalculateBackoff(retries)
								time.Sleep(backoff)
							}

							// Incrementar contador de reintentos en headers
							if msg.Headers == nil {
								msg.Headers = amqp.Table{}
							}

							// Nack con requeue
							_ = msg.Nack(false, true)
						}
					} else {
						// Procesado exitosamente
						_ = msg.Ack(false)
					}
				}
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
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	); err != nil {
		return fmt.Errorf("failed to declare DLX: %w", err)
	}

	// Declarar DLQ (queue para mensajes fallidos)
	_, err := ch.QueueDeclare(
		c.config.DLQ.DLXRoutingKey, // name (usa routing key como nombre)
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
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
	return 0
}
