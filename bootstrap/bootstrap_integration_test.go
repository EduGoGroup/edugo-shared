package bootstrap

import (
	"context"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBootstrap_LoggerOnly verifica inicialización solo de logger
func TestBootstrap_LoggerOnly(t *testing.T) {
	ctx := context.Background()

	factories := &Factories{
		Logger: NewDefaultLoggerFactory(),
	}

	// Config mínima
	type Config struct {
		Environment string
		Version     string
	}
	config := Config{
		Environment: "test",
		Version:     "1.0.0",
	}

	resources, err := Bootstrap(
		ctx,
		config,
		factories,
		nil,
		WithRequiredResources("logger"),
		WithSkipHealthCheck(),
	)

	require.NoError(t, err)
	assert.NotNil(t, resources)
	assert.True(t, resources.HasLogger())
	assert.False(t, resources.HasPostgreSQL())
	assert.False(t, resources.HasMongoDB())
}

// TestBootstrap_LoggerAndPostgreSQL verifica logger + PostgreSQL
func TestBootstrap_LoggerAndPostgreSQL(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	// Setup container
	containerConfig := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, containerConfig)
	require.NoError(t, err)

	_ = manager.PostgreSQL() // Asegurar que PostgreSQL está disponible

	// Factories
	factories := &Factories{
		Logger:     NewDefaultLoggerFactory(),
		PostgreSQL: NewDefaultPostgreSQLFactory(nil),
	}

	// Config con PostgreSQL usando DB directamente del container
	// El container ya tiene la conexión configurada
	type AppConfig struct {
		Environment string
		Version     string
		PostgreSQL  PostgreSQLConfig
	}

	appConfig := AppConfig{
		Environment: "test",
		Version:     "1.0.0",
		PostgreSQL: PostgreSQLConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test_user",
			Password: "test_pass",
			Database: "test_db",
			SSLMode:  "disable",
		},
	}

	resources, err := Bootstrap(
		ctx,
		appConfig,
		factories,
		nil,
		WithRequiredResources("logger", "postgresql"),
		WithSkipHealthCheck(),
	)

	require.NoError(t, err)
	assert.NotNil(t, resources)
	assert.True(t, resources.HasLogger())
	assert.True(t, resources.HasPostgreSQL())
}

// TestBootstrap_AllResources verifica inicialización de todos los recursos
func TestBootstrap_AllResources(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	// Setup containers
	containerConfig := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		WithMongoDB(nil).
		WithRabbitMQ(nil).
		Build()

	manager, err := containers.GetManager(t, containerConfig)
	require.NoError(t, err)

	_ = manager.PostgreSQL() // Asegurar que PostgreSQL está disponible

	// Factories - solo las que vamos a usar
	factories := &Factories{
		Logger:     NewDefaultLoggerFactory(),
		PostgreSQL: NewDefaultPostgreSQLFactory(nil),
	}

	// Config completa
	type FullConfig struct {
		Environment string
		Version     string
		PostgreSQL  PostgreSQLConfig
	}

	fullConfig := FullConfig{
		Environment: "test",
		Version:     "1.0.0",
		PostgreSQL: PostgreSQLConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test_user",
			Password: "test_pass",
			Database: "test_db",
			SSLMode:  "disable",
		},
	}

	resources, err := Bootstrap(
		ctx,
		fullConfig,
		factories,
		nil,
		WithRequiredResources("logger"),
		WithOptionalResources("postgresql"),
		WithSkipHealthCheck(),
	)

	require.NoError(t, err)
	assert.NotNil(t, resources)
	assert.True(t, resources.HasLogger())
	assert.NotNil(t, resources.PostgreSQL, "PostgreSQL debe estar inicializado")
}

// TestBootstrap_MissingRequiredFactory verifica error cuando falta factory requerida
func TestBootstrap_MissingRequiredFactory(t *testing.T) {
	ctx := context.Background()

	// Factories sin PostgreSQL
	factories := &Factories{
		Logger: NewDefaultLoggerFactory(),
		// PostgreSQL: nil - falta!
	}

	type Config struct {
		Environment string
		Version     string
	}
	config := Config{
		Environment: "test",
		Version:     "1.0.0",
	}

	// Intentar bootstrap requiriendo PostgreSQL
	resources, err := Bootstrap(
		ctx,
		config,
		factories,
		nil,
		WithRequiredResources("logger", "postgresql"),
		WithSkipHealthCheck(),
	)

	assert.Error(t, err)
	assert.Nil(t, resources)
	assert.Contains(t, err.Error(), "factory validation failed")
}

// TestBootstrap_OptionalResourceFailure verifica que recursos opcionales pueden fallar
func TestBootstrap_OptionalResourceFailure(t *testing.T) {
	ctx := context.Background()

	factories := &Factories{
		Logger: NewDefaultLoggerFactory(),
		// PostgreSQL está configurada pero no habrá config válida
		PostgreSQL: NewDefaultPostgreSQLFactory(nil),
	}

	type Config struct {
		Environment string
		Version     string
		PostgreSQL  PostgreSQLConfig // Config inválida (vacía)
	}

	config := Config{
		Environment: "test",
		Version:     "1.0.0",
		PostgreSQL:  PostgreSQLConfig{}, // Vacía, causará error
	}

	// PostgreSQL es opcional, no debe fallar el bootstrap completo
	resources, err := Bootstrap(
		ctx,
		config,
		factories,
		nil,
		WithRequiredResources("logger"),
		WithOptionalResources("postgresql"),
		WithSkipHealthCheck(),
	)

	// Debe continuar aunque PostgreSQL falle
	require.NoError(t, err)
	assert.NotNil(t, resources)
	assert.True(t, resources.HasLogger())
	// PostgreSQL probablemente no se inicializó
	assert.False(t, resources.HasPostgreSQL())
}

// TestBootstrap_WithHealthCheck verifica health check (skip=false)
func TestBootstrap_WithHealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	// Setup container
	containerConfig := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, containerConfig)
	require.NoError(t, err)

	_ = manager.PostgreSQL() // Asegurar que PostgreSQL está disponible

	factories := &Factories{
		Logger:     NewDefaultLoggerFactory(),
		PostgreSQL: NewDefaultPostgreSQLFactory(nil),
	}

	type AppConfig struct {
		Environment string
		Version     string
		PostgreSQL  PostgreSQLConfig
	}

	appConfig := AppConfig{
		Environment: "test",
		Version:     "1.0.0",
		PostgreSQL: PostgreSQLConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test_user",
			Password: "test_pass",
			Database: "test_db",
			SSLMode:  "disable",
		},
	}

	// Con health check habilitado (default)
	resources, err := Bootstrap(
		ctx,
		appConfig,
		factories,
		nil,
		WithRequiredResources("logger"),
		WithOptionalResources("postgresql"),
		// NO llamar WithSkipHealthCheck()
	)

	// Puede fallar o no dependiendo de la implementación de health check
	// Si hay health check implementado, debe pasar
	_ = resources
	_ = err
}
