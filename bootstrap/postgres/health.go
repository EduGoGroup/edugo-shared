package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	// DefaultHealthCheckTimeout es el timeout por defecto para health checks.
	DefaultHealthCheckTimeout = 5 * time.Second
)

// HealthCheck verifica si la conexion a la base de datos esta activa.
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultHealthCheckTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetStats retorna estadisticas del pool de conexiones.
func GetStats(db *sql.DB) sql.DBStats {
	return db.Stats()
}
