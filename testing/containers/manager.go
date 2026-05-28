package containers

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

// Manager gestiona los containers de testing de forma centralizada.
// Implementa el patrón Singleton para crear los containers una sola vez
// y reutilizarlos entre múltiples tests, mejorando el rendimiento.

// Manager gestiona los containers de testing de forma centralizada
// Manager gestiona los containers de testing de forma centralizada.
// Implementa el patrón Singleton para crear los containers una sola vez
// y reutilizarlos entre múltiples tests, mejorando el rendimiento.
type Manager struct {
	postgres *PostgresContainer
	mongodb  *MongoDBContainer
	rabbitmq *RabbitMQContainer
	config   *Config
	mu       sync.Mutex
}

var (
	globalManager *Manager
	initMu        sync.Mutex
)

// GetManager obtiene o crea el manager global de containers.
// Usa el patrón singleton para crear los containers UNA sola vez.
// Los tests subsiguientes reutilizarán los mismos containers.
func GetManager(t *testing.T, config *Config) (*Manager, error) {
	initMu.Lock()
	defer initMu.Unlock()

	if globalManager != nil {
		return globalManager, nil
	}

	ctx := context.Background()
	m := &Manager{config: config}

	if t != nil {
		t.Logf("🚀 Iniciando containers de testing... Postgres: %v, Mongo: %v, Rabbit: %v", config.UsePostgreSQL, config.UseMongoDB, config.UseRabbitMQ)
	}

	// Crear PostgreSQL si está habilitado
	if config.UsePostgreSQL {
		pg, err := createPostgres(ctx, config.PostgresConfig)
		if err != nil {
			return nil, fmt.Errorf("error creando PostgreSQL: %w", err)
		}
		m.postgres = pg
		if t != nil {
			t.Log("✅ PostgreSQL container listo")
		}
	}

	// Crear MongoDB si está habilitado
	if config.UseMongoDB {
		mongo, err := createMongoDB(ctx, config.MongoConfig)
		if err != nil {
			_ = m.cleanup(ctx, t) //nolint:errcheck // En error de setup, cleanup es best-effort
			return nil, fmt.Errorf("error creando MongoDB: %w", err)
		}
		m.mongodb = mongo
		if t != nil {
			t.Log("✅ MongoDB container listo")
		}
	}

	// Crear RabbitMQ si está habilitado
	if config.UseRabbitMQ {
		rabbit, err := createRabbitMQ(ctx, config.RabbitConfig)
		if err != nil {
			_ = m.cleanup(ctx, t) //nolint:errcheck // En error de setup, cleanup es best-effort
			return nil, fmt.Errorf("error creando RabbitMQ: %w", err)
		}
		m.rabbitmq = rabbit
		if t != nil {
			t.Log("✅ RabbitMQ container listo")
		}
	}

	globalManager = m
	if t != nil {
		t.Log("🎉 Todos los containers listos")
	}

	return globalManager, nil
}

// PostgreSQL retorna el container de PostgreSQL.
// Retorna nil si PostgreSQL no fue habilitado en la Config.
func (m *Manager) PostgreSQL() *PostgresContainer {
	return m.postgres
}

// MongoDB retorna el container de MongoDB.
// Retorna nil si MongoDB no fue habilitado en la Config.
func (m *Manager) MongoDB() *MongoDBContainer {
	return m.mongodb
}

// RabbitMQ retorna el container de RabbitMQ.
// Retorna nil si RabbitMQ no fue habilitado en la Config.
func (m *Manager) RabbitMQ() *RabbitMQContainer {
	return m.rabbitmq
}

// Cleanup limpia todos los containers creados.
// Debe llamarse al final de los tests (típicamente en TestMain).
func (m *Manager) Cleanup(ctx context.Context) error {
	return m.cleanup(ctx, nil)
}

// cleanup es el método interno para limpiar containers
func (m *Manager) cleanup(ctx context.Context, t *testing.T) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	if m.rabbitmq != nil {
		if t != nil {
			t.Log("🧹 Limpiando RabbitMQ...")
		}
		if err := m.rabbitmq.Terminate(ctx); err != nil {
			errors = append(errors, fmt.Errorf("error limpiando RabbitMQ: %w", err))
		}
	}

	if m.mongodb != nil {
		if t != nil {
			t.Log("🧹 Limpiando MongoDB...")
		}
		if err := m.mongodb.Terminate(ctx); err != nil {
			errors = append(errors, fmt.Errorf("error limpiando MongoDB: %w", err))
		}
	}

	if m.postgres != nil {
		if t != nil {
			t.Log("🧹 Limpiando PostgreSQL...")
		}
		if err := m.postgres.Terminate(ctx); err != nil {
			errors = append(errors, fmt.Errorf("error limpiando PostgreSQL: %w", err))
		}
	}

	if len(errors) > 0 {
		errMsg := "errores durante cleanup:"
		for i, err := range errors {
			errMsg += fmt.Sprintf("\n  %d. %v", i+1, err)
		}
		return fmt.Errorf("%s", errMsg)
	}

	if t != nil {
		t.Log("✅ Cleanup completado")
	}
	return nil
}

// CleanPostgreSQL trunca las tablas especificadas de PostgreSQL.
// Útil para limpiar datos entre tests sin recrear el container.
func (m *Manager) CleanPostgreSQL(ctx context.Context, tables ...string) error {
	if m.postgres == nil {
		return fmt.Errorf("PostgreSQL no está habilitado")
	}
	return m.postgres.Truncate(ctx, tables...)
}

// CleanMongoDB elimina todas las colecciones de MongoDB.
// Útil para limpiar datos entre tests sin recrear el container.
func (m *Manager) CleanMongoDB(ctx context.Context) error {
	if m.mongodb == nil {
		return fmt.Errorf("MongoDB no está habilitado")
	}
	return m.mongodb.DropAllCollections(ctx)
}

// PurgeRabbitMQ elimina todas las colas y exchanges de RabbitMQ.
// Útil para limpiar datos entre tests sin recrear el container.
func (m *Manager) PurgeRabbitMQ(ctx context.Context) error {
	if m.rabbitmq == nil {
		return fmt.Errorf("RabbitMQ no está habilitado")
	}
	return m.rabbitmq.PurgeAll(ctx)
}
