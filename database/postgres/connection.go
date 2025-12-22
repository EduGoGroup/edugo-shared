package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Importar el driver de PostgreSQL
)

const (
	// DefaultHealthCheckTimeout is the default timeout for health checks
	DefaultHealthCheckTimeout = 5 * time.Second
)

// Connect establece una conexión a PostgreSQL y retorna un pool de conexiones
func Connect(cfg *Config) (*sql.DB, error) {
	// Construir DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
		int(cfg.ConnectTimeout.Seconds()),
	)

	// Abrir conexión
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Verificar conectividad con timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close() //nolint:errcheck // Ignorar error en cleanup, el error principal es el de ping
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// HealthCheck verifica si la conexión a la base de datos está activa
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultHealthCheckTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetStats retorna estadísticas del pool de conexiones
func GetStats(db *sql.DB) sql.DBStats {
	return db.Stats()
}

// Close cierra todas las conexiones del pool
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
