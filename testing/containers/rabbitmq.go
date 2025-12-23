package containers

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RabbitMQContainer envuelve el container de RabbitMQ
// RabbitMQContainer envuelve el container de RabbitMQ de testcontainers.
// Proporciona acceso a la conexión AMQP y métodos para crear canales
// y limpiar colas entre tests.
type RabbitMQContainer struct {
	container  *rabbitmq.RabbitMQContainer
	connection *amqp.Connection
	config     *RabbitConfig
}

// createRabbitMQ crea y configura un container de RabbitMQ
func createRabbitMQ(ctx context.Context, cfg *RabbitConfig) (*RabbitMQContainer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("RabbitConfig no puede ser nil")
	}

	// Crear container con configuración
	container, err := rabbitmq.Run(ctx,
		cfg.Image,
		rabbitmq.WithAdminUsername(cfg.Username),
		rabbitmq.WithAdminPassword(cfg.Password),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5672/tcp").
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error creando container RabbitMQ: %w", err)
	}

	// Obtener connection string
	connStr, err := container.AmqpURL(ctx)
	if err != nil {
		_ = container.Terminate(ctx) //nolint:errcheck // Cleanup en error, el error principal es el de AMQP URL
		return nil, fmt.Errorf("error obteniendo AMQP URL: %w", err)
	}

	// Conectar con retry
	var conn *amqp.Connection
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(connStr)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		_ = container.Terminate(ctx) //nolint:errcheck // Cleanup en error, el error principal es el de conexión
		return nil, fmt.Errorf("error conectando a RabbitMQ después de %d intentos: %w", maxRetries, err)
	}

	return &RabbitMQContainer{
		container:  container,
		connection: conn,
		config:     cfg,
	}, nil
}

// ConnectionString retorna el AMQP URL del container
func (rc *RabbitMQContainer) ConnectionString(ctx context.Context) (string, error) {
	return rc.container.AmqpURL(ctx)
}

// Connection retorna la conexión AMQP
func (rc *RabbitMQContainer) Connection() *amqp.Connection {
	return rc.connection
}

// Channel crea y retorna un nuevo canal AMQP
func (rc *RabbitMQContainer) Channel() (*amqp.Channel, error) {
	return rc.connection.Channel()
}

// PurgeAll elimina todas las colas del vhost por defecto
// Nota: Esto NO elimina exchanges, solo colas
func (rc *RabbitMQContainer) PurgeAll(ctx context.Context) error {
	ch, err := rc.Channel()
	if err != nil {
		return fmt.Errorf("error creando canal: %w", err)
	}
	defer func() { _ = ch.Close() }() //nolint:errcheck // Close en defer es best-effort

	// Nota: La API de AMQP no proporciona un método directo para listar todas las colas
	// En un entorno de testing, típicamente se conocen los nombres de las colas
	// o se crean colas temporales con nombres únicos
	// Por ahora, este método es un placeholder

	return nil
}

// PurgeQueue elimina todos los mensajes de una cola específica
func (rc *RabbitMQContainer) PurgeQueue(queueName string) error {
	ch, err := rc.Channel()
	if err != nil {
		return fmt.Errorf("error creando canal: %w", err)
	}
	defer func() { _ = ch.Close() }() //nolint:errcheck // Close en defer es best-effort

	_, err = ch.QueuePurge(queueName, false)
	if err != nil {
		return fmt.Errorf("error purgando cola %s: %w", queueName, err)
	}

	return nil
}

// DeleteQueue elimina una cola específica
func (rc *RabbitMQContainer) DeleteQueue(queueName string) error {
	ch, err := rc.Channel()
	if err != nil {
		return fmt.Errorf("error creando canal: %w", err)
	}
	defer func() { _ = ch.Close() }() //nolint:errcheck // Close en defer es best-effort

	_, err = ch.QueueDelete(queueName, false, false, false)
	if err != nil {
		return fmt.Errorf("error eliminando cola %s: %w", queueName, err)
	}

	return nil
}

// Terminate termina el container y cierra las conexiones
func (rc *RabbitMQContainer) Terminate(ctx context.Context) error {
	if rc.connection != nil {
		// Ignorar error de close, el container será terminado de todos modos
		_ = rc.connection.Close() //nolint:errcheck // Cleanup, container será terminado
	}
	if rc.container != nil {
		return rc.container.Terminate(ctx)
	}
	return nil
}
