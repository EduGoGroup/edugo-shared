package bootstrap

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestPostgreSQLFactory_CreateConnection_Success verifica creación exitosa de conexión
func TestPostgreSQLFactory_CreateConnection_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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
	require.NotNil(t, pg)

	// Crear factory
	factory := NewDefaultPostgreSQLFactory(nil)

	// Obtener datos de conexión
	host, err := pg.Host(ctx)
	require.NoError(t, err)
	port, err := pg.MappedPort(ctx)
	require.NoError(t, err)

	// Configuración de PostgreSQL
	pgConfig := PostgreSQLConfig{
		Host:     host,
		Port:     port,
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	}

	// Crear conexión
	db, err := factory.CreateConnection(ctx, pgConfig)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() {
		if err := factory.Close(db); err != nil {
			t.Logf("Failed to close PostgreSQL connection: %v", err)
		}
	}()

	// Verificar que la conexión funciona
	sqlDB, err := db.DB()
	require.NoError(t, err)
	assert.NotNil(t, sqlDB)
}

// TestPostgreSQLFactory_CreateConnection_InvalidConfig verifica error con config inválida
func TestPostgreSQLFactory_CreateConnection_InvalidConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()
	factory := NewDefaultPostgreSQLFactory(nil)

	// Configuración inválida
	invalidConfig := PostgreSQLConfig{
		Host:     "invalid-host-that-does-not-exist",
		Port:     5432,
		User:     "invalid_user",
		Password: "invalid_pass",
		Database: "invalid_db",
		SSLMode:  "disable",
	}

	// Debe fallar al crear conexión
	db, err := factory.CreateConnection(ctx, invalidConfig)
	assert.Error(t, err)
	assert.Nil(t, db)
}

// TestPostgreSQLFactory_CreateConnection_WithSSLMode verifica diferentes modos SSL
func TestPostgreSQLFactory_CreateConnection_WithSSLMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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
	require.NotNil(t, pg)

	factory := NewDefaultPostgreSQLFactory(nil)

	host, err := pg.Host(ctx)
	require.NoError(t, err)
	port, err := pg.MappedPort(ctx)
	require.NoError(t, err)

	tests := []struct {
		name    string
		sslMode string
		wantErr bool
	}{
		{
			name:    "SSL disable",
			sslMode: "disable",
			wantErr: false,
		},
		{
			name:    "SSL prefer",
			sslMode: "prefer",
			wantErr: false,
		},
		{
			name:    "empty SSL mode (default to disable)",
			sslMode: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgConfig := PostgreSQLConfig{
				Host:     host,
				Port:     port,
				User:     pg.Username(),
				Password: pg.Password(),
				Database: pg.Database(),
				SSLMode:  tt.sslMode,
			}

			db, err := factory.CreateConnection(ctx, pgConfig)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				if db != nil {
					if closeErr := factory.Close(db); closeErr != nil {
						t.Logf("Failed to close PostgreSQL connection: %v", closeErr)
					}
				}
			}
		})
	}
}

// TestPostgreSQLFactory_Ping_Success verifica ping exitoso
func TestPostgreSQLFactory_Ping_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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
	require.NotNil(t, pg)

	factory := NewDefaultPostgreSQLFactory(nil)

	host, err := pg.Host(ctx)
	require.NoError(t, err)
	port, err := pg.MappedPort(ctx)
	require.NoError(t, err)

	pgConfig := PostgreSQLConfig{
		Host:     host,
		Port:     port,
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	}

	db, err := factory.CreateConnection(ctx, pgConfig)
	require.NoError(t, err)
	defer func() {
		if err := factory.Close(db); err != nil {
			t.Logf("Failed to close PostgreSQL connection: %v", err)
		}
	}()

	// Ping debe ser exitoso
	err = factory.Ping(ctx, db)
	assert.NoError(t, err)
}

// TestPostgreSQLFactory_Close_Success verifica cierre exitoso
func TestPostgreSQLFactory_Close_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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
	require.NotNil(t, pg)

	factory := NewDefaultPostgreSQLFactory(nil)

	host, err := pg.Host(ctx)
	require.NoError(t, err)
	port, err := pg.MappedPort(ctx)
	require.NoError(t, err)

	pgConfig := PostgreSQLConfig{
		Host:     host,
		Port:     port,
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	}

	db, err := factory.CreateConnection(ctx, pgConfig)
	require.NoError(t, err)

	// Close debe ser exitoso
	err = factory.Close(db)
	assert.NoError(t, err)
}

// TestPostgreSQLFactory_CreateRawConnection_Success verifica creación de conexión SQL raw
func TestPostgreSQLFactory_CreateRawConnection_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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
	require.NotNil(t, pg)

	factory := NewDefaultPostgreSQLFactory(nil)

	host, err := pg.Host(ctx)
	require.NoError(t, err)
	port, err := pg.MappedPort(ctx)
	require.NoError(t, err)

	pgConfig := PostgreSQLConfig{
		Host:     host,
		Port:     port,
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	}

	// Crear raw connection
	db, err := factory.CreateRawConnection(ctx, pgConfig)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close raw PostgreSQL connection: %v", err)
		}
	}()

	// Verificar que funciona
	err = db.PingContext(ctx)
	assert.NoError(t, err)
}

// TestPostgreSQLFactory_CreateRawConnection_InvalidConfig verifica error con config inválida
func TestPostgreSQLFactory_CreateRawConnection_InvalidConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()
	factory := NewDefaultPostgreSQLFactory(nil)

	invalidConfig := PostgreSQLConfig{
		Host:     "invalid-host",
		Port:     5432,
		User:     "invalid",
		Password: "invalid",
		Database: "invalid",
		SSLMode:  "disable",
	}

	db, err := factory.CreateRawConnection(ctx, invalidConfig)
	assert.Error(t, err)
	assert.Nil(t, db)
}

// TestPostgreSQLFactory_ConnectionPoolSettings verifica configuración del pool
func TestPostgreSQLFactory_ConnectionPoolSettings(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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
	require.NotNil(t, pg)

	factory := NewDefaultPostgreSQLFactory(nil)

	host, err := pg.Host(ctx)
	require.NoError(t, err)
	port, err := pg.MappedPort(ctx)
	require.NoError(t, err)

	pgConfig := PostgreSQLConfig{
		Host:     host,
		Port:     port,
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	}

	db, err := factory.CreateConnection(ctx, pgConfig)
	require.NoError(t, err)
	defer func() {
		if err := factory.Close(db); err != nil {
			t.Logf("Failed to close PostgreSQL connection: %v", err)
		}
	}()

	// Obtener stats del pool
	sqlDB, err := db.DB()
	require.NoError(t, err)

	stats := sqlDB.Stats()

	// Verificar que el pool está configurado
	assert.Equal(t, 25, stats.MaxOpenConnections, "MaxOpenConns debe ser 25")
	assert.GreaterOrEqual(t, stats.Idle, 0, "Debe tener conexiones idle configuradas")
}

// TestPostgreSQLFactory_WithCustomLogger verifica creación con logger custom
func TestPostgreSQLFactory_WithCustomLogger(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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
	require.NotNil(t, pg)

	// Crear factory con logger custom
	customLogger := logger.Default.LogMode(logger.Silent)
	factory := NewDefaultPostgreSQLFactory(customLogger)

	host, err := pg.Host(ctx)
	require.NoError(t, err)
	port, err := pg.MappedPort(ctx)
	require.NoError(t, err)

	pgConfig := PostgreSQLConfig{
		Host:     host,
		Port:     port,
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	}

	db, err := factory.CreateConnection(ctx, pgConfig)
	require.NoError(t, err)
	assert.NotNil(t, db)
	defer func() {
		if err := factory.Close(db); err != nil {
			t.Logf("Failed to close PostgreSQL connection: %v", err)
		}
	}()
}

// TestPostgreSQLFactory_MultipleConnections verifica múltiples conexiones
func TestPostgreSQLFactory_MultipleConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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
	require.NotNil(t, pg)

	factory := NewDefaultPostgreSQLFactory(nil)

	host, err := pg.Host(ctx)
	require.NoError(t, err)
	port, err := pg.MappedPort(ctx)
	require.NoError(t, err)

	pgConfig := PostgreSQLConfig{
		Host:     host,
		Port:     port,
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	}

	// Crear múltiples conexiones
	connections := make([]*gorm.DB, 3)
	for i := 0; i < 3; i++ {
		db, err := factory.CreateConnection(ctx, pgConfig)
		require.NoError(t, err)
		connections[i] = db
	}

	// Verificar que todas funcionan
	for i, db := range connections {
		err := factory.Ping(ctx, db)
		assert.NoError(t, err, "Conexión %d debe funcionar", i)
	}

	// Cerrar todas
	for _, db := range connections {
		err := factory.Close(db)
		assert.NoError(t, err)
	}
}

// TestPostgreSQLFactory_BuildDSN_Scenarios verifica construcción de DSN
func TestPostgreSQLFactory_BuildDSN_Scenarios(t *testing.T) {
	factory := NewDefaultPostgreSQLFactory(nil)

	tests := []struct {
		name     string
		config   PostgreSQLConfig
		contains []string
	}{
		{
			name: "basic config",
			config: PostgreSQLConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "pass",
				Database: "db",
				SSLMode:  "disable",
			},
			contains: []string{"host=localhost", "port=5432", "user=user", "password=pass", "dbname=db", "sslmode=disable"},
		},
		{
			name: "config without sslmode (defaults to disable)",
			config: PostgreSQLConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "pass",
				Database: "db",
				SSLMode:  "",
			},
			contains: []string{"sslmode=disable"},
		},
		{
			name: "config with require sslmode",
			config: PostgreSQLConfig{
				Host:     "remote.host.com",
				Port:     5432,
				User:     "admin",
				Password: "secret",
				Database: "production",
				SSLMode:  "require",
			},
			contains: []string{"host=remote.host.com", "sslmode=require"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := factory.buildDSN(tt.config)

			for _, substr := range tt.contains {
				assert.Contains(t, dsn, substr)
			}
		})
	}
}

// TestPostgreSQLFactory_QueryExecution verifica ejecución de queries
func TestPostgreSQLFactory_QueryExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
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
	require.NotNil(t, pg)

	factory := NewDefaultPostgreSQLFactory(nil)

	host, err := pg.Host(ctx)
	require.NoError(t, err)
	port, err := pg.MappedPort(ctx)
	require.NoError(t, err)

	pgConfig := PostgreSQLConfig{
		Host:     host,
		Port:     port,
		User:     pg.Username(),
		Password: pg.Password(),
		Database: pg.Database(),
		SSLMode:  "disable",
	}

	db, err := factory.CreateConnection(ctx, pgConfig)
	require.NoError(t, err)
	defer func() {
		if err := factory.Close(db); err != nil {
			t.Logf("Failed to close PostgreSQL connection: %v", err)
		}
	}()

	// Ejecutar query simple
	var result int
	err = db.Raw("SELECT 1 as value").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 1, result)

	// Crear tabla temporal
	err = db.Exec("CREATE TEMP TABLE test_table (id SERIAL PRIMARY KEY, name VARCHAR(100))").Error
	require.NoError(t, err)

	// Insertar datos
	err = db.Exec("INSERT INTO test_table (name) VALUES (?)", "test").Error
	require.NoError(t, err)

	// Consultar datos
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM test_table").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// TestPostgreSQLFactory_ConnectionTimeout verifica manejo de timeout
func TestPostgreSQLFactory_ConnectionTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	// Contexto con timeout muy corto
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	factory := NewDefaultPostgreSQLFactory(nil)

	// Configuración con host que no responde rápido
	pgConfig := PostgreSQLConfig{
		Host:     "192.0.2.1", // TEST-NET-1 (no routable)
		Port:     5432,
		User:     "user",
		Password: "pass",
		Database: "db",
		SSLMode:  "disable",
	}

	// Debe fallar por timeout (o por host inválido)
	db, err := factory.CreateConnection(ctx, pgConfig)
	assert.Error(t, err)
	assert.Nil(t, db)
}
