//go:build integration

//nolint:errcheck // Tests: errores de Close() en cleanup se ignoran intencionalmente
package containers

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	err = os.WriteFile(sqlFile, []byte(sqlContent), 0600)
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
	err = os.WriteFile(sqlFile, []byte(invalidSQL), 0600)
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

	err = os.WriteFile(sqlFile, []byte(""), 0600)
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

	err = os.WriteFile(sqlFile, []byte(multiSQL), 0600)
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

// Los tests de logica pura de WaitForHealthy y RetryOperation (sin Docker) viven en
// helpers_unit_test.go, que no lleva build tag y corre tambien con -tags integration.

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

	err = os.WriteFile(schemaFile, []byte(schemaSQL), 0600)
	require.NoError(t, err)

	err = ExecSQLFile(ctx, db, schemaFile)
	require.NoError(t, err)

	// Verificar resultado final
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM integration_test").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
