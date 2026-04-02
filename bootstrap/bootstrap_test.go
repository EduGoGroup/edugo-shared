package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

// =============================================================================
// MOCK FACTORIES
// =============================================================================

// mockLoggerFactory es un mock de LoggerFactory
type mockLoggerFactory struct {
	shouldFail bool
	logger     logger.Logger
}

func (m *mockLoggerFactory) CreateLogger(ctx context.Context, env string, version string) (logger.Logger, error) {
	if m.shouldFail {
		return nil, errors.New("mock logger creation failed")
	}
	if m.logger != nil {
		return m.logger, nil
	}
	return logger.NewLogrusLogger(logrus.New()), nil
}

// mockPostgreSQLFactory es un mock de PostgreSQLFactory
type mockPostgreSQLFactory struct {
	shouldFail bool
	db         *gorm.DB
}

func (m *mockPostgreSQLFactory) CreateConnection(ctx context.Context, config PostgreSQLConfig) (*gorm.DB, error) {
	if m.shouldFail {
		return nil, errors.New("mock postgresql connection failed")
	}
	return m.db, nil
}

func (m *mockPostgreSQLFactory) CreateRawConnection(ctx context.Context, config PostgreSQLConfig) (*sql.DB, error) {
	return nil, errors.New("not implemented in mock")
}

func (m *mockPostgreSQLFactory) Ping(ctx context.Context, db *gorm.DB) error {
	return nil
}

func (m *mockPostgreSQLFactory) Close(db *gorm.DB) error {
	return nil
}

// mockMongoDBFactory es un mock de MongoDBFactory
type mockMongoDBFactory struct {
	shouldFail bool
	client     *mongo.Client
	database   *mongo.Database
}

func (m *mockMongoDBFactory) CreateConnection(ctx context.Context, config MongoDBConfig) (*mongo.Client, error) {
	if m.shouldFail {
		return nil, errors.New("mock mongodb connection failed")
	}
	return m.client, nil
}

func (m *mockMongoDBFactory) GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return m.database
}

func (m *mockMongoDBFactory) Ping(ctx context.Context, client *mongo.Client) error {
	return nil
}

func (m *mockMongoDBFactory) Close(ctx context.Context, client *mongo.Client) error {
	return nil
}

// mockRabbitMQFactory es un mock de RabbitMQFactory
type mockRabbitMQFactory struct {
	shouldFail bool
}

func (m *mockRabbitMQFactory) CreateConnection(ctx context.Context, config RabbitMQConfig) (*amqp.Connection, error) {
	if m.shouldFail {
		return nil, errors.New("mock rabbitmq connection failed")
	}
	return nil, nil //nolint:nilnil // Mock: no podemos crear una conexión real en tests
}

func (m *mockRabbitMQFactory) CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	if m.shouldFail {
		return nil, errors.New("mock rabbitmq channel failed")
	}
	return nil, nil //nolint:nilnil // Mock: no podemos crear un canal real en tests
}

func (m *mockRabbitMQFactory) DeclareQueue(channel *amqp.Channel, queueName string) (amqp.Queue, error) {
	return amqp.Queue{Name: queueName}, nil
}

func (m *mockRabbitMQFactory) Close(channel *amqp.Channel, conn *amqp.Connection) error {
	return nil
}

// mockS3Factory es un mock de S3Factory
type mockS3Factory struct {
	shouldFail bool
}

func (m *mockS3Factory) CreateClient(ctx context.Context, config S3Config) (*s3.Client, error) {
	if m.shouldFail {
		return nil, errors.New("mock s3 client failed")
	}
	return nil, nil //nolint:nilnil // Mock: no podemos crear un cliente S3 real en tests
}

func (m *mockS3Factory) CreatePresignClient(client *s3.Client) *s3.PresignClient {
	if client == nil {
		return nil
	}
	return s3.NewPresignClient(client)
}

func (m *mockS3Factory) ValidateBucket(ctx context.Context, client *s3.Client, bucket string) error {
	return nil
}

// mockLifecycleManager es un mock de lifecycle manager para cleanup tests
type mockLifecycleManager struct {
	registered map[string]func() error
}

func (m *mockLifecycleManager) RegisterSimple(name string, cleanup func() error) {
	m.registered[name] = cleanup
}

// mockMessagePublisher es un mock de MessagePublisher para resources tests
type mockMessagePublisher struct{}

func (m *mockMessagePublisher) Publish(ctx context.Context, queueName string, body []byte) error {
	return nil
}

func (m *mockMessagePublisher) PublishWithPriority(ctx context.Context, queueName string, body []byte, priority uint8) error {
	return nil
}

func (m *mockMessagePublisher) Close() error {
	return nil
}

// mockStorageClient es un mock de StorageClient para resources tests
type mockStorageClient struct{}

func (m *mockStorageClient) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	return "mock-url", nil
}

func (m *mockStorageClient) Download(ctx context.Context, key string) ([]byte, error) {
	return []byte("mock-data"), nil
}

func (m *mockStorageClient) Delete(ctx context.Context, key string) error {
	return nil
}

func (m *mockStorageClient) GetPresignedURL(ctx context.Context, key string, expirationMinutes int) (string, error) {
	return "mock-presigned-url", nil
}

func (m *mockStorageClient) Exists(ctx context.Context, key string) (bool, error) {
	return true, nil
}

// =============================================================================
// HELPER FUNCTIONS FOR TESTS
// =============================================================================

func createMockFactories(loggerFail, pgFail, mongoFail, rabbitFail, s3Fail bool) *MockFactories {
	return &MockFactories{
		Logger:     &mockLoggerFactory{shouldFail: loggerFail},
		PostgreSQL: &mockPostgreSQLFactory{shouldFail: pgFail},
		MongoDB:    &mockMongoDBFactory{shouldFail: mongoFail},
		RabbitMQ:   &mockRabbitMQFactory{shouldFail: rabbitFail},
		S3:         &mockS3Factory{shouldFail: s3Fail},
	}
}

// =============================================================================
// BOOTSTRAP TESTS
// =============================================================================

func TestBootstrap_LoggerRequired(t *testing.T) {
	ctx := context.Background()

	// Test: Logger factory es obligatoria
	factories := &Factories{}

	_, err := Bootstrap(ctx, nil, factories, nil, WithRequiredResources("logger"))
	if err == nil {
		t.Fatal("Expected error when logger factory is missing, got nil")
	}
}

func TestBootstrap_LoggerCreationSuccess(t *testing.T) {
	ctx := context.Background()

	// Setup mocks
	mocks := createMockFactories(false, false, false, false, false)
	factories := &Factories{}

	// Test: Logger se crea exitosamente
	resources, err := Bootstrap(
		ctx,
		nil,
		factories,
		nil,
		WithRequiredResources("logger"),
		WithMockFactories(mocks),
		WithSkipHealthCheck(),
	)

	if err != nil {
		t.Fatalf("Expected successful bootstrap, got error: %v", err)
	}

	if resources == nil {
		t.Fatal("Expected resources to be non-nil")
	}

	if resources.Logger == nil {
		t.Fatal("Expected logger to be initialized")
	}
}

func TestBootstrap_LoggerCreationFailure(t *testing.T) {
	ctx := context.Background()

	// Setup mocks con fallo en logger
	mocks := createMockFactories(true, false, false, false, false)
	factories := &Factories{}

	// Test: Fallo en creación de logger
	_, err := Bootstrap(
		ctx,
		nil,
		factories,
		nil,
		WithRequiredResources("logger"),
		WithMockFactories(mocks),
	)

	if err == nil {
		t.Fatal("Expected error when logger creation fails, got nil")
	}
}

func TestBootstrap_RequiredResourceFailure(t *testing.T) {
	ctx := context.Background()

	// Setup mocks: logger ok, postgresql falla
	mocks := createMockFactories(false, true, false, false, false)
	factories := &Factories{}

	// Test: PostgreSQL es requerido, debe fallar el bootstrap
	_, err := Bootstrap(
		ctx,
		nil,
		factories,
		nil,
		WithRequiredResources("logger", "postgresql"),
		WithMockFactories(mocks),
		WithSkipHealthCheck(),
	)

	if err == nil {
		t.Fatal("Expected error when required resource fails, got nil")
	}
}

func TestBootstrap_AllResourcesSuccess(t *testing.T) {
	ctx := context.Background()

	// Setup mocks: todos exitosos
	mocks := createMockFactories(false, false, false, false, false)
	factories := &Factories{}

	// Test: Todos los recursos se inicializan exitosamente
	resources, err := Bootstrap(
		ctx,
		nil,
		factories,
		nil,
		WithRequiredResources("logger"),
		WithOptionalResources("postgresql", "mongodb", "rabbitmq", "s3"),
		WithMockFactories(mocks),
		WithSkipHealthCheck(),
	)

	if err != nil {
		t.Fatalf("Expected successful bootstrap, got error: %v", err)
	}

	// Verificar que logger está inicializado
	if resources.Logger == nil {
		t.Error("Expected logger to be initialized")
	}

	// Los demás recursos pueden o no estar inicializados dependiendo de los mocks
	// (no verificamos porque los mocks retornan nil en algunos casos)
}

func TestBootstrap_StopOnFirstError(t *testing.T) {
	ctx := context.Background()

	// Setup mocks: postgresql falla
	mocks := createMockFactories(false, true, false, false, false)
	factories := &Factories{}

	// Test: Con StopOnFirstError, debe detenerse en el primer error
	_, err := Bootstrap(
		ctx,
		nil,
		factories,
		nil,
		WithRequiredResources("logger", "postgresql"),
		WithMockFactories(mocks),
		WithStopOnFirstError(true),
	)

	if err == nil {
		t.Fatal("Expected error with StopOnFirstError, got nil")
	}
}

func TestBootstrap_SkipHealthCheck(t *testing.T) {
	ctx := context.Background()

	// Setup mocks
	mocks := createMockFactories(false, false, false, false, false)
	factories := &Factories{}

	// Test: Con SkipHealthCheck, no debe ejecutar health checks
	resources, err := Bootstrap(
		ctx,
		nil,
		factories,
		nil,
		WithRequiredResources("logger"),
		WithMockFactories(mocks),
		WithSkipHealthCheck(),
	)

	if err != nil {
		t.Fatalf("Expected successful bootstrap, got error: %v", err)
	}

	if resources == nil {
		t.Fatal("Expected resources to be non-nil")
	}
}

func TestFactoriesValidate(t *testing.T) {
	// Test: Validate con factory faltante
	factories := &Factories{
		Logger: &mockLoggerFactory{},
	}

	err := factories.Validate([]string{"logger", "postgresql"})
	if err == nil {
		t.Error("Expected validation error for missing postgresql factory")
	}

	// Test: Validate exitoso
	factories.PostgreSQL = &mockPostgreSQLFactory{}
	err = factories.Validate([]string{"logger", "postgresql"})
	if err != nil {
		t.Errorf("Expected successful validation, got error: %v", err)
	}
}

func TestResources_HasMethods(t *testing.T) {
	resources := &Resources{}

	// Test: Todos los recursos vacíos
	if resources.HasLogger() {
		t.Error("Expected HasLogger to return false for nil logger")
	}
	if resources.HasPostgreSQL() {
		t.Error("Expected HasPostgreSQL to return false for nil db")
	}
	if resources.HasMongoDB() {
		t.Error("Expected HasMongoDB to return false for nil client")
	}
	if resources.HasMessagePublisher() {
		t.Error("Expected HasMessagePublisher to return false for nil publisher")
	}
	if resources.HasStorageClient() {
		t.Error("Expected HasStorageClient to return false for nil client")
	}

	// Test: Logger inicializado
	resources.Logger = logger.NewLogrusLogger(logrus.New())
	if !resources.HasLogger() {
		t.Error("Expected HasLogger to return true after initialization")
	}
}

// =============================================================================
// TESTS FOR performHealthChecks
// =============================================================================

func TestPerformHealthChecks_AllPass(t *testing.T) {
	logger := logger.NewLogrusLogger(logrus.New())

	resources := &Resources{
		Logger: logger,
	}

	opts := DefaultBootstrapOptions()
	ctx := context.Background()

	err := performHealthChecks(ctx, resources, opts)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestPerformHealthChecks_WithoutLogger(t *testing.T) {
	resources := &Resources{}
	opts := DefaultBootstrapOptions()
	ctx := context.Background()

	err := performHealthChecks(ctx, resources, opts)
	if err != nil {
		t.Errorf("Expected no error without logger, got: %v", err)
	}
}

func TestPerformHealthChecks_ContextTimeout(t *testing.T) {
	logger := logger.NewLogrusLogger(logrus.New())
	resources := &Resources{Logger: logger}
	opts := DefaultBootstrapOptions()

	// Crear un contexto ya cancelado
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// El health check debería manejar el contexto cancelado correctamente
	err := performHealthChecks(ctx, resources, opts)
	// No debería fallar porque no hay recursos que validar
	if err != nil {
		t.Errorf("Expected no error with cancelled context and no resources, got: %v", err)
	}
}

func TestExtractEnvAndVersion(t *testing.T) {
	tests := []struct {
		name        string
		config      any
		wantEnv     string
		wantVersion string
	}{
		{
			name:        "nil config returns defaults",
			config:      nil,
			wantEnv:     "unknown",
			wantVersion: "0.0.0",
		},
		{
			name: "struct with Environment and Version",
			config: struct {
				Environment string
				Version     string
			}{
				Environment: "prod",
				Version:     "1.2.3",
			},
			wantEnv:     "prod",
			wantVersion: "1.2.3",
		},
		{
			name: "struct with only Environment",
			config: struct {
				Environment string
			}{
				Environment: "dev",
			},
			wantEnv:     "dev",
			wantVersion: "0.0.0",
		},
		{
			name: "pointer to struct",
			config: &struct {
				Environment string
				Version     string
			}{
				Environment: "qa",
				Version:     "2.0.0",
			},
			wantEnv:     "qa",
			wantVersion: "2.0.0",
		},
		{
			name: "empty strings use defaults",
			config: struct {
				Environment string
				Version     string
			}{
				Environment: "",
				Version:     "",
			},
			wantEnv:     "unknown",
			wantVersion: "0.0.0",
		},
		{
			name: "struct without fields",
			config: struct {
				Other string
			}{
				Other: "value",
			},
			wantEnv:     "unknown",
			wantVersion: "0.0.0",
		},
		{
			name:        "non-struct config",
			config:      "string",
			wantEnv:     "unknown",
			wantVersion: "0.0.0",
		},
		{
			name: "nil pointer to struct",
			config: (*struct {
				Environment string
				Version     string
			})(nil),
			wantEnv:     "unknown",
			wantVersion: "0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEnv, gotVersion := extractEnvAndVersion(tt.config)
			if gotEnv != tt.wantEnv {
				t.Errorf("extractEnvAndVersion() env = %v, want %v", gotEnv, tt.wantEnv)
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("extractEnvAndVersion() version = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
}

// =============================================================================
// OPTIONS TESTS
// =============================================================================

// TestDefaultBootstrapOptions verifica opciones por defecto
func TestDefaultBootstrapOptions(t *testing.T) {
	opts := DefaultBootstrapOptions()

	assert.NotNil(t, opts)
	assert.Equal(t, []string{"logger"}, opts.RequiredResources)
	assert.Equal(t, []string{}, opts.OptionalResources)
	assert.False(t, opts.SkipHealthCheck)
	assert.Nil(t, opts.MockFactories)
	assert.True(t, opts.StopOnFirstError)
}

// TestWithRequiredResources verifica configuración de recursos requeridos
func TestWithRequiredResources(t *testing.T) {
	opts := DefaultBootstrapOptions()

	// Aplicar opción
	option := WithRequiredResources("postgresql", "mongodb", "rabbitmq")
	option(opts)

	assert.Equal(t, []string{"postgresql", "mongodb", "rabbitmq"}, opts.RequiredResources)
}

// TestWithOptionalResources verifica configuración de recursos opcionales
func TestWithOptionalResources(t *testing.T) {
	opts := DefaultBootstrapOptions()

	// Aplicar opción
	option := WithOptionalResources("s3", "redis")
	option(opts)

	assert.Equal(t, []string{"s3", "redis"}, opts.OptionalResources)
}

// TestWithSkipHealthCheck verifica configuración de skip health check
func TestWithSkipHealthCheck(t *testing.T) {
	opts := DefaultBootstrapOptions()

	// Por defecto no está skip
	assert.False(t, opts.SkipHealthCheck)

	// Aplicar opción
	option := WithSkipHealthCheck()
	option(opts)

	assert.True(t, opts.SkipHealthCheck)
}

// TestWithMockFactories verifica inyección de factories simuladas
func TestWithMockFactories(t *testing.T) {
	opts := DefaultBootstrapOptions()

	// Por defecto no hay mocks
	assert.Nil(t, opts.MockFactories)

	// Crear mocks
	mocks := &MockFactories{
		Logger: NewDefaultLoggerFactory(),
	}

	// Aplicar opción
	option := WithMockFactories(mocks)
	option(opts)

	assert.NotNil(t, opts.MockFactories)
	assert.NotNil(t, opts.MockFactories.Logger)
}

// TestWithStopOnFirstError verifica configuración de stop on error
func TestWithStopOnFirstError(t *testing.T) {
	opts := DefaultBootstrapOptions()

	// Por defecto es true
	assert.True(t, opts.StopOnFirstError)

	// Cambiar a false
	option := WithStopOnFirstError(false)
	option(opts)

	assert.False(t, opts.StopOnFirstError)

	// Cambiar de vuelta a true
	option2 := WithStopOnFirstError(true)
	option2(opts)

	assert.True(t, opts.StopOnFirstError)
}

// TestApplyOptions verifica aplicación de múltiples opciones
func TestApplyOptions(t *testing.T) {
	opts := DefaultBootstrapOptions()

	options := []BootstrapOption{
		WithRequiredResources("postgresql", "mongodb"),
		WithOptionalResources("s3"),
		WithSkipHealthCheck(),
		WithStopOnFirstError(false),
	}

	ApplyOptions(opts, options...)

	assert.Equal(t, []string{"postgresql", "mongodb"}, opts.RequiredResources)
	assert.Equal(t, []string{"s3"}, opts.OptionalResources)
	assert.True(t, opts.SkipHealthCheck)
	assert.False(t, opts.StopOnFirstError)
}

// TestApplyOptions_Empty verifica aplicación de opciones vacías
func TestApplyOptions_Empty(t *testing.T) {
	opts := DefaultBootstrapOptions()
	original := *opts

	ApplyOptions(opts)

	// No debe cambiar nada
	assert.Equal(t, original, *opts)
}

// TestApplyOptions_OrderMatters verifica que el orden importa
func TestApplyOptions_OrderMatters(t *testing.T) {
	opts := DefaultBootstrapOptions()

	// Aplicar opciones en orden: primero true, luego false
	ApplyOptions(opts,
		WithStopOnFirstError(true),
		WithStopOnFirstError(false),
	)

	// El último valor debe prevalecer
	assert.False(t, opts.StopOnFirstError)
}

// TestMockFactories_Structure verifica estructura de MockFactories
func TestMockFactories_Structure(t *testing.T) {
	mocks := &MockFactories{
		Logger:     NewDefaultLoggerFactory(),
		PostgreSQL: NewDefaultPostgreSQLFactory(nil),
		MongoDB:    NewDefaultMongoDBFactory(),
		RabbitMQ:   NewDefaultRabbitMQFactory(),
		S3:         NewDefaultS3Factory(),
	}

	assert.NotNil(t, mocks.Logger)
	assert.NotNil(t, mocks.PostgreSQL)
	assert.NotNil(t, mocks.MongoDB)
	assert.NotNil(t, mocks.RabbitMQ)
	assert.NotNil(t, mocks.S3)
}

// TestBootstrapOptions_MultipleResourceTypes verifica configuración de múltiples tipos
func TestBootstrapOptions_MultipleResourceTypes(t *testing.T) {
	opts := DefaultBootstrapOptions()

	ApplyOptions(opts,
		WithRequiredResources("logger", "postgresql"),
		WithOptionalResources("mongodb", "rabbitmq", "s3"),
	)

	assert.Len(t, opts.RequiredResources, 2)
	assert.Len(t, opts.OptionalResources, 3)
	assert.Contains(t, opts.RequiredResources, "logger")
	assert.Contains(t, opts.RequiredResources, "postgresql")
	assert.Contains(t, opts.OptionalResources, "mongodb")
	assert.Contains(t, opts.OptionalResources, "rabbitmq")
	assert.Contains(t, opts.OptionalResources, "s3")
}

// =============================================================================
// CLEANUP TESTS
// =============================================================================

func TestRegisterPostgreSQLCleanup(t *testing.T) {
	t.Run("with nil lifecycle manager", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())

		assert.NotPanics(t, func() {
			registerPostgreSQLCleanup(nil, nil, nil, logger)
		})
	})

	t.Run("with lifecycle manager that doesn't implement RegisterSimple", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := struct{}{}

		assert.NotPanics(t, func() {
			registerPostgreSQLCleanup(lifecycleManager, nil, nil, logger)
		})
	})

	t.Run("with nil factory", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}

		registerPostgreSQLCleanup(lifecycleManager, nil, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with nil db", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockPostgreSQLFactory{}

		registerPostgreSQLCleanup(lifecycleManager, factory, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with wrong type of db", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockPostgreSQLFactory{}
		wrongDB := "not a *gorm.DB"

		registerPostgreSQLCleanup(lifecycleManager, factory, wrongDB, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("successful registration", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockPostgreSQLFactory{}
		db := &gorm.DB{}

		registerPostgreSQLCleanup(lifecycleManager, factory, db, logger)

		assert.Contains(t, lifecycleManager.registered, "postgresql")
		assert.NotNil(t, lifecycleManager.registered["postgresql"])
	})
}

func TestRegisterMongoDBCleanup(t *testing.T) {
	t.Run("with nil lifecycle manager", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())

		assert.NotPanics(t, func() {
			registerMongoDBCleanup(nil, nil, nil, logger)
		})
	})

	t.Run("with nil factory", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}

		registerMongoDBCleanup(lifecycleManager, nil, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with nil client", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockMongoDBFactory{}

		registerMongoDBCleanup(lifecycleManager, factory, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with wrong type of client", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockMongoDBFactory{}
		wrongClient := "not a *mongo.Client"

		registerMongoDBCleanup(lifecycleManager, factory, wrongClient, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("successful registration", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockMongoDBFactory{}
		client := &mongo.Client{}

		registerMongoDBCleanup(lifecycleManager, factory, client, logger)

		assert.Contains(t, lifecycleManager.registered, "mongodb")
		assert.NotNil(t, lifecycleManager.registered["mongodb"])
	})
}

func TestRegisterRabbitMQCleanup(t *testing.T) {
	t.Run("with nil lifecycle manager", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())

		assert.NotPanics(t, func() {
			registerRabbitMQCleanup(nil, nil, nil, nil, logger)
		})
	})

	t.Run("with nil factory", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}

		registerRabbitMQCleanup(lifecycleManager, nil, nil, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with nil channel", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockRabbitMQFactory{}

		registerRabbitMQCleanup(lifecycleManager, factory, nil, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with wrong type of channel", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockRabbitMQFactory{}
		wrongChannel := "not a *amqp.Channel"
		conn := &amqp.Connection{}

		registerRabbitMQCleanup(lifecycleManager, factory, wrongChannel, conn, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with wrong type of connection", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockRabbitMQFactory{}
		channel := &amqp.Channel{}
		wrongConn := "not a *amqp.Connection"

		registerRabbitMQCleanup(lifecycleManager, factory, channel, wrongConn, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("successful registration", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockRabbitMQFactory{}
		channel := &amqp.Channel{}
		conn := &amqp.Connection{}

		registerRabbitMQCleanup(lifecycleManager, factory, channel, conn, logger)

		assert.Contains(t, lifecycleManager.registered, "rabbitmq")
		assert.NotNil(t, lifecycleManager.registered["rabbitmq"])
	})
}

// =============================================================================
// CONFIG EXTRACTION TESTS
// =============================================================================

func TestExtractPostgreSQLConfig(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		_, err := extractPostgreSQLConfig(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with non-struct config", func(t *testing.T) {
		_, err := extractPostgreSQLConfig("not a struct")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with struct without PostgreSQL field", func(t *testing.T) {
		config := struct {
			OtherField string
		}{
			OtherField: "value",
		}

		_, err := extractPostgreSQLConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PostgreSQL field not found")
	})

	t.Run("successful extraction", func(t *testing.T) {
		expected := PostgreSQLConfig{Host: "localhost", Port: 5432}
		config := struct{ PostgreSQL PostgreSQLConfig }{PostgreSQL: expected}
		res, err := extractPostgreSQLConfig(config)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("successful extraction from pointer", func(t *testing.T) {
		expected := PostgreSQLConfig{Host: "localhost", Port: 5432}
		config := &struct{ PostgreSQL PostgreSQLConfig }{PostgreSQL: expected}
		res, err := extractPostgreSQLConfig(config)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestExtractMongoDBConfig(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		_, err := extractMongoDBConfig(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with non-struct config", func(t *testing.T) {
		_, err := extractMongoDBConfig(123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("successful extraction", func(t *testing.T) {
		expected := MongoDBConfig{URI: "mongodb://localhost"}
		config := struct{ MongoDB MongoDBConfig }{MongoDB: expected}
		res, err := extractMongoDBConfig(config)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestExtractRabbitMQConfig(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		_, err := extractRabbitMQConfig(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with non-struct config", func(t *testing.T) {
		_, err := extractRabbitMQConfig([]string{"not", "a", "struct"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("successful extraction", func(t *testing.T) {
		expected := RabbitMQConfig{URL: "amqp://localhost"}
		config := struct{ RabbitMQ RabbitMQConfig }{RabbitMQ: expected}
		res, err := extractRabbitMQConfig(config)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestExtractS3Config(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		_, err := extractS3Config(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with non-struct config", func(t *testing.T) {
		_, err := extractS3Config(make(map[string]string))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("successful extraction", func(t *testing.T) {
		expected := S3Config{Bucket: "my-bucket"}
		config := struct{ S3 S3Config }{S3: expected}
		res, err := extractS3Config(config)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

// =============================================================================
// RESOURCES TESTS
// =============================================================================

// TestResources_HasMessagePublisher verifica detección de message publisher
func TestResources_HasMessagePublisher(t *testing.T) {
	resources := &Resources{}

	// Sin publisher
	assert.False(t, resources.HasMessagePublisher())

	// Con publisher mock
	resources.MessagePublisher = &mockMessagePublisher{}
	assert.True(t, resources.HasMessagePublisher())
}

// TestResources_HasStorageClient verifica detección de storage client
func TestResources_HasStorageClient(t *testing.T) {
	resources := &Resources{}

	// Sin storage
	assert.False(t, resources.HasStorageClient())

	// Con storage mock
	resources.StorageClient = &mockStorageClient{}
	assert.True(t, resources.HasStorageClient())
}

// TestResources_PartialInitialization verifica inicialización parcial
func TestResources_PartialInitialization(t *testing.T) {
	// Solo logger y PostgreSQL
	resources := &Resources{
		Logger: logger.NewLogrusLogger(logrus.New()),
		// PostgreSQL se inicializaría en un test real
	}

	assert.True(t, resources.HasLogger())
	assert.False(t, resources.HasPostgreSQL())
	assert.False(t, resources.HasMongoDB())
	assert.False(t, resources.HasMessagePublisher())
	assert.False(t, resources.HasStorageClient())
}
