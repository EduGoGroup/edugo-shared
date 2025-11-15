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

	// Crear tabla de prueba
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100)
		)
	`)
	if err != nil {
		t.Fatalf("Error creando tabla: %v", err)
	}

	defer func() {
		// Limpiar al final
		_, _ = db.ExecContext(ctx, "DROP TABLE IF EXISTS test_users")
	}()

	t.Run("WithTransaction_CommitExitoso", func(t *testing.T) {
		// Limpiar tabla
		_, err := db.ExecContext(ctx, "TRUNCATE TABLE test_users")
		if err != nil {
			t.Fatalf("Error truncando tabla: %v", err)
		}

		// Ejecutar transacción exitosa
		err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test_users (name) VALUES ('Alice')")
			if err != nil {
				return err
			}
			_, err = tx.ExecContext(ctx, "INSERT INTO test_users (name) VALUES ('Bob')")
			return err
		})

		if err != nil {
			t.Fatalf("WithTransaction falló: %v", err)
		}

		// Verificar que los datos fueron insertados
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_users").Scan(&count)
		if err != nil {
			t.Fatalf("Error contando registros: %v", err)
		}

		if count != 2 {
			t.Errorf("Esperado 2 registros, obtenido %d", count)
		}
	})

	t.Run("WithTransaction_RollbackEnError", func(t *testing.T) {
		// Limpiar tabla
		_, err := db.ExecContext(ctx, "TRUNCATE TABLE test_users")
		if err != nil {
			t.Fatalf("Error truncando tabla: %v", err)
		}

		// Ejecutar transacción que falla
		expectedErr := errors.New("error intencional")
		err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test_users (name) VALUES ('Charlie')")
			if err != nil {
				return err
			}
			// Retornar error para forzar rollback
			return expectedErr
		})

		if err == nil {
			t.Fatal("Esperaba error de la transacción")
		}

		if err != expectedErr {
			t.Errorf("Esperaba error %v, obtenido %v", expectedErr, err)
		}

		// Verificar que NO se insertaron datos (rollback exitoso)
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_users").Scan(&count)
		if err != nil {
			t.Fatalf("Error contando registros: %v", err)
		}

		if count != 0 {
			t.Errorf("Esperado 0 registros después de rollback, obtenido %d", count)
		}
	})

	t.Run("WithTransaction_RollbackEnPanic", func(t *testing.T) {
		// Limpiar tabla
		_, err := db.ExecContext(ctx, "TRUNCATE TABLE test_users")
		if err != nil {
			t.Fatalf("Error truncando tabla: %v", err)
		}

		// Ejecutar transacción que hace panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("Esperaba panic")
			}
		}()

		_ = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test_users (name) VALUES ('Dave')")
			if err != nil {
				return err
			}
			// Provocar panic
			panic("panic intencional")
		})

		// Verificar que NO se insertaron datos (rollback en panic)
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_users").Scan(&count)
		if err != nil {
			t.Fatalf("Error contando registros: %v", err)
		}

		if count != 0 {
			t.Errorf("Esperado 0 registros después de rollback por panic, obtenido %d", count)
		}
	})
}

func TestWithTransactionIsolation_Integration(t *testing.T) {
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

	// Crear tabla de prueba
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_accounts (
			id SERIAL PRIMARY KEY,
			balance INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("Error creando tabla: %v", err)
	}

	defer func() {
		_, _ = db.ExecContext(ctx, "DROP TABLE IF EXISTS test_accounts")
	}()

	t.Run("WithTransactionIsolation_ReadCommitted", func(t *testing.T) {
		// Limpiar tabla
		_, err := db.ExecContext(ctx, "TRUNCATE TABLE test_accounts")
		if err != nil {
			t.Fatalf("Error truncando tabla: %v", err)
		}

		// Ejecutar transacción con nivel de aislamiento ReadCommitted
		err = postgres.WithTransactionIsolation(ctx, db, sql.LevelReadCommitted, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test_accounts (balance) VALUES (1000)")
			return err
		})

		if err != nil {
			t.Fatalf("WithTransactionIsolation falló: %v", err)
		}

		// Verificar que los datos fueron insertados
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_accounts").Scan(&count)
		if err != nil {
			t.Fatalf("Error contando registros: %v", err)
		}

		if count != 1 {
			t.Errorf("Esperado 1 registro, obtenido %d", count)
		}
	})

	t.Run("WithTransactionIsolation_Serializable", func(t *testing.T) {
		// Limpiar tabla
		_, err := db.ExecContext(ctx, "TRUNCATE TABLE test_accounts")
		if err != nil {
			t.Fatalf("Error truncando tabla: %v", err)
		}

		// Ejecutar transacción con nivel de aislamiento Serializable
		err = postgres.WithTransactionIsolation(ctx, db, sql.LevelSerializable, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test_accounts (balance) VALUES (2000)")
			return err
		})

		if err != nil {
			t.Fatalf("WithTransactionIsolation con Serializable falló: %v", err)
		}

		// Verificar que los datos fueron insertados
		var balance int
		err = db.QueryRowContext(ctx, "SELECT balance FROM test_accounts LIMIT 1").Scan(&balance)
		if err != nil {
			t.Fatalf("Error leyendo balance: %v", err)
		}

		if balance != 2000 {
			t.Errorf("Esperado balance 2000, obtenido %d", balance)
		}
	})

	t.Run("WithTransactionIsolation_RollbackEnError", func(t *testing.T) {
		// Limpiar tabla
		_, err := db.ExecContext(ctx, "TRUNCATE TABLE test_accounts")
		if err != nil {
			t.Fatalf("Error truncando tabla: %v", err)
		}

		// Ejecutar transacción que falla
		expectedErr := errors.New("error en transacción")
		err = postgres.WithTransactionIsolation(ctx, db, sql.LevelReadCommitted, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test_accounts (balance) VALUES (3000)")
			if err != nil {
				return err
			}
			return expectedErr
		})

		if err != expectedErr {
			t.Errorf("Esperaba error %v, obtenido %v", expectedErr, err)
		}

		// Verificar que NO se insertaron datos (rollback)
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_accounts").Scan(&count)
		if err != nil {
			t.Fatalf("Error contando registros: %v", err)
		}

		if count != 0 {
			t.Errorf("Esperado 0 registros después de rollback, obtenido %d", count)
		}
	})

	t.Run("WithTransactionIsolation_RollbackEnPanic", func(t *testing.T) {
		// Limpiar tabla
		_, err := db.ExecContext(ctx, "TRUNCATE TABLE test_accounts")
		if err != nil {
			t.Fatalf("Error truncando tabla: %v", err)
		}

		// Ejecutar transacción que hace panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("Esperaba panic")
			}
		}()

		_ = postgres.WithTransactionIsolation(ctx, db, sql.LevelSerializable, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test_accounts (balance) VALUES (4000)")
			if err != nil {
				return err
			}
			panic("panic en transacción")
		})

		// Verificar que NO se insertaron datos (rollback por panic)
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_accounts").Scan(&count)
		if err != nil {
			t.Fatalf("Error contando registros: %v", err)
		}

		if count != 0 {
			t.Errorf("Esperado 0 registros después de rollback por panic, obtenido %d", count)
		}
	})
}
