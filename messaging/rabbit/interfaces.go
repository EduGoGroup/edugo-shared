package rabbit

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ChannelInterface define la interfaz para las operaciones del canal AMQP
// que son utilizadas por el Publisher y Consumer.
// Esto permite mockear el canal para tests unitarios.
type ChannelInterface interface {
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Qos(prefetchCount, prefetchSize int, global bool) error
	ExchangeDelete(name string, ifUnused, noWait bool) error
	Close() error
}
