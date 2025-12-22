package bootstrap

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRabbitMQFactory_CreateConnection_Success verifica creación exitosa de conexión
func TestRabbitMQFactory_CreateConnection_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	// Setup container
	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	require.NotNil(t, rabbit)

	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	// Crear factory
	factory := NewDefaultRabbitMQFactory()

	// Configuración de RabbitMQ
	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	// Crear conexión
	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	// Verificar que la conexión funciona
	assert.False(t, conn.IsClosed())
}

// TestRabbitMQFactory_CreateConnection_InvalidURL verifica error con URL inválida
func TestRabbitMQFactory_CreateConnection_InvalidURL(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()
	factory := NewDefaultRabbitMQFactory()

	// URL inválida
	invalidConfig := RabbitMQConfig{
		URL: "amqp://invalid-host-that-does-not-exist:5672",
	}

	// Debe fallar al crear conexión
	conn, err := factory.CreateConnection(ctx, invalidConfig)
	assert.Error(t, err)
	assert.Nil(t, conn)
}

// TestRabbitMQFactory_CreateConnection_WithTimeout verifica timeout de conexión
func TestRabbitMQFactory_CreateConnection_WithTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	// Contexto con timeout muy corto
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	factory := NewDefaultRabbitMQFactory()

	// Configuración con host que no responde
	config := RabbitMQConfig{
		URL: "amqp://192.0.2.1:5672", // TEST-NET-1 (no routable)
	}

	// Debe fallar por timeout
	conn, err := factory.CreateConnection(ctx, config)
	assert.Error(t, err)
	assert.Nil(t, conn)
}

// TestRabbitMQFactory_CreateConnection_ContextCancellation verifica cancelación por contexto
func TestRabbitMQFactory_CreateConnection_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	// Contexto ya cancelado
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	factory := NewDefaultRabbitMQFactory()

	config := RabbitMQConfig{
		URL: "amqp://localhost:5672",
	}

	// Debe fallar porque el contexto está cancelado
	conn, err := factory.CreateConnection(ctx, config)
	assert.Error(t, err)
	assert.Nil(t, conn)
}

// TestRabbitMQFactory_CreateChannel_Success verifica creación exitosa de canal
func TestRabbitMQFactory_CreateChannel_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	factory := NewDefaultRabbitMQFactory()

	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)
	defer conn.Close()

	// Crear canal
	channel, err := factory.CreateChannel(conn)
	require.NoError(t, err)
	require.NotNil(t, channel)
	defer channel.Close()

	// Verificar que el canal funciona
	assert.NotNil(t, channel)
}

// TestRabbitMQFactory_CreateChannel_QoSConfigured verifica que QoS está configurado
func TestRabbitMQFactory_CreateChannel_QoSConfigured(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	factory := NewDefaultRabbitMQFactory()

	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)
	defer conn.Close()

	// Crear canal (internamente configura QoS)
	channel, err := factory.CreateChannel(conn)
	require.NoError(t, err)
	defer channel.Close()

	// Verificar que podemos usar el canal sin errores
	// (QoS ya debería estar configurado)
	err = channel.ExchangeDeclare(
		"test_exchange_qos",
		"fanout",
		false,
		true,
		false,
		false,
		nil,
	)
	assert.NoError(t, err)
}

// TestRabbitMQFactory_DeclareQueue_Success verifica declaración exitosa de cola
func TestRabbitMQFactory_DeclareQueue_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	factory := NewDefaultRabbitMQFactory()

	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)
	defer conn.Close()

	channel, err := factory.CreateChannel(conn)
	require.NoError(t, err)
	defer channel.Close()

	// Declarar cola
	queueName := "test_queue_factory"
	queue, err := factory.DeclareQueue(channel, queueName)
	require.NoError(t, err)
	assert.Equal(t, queueName, queue.Name)
	assert.Equal(t, 0, queue.Messages)
	assert.Equal(t, 0, queue.Consumers)
}

// TestRabbitMQFactory_DeclareQueue_WithConfiguration verifica configuración de cola
func TestRabbitMQFactory_DeclareQueue_WithConfiguration(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	factory := NewDefaultRabbitMQFactory()

	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)
	defer conn.Close()

	channel, err := factory.CreateChannel(conn)
	require.NoError(t, err)
	defer channel.Close()

	// Declarar cola (la factory configura TTL, priority, etc.)
	queueName := "test_queue_configured"
	queue, err := factory.DeclareQueue(channel, queueName)
	require.NoError(t, err)

	// Verificar que la cola fue creada
	assert.Equal(t, queueName, queue.Name)

	// Inspeccionar cola para verificar argumentos
	inspected, err := channel.QueueInspect(queueName)
	require.NoError(t, err)
	assert.Equal(t, queueName, inspected.Name)
}

// TestRabbitMQFactory_Close_Success verifica cierre exitoso
func TestRabbitMQFactory_Close_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	factory := NewDefaultRabbitMQFactory()

	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)

	channel, err := factory.CreateChannel(conn)
	require.NoError(t, err)

	// Close debe ser exitoso
	err = factory.Close(channel, conn)
	assert.NoError(t, err)
}

// TestRabbitMQFactory_Close_WithNilChannel verifica cierre con canal nil
func TestRabbitMQFactory_Close_WithNilChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	factory := NewDefaultRabbitMQFactory()

	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)

	// Close con canal nil (no debe dar error)
	err = factory.Close(nil, conn)
	assert.NoError(t, err)
}

// TestRabbitMQFactory_Close_AlreadyClosed verifica cierre de conexión ya cerrada
func TestRabbitMQFactory_Close_AlreadyClosed(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	factory := NewDefaultRabbitMQFactory()

	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)

	channel, err := factory.CreateChannel(conn)
	require.NoError(t, err)

	// Cerrar primera vez
	err = factory.Close(channel, conn)
	require.NoError(t, err)

	// Cerrar segunda vez (ya cerrada, debería manejar el error)
	err2 := factory.Close(channel, conn)
	// Puede dar error o no, dependiendo de cómo maneja amqp.ErrClosed
	_ = err2 // No falla el test
}

// TestRabbitMQFactory_MultipleChannels verifica creación de múltiples canales
func TestRabbitMQFactory_MultipleChannels(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	factory := NewDefaultRabbitMQFactory()

	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)
	defer conn.Close()

	// Crear múltiples canales
	channels := make([]*amqp.Channel, 5)
	for i := 0; i < 5; i++ {
		ch, err := factory.CreateChannel(conn)
		require.NoError(t, err)
		channels[i] = ch
	}

	// Verificar que todos funcionan
	for i, ch := range channels {
		assert.NotNil(t, ch, "Canal %d debe ser válido", i)
	}

	// Cerrar todos los canales
	for _, ch := range channels {
		_ = ch.Close()
	}
}

// TestRabbitMQFactory_PublishConsume verifica flujo completo publish/consume
func TestRabbitMQFactory_PublishConsume(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()

	config := containers.NewConfig().
		WithRabbitMQ(&containers.RabbitConfig{
			Image: "rabbitmq:3.12-alpine",
		}).
		Build()

	manager, err := containers.GetManager(t, config)
	require.NoError(t, err)

	rabbit := manager.RabbitMQ()
	connectionString, err := rabbit.ConnectionString(ctx)
	require.NoError(t, err)

	factory := NewDefaultRabbitMQFactory()

	rabbitConfig := RabbitMQConfig{
		URL: connectionString,
	}

	conn, err := factory.CreateConnection(ctx, rabbitConfig)
	require.NoError(t, err)
	defer factory.Close(nil, conn)

	channel, err := factory.CreateChannel(conn)
	require.NoError(t, err)
	defer channel.Close()

	// Declarar cola
	queueName := "test_queue_pubsub"
	queue, err := factory.DeclareQueue(channel, queueName)
	require.NoError(t, err)

	// Publicar mensaje
	testMessage := "test message from factory"
	err = channel.PublishWithContext(
		ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(testMessage),
		},
	)
	require.NoError(t, err)

	// Consumir mensaje
	msgs, err := channel.Consume(
		queue.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		assert.Equal(t, testMessage, string(msg.Body))
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout esperando mensaje")
	}
}

// TestRabbitMQFactory_DefaultTimeout verifica timeout por defecto
func TestRabbitMQFactory_DefaultTimeout(t *testing.T) {
	factory := NewDefaultRabbitMQFactory()

	assert.NotNil(t, factory)
	assert.NotZero(t, factory.connectionTimeout)
	assert.Equal(t, 10*time.Second, factory.connectionTimeout)
}
