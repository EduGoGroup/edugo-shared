package containers

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// skipIfNotIntegration skipea el test si no está habilitada INTEGRATION_TESTS
func skipIfNotIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test - set INTEGRATION_TESTS=true to run")
	}
}

// TestExecSQLFile_Success verifica ejecución exitosa de archivo SQL
func TestExecSQLFile_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	// Setup container
	config := NewConfig().
		WithPostgreSQL(&PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	// Obtener conexión raw
	dsn, err := pg.ConnectionString(context.Background())
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	// Crear archivo SQL temporal
	tmpDir := t.TempDir()
	sqlFile := filepath.Join(tmpDir, "test.sql")

	sqlContent := `
		CREATE TABLE test_table (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			created_at TIMESTAMP DEFAULT NOW()
		);

		INSERT INTO test_table (name) VALUES ('test1');
		INSERT INTO test_table (name) VALUES ('test2');
		INSERT INTO test_table (name) VALUES ('test3');
	`

	err = os.WriteFile(sqlFile, []byte(sqlContent), 0644)
	require.NoError(t, err)

	// Ejecutar archivo SQL
	err = ExecSQLFile(ctx, db, sqlFile)
	require.NoError(t, err)

	// Verificar que se creó la tabla y se insertaron los datos
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_table").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

// TestExecSQLFile_FileNotFound verifica error cuando archivo no existe
func TestExecSQLFile_FileNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := NewConfig().
		WithPostgreSQL(&PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	dsn, err := pg.ConnectionString(context.Background())
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	// Intentar ejecutar archivo inexistente
	err = ExecSQLFile(ctx, db, "/path/to/nonexistent/file.sql")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error leyendo archivo SQL")
}

// TestExecSQLFile_InvalidSQL verifica error con SQL inválido
func TestExecSQLFile_InvalidSQL(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := NewConfig().
		WithPostgreSQL(&PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	dsn, err := pg.ConnectionString(context.Background())
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	// Crear archivo con SQL inválido
	tmpDir := t.TempDir()
	sqlFile := filepath.Join(tmpDir, "invalid.sql")

	invalidSQL := "THIS IS NOT VALID SQL SYNTAX!!!"
	err = os.WriteFile(sqlFile, []byte(invalidSQL), 0644)
	require.NoError(t, err)

	// Ejecutar debe fallar
	err = ExecSQLFile(ctx, db, sqlFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error ejecutando SQL")
}

// TestExecSQLFile_EmptyFile verifica ejecución de archivo vacío
func TestExecSQLFile_EmptyFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := NewConfig().
		WithPostgreSQL(&PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	dsn, err := pg.ConnectionString(context.Background())
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	// Crear archivo vacío
	tmpDir := t.TempDir()
	sqlFile := filepath.Join(tmpDir, "empty.sql")

	err = os.WriteFile(sqlFile, []byte(""), 0644)
	require.NoError(t, err)

	// Ejecutar archivo vacío (no debe dar error)
	err = ExecSQLFile(ctx, db, sqlFile)
	assert.NoError(t, err)
}

// TestExecSQLFile_MultipleStatements verifica múltiples statements
func TestExecSQLFile_MultipleStatements(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := NewConfig().
		WithPostgreSQL(&PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()
	dsn, err := pg.ConnectionString(context.Background())
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	tmpDir := t.TempDir()
	sqlFile := filepath.Join(tmpDir, "multi.sql")

	multiSQL := `
		CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100));
		CREATE TABLE posts (id SERIAL PRIMARY KEY, user_id INT, title VARCHAR(200));

		INSERT INTO users (name) VALUES ('Alice');
		INSERT INTO users (name) VALUES ('Bob');

		INSERT INTO posts (user_id, title) VALUES (1, 'Post 1');
		INSERT INTO posts (user_id, title) VALUES (1, 'Post 2');
		INSERT INTO posts (user_id, title) VALUES (2, 'Post 3');
	`

	err = os.WriteFile(sqlFile, []byte(multiSQL), 0644)
	require.NoError(t, err)

	err = ExecSQLFile(ctx, db, sqlFile)
	require.NoError(t, err)

	// Verificar usuarios
	var userCount int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	require.NoError(t, err)
	assert.Equal(t, 2, userCount)

	// Verificar posts
	var postCount int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM posts").Scan(&postCount)
	require.NoError(t, err)
	assert.Equal(t, 3, postCount)
}

// TestWaitForHealthy_ImmediateSuccess verifica success inmediato
func TestWaitForHealthy_ImmediateSuccess(t *testing.T) {
	ctx := context.Background()

	// Health check que siempre pasa
	healthCheck := func() error {
		return nil
	}

	err := WaitForHealthy(ctx, healthCheck, 5*time.Second, 100*time.Millisecond)
	assert.NoError(t, err)
}

// TestWaitForHealthy_SuccessAfterRetries verifica success después de reintentos
func TestWaitForHealthy_SuccessAfterRetries(t *testing.T) {
	ctx := context.Background()

	attempts := atomic.Int32{}
	healthCheck := func() error {
		count := attempts.Add(1)
		if count < 3 {
			return errors.New("not ready yet")
		}
		return nil
	}

	start := time.Now()
	err := WaitForHealthy(ctx, healthCheck, 5*time.Second, 100*time.Millisecond)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, attempts.Load(), int32(3))
	assert.GreaterOrEqual(t, duration, 200*time.Millisecond, "Debe esperar al menos 2 intervalos")
}

// TestWaitForHealthy_Timeout verifica timeout
func TestWaitForHealthy_Timeout(t *testing.T) {
	ctx := context.Background()

	// Health check que siempre falla
	healthCheck := func() error {
		return errors.New("service not healthy")
	}

	start := time.Now()
	err := WaitForHealthy(ctx, healthCheck, 500*time.Millisecond, 100*time.Millisecond)
	duration := time.Since(start)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout esperando health check")
	assert.GreaterOrEqual(t, duration, 500*time.Millisecond)
	assert.Less(t, duration, 1*time.Second)
}

// TestWaitForHealthy_ContextCancellation verifica cancelación por contexto
func TestWaitForHealthy_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Health check que siempre falla
	healthCheck := func() error {
		return errors.New("not ready")
	}

	// Cancelar después de 200ms
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := WaitForHealthy(ctx, healthCheck, 5*time.Second, 100*time.Millisecond)
	duration := time.Since(start)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contexto cancelado")
	assert.Less(t, duration, 1*time.Second, "Debe cancelarse rápidamente")
}

// TestWaitForHealthy_ImmediateCancellation verifica cancelación inmediata
func TestWaitForHealthy_ImmediateCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	healthCheck := func() error {
		return errors.New("not ready")
	}

	err := WaitForHealthy(ctx, healthCheck, 5*time.Second, 100*time.Millisecond)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contexto cancelado")
}

// TestRetryOperation_ImmediateSuccess verifica success inmediato
func TestRetryOperation_ImmediateSuccess(t *testing.T) {
	attempts := atomic.Int32{}

	operation := func() error {
		attempts.Add(1)
		return nil
	}

	err := RetryOperation(operation, 3, 100*time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), attempts.Load(), "Solo debe intentar una vez")
}

// TestRetryOperation_SuccessAfterRetries verifica success después de reintentos
func TestRetryOperation_SuccessAfterRetries(t *testing.T) {
	attempts := atomic.Int32{}

	operation := func() error {
		count := attempts.Add(1)
		if count < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	start := time.Now()
	err := RetryOperation(operation, 5, 50*time.Millisecond)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, int32(3), attempts.Load())
	assert.GreaterOrEqual(t, duration, 100*time.Millisecond, "Debe esperar delay entre intentos")
}

// TestRetryOperation_AllFailed verifica fallo después de todos los intentos
func TestRetryOperation_AllFailed(t *testing.T) {
	attempts := atomic.Int32{}
	expectedError := errors.New("permanent error")

	operation := func() error {
		attempts.Add(1)
		return expectedError
	}

	maxRetries := 5
	err := RetryOperation(operation, maxRetries, 10*time.Millisecond)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operación falló después de")
	assert.Contains(t, err.Error(), "intentos")
	assert.Equal(t, int32(maxRetries), attempts.Load())
}

// TestRetryOperation_ZeroRetries verifica comportamiento con cero reintentos
func TestRetryOperation_ZeroRetries(t *testing.T) {
	attempts := atomic.Int32{}

	operation := func() error {
		attempts.Add(1)
		return errors.New("error")
	}

	err := RetryOperation(operation, 0, 100*time.Millisecond)

	// Con 0 reintentos, no debe ejecutarse
	assert.Error(t, err)
	assert.Equal(t, int32(0), attempts.Load())
}

// TestRetryOperation_OneRetry verifica con un solo reintento
func TestRetryOperation_OneRetry(t *testing.T) {
	attempts := atomic.Int32{}

	operation := func() error {
		attempts.Add(1)
		return errors.New("always fails")
	}

	err := RetryOperation(operation, 1, 10*time.Millisecond)

	assert.Error(t, err)
	assert.Equal(t, int32(1), attempts.Load())
}

// TestRetryOperation_DelayBetweenAttempts verifica delay entre intentos
func TestRetryOperation_DelayBetweenAttempts(t *testing.T) {
	attempts := atomic.Int32{}
	timestamps := []time.Time{}

	operation := func() error {
		attempts.Add(1)
		timestamps = append(timestamps, time.Now())
		return errors.New("error")
	}

	delay := 100 * time.Millisecond
	_ = RetryOperation(operation, 3, delay)

	// Debe haber 3 intentos
	assert.Equal(t, int32(3), attempts.Load())
	assert.Len(t, timestamps, 3)

	// Verificar delays entre intentos
	if len(timestamps) >= 2 {
		diff1 := timestamps[1].Sub(timestamps[0])
		assert.GreaterOrEqual(t, diff1, delay)
	}

	if len(timestamps) >= 3 {
		diff2 := timestamps[2].Sub(timestamps[1])
		assert.GreaterOrEqual(t, diff2, delay)
	}
}

// TestRetryOperation_VariousErrors verifica diferentes tipos de errores
func TestRetryOperation_VariousErrors(t *testing.T) {
	tests := []struct {
		name          string
		errors        []error
		maxRetries    int
		expectedError bool
		expectedCalls int
	}{
		{
			name:          "success on first try",
			errors:        []error{nil},
			maxRetries:    3,
			expectedError: false,
			expectedCalls: 1,
		},
		{
			name:          "success on second try",
			errors:        []error{errors.New("err1"), nil},
			maxRetries:    3,
			expectedError: false,
			expectedCalls: 2,
		},
		{
			name:          "all attempts fail",
			errors:        []error{errors.New("err1"), errors.New("err2"), errors.New("err3")},
			maxRetries:    3,
			expectedError: true,
			expectedCalls: 3,
		},
		{
			name:          "success on last attempt",
			errors:        []error{errors.New("err1"), errors.New("err2"), nil},
			maxRetries:    3,
			expectedError: false,
			expectedCalls: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := atomic.Int32{}
			operation := func() error {
				idx := int(callCount.Load())
				callCount.Add(1)
				if idx < len(tt.errors) {
					return tt.errors[idx]
				}
				return errors.New("unexpected call")
			}

			err := RetryOperation(operation, tt.maxRetries, 1*time.Millisecond)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, int32(tt.expectedCalls), callCount.Load())
		})
	}
}

// TestHelpers_Integration verifica integración de todas las funciones helper
func TestHelpers_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := NewConfig().
		WithPostgreSQL(&PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	// 1. Usar RetryOperation para conectar a la DB (con reintentos)
	var db *sql.DB
	err = RetryOperation(func() error {
		var err error
		dsn, err := pg.ConnectionString(context.Background())
		if err != nil {
			return err
		}
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}, 5, 100*time.Millisecond)
	require.NoError(t, err)
	defer db.Close()

	// 2. Usar WaitForHealthy para verificar que la DB está lista
	err = WaitForHealthy(ctx, func() error {
		return db.Ping()
	}, 10*time.Second, 500*time.Millisecond)
	require.NoError(t, err)

	// 3. Usar ExecSQLFile para inicializar schema
	tmpDir := t.TempDir()
	schemaFile := filepath.Join(tmpDir, "schema.sql")

	schemaSQL := `
		CREATE TABLE integration_test (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			value INT
		);
		INSERT INTO integration_test (name, value) VALUES ('test', 42);
	`

	err = os.WriteFile(schemaFile, []byte(schemaSQL), 0644)
	require.NoError(t, err)

	err = ExecSQLFile(ctx, db, schemaFile)
	require.NoError(t, err)

	// Verificar resultado final
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM integration_test").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
