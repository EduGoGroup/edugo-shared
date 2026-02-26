package rabbit

import (
	"context"
	"errors"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConsumer_Consume_Unit(t *testing.T) {
	mockChannel := new(MockChannel)
	conn := &Connection{
		channel: mockChannel,
	}

	config := ConsumerConfig{
		Name:    "test_consumer",
		AutoAck: true,
	}
	consumer := NewConsumer(conn, config)

	// Create a channel to simulate message delivery
	deliveryChan := make(chan amqp.Delivery)

	// Expect Consume call
	mockChannel.On("Consume",
		"test_queue",
		config.Name,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		false,         // noWait
		mock.Anything, // args
	).Return((<-chan amqp.Delivery)(deliveryChan), nil)

	handlerCalled := make(chan bool)
	handler := func(ctx context.Context, body []byte) error {
		assert.Equal(t, []byte("test_message"), body)
		handlerCalled <- true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := consumer.Consume(ctx, "test_queue", handler)
	assert.NoError(t, err)

	// Simulate receiving a message
	go func() {
		deliveryChan <- amqp.Delivery{
			Body: []byte("test_message"),
		}
	}()

	select {
	case <-handlerCalled:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Handler was not called")
	}

	// Clean stop
	cancel()
	consumer.Stop()
	_ = consumer.Wait() // Ignore error from Wait() as we just want to ensure it finishes
	mockChannel.AssertExpectations(t)
}

func TestConsumer_Consume_StartError_Unit(t *testing.T) {
	mockChannel := new(MockChannel)
	conn := &Connection{
		channel: mockChannel,
	}

	config := ConsumerConfig{Name: "test_consumer"}
	consumer := NewConsumer(conn, config)

	expectedErr := errors.New("consume error")

	// Return a nil channel (casted correctly) and the error
	var nilChan <-chan amqp.Delivery
	mockChannel.On("Consume",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(nilChan, expectedErr)

	err := consumer.Consume(context.Background(), "queue", func(ctx context.Context, body []byte) error { return nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to start consuming")
}
