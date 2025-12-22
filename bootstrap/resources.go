package bootstrap

import (
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// =============================================================================
// RESOURCES CONTAINER
// =============================================================================

// Resources contiene todos los recursos inicializados de la aplicación
type Resources struct {
	// Logger es el logger configurado de la aplicación
	Logger logger.Logger

	// PostgreSQL es la conexión a la base de datos PostgreSQL
	PostgreSQL *gorm.DB

	// MongoDB es el cliente de MongoDB
	MongoDB *mongo.Client

	// MongoDatabase es la base de datos específica de MongoDB
	MongoDatabase *mongo.Database

	// MessagePublisher es el publicador de mensajes (RabbitMQ)
	MessagePublisher MessagePublisher

	// StorageClient es el cliente de almacenamiento (S3)
	StorageClient StorageClient
}

// HasLogger verifica si el logger está inicializado
func (r *Resources) HasLogger() bool {
	return r.Logger != nil
}

// HasPostgreSQL verifica si PostgreSQL está inicializado
func (r *Resources) HasPostgreSQL() bool {
	return r.PostgreSQL != nil
}

// HasMongoDB verifica si MongoDB está inicializado
func (r *Resources) HasMongoDB() bool {
	return r.MongoDB != nil && r.MongoDatabase != nil
}

// HasMessagePublisher verifica si el publicador de mensajes está inicializado
func (r *Resources) HasMessagePublisher() bool {
	return r.MessagePublisher != nil
}

// HasStorageClient verifica si el cliente de almacenamiento está inicializado
func (r *Resources) HasStorageClient() bool {
	return r.StorageClient != nil
}
