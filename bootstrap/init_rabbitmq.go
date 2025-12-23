package bootstrap

import (
	"context"
	"errors"
	"fmt"
)

// initRabbitMQ inicializa la conexión a RabbitMQ.
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
func initRabbitMQ(
	ctx context.Context,
	config interface{},
	factories *Factories,
	resources *Resources,
	lifecycleManager interface{},
	opts *BootstrapOptions,
) error {
	if factories.RabbitMQ == nil {
		return fmt.Errorf("rabbitmq factory not provided")
	}

	// Extraer configuración de RabbitMQ
	rabbitConfig, err := extractRabbitMQConfig(config)
	if err != nil {
		return fmt.Errorf("failed to extract RabbitMQ config: %w", err)
	}

	// Log inicio
	if resources.Logger != nil {
		resources.Logger.Info("Initializing RabbitMQ connection...")
	}

	// Crear conexión
	conn, err := factories.RabbitMQ.CreateConnection(ctx, rabbitConfig)
	if err != nil {
		return fmt.Errorf("failed to create RabbitMQ connection: %w", err)
	}

	// Crear canal
	channel, err := factories.RabbitMQ.CreateChannel(conn)
	if err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			return errors.Join(
				fmt.Errorf("failed to create RabbitMQ channel: %w", err),
				fmt.Errorf("failed to close connection: %w", closeErr),
			)
		}
		return fmt.Errorf("failed to create RabbitMQ channel: %w", err)
	}

	// Crear MessagePublisher (implementación simple por ahora)
	resources.MessagePublisher = &defaultMessagePublisher{
		channel: channel,
		factory: factories.RabbitMQ,
	}

	// Registrar cleanup en lifecycle manager si está disponible
	if lifecycleManager != nil {
		registerRabbitMQCleanup(lifecycleManager, factories.RabbitMQ, channel, conn, resources.Logger)
	}

	// Log éxito
	if resources.Logger != nil {
		resources.Logger.Info("RabbitMQ connection established")
	}

	return nil
}
