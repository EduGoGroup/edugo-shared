// Package rabbit provee funcionalidad de mensajeria RabbitMQ incluyendo
// publishers, consumers, y gestion de conexiones para la libreria compartida EduGo.
package rabbit

import "time"

const (
	// DefaultPrefetchCount es el numero de mensajes a prefetch por defecto
	DefaultPrefetchCount = 5
)

// ReconnectConfig configura el comportamiento de reconexion automatica.
type ReconnectConfig struct {
	Enabled      bool          // Reconectar automaticamente al desconectarse
	InitialDelay time.Duration // Delay inicial antes del primer intento de reconexion (default: 1s)
	MaxDelay     time.Duration // Delay maximo entre intentos de reconexion (default: 30s)
	MaxAttempts  int           // Intentos maximos de reconexion, 0 = ilimitado (default: 0)
}

// DefaultReconnectConfig retorna valores por defecto para la reconexion.
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
	Durable    bool   // Persistente entre reinicios
	AutoDelete bool   // Auto-eliminar cuando no hay bindings
}

// QueueConfig configuracion de la cola
type QueueConfig struct {
	Args       map[string]any // Argumentos adicionales (prioridad, TTL, etc.)
	Name       string                 // Nombre de la cola
	Durable    bool                   // Persistente entre reinicios
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
	DLQ           DLQConfig // Configuracion Dead Letter Queue
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
			Args:       make(map[string]any),
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
