package rabbit

import (
	"context"
	"errors"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetRetryCount_ZeroWhenNil verifica que retorna 0 cuando headers es nil
func TestGetRetryCount_ZeroWhenNil(t *testing.T) {
	count := getRetryCount(nil)
	assert.Equal(t, 0, count)
}

// TestGetRetryCount_ZeroWhenMissing verifica que retorna 0 cuando no existe el header
func TestGetRetryCount_ZeroWhenMissing(t *testing.T) {
	headers := amqp.Table{
		"other-header": "value",
	}
	count := getRetryCount(headers)
	assert.Equal(t, 0, count)
}

// TestGetRetryCount_Int verifica lectura correcta con tipo int
func TestGetRetryCount_Int(t *testing.T) {
	headers := amqp.Table{
		"x-retry-count": 5,
	}
	count := getRetryCount(headers)
	assert.Equal(t, 5, count)
}

// TestGetRetryCount_Int32 verifica lectura correcta con tipo int32
func TestGetRetryCount_Int32(t *testing.T) {
	headers := amqp.Table{
		"x-retry-count": int32(10),
	}
	count := getRetryCount(headers)
	assert.Equal(t, 10, count)
}

// TestGetRetryCount_Int64 verifica lectura correcta con tipo int64
func TestGetRetryCount_Int64(t *testing.T) {
	headers := amqp.Table{
		"x-retry-count": int64(15),
	}
	count := getRetryCount(headers)
	assert.Equal(t, 15, count)
}

// TestCloneHeaders_NilInput verifica que clone de nil retorna tabla vacía
func TestCloneHeaders_NilInput(t *testing.T) {
	cloned := cloneHeaders(nil)
	assert.NotNil(t, cloned)
	assert.Equal(t, 0, len(cloned))
}

// TestCloneHeaders_EmptyInput verifica clone de tabla vacía
func TestCloneHeaders_EmptyInput(t *testing.T) {
	original := amqp.Table{}
	cloned := cloneHeaders(original)

	assert.NotNil(t, cloned)
	assert.Equal(t, 0, len(cloned))
}

// TestCloneHeaders_WithData verifica que clona correctamente los datos
func TestCloneHeaders_WithData(t *testing.T) {
	original := amqp.Table{
		"header1": "value1",
		"header2": 42,
		"header3": true,
		"header4": []byte("bytes"),
	}

	cloned := cloneHeaders(original)

	// Verificar que tiene los mismos valores
	assert.Equal(t, len(original), len(cloned))
	assert.Equal(t, original["header1"], cloned["header1"])
	assert.Equal(t, original["header2"], cloned["header2"])
	assert.Equal(t, original["header3"], cloned["header3"])
	assert.Equal(t, original["header4"], cloned["header4"])
}

// TestCloneHeaders_IndependentCopy verifica que es una copia independiente
func TestCloneHeaders_IndependentCopy(t *testing.T) {
	original := amqp.Table{
		"header1": "value1",
	}

	cloned := cloneHeaders(original)

	// Modificar el clone no debe afectar el original
	cloned["header2"] = "value2"

	assert.Contains(t, cloned, "header2")
	assert.NotContains(t, original, "header2")
}

// TestDLQConfig_CalculateBackoff_Linear verifica backoff lineal
func TestDLQConfig_CalculateBackoff_Linear(t *testing.T) {
	config := DLQConfig{
		RetryDelay:            5 * time.Second,
		UseExponentialBackoff: false,
	}

	// Con backoff lineal, siempre debe retornar el mismo delay
	for attempt := 0; attempt < 5; attempt++ {
		backoff := config.CalculateBackoff(attempt)
		assert.Equal(t, 5*time.Second, backoff, "Backoff lineal debe ser constante")
	}
}

// TestDLQConfig_CalculateBackoff_Exponential verifica backoff exponencial
func TestDLQConfig_CalculateBackoff_Exponential(t *testing.T) {
	config := DLQConfig{
		RetryDelay:            5 * time.Second,
		UseExponentialBackoff: true,
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 5 * time.Second},   // 5 * 2^0 = 5
		{1, 10 * time.Second},  // 5 * 2^1 = 10
		{2, 20 * time.Second},  // 5 * 2^2 = 20
		{3, 40 * time.Second},  // 5 * 2^3 = 40
		{4, 80 * time.Second},  // 5 * 2^4 = 80
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			backoff := config.CalculateBackoff(tt.attempt)
			assert.Equal(t, tt.expected, backoff)
		})
	}
}

// TestDefaultDLQConfig verifica la configuración por defecto
func TestDefaultDLQConfig(t *testing.T) {
	config := DefaultDLQConfig()

	assert.True(t, config.Enabled)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 5*time.Second, config.RetryDelay)
	assert.Equal(t, "dlx", config.DLXExchange)
	assert.Equal(t, "dlq", config.DLXRoutingKey)
	assert.True(t, config.UseExponentialBackoff)
}

// TestSendToDLQ_Integration verifica envío a DLQ
func TestSendToDLQ_Integration(t *testing.T) {
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
		MaxRetries:            1,
		RetryDelay:            50 * time.Millisecond,
		DLXExchange:           "test.dlx.send",
		DLXRoutingKey:         "test.dlq.send",
		UseExponentialBackoff: false,
	}

	queueName := "test_queue_send_dlq"
	ch := conn.GetChannel()

	// Crear cola principal
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	consumer := &RabbitMQConsumer{
		conn: conn,
		config: ConsumerConfig{
			Name:          "test_consumer",
			AutoAck:       false,
			PrefetchCount: 1,
			DLQ:           dlqConfig,
		},
	}

	// Setup DLQ infrastructure
	err = consumer.setupDLQ(ch)
	require.NoError(t, err)

	// Crear mensaje de test
	testMessage := amqp.Delivery{
		Body:        []byte("test message for DLQ"),
		ContentType: "text/plain",
		Exchange:    "",
		RoutingKey:  queueName,
		Headers: amqp.Table{
			"custom-header": "custom-value",
		},
	}

	// Enviar a DLQ
	err = consumer.sendToDLQ(ch, testMessage)
	require.NoError(t, err)

	// Verificar que llegó a DLQ
	time.Sleep(200 * time.Millisecond)

	channel, err := rabbitContainer.Channel()
	require.NoError(t, err)
	defer channel.Close()

	msg, ok, err := channel.Get(dlqConfig.DLXRoutingKey, false)
	require.NoError(t, err)
	require.True(t, ok, "Debe haber mensaje en DLQ")

	assert.Equal(t, testMessage.Body, msg.Body)
	assert.Contains(t, msg.Headers, "x-original-exchange")
	assert.Contains(t, msg.Headers, "x-original-routing-key")
	assert.Contains(t, msg.Headers, "x-failed-at")
	assert.Contains(t, msg.Headers, "custom-header")

	_ = msg.Ack(false)
}

// TestSetupDLQ_Integration verifica setup de DLQ infrastructure
func TestSetupDLQ_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	rabbitContainer, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            3,
		RetryDelay:            5 * time.Second,
		DLXExchange:           "test.dlx.setup",
		DLXRoutingKey:         "test.dlq.setup",
		UseExponentialBackoff: true,
	}

	consumer := &RabbitMQConsumer{
		conn: conn,
		config: ConsumerConfig{
			Name:          "test_consumer_setup",
			AutoAck:       false,
			PrefetchCount: 1,
			DLQ:           dlqConfig,
		},
	}

	ch := conn.GetChannel()

	// Setup DLQ
	err = consumer.setupDLQ(ch)
	require.NoError(t, err)

	// Verificar que el DLX (exchange) existe
	channel, err := rabbitContainer.Channel()
	require.NoError(t, err)
	defer channel.Close()

	// Intentar declarar el exchange de nuevo (debe funcionar si existe)
	err = channel.ExchangeDeclarePassive(
		dlqConfig.DLXExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	assert.NoError(t, err, "DLX exchange debe existir")

	// Verificar que DLQ (queue) existe
	_, err = channel.QueueInspect(dlqConfig.DLXRoutingKey)
	assert.NoError(t, err, "DLQ queue debe existir")
}

// TestProcessMessage_Success verifica procesamiento exitoso
func TestProcessMessage_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_process_success"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	consumer := &RabbitMQConsumer{
		conn: conn,
		config: ConsumerConfig{
			Name:          "test_consumer_process",
			AutoAck:       false,
			PrefetchCount: 1,
			DLQ: DLQConfig{
				Enabled:    false, // Sin DLQ para este test
				MaxRetries: 0,
			},
		},
	}

	// Publicar mensaje
	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			Body:         []byte("test"),
			DeliveryMode: amqp.Persistent,
		},
	)
	require.NoError(t, err)

	// Consumir mensaje
	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		// Handler exitoso
		handlerCalled := false
		handler := func(ctx context.Context, body []byte) error {
			handlerCalled = true
			assert.Equal(t, []byte("test"), body)
			return nil
		}

		// Procesar mensaje
		consumer.processMessage(ctx, ch, queueName, handler, msg)

		assert.True(t, handlerCalled, "Handler debe ser llamado")

	case <-time.After(2 * time.Second):
		t.Fatal("Timeout esperando mensaje")
	}
}

// TestProcessMessage_Failure verifica procesamiento con error
func TestProcessMessage_Failure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	queueName := "test_queue_process_fail"
	ch := conn.GetChannel()
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	require.NoError(t, err)

	consumer := &RabbitMQConsumer{
		conn: conn,
		config: ConsumerConfig{
			Name:          "test_consumer_process_fail",
			AutoAck:       false,
			PrefetchCount: 1,
			DLQ: DLQConfig{
				Enabled:    false,
				MaxRetries: 0,
			},
		},
	}

	// Publicar mensaje
	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			Body:         []byte("fail test"),
			DeliveryMode: amqp.Persistent,
		},
	)
	require.NoError(t, err)

	// Consumir mensaje
	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		// Handler que falla
		handler := func(ctx context.Context, body []byte) error {
			return errors.New("processing error")
		}

		// Procesar mensaje (debe hacer NACK internamente)
		consumer.processMessage(ctx, ch, queueName, handler, msg)

		// El mensaje debe volver a la cola (NACK with requeue)
		// Nota: Esto es difícil de verificar sin consumir de nuevo

	case <-time.After(2 * time.Second):
		t.Fatal("Timeout esperando mensaje")
	}
}
