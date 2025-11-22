package rabbit

import (
	"context"
	"fmt"
	"testing"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRabbitContainer obtiene el container de RabbitMQ desde el Manager
func setupRabbitContainer(t *testing.T) (*containers.RabbitMQContainer, string) {
	t.Helper()

	ctx := context.Background()

	// Configurar manager con RabbitMQ habilitado
	config := containers.NewConfig().
		WithRabbitMQ(nil).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err, "Error creando manager de containers")

	rabbitContainer := manager.RabbitMQ()
	require.NotNil(t, rabbitContainer, "RabbitMQ container no disponible")

	connectionString, err := rabbitContainer.ConnectionString(ctx)
	require.NoError(t, err, "Error obteniendo connection string")

	return rabbitContainer, connectionString
}

func TestConnect_Success(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() { _ = conn.Close() }()

	assert.NotNil(t, conn.conn)
	assert.NotNil(t, conn.channel)
	assert.Equal(t, connectionString, conn.url)
}

func TestConnect_InvalidURL(t *testing.T) {
	conn, err := Connect("amqp://invalid:invalid@localhost:9999/")
	assert.Error(t, err)
	assert.Nil(t, conn)
	assert.Contains(t, err.Error(), "failed to connect to RabbitMQ")
}

func TestConnection_GetChannel(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	channel := conn.GetChannel()
	assert.NotNil(t, channel)
	assert.IsType(t, &amqp.Channel{}, channel)
}

func TestConnection_GetConnection(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	amqpConn := conn.GetConnection()
	assert.NotNil(t, amqpConn)
	assert.IsType(t, &amqp.Connection{}, amqpConn)
	assert.False(t, amqpConn.IsClosed())
}

func TestConnection_Close(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.Close()
	assert.NoError(t, err)
	assert.True(t, conn.IsClosed())
}

func TestConnection_IsClosed(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	require.NotNil(t, conn)

	// Inicialmente no debe estar cerrada
	assert.False(t, conn.IsClosed())

	// Cerrar y verificar
	err = conn.Close()
	require.NoError(t, err)
	assert.True(t, conn.IsClosed())
}

func TestConnection_DeclareExchange(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	tests := []struct {
		name   string
		config ExchangeConfig
	}{
		{
			name: "direct exchange",
			config: ExchangeConfig{
				Name:       "test_direct",
				Type:       "direct",
				Durable:    true,
				AutoDelete: false,
			},
		},
		{
			name: "topic exchange",
			config: ExchangeConfig{
				Name:       "test_topic",
				Type:       "topic",
				Durable:    true,
				AutoDelete: false,
			},
		},
		{
			name: "fanout exchange",
			config: ExchangeConfig{
				Name:       "test_fanout",
				Type:       "fanout",
				Durable:    false,
				AutoDelete: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conn.DeclareExchange(tt.config)
			assert.NoError(t, err)
		})
	}
}

func TestConnection_DeclareQueue(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	tests := []struct {
		name   string
		config QueueConfig
	}{
		{
			name: "durable queue",
			config: QueueConfig{
				Name:       "test_durable",
				Durable:    true,
				AutoDelete: false,
				Exclusive:  false,
				Args:       nil,
			},
		},
		{
			name: "auto-delete queue",
			config: QueueConfig{
				Name:       "test_autodelete",
				Durable:    false,
				AutoDelete: true,
				Exclusive:  false,
				Args:       nil,
			},
		},
		{
			name: "queue with TTL",
			config: QueueConfig{
				Name:       "test_ttl",
				Durable:    true,
				AutoDelete: false,
				Exclusive:  false,
				Args: map[string]interface{}{
					"x-message-ttl": 60000,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue, err := conn.DeclareQueue(tt.config)
			assert.NoError(t, err)
			assert.Equal(t, tt.config.Name, queue.Name)
		})
	}
}

func TestConnection_BindQueue(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Declarar exchange
	exchangeConfig := ExchangeConfig{
		Name:       "test_bind_exchange",
		Type:       "direct",
		Durable:    false,
		AutoDelete: true,
	}
	err = conn.DeclareExchange(exchangeConfig)
	require.NoError(t, err)

	// Declarar queue
	queueConfig := QueueConfig{
		Name:       "test_bind_queue",
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	queue, err := conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Bind queue to exchange
	err = conn.BindQueue(queue.Name, "test_routing_key", exchangeConfig.Name)
	assert.NoError(t, err)
}

func TestConnection_SetPrefetchCount(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	tests := []struct {
		name  string
		count int
	}{
		{"prefetch 1", 1},
		{"prefetch 5", 5},
		{"prefetch 10", 10},
		{"prefetch 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conn.SetPrefetchCount(tt.count)
			assert.NoError(t, err)
		})
	}
}

func TestConnection_HealthCheck_Success(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	err = conn.HealthCheck()
	assert.NoError(t, err)
}

func TestConnection_HealthCheck_ClosedConnection(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)

	// Cerrar conexión
	err = conn.Close()
	require.NoError(t, err)

	// Health check debe fallar
	err = conn.HealthCheck()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection is closed")
}

func TestConnection_MultipleChannels(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Obtener el canal principal
	channel1 := conn.GetChannel()
	assert.NotNil(t, channel1)

	// El canal principal debe ser el mismo
	channel2 := conn.GetChannel()
	assert.Equal(t, channel1, channel2)
}

func TestConnection_DeclareQueue_WithDLX(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Declarar DLX exchange primero
	dlxConfig := ExchangeConfig{
		Name:       "test_dlx",
		Type:       "direct",
		Durable:    true,
		AutoDelete: false,
	}
	err = conn.DeclareExchange(dlxConfig)
	require.NoError(t, err)

	// Declarar queue con DLX
	queueConfig := QueueConfig{
		Name:       "test_queue_with_dlx",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		Args: map[string]interface{}{
			"x-dead-letter-exchange":    "test_dlx",
			"x-dead-letter-routing-key": "dlq",
		},
	}
	queue, err := conn.DeclareQueue(queueConfig)
	assert.NoError(t, err)
	assert.Equal(t, queueConfig.Name, queue.Name)
}

func TestConnection_CloseWithNilChannel(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)

	// Forzar canal a nil para test edge case
	conn.channel = nil

	err = conn.Close()
	assert.NoError(t, err)
}

func TestConnection_CloseWithNilConnection(t *testing.T) {
	// Crear una conexión parcialmente inicializada
	conn := &Connection{
		conn:    nil,
		channel: nil,
		url:     "amqp://test",
	}

	err := conn.Close()
	assert.NoError(t, err)
	assert.True(t, conn.IsClosed())
}

func TestConnection_ConcurrentHealthChecks(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Ejecutar múltiples health checks concurrentemente
	const concurrency = 10
	errChan := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			errChan <- conn.HealthCheck()
		}()
	}

	// Verificar que todos los health checks pasaron
	for i := 0; i < concurrency; i++ {
		err := <-errChan
		assert.NoError(t, err, fmt.Sprintf("Health check %d falló", i))
	}
}

func TestConnection_ReconnectAfterClose(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	// Primera conexión
	conn1, err := Connect(connectionString)
	require.NoError(t, err)
	require.NotNil(t, conn1)

	err = conn1.Close()
	require.NoError(t, err)

	// Segunda conexión debe funcionar
	conn2, err := Connect(connectionString)
	require.NoError(t, err)
	require.NotNil(t, conn2)
	defer func() { _ = conn2.Close() }()

	assert.False(t, conn2.IsClosed())
}

func TestConnection_HealthCheck_CreatesTemporaryExchange(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Ejecutar health check
	err = conn.HealthCheck()
	require.NoError(t, err)

	// El exchange temporal debe haberse creado y eliminado automáticamente
	// No podemos verificar directamente, pero podemos ejecutar otro health check
	err = conn.HealthCheck()
	assert.NoError(t, err, "El segundo health check debe funcionar")
}

func TestConnection_BindQueue_WithMultipleRoutingKeys(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Declarar exchange
	exchangeConfig := ExchangeConfig{
		Name:       "test_multi_bind",
		Type:       "topic",
		Durable:    false,
		AutoDelete: true,
	}
	err = conn.DeclareExchange(exchangeConfig)
	require.NoError(t, err)

	// Declarar queue
	queueConfig := QueueConfig{
		Name:       "test_multi_bind_queue",
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	queue, err := conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Bind con múltiples routing keys
	routingKeys := []string{"test.key.1", "test.key.2", "test.key.*"}
	for _, key := range routingKeys {
		err = conn.BindQueue(queue.Name, key, exchangeConfig.Name)
		assert.NoError(t, err, fmt.Sprintf("Error binding queue con routing key %s", key))
	}
}

func TestConnection_SetPrefetchCount_Zero(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Prefetch count 0 significa "sin límite"
	err = conn.SetPrefetchCount(0)
	assert.NoError(t, err)
}

func TestConnection_DeclareExchange_DuplicateName(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	config := ExchangeConfig{
		Name:       "test_duplicate",
		Type:       "direct",
		Durable:    true,
		AutoDelete: false,
	}

	// Primera declaración
	err = conn.DeclareExchange(config)
	require.NoError(t, err)

	// Segunda declaración con misma configuración debe funcionar
	err = conn.DeclareExchange(config)
	assert.NoError(t, err)
}

func TestConnection_Lifecycle(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	// Test del ciclo completo de vida de una conexión
	conn, err := Connect(connectionString)
	require.NoError(t, err)
	require.NotNil(t, conn)

	// Verificar estado inicial
	assert.False(t, conn.IsClosed())
	assert.NotNil(t, conn.GetChannel())
	assert.NotNil(t, conn.GetConnection())

	// Usar la conexión
	err = conn.HealthCheck()
	require.NoError(t, err)

	// Declarar exchange
	err = conn.DeclareExchange(ExchangeConfig{
		Name:       "lifecycle_test",
		Type:       "direct",
		Durable:    false,
		AutoDelete: true,
	})
	require.NoError(t, err)

	// Declarar queue
	queue, err := conn.DeclareQueue(QueueConfig{
		Name:       "lifecycle_queue",
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, queue.Name)

	// Bind queue
	err = conn.BindQueue(queue.Name, "test_key", "lifecycle_test")
	require.NoError(t, err)

	// Set prefetch
	err = conn.SetPrefetchCount(10)
	require.NoError(t, err)

	// Health check antes de cerrar
	err = conn.HealthCheck()
	require.NoError(t, err)

	// Cerrar
	err = conn.Close()
	require.NoError(t, err)
	assert.True(t, conn.IsClosed())

	// Health check debe fallar después de cerrar
	err = conn.HealthCheck()
	assert.Error(t, err)
}
