package bootstrap

import (
	"context"
	"fmt"
)

// initPostgreSQL inicializa la conexión a PostgreSQL.
//
// Parámetros:
//   - ctx: Contexto para cancelación
//   - config: Configuración de la aplicación
//   - factories: Fábricas disponibles
//   - resources: Recursos a inicializar
//   - lifecycleManager: Manager de lifecycle para cleanup
//   - opts: Opciones de bootstrap
//
// Retorna error si el recurso es requerido y falla la inicialización.
func initPostgreSQL(
	ctx context.Context,
	config interface{},
	factories *Factories,
	resources *Resources,
	lifecycleManager interface{},
	opts *BootstrapOptions,
) error {
	if factories.PostgreSQL == nil {
		return fmt.Errorf("postgresql factory not provided")
	}

	// Extraer configuración de PostgreSQL
	pgConfig, err := extractPostgreSQLConfig(config)
	if err != nil {
		return fmt.Errorf("failed to extract PostgreSQL config: %w", err)
	}

	// Log inicio
	if resources.Logger != nil {
		resources.Logger.Info("Initializing PostgreSQL connection...")
	}

	// Crear conexión
	db, err := factories.PostgreSQL.CreateConnection(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL connection: %w", err)
	}

	resources.PostgreSQL = db

	// Registrar cleanup en lifecycle manager si está disponible
	if lifecycleManager != nil {
		registerPostgreSQLCleanup(lifecycleManager, factories.PostgreSQL, db, resources.Logger)
	}

	// Log éxito
	if resources.Logger != nil {
		resources.Logger.With(
			"host", pgConfig.Host,
			"port", pgConfig.Port,
			"database", pgConfig.Database,
		).Info("PostgreSQL connection established")
	}

	return nil
}
