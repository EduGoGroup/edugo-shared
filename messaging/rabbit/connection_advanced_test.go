package rabbit

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConnection_MultipleChannels verifica creación de múltiples canales
func TestConnection_MultipleChannels(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Crear múltiples canales concurrentemente
	channelCount := 10
	channels := make([]*amqp.Channel, channelCount)
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make([]error, 0)

	for i := 0; i < channelCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			ch, err := conn.GetConnection().Channel()
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			channels[index] = ch
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 0, len(errors), "No debe haber errores al crear canales")

	// Verificar que todos los canales son válidos
	for i, ch := range channels {
		assert.NotNil(t, ch, "Canal %d debe ser válido", i)
		if ch != nil {
			_ = ch.Close()
		}
	}
}

// TestConnection_ConcurrentOperations verifica operaciones concurrentes
func TestConnection_ConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Operaciones concurrentes: declarar colas
	queueCount := 20
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make([]error, 0)

	for i := 0; i < queueCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Cada goroutine crea su propio canal
			ch, err := conn.GetConnection().Channel()
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}
			defer ch.Close()

			queueName := fmt.Sprintf("test_concurrent_queue_%d", index)
			_, err = ch.QueueDeclare(
				queueName,
				false, // no durable para facilitar limpieza
				true,  // auto-delete
				false,
				false,
				nil,
			)

			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 0, len(errors), "No debe haber errores en operaciones concurrentes")
}

// TestConnection_HealthCheck_Success verifica health check exitoso
func TestConnection_HealthCheck_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Health check debe ser exitoso
	err = conn.HealthCheck()
	assert.NoError(t, err)
}

// TestConnection_HealthCheck_AfterClose verifica health check falla después de cerrar
func TestConnection_HealthCheck_AfterClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)

	// Cerrar conexión
	err = conn.Close()
	require.NoError(t, err)

	// Health check debe fallar
	err = conn.HealthCheck()
	assert.Error(t, err, "Health check debe fallar después de cerrar conexión")
	assert.Contains(t, err.Error(), "closed")
}

// TestConnection_SetPrefetchCount verifica configuración de prefetch
func TestConnection_SetPrefetchCount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	tests := []struct {
		name         string
		prefetchCount int
		expectError  bool
	}{
		{
			name:          "prefetch 1",
			prefetchCount: 1,
			expectError:   false,
		},
		{
			name:          "prefetch 10",
			prefetchCount: 10,
			expectError:   false,
		},
		{
			name:          "prefetch 100",
			prefetchCount: 100,
			expectError:   false,
		},
		{
			name:          "prefetch 0",
			prefetchCount: 0,
			expectError:   false, // 0 es válido (sin límite)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conn.SetPrefetchCount(tt.prefetchCount)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestConnection_DeclareExchange_VariousTypes verifica declaración de exchanges
func TestConnection_DeclareExchange_VariousTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	tests := []struct {
		name         string
		exchangeType string
		durable      bool
		autoDelete   bool
	}{
		{
			name:         "direct exchange",
			exchangeType: "direct",
			durable:      true,
			autoDelete:   false,
		},
		{
			name:         "topic exchange",
			exchangeType: "topic",
			durable:      true,
			autoDelete:   false,
		},
		{
			name:         "fanout exchange",
			exchangeType: "fanout",
			durable:      false,
			autoDelete:   true,
		},
		{
			name:         "headers exchange",
			exchangeType: "headers",
			durable:      true,
			autoDelete:   false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ExchangeConfig{
				Name:       fmt.Sprintf("test_exchange_%s_%d", tt.exchangeType, i),
				Type:       tt.exchangeType,
				Durable:    tt.durable,
				AutoDelete: tt.autoDelete,
			}

			err := conn.DeclareExchange(cfg)
			assert.NoError(t, err)
		})
	}
}

// TestConnection_DeclareQueue_VariousConfigs verifica declaración de colas
func TestConnection_DeclareQueue_VariousConfigs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	tests := []struct {
		name       string
		durable    bool
		autoDelete bool
		exclusive  bool
		args       map[string]interface{}
	}{
		{
			name:       "durable queue",
			durable:    true,
			autoDelete: false,
			exclusive:  false,
			args:       nil,
		},
		{
			name:       "temporary queue",
			durable:    false,
			autoDelete: true,
			exclusive:  false,
			args:       nil,
		},
		{
			name:       "exclusive queue",
			durable:    false,
			autoDelete: true,
			exclusive:  true,
			args:       nil,
		},
		{
			name:       "queue with TTL",
			durable:    true,
			autoDelete: false,
			exclusive:  false,
			args: map[string]interface{}{
				"x-message-ttl": int32(60000), // 60 segundos
			},
		},
		{
			name:       "queue with max length",
			durable:    true,
			autoDelete: false,
			exclusive:  false,
			args: map[string]interface{}{
				"x-max-length": int32(1000),
			},
		},
		{
			name:       "queue with priority",
			durable:    true,
			autoDelete: false,
			exclusive:  false,
			args: map[string]interface{}{
				"x-max-priority": int32(10),
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := QueueConfig{
				Name:       fmt.Sprintf("test_queue_config_%d", i),
				Durable:    tt.durable,
				AutoDelete: tt.autoDelete,
				Exclusive:  tt.exclusive,
				Args:       tt.args,
			}

			queue, err := conn.DeclareQueue(cfg)
			assert.NoError(t, err)
			assert.NotEmpty(t, queue.Name)
		})
	}
}

// TestConnection_BindQueue verifica binding de cola a exchange
func TestConnection_BindQueue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Crear exchange
	exchangeName := "test_exchange_bind"
	err = conn.DeclareExchange(ExchangeConfig{
		Name:       exchangeName,
		Type:       "topic",
		Durable:    false,
		AutoDelete: true,
	})
	require.NoError(t, err)

	// Crear queue
	queueName := "test_queue_bind"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	// Bind queue a exchange
	routingKey := "test.routing.key"
	err = conn.BindQueue(queueName, routingKey, exchangeName)
	assert.NoError(t, err)

	// Publicar mensaje al exchange
	ctx := context.Background()
	err = conn.GetChannel().PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("test message"),
		},
	)
	require.NoError(t, err)

	// Consumir de la cola para verificar binding
	msgs, err := conn.GetChannel().Consume(
		queueName,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		assert.Equal(t, []byte("test message"), msg.Body)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout esperando mensaje (binding puede haber fallado)")
	}
}

// TestConnection_IsClosed verifica detección de estado cerrado
func TestConnection_IsClosed(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)

	// Antes de cerrar, debe estar abierta
	assert.False(t, conn.IsClosed())

	// Cerrar conexión
	err = conn.Close()
	require.NoError(t, err)

	// Después de cerrar, debe estar cerrada
	assert.True(t, conn.IsClosed())
}

// TestConnection_GetConnection verifica obtención de conexión raw
func TestConnection_GetConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	rawConn := conn.GetConnection()
	assert.NotNil(t, rawConn)
	assert.False(t, rawConn.IsClosed())
}

// TestConnection_GetChannel verifica obtención de canal
func TestConnection_GetChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	channel := conn.GetChannel()
	assert.NotNil(t, channel)

	// Usar canal para verificar que funciona
	err = channel.ExchangeDeclare(
		"test_channel_verify",
		"fanout",
		false,
		true,
		false,
		false,
		nil,
	)
	assert.NoError(t, err)
}

// TestConnection_DoubleClose verifica que cerrar dos veces no causa panic
func TestConnection_DoubleClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)

	// Primera vez
	err = conn.Close()
	assert.NoError(t, err)

	// Segunda vez (puede dar error pero no debe panic)
	err2 := conn.Close()
	// Puede ser error o nil, dependiendo de la implementación de amqp
	// Lo importante es que no haya panic
	_ = err2
}

// TestConnection_LargeMessageHandling verifica manejo de mensajes grandes
func TestConnection_LargeMessageHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_large_message"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	// Crear mensaje grande (1MB)
	largeMessage := make([]byte, 1024*1024)
	for i := range largeMessage {
		largeMessage[i] = byte(i % 256)
	}

	ctx := context.Background()

	// Publicar mensaje grande
	err = conn.GetChannel().PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/octet-stream",
			Body:         largeMessage,
			DeliveryMode: amqp.Persistent,
		},
	)
	assert.NoError(t, err, "Debe poder publicar mensaje grande")

	// Consumir mensaje grande
	msgs, err := conn.GetChannel().Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		assert.Equal(t, len(largeMessage), len(msg.Body))
		assert.Equal(t, largeMessage, msg.Body)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout esperando mensaje grande")
	}
}
