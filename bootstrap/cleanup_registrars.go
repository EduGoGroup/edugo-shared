package bootstrap

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

// registerPostgreSQLCleanup registra la función de cleanup de PostgreSQL en el lifecycle manager.
//
// Parámetros:
//   - lifecycleManager: Manager de lifecycle
//   - factory: Factory de PostgreSQL
//   - db: Instancia de la base de datos
//   - logger: Logger para mensajes de cleanup
func registerPostgreSQLCleanup(lifecycleManager interface{}, factory PostgreSQLFactory, db interface{}, logger logger.Logger) {
	registrar, ok := lifecycleManager.(interface {
		RegisterSimple(name string, cleanup func() error)
	})
	if !ok || factory == nil || db == nil {
		return
	}

	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return
	}

	registrar.RegisterSimple("postgresql", func() error {
		if logger != nil {
			logger.Info("Closing PostgreSQL connection via lifecycle manager")
		}
		return factory.Close(gormDB)
	})
}

// registerMongoDBCleanup registra la función de cleanup de MongoDB en el lifecycle manager.
//
// Parámetros:
//   - lifecycleManager: Manager de lifecycle
//   - factory: Factory de MongoDB
//   - client: Cliente de MongoDB
//   - logger: Logger para mensajes de cleanup
func registerMongoDBCleanup(lifecycleManager interface{}, factory MongoDBFactory, client interface{}, logger logger.Logger) {
	registrar, ok := lifecycleManager.(interface {
		RegisterSimple(name string, cleanup func() error)
	})
	if !ok || factory == nil || client == nil {
		return
	}

	mongoClient, ok := client.(*mongo.Client)
	if !ok {
		return
	}

	registrar.RegisterSimple("mongodb", func() error {
		if logger != nil {
			logger.Info("Closing MongoDB client via lifecycle manager")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return factory.Close(ctx, mongoClient)
	})
}

// registerRabbitMQCleanup registra la función de cleanup de RabbitMQ en el lifecycle manager.
//
// Parámetros:
//   - lifecycleManager: Manager de lifecycle
//   - factory: Factory de RabbitMQ
//   - channel: Canal de RabbitMQ
//   - conn: Conexión de RabbitMQ
//   - logger: Logger para mensajes de cleanup
func registerRabbitMQCleanup(lifecycleManager interface{}, factory RabbitMQFactory, channel interface{}, conn interface{}, logger logger.Logger) {
	registrar, ok := lifecycleManager.(interface {
		RegisterSimple(name string, cleanup func() error)
	})
	if !ok || factory == nil || channel == nil || conn == nil {
		return
	}

	amqpChannel, ok := channel.(*amqp.Channel)
	if !ok {
		return
	}
	amqpConn, ok := conn.(*amqp.Connection)
	if !ok {
		return
	}

	registrar.RegisterSimple("rabbitmq", func() error {
		if logger != nil {
			logger.Info("Closing RabbitMQ channel and connection via lifecycle manager")
		}
		return factory.Close(amqpChannel, amqpConn)
	})
}
