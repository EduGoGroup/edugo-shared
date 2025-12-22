package bootstrap

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Bootstrap inicializa todos los recursos de infraestructura de la aplicación.
//
// Esta función coordina la inicialización ordenada de todos los recursos necesarios
// para la aplicación, incluyendo logger, bases de datos, colas de mensajes y almacenamiento.
//
// Parámetros:
//   - ctx: Contexto para cancelación y timeouts
//   - config: Struct de configuración con campos para cada recurso
//   - factories: Fábricas para crear recursos
//   - lifecycleManager: Manager de lifecycle para cleanup ordenado
//   - options: Opciones adicionales de configuración
//
// Retorna los recursos inicializados o error si falla algún recurso requerido.
func Bootstrap(
	ctx context.Context,
	config interface{},
	factories *Factories,
	lifecycleManager interface{},
	options ...BootstrapOption,
) (*Resources, error) {
	// Aplicar opciones
	opts := DefaultBootstrapOptions()
	ApplyOptions(opts, options...)

	// Usar factories mock si están configuradas
	if opts.MockFactories != nil {
		factories = mergeFactories(factories, opts.MockFactories)
	}

	// Validar factories requeridas
	if err := factories.Validate(opts.RequiredResources); err != nil {
		return nil, fmt.Errorf("factory validation failed: %w", err)
	}

	// Inicializar recursos
	resources := &Resources{}

	// Inicializar Logger (siempre primero)
	if err := initLogger(ctx, config, factories, resources, opts); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Log inicio del bootstrap
	if resources.Logger != nil {
		resources.Logger.Info("Starting application bootstrap...")
		resources.Logger.WithFields(logrus.Fields{
			"required_resources": opts.RequiredResources,
			"optional_resources": opts.OptionalResources,
		}).Debug("Bootstrap configuration")
	}

	// Inicializar PostgreSQL
	if err := initPostgreSQL(ctx, config, factories, resources, lifecycleManager, opts); err != nil {
		if isRequired("postgresql", opts) {
			return nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
		}
		logWarning(resources.Logger, "PostgreSQL initialization skipped", err)
	}

	// Inicializar MongoDB
	if err := initMongoDB(ctx, config, factories, resources, lifecycleManager, opts); err != nil {
		if isRequired("mongodb", opts) {
			return nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
		}
		logWarning(resources.Logger, "MongoDB initialization skipped", err)
	}

	// Inicializar RabbitMQ
	if err := initRabbitMQ(ctx, config, factories, resources, lifecycleManager, opts); err != nil {
		if isRequired("rabbitmq", opts) {
			return nil, fmt.Errorf("failed to initialize RabbitMQ: %w", err)
		}
		logWarning(resources.Logger, "RabbitMQ initialization skipped", err)
	}

	// Inicializar S3
	if err := initS3(ctx, config, factories, resources, lifecycleManager, opts); err != nil {
		if isRequired("s3", opts) {
			return nil, fmt.Errorf("failed to initialize S3: %w", err)
		}
		logWarning(resources.Logger, "S3 initialization skipped", err)
	}

	// Health checks (si no están deshabilitados)
	if !opts.SkipHealthCheck {
		if err := performHealthChecks(ctx, resources, opts); err != nil {
			return nil, fmt.Errorf("health checks failed: %w", err)
		}
	}

	// Log finalización exitosa
	if resources.Logger != nil {
		resources.Logger.Info("Application bootstrap completed successfully")
	}

	return resources, nil
}
