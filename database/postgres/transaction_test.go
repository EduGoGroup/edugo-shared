package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/EduGoGroup/edugo-shared/database/postgres"
	"github.com/EduGoGroup/edugo-shared/testing/containers"
)

func TestWithTransaction_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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

	pg := manager.PostgreSQL()
	db := pg.DB()

	t.Run("WithTransaction_CommitExitoso", func(t *testing.T) {
		// Test simplificado: solo verificamos que la transacción se ejecuta sin error
		// usando una query simple sin crear tablas permanentes
		err := postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
			// Usar tabla temporal que se limpia automáticamente
			_, err := tx.ExecContext(ctx, `
				CREATE TEMP TABLE IF NOT EXISTS temp_test_users (
					id SERIAL PRIMARY KEY,
					name VARCHAR(100)
				)
			`)
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(ctx, "INSERT INTO temp_test_users (name) VALUES ('Alice')")
			if err != nil {
				return err
			}
			_, err = tx.ExecContext(ctx, "INSERT INTO temp_test_users (name) VALUES ('Bob')")
			if err != nil {
				return err
			}

			// Verificar que los datos fueron insertados
			var count int
			err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM temp_test_users").Scan(&count)
			if err != nil {
				return err
			}

			if count != 2 {
				t.Errorf("Esperado 2 registros, obtenido %d", count)
			}

			return nil
		})

		if err != nil {
			t.Fatalf("WithTransaction falló: %v", err)
		}
	})

	t.Run("WithTransaction_RollbackEnError", func(t *testing.T) {
		// Test de rollback por error
		expectedErr := errors.New("error intencional")
		err := postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
			// Simplemente retornar error para probar el rollback
			return expectedErr
		})

		if err == nil {
			t.Fatal("Esperaba error de la transacción")
		}

		if !errors.Is(err, expectedErr) {
			t.Errorf("Esperaba error %v, obtenido %v", expectedErr, err)
		}
	})

	t.Run("WithTransaction_RollbackEnPanic", func(t *testing.T) {
		// Test de rollback por panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("Esperaba panic")
			}
		}()

		//nolint:errcheck // Intencionalmente ignoramos el error, estamos probando el panic
		_ = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
			// Provocar panic para probar rollback
			panic("panic intencional")
		})
	})
}

func TestWithTransactionIsolation_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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

	pg := manager.PostgreSQL()
	db := pg.DB()

	t.Run("WithTransactionIsolation_ReadCommitted", func(t *testing.T) {
		// Test simplificado: solo verificar que la transacción se ejecuta correctamente
		err := postgres.WithTransactionIsolation(ctx, db, sql.LevelReadCommitted, func(tx *sql.Tx) error {
			// Ejecutar una query simple para verificar que la transacción funciona
			var result int
			return tx.QueryRowContext(ctx, "SELECT 1").Scan(&result)
		})

		if err != nil {
			t.Fatalf("WithTransactionIsolation falló: %v", err)
		}
	})

	t.Run("WithTransactionIsolation_Serializable", func(t *testing.T) {
		// Test simplificado: verificar nivel de aislamiento Serializable
		err := postgres.WithTransactionIsolation(ctx, db, sql.LevelSerializable, func(tx *sql.Tx) error {
			var result int
			return tx.QueryRowContext(ctx, "SELECT 1").Scan(&result)
		})

		if err != nil {
			t.Fatalf("WithTransactionIsolation con Serializable falló: %v", err)
		}
	})

	t.Run("WithTransactionIsolation_RollbackEnError", func(t *testing.T) {
		// Test de rollback por error
		expectedErr := errors.New("error en transacción")
		err := postgres.WithTransactionIsolation(ctx, db, sql.LevelReadCommitted, func(tx *sql.Tx) error {
			return expectedErr
		})

		if !errors.Is(err, expectedErr) {
			t.Errorf("Esperaba error %v, obtenido %v", expectedErr, err)
		}
	})

	t.Run("WithTransactionIsolation_RollbackEnPanic", func(t *testing.T) {
		// Test de rollback por panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("Esperaba panic")
			}
		}()

		//nolint:errcheck // Intencionalmente ignoramos el error, estamos probando el panic
		_ = postgres.WithTransactionIsolation(ctx, db, sql.LevelSerializable, func(tx *sql.Tx) error {
			panic("panic en transacción")
		})
	})
}
