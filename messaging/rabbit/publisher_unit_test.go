package rabbit

import (
	"context"
	"errors"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPublisher_Publish_Unit(t *testing.T) {
	mockChannel := new(MockChannel)
	// Connection with mock channel
	conn := &Connection{
		channel: mockChannel,
	}
	publisher := NewPublisher(conn)

	ctx := context.Background()
	exchange := "test_exchange"
	routingKey := "test_key"
	msgBody := map[string]string{"msg": "hello"}

	// Expect PublishWithContext call
	mockChannel.On("PublishWithContext",
		mock.Anything, // ctx
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		mock.MatchedBy(func(msg amqp.Publishing) bool {
			return msg.ContentType == "application/json" && len(msg.Body) > 0
		}),
	).Return(nil)

	err := publisher.Publish(ctx, exchange, routingKey, msgBody)
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
}

func TestPublisher_PublishWithPriority_Unit(t *testing.T) {
	mockChannel := new(MockChannel)
	conn := &Connection{
		channel: mockChannel,
	}
	publisher := NewPublisher(conn)

	ctx := context.Background()
	exchange := "test_exchange"
	routingKey := "test_key"
	msgBody := map[string]string{"msg": "priority"}
	priority := uint8(5)

	mockChannel.On("PublishWithContext",
		mock.Anything,
		exchange,
		routingKey,
		false,
		false,
		mock.MatchedBy(func(msg amqp.Publishing) bool {
			return msg.Priority == priority
		}),
	).Return(nil)

	err := publisher.PublishWithPriority(ctx, exchange, routingKey, msgBody, priority)
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
}

func TestPublisher_Publish_Error_Unit(t *testing.T) {
	mockChannel := new(MockChannel)
	conn := &Connection{
		channel: mockChannel,
	}
	publisher := NewPublisher(conn)

	ctx := context.Background()
	expectedErr := errors.New("amqp error")

	mockChannel.On("PublishWithContext",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		false,
		false,
		mock.Anything,
	).Return(expectedErr)

	err := publisher.Publish(ctx, "ex", "rk", "body")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish message")
	mockChannel.AssertExpectations(t)
}
