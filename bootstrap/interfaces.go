package bootstrap

import (
	"context"
	"database/sql"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// =============================================================================
// CONFIGURATION STRUCTS
// =============================================================================

// PostgreSQLConfig define la configuración para PostgreSQL
type PostgreSQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// MongoDBConfig define la configuración para MongoDB
type MongoDBConfig struct {
	URI      string
	Database string
}

// RabbitMQConfig define la configuración para RabbitMQ
type RabbitMQConfig struct {
	URL string
}

// S3Config define la configuración para AWS S3
type S3Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

// =============================================================================
// FACTORY INTERFACES
// =============================================================================

// LoggerFactory crea y configura instancias de logger
type LoggerFactory interface {
	// CreateLogger crea un logger con la configuración especificada
	CreateLogger(ctx context.Context, env string, version string) (*logrus.Logger, error)
}

// PostgreSQLFactory crea y gestiona conexiones a PostgreSQL
type PostgreSQLFactory interface {
	// CreateConnection crea una conexión GORM a PostgreSQL
	CreateConnection(ctx context.Context, config PostgreSQLConfig) (*gorm.DB, error)

	// CreateRawConnection crea una conexión SQL nativa (para casos específicos)
	CreateRawConnection(ctx context.Context, config PostgreSQLConfig) (*sql.DB, error)

	// Ping verifica la conectividad con PostgreSQL
	Ping(ctx context.Context, db *gorm.DB) error

	// Close cierra la conexión
	Close(db *gorm.DB) error
}

// MongoDBFactory crea y gestiona conexiones a MongoDB
type MongoDBFactory interface {
	// CreateConnection crea una conexión a MongoDB
	CreateConnection(ctx context.Context, config MongoDBConfig) (*mongo.Client, error)

	// GetDatabase obtiene una base de datos específica
	GetDatabase(client *mongo.Client, dbName string) *mongo.Database

	// Ping verifica la conectividad con MongoDB
	Ping(ctx context.Context, client *mongo.Client) error

	// Close cierra la conexión
	Close(ctx context.Context, client *mongo.Client) error
}

// RabbitMQFactory crea y gestiona conexiones a RabbitMQ
type RabbitMQFactory interface {
	// CreateConnection crea una conexión a RabbitMQ
	CreateConnection(ctx context.Context, config RabbitMQConfig) (*amqp.Connection, error)

	// CreateChannel crea un canal de comunicación
	CreateChannel(conn *amqp.Connection) (*amqp.Channel, error)

	// DeclareQueue declara una cola
	DeclareQueue(channel *amqp.Channel, queueName string) (amqp.Queue, error)

	// Close cierra el canal y la conexión
	Close(channel *amqp.Channel, conn *amqp.Connection) error
}

// S3Factory crea y gestiona clientes de AWS S3
type S3Factory interface {
	// CreateClient crea un cliente de S3
	CreateClient(ctx context.Context, config S3Config) (*s3.Client, error)

	// CreatePresignClient crea un cliente para URLs pre-firmadas
	CreatePresignClient(client *s3.Client) *s3.PresignClient

	// ValidateBucket verifica que el bucket existe y es accesible
	ValidateBucket(ctx context.Context, client *s3.Client, bucket string) error
}

// =============================================================================
// RESOURCE INTERFACES
// =============================================================================

// MessagePublisher define la interfaz para publicar mensajes
type MessagePublisher interface {
	// Publish publica un mensaje en una cola
	Publish(ctx context.Context, queueName string, body []byte) error

	// PublishWithPriority publica un mensaje con prioridad
	PublishWithPriority(ctx context.Context, queueName string, body []byte, priority uint8) error

	// Close cierra el publicador
	Close() error
}

// StorageClient define la interfaz para operaciones de almacenamiento
type StorageClient interface {
	// Upload sube un archivo al storage
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)

	// Download descarga un archivo del storage
	Download(ctx context.Context, key string) ([]byte, error)

	// Delete elimina un archivo del storage
	Delete(ctx context.Context, key string) error

	// GetPresignedURL genera una URL pre-firmada para acceso temporal
	GetPresignedURL(ctx context.Context, key string, expirationMinutes int) (string, error)

	// Exists verifica si un archivo existe
	Exists(ctx context.Context, key string) (bool, error)
}

// DatabaseClient define operaciones básicas de base de datos
type DatabaseClient interface {
	// Ping verifica la conectividad
	Ping(ctx context.Context) error

	// Close cierra la conexión
	Close(ctx context.Context) error

	// GetStats obtiene estadísticas de la conexión
	GetStats(ctx context.Context) (map[string]interface{}, error)
}

// =============================================================================
// HEALTH CHECK INTERFACES
// =============================================================================

// HealthChecker define la interfaz para health checks de recursos
type HealthChecker interface {
	// Check verifica el estado de salud de un recurso
	Check(ctx context.Context) error

	// GetResourceName retorna el nombre del recurso
	GetResourceName() string

	// IsRequired indica si el recurso es obligatorio
	IsRequired() bool
}

// =============================================================================
// FACTORY COLLECTION
// =============================================================================

// Factories agrupa todas las factories disponibles
type Factories struct {
	Logger     LoggerFactory
	PostgreSQL PostgreSQLFactory
	MongoDB    MongoDBFactory
	RabbitMQ   RabbitMQFactory
	S3         S3Factory
}

// Validate verifica que todas las factories requeridas estén presentes
func (f *Factories) Validate(requiredResources []string) error {
	for _, resource := range requiredResources {
		switch resource {
		case "logger":
			if f.Logger == nil {
				return ErrMissingFactory{Resource: resource}
			}
		case "postgresql":
			if f.PostgreSQL == nil {
				return ErrMissingFactory{Resource: resource}
			}
		case "mongodb":
			if f.MongoDB == nil {
				return ErrMissingFactory{Resource: resource}
			}
		case "rabbitmq":
			if f.RabbitMQ == nil {
				return ErrMissingFactory{Resource: resource}
			}
		case "s3":
			if f.S3 == nil {
				return ErrMissingFactory{Resource: resource}
			}
		}
	}
	return nil
}

// =============================================================================
// ERRORS
// =============================================================================

// ErrMissingFactory se lanza cuando falta una factory requerida
type ErrMissingFactory struct {
	Resource string
}

func (e ErrMissingFactory) Error() string {
	return "missing required factory: " + e.Resource
}
