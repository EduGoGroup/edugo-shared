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

// ConnectWithReconnect establece una conexion con reconexion automatica al desconectarse.
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

// reconnectLoop observa notificaciones de cierre de conexion y dispara la reconexion.
func (c *Connection) reconnectLoop() {
	for {
		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil || conn.IsClosed() {
			// La conexion no esta disponible, intentar reconectar
			if !c.doReconnect() {
				// MaxAttempts agotados, salir del loop
				c.log("reconnect loop saliendo: intentos maximos agotados")
				return
			}
			continue
		}

		// Registrar notificacion de cierre
		notifyClose := make(chan *amqp.Error, 1)
		conn.NotifyClose(notifyClose)

		select {
		case <-c.closeCh:
			return
		case amqpErr, ok := <-notifyClose:
			if !ok {
				// Canal cerrado sin error (cierre graceful)
				return
			}
			c.log("conexion RabbitMQ perdida, reconectando...", "error", amqpErr)
			if !c.doReconnect() {
				// MaxAttempts agotados, salir del loop
				c.log("reconnect loop saliendo: intentos maximos agotados")
				return
			}
		}
	}
}

// doReconnect intenta restablecer la conexion con backoff exponencial.
// Retorna true si la reconexion fue exitosa, false si se agotaron los intentos o se cerro.
func (c *Connection) doReconnect() bool {
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
			return false
		case <-time.After(delay):
		}

		if c.tryReconnect(attempt) {
			return true
		}

		delay = nextBackoff(delay, maxDelay)
	}

	c.log("reconexion RabbitMQ fallida tras intentos maximos", "max_attempts", c.reconnectCfg.MaxAttempts)
	return false
}

// tryReconnect intenta una reconexion individual. Retorna true si fue exitosa.
func (c *Connection) tryReconnect(attempt int) bool {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		c.log("reconexion RabbitMQ fallida", "attempt", attempt, "error", err)
		return false
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close() //nolint:errcheck // cleanup en error de canal
		c.log("creacion de canal fallida tras reconexion", "attempt", attempt, "error", err)
		return false
	}

	// Verificar closeCh antes de asignar para evitar race con Close()
	select {
	case <-c.closeCh:
		// La conexion fue cerrada mientras reconectabamos, limpiar recursos nuevos
		_ = ch.Close()   //nolint:errcheck // cleanup de recursos huerfanos
		_ = conn.Close() //nolint:errcheck // cleanup de recursos huerfanos
		return false
	default:
	}

	c.mu.Lock()
	c.conn = conn
	c.channel = ch
	c.mu.Unlock()

	c.log("RabbitMQ reconectado exitosamente", "attempt", attempt)

	// Notificar a los listeners
	select {
	case c.reconnectCh <- struct{}{}:
	default:
	}

	return true
}

// log escribe un mensaje de log si el logger esta configurado.
// Accede a c.logger bajo RLock para evitar data race con SetLogger.
func (c *Connection) log(msg string, args ...any) {
	c.mu.RLock()
	lgr := c.logger
	c.mu.RUnlock()

	if lgr != nil {
		lgr(msg, args...)
	}
}

// nextBackoff duplica el delay hasta maxDelay.
func nextBackoff(delay, maxDelay time.Duration) time.Duration {
	delay *= 2
	if delay > maxDelay {
		return maxDelay
	}
	return delay
}

// NotifyReconnect retorna un canal que recibe una senal cuando la conexion se restablece.
// Retorna nil si la reconexion no esta habilitada.
func (c *Connection) NotifyReconnect() <-chan struct{} {
	return c.reconnectCh
}

// SetLogger configura un logger opcional para eventos de reconexion.
// Es seguro para uso concurrente.
func (c *Connection) SetLogger(logger func(msg string, args ...any)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
}

// GetChannel retorna el canal de RabbitMQ de forma segura para concurrencia.
func (c *Connection) GetChannel() *amqp.Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channel
}

// GetConnection retorna la conexion de RabbitMQ de forma segura para concurrencia.
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
	ch := c.GetChannel()
	if ch == nil {
		return fmt.Errorf("channel is nil")
	}
	return ch.ExchangeDeclare(
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
	ch := c.GetChannel()
	if ch == nil {
		return amqp.Queue{}, fmt.Errorf("channel is nil")
	}
	return ch.QueueDeclare(
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
	ch := c.GetChannel()
	if ch == nil {
		return fmt.Errorf("channel is nil")
	}
	return ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false, // no-wait
		nil,   // arguments
	)
}

// SetPrefetchCount establece el prefetch count
func (c *Connection) SetPrefetchCount(count int) error {
	ch := c.GetChannel()
	if ch == nil {
		return fmt.Errorf("channel is nil")
	}
	return ch.Qos(
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

	// Obtener conexion bajo RLock
	conn := c.GetConnection()
	if conn == nil {
		return fmt.Errorf("connection is closed")
	}

	// Crear un canal temporal para este health check
	// Esto evita race conditions cuando multiples goroutines llaman HealthCheck concurrentemente
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
