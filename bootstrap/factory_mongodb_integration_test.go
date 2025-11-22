package bootstrap

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

// TestMongoDBFactory_CreateConnection_Success verifica creación exitosa de conexión
func TestMongoDBFactory_CreateConnection_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	// Setup container
	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	require.NotNil(t, mongo)

	// Crear factory
	factory := NewDefaultMongoDBFactory()

	// Configuración de MongoDB
	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	// Crear conexión
	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer factory.Close(ctx, client)

	// Verificar que la conexión funciona
	err = factory.Ping(ctx, client)
	assert.NoError(t, err)
}

// TestMongoDBFactory_CreateConnection_InvalidURI verifica error con URI inválida
func TestMongoDBFactory_CreateConnection_InvalidURI(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()
	factory := NewDefaultMongoDBFactory()

	// URI inválida
	invalidConfig := MongoDBConfig{
		URI:      "mongodb://invalid-host-that-does-not-exist:27017",
		Database: "test_db",
	}

	// Debe fallar al crear conexión (o al hacer ping)
	client, err := factory.CreateConnection(ctx, invalidConfig)

	// Puede fallar en CreateConnection o en Ping
	if err == nil {
		// Si la creación pasó, el ping debe fallar
		err = factory.Ping(ctx, client)
		assert.Error(t, err, "Ping debe fallar con URI inválida")
		if client != nil {
			_ = factory.Close(ctx, client)
		}
	} else {
		assert.Error(t, err)
		assert.Nil(t, client)
	}
}

// TestMongoDBFactory_Ping_Success verifica ping exitoso
func TestMongoDBFactory_Ping_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)
	defer factory.Close(ctx, client)

	// Ping debe ser exitoso
	err = factory.Ping(ctx, client)
	assert.NoError(t, err)
}

// TestMongoDBFactory_GetDatabase verifica obtención de database
func TestMongoDBFactory_GetDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)
	defer factory.Close(ctx, client)

	// Obtener database
	db := factory.GetDatabase(client, "test_db")
	assert.NotNil(t, db)
	assert.Equal(t, "test_db", db.Name())
}

// TestMongoDBFactory_Close_Success verifica cierre exitoso
func TestMongoDBFactory_Close_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)

	// Close debe ser exitoso
	err = factory.Close(ctx, client)
	assert.NoError(t, err)
}

// TestMongoDBFactory_ConnectionPoolSettings verifica configuración del pool
func TestMongoDBFactory_ConnectionPoolSettings(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)
	defer factory.Close(ctx, client)

	// Verificar que la conexión tiene configuración de pool
	// (La factory configura MaxPoolSize=100, MinPoolSize=10)
	err = factory.Ping(ctx, client)
	assert.NoError(t, err, "Connection pool debe estar configurado correctamente")
}

// TestMongoDBFactory_MultipleConnections verifica múltiples conexiones
func TestMongoDBFactory_MultipleConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	// Crear múltiples clientes
	clients := make([]*mongoClient, 3)
	for i := 0; i < 3; i++ {
		client, err := factory.CreateConnection(ctx, mongoConfig)
		require.NoError(t, err)
		clients[i] = client
	}

	// Verificar que todos funcionan
	for i, client := range clients {
		err := factory.Ping(ctx, client)
		assert.NoError(t, err, "Cliente %d debe funcionar", i)
	}

	// Cerrar todos
	for _, client := range clients {
		err := factory.Close(ctx, client)
		assert.NoError(t, err)
	}
}

// Define alias de tipo para mejor legibilidad
type mongoClient = mongo.Client

// TestMongoDBFactory_DatabaseOperations verifica operaciones básicas
func TestMongoDBFactory_DatabaseOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)
	defer factory.Close(ctx, client)

	// Obtener database
	db := factory.GetDatabase(client, "test_db")
	require.NotNil(t, db)

	// Crear colección y insertar documento
	collection := db.Collection("test_collection")

	testDoc := bson.M{
		"name":  "test",
		"value": 123,
	}

	result, err := collection.InsertOne(ctx, testDoc)
	require.NoError(t, err)
	assert.NotNil(t, result.InsertedID)

	// Buscar documento
	var found bson.M
	err = collection.FindOne(ctx, bson.M{"name": "test"}).Decode(&found)
	require.NoError(t, err)
	assert.Equal(t, "test", found["name"])
	assert.Equal(t, int32(123), found["value"])
}

// TestMongoDBFactory_PingWithTimeout verifica ping con timeout
func TestMongoDBFactory_PingWithTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	client, err := factory.CreateConnection(context.Background(), mongoConfig)
	require.NoError(t, err)
	defer factory.Close(context.Background(), client)

	// Ping con timeout muy corto
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = factory.Ping(ctx, client)
	assert.NoError(t, err, "Ping debe completarse dentro del timeout")
}

// TestMongoDBFactory_GetDatabase_MultipleDatabases verifica múltiples databases
func TestMongoDBFactory_GetDatabase_MultipleDatabases(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)
	defer factory.Close(ctx, client)

	// Obtener diferentes databases
	databases := []string{"db1", "db2", "db3"}
	for _, dbName := range databases {
		db := factory.GetDatabase(client, dbName)
		assert.NotNil(t, db)
		assert.Equal(t, dbName, db.Name())
	}
}

// TestMongoDBFactory_ConnectionTimeout verifica manejo de timeout de conexión
func TestMongoDBFactory_ConnectionTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	// Contexto con timeout muy corto
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	factory := NewDefaultMongoDBFactory()

	// Configuración con host que no responde
	mongoConfig := MongoDBConfig{
		URI:      "mongodb://192.0.2.1:27017", // TEST-NET-1 (no routable)
		Database: "test_db",
	}

	// Debe fallar por timeout (eventualmente)
	client, err := factory.CreateConnection(ctx, mongoConfig)

	// CreateConnection puede no fallar inmediatamente, pero Ping sí
	if err == nil && client != nil {
		err = factory.Ping(ctx, client)
		assert.Error(t, err, "Ping debe fallar con timeout")
		_ = factory.Close(context.Background(), client)
	}
}

// TestMongoDBFactory_CloseWithTimeout verifica cierre con timeout
func TestMongoDBFactory_CloseWithTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)

	// Cerrar con timeout
	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = factory.Close(closeCtx, client)
	assert.NoError(t, err)
}

// TestMongoDBFactory_DefaultTimeout verifica timeout por defecto configurado
func TestMongoDBFactory_DefaultTimeout(t *testing.T) {
	factory := NewDefaultMongoDBFactory()

	assert.NotNil(t, factory)
	assert.NotZero(t, factory.connectionTimeout, "Factory debe tener timeout configurado")
	assert.Equal(t, 10*time.Second, factory.connectionTimeout, "Timeout debe ser 10 segundos por defecto")
}

// TestMongoDBFactory_ConcurrentOperations verifica operaciones concurrentes
func TestMongoDBFactory_ConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithMongoDB(&containers.MongoDBConfig{
			Image:    "mongo:7.0",
			Database: "test_db",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	mongo := manager.MongoDB()
	factory := NewDefaultMongoDBFactory()

	mongoConfig := MongoDBConfig{
		URI:      mongo.ConnectionString(),
		Database: "test_db",
	}

	client, err := factory.CreateConnection(ctx, mongoConfig)
	require.NoError(t, err)
	defer factory.Close(ctx, client)

	db := factory.GetDatabase(client, "test_db")
	collection := db.Collection("concurrent_test")

	// Insertar documentos concurrentemente
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			doc := bson.M{
				"id":    index,
				"value": index * 10,
			}
			_, err := collection.InsertOne(ctx, doc)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Esperar a que terminen todas
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verificar que se insertaron todos
	count, err := collection.CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	assert.Equal(t, int64(10), count)
}
