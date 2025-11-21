package rabbit

import (
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

func waitForQueueMessages(t *testing.T, container *containers.RabbitMQContainer, queueName string, expected int, timeout time.Duration) amqp.Queue {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for {
		channel, err := container.Channel()
		require.NoError(t, err)

		queueInfo, err := channel.QueueInspect(queueName)
		_ = channel.Close()
		require.NoError(t, err)

		if queueInfo.Messages == expected {
			return queueInfo
		}

		if time.Now().After(deadline) {
			t.Fatalf("Timeout esperando %d mensajes en %s (último valor: %d)", expected, queueName, queueInfo.Messages)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func waitForQueueConsumers(t *testing.T, container *containers.RabbitMQContainer, queueName string, minConsumers int, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for {
		channel, err := container.Channel()
		require.NoError(t, err)

		queueInfo, err := channel.QueueInspect(queueName)
		_ = channel.Close()
		require.NoError(t, err)

		if queueInfo.Consumers >= minConsumers {
			return
		}

		if time.Now().After(deadline) {
			t.Fatalf("Timeout esperando %d consumidores en %s", minConsumers, queueName)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
