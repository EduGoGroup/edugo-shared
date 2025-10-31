package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// DefaultHealthCheckTimeout is the default timeout for health checks
	DefaultHealthCheckTimeout = 5 * time.Second
	// DefaultDisconnectTimeout is the default timeout for disconnection
	DefaultDisconnectTimeout = 10 * time.Second
)

// Connect establece una conexión a MongoDB y retorna el cliente
func Connect(cfg Config) (*mongo.Client, error) {
	// Crear opciones de cliente
	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetConnectTimeout(cfg.Timeout).
		SetServerSelectionTimeout(cfg.Timeout)

	// Conectar con timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// Verificar conectividad
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx) // Ignore error on cleanup
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return client, nil
}

// GetDatabase retorna una instancia de base de datos
func GetDatabase(client *mongo.Client, databaseName string) *mongo.Database {
	return client.Database(databaseName)
}

// HealthCheck verifica si la conexión a MongoDB está activa
func HealthCheck(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultHealthCheckTimeout)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("mongodb health check failed: %w", err)
	}

	return nil
}

// Close cierra la conexión a MongoDB
func Close(client *mongo.Client) error {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultDisconnectTimeout)
		defer cancel()

		return client.Disconnect(ctx)
	}
	return nil
}
