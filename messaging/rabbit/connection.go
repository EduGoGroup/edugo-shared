package rabbit

import (
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection encapsula la conexion a RabbitMQ
type Connection struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	url          string
	mu           sync.RWMutex
	reconnectCfg ReconnectConfig
	closeCh      chan struct{}
	closeOnce    sync.Once
	reconnectCh  chan struct{}
	logger       func(msg string, args ...any)
}

// Connect establece una conexion a RabbitMQ
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
		closeCh: make(chan struct{}),
	}, nil
}

// ConnectWithReconnect establishes a connection with automatic reconnection on disconnect.
func ConnectWithReconnect(url string, cfg ReconnectConfig) (*Connection, error) {
	c, err := Connect(url)
	if err != nil {
		return nil, err
	}
	c.reconnectCfg = cfg
	c.reconnectCh = make(chan struct{}, 1)

	if cfg.Enabled {
		go c.reconnectLoop()
	}

	return c, nil
}

// reconnectLoop watches for connection close notifications and triggers reconnection.
func (c *Connection) reconnectLoop() {
	for {
		// Register for close notification
		notifyClose := make(chan *amqp.Error, 1)
		c.mu.RLock()
		if c.conn != nil && !c.conn.IsClosed() {
			c.conn.NotifyClose(notifyClose)
		}
		c.mu.RUnlock()

		select {
		case <-c.closeCh:
			return
		case amqpErr, ok := <-notifyClose:
			if !ok {
				// Channel closed without error (graceful close)
				return
			}
			// Log reconnection attempt
			if c.logger != nil {
				c.logger("RabbitMQ connection lost, reconnecting...", "error", amqpErr)
			}
			c.doReconnect()
		}
	}
}

// doReconnect attempts to re-establish the connection with exponential backoff.
func (c *Connection) doReconnect() {
	delay := c.reconnectCfg.InitialDelay
	if delay == 0 {
		delay = time.Second
	}
	maxDelay := c.reconnectCfg.MaxDelay
	if maxDelay == 0 {
		maxDelay = 30 * time.Second
	}

	for attempt := 1; c.reconnectCfg.MaxAttempts == 0 || attempt <= c.reconnectCfg.MaxAttempts; attempt++ {
		select {
		case <-c.closeCh:
			return
		case <-time.After(delay):
		}

		if c.tryReconnect(attempt) {
			return
		}

		delay = nextBackoff(delay, maxDelay)
	}

	if c.logger != nil {
		c.logger("RabbitMQ reconnect failed after max attempts", "max_attempts", c.reconnectCfg.MaxAttempts)
	}
}

// tryReconnect attempts a single reconnection. Returns true on success.
func (c *Connection) tryReconnect(attempt int) bool {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		if c.logger != nil {
			c.logger("RabbitMQ reconnect failed", "attempt", attempt, "error", err)
		}
		return false
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close() //nolint:errcheck // cleanup on channel error
		if c.logger != nil {
			c.logger("RabbitMQ channel creation failed after reconnect", "attempt", attempt, "error", err)
		}
		return false
	}

	c.mu.Lock()
	c.conn = conn
	c.channel = ch
	c.mu.Unlock()

	if c.logger != nil {
		c.logger("RabbitMQ reconnected successfully", "attempt", attempt)
	}

	// Notify listeners
	select {
	case c.reconnectCh <- struct{}{}:
	default:
	}

	return true
}

// nextBackoff doubles the delay up to maxDelay.
func nextBackoff(delay, maxDelay time.Duration) time.Duration {
	delay *= 2
	if delay > maxDelay {
		return maxDelay
	}
	return delay
}

// NotifyReconnect returns a channel that receives a signal when the connection is re-established.
func (c *Connection) NotifyReconnect() <-chan struct{} {
	if c.reconnectCh == nil {
		ch := make(chan struct{})
		return ch
	}
	return c.reconnectCh
}

// SetLogger sets an optional logger for reconnection events.
func (c *Connection) SetLogger(logger func(msg string, args ...any)) {
	c.logger = logger
}

// GetChannel retorna el canal de RabbitMQ
func (c *Connection) GetChannel() *amqp.Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channel
}

// GetConnection retorna la conexion de RabbitMQ
func (c *Connection) GetConnection() *amqp.Connection {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

// Close cierra la conexion y el canal
func (c *Connection) Close() error {
	if c.closeCh != nil {
		c.closeOnce.Do(func() {
			close(c.closeCh)
		})
	}

	c.mu.RLock()
	ch := c.channel
	conn := c.conn
	c.mu.RUnlock()

	var errs []error
	if ch != nil {
		if err := ch.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close channel: %w", err))
		}
	}
	if conn != nil {
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection: %w", err))
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// IsClosed verifica si la conexion esta cerrada
func (c *Connection) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
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

// HealthCheck verifica si la conexion esta activa
// Crea un canal temporal para evitar race conditions cuando se llama concurrentemente
func (c *Connection) HealthCheck() error {
	if c.IsClosed() {
		return fmt.Errorf("connection is closed")
	}

	// Crear un canal temporal para este health check
	// Esto evita race conditions cuando multiples goroutines llaman HealthCheck concurrentemente
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	ch, err := conn.Channel()
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
