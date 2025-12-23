package bootstrap

import (
	"context"
	"fmt"
)

// initLogger inicializa el logger para la aplicación.
//
// Parámetros:
//   - ctx: Contexto para cancelación
//   - config: Configuración de la aplicación
//   - factories: Fábricas disponibles
//   - resources: Recursos a inicializar
//   - opts: Opciones de bootstrap
//
// Retorna error si la inicialización falla.
func initLogger(
	ctx context.Context,
	config interface{},
	factories *Factories,
	resources *Resources,
	opts *BootstrapOptions,
) error {
	if factories.Logger == nil {
		return fmt.Errorf("logger factory is required but not provided")
	}

	// Extraer configuración
	env, version := extractEnvAndVersion(config)

	// Crear logger
	logger, err := factories.Logger.CreateLogger(ctx, env, version)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	resources.Logger = logger
	return nil
}
