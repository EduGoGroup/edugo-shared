package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher interface para publicar mensajes
type Publisher interface {
	Publish(ctx context.Context, exchange, routingKey string, body interface{}) error
	PublishWithPriority(ctx context.Context, exchange, routingKey string, body interface{}, priority uint8) error
	Close() error
}

// RabbitMQPublisher implementación de Publisher para RabbitMQ
type RabbitMQPublisher struct {
	conn *Connection
}

// NewPublisher crea un nuevo Publisher
func NewPublisher(conn *Connection) Publisher {
	return &RabbitMQPublisher{
		conn: conn,
	}
}

// Publish publica un mensaje en un exchange con routing key
func (p *RabbitMQPublisher) Publish(ctx context.Context, exchange, routingKey string, body interface{}) error {
	return p.PublishWithPriority(ctx, exchange, routingKey, body, 0)
}

// PublishWithPriority publica un mensaje con prioridad específica
func (p *RabbitMQPublisher) PublishWithPriority(ctx context.Context, exchange, routingKey string, body interface{}, priority uint8) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// Serializar el body a JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Crear mensaje
	msg := amqp.Publishing{
		ContentType:  "application/json",
		Body:         bodyBytes,
		DeliveryMode: amqp.Persistent, // Persistent message
		Timestamp:    time.Now(),
		Priority:     priority,
	}

	// Publicar con contexto
	err = p.conn.GetChannel().PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		msg,
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Close cierra el publisher (mantiene la conexión abierta)
func (p *RabbitMQPublisher) Close() error {
	// El publisher no cierra la conexión, eso lo maneja el caller
	return nil
}
