package containers

import (
	"context"
	"fmt"
)

// RabbitMQContainer envuelve el container de RabbitMQ
type RabbitMQContainer struct {
	// container *rabbitmq.RabbitMQContainer
	config *RabbitConfig
}

// createRabbitMQ crea y configura un container de RabbitMQ
func createRabbitMQ(ctx context.Context, cfg *RabbitConfig) (*RabbitMQContainer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("RabbitConfig no puede ser nil")
	}

	// TODO: Implementar creación de container RabbitMQ
	// Por ahora retorna un stub

	return &RabbitMQContainer{
		config: cfg,
	}, nil
}

// PurgeAll elimina todas las colas y exchanges
func (rc *RabbitMQContainer) PurgeAll(ctx context.Context) error {
	// TODO: Implementar
	return nil
}

// Terminate termina el container
func (rc *RabbitMQContainer) Terminate(ctx context.Context) error {
	// TODO: Implementar
	return nil
}
