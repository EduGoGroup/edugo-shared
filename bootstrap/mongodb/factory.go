package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-shared/bootstrap"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Factory implementa la creacion de conexiones MongoDB.
type Factory struct {
	connectionTimeout time.Duration
}

// NewFactory crea una nueva Factory de MongoDB.
func NewFactory() *Factory {
	return &Factory{
		connectionTimeout: 10 * time.Second,
	}
}

// CreateConnection crea una conexion a MongoDB con pool configurado.
func (f *Factory) CreateConnection(ctx context.Context, cfg bootstrap.MongoDBConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, f.connectionTimeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Minute).
		SetServerSelectionTimeout(5 * time.Second).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("bootstrap/mongodb: connect: %w", err)
	}

	if err := f.Ping(ctx, client); err != nil {
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			return nil, errors.Join(
				fmt.Errorf("bootstrap/mongodb: ping: %w", err),
				fmt.Errorf("bootstrap/mongodb: disconnect: %w", disconnectErr),
			)
		}
		return nil, fmt.Errorf("bootstrap/mongodb: ping: %w", err)
	}

	return client, nil
}

// GetDatabase obtiene una base de datos especifica.
func (f *Factory) GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}

// Ping verifica la conectividad con MongoDB.
func (f *Factory) Ping(ctx context.Context, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

// Close cierra la conexion a MongoDB.
func (f *Factory) Close(ctx context.Context, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("bootstrap/mongodb: disconnect: %w", err)
	}
	return nil
}
