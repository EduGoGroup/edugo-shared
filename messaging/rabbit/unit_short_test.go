package rabbit

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ackRecorder struct {
	ackCalls  int32
	nackCalls int32
	ackErr    error
	nackErr   error
}

func (a *ackRecorder) Ack(uint64, bool) error {
	atomic.AddInt32(&a.ackCalls, 1)
	return a.ackErr
}

func (a *ackRecorder) Nack(uint64, bool, bool) error {
	atomic.AddInt32(&a.nackCalls, 1)
	return a.nackErr
}

func (a *ackRecorder) Reject(uint64, bool) error {
	return nil
}

func newUnitConsumer(cfg ConsumerConfig) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		config:  cfg,
		errChan: make(chan error, 1),
		stopCh:  make(chan struct{}),
	}
}

func TestConnection_Unit_NilSafetyAndClosedHealthCheck(t *testing.T) {
	conn := &Connection{}

	if conn.GetChannel() != nil {
		t.Fatal("expected nil channel")
	}
	if conn.GetConnection() != nil {
		t.Fatal("expected nil connection")
	}
	if !conn.IsClosed() {
		t.Fatal("expected nil connection to be considered closed")
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("expected close to succeed with nil fields, got: %v", err)
	}
	if err := conn.HealthCheck(); err == nil {
		t.Fatal("expected health check error for closed connection")
	}
}

func TestNewConsumer_Unit(t *testing.T) {
	cfg := ConsumerConfig{Name: "unit-consumer", AutoAck: false}
	consumer := NewConsumer(nil, cfg)

	rabbitConsumer, ok := consumer.(*RabbitMQConsumer)
	if !ok {
		t.Fatalf("expected *RabbitMQConsumer, got %T", consumer)
	}
	if rabbitConsumer.IsRunning() {
		t.Fatal("new consumer should not be running")
	}
	if rabbitConsumer.Errors() == nil {
		t.Fatal("expected errors channel to be initialized")
	}
}

func TestConsumer_Consume_WhenAlreadyRunning(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{})
	consumer.running = true

	err := consumer.Consume(context.Background(), "q", func(context.Context, []byte) error { return nil })
	if err == nil {
		t.Fatal("expected error when consumer is already running")
	}
	if !strings.Contains(err.Error(), "consumer already running") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConsumer_ConsumeWithDLQ_WhenAlreadyRunning(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{})
	consumer.running = true

	err := consumer.ConsumeWithDLQ(context.Background(), "q", func(context.Context, []byte) error { return nil })
	if err == nil {
		t.Fatal("expected error when consumer is already running")
	}
	if !strings.Contains(err.Error(), "consumer already running") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConsumer_ProcessBasicMessage_SuccessWithManualAck(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{AutoAck: false})

	var ackCalls int32
	var nackCalls int32
	delivery := amqpDelivery{
		Body:        []byte(`{"ok":true}`),
		DeliveryTag: 1,
		Ack: func(bool) error {
			atomic.AddInt32(&ackCalls, 1)
			return nil
		},
		Nack: func(bool, bool) error {
			atomic.AddInt32(&nackCalls, 1)
			return nil
		},
	}

	consumer.processBasicMessage(context.Background(), "queue.unit", delivery, func(context.Context, []byte) error {
		return nil
	})

	if got := atomic.LoadInt32(&ackCalls); got != 1 {
		t.Fatalf("expected ack to be called once, got %d", got)
	}
	if got := atomic.LoadInt32(&nackCalls); got != 0 {
		t.Fatalf("expected nack to not be called, got %d", got)
	}
}

func TestConsumer_ProcessBasicMessage_HandlerErrorTriggersNack(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{AutoAck: false})

	var ackCalls int32
	var nackCalls int32
	delivery := amqpDelivery{
		Body:        []byte("fail"),
		DeliveryTag: 2,
		Ack: func(bool) error {
			atomic.AddInt32(&ackCalls, 1)
			return nil
		},
		Nack: func(bool, bool) error {
			atomic.AddInt32(&nackCalls, 1)
			return nil
		},
	}

	consumer.processBasicMessage(context.Background(), "queue.unit", delivery, func(context.Context, []byte) error {
		return errors.New("handler failed")
	})

	if got := atomic.LoadInt32(&ackCalls); got != 0 {
		t.Fatalf("expected ack to not be called, got %d", got)
	}
	if got := atomic.LoadInt32(&nackCalls); got != 1 {
		t.Fatalf("expected nack to be called once, got %d", got)
	}
}

func TestConsumer_ProcessBasicMessage_AutoAckSkipsAckAndNack(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{AutoAck: true})

	var ackCalls int32
	var nackCalls int32
	delivery := amqpDelivery{
		Body:        []byte("auto"),
		DeliveryTag: 3,
		Ack: func(bool) error {
			atomic.AddInt32(&ackCalls, 1)
			return nil
		},
		Nack: func(bool, bool) error {
			atomic.AddInt32(&nackCalls, 1)
			return nil
		},
	}

	consumer.processBasicMessage(context.Background(), "queue.unit", delivery, func(context.Context, []byte) error {
		return nil
	})

	if got := atomic.LoadInt32(&ackCalls); got != 0 {
		t.Fatalf("expected ack to not be called when auto-ack is enabled, got %d", got)
	}
	if got := atomic.LoadInt32(&nackCalls); got != 0 {
		t.Fatalf("expected nack to not be called when auto-ack is enabled, got %d", got)
	}
}

func TestConsumer_Wait_ReturnsAsyncError_Unit(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{})
	expectedErr := errors.New("async failure")
	consumer.errChan <- expectedErr

	err := consumer.Wait()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected async error, got: %v", err)
	}
}

func TestConsumer_Wait_NoError(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{})

	if err := consumer.Wait(); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestConsumer_Stop_Close_Errors_IsRunning_Unit(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{})
	if consumer.IsRunning() {
		t.Fatal("expected consumer not running by default")
	}
	if consumer.Errors() == nil {
		t.Fatal("expected non-nil error channel")
	}

	consumer.wg.Add(1)
	go func() {
		defer consumer.wg.Done()
		<-consumer.stopCh
	}()

	consumer.Stop()
	consumer.Stop() // idempotent

	select {
	case <-consumer.stopCh:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected stop channel to be closed")
	}

	if err := consumer.Close(); err != nil {
		t.Fatalf("expected close to succeed, got: %v", err)
	}
}

func TestConsumerDLQ_ProcessMessage_AutoAckReturnsEarly(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{
		AutoAck: true,
	})

	rec := &ackRecorder{}
	msg := amqp.Delivery{
		Body:         []byte("test"),
		DeliveryTag:  10,
		Acknowledger: rec,
	}

	consumer.processMessage(context.Background(), nil, "queue.unit", func(context.Context, []byte) error {
		return errors.New("ignored")
	}, msg)

	if got := atomic.LoadInt32(&rec.ackCalls); got != 0 {
		t.Fatalf("expected no ack calls, got %d", got)
	}
	if got := atomic.LoadInt32(&rec.nackCalls); got != 0 {
		t.Fatalf("expected no nack calls, got %d", got)
	}
}

func TestConsumerDLQ_ProcessMessage_SuccessAcks(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{
		AutoAck: false,
		DLQ:     DLQConfig{Enabled: false},
	})

	rec := &ackRecorder{}
	msg := amqp.Delivery{
		Body:         []byte("ok"),
		DeliveryTag:  11,
		Acknowledger: rec,
	}

	consumer.processMessage(context.Background(), nil, "queue.unit", func(context.Context, []byte) error {
		return nil
	}, msg)

	if got := atomic.LoadInt32(&rec.ackCalls); got != 1 {
		t.Fatalf("expected ack once, got %d", got)
	}
	if got := atomic.LoadInt32(&rec.nackCalls); got != 0 {
		t.Fatalf("expected no nack calls, got %d", got)
	}
}

func TestConsumerDLQ_ProcessMessage_ErrorWithoutDLQNacks(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{
		AutoAck: false,
		DLQ:     DLQConfig{Enabled: false},
	})

	rec := &ackRecorder{}
	msg := amqp.Delivery{
		Body:         []byte("boom"),
		DeliveryTag:  12,
		Acknowledger: rec,
	}

	consumer.processMessage(context.Background(), nil, "queue.unit", func(context.Context, []byte) error {
		return errors.New("handler error")
	}, msg)

	if got := atomic.LoadInt32(&rec.ackCalls); got != 0 {
		t.Fatalf("expected no ack calls, got %d", got)
	}
	if got := atomic.LoadInt32(&rec.nackCalls); got != 1 {
		t.Fatalf("expected nack once, got %d", got)
	}
}

func TestConsumerDLQ_HandleProcessingError_CanceledContextNacks(t *testing.T) {
	consumer := newUnitConsumer(ConsumerConfig{
		AutoAck: false,
		DLQ: DLQConfig{
			Enabled:               true,
			MaxRetries:            5,
			RetryDelay:            100 * time.Millisecond,
			UseExponentialBackoff: false,
		},
	})

	rec := &ackRecorder{}
	msg := amqp.Delivery{
		Body:         []byte("retry"),
		DeliveryTag:  13,
		Acknowledger: rec,
		Headers: amqp.Table{
			"x-retry-count": int32(0),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // fuerza la rama ctx.Done antes del publish de retry

	consumer.handleProcessingError(ctx, nil, "queue.unit", msg, 0)

	if got := atomic.LoadInt32(&rec.nackCalls); got != 1 {
		t.Fatalf("expected nack once after context cancelation, got %d", got)
	}
}

func TestNewPublisherAndClose_Unit(t *testing.T) {
	publisher := NewPublisher(nil)
	if publisher == nil {
		t.Fatal("expected publisher instance")
	}
	if err := publisher.Close(); err != nil {
		t.Fatalf("expected close to return nil, got: %v", err)
	}
}

func TestPublisher_Publish_ContextCanceled(t *testing.T) {
	publisher := &RabbitMQPublisher{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := publisher.Publish(ctx, "ex", "rk", map[string]string{"k": "v"})
	if err == nil {
		t.Fatal("expected context canceled error")
	}
	if !strings.Contains(err.Error(), "failed to publish message") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPublisher_PublishWithPriority_MarshalError(t *testing.T) {
	publisher := &RabbitMQPublisher{}

	err := publisher.PublishWithPriority(context.Background(), "ex", "rk", make(chan int), 1)
	if err == nil {
		t.Fatal("expected marshal error")
	}
	if !strings.Contains(err.Error(), "failed to marshal message") {
		t.Fatalf("unexpected error: %v", err)
	}
}
