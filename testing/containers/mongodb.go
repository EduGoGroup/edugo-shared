package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBContainer envuelve el container de MongoDB
// MongoDBContainer envuelve el container de MongoDB de testcontainers.
// Proporciona acceso al cliente MongoDB y métodos de utilidad para
// limpiar colecciones entre tests.
type MongoDBContainer struct {
	container *mongodb.MongoDBContainer
	client    *mongo.Client
	config    *MongoConfig
}

// createMongoDB crea y configura un container de MongoDB
func createMongoDB(ctx context.Context, cfg *MongoConfig) (*MongoDBContainer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("MongoConfig no puede ser nil")
	}

	// Opciones del container
	opts := []testcontainers.ContainerCustomizer{
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections").
				WithStartupTimeout(60 * time.Second),
		),
	}

	// Agregar autenticación si está configurada
	if cfg.Username != "" && cfg.Password != "" {
		opts = append(opts,
			mongodb.WithUsername(cfg.Username),
			mongodb.WithPassword(cfg.Password),
		)
	}

	// Crear container
	container, err := mongodb.Run(ctx, cfg.Image, opts...)
	if err != nil {
		return nil, fmt.Errorf("error creando container MongoDB: %w", err)
	}

	// Obtener connection string
	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("error obteniendo connection string: %w", err)
	}

	// Conectar al cliente
	clientOpts := options.Client().
		ApplyURI(connStr).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("error conectando a MongoDB: %w", err)
	}

	// Verificar conexión
	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		container.Terminate(ctx)
		return nil, fmt.Errorf("error haciendo ping a MongoDB: %w", err)
	}

	return &MongoDBContainer{
		container: container,
		client:    client,
		config:    cfg,
	}, nil
}

// ConnectionString retorna el connection string del container
func (mc *MongoDBContainer) ConnectionString(ctx context.Context) (string, error) {
	return mc.container.ConnectionString(ctx)
}

// Client retorna el cliente de MongoDB
func (mc *MongoDBContainer) Client() *mongo.Client {
	return mc.client
}

// Database retorna una referencia a la base de datos configurada
func (mc *MongoDBContainer) Database() *mongo.Database {
	return mc.client.Database(mc.config.Database)
}

// DropAllCollections elimina todas las colecciones de la base de datos
func (mc *MongoDBContainer) DropAllCollections(ctx context.Context) error {
	db := mc.Database()

	// Listar todas las colecciones
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("error listando colecciones: %w", err)
	}

	// Eliminar cada colección
	for _, collName := range collections {
		if err := db.Collection(collName).Drop(ctx); err != nil {
			return fmt.Errorf("error eliminando colección %s: %w", collName, err)
		}
	}

	return nil
}

// DropCollections elimina colecciones específicas
func (mc *MongoDBContainer) DropCollections(ctx context.Context, collections ...string) error {
	db := mc.Database()

	for _, collName := range collections {
		if err := db.Collection(collName).Drop(ctx); err != nil {
			return fmt.Errorf("error eliminando colección %s: %w", collName, err)
		}
	}

	return nil
}

// Terminate termina el container y cierra las conexiones
func (mc *MongoDBContainer) Terminate(ctx context.Context) error {
	if mc.client != nil {
		if err := mc.client.Disconnect(ctx); err != nil {
			// Log error pero continuar con termination
		}
	}
	if mc.container != nil {
		return mc.container.Terminate(ctx)
	}
	return nil
}
