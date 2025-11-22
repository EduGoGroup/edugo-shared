package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
