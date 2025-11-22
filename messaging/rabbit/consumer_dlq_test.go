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

// TestConsumeWithDLQ_BasicFunctionality verifica que el consumer DLQ básico funciona
func TestConsumeWithDLQ_BasicFunctionality(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Configurar DLQ
	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            3,
		RetryDelay:            100 * time.Millisecond,
		DLXExchange:           "test.dlx",
		DLXRoutingKey:         "test.dlq",
		UseExponentialBackoff: false,
	}

	// Crear cola de test
	queueName := "test_queue_dlq_basic"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)

	// Crear consumer con DLQ
	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_dlq",
		AutoAck:       false,
		PrefetchCount: 5,
		DLQ:           dlqConfig,
	}

	consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)

	// Contador de mensajes procesados
	var processed atomic.Int32
	handler := func(ctx context.Context, body []byte) error {
		processed.Add(1)
		return nil
	}

	// Iniciar consumo
	err = consumer.ConsumeWithDLQ(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar mensaje de test
	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("test message"),
			DeliveryMode: amqp.Persistent,
		},
	)
	require.NoError(t, err)

	// Esperar procesamiento
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, int32(1), processed.Load())

	// Verificar que DLQ y DLX fueron creados
	channel, err := rabbitContainer.Channel()
	require.NoError(t, err)
	defer channel.Close()

	// Verificar DLQ existe
	_, err = channel.QueueInspect(dlqConfig.DLXRoutingKey)
	assert.NoError(t, err, "DLQ debe haber sido creada")
}

// TestConsumeWithDLQ_MessageSentToDLQ verifica que mensajes fallidos van a DLQ
func TestConsumeWithDLQ_MessageSentToDLQ(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Configurar DLQ con pocas retries
	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            1, // Solo 1 retry
		RetryDelay:            50 * time.Millisecond,
		DLXExchange:           "test.dlx.fail",
		DLXRoutingKey:         "test.dlq.fail",
		UseExponentialBackoff: false,
	}

	// Crear cola
	queueName := "test_queue_dlq_fail"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	// Consumer config
	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_fail",
		AutoAck:       false,
		PrefetchCount: 1,
		DLQ:           dlqConfig,
	}

	consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)

	// Handler que siempre falla
	var attempts atomic.Int32
	handler := func(ctx context.Context, body []byte) error {
		attempts.Add(1)
		return errors.New("processing failed")
	}

	// Iniciar consumo
	err = consumer.ConsumeWithDLQ(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar mensaje
	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("message that will fail"),
			DeliveryMode: amqp.Persistent,
		},
	)
	require.NoError(t, err)

	// Esperar procesamiento y retries
	time.Sleep(2 * time.Second)

	// Verificar que se intentó procesar más de una vez (original + retry)
	assert.GreaterOrEqual(t, attempts.Load(), int32(2), "Debe haber al menos 2 intentos")

	// Verificar que el mensaje está en DLQ
	dlqInfo := waitForQueueMessages(t, rabbitContainer, dlqConfig.DLXRoutingKey, 1, 5*time.Second)
	assert.Equal(t, 1, dlqInfo.Messages, "Debe haber 1 mensaje en DLQ")
}

// TestConsumeWithDLQ_SuccessfulRetry verifica que los reintentos funcionan
func TestConsumeWithDLQ_SuccessfulRetry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            3,
		RetryDelay:            100 * time.Millisecond,
		DLXExchange:           "test.dlx.retry",
		DLXRoutingKey:         "test.dlq.retry",
		UseExponentialBackoff: false,
	}

	queueName := "test_queue_dlq_retry"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_retry",
		AutoAck:       false,
		PrefetchCount: 1,
		DLQ:           dlqConfig,
	}

	consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)

	// Handler que falla las primeras 2 veces, luego tiene éxito
	var attempts atomic.Int32
	handler := func(ctx context.Context, body []byte) error {
		attempt := attempts.Add(1)
		if attempt <= 2 {
			return fmt.Errorf("attempt %d failed", attempt)
		}
		return nil // Éxito en el tercer intento
	}

	err = consumer.ConsumeWithDLQ(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar mensaje
	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("retry message"),
			DeliveryMode: amqp.Persistent,
		},
	)
	require.NoError(t, err)

	// Esperar procesamiento con retries
	time.Sleep(2 * time.Second)

	// Verificar que se procesó exitosamente después de retries
	assert.Equal(t, int32(3), attempts.Load(), "Debe haber 3 intentos (2 fallos + 1 éxito)")
}

// TestConsumeWithDLQ_ExponentialBackoff verifica que el backoff exponencial funciona
func TestConsumeWithDLQ_ExponentialBackoff(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            3,
		RetryDelay:            100 * time.Millisecond,
		DLXExchange:           "test.dlx.backoff",
		DLXRoutingKey:         "test.dlq.backoff",
		UseExponentialBackoff: true, // Activar backoff exponencial
	}

	queueName := "test_queue_backoff"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_backoff",
		AutoAck:       false,
		PrefetchCount: 1,
		DLQ:           dlqConfig,
	}

	consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)

	// Registrar tiempos de cada intento
	var attemptTimes []time.Time
	var mu sync.Mutex

	handler := func(ctx context.Context, body []byte) error {
		mu.Lock()
		attemptTimes = append(attemptTimes, time.Now())
		mu.Unlock()
		return errors.New("always fail")
	}

	err = consumer.ConsumeWithDLQ(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar mensaje
	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("backoff test"),
			DeliveryMode: amqp.Persistent,
		},
	)
	require.NoError(t, err)

	// Esperar suficiente tiempo para retries con backoff
	// Backoff: 100ms, 200ms, 400ms, 800ms
	time.Sleep(5 * time.Second)

	mu.Lock()
	defer mu.Unlock()

	// Debe haber múltiples intentos
	assert.GreaterOrEqual(t, len(attemptTimes), 2, "Debe haber al menos 2 intentos")

	// Verificar que los intervalos aumentan (backoff exponencial)
	if len(attemptTimes) >= 3 {
		interval1 := attemptTimes[1].Sub(attemptTimes[0])
		interval2 := attemptTimes[2].Sub(attemptTimes[1])
		assert.Less(t, interval1, interval2, "El segundo intervalo debe ser mayor (backoff exponencial)")
	}
}

// TestConsumeWithDLQ_HeadersPreserved verifica que los headers se preservan
func TestConsumeWithDLQ_HeadersPreserved(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            1,
		RetryDelay:            50 * time.Millisecond,
		DLXExchange:           "test.dlx.headers",
		DLXRoutingKey:         "test.dlq.headers",
		UseExponentialBackoff: false,
	}

	queueName := "test_queue_headers"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_headers",
		AutoAck:       false,
		PrefetchCount: 1,
		DLQ:           dlqConfig,
	}

	consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)

	// Handler que siempre falla
	handler := func(ctx context.Context, body []byte) error {
		return errors.New("fail to DLQ")
	}

	err = consumer.ConsumeWithDLQ(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar mensaje con headers custom
	customHeaders := amqp.Table{
		"custom-header-1": "value1",
		"custom-header-2": 42,
	}

	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("message with headers"),
			Headers:      customHeaders,
			DeliveryMode: amqp.Persistent,
		},
	)
	require.NoError(t, err)

	// Esperar a que llegue a DLQ
	time.Sleep(2 * time.Second)

	// Consumir de DLQ y verificar headers
	channel, err := rabbitContainer.Channel()
	require.NoError(t, err)
	defer channel.Close()

	// Consumir un mensaje de DLQ
	msg, ok, err := channel.Get(dlqConfig.DLXRoutingKey, false)
	require.NoError(t, err)
	require.True(t, ok, "Debe haber mensaje en DLQ")

	// Verificar headers DLQ
	assert.Contains(t, msg.Headers, "x-original-exchange")
	assert.Contains(t, msg.Headers, "x-original-routing-key")
	assert.Contains(t, msg.Headers, "x-failed-at")
	assert.Contains(t, msg.Headers, "x-retry-count")

	// Limpiar
	_ = msg.Ack(false)
}

// TestConsumeWithDLQ_MultipleConsumers verifica múltiples consumidores DLQ
func TestConsumeWithDLQ_MultipleConsumers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            3,
		RetryDelay:            100 * time.Millisecond,
		DLXExchange:           "test.dlx.multi",
		DLXRoutingKey:         "test.dlq.multi",
		UseExponentialBackoff: false,
	}

	queueName := "test_queue_multi_consumer"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	// Crear múltiples consumers
	var processed atomic.Int32
	handler := func(ctx context.Context, body []byte) error {
		processed.Add(1)
		return nil
	}

	consumers := 3
	for i := 0; i < consumers; i++ {
		consumerConfig := ConsumerConfig{
			Name:          fmt.Sprintf("test_consumer_%d", i),
			AutoAck:       false,
			PrefetchCount: 2,
			DLQ:           dlqConfig,
		}

		consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)
		err = consumer.ConsumeWithDLQ(ctx, queueName, handler)
		require.NoError(t, err)
	}

	// Publicar múltiples mensajes
	messageCount := 10
	for i := 0; i < messageCount; i++ {
		err = ch.PublishWithContext(
			ctx,
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(fmt.Sprintf("message %d", i)),
				DeliveryMode: amqp.Persistent,
			},
		)
		require.NoError(t, err)
	}

	// Esperar procesamiento
	time.Sleep(2 * time.Second)

	// Todos los mensajes deben ser procesados
	assert.Equal(t, int32(messageCount), processed.Load())
}

// TestConsumeWithDLQ_ContextCancellation verifica cancelación por contexto
func TestConsumeWithDLQ_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            5,
		RetryDelay:            200 * time.Millisecond,
		DLXExchange:           "test.dlx.cancel",
		DLXRoutingKey:         "test.dlq.cancel",
		UseExponentialBackoff: false,
	}

	queueName := "test_queue_cancel"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_cancel",
		AutoAck:       false,
		PrefetchCount: 1,
		DLQ:           dlqConfig,
	}

	consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)

	// Contexto con cancelación
	ctx, cancel := context.WithCancel(context.Background())

	var attempts atomic.Int32
	handler := func(ctx context.Context, body []byte) error {
		attempts.Add(1)
		return errors.New("fail")
	}

	err = consumer.ConsumeWithDLQ(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar mensaje
	err = ch.PublishWithContext(
		context.Background(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("cancel test"),
			DeliveryMode: amqp.Persistent,
		},
	)
	require.NoError(t, err)

	// Esperar un poco
	time.Sleep(300 * time.Millisecond)

	// Cancelar contexto
	cancel()

	// Esperar un poco más
	time.Sleep(500 * time.Millisecond)

	// No debe haber más intentos después de cancelar
	firstAttempts := attempts.Load()
	time.Sleep(500 * time.Millisecond)
	finalAttempts := attempts.Load()

	assert.Equal(t, firstAttempts, finalAttempts, "No debe haber más intentos después de cancelar")
}

// TestConsumeWithDLQ_CleanShutdown verifica cierre limpio del consumer
func TestConsumeWithDLQ_CleanShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            3,
		RetryDelay:            100 * time.Millisecond,
		DLXExchange:           "test.dlx.shutdown",
		DLXRoutingKey:         "test.dlq.shutdown",
		UseExponentialBackoff: false,
	}

	queueName := "test_queue_shutdown"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_shutdown",
		AutoAck:       false,
		PrefetchCount: 5,
		DLQ:           dlqConfig,
	}

	consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)

	handler := func(ctx context.Context, body []byte) error {
		time.Sleep(100 * time.Millisecond) // Simular procesamiento
		return nil
	}

	err = consumer.ConsumeWithDLQ(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar algunos mensajes
	for i := 0; i < 5; i++ {
		err = ch.PublishWithContext(
			ctx,
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(fmt.Sprintf("message %d", i)),
				DeliveryMode: amqp.Persistent,
			},
		)
		require.NoError(t, err)
	}

	// Esperar un poco
	time.Sleep(200 * time.Millisecond)

	// Cerrar consumer
	err = consumer.Close()
	assert.NoError(t, err, "Close debe ser exitoso")
}

// TestConsumeWithDLQ_QueueNotExists verifica comportamiento cuando cola no existe
func TestConsumeWithDLQ_QueueNotExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            3,
		RetryDelay:            100 * time.Millisecond,
		DLXExchange:           "test.dlx.notexist",
		DLXRoutingKey:         "test.dlq.notexist",
		UseExponentialBackoff: false,
	}

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_notexist",
		AutoAck:       false,
		PrefetchCount: 1,
		DLQ:           dlqConfig,
	}

	consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)

	handler := func(ctx context.Context, body []byte) error {
		return nil
	}

	// Intentar consumir de cola que no existe
	err = consumer.ConsumeWithDLQ(ctx, "non_existent_queue", handler)
	assert.Error(t, err, "Debe dar error al consumir cola inexistente")
}

// TestConsumeWithDLQ_AutoAckDisabled verifica que AutoAck=false funciona con DLQ
func TestConsumeWithDLQ_AutoAckDisabled(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            2,
		RetryDelay:            100 * time.Millisecond,
		DLXExchange:           "test.dlx.autoack",
		DLXRoutingKey:         "test.dlq.autoack",
		UseExponentialBackoff: false,
	}

	queueName := "test_queue_autoack"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	consumerConfig := ConsumerConfig{
		Name:          "test_consumer_autoack",
		AutoAck:       false, // Manualmente ACK/NACK
		PrefetchCount: 1,
		DLQ:           dlqConfig,
	}

	consumer := NewConsumer(conn, consumerConfig).(*RabbitMQConsumer)

	var processed atomic.Int32
	handler := func(ctx context.Context, body []byte) error {
		processed.Add(1)
		return nil
	}

	err = consumer.ConsumeWithDLQ(ctx, queueName, handler)
	require.NoError(t, err)

	// Publicar mensaje
	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("autoack test"),
			DeliveryMode: amqp.Persistent,
		},
	)
	require.NoError(t, err)

	// Esperar procesamiento
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, int32(1), processed.Load(), "Mensaje debe procesarse una vez")
}
