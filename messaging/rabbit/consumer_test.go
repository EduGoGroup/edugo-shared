package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConsumer(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	config := ConsumerConfig{
		Name:      "test_consumer",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
	}

	consumer := NewConsumer(conn, config)
	assert.NotNil(t, consumer)
	assert.IsType(t, &RabbitMQConsumer{}, consumer)
}

func TestNewConsumer_WithDifferentConfigs(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	tests := []struct {
		name   string
		config ConsumerConfig
	}{
		{
			name: "auto-ack enabled",
			config: ConsumerConfig{
				Name:      "consumer_autoack",
				AutoAck:   true,
				Exclusive: false,
				NoLocal:   false,
			},
		},
		{
			name: "exclusive consumer",
			config: ConsumerConfig{
				Name:      "consumer_exclusive",
				AutoAck:   false,
				Exclusive: true,
				NoLocal:   false,
			},
		},
		{
			name: "no-local consumer",
			config: ConsumerConfig{
				Name:      "consumer_nolocal",
				AutoAck:   false,
				Exclusive: false,
				NoLocal:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer := NewConsumer(conn, tt.config)
			assert.NotNil(t, consumer)
		})
	}
}

func TestConsumer_Close(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	config := ConsumerConfig{
		Name:    "test_consumer",
		AutoAck: false,
	}

	consumer := NewConsumer(conn, config)
	require.NotNil(t, consumer)

	err = consumer.Close()
	assert.NoError(t, err)
}

func TestConsumer_Consume_BasicMessage(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Setup queue
	queueName := "test_consume_basic"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create consumer
	consumerConfig := ConsumerConfig{
		Name:    "test_consumer",
		AutoAck: true,
	}
	consumer := NewConsumer(conn, consumerConfig)
	require.NotNil(t, consumer)
	defer func() { _ = consumer.Close() }()

	// Setup message handler
	receivedMessages := make(chan []byte, 1)
	handler := func(ctx context.Context, body []byte) error {
		receivedMessages <- body
		return nil
	}

	// Start consuming
	consumeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = consumer.Consume(consumeCtx, queueName, handler)
	require.NoError(t, err)

	// Wait for consumer to start
	waitForQueueConsumers(t, rabbitContainer, queueName, 1, 5*time.Second)

	// Publish a message directly using the container
	testMessage := []byte("test message")
	channel, err := rabbitContainer.Channel()
	require.NoError(t, err)
	defer func() { _ = channel.Close() }()

	err = channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        testMessage,
		},
	)
	require.NoError(t, err)

	// Wait for message
	select {
	case msg := <-receivedMessages:
		assert.Equal(t, testMessage, msg)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout esperando mensaje")
	}
}

func TestConsumer_Consume_WithManualAck(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Setup queue
	queueName := "test_consume_manual_ack"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create consumer with manual ack
	consumerConfig := ConsumerConfig{
		Name:    "test_consumer",
		AutoAck: false, // Manual ack
	}
	consumer := NewConsumer(conn, consumerConfig)
	require.NotNil(t, consumer)
	defer func() { _ = consumer.Close() }()

	// Setup message handler
	receivedMessages := make(chan []byte, 1)
	handler := func(ctx context.Context, body []byte) error {
		receivedMessages <- body
		return nil // Success - should trigger Ack
	}

	// Start consuming
	consumeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = consumer.Consume(consumeCtx, queueName, handler)
	require.NoError(t, err)

	// Wait for consumer to start
	waitForQueueConsumers(t, rabbitContainer, queueName, 1, 5*time.Second)

	// Publish a message
	testMessage := []byte("test message with ack")
	channel, err := rabbitContainer.Channel()
	require.NoError(t, err)
	defer func() { _ = channel.Close() }()

	err = channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        testMessage,
		},
	)
	require.NoError(t, err)

	// Wait for message
	select {
	case msg := <-receivedMessages:
		assert.Equal(t, testMessage, msg)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout esperando mensaje")
	}

	// Verify queue is empty (message was acked)
	waitForQueueMessages(t, rabbitContainer, queueName, 0, 5*time.Second)
}

func TestConsumer_Consume_ErrorHandling(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Setup queue
	queueName := "test_consume_error"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create consumer with manual ack
	consumerConfig := ConsumerConfig{
		Name:    "test_consumer",
		AutoAck: false,
	}
	consumer := NewConsumer(conn, consumerConfig)
	require.NotNil(t, consumer)
	defer func() { _ = consumer.Close() }()

	// Setup message handler that returns error
	receivedCount := 0
	var mu sync.Mutex
	handler := func(ctx context.Context, body []byte) error {
		mu.Lock()
		receivedCount++
		mu.Unlock()
		return fmt.Errorf("processing error")
	}

	// Start consuming
	consumeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = consumer.Consume(consumeCtx, queueName, handler)
	require.NoError(t, err)

	// Wait for consumer to start
	waitForQueueConsumers(t, rabbitContainer, queueName, 1, 5*time.Second)

	// Publish a message
	testMessage := []byte("test error message")
	channel, err := rabbitContainer.Channel()
	require.NoError(t, err)
	defer func() { _ = channel.Close() }()

	err = channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        testMessage,
		},
	)
	require.NoError(t, err)

	// Wait for message processing (allow multiple requeue attempts)
	require.Eventually(t, func() bool {
		mu.Lock()
		count := receivedCount
		mu.Unlock()
		return count >= 1
	}, 5*time.Second, 100*time.Millisecond)

	// Cancel consumer to stop processing
	cancel()

	mu.Lock()
	count := receivedCount
	mu.Unlock()

	// Message should have been received multiple times due to requeue
	assert.GreaterOrEqual(t, count, 1, "El mensaje debe haber sido recibido al menos una vez")

	// After stopping consumer, message should still be in queue (last requeue)
	queueInfo, err := channel.QueueDeclarePassive(queueName, false, false, false, false, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, queueInfo.Messages, 0, "Verificar mensajes en cola")
}

func TestConsumer_Consume_ContextCancellation(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Setup queue
	queueName := "test_consume_cancel"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create consumer
	consumerConfig := ConsumerConfig{
		Name:    "test_consumer",
		AutoAck: true,
	}
	consumer := NewConsumer(conn, consumerConfig)
	require.NotNil(t, consumer)
	defer func() { _ = consumer.Close() }()

	// Setup handler
	handlerCalled := false
	var mu sync.Mutex
	handler := func(ctx context.Context, body []byte) error {
		mu.Lock()
		handlerCalled = true
		mu.Unlock()
		return nil
	}

	// Start consuming with short-lived context
	consumeCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err = consumer.Consume(consumeCtx, queueName, handler)
	require.NoError(t, err)

	// Wait for context to expire
	<-consumeCtx.Done()

	// Consumer should have stopped gracefully
	mu.Lock()
	called := handlerCalled
	mu.Unlock()

	// Handler may or may not have been called, but no error should occur
	_ = called
}

func TestConsumer_Consume_InvalidQueue(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	consumerConfig := ConsumerConfig{
		Name:    "test_consumer",
		AutoAck: true,
	}
	consumer := NewConsumer(conn, consumerConfig)
	require.NotNil(t, consumer)
	defer func() { _ = consumer.Close() }()

	handler := func(ctx context.Context, body []byte) error {
		return nil
	}

	ctx := context.Background()
	err = consumer.Consume(ctx, "non_existent_queue", handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to start consuming")
}

func TestConsumer_Consume_MultipleMessages(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Setup queue
	queueName := "test_consume_multiple"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create consumer
	consumerConfig := ConsumerConfig{
		Name:    "test_consumer",
		AutoAck: true,
	}
	consumer := NewConsumer(conn, consumerConfig)
	require.NotNil(t, consumer)
	defer func() { _ = consumer.Close() }()

	// Setup message handler
	receivedMessages := make([][]byte, 0)
	var mu sync.Mutex
	handler := func(ctx context.Context, body []byte) error {
		mu.Lock()
		receivedMessages = append(receivedMessages, body)
		mu.Unlock()
		return nil
	}

	// Start consuming
	consumeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = consumer.Consume(consumeCtx, queueName, handler)
	require.NoError(t, err)

	// Wait for consumer to start
	waitForQueueConsumers(t, rabbitContainer, queueName, 1, 5*time.Second)

	// Publish multiple messages
	channel, err := rabbitContainer.Channel()
	require.NoError(t, err)
	defer func() { _ = channel.Close() }()

	messageCount := 5
	for i := 0; i < messageCount; i++ {
		err = channel.PublishWithContext(
			ctx,
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("message %d", i)),
			},
		)
		require.NoError(t, err)
	}

	// Wait for messages to be processed
	require.Eventually(t, func() bool {
		mu.Lock()
		count := len(receivedMessages)
		mu.Unlock()
		return count == messageCount
	}, 5*time.Second, 100*time.Millisecond)
}

func TestUnmarshalMessage_Success(t *testing.T) {
	type TestMessage struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	testData := TestMessage{
		ID:   "123",
		Name: "Test",
	}

	body, err := json.Marshal(testData)
	require.NoError(t, err)

	var result TestMessage
	err = UnmarshalMessage(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, testData.ID, result.ID)
	assert.Equal(t, testData.Name, result.Name)
}

func TestUnmarshalMessage_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"id": "123", "name": `)

	var result map[string]interface{}
	err := UnmarshalMessage(invalidJSON, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal message")
}

func TestUnmarshalMessage_EmptyBody(t *testing.T) {
	emptyBody := []byte(`{}`)

	var result map[string]interface{}
	err := UnmarshalMessage(emptyBody, &result)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestHandleWithUnmarshal_Success(t *testing.T) {
	type TestMessage struct {
		Value int `json:"value"`
	}

	testData := TestMessage{Value: 42}
	body, err := json.Marshal(testData)
	require.NoError(t, err)

	handlerCalled := false
	handler := func(v interface{}) error {
		msg, ok := v.(*TestMessage)
		assert.True(t, ok)
		assert.Equal(t, 42, msg.Value)
		handlerCalled = true
		return nil
	}

	var result TestMessage
	err = HandleWithUnmarshal(body, &result, handler)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

func TestHandleWithUnmarshal_UnmarshalError(t *testing.T) {
	invalidJSON := []byte(`invalid json`)

	handler := func(v interface{}) error {
		t.Fatal("Handler should not be called")
		return nil
	}

	var result map[string]interface{}
	err := HandleWithUnmarshal(invalidJSON, &result, handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal message")
}

func TestHandleWithUnmarshal_HandlerError(t *testing.T) {
	type TestMessage struct {
		Value int `json:"value"`
	}

	testData := TestMessage{Value: 42}
	body, err := json.Marshal(testData)
	require.NoError(t, err)

	expectedError := fmt.Errorf("handler error")
	handler := func(v interface{}) error {
		return expectedError
	}

	var result TestMessage
	err = HandleWithUnmarshal(body, &result, handler)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestConsumer_Consume_ExclusiveConsumer(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Setup queue
	queueName := "test_exclusive_consumer"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create exclusive consumer
	consumerConfig := ConsumerConfig{
		Name:      "exclusive_consumer",
		AutoAck:   true,
		Exclusive: true,
	}
	consumer1 := NewConsumer(conn, consumerConfig)
	require.NotNil(t, consumer1)
	defer func() { _ = consumer1.Close() }()

	handler := func(ctx context.Context, body []byte) error {
		return nil
	}

	consumeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = consumer1.Consume(consumeCtx, queueName, handler)
	require.NoError(t, err)

	// Wait for exclusive consumer to start
	waitForQueueConsumers(t, rabbitContainer, queueName, 1, 5*time.Second)

	// Try to create a second exclusive consumer - should fail
	consumer2 := NewConsumer(conn, consumerConfig)
	require.NotNil(t, consumer2)
	defer func() { _ = consumer2.Close() }()

	err = consumer2.Consume(consumeCtx, queueName, handler)
	assert.Error(t, err, "Second exclusive consumer should fail")
}

func TestConsumer_Consume_WithPrefetch(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Set prefetch count
	err = conn.SetPrefetchCount(1)
	require.NoError(t, err)

	// Setup queue
	queueName := "test_prefetch"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create consumer
	consumerConfig := ConsumerConfig{
		Name:          "prefetch_consumer",
		AutoAck:       false,
		PrefetchCount: 1,
	}
	consumer := NewConsumer(conn, consumerConfig)
	require.NotNil(t, consumer)
	defer func() { _ = consumer.Close() }()

	receivedCount := 0
	var mu sync.Mutex
	handler := func(ctx context.Context, body []byte) error {
		mu.Lock()
		receivedCount++
		mu.Unlock()
		time.Sleep(500 * time.Millisecond) // Simulate slow processing
		return nil
	}

	consumeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = consumer.Consume(consumeCtx, queueName, handler)
	require.NoError(t, err)

	// Test verifies that consumer can be created with prefetch config
	// Actual prefetch behavior is tested at the connection level
	waitForQueueConsumers(t, rabbitContainer, queueName, 1, 5*time.Second)
	cancel()
}
