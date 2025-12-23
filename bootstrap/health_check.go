package bootstrap

import (
	"context"
	"fmt"
	"time"
)

// performHealthChecks ejecuta verificaciones de salud en los recursos inicializados.
//
// Parámetros:
//   - ctx: Contexto para cancelación
//   - resources: Recursos a verificar
//   - opts: Opciones de bootstrap
//
// Retorna error si alguna verificación falla.
func performHealthChecks(ctx context.Context, resources *Resources, opts *BootstrapOptions) error {
	if resources.Logger != nil {
		resources.Logger.Info("Performing health checks...")
	}

	// Health check timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// PostgreSQL health check
	if resources.PostgreSQL != nil {
		if err := resources.PostgreSQL.Raw("SELECT 1").Error; err != nil {
			return fmt.Errorf("postgresql health check failed: %w", err)
		}
		if resources.Logger != nil {
			resources.Logger.Debug("PostgreSQL health check passed")
		}
	}

	// MongoDB health check
	if resources.MongoDB != nil {
		if err := resources.MongoDB.Ping(ctx, nil); err != nil {
			return fmt.Errorf("mongodb health check failed: %w", err)
		}
		if resources.Logger != nil {
			resources.Logger.Debug("MongoDB health check passed")
		}
	}

	if resources.Logger != nil {
		resources.Logger.Info("All health checks passed")
	}

	return nil
}
