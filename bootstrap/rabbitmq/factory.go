package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-shared/bootstrap"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Factory implementa la creacion de conexiones RabbitMQ.
type Factory struct {
	connectionTimeout time.Duration
}

// NewFactory crea una nueva Factory de RabbitMQ.
func NewFactory() *Factory {
	return &Factory{
		connectionTimeout: 10 * time.Second,
	}
}

// CreateConnection crea una conexion a RabbitMQ con timeout.
func (f *Factory) CreateConnection(ctx context.Context, cfg bootstrap.RabbitMQConfig) (*amqp.Connection, error) {
	connChan := make(chan *amqp.Connection, 1)
	errChan := make(chan error, 1)

	go func() {
		conn, err := amqp.Dial(cfg.URL)
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
		return nil, fmt.Errorf("bootstrap/rabbitmq: connect: %w", err)
	case <-time.After(f.connectionTimeout):
		return nil, fmt.Errorf("bootstrap/rabbitmq: connection timeout after %v", f.connectionTimeout)
	case <-ctx.Done():
		return nil, fmt.Errorf("bootstrap/rabbitmq: cancelled: %w", ctx.Err())
	}
}

// CreateChannel crea un canal AMQP con QoS configurado.
func (f *Factory) CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("bootstrap/rabbitmq: create channel: %w", err)
	}

	if err := channel.Qos(10, 0, false); err != nil {
		if closeErr := channel.Close(); closeErr != nil {
			return nil, errors.Join(
				fmt.Errorf("bootstrap/rabbitmq: set QoS: %w", err),
				fmt.Errorf("bootstrap/rabbitmq: close channel: %w", closeErr),
			)
		}
		return nil, fmt.Errorf("bootstrap/rabbitmq: set QoS: %w", err)
	}

	return channel, nil
}

// DeclareQueue declara una cola con configuracion por defecto.
func (f *Factory) DeclareQueue(channel *amqp.Channel, queueName string) (amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-message-ttl":             int32(3600000), // 1 hour
			"x-max-priority":            int32(10),
			"x-queue-mode":              "lazy",
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": "",
		},
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("bootstrap/rabbitmq: declare queue %s: %w", queueName, err)
	}
	return queue, nil
}

// Close cierra el canal y la conexion.
func (f *Factory) Close(channel *amqp.Channel, conn *amqp.Connection) error {
	var errs []error

	if channel != nil {
		if err := channel.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
			errs = append(errs, fmt.Errorf("bootstrap/rabbitmq: close channel: %w", err))
		}
	}

	if conn != nil {
		if err := conn.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
			errs = append(errs, fmt.Errorf("bootstrap/rabbitmq: close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}
	return nil
}
