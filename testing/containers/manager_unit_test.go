package containers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Tests unitarios para Manager (sin Docker)
// Estos tests verifican la lógica del Manager usando structs construidos
// manualmente, sin crear containers reales.
// =============================================================================

// TestManager_Accessors_AllNil verifica que un Manager vacío retorna nil en todos los accessors
func TestManager_Accessors_AllNil(t *testing.T) {
	m := &Manager{
		config: &Config{},
	}

	assert.Nil(t, m.PostgreSQL(), "PostgreSQL() debe ser nil en manager vacío")
	assert.Nil(t, m.MongoDB(), "MongoDB() debe ser nil en manager vacío")
	assert.Nil(t, m.RabbitMQ(), "RabbitMQ() debe ser nil en manager vacío")
}

// TestManager_Accessors_WithPostgres verifica accessor de PostgreSQL
func TestManager_Accessors_WithPostgres(t *testing.T) {
	pg := &PostgresContainer{
		config: &PostgresConfig{
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		},
	}

	m := &Manager{
		postgres: pg,
		config:   &Config{UsePostgreSQL: true},
	}

	assert.NotNil(t, m.PostgreSQL())
	assert.Same(t, pg, m.PostgreSQL())
	assert.Nil(t, m.MongoDB())
	assert.Nil(t, m.RabbitMQ())
}

// TestManager_Accessors_WithMongoDB verifica accessor de MongoDB
func TestManager_Accessors_WithMongoDB(t *testing.T) {
	mongo := &MongoDBContainer{
		config: &MongoConfig{
			Database: "test_db",
		},
	}

	m := &Manager{
		mongodb: mongo,
		config:  &Config{UseMongoDB: true},
	}

	assert.Nil(t, m.PostgreSQL())
	assert.NotNil(t, m.MongoDB())
	assert.Same(t, mongo, m.MongoDB())
	assert.Nil(t, m.RabbitMQ())
}

// TestManager_Accessors_WithRabbitMQ verifica accessor de RabbitMQ
func TestManager_Accessors_WithRabbitMQ(t *testing.T) {
	rabbit := &RabbitMQContainer{
		config: &RabbitConfig{
			Username: "rabbit_user",
			Password: "rabbit_pass",
		},
	}

	m := &Manager{
		rabbitmq: rabbit,
		config:   &Config{UseRabbitMQ: true},
	}

	assert.Nil(t, m.PostgreSQL())
	assert.Nil(t, m.MongoDB())
	assert.NotNil(t, m.RabbitMQ())
	assert.Same(t, rabbit, m.RabbitMQ())
}

// TestManager_Accessors_AllSet verifica accessors cuando todos los containers están seteados
func TestManager_Accessors_AllSet(t *testing.T) {
	pg := &PostgresContainer{config: &PostgresConfig{Database: "pg"}}
	mongo := &MongoDBContainer{config: &MongoConfig{Database: "mongo"}}
	rabbit := &RabbitMQContainer{config: &RabbitConfig{Username: "rabbit"}}

	m := &Manager{
		postgres: pg,
		mongodb:  mongo,
		rabbitmq: rabbit,
		config: &Config{
			UsePostgreSQL: true,
			UseMongoDB:    true,
			UseRabbitMQ:   true,
		},
	}

	assert.NotNil(t, m.PostgreSQL())
	assert.NotNil(t, m.MongoDB())
	assert.NotNil(t, m.RabbitMQ())
}

// TestManager_CleanPostgreSQL_NotEnabled verifica error cuando PostgreSQL no está habilitado
func TestManager_CleanPostgreSQL_NotEnabled(t *testing.T) {
	m := &Manager{
		config: &Config{},
	}

	err := m.CleanPostgreSQL(context.Background(), "users")
	require.Error(t, err)
	assert.Equal(t, "PostgreSQL no está habilitado", err.Error())
}

// TestManager_CleanMongoDB_NotEnabled verifica error cuando MongoDB no está habilitado
func TestManager_CleanMongoDB_NotEnabled(t *testing.T) {
	m := &Manager{
		config: &Config{},
	}

	err := m.CleanMongoDB(context.Background())
	require.Error(t, err)
	assert.Equal(t, "MongoDB no está habilitado", err.Error())
}

// TestManager_PurgeRabbitMQ_NotEnabled verifica error cuando RabbitMQ no está habilitado
func TestManager_PurgeRabbitMQ_NotEnabled(t *testing.T) {
	m := &Manager{
		config: &Config{},
	}

	err := m.PurgeRabbitMQ(context.Background())
	require.Error(t, err)
	assert.Equal(t, "RabbitMQ no está habilitado", err.Error())
}

// TestManager_Cleanup_AllNilContainers verifica cleanup sin containers
func TestManager_Cleanup_AllNilContainers(t *testing.T) {
	m := &Manager{
		config: &Config{},
	}

	err := m.Cleanup(context.Background())
	assert.NoError(t, err, "Cleanup sin containers no debe dar error")
}

// TestManager_CleanPostgreSQL_EmptyTableList verifica que CleanPostgreSQL con lista vacía no da error
func TestManager_CleanPostgreSQL_EmptyTableList(t *testing.T) {
	pg := &PostgresContainer{
		config: &PostgresConfig{Database: "test_db"},
	}

	m := &Manager{
		postgres: pg,
		config:   &Config{UsePostgreSQL: true},
	}

	// Truncate con lista vacía retorna nil (no hace nada)
	err := m.CleanPostgreSQL(context.Background())
	assert.NoError(t, err, "CleanPostgreSQL con lista vacía no debe dar error")
}

// TestManager_CleanMethods_CombinedNotEnabled verifica mensajes de error para cada servicio
func TestManager_CleanMethods_CombinedNotEnabled(t *testing.T) {
	tests := []struct {
		name        string
		cleanFunc   func(*Manager) error
		expectedMsg string
	}{
		{
			name: "CleanPostgreSQL when disabled",
			cleanFunc: func(m *Manager) error {
				return m.CleanPostgreSQL(context.Background(), "table1")
			},
			expectedMsg: "PostgreSQL no está habilitado",
		},
		{
			name: "CleanMongoDB when disabled",
			cleanFunc: func(m *Manager) error {
				return m.CleanMongoDB(context.Background())
			},
			expectedMsg: "MongoDB no está habilitado",
		},
		{
			name: "PurgeRabbitMQ when disabled",
			cleanFunc: func(m *Manager) error {
				return m.PurgeRabbitMQ(context.Background())
			},
			expectedMsg: "RabbitMQ no está habilitado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{config: &Config{}}
			err := tt.cleanFunc(m)
			require.Error(t, err)
			assert.Equal(t, tt.expectedMsg, err.Error())
		})
	}
}
