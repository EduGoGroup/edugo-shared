//nolint:errcheck // Tests: errores de Close() en cleanup se ignoran intencionalmente
package rabbit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unit tests for Publisher logic that doesn't strictly depend on AMQP connection
// (mostly validation or pre-processing, although NewPublisher is simple)

func TestNewPublisher_Unit(t *testing.T) {
	conn := &Connection{}
	pub := NewPublisher(conn)
	assert.NotNil(t, pub)
	assert.IsType(t, &RabbitMQPublisher{}, pub)
}

func TestRabbitMQPublisher_Close_Unit(t *testing.T) {
	// Close is a no-op currently, but good to verify it doesn't panic
	pub := &RabbitMQPublisher{}
	err := pub.Close()
	assert.NoError(t, err)
}

func TestRabbitMQPublisher_PublishWithPriority_ContextCheck_Unit(t *testing.T) {
	// We can't easily mock the connection/channel without an interface,
	// but we can check if it fails early on context error before calling the channel

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	pub := &RabbitMQPublisher{}

	// Should fail because context is canceled
	err := pub.PublishWithPriority(ctx, "ex", "rk", "body", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish message")
	assert.Contains(t, err.Error(), "context canceled")
}

func TestRabbitMQPublisher_PublishWithPriority_MarshalError_Unit(t *testing.T) {
	// Test JSON marshal failure

	pub := &RabbitMQPublisher{}

	// Channel with unsupported type for JSON (e.g. channel)
	badBody := make(chan int)

	err := pub.PublishWithPriority(context.Background(), "ex", "rk", badBody, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal message")
}
