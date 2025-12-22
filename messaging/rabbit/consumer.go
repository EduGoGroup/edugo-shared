package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// MessageHandler es la función que procesa un mensaje
type MessageHandler func(ctx context.Context, body []byte) error

// Consumer interface para consumir mensajes
type Consumer interface {
	Consume(ctx context.Context, queueName string, handler MessageHandler) error
	Wait() error
	Stop()
	Errors() <-chan error
	IsRunning() bool
	Close() error
}

// RabbitMQConsumer implementación de Consumer para RabbitMQ con control de goroutines.
//
// Nota sobre errChan: El canal de errores tiene buffer de tamaño 1 y solo reporta
// el primer error asíncrono que ocurra. Errores adicionales se descartan silenciosamente.
// Esto es adecuado para el caso de uso actual donde un error fatal detiene el consumer.
type RabbitMQConsumer struct {
	conn     *Connection
	config   ConsumerConfig
	wg       sync.WaitGroup
	errChan  chan error // Buffer de 1: solo el primer error asíncrono se reporta
	stopCh   chan struct{}
	stopOnce sync.Once
	mu       sync.Mutex
	running  bool
}

// NewConsumer crea un nuevo Consumer con control de goroutines.
func NewConsumer(conn *Connection, config ConsumerConfig) Consumer {
	return &RabbitMQConsumer{
		conn:    conn,
		config:  config,
		errChan: make(chan error, 1),
		stopCh:  make(chan struct{}),
	}
}

// Consume inicia el consumo de mensajes de una cola.
func (c *RabbitMQConsumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("consumer already running")
	}
	c.running = true
	c.mu.Unlock()

	// Obtener canal de mensajes
	msgs, err := c.conn.GetChannel().Consume(
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
					select {
					case c.errChan <- fmt.Errorf("message channel closed unexpectedly"):
					default:
					}
					return
				}

				// Convertir amqp.Delivery a nuestro struct interno
				delivery := amqpDelivery{
					Body:        msg.Body,
					DeliveryTag: msg.DeliveryTag,
					Ack:         msg.Ack,
					Nack:        msg.Nack,
				}
				c.processBasicMessage(ctx, queueName, delivery, handler)
			}
		}
	}()

	return nil
}

// processBasicMessage procesa un mensaje individual con manejo de acknowledgment (sin soporte DLQ).
func (c *RabbitMQConsumer) processBasicMessage(ctx context.Context, queueName string, delivery amqpDelivery, handler MessageHandler) {
	// Procesar mensaje
	err := handler(ctx, delivery.Body)

	// Manejar acknowledgment si no es auto-ack
	if !c.config.AutoAck {
		if err != nil {
			// Nack con requeue si hubo error en el handler
			if nackErr := delivery.Nack(false, true); nackErr != nil {
				log.Printf("[ERROR] failed to nack message (delivery_tag=%d, queue=%s): %v (original error: %v)",
					delivery.DeliveryTag, queueName, nackErr, err)
			}
		} else {
			// Ack si el procesamiento fue exitoso
			if ackErr := delivery.Ack(false); ackErr != nil {
				log.Printf("[ERROR] failed to ack message (delivery_tag=%d, queue=%s): %v",
					delivery.DeliveryTag, queueName, ackErr)
			}
		}
	}
}

// amqpDelivery es un struct interno para abstraer mensajes AMQP
type amqpDelivery struct {
	Body        []byte
	DeliveryTag uint64
	Ack         func(multiple bool) error
	Nack        func(multiple bool, requeue bool) error
}

// Wait bloquea hasta que el consumer se detenga y retorna cualquier error asíncrono.
func (c *RabbitMQConsumer) Wait() error {
	done := make(chan struct{})

	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case err := <-c.errChan:
		<-done // Asegurar que WaitGroup completó
		return err
	case <-done:
		// Intentar leer error después de que goroutines terminaron
		select {
		case err := <-c.errChan:
			return err
		default:
			return nil
		}
	}
}

// Stop detiene el consumer de forma graceful.
func (c *RabbitMQConsumer) Stop() {
	c.stopOnce.Do(func() {
		close(c.stopCh)
	})
}

// Errors retorna un canal para recibir errores del consumer.
func (c *RabbitMQConsumer) Errors() <-chan error {
	return c.errChan
}

// IsRunning retorna true si el consumer está actualmente ejecutándose.
func (c *RabbitMQConsumer) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

// Close detiene y limpia el consumer.
func (c *RabbitMQConsumer) Close() error {
	c.Stop()
	c.wg.Wait()
	return nil
}

// UnmarshalMessage helper para deserializar un mensaje JSON.
func UnmarshalMessage(body []byte, v interface{}) error {
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return nil
}

// HandleWithUnmarshal helper que combina unmarshal y handling.
func HandleWithUnmarshal(body []byte, v interface{}, handler func(interface{}) error) error {
	if err := UnmarshalMessage(body, v); err != nil {
		return err
	}
	return handler(v)
}
