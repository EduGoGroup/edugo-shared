package containers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewConfig_EmptyConfiguration verifica creación de config vacía
func TestNewConfig_EmptyConfiguration(t *testing.T) {
	builder := NewConfig()
	require.NotNil(t, builder)
	require.NotNil(t, builder.config)

	config := builder.Build()
	assert.NotNil(t, config)
	assert.False(t, config.UsePostgreSQL)
	assert.False(t, config.UseMongoDB)
	assert.False(t, config.UseRabbitMQ)
	assert.Nil(t, config.PostgresConfig)
	assert.Nil(t, config.MongoConfig)
	assert.Nil(t, config.RabbitConfig)
}

// TestConfigBuilder_WithPostgreSQL_Defaults verifica defaults de PostgreSQL
func TestConfigBuilder_WithPostgreSQL_Defaults(t *testing.T) {
	config := NewConfig().
		WithPostgreSQL(nil).
		Build()

	require.NotNil(t, config)
	assert.True(t, config.UsePostgreSQL)
	require.NotNil(t, config.PostgresConfig)

	// Verificar defaults
	assert.Equal(t, "postgres:15-alpine", config.PostgresConfig.Image)
	assert.Equal(t, "edugo_test", config.PostgresConfig.Database)
	assert.Equal(t, "edugo_user", config.PostgresConfig.Username)
	assert.Equal(t, "edugo_pass", config.PostgresConfig.Password)
	assert.Equal(t, "5432", config.PostgresConfig.Port)
	assert.Empty(t, config.PostgresConfig.InitScripts)
}

// TestConfigBuilder_WithPostgreSQL_CustomConfig verifica configuración custom
func TestConfigBuilder_WithPostgreSQL_CustomConfig(t *testing.T) {
	customConfig := &PostgresConfig{
		Image:       "postgres:16-alpine",
		Database:    "custom_db",
		Username:    "custom_user",
		Password:    "custom_pass",
		Port:        "5433",
		InitScripts: []string{"/path/to/init.sql"},
	}

	config := NewConfig().
		WithPostgreSQL(customConfig).
		Build()

	require.NotNil(t, config)
	assert.True(t, config.UsePostgreSQL)
	require.NotNil(t, config.PostgresConfig)

	// Verificar custom values
	assert.Equal(t, "postgres:16-alpine", config.PostgresConfig.Image)
	assert.Equal(t, "custom_db", config.PostgresConfig.Database)
	assert.Equal(t, "custom_user", config.PostgresConfig.Username)
	assert.Equal(t, "custom_pass", config.PostgresConfig.Password)
	assert.Equal(t, "5433", config.PostgresConfig.Port)
	assert.Equal(t, []string{"/path/to/init.sql"}, config.PostgresConfig.InitScripts)
}

// TestConfigBuilder_WithPostgreSQL_PartialConfig verifica defaults parciales
func TestConfigBuilder_WithPostgreSQL_PartialConfig(t *testing.T) {
	partialConfig := &PostgresConfig{
		Database: "partial_db",
		Username: "partial_user",
		// Otros campos vacíos, deben usar defaults
	}

	config := NewConfig().
		WithPostgreSQL(partialConfig).
		Build()

	require.NotNil(t, config)
	assert.True(t, config.UsePostgreSQL)
	require.NotNil(t, config.PostgresConfig)

	// Verificar que valores especificados se mantienen
	assert.Equal(t, "partial_db", config.PostgresConfig.Database)
	assert.Equal(t, "partial_user", config.PostgresConfig.Username)

	// Verificar que valores vacíos reciben defaults
	assert.Equal(t, "postgres:15-alpine", config.PostgresConfig.Image)
	assert.Equal(t, "edugo_pass", config.PostgresConfig.Password)
	assert.Equal(t, "5432", config.PostgresConfig.Port)
}

// TestConfigBuilder_WithMongoDB_Defaults verifica defaults de MongoDB
func TestConfigBuilder_WithMongoDB_Defaults(t *testing.T) {
	config := NewConfig().
		WithMongoDB(nil).
		Build()

	require.NotNil(t, config)
	assert.True(t, config.UseMongoDB)
	require.NotNil(t, config.MongoConfig)

	// Verificar defaults
	assert.Equal(t, "mongo:7.0", config.MongoConfig.Image)
	assert.Equal(t, "edugo_test", config.MongoConfig.Database)
	assert.Empty(t, config.MongoConfig.Username, "MongoDB sin auth por defecto")
	assert.Empty(t, config.MongoConfig.Password, "MongoDB sin auth por defecto")
}

// TestConfigBuilder_WithMongoDB_CustomConfig verifica configuración custom de MongoDB
func TestConfigBuilder_WithMongoDB_CustomConfig(t *testing.T) {
	customConfig := &MongoConfig{
		Image:    "mongo:6.0",
		Database: "custom_mongo_db",
		Username: "mongo_admin",
		Password: "mongo_secret",
	}

	config := NewConfig().
		WithMongoDB(customConfig).
		Build()

	require.NotNil(t, config)
	assert.True(t, config.UseMongoDB)
	require.NotNil(t, config.MongoConfig)

	// Verificar custom values
	assert.Equal(t, "mongo:6.0", config.MongoConfig.Image)
	assert.Equal(t, "custom_mongo_db", config.MongoConfig.Database)
	assert.Equal(t, "mongo_admin", config.MongoConfig.Username)
	assert.Equal(t, "mongo_secret", config.MongoConfig.Password)
}

// TestConfigBuilder_WithRabbitMQ_Defaults verifica defaults de RabbitMQ
func TestConfigBuilder_WithRabbitMQ_Defaults(t *testing.T) {
	config := NewConfig().
		WithRabbitMQ(nil).
		Build()

	require.NotNil(t, config)
	assert.True(t, config.UseRabbitMQ)
	require.NotNil(t, config.RabbitConfig)

	// Verificar defaults
	assert.Equal(t, "rabbitmq:3.12-management-alpine", config.RabbitConfig.Image)
	assert.Equal(t, "edugo_user", config.RabbitConfig.Username)
	assert.Equal(t, "edugo_pass", config.RabbitConfig.Password)
}

// TestConfigBuilder_WithRabbitMQ_CustomConfig verifica configuración custom de RabbitMQ
func TestConfigBuilder_WithRabbitMQ_CustomConfig(t *testing.T) {
	customConfig := &RabbitConfig{
		Image:    "rabbitmq:3.11-alpine",
		Username: "rabbit_admin",
		Password: "rabbit_secret",
	}

	config := NewConfig().
		WithRabbitMQ(customConfig).
		Build()

	require.NotNil(t, config)
	assert.True(t, config.UseRabbitMQ)
	require.NotNil(t, config.RabbitConfig)

	// Verificar custom values
	assert.Equal(t, "rabbitmq:3.11-alpine", config.RabbitConfig.Image)
	assert.Equal(t, "rabbit_admin", config.RabbitConfig.Username)
	assert.Equal(t, "rabbit_secret", config.RabbitConfig.Password)
}

// TestConfigBuilder_ChainedCalls verifica llamadas encadenadas (builder pattern)
func TestConfigBuilder_ChainedCalls(t *testing.T) {
	config := NewConfig().
		WithPostgreSQL(nil).
		WithMongoDB(nil).
		WithRabbitMQ(nil).
		Build()

	require.NotNil(t, config)

	// Todos deben estar habilitados
	assert.True(t, config.UsePostgreSQL)
	assert.True(t, config.UseMongoDB)
	assert.True(t, config.UseRabbitMQ)

	// Todos deben tener sus configs
	assert.NotNil(t, config.PostgresConfig)
	assert.NotNil(t, config.MongoConfig)
	assert.NotNil(t, config.RabbitConfig)
}

// TestConfigBuilder_SelectiveConfiguration verifica configuración selectiva
func TestConfigBuilder_SelectiveConfiguration(t *testing.T) {
	tests := []struct {
		name               string
		build              func(*ConfigBuilder) *Config
		expectPostgres     bool
		expectMongo        bool
		expectRabbit       bool
		expectPgConfig     bool
		expectMongoConfig  bool
		expectRabbitConfig bool
	}{
		{
			name: "only postgres",
			build: func(b *ConfigBuilder) *Config {
				return b.WithPostgreSQL(nil).Build()
			},
			expectPostgres:     true,
			expectMongo:        false,
			expectRabbit:       false,
			expectPgConfig:     true,
			expectMongoConfig:  false,
			expectRabbitConfig: false,
		},
		{
			name: "only mongo",
			build: func(b *ConfigBuilder) *Config {
				return b.WithMongoDB(nil).Build()
			},
			expectPostgres:     false,
			expectMongo:        true,
			expectRabbit:       false,
			expectPgConfig:     false,
			expectMongoConfig:  true,
			expectRabbitConfig: false,
		},
		{
			name: "only rabbit",
			build: func(b *ConfigBuilder) *Config {
				return b.WithRabbitMQ(nil).Build()
			},
			expectPostgres:     false,
			expectMongo:        false,
			expectRabbit:       true,
			expectPgConfig:     false,
			expectMongoConfig:  false,
			expectRabbitConfig: true,
		},
		{
			name: "postgres and mongo",
			build: func(b *ConfigBuilder) *Config {
				return b.WithPostgreSQL(nil).WithMongoDB(nil).Build()
			},
			expectPostgres:     true,
			expectMongo:        true,
			expectRabbit:       false,
			expectPgConfig:     true,
			expectMongoConfig:  true,
			expectRabbitConfig: false,
		},
		{
			name: "mongo and rabbit",
			build: func(b *ConfigBuilder) *Config {
				return b.WithMongoDB(nil).WithRabbitMQ(nil).Build()
			},
			expectPostgres:     false,
			expectMongo:        true,
			expectRabbit:       true,
			expectPgConfig:     false,
			expectMongoConfig:  true,
			expectRabbitConfig: true,
		},
		{
			name: "all three",
			build: func(b *ConfigBuilder) *Config {
				return b.WithPostgreSQL(nil).WithMongoDB(nil).WithRabbitMQ(nil).Build()
			},
			expectPostgres:     true,
			expectMongo:        true,
			expectRabbit:       true,
			expectPgConfig:     true,
			expectMongoConfig:  true,
			expectRabbitConfig: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.build(NewConfig())

			assert.Equal(t, tt.expectPostgres, config.UsePostgreSQL)
			assert.Equal(t, tt.expectMongo, config.UseMongoDB)
			assert.Equal(t, tt.expectRabbit, config.UseRabbitMQ)

			if tt.expectPgConfig {
				assert.NotNil(t, config.PostgresConfig)
			} else {
				assert.Nil(t, config.PostgresConfig)
			}

			if tt.expectMongoConfig {
				assert.NotNil(t, config.MongoConfig)
			} else {
				assert.Nil(t, config.MongoConfig)
			}

			if tt.expectRabbitConfig {
				assert.NotNil(t, config.RabbitConfig)
			} else {
				assert.Nil(t, config.RabbitConfig)
			}
		})
	}
}

// TestConfigBuilder_MultipleCallsSameResource verifica múltiples llamadas al mismo recurso
func TestConfigBuilder_MultipleCallsSameResource(t *testing.T) {
	// Primera configuración
	firstConfig := &PostgresConfig{
		Database: "first_db",
		Username: "first_user",
	}

	// Segunda configuración (debería sobrescribir la primera)
	secondConfig := &PostgresConfig{
		Database: "second_db",
		Username: "second_user",
		Password: "second_pass",
	}

	config := NewConfig().
		WithPostgreSQL(firstConfig).
		WithPostgreSQL(secondConfig).
		Build()

	require.NotNil(t, config)
	assert.True(t, config.UsePostgreSQL)
	require.NotNil(t, config.PostgresConfig)

	// Debe usar la última configuración
	assert.Equal(t, "second_db", config.PostgresConfig.Database)
	assert.Equal(t, "second_user", config.PostgresConfig.Username)
	assert.Equal(t, "second_pass", config.PostgresConfig.Password)
}

// TestConfigBuilder_Immutability verifica que Build() retorna la misma referencia
func TestConfigBuilder_Immutability(t *testing.T) {
	builder := NewConfig().WithPostgreSQL(nil)

	config1 := builder.Build()
	config2 := builder.Build()

	// Deben ser la misma referencia (no una copia)
	assert.Same(t, config1, config2)
}

// TestConfigBuilder_BuildWithoutCalls verifica Build() sin configurar nada
func TestConfigBuilder_BuildWithoutCalls(t *testing.T) {
	config := NewConfig().Build()

	require.NotNil(t, config)
	assert.False(t, config.UsePostgreSQL)
	assert.False(t, config.UseMongoDB)
	assert.False(t, config.UseRabbitMQ)
	assert.Nil(t, config.PostgresConfig)
	assert.Nil(t, config.MongoConfig)
	assert.Nil(t, config.RabbitConfig)
}

// TestConfigBuilder_AllDefaultsCombined verifica todos los defaults combinados
func TestConfigBuilder_AllDefaultsCombined(t *testing.T) {
	config := NewConfig().
		WithPostgreSQL(nil).
		WithMongoDB(nil).
		WithRabbitMQ(nil).
		Build()

	require.NotNil(t, config)

	// PostgreSQL defaults
	assert.Equal(t, "postgres:15-alpine", config.PostgresConfig.Image)
	assert.Equal(t, "edugo_test", config.PostgresConfig.Database)
	assert.Equal(t, "edugo_user", config.PostgresConfig.Username)
	assert.Equal(t, "edugo_pass", config.PostgresConfig.Password)
	assert.Equal(t, "5432", config.PostgresConfig.Port)

	// MongoDB defaults
	assert.Equal(t, "mongo:7.0", config.MongoConfig.Image)
	assert.Equal(t, "edugo_test", config.MongoConfig.Database)
	assert.Empty(t, config.MongoConfig.Username)
	assert.Empty(t, config.MongoConfig.Password)

	// RabbitMQ defaults
	assert.Equal(t, "rabbitmq:3.12-management-alpine", config.RabbitConfig.Image)
	assert.Equal(t, "edugo_user", config.RabbitConfig.Username)
	assert.Equal(t, "edugo_pass", config.RabbitConfig.Password)
}

// TestConfigBuilder_PostgreSQLInitScripts verifica manejo de init scripts
func TestConfigBuilder_PostgreSQLInitScripts(t *testing.T) {
	scripts := []string{
		"/path/to/schema.sql",
		"/path/to/seed.sql",
		"/path/to/fixtures.sql",
	}

	config := NewConfig().
		WithPostgreSQL(&PostgresConfig{
			InitScripts: scripts,
		}).
		Build()

	require.NotNil(t, config)
	require.NotNil(t, config.PostgresConfig)
	assert.Equal(t, scripts, config.PostgresConfig.InitScripts)
}

// TestConfigBuilder_EmptyStringsUseDefaults verifica que strings vacíos usan defaults
func TestConfigBuilder_EmptyStringsUseDefaults(t *testing.T) {
	tests := []struct {
		name   string
		config *PostgresConfig
	}{
		{
			name: "empty image",
			config: &PostgresConfig{
				Image:    "",
				Database: "custom_db",
			},
		},
		{
			name: "empty database",
			config: &PostgresConfig{
				Image:    "postgres:14",
				Database: "",
			},
		},
		{
			name: "all empty",
			config: &PostgresConfig{
				Image:    "",
				Database: "",
				Username: "",
				Password: "",
				Port:     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig().
				WithPostgreSQL(tt.config).
				Build()

			require.NotNil(t, config.PostgresConfig)

			// Verificar que campos vacíos reciben defaults
			if tt.config.Image == "" {
				assert.Equal(t, "postgres:15-alpine", config.PostgresConfig.Image)
			}
			if tt.config.Database == "" {
				assert.Equal(t, "edugo_test", config.PostgresConfig.Database)
			}
			if tt.config.Username == "" {
				assert.Equal(t, "edugo_user", config.PostgresConfig.Username)
			}
			if tt.config.Password == "" {
				assert.Equal(t, "edugo_pass", config.PostgresConfig.Password)
			}
			if tt.config.Port == "" {
				assert.Equal(t, "5432", config.PostgresConfig.Port)
			}
		})
	}
}

// TestConfigBuilder_ComplexScenario verifica escenario complejo combinado
func TestConfigBuilder_ComplexScenario(t *testing.T) {
	// Escenario: PostgreSQL custom, MongoDB default, sin RabbitMQ
	config := NewConfig().
		WithPostgreSQL(&PostgresConfig{
			Image:    "postgres:16-alpine",
			Database: "production_db",
			Username: "prod_user",
			Password: "super_secret",
			Port:     "5433",
			InitScripts: []string{
				"/migrations/001_initial.sql",
				"/migrations/002_add_users.sql",
			},
		}).
		WithMongoDB(nil). // Usar defaults
		Build()

	require.NotNil(t, config)

	// PostgreSQL custom
	assert.True(t, config.UsePostgreSQL)
	assert.Equal(t, "postgres:16-alpine", config.PostgresConfig.Image)
	assert.Equal(t, "production_db", config.PostgresConfig.Database)
	assert.Equal(t, "prod_user", config.PostgresConfig.Username)
	assert.Equal(t, "super_secret", config.PostgresConfig.Password)
	assert.Equal(t, "5433", config.PostgresConfig.Port)
	assert.Len(t, config.PostgresConfig.InitScripts, 2)

	// MongoDB defaults
	assert.True(t, config.UseMongoDB)
	assert.Equal(t, "mongo:7.0", config.MongoConfig.Image)
	assert.Equal(t, "edugo_test", config.MongoConfig.Database)

	// RabbitMQ no habilitado
	assert.False(t, config.UseRabbitMQ)
	assert.Nil(t, config.RabbitConfig)
}

// TestConfigBuilder_NilSafety verifica que nil configs son manejados correctamente
func TestConfigBuilder_NilSafety(t *testing.T) {
	// Todas las llamadas con nil
	config := NewConfig().
		WithPostgreSQL(nil).
		WithMongoDB(nil).
		WithRabbitMQ(nil).
		Build()

	require.NotNil(t, config)

	// Todos habilitados con defaults
	assert.True(t, config.UsePostgreSQL)
	assert.True(t, config.UseMongoDB)
	assert.True(t, config.UseRabbitMQ)

	// Todas las configs deben estar presentes (no nil)
	assert.NotNil(t, config.PostgresConfig)
	assert.NotNil(t, config.MongoConfig)
	assert.NotNil(t, config.RabbitConfig)

	// Todas con defaults correctos
	assert.NotEmpty(t, config.PostgresConfig.Image)
	assert.NotEmpty(t, config.MongoConfig.Image)
	assert.NotEmpty(t, config.RabbitConfig.Image)
}

// TestConfigBuilder_MethodChaining verifica que los métodos retornan el builder
func TestConfigBuilder_MethodChaining(t *testing.T) {
	builder := NewConfig()

	// Verificar que cada método retorna *ConfigBuilder
	pgBuilder := builder.WithPostgreSQL(nil)
	assert.Same(t, builder, pgBuilder, "WithPostgreSQL debe retornar el mismo builder")

	mongoBuilder := builder.WithMongoDB(nil)
	assert.Same(t, builder, mongoBuilder, "WithMongoDB debe retornar el mismo builder")

	rabbitBuilder := builder.WithRabbitMQ(nil)
	assert.Same(t, builder, rabbitBuilder, "WithRabbitMQ debe retornar el mismo builder")
}
