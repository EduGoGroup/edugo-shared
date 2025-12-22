package bootstrap

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// initMongoDB inicializa la conexión a MongoDB.
//
// Parámetros:
//   - ctx: Contexto para cancelación
//   - config: Configuración de la aplicación
//   - factories: Fábricas disponibles
//   - resources: Recursos a inicializar
//   - lifecycleManager: Manager de lifecycle para cleanup
//   - opts: Opciones de bootstrap
//
// Retorna error si el recurso es requerido y falla la inicialización.
func initMongoDB(
	ctx context.Context,
	config interface{},
	factories *Factories,
	resources *Resources,
	lifecycleManager interface{},
	opts *BootstrapOptions,
) error {
	if factories.MongoDB == nil {
		return fmt.Errorf("mongodb factory not provided")
	}

	// Extraer configuración de MongoDB
	mongoConfig, err := extractMongoDBConfig(config)
	if err != nil {
		return fmt.Errorf("failed to extract MongoDB config: %w", err)
	}

	// Log inicio
	if resources.Logger != nil {
		resources.Logger.Info("Initializing MongoDB connection...")
	}

	// Crear conexión
	client, err := factories.MongoDB.CreateConnection(ctx, mongoConfig)
	if err != nil {
		return fmt.Errorf("failed to create MongoDB connection: %w", err)
	}

	resources.MongoDB = client
	resources.MongoDatabase = factories.MongoDB.GetDatabase(client, mongoConfig.Database)

	// Registrar cleanup en lifecycle manager si está disponible
	if lifecycleManager != nil {
		registerMongoDBCleanup(lifecycleManager, factories.MongoDB, client, resources.Logger)
	}

	// Log éxito
	if resources.Logger != nil {
		resources.Logger.WithFields(logrus.Fields{
			"database": mongoConfig.Database,
		}).Info("MongoDB connection established")
	}

	return nil
}
