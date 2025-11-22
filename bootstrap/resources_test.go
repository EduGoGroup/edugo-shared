package bootstrap

import (
	"context"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResources_EmptyInitialization verifica inicialización vacía
func TestResources_EmptyInitialization(t *testing.T) {
	resources := &Resources{}

	assert.False(t, resources.HasLogger())
	assert.False(t, resources.HasPostgreSQL())
	assert.False(t, resources.HasMongoDB())
	assert.False(t, resources.HasMessagePublisher())
	assert.False(t, resources.HasStorageClient())
}

// TestResources_HasLogger verifica detección de logger
func TestResources_HasLogger(t *testing.T) {
	resources := &Resources{}

	// Sin logger
	assert.False(t, resources.HasLogger())

	// Con logger
	resources.Logger = logrus.New()
	assert.True(t, resources.HasLogger())
}

// TestResources_HasPostgreSQL verifica detección de PostgreSQL
func TestResources_HasPostgreSQL(t *testing.T) {
	skipIfNotIntegration(t)

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	factory := NewDefaultPostgreSQLFactory(nil)

	// Obtener config dinámico del container
	connStr, err := pg.ConnectionString(ctx)
	require.NoError(t, err)

	pgConfig, err := parsePostgresConnectionString(connStr)
	require.NoError(t, err)

	db, err := factory.CreateConnection(ctx, pgConfig)
	require.NoError(t, err)
	defer factory.Close(db)

	resources := &Resources{}

	// Sin PostgreSQL
	assert.False(t, resources.HasPostgreSQL())

	// Con PostgreSQL
	resources.PostgreSQL = db
	assert.True(t, resources.HasPostgreSQL())
}

// TestResources_HasMongoDB verifica detección de MongoDB
func TestResources_HasMongoDB(t *testing.T) {
	skipIfNotIntegration(t)

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoURI, err := mongo.ConnectionString(ctx)
	require.NoError(t, err)

	mongoConfig := MongoDBConfig{
		URI:      mongoURI,
		Database: "test_db",
	}

	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)
	defer factory.Close(ctx, client)

	db := factory.GetDatabase(client, "test_db")

	resources := &Resources{}

	// Sin MongoDB
	assert.False(t, resources.HasMongoDB())

	// Con client pero sin database
	resources.MongoDB = client
	assert.False(t, resources.HasMongoDB(), "Debe tener ambos client y database")

	// Con ambos
	resources.MongoDatabase = db
	assert.True(t, resources.HasMongoDB())
}

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

// TestResources_AllResourcesPresent verifica todos los recursos presentes
func TestResources_AllResourcesPresent(t *testing.T) {
	skipIfNotIntegration(t)

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		WithMongoDB(&containers.MongoConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	// Setup PostgreSQL
	pg := manager.PostgreSQL()
	pgFactory := NewDefaultPostgreSQLFactory(nil)

	connStr, err := pg.ConnectionString(ctx)
	require.NoError(t, err)

	pgConfig, err := parsePostgresConnectionString(connStr)
	require.NoError(t, err)

	pgDB, err := pgFactory.CreateConnection(ctx, pgConfig)
	require.NoError(t, err)
	defer pgFactory.Close(pgDB)

	// Setup MongoDB
	mongo := manager.MongoDB()
	mongoFactory := NewDefaultMongoDBFactory()

	mongoURI, err := mongo.ConnectionString(ctx)
	require.NoError(t, err)

	mongoConfig := MongoDBConfig{
		URI:      mongoURI,
		Database: "test_db",
	}
	mongoClient, err := mongoFactory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)
	defer mongoFactory.Close(ctx, mongoClient)
	mongoDB := mongoFactory.GetDatabase(mongoClient, "test_db")

	// Crear resources con todo
	resources := &Resources{
		Logger:           logrus.New(),
		PostgreSQL:       pgDB,
		MongoDB:          mongoClient,
		MongoDatabase:    mongoDB,
		MessagePublisher: &mockMessagePublisher{},
		StorageClient:    &mockStorageClient{},
	}

	assert.True(t, resources.HasLogger())
	assert.True(t, resources.HasPostgreSQL())
	assert.True(t, resources.HasMongoDB())
	assert.True(t, resources.HasMessagePublisher())
	assert.True(t, resources.HasStorageClient())
}

// TestResources_PartialInitialization verifica inicialización parcial
func TestResources_PartialInitialization(t *testing.T) {
	// Solo logger y PostgreSQL
	resources := &Resources{
		Logger: logrus.New(),
		// PostgreSQL se inicializaría en un test real
	}

	assert.True(t, resources.HasLogger())
	assert.False(t, resources.HasPostgreSQL())
	assert.False(t, resources.HasMongoDB())
	assert.False(t, resources.HasMessagePublisher())
	assert.False(t, resources.HasStorageClient())
}

// Mocks para testing
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
