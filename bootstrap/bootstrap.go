package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// =============================================================================
// BOOTSTRAP FUNCTION
// =============================================================================

// Bootstrap inicializa todos los recursos de la aplicación
// Retorna Resources con los recursos inicializados y error si falla alguno
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

// =============================================================================
// LOGGER INITIALIZATION
// =============================================================================

func initLogger(
	ctx context.Context,
	config interface{},
	factories *Factories,
	resources *Resources,
	opts *BootstrapOptions,
) error {
	if factories.Logger == nil {
		return fmt.Errorf("logger factory is required but not provided")
	}

	// Extraer configuración
	env, version := extractEnvAndVersion(config)

	// Crear logger
	logger, err := factories.Logger.CreateLogger(ctx, env, version)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	resources.Logger = logger
	return nil
}

// =============================================================================
// POSTGRESQL INITIALIZATION
// =============================================================================

func initPostgreSQL(
	ctx context.Context,
	config interface{},
	factories *Factories,
	resources *Resources,
	lifecycleManager interface{},
	opts *BootstrapOptions,
) error {
	if factories.PostgreSQL == nil {
		return fmt.Errorf("postgresql factory not provided")
	}

	// Extraer configuración de PostgreSQL
	pgConfig, err := extractPostgreSQLConfig(config)
	if err != nil {
		return fmt.Errorf("failed to extract PostgreSQL config: %w", err)
	}

	// Log inicio
	if resources.Logger != nil {
		resources.Logger.Info("Initializing PostgreSQL connection...")
	}

	// Crear conexión
	db, err := factories.PostgreSQL.CreateConnection(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL connection: %w", err)
	}

	resources.PostgreSQL = db

	// Registrar cleanup en lifecycle manager si está disponible
	if lifecycleManager != nil {
		registerPostgreSQLCleanup(lifecycleManager, factories.PostgreSQL, db, resources.Logger)
	}

	// Log éxito
	if resources.Logger != nil {
		resources.Logger.WithFields(logrus.Fields{
			"host":     pgConfig.Host,
			"port":     pgConfig.Port,
			"database": pgConfig.Database,
		}).Info("PostgreSQL connection established")
	}

	return nil
}

// =============================================================================
// MONGODB INITIALIZATION
// =============================================================================

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

// =============================================================================
// RABBITMQ INITIALIZATION
// =============================================================================

func initRabbitMQ(
	ctx context.Context,
	config interface{},
	factories *Factories,
	resources *Resources,
	lifecycleManager interface{},
	opts *BootstrapOptions,
) error {
	if factories.RabbitMQ == nil {
		return fmt.Errorf("rabbitmq factory not provided")
	}

	// Extraer configuración de RabbitMQ
	rabbitConfig, err := extractRabbitMQConfig(config)
	if err != nil {
		return fmt.Errorf("failed to extract RabbitMQ config: %w", err)
	}

	// Log inicio
	if resources.Logger != nil {
		resources.Logger.Info("Initializing RabbitMQ connection...")
	}

	// Crear conexión
	conn, err := factories.RabbitMQ.CreateConnection(ctx, rabbitConfig)
	if err != nil {
		return fmt.Errorf("failed to create RabbitMQ connection: %w", err)
	}

	// Crear canal
	channel, err := factories.RabbitMQ.CreateChannel(conn)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create RabbitMQ channel: %w", err)
	}

	// Crear MessagePublisher (implementación simple por ahora)
	resources.MessagePublisher = &defaultMessagePublisher{
		channel: channel,
		factory: factories.RabbitMQ,
	}

	// Registrar cleanup en lifecycle manager si está disponible
	if lifecycleManager != nil {
		registerRabbitMQCleanup(lifecycleManager, factories.RabbitMQ, channel, conn, resources.Logger)
	}

	// Log éxito
	if resources.Logger != nil {
		resources.Logger.Info("RabbitMQ connection established")
	}

	return nil
}

// =============================================================================
// S3 INITIALIZATION
// =============================================================================

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
		client:       client,
		presignClient: factories.S3.CreatePresignClient(client),
		bucket:       s3Config.Bucket,
	}

	// Log éxito
	if resources.Logger != nil {
		resources.Logger.WithFields(logrus.Fields{
			"bucket": s3Config.Bucket,
			"region": s3Config.Region,
		}).Info("S3 client initialized")
	}

	return nil
}

// =============================================================================
// HEALTH CHECKS
// =============================================================================

func performHealthChecks(ctx context.Context, resources *Resources, opts *BootstrapOptions) error {
	if resources.Logger != nil {
		resources.Logger.Info("Performing health checks...")
	}

	// Health check timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// PostgreSQL health check
	if resources.PostgreSQL != nil {
		if err := resources.PostgreSQL.Raw("SELECT 1").Error; err != nil {
			return fmt.Errorf("postgresql health check failed: %w", err)
		}
		if resources.Logger != nil {
			resources.Logger.Debug("PostgreSQL health check passed")
		}
	}

	// MongoDB health check
	if resources.MongoDB != nil {
		if err := resources.MongoDB.Ping(ctx, nil); err != nil {
			return fmt.Errorf("mongodb health check failed: %w", err)
		}
		if resources.Logger != nil {
			resources.Logger.Debug("MongoDB health check passed")
		}
	}

	if resources.Logger != nil {
		resources.Logger.Info("All health checks passed")
	}

	return nil
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func mergeFactories(base *Factories, mocks *MockFactories) *Factories {
	result := &Factories{}
	if base != nil {
		*result = *base
	}
	if mocks != nil {
		if mocks.Logger != nil {
			result.Logger = mocks.Logger
		}
		if mocks.PostgreSQL != nil {
			result.PostgreSQL = mocks.PostgreSQL
		}
		if mocks.MongoDB != nil {
			result.MongoDB = mocks.MongoDB
		}
		if mocks.RabbitMQ != nil {
			result.RabbitMQ = mocks.RabbitMQ
		}
		if mocks.S3 != nil {
			result.S3 = mocks.S3
		}
	}
	return result
}

func isRequired(resource string, opts *BootstrapOptions) bool {
	for _, r := range opts.RequiredResources {
		if r == resource {
			return true
		}
	}
	return false
}

func logWarning(logger *logrus.Logger, msg string, err error) {
	if logger != nil {
		logger.WithError(err).Warn(msg)
	}
}

// extractEnvAndVersion extrae environment y version de la configuración
// Por ahora retorna valores por defecto, será implementado según BaseConfig
func extractEnvAndVersion(config interface{}) (string, string) {
	// TODO: Implementar extracción real cuando BaseConfig esté integrado
	return "local", "0.0.0"
}

// extractPostgreSQLConfig extrae configuración de PostgreSQL
func extractPostgreSQLConfig(config interface{}) (PostgreSQLConfig, error) {
	// TODO: Implementar extracción real cuando BaseConfig esté integrado
	return PostgreSQLConfig{}, fmt.Errorf("not implemented")
}

// extractMongoDBConfig extrae configuración de MongoDB
func extractMongoDBConfig(config interface{}) (MongoDBConfig, error) {
	// TODO: Implementar extracción real cuando BaseConfig esté integrado
	return MongoDBConfig{}, fmt.Errorf("not implemented")
}

// extractRabbitMQConfig extrae configuración de RabbitMQ
func extractRabbitMQConfig(config interface{}) (RabbitMQConfig, error) {
	// TODO: Implementar extracción real cuando BaseConfig esté integrado
	return RabbitMQConfig{}, fmt.Errorf("not implemented")
}

// extractS3Config extrae configuración de S3
func extractS3Config(config interface{}) (S3Config, error) {
	// TODO: Implementar extracción real cuando BaseConfig esté integrado
	return S3Config{}, fmt.Errorf("not implemented")
}

// Funciones de registro de cleanup (stubs por ahora)
func registerPostgreSQLCleanup(lifecycleManager interface{}, factory PostgreSQLFactory, db interface{}, logger *logrus.Logger) {
	// TODO: Implementar cuando lifecycle esté integrado
}

func registerMongoDBCleanup(lifecycleManager interface{}, factory MongoDBFactory, client interface{}, logger *logrus.Logger) {
	// TODO: Implementar cuando lifecycle esté integrado
}

func registerRabbitMQCleanup(lifecycleManager interface{}, factory RabbitMQFactory, channel interface{}, conn interface{}, logger *logrus.Logger) {
	// TODO: Implementar cuando lifecycle esté integrado
}
