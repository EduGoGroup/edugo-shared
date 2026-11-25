//go:build integration

//nolint:errcheck // Tests: errores de Close() en cleanup se ignoran intencionalmente
package rabbit

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/testing/containers"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dlqTestTimeout = 15 * time.Second

var errPermanentTest = errors.New("permanent failure")

// dlqHarness arma una cola con su cola muerta sobre el broker de prueba.
type dlqHarness struct {
	container *containers.RabbitMQContainer
	conn      *Connection
	queue     string
	dlq       string
}

func newDLQHarness(t *testing.T, name string) *dlqHarness {
	t.Helper()
	container, connectionString := setupRabbitContainer(t)

	conn, err := Connect(connectionString)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	h := &dlqHarness{container: container, conn: conn, queue: "hu_ia3." + name, dlq: "hu_ia3." + name + ".dlq"}
	_, err = conn.GetChannel().QueueDeclare(h.queue, true, false, false, false, nil)
	require.NoError(t, err)
	return h
}

func (h *dlqHarness) config(maxRetries int) ConsumerConfig {
	return ConsumerConfig{
		Name: h.queue + ".consumer",
		DLQ: DLQConfig{
			Enabled:       true,
			MaxRetries:    maxRetries,
			RetryDelay:    50 * time.Millisecond,
			DLXExchange:   h.queue + ".dlx",
			DLXRoutingKey: h.dlq,
		},
	}
}

func (h *dlqHarness) start(t *testing.T, cfg ConsumerConfig, handler MessageHandler) {
	t.Helper()
	consumer, ok := NewConsumer(h.conn, cfg).(*RabbitMQConsumer)
	require.True(t, ok)

	ctx, cancel := context.WithCancel(context.Background())
	require.NoError(t, consumer.ConsumeWithDLQ(ctx, h.queue, handler))
	t.Cleanup(func() {
		cancel()
		_ = consumer.Close()
	})
}

func (h *dlqHarness) publish(t *testing.T, exchange, routingKey string, body []byte) {
	t.Helper()
	channel, err := h.container.Channel()
	require.NoError(t, err)
	defer func() { _ = channel.Close() }()

	err = channel.PublishWithContext(context.Background(), exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	require.NoError(t, err)
}

// deadMessage espera a que haya un mensaje en la cola muerta y lo saca.
func (h *dlqHarness) deadMessage(t *testing.T) amqp.Delivery {
	t.Helper()
	waitForQueueMessages(t, h.container, h.dlq, 1, dlqTestTimeout)

	channel, err := h.container.Channel()
	require.NoError(t, err)
	defer func() { _ = channel.Close() }()

	msg, ok, err := channel.Get(h.dlq, true)
	require.NoError(t, err)
	require.True(t, ok)
	return msg
}

func TestConsumeWithDLQ_DeadLetterIsPersistentAndCarriesError(t *testing.T) {
	h := newDLQHarness(t, "persistent")
	body := []byte(`{"event_type":"x","n":1}`)

	var calls int32
	h.start(t, h.config(2), func(context.Context, []byte) error {
		atomic.AddInt32(&calls, 1)
		return errors.New("model is down")
	})
	h.publish(t, "", h.queue, body)

	dead := h.deadMessage(t)

	assert.Equal(t, int32(3), atomic.LoadInt32(&calls), "MaxRetries+1 ejecuciones")
	assert.Equal(t, uint8(amqp.Persistent), dead.DeliveryMode)
	assert.Equal(t, "model is down", dead.Headers[HeaderError])
	assert.Equal(t, DeadReasonMaxRetries, dead.Headers[HeaderDeadReason])
	assert.Equal(t, 2, getRetryCount(dead.Headers))
	assert.Equal(t, body, dead.Body)
	waitForQueueMessages(t, h.container, h.queue, 0, dlqTestTimeout)
}

func TestConsumeWithDLQ_PermanentErrorSkipsRetries(t *testing.T) {
	h := newDLQHarness(t, "permanent")

	cfg := h.config(3)
	cfg.IsPermanent = func(err error) bool { return errors.Is(err, errPermanentTest) }

	var calls int32
	h.start(t, cfg, func(context.Context, []byte) error {
		atomic.AddInt32(&calls, 1)
		return errPermanentTest
	})
	h.publish(t, "", h.queue, []byte(`not json`))

	dead := h.deadMessage(t)

	assert.Equal(t, 0, getRetryCount(dead.Headers))
	assert.Equal(t, DeadReasonPermanent, dead.Headers[HeaderDeadReason])
	assert.Equal(t, errPermanentTest.Error(), dead.Headers[HeaderError])
	assert.Equal(t, uint8(amqp.Persistent), dead.DeliveryMode)
	assert.Equal(t, []byte(`not json`), dead.Body)

	// Deja pasar varias esperas de reintento: no debe haber ninguna ejecución más.
	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls), "una sola ejecución")
	waitForQueueMessages(t, h.container, h.queue, 0, dlqTestTimeout)
}

func TestConsumeWithDLQ_NilClassifierKeepsRetrying(t *testing.T) {
	h := newDLQHarness(t, "nil_classifier")

	var calls int32
	h.start(t, h.config(2), func(context.Context, []byte) error {
		atomic.AddInt32(&calls, 1)
		return errPermanentTest
	})
	h.publish(t, "", h.queue, []byte(`{}`))

	dead := h.deadMessage(t)

	assert.Equal(t, int32(3), atomic.LoadInt32(&calls))
	assert.Equal(t, 2, getRetryCount(dead.Headers))
	assert.Equal(t, DeadReasonMaxRetries, dead.Headers[HeaderDeadReason])
}

func TestConsumeWithDLQ_HooksFire(t *testing.T) {
	h := newDLQHarness(t, "hooks")
	body := []byte(`{"hook":true}`)
	cause := errors.New("transient failure")

	var mu sync.Mutex
	var retries, deads []FailureEvent

	cfg := h.config(2)
	cfg.OnRetry = func(ev FailureEvent) {
		mu.Lock()
		defer mu.Unlock()
		retries = append(retries, ev)
	}
	cfg.OnDeadLetter = func(ev FailureEvent) {
		mu.Lock()
		defer mu.Unlock()
		deads = append(deads, ev)
	}

	h.start(t, cfg, func(context.Context, []byte) error { return cause })
	h.publish(t, "", h.queue, body)
	_ = h.deadMessage(t)

	require.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(deads) == 1
	}, dlqTestTimeout, 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	require.Len(t, retries, 2, "una llamada por reintento")
	for i, ev := range retries {
		assert.Equal(t, h.queue, ev.Queue)
		assert.Equal(t, i, ev.Retry)
		assert.Equal(t, 2, ev.MaxRetries)
		assert.Equal(t, 50*time.Millisecond, ev.Backoff)
		assert.Empty(t, ev.Reason)
		assert.ErrorIs(t, ev.Err, cause)
	}

	dead := deads[0]
	assert.Equal(t, h.queue, dead.Queue)
	assert.Equal(t, 2, dead.Retry)
	assert.Equal(t, DeadReasonMaxRetries, dead.Reason)
	assert.ErrorIs(t, dead.Err, cause)
	assert.Equal(t, len(body), dead.BodySize)
	assert.Equal(t, bodyHash(body), dead.BodyHash)
	assert.Zero(t, dead.Backoff)
}

func TestConsumeWithDLQ_HandlerSeesRetryCount(t *testing.T) {
	h := newDLQHarness(t, "retry_count")

	var mu sync.Mutex
	var seen []DeliveryInfo
	var missing int32

	h.start(t, h.config(2), func(ctx context.Context, _ []byte) error {
		info, ok := DeliveryFromContext(ctx)
		if !ok {
			atomic.AddInt32(&missing, 1)
		}
		mu.Lock()
		seen = append(seen, info)
		mu.Unlock()
		return errors.New("always fails")
	})
	h.publish(t, "", h.queue, []byte(`{}`))
	_ = h.deadMessage(t)

	mu.Lock()
	defer mu.Unlock()

	assert.Zero(t, atomic.LoadInt32(&missing))
	require.Len(t, seen, 3)
	for i, info := range seen {
		assert.Equal(t, h.queue, info.Queue)
		assert.Equal(t, i, info.Retry)
		assert.Equal(t, 2, info.MaxRetries)
		assert.Equal(t, i == 2, info.IsLastAttempt(), "ejecución %d", i)
	}
}

func TestConsumeWithDLQ_FirstRouteSurvivesRetries(t *testing.T) {
	h := newDLQHarness(t, "first_route")
	exchange := h.queue + ".exchange"
	routingKey := "first.route.key"

	channel := h.conn.GetChannel()
	require.NoError(t, channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil))
	require.NoError(t, channel.QueueBind(h.queue, routingKey, exchange, false, nil))

	h.start(t, h.config(2), func(context.Context, []byte) error { return errors.New("always fails") })
	h.publish(t, exchange, routingKey, []byte(`{}`))

	dead := h.deadMessage(t)

	assert.Equal(t, exchange, dead.Headers[HeaderFirstExchange])
	assert.Equal(t, routingKey, dead.Headers[HeaderFirstRoutingKey])
	// Las de siempre no cambian: tras un reintento apuntan al exchange por defecto y a la cola.
	assert.Equal(t, "", dead.Headers[HeaderOriginalExchange])
	assert.Equal(t, h.queue, dead.Headers[HeaderOriginalRoutingKey])
	assert.Equal(t, 2, getRetryCount(dead.Headers))
}
