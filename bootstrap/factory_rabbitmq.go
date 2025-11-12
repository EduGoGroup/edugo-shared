package bootstrap

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// =============================================================================
// RABBITMQ FACTORY IMPLEMENTATION
// =============================================================================

// DefaultRabbitMQFactory implementa RabbitMQFactory
type DefaultRabbitMQFactory struct {
	connectionTimeout time.Duration
}

// NewDefaultRabbitMQFactory crea una nueva instancia de DefaultRabbitMQFactory
func NewDefaultRabbitMQFactory() *DefaultRabbitMQFactory {
	return &DefaultRabbitMQFactory{
		connectionTimeout: 10 * time.Second,
	}
}

// CreateConnection crea una conexión a RabbitMQ
func (f *DefaultRabbitMQFactory) CreateConnection(ctx context.Context, config RabbitMQConfig) (*amqp.Connection, error) {
	// Intentar conexión con timeout
	connChan := make(chan *amqp.Connection, 1)
	errChan := make(chan error, 1)

	go func() {
		conn, err := amqp.Dial(config.URL)
		if err != nil {
			errChan <- err
			return
		}
		connChan <- conn
	}()

	select {
	case conn := <-connChan:
		return conn, nil
	case err := <-errChan:
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	case <-time.After(f.connectionTimeout):
		return nil, fmt.Errorf("connection timeout after %v", f.connectionTimeout)
	case <-ctx.Done():
		return nil, fmt.Errorf("connection cancelled: %w", ctx.Err())
	}
}

// CreateChannel crea un canal de comunicación
func (f *DefaultRabbitMQFactory) CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	// Configurar QoS para el canal
	if err := channel.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		channel.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return channel, nil
}

// DeclareQueue declara una cola con configuración por defecto
func (f *DefaultRabbitMQFactory) DeclareQueue(channel *amqp.Channel, queueName string) (amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-message-ttl":             int32(3600000), // 1 hour
			"x-max-priority":            int32(10),      // max priority
			"x-queue-mode":              "lazy",         // lazy mode
			"x-dead-letter-exchange":    "",             // sin DLX por defecto
			"x-dead-letter-routing-key": "",
		},
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	return queue, nil
}

// Close cierra el canal y la conexión
func (f *DefaultRabbitMQFactory) Close(channel *amqp.Channel, conn *amqp.Connection) error {
	var errs []error

	if channel != nil {
		if err := channel.Close(); err != nil && err != amqp.ErrClosed {
			errs = append(errs, fmt.Errorf("failed to close channel: %w", err))
		}
	}

	if conn != nil {
		if err := conn.Close(); err != nil && err != amqp.ErrClosed {
			errs = append(errs, fmt.Errorf("failed to close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}

	return nil
}

// Verificar que DefaultRabbitMQFactory implementa RabbitMQFactory
var _ RabbitMQFactory = (*DefaultRabbitMQFactory)(nil)
