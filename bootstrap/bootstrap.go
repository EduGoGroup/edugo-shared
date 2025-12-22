package bootstrap

import (
	"context"
	"fmt"
	"reflect"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
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
		client:        client,
		presignClient: factories.S3.CreatePresignClient(client),
		bucket:        s3Config.Bucket,
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

// extractEnvAndVersion extrae los campos Environment y Version de una configuración.
//
// Busca campos llamados "Environment" y "Version" en el struct proporcionado.
// Si no los encuentra o el config es nil, retorna valores por defecto.
//
// Parámetros:
//   - config: Struct de configuración (puede ser valor o puntero)
//
// Retorna:
//   - environment: Valor del campo Environment o "unknown"
//   - version: Valor del campo Version o "0.0.0"
func extractEnvAndVersion(config interface{}) (string, string) {
	if config == nil {
		return "unknown", "0.0.0"
	}

	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "unknown", "0.0.0"
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "unknown", "0.0.0"
	}

	// Buscar campo Environment
	env := "unknown"
	envField := v.FieldByName("Environment")
	if envField.IsValid() && envField.Kind() == reflect.String {
		env = envField.String()
		if env == "" {
			env = "unknown"
		}
	}

	// Buscar campo Version
	version := "0.0.0"
	versionField := v.FieldByName("Version")
	if versionField.IsValid() && versionField.Kind() == reflect.String {
		ver := versionField.String()
		if ver != "" {
			version = ver
		}
	}

	return env, version
}

// extractPostgreSQLConfig extrae configuración de PostgreSQL usando reflection
func extractPostgreSQLConfig(config interface{}) (PostgreSQLConfig, error) {
	// Intentar type assertion directo primero
	if pgConfig, ok := config.(PostgreSQLConfig); ok {
		return pgConfig, nil
	}

	// Usar reflection para extraer campo PostgreSQL
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return PostgreSQLConfig{}, fmt.Errorf("config must be a struct, got %T", config)
	}

	pgField := v.FieldByName("PostgreSQL")
	if !pgField.IsValid() {
		return PostgreSQLConfig{}, fmt.Errorf("PostgreSQL field not found in config")
	}

	if pgConfig, ok := pgField.Interface().(PostgreSQLConfig); ok {
		return pgConfig, nil
	}

	return PostgreSQLConfig{}, fmt.Errorf("PostgreSQL field is not of type PostgreSQLConfig")
}

// extractMongoDBConfig extrae configuración de MongoDB usando reflection
func extractMongoDBConfig(config interface{}) (MongoDBConfig, error) {
	// Intentar type assertion directo
	if mongoConfig, ok := config.(MongoDBConfig); ok {
		return mongoConfig, nil
	}

	// Usar reflection
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return MongoDBConfig{}, fmt.Errorf("config must be a struct, got %T", config)
	}

	mongoField := v.FieldByName("MongoDB")
	if !mongoField.IsValid() {
		return MongoDBConfig{}, fmt.Errorf("MongoDB field not found in config")
	}

	if mongoConfig, ok := mongoField.Interface().(MongoDBConfig); ok {
		return mongoConfig, nil
	}

	return MongoDBConfig{}, fmt.Errorf("MongoDB field is not of type MongoDBConfig")
}

// extractRabbitMQConfig extrae configuración de RabbitMQ usando reflection
func extractRabbitMQConfig(config interface{}) (RabbitMQConfig, error) {
	// Intentar type assertion directo
	if rabbitConfig, ok := config.(RabbitMQConfig); ok {
		return rabbitConfig, nil
	}

	// Usar reflection
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return RabbitMQConfig{}, fmt.Errorf("config must be a struct, got %T", config)
	}

	rabbitField := v.FieldByName("RabbitMQ")
	if !rabbitField.IsValid() {
		return RabbitMQConfig{}, fmt.Errorf("RabbitMQ field not found in config")
	}

	if rabbitConfig, ok := rabbitField.Interface().(RabbitMQConfig); ok {
		return rabbitConfig, nil
	}

	return RabbitMQConfig{}, fmt.Errorf("RabbitMQ field is not of type RabbitMQConfig")
}

// extractS3Config extrae configuración de S3 usando reflection
func extractS3Config(config interface{}) (S3Config, error) {
	// Intentar type assertion directo
	if s3Config, ok := config.(S3Config); ok {
		return s3Config, nil
	}

	// Usar reflection
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return S3Config{}, fmt.Errorf("config must be a struct, got %T", config)
	}

	s3Field := v.FieldByName("S3")
	if !s3Field.IsValid() {
		return S3Config{}, fmt.Errorf("S3 field not found in config")
	}

	if s3Config, ok := s3Field.Interface().(S3Config); ok {
		return s3Config, nil
	}

	return S3Config{}, fmt.Errorf("S3 field is not of type S3Config")
}

// Funciones de registro de cleanup (stubs por ahora)
func registerPostgreSQLCleanup(lifecycleManager interface{}, factory PostgreSQLFactory, db interface{}, logger *logrus.Logger) {
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

func registerMongoDBCleanup(lifecycleManager interface{}, factory MongoDBFactory, client interface{}, logger *logrus.Logger) {
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

func registerRabbitMQCleanup(lifecycleManager interface{}, factory RabbitMQFactory, channel interface{}, conn interface{}, logger *logrus.Logger) {
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
