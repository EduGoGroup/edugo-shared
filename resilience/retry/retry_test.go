package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (m *mockLogger) Info(_ string, _ ...any)  {}
func (m *mockLogger) Warn(_ string, _ ...any)  {}
func (m *mockLogger) Error(_ string, _ ...any) {}

var errTransient = errors.New("transient error")
var errPermanent = errors.New("permanent error")

func classifierForTest(err error) ErrorType {
	if errors.Is(err, errPermanent) {
		return ErrorTypePermanent
	}
	return ErrorTypeTransient
}

func TestWithRetry_Success(t *testing.T) {
	cfg := Config{
		MaxRetries:      3,
		InitialBackoff:  10 * time.Millisecond,
		MaxBackoff:      100 * time.Millisecond,
		BackoffMultiple: 2.0,
		Logger:          &mockLogger{},
	}

	callCount := 0
	err := WithRetry(context.Background(), cfg, func() error {
		callCount++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestWithRetry_SuccessAfterRetries(t *testing.T) {
	cfg := Config{
		MaxRetries:      3,
		InitialBackoff:  10 * time.Millisecond,
		MaxBackoff:      100 * time.Millisecond,
		BackoffMultiple: 2.0,
		Logger:          &mockLogger{},
	}

	callCount := 0
	err := WithRetry(context.Background(), cfg, func() error {
		callCount++
		if callCount < 3 {
			return errTransient
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, callCount)
}

func TestWithRetry_PermanentError(t *testing.T) {
	cfg := Config{
		MaxRetries:      3,
		InitialBackoff:  10 * time.Millisecond,
		MaxBackoff:      100 * time.Millisecond,
		BackoffMultiple: 2.0,
		Logger:          &mockLogger{},
		Classifier:      classifierForTest,
	}

	callCount := 0
	err := WithRetry(context.Background(), cfg, func() error {
		callCount++
		return errPermanent
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, errPermanent)
	assert.Equal(t, 1, callCount, "should not retry permanent errors")
}

func TestWithRetry_MaxRetriesExceeded(t *testing.T) {
	cfg := Config{
		MaxRetries:      3,
		InitialBackoff:  10 * time.Millisecond,
		MaxBackoff:      100 * time.Millisecond,
		BackoffMultiple: 2.0,
		Logger:          &mockLogger{},
	}

	callCount := 0
	err := WithRetry(context.Background(), cfg, func() error {
		callCount++
		return errTransient
	})

	assert.Error(t, err)
	assert.Equal(t, errTransient, err)
	assert.Equal(t, 4, callCount, "should attempt 1 initial + 3 retries")
}

func TestWithRetry_ContextCancellation(t *testing.T) {
	cfg := Config{
		MaxRetries:      3,
		InitialBackoff:  100 * time.Millisecond,
		MaxBackoff:      1 * time.Second,
		BackoffMultiple: 2.0,
		Logger:          &mockLogger{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	callCount := 0

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := WithRetry(ctx, cfg, func() error {
		callCount++
		return errTransient
	})

	assert.Error(t, err)
	assert.True(t, IsContextError(err))
	assert.LessOrEqual(t, callCount, 2)
}

func TestWithRetry_ExponentialBackoff(t *testing.T) {
	cfg := Config{
		MaxRetries:      3,
		InitialBackoff:  50 * time.Millisecond,
		MaxBackoff:      500 * time.Millisecond,
		BackoffMultiple: 2.0,
		Logger:          &mockLogger{},
	}

	callTimes := []time.Time{}
	err := WithRetry(context.Background(), cfg, func() error {
		callTimes = append(callTimes, time.Now())
		return errTransient
	})

	assert.Error(t, err)
	assert.Equal(t, 4, len(callTimes))

	if len(callTimes) >= 2 {
		diff1 := callTimes[1].Sub(callTimes[0])
		assert.GreaterOrEqual(t, diff1, 40*time.Millisecond)
	}
	if len(callTimes) >= 3 {
		diff2 := callTimes[2].Sub(callTimes[1])
		assert.GreaterOrEqual(t, diff2, 90*time.Millisecond)
	}
}

func TestWithRetry_BackoffCap(t *testing.T) {
	cfg := Config{
		MaxRetries:      5,
		InitialBackoff:  10 * time.Millisecond,
		MaxBackoff:      50 * time.Millisecond,
		BackoffMultiple: 2.0,
		Logger:          &mockLogger{},
	}

	callTimes := []time.Time{}
	err := WithRetry(context.Background(), cfg, func() error {
		callTimes = append(callTimes, time.Now())
		return errTransient
	})

	assert.Error(t, err)
	assert.Equal(t, 6, len(callTimes))

	if len(callTimes) >= 5 {
		diff := callTimes[4].Sub(callTimes[3])
		assert.LessOrEqual(t, diff, 70*time.Millisecond)
	}
}

func TestWithRetry_NilLogger(t *testing.T) {
	cfg := Config{
		MaxRetries:      2,
		InitialBackoff:  10 * time.Millisecond,
		MaxBackoff:      50 * time.Millisecond,
		BackoffMultiple: 2.0,
	}

	callCount := 0
	err := WithRetry(context.Background(), cfg, func() error {
		callCount++
		if callCount < 2 {
			return errTransient
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func TestWithRetry_NilClassifier_AllTransient(t *testing.T) {
	cfg := Config{
		MaxRetries:      3,
		InitialBackoff:  10 * time.Millisecond,
		MaxBackoff:      50 * time.Millisecond,
		BackoffMultiple: 2.0,
		Logger:          &mockLogger{},
	}

	callCount := 0
	_ = WithRetry(context.Background(), cfg, func() error {
		callCount++
		return errPermanent // sin classifier, se trata como transient
	})

	assert.Equal(t, 4, callCount, "without classifier all errors are transient")
}

func TestIsContextError(t *testing.T) {
	assert.True(t, IsContextError(context.Canceled))
	assert.True(t, IsContextError(context.DeadlineExceeded))
	assert.False(t, IsContextError(errors.New("generic")))
	assert.False(t, IsContextError(nil))
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, 3, cfg.MaxRetries)
	assert.Equal(t, 500*time.Millisecond, cfg.InitialBackoff)
	assert.Equal(t, 10*time.Second, cfg.MaxBackoff)
	assert.Equal(t, 2.0, cfg.BackoffMultiple)
	assert.Nil(t, cfg.Logger)
	assert.Nil(t, cfg.Classifier)
}
