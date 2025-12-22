package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// MessageHandler es la función que procesa un mensaje
type MessageHandler func(ctx context.Context, body []byte) error

// Consumer interface para consumir mensajes
type Consumer interface {
	Consume(ctx context.Context, queueName string, handler MessageHandler) error
	Close() error
}

// RabbitMQConsumer implementación de Consumer para RabbitMQ
type RabbitMQConsumer struct {
	conn   *Connection
	config ConsumerConfig
}

// NewConsumer crea un nuevo Consumer
func NewConsumer(conn *Connection, config ConsumerConfig) Consumer {
	return &RabbitMQConsumer{
		conn:   conn,
		config: config,
	}
}

// Consume inicia el consumo de mensajes de una cola
func (c *RabbitMQConsumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
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

				// Procesar mensaje
				err := handler(ctx, msg.Body)

				// Manejar acknowledgment si no es auto-ack
				if !c.config.AutoAck {
					if err != nil {
						// Nack con requeue si hubo error en el handler
						if nackErr := msg.Nack(false, true); nackErr != nil {
							log.Printf("[ERROR] failed to nack message (delivery_tag=%d, queue=%s): %v (original error: %v)",
								msg.DeliveryTag, queueName, nackErr, err)
						}
					} else {
						// Ack si el procesamiento fue exitoso
						if ackErr := msg.Ack(false); ackErr != nil {
							log.Printf("[ERROR] failed to ack message (delivery_tag=%d, queue=%s): %v",
								msg.DeliveryTag, queueName, ackErr)
						}
					}
				}
			}
		}
	}()

	return nil
}

// Close cierra el consumer (mantiene la conexión abierta)
func (c *RabbitMQConsumer) Close() error {
	// El consumer no cierra la conexión, eso lo maneja el caller
	return nil
}

// UnmarshalMessage helper para deserializar un mensaje JSON
func UnmarshalMessage(body []byte, v interface{}) error {
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return nil
}

// HandleWithUnmarshal helper que combina unmarshal y handling
func HandleWithUnmarshal(body []byte, v interface{}, handler func(interface{}) error) error {
	if err := UnmarshalMessage(body, v); err != nil {
		return err
	}
	return handler(v)
}
