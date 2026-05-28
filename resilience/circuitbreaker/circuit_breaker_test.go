package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	// testCircuitBreakerTimeout es el timeout configurado para el circuit breaker en tests
	testCircuitBreakerTimeout = 100 * time.Millisecond

	// testWaitForTimeout es el tiempo de espera para asegurar que el timeout del circuit breaker haya expirado
	// Debe ser mayor que testCircuitBreakerTimeout para garantizar la expiracion
	testWaitForTimeout = testCircuitBreakerTimeout + 50*time.Millisecond
)

func TestNew(t *testing.T) {
	// Arrange
	config := DefaultConfig("test")

	// Act
	cb := New(config)

	// Assert
	assert.NotNil(t, cb)
	assert.Equal(t, StateClosed, cb.State())
	assert.Equal(t, uint32(0), cb.Failures())
	assert.Equal(t, uint32(0), cb.Successes())
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	// Arrange
	cb := New(DefaultConfig("test"))
	ctx := context.Background()

	// Act
	err := cb.Execute(ctx, func(ctx context.Context) error {
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, StateClosed, cb.State())
	assert.Equal(t, uint32(0), cb.Failures())
}

func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	// Arrange
	cb := New(DefaultConfig("test"))
	ctx := context.Background()
	expectedErr := errors.New("test error")

	// Act
	err := cb.Execute(ctx, func(ctx context.Context) error {
		return expectedErr
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, StateClosed, cb.State())
	assert.Equal(t, uint32(1), cb.Failures())
}

func TestCircuitBreaker_OpensAfterMaxFailures(t *testing.T) {
	// Arrange
	config := DefaultConfig("test")
	config.MaxFailures = 3
	cb := New(config)
	ctx := context.Background()

	// Act - Generar 3 fallos
	for i := 0; i < 3; i++ {
		_ = cb.Execute(ctx, func(ctx context.Context) error {
			return errors.New("test error")
		})
	}

	// Assert
	assert.Equal(t, StateOpen, cb.State())
	assert.Equal(t, uint32(3), cb.Failures())
}

func TestCircuitBreaker_RejectsWhenOpen(t *testing.T) {
	// Arrange
	config := DefaultConfig("test")
	config.MaxFailures = 1
	cb := New(config)
	ctx := context.Background()

	// Abrir el circuit
	_ = cb.Execute(ctx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	// Act - Intentar ejecutar con circuit abierto
	err := cb.Execute(ctx, func(ctx context.Context) error {
		return nil
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrCircuitOpen, err)
	assert.Equal(t, StateOpen, cb.State())
}

func TestCircuitBreaker_TransitionsToHalfOpen(t *testing.T) {
	// Arrange
	config := DefaultConfig("test")
	config.MaxFailures = 1
	config.Timeout = testCircuitBreakerTimeout
	cb := New(config)
	ctx := context.Background()

	// Abrir el circuit
	_ = cb.Execute(ctx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	assert.Equal(t, StateOpen, cb.State())

	// Esperar a que pase el timeout
	time.Sleep(testWaitForTimeout)

	// Act - Ejecutar despues del timeout
	err := cb.Execute(ctx, func(ctx context.Context) error {
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, StateHalfOpen, cb.State())
}

func TestCircuitBreaker_ClosesAfterSuccessesInHalfOpen(t *testing.T) {
	// Arrange
	config := DefaultConfig("test")
	config.MaxFailures = 1
	config.Timeout = testCircuitBreakerTimeout
	config.SuccessThreshold = 2
	config.MaxRequests = 5
	cb := New(config)
	ctx := context.Background()

	// Abrir el circuit
	_ = cb.Execute(ctx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	// Esperar a que pase el timeout
	time.Sleep(testWaitForTimeout)

	// Act - Ejecutar con exito en half-open
	for i := 0; i < 2; i++ {
		err := cb.Execute(ctx, func(ctx context.Context) error {
			return nil
		})
		assert.NoError(t, err)
	}

	// Assert
	assert.Equal(t, StateClosed, cb.State())
	assert.Equal(t, uint32(0), cb.Failures())
}

func TestCircuitBreaker_ReopensOnFailureInHalfOpen(t *testing.T) {
	// Arrange
	config := DefaultConfig("test")
	config.MaxFailures = 1
	config.Timeout = testCircuitBreakerTimeout
	config.MaxRequests = 2 // Permitir 2 peticiones en half-open para este test
	cb := New(config)
	ctx := context.Background()

	// Abrir el circuit
	_ = cb.Execute(ctx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	// Esperar a que pase el timeout
	time.Sleep(testWaitForTimeout)

	// Transicionar a half-open con una ejecucion exitosa
	_ = cb.Execute(ctx, func(ctx context.Context) error {
		return nil
	})
	assert.Equal(t, StateHalfOpen, cb.State())

	// Act - Fallar en half-open
	_ = cb.Execute(ctx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	// Assert
	assert.Equal(t, StateOpen, cb.State())
}

func TestCircuitBreaker_LimitsRequestsInHalfOpen(t *testing.T) {
	// Arrange
	config := DefaultConfig("test")
	config.MaxFailures = 1
	config.Timeout = testCircuitBreakerTimeout
	config.MaxRequests = 1
	config.SuccessThreshold = 10 // Alto para que no se cierre inmediatamente
	cb := New(config)
	ctx := context.Background()

	// Abrir el circuit
	_ = cb.Execute(ctx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	// Esperar a que pase el timeout
	time.Sleep(testWaitForTimeout)

	// Transicionar a half-open con primera peticion
	err1 := cb.Execute(ctx, func(ctx context.Context) error {
		return nil
	})
	assert.NoError(t, err1)
	assert.Equal(t, StateHalfOpen, cb.State())

	// Act - La segunda peticion debe ser rechazada porque MaxRequests=1
	err2 := cb.Execute(ctx, func(ctx context.Context) error {
		return nil
	})

	// Assert
	assert.Error(t, err2)
	assert.Equal(t, ErrTooManyRequests, err2)
	assert.Equal(t, StateHalfOpen, cb.State())
}

func TestCircuitBreaker_WithMetrics(t *testing.T) {
	// Arrange
	hook := &mockMetricsHook{}
	config := DefaultConfig("test")
	config.MaxFailures = 1
	cb := New(config).WithMetrics(hook)
	ctx := context.Background()

	// Assert - WithMetrics should have called SetState for initial state
	assert.Equal(t, 1, hook.setStateCalls)
	assert.Equal(t, StateClosed, hook.lastState)

	// Act - Trigger a transition to Open
	_ = cb.Execute(ctx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	// Assert
	assert.Equal(t, StateOpen, cb.State())
	assert.Equal(t, 2, hook.setStateCalls)       // initial + open
	assert.Equal(t, 1, hook.transitionCalls)      // closed -> open
	assert.Equal(t, StateClosed, hook.lastFrom)
	assert.Equal(t, StateOpen, hook.lastTo)
}

func TestCircuitBreaker_NilMetrics(t *testing.T) {
	// Arrange - no metrics hook, should not panic
	config := DefaultConfig("test")
	config.MaxFailures = 1
	cb := New(config)
	ctx := context.Background()

	// Act
	_ = cb.Execute(ctx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	// Assert - just verify it didn't panic and state is correct
	assert.Equal(t, StateOpen, cb.State())
}

func TestCircuitBreaker_ConfigMetrics(t *testing.T) {
	// Arrange - set metrics via Config
	hook := &mockMetricsHook{}
	config := DefaultConfig("test")
	config.Metrics = hook

	// Act
	cb := New(config)

	// Assert - New should have called SetState for initial state
	assert.NotNil(t, cb)
	assert.Equal(t, 1, hook.setStateCalls)
	assert.Equal(t, StateClosed, hook.lastState)
}

// mockMetricsHook es un mock para verificar las llamadas a metricas
type mockMetricsHook struct {
	setStateCalls   int
	transitionCalls int
	lastState       State
	lastFrom        State
	lastTo          State
}

func (m *mockMetricsHook) SetState(name string, state State) {
	m.setStateCalls++
	m.lastState = state
}

func (m *mockMetricsHook) RecordTransition(name string, from, to State) {
	m.transitionCalls++
	m.lastFrom = from
	m.lastTo = to
}
