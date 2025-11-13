package containers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

// ExecSQLFile lee y ejecuta un archivo SQL en la base de datos PostgreSQL
// ExecSQLFile lee y ejecuta un archivo SQL en la base de datos PostgreSQL.
//
// Parámetros:
//   - ctx: Contexto para la operación
//   - db: Conexión a la base de datos PostgreSQL
//   - filePath: Ruta al archivo SQL a ejecutar
//
// Retorna un error si el archivo no puede leerse o si la ejecución del SQL falla.
func ExecSQLFile(ctx context.Context, db *sql.DB, filePath string) error {
	// Leer archivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error leyendo archivo SQL %s: %w", filePath, err)
	}

	// Ejecutar SQL
	_, err = db.ExecContext(ctx, string(content))
	if err != nil {
		return fmt.Errorf("error ejecutando SQL de %s: %w", filePath, err)
	}

	return nil
}

// WaitForHealthy espera a que un servicio esté saludable
// Intenta hacer ping cada intervalo hasta alcanzar timeout
// WaitForHealthy espera a que un servicio esté saludable.
// Intenta ejecutar el health check cada intervalo hasta alcanzar el timeout.
//
// Parámetros:
//   - ctx: Contexto que puede cancelar la operación
//   - healthCheck: Función que retorna nil cuando el servicio está saludable
//   - timeout: Tiempo máximo de espera
//   - interval: Intervalo entre intentos
func WaitForHealthy(ctx context.Context, healthCheck func() error, timeout time.Duration, interval time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := healthCheck(); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("contexto cancelado mientras esperaba health check")
		case <-time.After(interval):
			// Continuar intentando
		}
	}

	return fmt.Errorf("timeout esperando health check después de %v", timeout)
}

// RetryOperation ejecuta una operación con reintentos
// RetryOperation ejecuta una operación con reintentos.
//
// Parámetros:
//   - operation: Función a ejecutar que retorna error
//   - maxRetries: Número máximo de intentos
//   - delay: Tiempo de espera entre reintentos
func RetryOperation(operation func() error, maxRetries int, delay time.Duration) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := operation(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < maxRetries-1 {
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("operación falló después de %d intentos: %w", maxRetries, lastErr)
}
