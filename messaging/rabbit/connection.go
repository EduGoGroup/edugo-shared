package rabbit

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection encapsula la conexión a RabbitMQ
type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

// Connect establece una conexión a RabbitMQ
func Connect(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close() //nolint:errcheck // Ignorar error en cleanup, el error principal es el de channel
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Connection{
		conn:    conn,
		channel: channel,
		url:     url,
	}, nil
}

// GetChannel retorna el canal de RabbitMQ
func (c *Connection) GetChannel() *amqp.Channel {
	return c.channel
}

// GetConnection retorna la conexión de RabbitMQ
func (c *Connection) GetConnection() *amqp.Connection {
	return c.conn
}

// Close cierra la conexión y el canal
func (c *Connection) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}
	return nil
}

// IsClosed verifica si la conexión está cerrada
func (c *Connection) IsClosed() bool {
	return c.conn == nil || c.conn.IsClosed()
}

// DeclareExchange declara un exchange
func (c *Connection) DeclareExchange(cfg ExchangeConfig) error {
	return c.channel.ExchangeDeclare(
		cfg.Name,
		cfg.Type,
		cfg.Durable,
		cfg.AutoDelete,
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareQueue declara una cola
func (c *Connection) DeclareQueue(cfg QueueConfig) (amqp.Queue, error) {
	return c.channel.QueueDeclare(
		cfg.Name,
		cfg.Durable,
		cfg.AutoDelete,
		cfg.Exclusive,
		false, // no-wait
		cfg.Args,
	)
}

// BindQueue vincula una cola a un exchange con una routing key
func (c *Connection) BindQueue(queueName, routingKey, exchangeName string) error {
	return c.channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false, // no-wait
		nil,   // arguments
	)
}

// SetPrefetchCount establece el prefetch count
func (c *Connection) SetPrefetchCount(count int) error {
	return c.channel.Qos(
		count, // prefetch count
		0,     // prefetch size
		false, // global
	)
}

// HealthCheck verifica si la conexión está activa
// Crea un canal temporal para evitar race conditions cuando se llama concurrentemente
func (c *Connection) HealthCheck() error {
	if c.IsClosed() {
		return fmt.Errorf("connection is closed")
	}

	// Crear un canal temporal para este health check
	// Esto evita race conditions cuando múltiples goroutines llaman HealthCheck concurrentemente
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer func() { _ = ch.Close() }() //nolint:errcheck // Health check cleanup

	// Intentar declarar un exchange temporal para verificar conectividad
	tempExchange := fmt.Sprintf("health_check_%d", time.Now().UnixNano())
	err = ch.ExchangeDeclare(
		tempExchange,
		"fanout",
		false, // durable
		true,  // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	// Eliminar el exchange temporal
	return ch.ExchangeDelete(tempExchange, false, false)
}
