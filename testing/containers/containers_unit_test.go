package containers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Tests unitarios para PostgresContainer accessors (sin Docker)
// =============================================================================

// TestPostgresContainer_Username verifica que Username() retorna el valor configurado
func TestPostgresContainer_Username(t *testing.T) {
	pc := &PostgresContainer{
		config: &PostgresConfig{
			Username: "custom_user",
		},
	}

	assert.Equal(t, "custom_user", pc.Username())
}

// TestPostgresContainer_Password verifica que Password() retorna el valor configurado
func TestPostgresContainer_Password(t *testing.T) {
	pc := &PostgresContainer{
		config: &PostgresConfig{
			Password: "custom_pass",
		},
	}

	assert.Equal(t, "custom_pass", pc.Password())
}

// TestPostgresContainer_Database verifica que Database() retorna el valor configurado
func TestPostgresContainer_Database(t *testing.T) {
	pc := &PostgresContainer{
		config: &PostgresConfig{
			Database: "custom_db",
		},
	}

	assert.Equal(t, "custom_db", pc.Database())
}

// TestPostgresContainer_DB_Nil verifica que DB() retorna nil cuando no hay conexión
func TestPostgresContainer_DB_Nil(t *testing.T) {
	pc := &PostgresContainer{
		config: &PostgresConfig{},
	}

	assert.Nil(t, pc.DB())
}

// TestPostgresContainer_Accessors_AllFields verifica todos los accessors con una config completa
func TestPostgresContainer_Accessors_AllFields(t *testing.T) {
	cfg := &PostgresConfig{
		Image:    "postgres:16",
		Database: "prod_db",
		Username: "prod_user",
		Password: "prod_pass",
		Port:     "5433",
	}

	pc := &PostgresContainer{config: cfg}

	assert.Equal(t, "prod_db", pc.Database())
	assert.Equal(t, "prod_user", pc.Username())
	assert.Equal(t, "prod_pass", pc.Password())
}

// TestPostgresContainer_Truncate_EmptyList verifica que Truncate con lista vacía no da error
func TestPostgresContainer_Truncate_EmptyList(t *testing.T) {
	pc := &PostgresContainer{
		config: &PostgresConfig{Database: "test_db"},
	}

	err := pc.Truncate(context.Background())
	assert.NoError(t, err, "Truncate con lista vacía debe retornar nil")
}

// TestPostgresContainer_Terminate_NilFields verifica Terminate con campos nil
func TestPostgresContainer_Terminate_NilFields(t *testing.T) {
	pc := &PostgresContainer{
		config: &PostgresConfig{},
		// db y container son nil
	}

	err := pc.Terminate(context.Background())
	assert.NoError(t, err, "Terminate con campos nil no debe dar error")
}

// TestPostgresContainer_Accessors_DefaultValues verifica accessors con valores default
func TestPostgresContainer_Accessors_DefaultValues(t *testing.T) {
	// Simular un container creado con defaults
	cfg := &PostgresConfig{}
	// Aplicar los mismos defaults que WithPostgreSQL
	config := NewConfig().WithPostgreSQL(cfg).Build()

	pc := &PostgresContainer{config: config.PostgresConfig}

	assert.Equal(t, "edugo_test", pc.Database())
	assert.Equal(t, "edugo_user", pc.Username())
	assert.Equal(t, "edugo_pass", pc.Password())
}

// =============================================================================
// Tests unitarios para createPostgres/createMongoDB/createRabbitMQ nil config
// =============================================================================

// TestCreatePostgres_NilConfig_Unit verifica error con config nil
func TestCreatePostgres_NilConfig_Unit(t *testing.T) {
	_, err := createPostgres(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "PostgresConfig no puede ser nil")
}

// TestCreateMongoDB_NilConfig_Unit verifica error con config nil
func TestCreateMongoDB_NilConfig_Unit(t *testing.T) {
	_, err := createMongoDB(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MongoConfig no puede ser nil")
}

// TestCreateRabbitMQ_NilConfig_Unit verifica error con config nil
func TestCreateRabbitMQ_NilConfig_Unit(t *testing.T) {
	_, err := createRabbitMQ(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RabbitConfig no puede ser nil")
}

// =============================================================================
// Tests unitarios para MongoDBContainer (sin Docker)
// =============================================================================

// TestMongoDBContainer_Terminate_NilFields verifica Terminate con campos nil
func TestMongoDBContainer_Terminate_NilFields(t *testing.T) {
	mc := &MongoDBContainer{
		config: &MongoConfig{},
		// client y container son nil
	}

	err := mc.Terminate(context.Background())
	assert.NoError(t, err, "Terminate con campos nil no debe dar error")
}

// TestMongoDBContainer_Client_Nil verifica que Client() retorna nil cuando no hay conexión
func TestMongoDBContainer_Client_Nil(t *testing.T) {
	mc := &MongoDBContainer{
		config: &MongoConfig{},
	}

	assert.Nil(t, mc.Client())
}

// =============================================================================
// Tests unitarios para RabbitMQContainer (sin Docker)
// =============================================================================

// TestRabbitMQContainer_Terminate_NilFields verifica Terminate con campos nil
func TestRabbitMQContainer_Terminate_NilFields(t *testing.T) {
	rc := &RabbitMQContainer{
		config: &RabbitConfig{},
		// connection y container son nil
	}

	err := rc.Terminate(context.Background())
	assert.NoError(t, err, "Terminate con campos nil no debe dar error")
}

// TestRabbitMQContainer_Connection_Nil verifica que Connection() retorna nil cuando no hay conexión
func TestRabbitMQContainer_Connection_Nil(t *testing.T) {
	rc := &RabbitMQContainer{
		config: &RabbitConfig{},
	}

	assert.Nil(t, rc.Connection())
}
