package bootstrap

import (
	"testing"
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
