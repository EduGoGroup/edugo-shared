package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPublisher(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	publisher := NewPublisher(conn)
	assert.NotNil(t, publisher)
	assert.IsType(t, &RabbitMQPublisher{}, publisher)
}

func TestPublisher_Publish_BasicMessage(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup exchange and queue
	exchangeConfig := ExchangeConfig{
		Name:       "test_publish_exchange",
		Type:       "direct",
		Durable:    false,
		AutoDelete: true,
	}
	err = conn.DeclareExchange(exchangeConfig)
	require.NoError(t, err)

	queueConfig := QueueConfig{
		Name:       "test_publish_queue",
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	queue, err := conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	err = conn.BindQueue(queue.Name, "test_key", exchangeConfig.Name)
	require.NoError(t, err)

	// Create publisher
	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Publish message
	type TestMessage struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}

	testMsg := TestMessage{
		ID:   "123",
		Text: "Hello, RabbitMQ!",
	}

	err = publisher.Publish(ctx, exchangeConfig.Name, "test_key", testMsg)
	assert.NoError(t, err)

	// Verify message was published
	waitForQueueMessages(t, rabbitContainer, queue.Name, 1, 5*time.Second)
}

func TestPublisher_Publish_ToDefaultExchange(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue (no need for explicit exchange with default exchange)
	queueName := "test_default_exchange"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create publisher
	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Publish to default exchange (empty string) with queue name as routing key
	testMsg := map[string]interface{}{
		"message": "test message",
	}

	err = publisher.Publish(ctx, "", queueName, testMsg)
	assert.NoError(t, err)

	// Verify message
	waitForQueueMessages(t, rabbitContainer, queueName, 1, 5*time.Second)
}

func TestPublisher_PublishWithPriority(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue with max priority
	queueName := "test_priority_queue"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args: map[string]interface{}{
			"x-max-priority": 10,
		},
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create publisher
	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Publish messages with different priorities
	priorities := []uint8{0, 5, 10}
	for _, priority := range priorities {
		testMsg := map[string]interface{}{
			"priority": priority,
			"message":  fmt.Sprintf("Message with priority %d", priority),
		}

		err = publisher.PublishWithPriority(ctx, "", queueName, testMsg, priority)
		assert.NoError(t, err, fmt.Sprintf("Error publishing with priority %d", priority))
	}

	// Verify messages
	waitForQueueMessages(t, rabbitContainer, queueName, len(priorities), 5*time.Second)
}

func TestPublisher_Publish_InvalidJSON(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Try to publish something that can't be marshaled to JSON
	invalidMsg := make(chan int) // channels can't be marshaled to JSON

	err = publisher.Publish(ctx, "", "test_queue", invalidMsg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal message")
}

func TestPublisher_Publish_ToNonExistentExchange(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	testMsg := map[string]string{
		"test": "message",
	}

	// RabbitMQ NO retorna error cuando publicas a exchange inexistente con mandatory=false
	// El mensaje simplemente se descarta silenciosamente (comportamiento estándar de RabbitMQ)
	err = publisher.Publish(ctx, "non_existent_exchange", "test_key", testMsg)
	assert.NoError(t, err, "RabbitMQ permite publicar a exchange inexistente sin error cuando mandatory=false")
}

func TestPublisher_Close(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)

	err = publisher.Close()
	assert.NoError(t, err)
}

func TestPublisher_Publish_MultipleMessages(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue
	queueName := "test_multiple_publish"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create publisher
	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Publish multiple messages
	messageCount := 10
	for i := 0; i < messageCount; i++ {
		testMsg := map[string]interface{}{
			"index":   i,
			"message": fmt.Sprintf("Message %d", i),
		}

		err = publisher.Publish(ctx, "", queueName, testMsg)
		require.NoError(t, err, fmt.Sprintf("Error publishing message %d", i))
	}

	// Verify all messages were published
	waitForQueueMessages(t, rabbitContainer, queueName, messageCount, 5*time.Second)
}

func TestPublisher_Publish_ConcurrentPublishing(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue
	queueName := "test_concurrent_publish"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create publisher
	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Publish messages concurrently
	concurrency := 5
	messagesPerGoroutine := 10
	var wg sync.WaitGroup
	errors := make(chan error, concurrency*messagesPerGoroutine)

	for g := 0; g < concurrency; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for i := 0; i < messagesPerGoroutine; i++ {
				testMsg := map[string]interface{}{
					"goroutine": goroutineID,
					"index":     i,
				}
				if err := publisher.Publish(ctx, "", queueName, testMsg); err != nil {
					errors <- err
				}
			}
		}(g)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Logf("Error: %v", err)
		errorCount++
	}
	assert.Equal(t, 0, errorCount, "No debe haber errores en publicación concurrente")

	// Verify message count
	expectedMessages := concurrency * messagesPerGoroutine
	waitForQueueMessages(t, rabbitContainer, queueName, expectedMessages, 5*time.Second)
}

func TestPublisher_Publish_WithContextTimeout(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue
	queueName := "test_context_timeout"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testMsg := map[string]string{
		"message": "test with timeout",
	}

	err = publisher.Publish(ctx, "", queueName, testMsg)
	assert.NoError(t, err)
}

func TestPublisher_Publish_WithCancelledContext(t *testing.T) {
	_, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	testMsg := map[string]string{
		"message": "test with cancelled context",
	}

	err = publisher.Publish(ctx, "", "test_queue", testMsg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish message")
}

func TestPublisher_Publish_ComplexMessage(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue
	queueName := "test_complex_message"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Complex nested structure
	type Address struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Country string `json:"country"`
	}

	type Person struct {
		Name      string            `json:"name"`
		Age       int               `json:"age"`
		Address   Address           `json:"address"`
		Tags      []string          `json:"tags"`
		Metadata  map[string]string `json:"metadata"`
		Timestamp time.Time         `json:"timestamp"`
	}

	complexMsg := Person{
		Name: "John Doe",
		Age:  30,
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
		Tags: []string{"developer", "go", "rabbitmq"},
		Metadata: map[string]string{
			"department": "engineering",
			"team":       "backend",
		},
		Timestamp: time.Now(),
	}

	err = publisher.Publish(ctx, "", queueName, complexMsg)
	assert.NoError(t, err)

	// Verify message
	waitForQueueMessages(t, rabbitContainer, queueName, 1, 5*time.Second)
	channel, err := rabbitContainer.Channel()
	require.NoError(t, err)
	defer channel.Close()

	// Consume and verify the message content
	msgs, err := channel.Consume(
		queueName,
		"test_consumer",
		true,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		var receivedPerson Person
		err = json.Unmarshal(msg.Body, &receivedPerson)
		assert.NoError(t, err)
		assert.Equal(t, complexMsg.Name, receivedPerson.Name)
		assert.Equal(t, complexMsg.Age, receivedPerson.Age)
		assert.Equal(t, complexMsg.Address.City, receivedPerson.Address.City)
		assert.Equal(t, len(complexMsg.Tags), len(receivedPerson.Tags))
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout esperando mensaje")
	}
}

func TestPublisher_PublishWithPriority_ZeroPriority(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue
	queueName := "test_zero_priority"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	testMsg := map[string]string{
		"message": "zero priority message",
	}

	err = publisher.PublishWithPriority(ctx, "", queueName, testMsg, 0)
	assert.NoError(t, err)

	// Verify message
	waitForQueueMessages(t, rabbitContainer, queueName, 1, 5*time.Second)
}

func TestPublisher_PublishWithPriority_MaxPriority(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue with max priority
	queueName := "test_max_priority"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args: map[string]interface{}{
			"x-max-priority": 255,
		},
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	testMsg := map[string]string{
		"message": "max priority message",
	}

	err = publisher.PublishWithPriority(ctx, "", queueName, testMsg, 255)
	assert.NoError(t, err)

	// Verify message
	waitForQueueMessages(t, rabbitContainer, queueName, 1, 5*time.Second)
}

func TestPublisher_Publish_EmptyMessage(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue
	queueName := "test_empty_message"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Publish empty map
	emptyMsg := map[string]interface{}{}

	err = publisher.Publish(ctx, "", queueName, emptyMsg)
	assert.NoError(t, err)

	// Verify message
	waitForQueueMessages(t, rabbitContainer, queueName, 1, 5*time.Second)
}

func TestPublisher_Publish_StringMessage(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup queue
	queueName := "test_string_message"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)
	defer publisher.Close()

	// Publish string message (will be JSON-encoded as a string)
	stringMsg := "Hello, World!"

	err = publisher.Publish(ctx, "", queueName, stringMsg)
	assert.NoError(t, err)

	// Verify message
	waitForQueueMessages(t, rabbitContainer, queueName, 1, 5*time.Second)
}

func TestPublisher_Lifecycle(t *testing.T) {
	rabbitContainer, connectionString := setupRabbitContainer(t)
	ctx := context.Background()

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	defer conn.Close()

	// Setup
	queueName := "test_publisher_lifecycle"
	queueConfig := QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		Args:       nil,
	}
	_, err = conn.DeclareQueue(queueConfig)
	require.NoError(t, err)

	// Create publisher
	publisher := NewPublisher(conn)
	require.NotNil(t, publisher)

	// Publish multiple messages
	for i := 0; i < 5; i++ {
		err = publisher.Publish(ctx, "", queueName, map[string]int{"count": i})
		require.NoError(t, err)
	}

	// Verify messages
	waitForQueueMessages(t, rabbitContainer, queueName, 5, 5*time.Second)

	// Close publisher
	err = publisher.Close()
	assert.NoError(t, err)

	// Publisher can still be used after Close (it doesn't close the connection)
	err = publisher.Publish(ctx, "", queueName, map[string]string{"after": "close"})
	assert.NoError(t, err)
}
