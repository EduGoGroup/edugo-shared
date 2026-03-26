// Package messaging provides RabbitMQ messaging functionality including
// publishers, consumers, and connection management for the EduGo shared library.
package rabbit

import "time"

const (
	// DefaultPrefetchCount is the default number of messages to prefetch
	DefaultPrefetchCount = 5
)

// ReconnectConfig configures automatic reconnection behavior.
type ReconnectConfig struct {
	Enabled      bool          // Whether to automatically reconnect on disconnect
	InitialDelay time.Duration // Initial delay before first reconnect attempt (default: 1s)
	MaxDelay     time.Duration // Maximum delay between reconnect attempts (default: 30s)
	MaxAttempts  int           // Maximum reconnect attempts, 0 = unlimited (default: 0)
}

// DefaultReconnectConfig returns sensible defaults for reconnection.
func DefaultReconnectConfig() ReconnectConfig {
	return ReconnectConfig{
		Enabled:      true,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		MaxAttempts:  0,
	}
}

// Config contiene la configuración para RabbitMQ
type Config struct {
	// URL de conexión a RabbitMQ
	// Formato: amqp://user:password@host:port/vhost
	URL string

	// Exchange configuración del exchange
	Exchange ExchangeConfig

	// Queue configuración de la cola
	Queue QueueConfig

	// Consumer configuración del consumidor
	Consumer ConsumerConfig

	// PrefetchCount número de mensajes a pre-cargar
	PrefetchCount int
}

// ExchangeConfig configuración del exchange
type ExchangeConfig struct {
	Name       string // Nombre del exchange
	Type       string // Tipo: direct, topic, fanout, headers
	Durable    bool   // Persistent across restarts
	AutoDelete bool   // Auto-eliminar cuando no hay bindings
}

// QueueConfig configuración de la cola
type QueueConfig struct {
	Args       map[string]interface{} // Additional arguments (priority, TTL, etc.)
	Name       string                 // Nombre de la cola
	Durable    bool                   // Persistent across restarts
	AutoDelete bool                   // Auto-eliminar cuando no hay consumidores
	Exclusive  bool                   // Exclusiva para esta conexión
}

// ConsumerConfig configuración del consumidor
type ConsumerConfig struct {
	Name          string    // Nombre del consumidor
	AutoAck       bool      // Auto-acknowledge
	Exclusive     bool      // Exclusivo
	NoLocal       bool      // No recibir mensajes publicados en la misma conexión
	PrefetchCount int       // Número de mensajes a prefetch
	DLQ           DLQConfig // Dead Letter Queue configuration
}

// DefaultConfig retorna una configuración con valores por defecto
func DefaultConfig() Config {
	return Config{ //nolint:gosec // G101: Default local dev URL with well-known guest credentials, not a hardcoded secret
		URL: "amqp://guest:guest@localhost:5672/",
		Exchange: ExchangeConfig{
			Name:       "default_exchange",
			Type:       "topic",
			Durable:    true,
			AutoDelete: false,
		},
		Queue: QueueConfig{
			Name:       "default_queue",
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			Args:       make(map[string]interface{}),
		},
		Consumer: ConsumerConfig{
			Name:      "default_consumer",
			AutoAck:   false,
			Exclusive: false,
			NoLocal:   false,
		},
		PrefetchCount: DefaultPrefetchCount,
	}
}
