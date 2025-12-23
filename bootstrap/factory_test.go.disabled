package bootstrap

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests simples para verificar que las factories se pueden crear correctamente

func TestNewDefaultLoggerFactory(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	if factory == nil {
		t.Error("NewDefaultLoggerFactory retornó nil")
	}
}

func TestNewDefaultMongoDBFactory(t *testing.T) {
	factory := NewDefaultMongoDBFactory()
	if factory == nil {
		t.Error("NewDefaultMongoDBFactory retornó nil")
	}
}

func TestNewDefaultPostgreSQLFactory(t *testing.T) {
	// PostgreSQLFactory requiere un logger interface
	factory := NewDefaultPostgreSQLFactory(nil)
	if factory == nil {
		t.Error("NewDefaultPostgreSQLFactory retornó nil")
	}
}

func TestNewDefaultRabbitMQFactory(t *testing.T) {
	factory := NewDefaultRabbitMQFactory()
	if factory == nil {
		t.Error("NewDefaultRabbitMQFactory retornó nil")
	}
}

func TestNewDefaultS3Factory(t *testing.T) {
	factory := NewDefaultS3Factory()
	if factory == nil {
		t.Error("NewDefaultS3Factory retornó nil")
	}
}

// Tests para verificar que las factories implementan sus interfaces correctamente

func TestFactoriesImplementInterfaces(t *testing.T) {
	t.Run("LoggerFactory_implementa_interfaz", func(t *testing.T) {
		var _ LoggerFactory = NewDefaultLoggerFactory()
	})

	t.Run("MongoDBFactory_implementa_interfaz", func(t *testing.T) {
		var _ MongoDBFactory = NewDefaultMongoDBFactory()
	})

	t.Run("PostgreSQLFactory_implementa_interfaz", func(t *testing.T) {
		var _ PostgreSQLFactory = NewDefaultPostgreSQLFactory(nil)
	})

	t.Run("RabbitMQFactory_implementa_interfaz", func(t *testing.T) {
		var _ RabbitMQFactory = NewDefaultRabbitMQFactory()
	})

	t.Run("S3Factory_implementa_interfaz", func(t *testing.T) {
		var _ S3Factory = NewDefaultS3Factory()
	})
}

// =============================================================================
// LOGGER FACTORY TESTS
// =============================================================================

func TestLoggerFactory_CreateLogger_ProductionEnv(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	logger, err := factory.CreateLogger(ctx, "production", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Verificar que el logger está configurado para producción
	assert.Equal(t, logrus.InfoLevel, logger.Level)

	// Verificar que usa JSON formatter
	_, ok := logger.Formatter.(*logrus.JSONFormatter)
	assert.True(t, ok, "El logger debe usar JSONFormatter en producción")
}

func TestLoggerFactory_CreateLogger_ProdEnv(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	logger, err := factory.CreateLogger(ctx, "prod", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, logger)

	assert.Equal(t, logrus.InfoLevel, logger.Level)

	_, ok := logger.Formatter.(*logrus.JSONFormatter)
	assert.True(t, ok, "El logger debe usar JSONFormatter en prod")
}

func TestLoggerFactory_CreateLogger_DevelopmentEnv(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	logger, err := factory.CreateLogger(ctx, "development", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Verificar que el logger está configurado para desarrollo
	assert.Equal(t, logrus.DebugLevel, logger.Level)

	// Verificar que usa Text formatter
	_, ok := logger.Formatter.(*logrus.TextFormatter)
	assert.True(t, ok, "El logger debe usar TextFormatter en desarrollo")
}

func TestLoggerFactory_CreateLogger_DevEnv(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	logger, err := factory.CreateLogger(ctx, "dev", "2.0.0")
	require.NoError(t, err)
	require.NotNil(t, logger)

	assert.Equal(t, logrus.DebugLevel, logger.Level)
}

func TestLoggerFactory_CreateLogger_LocalEnv(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	logger, err := factory.CreateLogger(ctx, "local", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Local debe usar el nivel más detallado
	assert.Equal(t, logrus.TraceLevel, logger.Level)
}

func TestLoggerFactory_CreateLogger_QAEnv(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	logger, err := factory.CreateLogger(ctx, "qa", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, logger)

	assert.Equal(t, logrus.InfoLevel, logger.Level)
}

func TestLoggerFactory_CreateLogger_StagingEnv(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	logger, err := factory.CreateLogger(ctx, "staging", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, logger)

	assert.Equal(t, logrus.InfoLevel, logger.Level)
}

func TestLoggerFactory_CreateLogger_UnknownEnv(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	logger, err := factory.CreateLogger(ctx, "unknown", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Entorno desconocido debe usar nivel Info por defecto
	assert.Equal(t, logrus.InfoLevel, logger.Level)
}

func TestLoggerFactory_CreateLogger_EmptyEnv(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	logger, err := factory.CreateLogger(ctx, "", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, logger)

	assert.Equal(t, logrus.InfoLevel, logger.Level)
}

func TestLoggerFactory_CreateLogger_WithDifferentVersions(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	versions := []string{"1.0.0", "2.1.3", "v1.0.0", "latest", "dev-123"}

	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			logger, err := factory.CreateLogger(ctx, "production", version)
			require.NoError(t, err)
			require.NotNil(t, logger)
		})
	}
}

// =============================================================================
// MONGODB FACTORY TESTS
// =============================================================================

// TestMongoDBFactory_DefaultTimeout movido a factory_mongodb_integration_test.go

// =============================================================================
// POSTGRESQL FACTORY TESTS
// =============================================================================

// TestPostgreSQLFactory_WithCustomLogger movido a factory_postgresql_integration_test.go

func TestPostgreSQLFactory_BuildDSN(t *testing.T) {
	factory := NewDefaultPostgreSQLFactory(nil)
	require.NotNil(t, factory)

	tests := []struct {
		name     string
		config   PostgreSQLConfig
		expected string
	}{
		{
			name: "basic config",
			config: PostgreSQLConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "testpass",
				Database: "testdb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable",
		},
		{
			name: "config without sslmode",
			config: PostgreSQLConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "pass",
				Database: "db",
			},
			expected: "host=localhost port=5432 user=user password=pass dbname=db sslmode=disable",
		},
		{
			name: "config with require sslmode",
			config: PostgreSQLConfig{
				Host:     "remote.db.com",
				Port:     5432,
				User:     "admin",
				Password: "secret",
				Database: "production",
				SSLMode:  "require",
			},
			expected: "host=remote.db.com port=5432 user=admin password=secret dbname=production sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := factory.buildDSN(tt.config)
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

// =============================================================================
// RABBITMQ FACTORY TESTS
// =============================================================================

func TestRabbitMQFactory_Creation(t *testing.T) {
	factory := NewDefaultRabbitMQFactory()
	require.NotNil(t, factory)

	// RabbitMQFactory debe ser creada correctamente
	assert.IsType(t, &DefaultRabbitMQFactory{}, factory)
}

// =============================================================================
// S3 FACTORY TESTS
// =============================================================================

// TestS3Factory_Creation movido a factory_s3_test.go

// =============================================================================
// INTEGRATION TESTS
// =============================================================================

func TestMultipleFactories_CanBeCreatedSimultaneously(t *testing.T) {
	loggerFactory := NewDefaultLoggerFactory()
	mongoFactory := NewDefaultMongoDBFactory()
	pgFactory := NewDefaultPostgreSQLFactory(nil)
	rabbitFactory := NewDefaultRabbitMQFactory()
	s3Factory := NewDefaultS3Factory()

	assert.NotNil(t, loggerFactory)
	assert.NotNil(t, mongoFactory)
	assert.NotNil(t, pgFactory)
	assert.NotNil(t, rabbitFactory)
	assert.NotNil(t, s3Factory)
}

func TestLoggerFactory_MultipleLoggers(t *testing.T) {
	factory := NewDefaultLoggerFactory()
	ctx := context.Background()

	// Crear múltiples loggers con diferentes configuraciones
	logger1, err1 := factory.CreateLogger(ctx, "production", "1.0.0")
	logger2, err2 := factory.CreateLogger(ctx, "development", "1.0.0")
	logger3, err3 := factory.CreateLogger(ctx, "local", "2.0.0")

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)

	assert.NotNil(t, logger1)
	assert.NotNil(t, logger2)
	assert.NotNil(t, logger3)

	// Verificar que tienen configuraciones diferentes
	assert.Equal(t, logrus.InfoLevel, logger1.Level)
	assert.Equal(t, logrus.DebugLevel, logger2.Level)
	assert.Equal(t, logrus.TraceLevel, logger3.Level)
}
