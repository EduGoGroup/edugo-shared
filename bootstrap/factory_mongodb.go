package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// =============================================================================
// MONGODB FACTORY IMPLEMENTATION
// =============================================================================

// DefaultMongoDBFactory implementa MongoDBFactory
type DefaultMongoDBFactory struct {
	connectionTimeout time.Duration
}

// NewDefaultMongoDBFactory crea una nueva instancia de DefaultMongoDBFactory
func NewDefaultMongoDBFactory() *DefaultMongoDBFactory {
	return &DefaultMongoDBFactory{
		connectionTimeout: 10 * time.Second,
	}
}

// CreateConnection crea una conexión a MongoDB
func (f *DefaultMongoDBFactory) CreateConnection(ctx context.Context, config MongoDBConfig) (*mongo.Client, error) {
	// Crear context con timeout
	ctx, cancel := context.WithTimeout(ctx, f.connectionTimeout)
	defer cancel()

	// Configurar opciones del cliente
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Minute).
		SetServerSelectionTimeout(5 * time.Second).
		SetConnectTimeout(10 * time.Second)

	// Crear cliente
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verificar conexión
	if err := f.Ping(ctx, client); err != nil {
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to ping MongoDB: %w", err),
				fmt.Errorf("failed to disconnect: %w", disconnectErr),
			)
		}
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// GetDatabase obtiene una base de datos específica
func (f *DefaultMongoDBFactory) GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}

// Ping verifica la conectividad con MongoDB
func (f *DefaultMongoDBFactory) Ping(ctx context.Context, client *mongo.Client) error {
	// Crear context con timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Intentar ping con preferencia de lectura primaria
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	return nil
}

// Close cierra la conexión
func (f *DefaultMongoDBFactory) Close(ctx context.Context, client *mongo.Client) error {
	// Crear context con timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	return nil
}

// Verificar que DefaultMongoDBFactory implementa MongoDBFactory
var _ MongoDBFactory = (*DefaultMongoDBFactory)(nil)
