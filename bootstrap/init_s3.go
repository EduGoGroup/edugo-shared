package bootstrap

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-shared/logger"
)

// initS3 inicializa el cliente de S3.
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
func initS3(
	ctx context.Context,
	config interface{},
	factories *Factories,
	resources *Resources,
	lifecycleManager interface{},
	opts *BootstrapOptions,
) error {
	if factories.S3 == nil {
		return fmt.Errorf("s3 factory not provided")
	}

	// Extraer configuración de S3
	s3Config, err := extractS3Config(config)
	if err != nil {
		return fmt.Errorf("failed to extract S3 config: %w", err)
	}

	// Log inicio
	if resources.Logger != nil {
		resources.Logger.Info("Initializing S3 client...")
	}

	// Crear cliente
	client, err := factories.S3.CreateClient(ctx, s3Config)
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	// Crear StorageClient (implementación simple por ahora)
	resources.StorageClient = &defaultStorageClient{
		client:        client,
		presignClient: factories.S3.CreatePresignClient(client),
		bucket:        s3Config.Bucket,
	}

	// Log éxito
	if resources.Logger != nil {
		resources.Logger.With(
			logger.FieldBucket, s3Config.Bucket,
			logger.FieldRegion, s3Config.Region,
		).Info("S3 client initialized")
	}

	return nil
}
