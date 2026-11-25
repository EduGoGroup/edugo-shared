package rabbit

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTruncateError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		assert.Empty(t, truncateError(nil, MaxErrorHeaderBytes))
	})

	t.Run("empty error", func(t *testing.T) {
		assert.Empty(t, truncateError(errors.New(""), MaxErrorHeaderBytes))
	})

	t.Run("short error is untouched", func(t *testing.T) {
		assert.Equal(t, "boom", truncateError(errors.New("boom"), MaxErrorHeaderBytes))
	})

	t.Run("exact limit is untouched", func(t *testing.T) {
		text := strings.Repeat("a", MaxErrorHeaderBytes)
		assert.Equal(t, text, truncateError(errors.New(text), MaxErrorHeaderBytes))
	})

	t.Run("long error is cut and marked", func(t *testing.T) {
		got := truncateError(errors.New(strings.Repeat("a", 5000)), MaxErrorHeaderBytes)
		assert.Len(t, got, MaxErrorHeaderBytes)
		assert.True(t, strings.HasSuffix(got, truncationMark))
	})

	t.Run("cut lands on a rune boundary", func(t *testing.T) {
		// "ñ" ocupa 2 bytes: con un prefijo de 1 byte, el corte cae en mitad de una runa.
		got := truncateError(errors.New("a"+strings.Repeat("ñ", 2000)), MaxErrorHeaderBytes)
		assert.LessOrEqual(t, len(got), MaxErrorHeaderBytes)
		assert.True(t, utf8.ValidString(got))
		assert.True(t, strings.HasSuffix(got, truncationMark))
	})

	t.Run("limit smaller than the mark", func(t *testing.T) {
		assert.Equal(t, truncationMark, truncateError(errors.New("a long error"), 3))
	})
}

func TestDeadLetterHeaders(t *testing.T) {
	now := time.Date(2026, 11, 25, 10, 0, 0, 0, time.UTC)

	t.Run("first delivery", func(t *testing.T) {
		original := amqp.Table{"custom": "value"}
		msg := &amqp.Delivery{Headers: original, Exchange: "ex", RoutingKey: "rk"}

		got := deadLetterHeaders(msg, DeadReasonPermanent, errors.New("boom"), now)

		assert.Equal(t, "value", got["custom"])
		assert.Equal(t, "ex", got[HeaderOriginalExchange])
		assert.Equal(t, "rk", got[HeaderOriginalRoutingKey])
		assert.Equal(t, now.Unix(), got[HeaderFailedAt])
		assert.Equal(t, 0, got[HeaderRetryCount])
		assert.Equal(t, "boom", got[HeaderError])
		assert.Equal(t, DeadReasonPermanent, got[HeaderDeadReason])
		assert.Equal(t, "ex", got[HeaderFirstExchange])
		assert.Equal(t, "rk", got[HeaderFirstRoutingKey])
		assert.Equal(t, amqp.Table{"custom": "value"}, original, "no debe mutar las cabeceras de la entrega")
	})

	t.Run("after retries keeps the first route", func(t *testing.T) {
		msg := &amqp.Delivery{
			Headers: amqp.Table{
				HeaderRetryCount:      int32(3),
				HeaderFirstExchange:   "ex",
				HeaderFirstRoutingKey: "rk",
			},
			Exchange:   "",
			RoutingKey: "queue",
		}

		got := deadLetterHeaders(msg, DeadReasonMaxRetries, errors.New("still failing"), now)

		assert.Empty(t, got[HeaderOriginalExchange])
		assert.Equal(t, "queue", got[HeaderOriginalRoutingKey])
		assert.Equal(t, 3, got[HeaderRetryCount])
		assert.Equal(t, DeadReasonMaxRetries, got[HeaderDeadReason])
		assert.Equal(t, "ex", got[HeaderFirstExchange])
		assert.Equal(t, "rk", got[HeaderFirstRoutingKey])
	})

	t.Run("nil headers and nil error", func(t *testing.T) {
		msg := &amqp.Delivery{RoutingKey: "queue"}

		got := deadLetterHeaders(msg, DeadReasonMaxRetries, nil, now)

		assert.Empty(t, got[HeaderError])
		assert.Nil(t, msg.Headers)
		require.NoError(t, got.Validate())
	})
}

func TestFailureDecision(t *testing.T) {
	tests := []struct {
		name       string
		dlqEnabled bool
		permanent  bool
		retries    int
		maxRetries int
		wantAction failureAction
		wantReason string
	}{
		{"permanent dies at once", true, true, 0, 3, actionDeadLetter, DeadReasonPermanent},
		{"permanent dies even after retries", true, true, 2, 3, actionDeadLetter, DeadReasonPermanent},
		{"transient under the max retries", true, false, 0, 3, actionRetry, ""},
		{"transient on the last retry slot", true, false, 2, 3, actionRetry, ""},
		{"transient at the max dies", true, false, 3, 3, actionDeadLetter, DeadReasonMaxRetries},
		{"zero max retries dies at once", true, false, 0, 0, actionDeadLetter, DeadReasonMaxRetries},
		{"dlq disabled requeues a transient", false, false, 0, 3, actionRequeue, ""},
		{"dlq disabled ignores permanent", false, true, 0, 3, actionRequeue, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, reason := decideFailure(tt.dlqEnabled, tt.permanent, tt.retries, tt.maxRetries)
			assert.Equal(t, tt.wantAction, action)
			assert.Equal(t, tt.wantReason, reason)
		})
	}
}

func TestDeliveryInfo(t *testing.T) {
	t.Run("absent outside the consumer", func(t *testing.T) {
		_, ok := DeliveryFromContext(context.Background())
		assert.False(t, ok)
	})

	t.Run("round trip", func(t *testing.T) {
		want := DeliveryInfo{Queue: "q", Retry: 1, MaxRetries: 3, DeadLetterEnabled: true}
		got, ok := DeliveryFromContext(withDeliveryInfo(context.Background(), want))
		require.True(t, ok)
		assert.Equal(t, want, got)
	})

	t.Run("last attempt", func(t *testing.T) {
		assert.False(t, DeliveryInfo{Retry: 2, MaxRetries: 3, DeadLetterEnabled: true}.IsLastAttempt())
		assert.True(t, DeliveryInfo{Retry: 3, MaxRetries: 3, DeadLetterEnabled: true}.IsLastAttempt())
		assert.False(t, DeliveryInfo{Retry: 3, MaxRetries: 3}.IsLastAttempt())
	})
}

func TestNewFailureEvent(t *testing.T) {
	cause := errors.New("boom")
	msg := &amqp.Delivery{Body: []byte(`{"a":1}`)}

	got := newFailureEvent("q", msg, 2, 3, cause)

	assert.Equal(t, "q", got.Queue)
	assert.Equal(t, 2, got.Retry)
	assert.Equal(t, 3, got.MaxRetries)
	assert.Equal(t, 7, got.BodySize)
	assert.Len(t, got.BodyHash, bodyHashLength)
	assert.Equal(t, bodyHash([]byte(`{"a":1}`)), got.BodyHash)
	assert.NotEqual(t, bodyHash([]byte(`{"a":2}`)), got.BodyHash)
	assert.ErrorIs(t, got.Err, cause)
}
