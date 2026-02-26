package rabbit

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
)

// MockChannel es un mock de ChannelInterface
type MockChannel struct {
	mock.Mock
}

func (m *MockChannel) PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	args := m.Called(ctx, exchange, key, mandatory, immediate, msg)
	return args.Error(0)
}

func (m *MockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	callArgs := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	return callArgs.Get(0).(<-chan amqp.Delivery), callArgs.Error(1)
}

func (m *MockChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	callArgs := m.Called(name, kind, durable, autoDelete, internal, noWait, args)
	return callArgs.Error(0)
}

func (m *MockChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	callArgs := m.Called(name, durable, autoDelete, exclusive, noWait, args)
	return callArgs.Get(0).(amqp.Queue), callArgs.Error(1)
}

func (m *MockChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	callArgs := m.Called(name, key, exchange, noWait, args)
	return callArgs.Error(0)
}

func (m *MockChannel) Qos(prefetchCount, prefetchSize int, global bool) error {
	callArgs := m.Called(prefetchCount, prefetchSize, global)
	return callArgs.Error(0)
}

func (m *MockChannel) ExchangeDelete(name string, ifUnused, noWait bool) error {
	callArgs := m.Called(name, ifUnused, noWait)
	return callArgs.Error(0)
}

func (m *MockChannel) Close() error {
	args := m.Called()
	return args.Error(0)
}
