package containers

import (
	"context"
	"fmt"
)

// MongoDBContainer envuelve el container de MongoDB
type MongoDBContainer struct {
	// container *mongodb.MongoDBContainer
	config *MongoConfig
}

// createMongoDB crea y configura un container de MongoDB
func createMongoDB(ctx context.Context, cfg *MongoConfig) (*MongoDBContainer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("MongoConfig no puede ser nil")
	}

	// TODO: Implementar creación de container MongoDB
	// Por ahora retorna un stub

	return &MongoDBContainer{
		config: cfg,
	}, nil
}

// DropAllCollections elimina todas las colecciones de la base de datos
func (mc *MongoDBContainer) DropAllCollections(ctx context.Context) error {
	// TODO: Implementar
	return nil
}

// Terminate termina el container
func (mc *MongoDBContainer) Terminate(ctx context.Context) error {
	// TODO: Implementar
	return nil
}
