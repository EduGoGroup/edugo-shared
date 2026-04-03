package health

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQChannel define la interfaz para operaciones de RabbitMQ necesarias para health checks
type RabbitMQChannel interface {
	IsClosed() bool
}

// RabbitMQCheck implementa HealthCheck para RabbitMQ
type RabbitMQCheck struct {
	channel RabbitMQChannel
	timeout time.Duration
}

// NewRabbitMQCheck crea un nuevo RabbitMQ health check
// Panics if channel is nil. Defaults timeout to 5s if <= 0.
func NewRabbitMQCheck(channel *amqp.Channel, timeout time.Duration) *RabbitMQCheck {
	if channel == nil {
		panic("health: RabbitMQ channel must not be nil")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &RabbitMQCheck{
		channel: channel,
		timeout: timeout,
	}
}

// NewRabbitMQCheckWithChannel crea un health check con una interfaz RabbitMQChannel
// Útil para testing con mocks
// Panics if channel is nil. Defaults timeout to 5s if <= 0.
func NewRabbitMQCheckWithChannel(channel RabbitMQChannel, timeout time.Duration) *RabbitMQCheck {
	if channel == nil {
		panic("health: RabbitMQ channel must not be nil")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &RabbitMQCheck{
		channel: channel,
		timeout: timeout,
	}
}

// Name retorna el nombre del health check
func (c *RabbitMQCheck) Name() string {
	return "rabbitmq"
}

// Check ejecuta el health check de RabbitMQ
func (c *RabbitMQCheck) Check(ctx context.Context) CheckResult {
	result := CheckResult{
		Component: c.Name(),
		Timestamp: time.Now(),
		Metadata:  make(map[string]any),
	}

	// Crear contexto con timeout
	checkCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	start := time.Now()

	// Ejecutar IsClosed() en goroutine para respetar timeout/contexto
	type checkResult struct {
		closed bool
	}
	ch := make(chan checkResult, 1)
	go func() {
		ch <- checkResult{closed: c.channel.IsClosed()}
	}()

	select {
	case <-checkCtx.Done():
		result.Status = StatusUnhealthy
		result.Message = "RabbitMQ health check timed out"
		return result
	case res := <-ch:
		duration := time.Since(start)
		result.Metadata["response_time_ms"] = duration.Milliseconds()

		if res.closed {
			result.Status = StatusUnhealthy
			result.Message = "RabbitMQ channel is closed"
			return result
		}
	}

	result.Status = StatusHealthy
	result.Message = "RabbitMQ is healthy"
	return result
}
