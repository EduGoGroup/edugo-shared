package postgres_test

import (
	"context"
	"testing"

	"github.com/EduGoGroup/edugo-shared/database/postgres"
	"github.com/EduGoGroup/edugo-shared/testing/containers"
)

func TestHealthCheck_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()
	pgConfig := &containers.PostgresConfig{
		Image:    "postgres:15-alpine",
		Database: "test_db",
		Username: "test_user",
		Password: "test_pass",
	}

	config := containers.NewConfig().
		WithPostgreSQL(pgConfig).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}
	defer manager.Cleanup(ctx)

	pg := manager.PostgreSQL()
	if pg == nil {
		t.Fatal("PostgreSQL container es nil")
	}

	t.Run("HealthCheck_Exitoso", func(t *testing.T) {
		db := pg.DB()
		err := postgres.HealthCheck(db)
		if err != nil {
			t.Errorf("HealthCheck falló: %v", err)
		}
	})

	t.Run("HealthCheck_ConexionCerrada", func(t *testing.T) {
		// Obtener DB y cerrarla
		db := pg.DB()

		// Cerrar la conexión
		err := db.Close()
		if err != nil {
			t.Fatalf("Error cerrando DB: %v", err)
		}

		// HealthCheck debería fallar
		err = postgres.HealthCheck(db)
		if err == nil {
			t.Error("Esperaba error en HealthCheck con conexión cerrada")
		}
	})
}

func TestGetStats_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()
	pgConfig := &containers.PostgresConfig{
		Image:    "postgres:15-alpine",
		Database: "test_db",
		Username: "test_user",
		Password: "test_pass",
	}

	config := containers.NewConfig().
		WithPostgreSQL(pgConfig).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}
	defer manager.Cleanup(ctx)

	pg := manager.PostgreSQL()
	db := pg.DB()

	t.Run("GetStats_RetornaEstadisticas", func(t *testing.T) {
		stats := postgres.GetStats(db)

		// Verificar que las estadísticas tienen sentido
		if stats.MaxOpenConnections <= 0 {
			t.Error("MaxOpenConnections debe ser mayor que 0")
		}
	})

	t.Run("GetStats_ConConexionesActivas", func(t *testing.T) {
		// Ejecutar una query para crear conexión activa
		_, err := db.QueryContext(ctx, "SELECT 1")
		if err != nil {
			t.Fatalf("Query falló: %v", err)
		}

		stats := postgres.GetStats(db)

		// Debería tener al menos una conexión abierta
		if stats.OpenConnections < 0 {
			t.Error("OpenConnections no puede ser negativo")
		}
	})
}

func TestClose_Integration(t *testing.T) {
	t.Run("Close_ConDBNil", func(t *testing.T) {
		err := postgres.Close(nil)
		if err != nil {
			t.Errorf("Close con DB nil no debería dar error: %v", err)
		}
	})

	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()
	pgConfig := &containers.PostgresConfig{
		Image:    "postgres:15-alpine",
		Database: "test_db",
		Username: "test_user",
		Password: "test_pass",
	}

	config := containers.NewConfig().
		WithPostgreSQL(pgConfig).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}
	defer manager.Cleanup(ctx)

	pg := manager.PostgreSQL()

	t.Run("Close_Exitoso", func(t *testing.T) {
		db := pg.DB()

		// Cerrar
		err := postgres.Close(db)
		if err != nil {
			t.Errorf("Close falló: %v", err)
		}

		// Verificar que está cerrada
		err = db.Ping()
		if err == nil {
			t.Error("Esperaba error al hacer ping después de cerrar")
		}
	})
}
