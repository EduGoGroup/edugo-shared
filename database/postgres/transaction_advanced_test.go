package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/database/postgres"
	"github.com/EduGoGroup/edugo-shared/testing/containers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWithTransaction_MultipleConcurrent verifica múltiples transacciones concurrentes
func TestWithTransaction_MultipleConcurrent(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Crear tabla de test
	_, err = db.ExecContext(ctx, `
		CREATE TEMP TABLE concurrent_test (
			id SERIAL PRIMARY KEY,
			value INT
		)
	`)
	require.NoError(t, err)

	// Ejecutar múltiples transacciones concurrentes
	var wg sync.WaitGroup
	concurrency := 10

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()

			err := postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx, "INSERT INTO concurrent_test (value) VALUES ($1)", val)
				return err
			})
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Verificar que todas las inserciones se hicieron
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM concurrent_test").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, concurrency, count)
}

// TestWithTransaction_MultipleTablesOperation verifica operaciones en múltiples tablas
func TestWithTransaction_MultipleTablesOperation(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Crear tablas relacionadas
	_, err = db.ExecContext(ctx, `
		CREATE TEMP TABLE users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100)
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TEMP TABLE orders (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			amount DECIMAL(10,2)
		)
	`)
	require.NoError(t, err)

	// Transacción que modifica ambas tablas
	err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		// Insertar usuario
		var userID int
		err := tx.QueryRowContext(ctx, "INSERT INTO users (name) VALUES ($1) RETURNING id", "John Doe").Scan(&userID)
		if err != nil {
			return err
		}

		// Insertar orden para ese usuario
		_, err = tx.ExecContext(ctx, "INSERT INTO orders (user_id, amount) VALUES ($1, $2)", userID, 99.99)
		return err
	})

	require.NoError(t, err)

	// Verificar datos
	var userCount, orderCount int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders").Scan(&orderCount)

	assert.Equal(t, 1, userCount)
	assert.Equal(t, 1, orderCount)
}

// TestWithTransaction_ForeignKeyConstraints verifica manejo de constraints
func TestWithTransaction_ForeignKeyConstraints(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Crear tablas con foreign key
	_, err = db.ExecContext(ctx, `
		CREATE TEMP TABLE parent (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100)
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TEMP TABLE child (
			id SERIAL PRIMARY KEY,
			parent_id INT REFERENCES parent(id) ON DELETE CASCADE,
			value INT
		)
	`)
	require.NoError(t, err)

	// Intentar insertar child sin parent (debe fallar)
	err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		// Intentar insertar con parent_id que no existe
		_, err := tx.ExecContext(ctx, "INSERT INTO child (parent_id, value) VALUES ($1, $2)", 999, 100)
		return err
	})

	assert.Error(t, err, "Debe fallar por foreign key constraint")
}

// TestWithTransaction_RollbackOnCommitError verifica rollback cuando commit falla
func TestWithTransaction_RollbackOnCommitError(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Crear tabla
	_, err = db.ExecContext(ctx, `
		CREATE TEMP TABLE rollback_test (
			id SERIAL PRIMARY KEY,
			value INT UNIQUE
		)
	`)
	require.NoError(t, err)

	// Insertar valor inicial
	_, err = db.ExecContext(ctx, "INSERT INTO rollback_test (value) VALUES (1)")
	require.NoError(t, err)

	// Transacción que viola constraint al final
	err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		// Esta inserción es válida dentro de la transacción
		_, err := tx.ExecContext(ctx, "INSERT INTO rollback_test (value) VALUES (2)")
		if err != nil {
			return err
		}

		// Esta viola constraint UNIQUE (debe fallar al commit)
		_, err = tx.ExecContext(ctx, "INSERT INTO rollback_test (value) VALUES (1)")
		return err
	})

	assert.Error(t, err, "Debe fallar por UNIQUE constraint")

	// Verificar que se hizo rollback (solo debe haber 1 registro)
	var count int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM rollback_test").Scan(&count)
	assert.Equal(t, 1, count)
}

// TestWithTransactionIsolation_AllLevels verifica todos los niveles de aislamiento
func TestWithTransactionIsolation_AllLevels(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	levels := []struct {
		name  string
		level sql.IsolationLevel
	}{
		{"Default", sql.LevelDefault},
		{"ReadUncommitted", sql.LevelReadUncommitted},
		{"ReadCommitted", sql.LevelReadCommitted},
		{"RepeatableRead", sql.LevelRepeatableRead},
		{"Serializable", sql.LevelSerializable},
	}

	for _, tt := range levels {
		t.Run(tt.name, func(t *testing.T) {
			err := postgres.WithTransactionIsolation(ctx, db, tt.level, func(tx *sql.Tx) error {
				var result int
				return tx.QueryRowContext(ctx, "SELECT 1").Scan(&result)
			})
			assert.NoError(t, err)
		})
	}
}

// TestWithTransaction_ContextCancellation verifica cancelación por contexto
func TestWithTransaction_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Contexto con cancelación inmediata
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		// Nunca debería ejecutarse
		time.Sleep(1 * time.Second)
		return nil
	})

	assert.Error(t, err, "Debe fallar con contexto cancelado")
}

// TestWithTransaction_Timeout verifica timeout durante transacción
func TestWithTransaction_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Contexto con timeout muy corto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		// Operación larga que excede el timeout
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.Error(t, err, "Debe fallar por timeout")
}

// TestWithTransaction_DBClosed verifica error cuando DB está cerrada
func TestWithTransaction_DBClosed(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Cerrar la DB
	db.Close()

	// Intentar transacción con DB cerrada
	err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		return nil
	})

	assert.Error(t, err, "Debe fallar con DB cerrada")
}

// TestWithTransaction_NestedTransactionSimulation verifica savepoints (simulación)
func TestWithTransaction_NestedTransactionSimulation(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Crear tabla
	_, err = db.ExecContext(ctx, `
		CREATE TEMP TABLE savepoint_test (
			id SERIAL PRIMARY KEY,
			value INT
		)
	`)
	require.NoError(t, err)

	// Transacción con savepoint
	err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		// Insertar valor 1
		_, err := tx.ExecContext(ctx, "INSERT INTO savepoint_test (value) VALUES (1)")
		if err != nil {
			return err
		}

		// Crear savepoint
		_, err = tx.ExecContext(ctx, "SAVEPOINT sp1")
		if err != nil {
			return err
		}

		// Insertar valor 2
		_, err = tx.ExecContext(ctx, "INSERT INTO savepoint_test (value) VALUES (2)")
		if err != nil {
			return err
		}

		// Rollback al savepoint
		_, err = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT sp1")
		if err != nil {
			return err
		}

		// Insertar valor 3
		_, err = tx.ExecContext(ctx, "INSERT INTO savepoint_test (value) VALUES (3)")
		return err
	})

	require.NoError(t, err)

	// Verificar que solo hay 2 valores (1 y 3, no 2)
	var count int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM savepoint_test").Scan(&count)
	assert.Equal(t, 2, count)

	var values []int
	rows, _ := db.QueryContext(ctx, "SELECT value FROM savepoint_test ORDER BY value")
	defer rows.Close()
	for rows.Next() {
		var v int
		rows.Scan(&v)
		values = append(values, v)
	}
	assert.Equal(t, []int{1, 3}, values)
}

// TestWithTransaction_LargeDataset verifica transacción con dataset grande
func TestWithTransaction_LargeDataset(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Crear tabla
	_, err = db.ExecContext(ctx, `
		CREATE TEMP TABLE large_dataset (
			id SERIAL PRIMARY KEY,
			value INT
		)
	`)
	require.NoError(t, err)

	// Insertar 1000 registros en una transacción
	err = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, "INSERT INTO large_dataset (value) VALUES ($1)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for i := 0; i < 1000; i++ {
			_, err := stmt.ExecContext(ctx, i)
			if err != nil {
				return err
			}
		}
		return nil
	})

	require.NoError(t, err)

	// Verificar count
	var count int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM large_dataset").Scan(&count)
	assert.Equal(t, 1000, count)
}

// TestWithTransaction_ErrorTypes verifica diferentes tipos de errores
func TestWithTransaction_ErrorTypes(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	tests := []struct {
		name        string
		txFunc      postgres.TxFunc
		expectError bool
	}{
		{
			name: "syntax error",
			txFunc: func(tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx, "INVALID SQL SYNTAX")
				return err
			},
			expectError: true,
		},
		{
			name: "table not found",
			txFunc: func(tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx, "SELECT * FROM non_existent_table")
				return err
			},
			expectError: true,
		},
		{
			name: "custom error",
			txFunc: func(tx *sql.Tx) error {
				return errors.New("custom application error")
			},
			expectError: true,
		},
		{
			name: "success",
			txFunc: func(tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx, "SELECT 1")
				return err
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := postgres.WithTransaction(ctx, db, tt.txFunc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestWithTransaction_PanicRecovery verifica recuperación de panic con mensaje
func TestWithTransaction_PanicRecovery_WithMessage(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "test panic with specific message", r)
		} else {
			t.Error("Expected panic but didn't happen")
		}
	}()

	_ = postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		panic("test panic with specific message")
	})
}

// TestWithTransactionIsolation_Conflicts verifica conflictos de serialización
func TestWithTransactionIsolation_Conflicts(t *testing.T) {
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
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	db := pg.DB()

	// Crear tabla
	_, err = db.ExecContext(ctx, `
		CREATE TEMP TABLE isolation_test (
			id SERIAL PRIMARY KEY,
			counter INT DEFAULT 0
		)
	`)
	require.NoError(t, err)

	// Insertar registro inicial
	_, err = db.ExecContext(ctx, "INSERT INTO isolation_test (counter) VALUES (0)")
	require.NoError(t, err)

	// Dos transacciones concurrentes con nivel Serializable
	var wg sync.WaitGroup
	errors := make([]error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			errors[index] = postgres.WithTransactionIsolation(ctx, db, sql.LevelSerializable, func(tx *sql.Tx) error {
				// Leer valor actual
				var counter int
				err := tx.QueryRowContext(ctx, "SELECT counter FROM isolation_test WHERE id = 1").Scan(&counter)
				if err != nil {
					return err
				}

				// Simular procesamiento
				time.Sleep(10 * time.Millisecond)

				// Actualizar
				_, err = tx.ExecContext(ctx, "UPDATE isolation_test SET counter = $1 WHERE id = 1", counter+1)
				return err
			})
		}(i)
	}

	wg.Wait()

	// Al menos una debería tener error de serialización (o ambas exitosas en algunos casos)
	// Lo importante es que no se pierdan updates
	var finalCounter int
	db.QueryRowContext(ctx, "SELECT counter FROM isolation_test WHERE id = 1").Scan(&finalCounter)

	// El counter final debe ser consistente
	assert.GreaterOrEqual(t, finalCounter, 1, "Al menos una transacción debe haber tenido éxito")
}
