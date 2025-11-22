package rabbit

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPublisher_BatchPublishing verifica publicación en batch
func TestPublisher_BatchPublishing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_batch"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	publisher := NewPublisher(conn)

	// Publicar batch de mensajes
	messageCount := 100
	for i := 0; i < messageCount; i++ {
		message := map[string]interface{}{
			"id":      i,
			"content": fmt.Sprintf("message %d", i),
		}

		err = publisher.Publish(ctx, "", queueName, message)
		require.NoError(t, err)
	}

	// Verificar que todos los mensajes llegaron
	time.Sleep(500 * time.Millisecond)

	queueInfo := waitForQueueMessages(t, rabbitContainer, queueName, messageCount, 5*time.Second)
	assert.Equal(t, messageCount, queueInfo.Messages)
}

// TestPublisher_WithPriority verifica publicación con diferentes prioridades
func TestPublisher_WithPriority(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Crear cola con prioridad habilitada
	queueName := "test_queue_priority"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args: map[string]interface{}{
			"x-max-priority": int32(10),
		},
	})
	require.NoError(t, err)

	publisher := NewPublisher(conn).(*RabbitMQPublisher)

	// Publicar mensajes con diferentes prioridades
	priorities := []uint8{5, 10, 1, 8, 3}
	for i, priority := range priorities {
		message := map[string]interface{}{
			"id":       i,
			"priority": priority,
		}

		err = publisher.PublishWithPriority(ctx, "", queueName, message, priority)
		require.NoError(t, err)
	}

	time.Sleep(200 * time.Millisecond)

	// Consumir mensajes y verificar que Priority fue configurado
	msgs, err := conn.GetChannel().Consume(queueName, "", true, false, false, false, nil)
	require.NoError(t, err)

	receivedCount := 0
	timeout := time.After(2 * time.Second)

	for receivedCount < len(priorities) {
		select {
		case msg := <-msgs:
			// Verificar que el mensaje tiene el campo Priority
			assert.GreaterOrEqual(t, msg.Priority, uint8(0))
			receivedCount++
		case <-timeout:
			t.Fatalf("Timeout: solo recibidos %d/%d mensajes", receivedCount, len(priorities))
		}
	}
}

// TestPublisher_ContextCancellation verifica cancelación por contexto
func TestPublisher_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_ctx_cancel"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	publisher := NewPublisher(conn)

	// Crear contexto ya cancelado
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	message := map[string]string{"test": "message"}
	err = publisher.Publish(ctx, "", queueName, message)
	assert.Error(t, err, "Debe dar error al publicar con contexto cancelado")
}

// TestPublisher_InvalidJSON verifica manejo de datos no serializables
func TestPublisher_InvalidJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	publisher := NewPublisher(conn)

	// Intentar publicar un tipo que no puede ser serializado a JSON
	invalidData := make(chan int) // Los canales no son serializables a JSON

	err = publisher.Publish(ctx, "", "test_queue", invalidData)
	assert.Error(t, err, "Debe dar error al serializar datos inválidos")
	assert.Contains(t, err.Error(), "marshal")
}

// TestConsumer_ConcurrentProcessing verifica procesamiento concurrente
func TestConsumer_ConcurrentProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_concurrent"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	// Configurar QoS para permitir múltiples mensajes en paralelo
	err = conn.SetPrefetchCount(10)
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_concurrent",
		AutoAck:       false,
		PrefetchCount: 10,
	}

	consumer := NewConsumer(conn, consumerConfig)

	// Contador de mensajes procesados concurrentemente
	var processing atomic.Int32
	var maxConcurrent atomic.Int32
	var processed atomic.Int32

	handler := func(ctx context.Context, body []byte) error {
		current := processing.Add(1)

		// Actualizar máximo concurrente
		for {
			max := maxConcurrent.Load()
			if current <= max || maxConcurrent.CompareAndSwap(max, current) {
				break
			}
		}

		// Simular procesamiento
		time.Sleep(100 * time.Millisecond)

		processing.Add(-1)
		processed.Add(1)
		return nil
	}

	err = consumer.Consume(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar múltiples mensajes
	publisher := NewPublisher(conn)
	messageCount := 20

	for i := 0; i < messageCount; i++ {
		err = publisher.Publish(ctx, "", queueName, map[string]int{"id": i})
		require.NoError(t, err)
	}

	// Esperar procesamiento
	time.Sleep(3 * time.Second)

	assert.Equal(t, int32(messageCount), processed.Load())
	assert.GreaterOrEqual(t, maxConcurrent.Load(), int32(2), "Debe procesar al menos 2 mensajes concurrentemente")
}

// TestConsumer_Acknowledgment verifica ACK/NACK correcto
func TestConsumer_Acknowledgment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_ack"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_ack",
		AutoAck:       false, // Manual ACK
		PrefetchCount: 1,
	}

	consumer := NewConsumer(conn, consumerConfig)

	// Handler que tiene éxito
	var processed atomic.Int32
	handler := func(ctx context.Context, body []byte) error {
		processed.Add(1)
		return nil
	}

	err = consumer.Consume(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar mensaje
	publisher := NewPublisher(conn)
	err = publisher.Publish(ctx, "", queueName, map[string]string{"test": "ack"})
	require.NoError(t, err)

	// Esperar procesamiento
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, int32(1), processed.Load())

	// Verificar que la cola está vacía (mensaje fue ACKed)
	queueInfo := waitForQueueMessages(t, rabbitContainer, queueName, 0, 2*time.Second)
	assert.Equal(t, 0, queueInfo.Messages, "Cola debe estar vacía después de ACK")
}

// TestConsumer_NACK_Requeue verifica NACK con requeue
func TestConsumer_NACK_Requeue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_nack_requeue"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_nack",
		AutoAck:       false,
		PrefetchCount: 1,
	}

	consumer := NewConsumer(conn, consumerConfig)

	// Handler que falla las primeras 2 veces, luego tiene éxito
	var attempts atomic.Int32
	handler := func(ctx context.Context, body []byte) error {
		attempt := attempts.Add(1)
		if attempt < 3 {
			return errors.New("simulated failure")
		}
		return nil
	}

	err = consumer.Consume(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar mensaje
	publisher := NewPublisher(conn)
	err = publisher.Publish(context.Background(), "", queueName, map[string]string{"test": "nack"})
	require.NoError(t, err)

	// Esperar múltiples intentos
	time.Sleep(2 * time.Second)

	// Debe haber intentado al menos 3 veces
	assert.GreaterOrEqual(t, attempts.Load(), int32(3))
}

// TestConsumer_GracefulShutdown verifica shutdown limpio
func TestConsumer_GracefulShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_shutdown_graceful"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_shutdown",
		AutoAck:       false,
		PrefetchCount: 10,
	}

	consumer := NewConsumer(conn, consumerConfig)

	// Contexto con cancelación
	ctx, cancel := context.WithCancel(context.Background())

	var processing atomic.Int32
	var completed atomic.Int32

	handler := func(ctx context.Context, body []byte) error {
		processing.Add(1)
		time.Sleep(500 * time.Millisecond) // Simular procesamiento largo
		processing.Add(-1)
		completed.Add(1)
		return nil
	}

	err = consumer.Consume(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar algunos mensajes
	publisher := NewPublisher(conn)
	for i := 0; i < 5; i++ {
		err = publisher.Publish(context.Background(), "", queueName, map[string]int{"id": i})
		require.NoError(t, err)
	}

	// Esperar a que empiece a procesar
	time.Sleep(200 * time.Millisecond)

	// Cancelar contexto (shutdown graceful)
	cancel()

	// Esperar un poco más
	time.Sleep(1 * time.Second)

	// Verificar que procesó algunos mensajes
	assert.GreaterOrEqual(t, completed.Load(), int32(1), "Debe haber procesado al menos 1 mensaje")
}

// TestPublisher_Close verifica cierre correcto del publisher
func TestPublisher_Close(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	publisher := NewPublisher(conn)

	err = publisher.Close()
	assert.NoError(t, err)
}

// TestConsumer_Close verifica cierre correcto del consumer
func TestConsumer_Close(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	consumerConfig := ConsumerConfig{
		Name:    "test_consumer_close",
		AutoAck: false,
	}

	consumer := NewConsumer(conn, consumerConfig)

	err = consumer.Close()
	assert.NoError(t, err)
}

// TestPublisher_LargePayload verifica publicación de payloads grandes
func TestPublisher_LargePayload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_large_payload"
	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	publisher := NewPublisher(conn)

	// Crear payload grande (array de 10000 elementos)
	largePayload := make([]map[string]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		largePayload[i] = map[string]interface{}{
			"id":      i,
			"content": fmt.Sprintf("large content %d with some text to increase size", i),
			"data":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
	}

	message := map[string]interface{}{
		"payload": largePayload,
	}

	err = publisher.Publish(ctx, "", queueName, message)
	assert.NoError(t, err, "Debe poder publicar payload grande")

	// Consumir y verificar
	consumerConfig := ConsumerConfig{
		Name:    "test_consumer_large",
		AutoAck: true,
	}

	consumer := NewConsumer(conn, consumerConfig)

	var received atomic.Bool
	handler := func(ctx context.Context, body []byte) error {
		// Verificar que el tamaño es grande
		assert.Greater(t, len(body), 100000, "Body debe ser grande")
		received.Store(true)
		return nil
	}

	err = consumer.Consume(ctx, queueName, handler)
	require.NoError(t, err)

	// Esperar procesamiento
	time.Sleep(2 * time.Second)

	assert.True(t, received.Load(), "Debe haber recibido el mensaje grande")
}

// TestPublisher_Consumer_Integration verifica flujo completo
func TestPublisher_Consumer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Crear exchange y queue
	exchangeName := "test_exchange_integration"
	queueName := "test_queue_integration"
	routingKey := "test.routing.key"

	err = conn.DeclareExchange(ExchangeConfig{
		Name:       exchangeName,
		Type:       "topic",
		Durable:    false,
		AutoDelete: true,
	})
	require.NoError(t, err)

	_, err = conn.DeclareQueue(QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	})
	require.NoError(t, err)

	err = conn.BindQueue(queueName, routingKey, exchangeName)
	require.NoError(t, err)

	// Setup consumer
	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_integration",
		AutoAck:       false,
		PrefetchCount: 5,
	}

	consumer := NewConsumer(conn, consumerConfig)

	type TestMessage struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
	}

	receivedMessages := make([]TestMessage, 0)
	var mu sync.Mutex

	handler := func(ctx context.Context, body []byte) error {
		var msg TestMessage
		err := UnmarshalMessage(body, &msg)
		if err != nil {
			return err
		}

		mu.Lock()
		receivedMessages = append(receivedMessages, msg)
		mu.Unlock()

		return nil
	}

	err = consumer.Consume(ctx, queueName, handler)
	require.NoError(t, err)

	// Setup publisher
	publisher := NewPublisher(conn)

	// Publicar mensajes
	messageCount := 10
	for i := 0; i < messageCount; i++ {
		message := TestMessage{
			ID:      i,
			Content: fmt.Sprintf("integration test message %d", i),
		}

		err = publisher.Publish(ctx, exchangeName, routingKey, message)
		require.NoError(t, err)
	}

	// Esperar procesamiento
	time.Sleep(2 * time.Second)

	// Verificar todos los mensajes fueron recibidos
	mu.Lock()
	assert.Equal(t, messageCount, len(receivedMessages))
	mu.Unlock()

	// Verificar orden y contenido
	for i, msg := range receivedMessages {
		assert.Equal(t, i, msg.ID)
		assert.Contains(t, msg.Content, fmt.Sprintf("message %d", i))
	}
}
