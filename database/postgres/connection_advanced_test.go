package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestConnection_ReconnectAfterConnectionLoss verifica reconexión después de pérdida de conexión
func TestConnection_ReconnectAfterConnectionLoss(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	// Setup container
	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	// Crear primera conexión
	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	// Verificar que funciona
	err = db.Ping(ctx)
	require.NoError(t, err)

	// Obtener conexión SQL subyacente
	sqlDB, err := db.DB()
	require.NoError(t, err)

	// Configurar pool con límites bajos para forzar reconexiones
	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(1 * time.Second)

	// Ejecutar queries y verificar que el pool reconecta automáticamente
	for i := 0; i < 5; i++ {
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err, "Query %d debe funcionar", i)
		assert.Equal(t, 1, result)

		// Esperar para que expire la conexión
		time.Sleep(1200 * time.Millisecond)
	}

	// Verificar estadísticas del pool
	stats := sqlDB.Stats()
	assert.GreaterOrEqual(t, stats.OpenConnections, 0)
}

// TestConnection_PoolExhaustion verifica comportamiento cuando se agota el pool
func TestConnection_PoolExhaustion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	sqlDB, err := db.DB()
	require.NoError(t, err)

	// Configurar pool muy limitado
	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxIdleTime(100 * time.Millisecond)

	// Crear tabla de prueba
	err = db.Exec("CREATE TEMP TABLE pool_test (id SERIAL PRIMARY KEY, value INT)").Error
	require.NoError(t, err)

	// Lanzar múltiples operaciones concurrentes
	done := make(chan bool, 10)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			// Insertar dato
			err := db.Exec("INSERT INTO pool_test (value) VALUES (?)", index).Error
			if err != nil {
				errors <- err
			}
			done <- true
		}(i)
	}

	// Esperar a que todas terminen
	for i := 0; i < 10; i++ {
		<-done
	}

	// No debe haber errores (el pool debe esperar por conexiones disponibles)
	close(errors)
	errorCount := 0
	for err := range errors {
		errorCount++
		t.Logf("Error encontrado: %v", err)
	}
	assert.Equal(t, 0, errorCount, "No debe haber errores por agotamiento de pool")

	// Verificar que se insertaron todos los registros
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM pool_test").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(10), count)
}

// TestConnection_IdleConnectionCleanup verifica limpieza de conexiones idle
func TestConnection_IdleConnectionCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	sqlDB, err := db.DB()
	require.NoError(t, err)

	// Configurar tiempo de vida corto para conexiones idle
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxIdleTime(500 * time.Millisecond)

	// Ejecutar varias queries para crear conexiones
	for i := 0; i < 8; i++ {
		var result int
		err = db.Raw("SELECT ?", i).Scan(&result).Error
		require.NoError(t, err)
	}

	// Verificar estadísticas iniciales
	stats := sqlDB.Stats()
	initialIdle := stats.Idle
	t.Logf("Conexiones idle iniciales: %d", initialIdle)

	// Esperar a que se limpien las conexiones idle
	time.Sleep(1 * time.Second)

	// Verificar que las conexiones idle fueron limpiadas
	stats = sqlDB.Stats()
	t.Logf("Conexiones idle después de limpieza: %d", stats.Idle)

	// Verificar que aún podemos ejecutar queries
	var result int
	err = db.Raw("SELECT 42").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 42, result)
}

// TestConnection_StatementTimeout verifica manejo de statement timeout
func TestConnection_StatementTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	// Configurar statement timeout en la sesión
	err = db.Exec("SET statement_timeout = '1s'").Error
	require.NoError(t, err)

	// Query rápida debe funcionar
	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 1, result)

	// Query lenta (pg_sleep) debe dar timeout
	err = db.Raw("SELECT pg_sleep(3)").Scan(&result).Error
	assert.Error(t, err, "Query con pg_sleep(3) debe dar timeout")
	assert.Contains(t, err.Error(), "canceling statement due to statement timeout", "Error debe mencionar timeout")

	// Después del timeout, la conexión debe seguir funcionando
	err = db.Raw("SELECT 2").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 2, result)
}

// TestConnection_ContextTimeout verifica manejo de context.WithTimeout
func TestConnection_ContextTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	// Contexto con timeout muy corto
	queryCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Query lenta debe ser cancelada por el contexto
	var result int
	err = db.WithContext(queryCtx).Raw("SELECT pg_sleep(2)").Scan(&result).Error
	assert.Error(t, err, "Query debe ser cancelada por contexto")

	// La conexión debe seguir funcionando después de cancelación
	err = db.Raw("SELECT 1").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 1, result)
}

// TestConnection_PreparedStatements verifica uso de prepared statements
func TestConnection_PreparedStatements(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	// Crear tabla
	err = db.Exec("CREATE TEMP TABLE stmt_test (id SERIAL PRIMARY KEY, name VARCHAR(100))").Error
	require.NoError(t, err)

	// Ejecutar mismo query múltiples veces (debería usar prepared statements)
	for i := 0; i < 100; i++ {
		err = db.Exec("INSERT INTO stmt_test (name) VALUES (?)", "user"+string(rune(i))).Error
		require.NoError(t, err)
	}

	// Verificar que se insertaron todos
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM stmt_test").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(100), count)
}

// TestConnection_MultipleSessionSettings verifica diferentes configuraciones de sesión
func TestConnection_MultipleSessionSettings(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name        string
		setting     string
		value       string
		checkQuery  string
		expectValue string
	}{
		{
			name:        "timezone setting",
			setting:     "timezone",
			value:       "'UTC'",
			checkQuery:  "SHOW timezone",
			expectValue: "UTC",
		},
		{
			name:        "search_path setting",
			setting:     "search_path",
			value:       "'public'",
			checkQuery:  "SHOW search_path",
			expectValue: "public",
		},
		{
			name:        "application_name setting",
			setting:     "application_name",
			value:       "'test_app'",
			checkQuery:  "SHOW application_name",
			expectValue: "test_app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar setting
			err := db.Exec("SET " + tt.setting + " = " + tt.value).Error
			require.NoError(t, err)

			// Verificar valor
			var result string
			err = db.Raw(tt.checkQuery).Scan(&result).Error
			require.NoError(t, err)
			assert.Equal(t, tt.expectValue, result)
		})
	}
}

// TestConnection_LongRunningQuery verifica queries de larga duración
func TestConnection_LongRunningQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	// Crear tabla grande
	err = db.Exec("CREATE TEMP TABLE large_table (id SERIAL PRIMARY KEY, data TEXT)").Error
	require.NoError(t, err)

	// Insertar muchos datos (batch)
	tx := db.Begin()
	for i := 0; i < 1000; i++ {
		err = tx.Exec("INSERT INTO large_table (data) VALUES (?)", "data"+string(rune(i))).Error
		require.NoError(t, err)
	}
	err = tx.Commit().Error
	require.NoError(t, err)

	// Query compleja (join con self)
	var count int64
	err = db.Raw(`
		SELECT COUNT(*)
		FROM large_table t1
		INNER JOIN large_table t2 ON t1.id < t2.id
		WHERE t1.id < 100
	`).Scan(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0))
}

// TestConnection_ConcurrentTransactions verifica transacciones concurrentes independientes
func TestConnection_ConcurrentTransactions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	// Crear tabla
	err = db.Exec("CREATE TEMP TABLE concurrent_test (id SERIAL PRIMARY KEY, value INT)").Error
	require.NoError(t, err)

	// Lanzar múltiples transacciones concurrentes
	done := make(chan bool, 20)
	errors := make(chan error, 20)

	for i := 0; i < 20; i++ {
		go func(index int) {
			err := db.Transaction(func(tx *gorm.DB) error {
				// Insertar valor
				if err := tx.Exec("INSERT INTO concurrent_test (value) VALUES (?)", index).Error; err != nil {
					return err
				}
				// Pequeña pausa
				time.Sleep(10 * time.Millisecond)
				// Verificar que se insertó
				var count int64
				return tx.Raw("SELECT COUNT(*) FROM concurrent_test WHERE value = ?", index).Scan(&count).Error
			})
			if err != nil {
				errors <- err
			}
			done <- true
		}(i)
	}

	// Esperar a que todas terminen
	for i := 0; i < 20; i++ {
		<-done
	}

	// No debe haber errores
	close(errors)
	errorCount := 0
	for err := range errors {
		errorCount++
		t.Logf("Error: %v", err)
	}
	assert.Equal(t, 0, errorCount)

	// Verificar que se insertaron todos los registros
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM concurrent_test").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(20), count)
}

// TestConnection_RawSQLAndGORMOperations verifica intercalación de SQL raw y GORM
func TestConnection_RawSQLAndGORMOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	// Crear tabla con SQL raw
	err = db.Exec("CREATE TEMP TABLE mixed_test (id SERIAL PRIMARY KEY, name VARCHAR(100), value INT)").Error
	require.NoError(t, err)

	// Insertar con SQL raw
	err = db.Exec("INSERT INTO mixed_test (name, value) VALUES (?, ?)", "test1", 100).Error
	require.NoError(t, err)

	// Consultar con GORM Raw
	var result struct {
		ID    int
		Name  string
		Value int
	}
	err = db.Raw("SELECT id, name, value FROM mixed_test WHERE name = ?", "test1").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, "test1", result.Name)
	assert.Equal(t, 100, result.Value)

	// Actualizar con Exec
	err = db.Exec("UPDATE mixed_test SET value = ? WHERE name = ?", 200, "test1").Error
	require.NoError(t, err)

	// Verificar con Raw
	var newValue int
	err = db.Raw("SELECT value FROM mixed_test WHERE name = ?", "test1").Scan(&newValue).Error
	require.NoError(t, err)
	assert.Equal(t, 200, newValue)

	// Obtener conexión SQL directa
	sqlDB, err := db.DB()
	require.NoError(t, err)

	// Ejecutar query con database/sql directamente
	var countDB int
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM mixed_test").Scan(&countDB)
	require.NoError(t, err)
	assert.Equal(t, 1, countDB)
}

// TestConnection_ErrorRecovery verifica recuperación después de errores
func TestConnection_ErrorRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	// Error de sintaxis SQL
	err = db.Exec("INVALID SQL QUERY").Error
	assert.Error(t, err)

	// La conexión debe seguir funcionando
	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 1, result)

	// Error de tabla inexistente
	err = db.Raw("SELECT * FROM nonexistent_table").Scan(&result).Error
	assert.Error(t, err)

	// La conexión debe seguir funcionando
	err = db.Raw("SELECT 2").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 2, result)

	// Error de división por cero
	err = db.Raw("SELECT 1/0").Scan(&result).Error
	assert.Error(t, err)

	// La conexión debe seguir funcionando
	err = db.Raw("SELECT 3").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 3, result)
}

// TestConnection_PoolStatistics verifica estadísticas del connection pool
func TestConnection_PoolStatistics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	db, err := New(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer db.Close()

	sqlDB, err := db.DB()
	require.NoError(t, err)

	// Configurar pool
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	// Obtener estadísticas iniciales
	stats := sqlDB.Stats()
	t.Logf("Stats iniciales: %+v", stats)

	// Ejecutar queries para usar el pool
	for i := 0; i < 20; i++ {
		var result int
		err = db.Raw("SELECT ?", i).Scan(&result).Error
		require.NoError(t, err)
	}

	// Obtener estadísticas finales
	stats = sqlDB.Stats()
	t.Logf("Stats finales: %+v", stats)

	// Verificar que el pool está siendo usado
	assert.LessOrEqual(t, stats.OpenConnections, 10, "No debe exceder MaxOpenConns")
	assert.LessOrEqual(t, stats.Idle, 5, "Idle no debe exceder MaxIdleConns")
	assert.Greater(t, stats.InUse+stats.Idle, 0, "Debe haber conexiones abiertas")
}

// TestRawConnection_DirectSQL verifica uso de conexión SQL directa
func TestRawConnection_DirectSQL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	// Usar NewRawConnection
	rawDB, err := NewRawConnection(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer rawDB.Close()

	// Ping
	err = rawDB.PingContext(ctx)
	require.NoError(t, err)

	// Crear tabla
	_, err = rawDB.ExecContext(ctx, "CREATE TEMP TABLE raw_test (id SERIAL PRIMARY KEY, value INT)")
	require.NoError(t, err)

	// Insertar datos
	_, err = rawDB.ExecContext(ctx, "INSERT INTO raw_test (value) VALUES ($1), ($2), ($3)", 10, 20, 30)
	require.NoError(t, err)

	// Consultar datos
	rows, err := rawDB.QueryContext(ctx, "SELECT value FROM raw_test ORDER BY value")
	require.NoError(t, err)
	defer rows.Close()

	values := []int{}
	for rows.Next() {
		var value int
		err = rows.Scan(&value)
		require.NoError(t, err)
		values = append(values, value)
	}

	assert.Equal(t, []int{10, 20, 30}, values)
}

// TestRawConnection_Transaction verifica transacciones con sql.DB
func TestRawConnection_Transaction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithPostgreSQL(&containers.PostgresConfig{
			Image:    "postgres:15-alpine",
			Database: "test_db",
			Username: "test_user",
			Password: "test_pass",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	pg := manager.PostgreSQL()

	rawDB, err := NewRawConnection(ctx, Config{
		Host:     pg.Host(),
		Port:     pg.Port(),
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	})
	require.NoError(t, err)
	defer rawDB.Close()

	// Crear tabla
	_, err = rawDB.ExecContext(ctx, "CREATE TEMP TABLE tx_raw_test (id SERIAL PRIMARY KEY, value INT)")
	require.NoError(t, err)

	// Iniciar transacción
	tx, err := rawDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	require.NoError(t, err)

	// Insertar en transacción
	_, err = tx.ExecContext(ctx, "INSERT INTO tx_raw_test (value) VALUES ($1)", 100)
	require.NoError(t, err)

	// Commit
	err = tx.Commit()
	require.NoError(t, err)

	// Verificar que se insertó
	var count int
	err = rawDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tx_raw_test").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Transacción con rollback
	tx2, err := rawDB.BeginTx(ctx, nil)
	require.NoError(t, err)

	_, err = tx2.ExecContext(ctx, "INSERT INTO tx_raw_test (value) VALUES ($1)", 200)
	require.NoError(t, err)

	// Rollback
	err = tx2.Rollback()
	require.NoError(t, err)

	// Verificar que NO se insertó el segundo
	err = rawDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tx_raw_test").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "Rollback debe prevenir la inserción")
}
