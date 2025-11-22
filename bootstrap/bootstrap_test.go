package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// =============================================================================
// MOCK FACTORIES
// =============================================================================

// mockLoggerFactory es un mock de LoggerFactory
type mockLoggerFactory struct {
	shouldFail bool
	logger     *logrus.Logger
}

func (m *mockLoggerFactory) CreateLogger(ctx context.Context, env string, version string) (*logrus.Logger, error) {
	if m.shouldFail {
		return nil, errors.New("mock logger creation failed")
	}
	if m.logger != nil {
		return m.logger, nil
	}
	return logrus.New(), nil
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
	return nil, nil // Retornamos nil porque no podemos crear una conexión real
}

func (m *mockRabbitMQFactory) CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	if m.shouldFail {
		return nil, errors.New("mock rabbitmq channel failed")
	}
	return nil, nil
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
	return nil, nil
}

func (m *mockS3Factory) CreatePresignClient(client *s3.Client) interface{} {
	return nil
}

func (m *mockS3Factory) ValidateBucket(ctx context.Context, client *s3.Client, bucket string) error {
	return nil
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
// TESTS
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

// TestBootstrap_OptionalResourceFailure movido a bootstrap_integration_test.go

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

// TestDefaultBootstrapOptions duplicado - ver options_test.go

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
	resources.Logger = logrus.New()
	if !resources.HasLogger() {
		t.Error("Expected HasLogger to return true after initialization")
	}
}

// =============================================================================
// TESTS FOR performHealthChecks
// =============================================================================

func TestPerformHealthChecks_AllPass(t *testing.T) {
	logger := logrus.New()

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
	logger := logrus.New()
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
