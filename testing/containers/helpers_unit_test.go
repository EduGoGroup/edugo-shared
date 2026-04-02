package containers

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Tests para WaitForHealthy (lógica pura, sin Docker)
// =============================================================================

// TestWaitForHealthy_ImmediateSuccess verifica success inmediato
func TestWaitForHealthy_ImmediateSuccess(t *testing.T) {
	ctx := context.Background()

	healthCheck := func() error {
		return nil
	}

	err := WaitForHealthy(ctx, healthCheck, 5*time.Second, 100*time.Millisecond)
	assert.NoError(t, err)
}

// TestWaitForHealthy_SuccessAfterRetries verifica success después de reintentos
func TestWaitForHealthy_SuccessAfterRetries(t *testing.T) {
	ctx := context.Background()

	attempts := atomic.Int32{}
	healthCheck := func() error {
		count := attempts.Add(1)
		if count < 3 {
			return errors.New("not ready yet")
		}
		return nil
	}

	start := time.Now()
	err := WaitForHealthy(ctx, healthCheck, 5*time.Second, 100*time.Millisecond)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, attempts.Load(), int32(3))
	assert.GreaterOrEqual(t, duration, 200*time.Millisecond, "Debe esperar al menos 2 intervalos")
}

// TestWaitForHealthy_Timeout verifica timeout
func TestWaitForHealthy_Timeout(t *testing.T) {
	ctx := context.Background()

	healthCheck := func() error {
		return errors.New("service not healthy")
	}

	start := time.Now()
	err := WaitForHealthy(ctx, healthCheck, 500*time.Millisecond, 100*time.Millisecond)
	duration := time.Since(start)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout esperando health check")
	assert.GreaterOrEqual(t, duration, 500*time.Millisecond)
	assert.Less(t, duration, 2*time.Second)
}

// TestWaitForHealthy_ContextCancellation verifica cancelación por contexto
func TestWaitForHealthy_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	healthCheck := func() error {
		return errors.New("not ready")
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := WaitForHealthy(ctx, healthCheck, 5*time.Second, 100*time.Millisecond)
	duration := time.Since(start)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contexto cancelado")
	assert.Less(t, duration, 1*time.Second, "Debe cancelarse rápidamente")
}

// TestWaitForHealthy_ImmediateCancellation verifica cancelación inmediata
func TestWaitForHealthy_ImmediateCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	healthCheck := func() error {
		return errors.New("not ready")
	}

	err := WaitForHealthy(ctx, healthCheck, 5*time.Second, 100*time.Millisecond)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contexto cancelado")
}

// TestWaitForHealthy_ZeroTimeout verifica comportamiento con timeout cero
func TestWaitForHealthy_ZeroTimeout(t *testing.T) {
	ctx := context.Background()

	healthCheck := func() error {
		return errors.New("not ready")
	}

	err := WaitForHealthy(ctx, healthCheck, 0, 100*time.Millisecond)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout esperando health check")
}

// =============================================================================
// Tests para RetryOperation (lógica pura, sin Docker)
// =============================================================================

// TestRetryOperation_ImmediateSuccess verifica success inmediato
func TestRetryOperation_ImmediateSuccess(t *testing.T) {
	attempts := atomic.Int32{}

	operation := func() error {
		attempts.Add(1)
		return nil
	}

	err := RetryOperation(operation, 3, 100*time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), attempts.Load(), "Solo debe intentar una vez")
}

// TestRetryOperation_SuccessAfterRetries verifica success después de reintentos
func TestRetryOperation_SuccessAfterRetries(t *testing.T) {
	attempts := atomic.Int32{}

	operation := func() error {
		count := attempts.Add(1)
		if count < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	start := time.Now()
	err := RetryOperation(operation, 5, 50*time.Millisecond)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, int32(3), attempts.Load())
	assert.GreaterOrEqual(t, duration, 100*time.Millisecond, "Debe esperar delay entre intentos")
}

// TestRetryOperation_AllFailed verifica fallo después de todos los intentos
func TestRetryOperation_AllFailed(t *testing.T) {
	attempts := atomic.Int32{}
	expectedError := errors.New("permanent error")

	operation := func() error {
		attempts.Add(1)
		return expectedError
	}

	maxRetries := 5
	err := RetryOperation(operation, maxRetries, 10*time.Millisecond)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operación falló después de")
	assert.Contains(t, err.Error(), "intentos")
	assert.Equal(t, int32(maxRetries), attempts.Load())
}

// TestRetryOperation_ZeroRetries verifica comportamiento con cero reintentos
func TestRetryOperation_ZeroRetries(t *testing.T) {
	attempts := atomic.Int32{}

	operation := func() error {
		attempts.Add(1)
		return errors.New("error")
	}

	err := RetryOperation(operation, 0, 100*time.Millisecond)

	assert.Error(t, err)
	assert.Equal(t, int32(0), attempts.Load())
}

// TestRetryOperation_OneRetry verifica con un solo reintento
func TestRetryOperation_OneRetry(t *testing.T) {
	attempts := atomic.Int32{}

	operation := func() error {
		attempts.Add(1)
		return errors.New("always fails")
	}

	err := RetryOperation(operation, 1, 10*time.Millisecond)

	assert.Error(t, err)
	assert.Equal(t, int32(1), attempts.Load())
}

// TestRetryOperation_DelayBetweenAttempts verifica delay entre intentos
func TestRetryOperation_DelayBetweenAttempts(t *testing.T) {
	attempts := atomic.Int32{}
	timestamps := []time.Time{}

	operation := func() error {
		attempts.Add(1)
		timestamps = append(timestamps, time.Now())
		return errors.New("error")
	}

	delay := 100 * time.Millisecond
	_ = RetryOperation(operation, 3, delay) //nolint:errcheck // error is intentionally discarded; test validates timing

	assert.Equal(t, int32(3), attempts.Load())
	assert.Len(t, timestamps, 3)

	if len(timestamps) >= 2 {
		diff1 := timestamps[1].Sub(timestamps[0])
		assert.GreaterOrEqual(t, diff1, delay)
	}

	if len(timestamps) >= 3 {
		diff2 := timestamps[2].Sub(timestamps[1])
		assert.GreaterOrEqual(t, diff2, delay)
	}
}

// TestRetryOperation_VariousErrors verifica diferentes tipos de errores
func TestRetryOperation_VariousErrors(t *testing.T) {
	tests := []struct {
		name          string
		errors        []error
		maxRetries    int
		expectedError bool
		expectedCalls int
	}{
		{
			name:          "success on first try",
			errors:        []error{nil},
			maxRetries:    3,
			expectedError: false,
			expectedCalls: 1,
		},
		{
			name:          "success on second try",
			errors:        []error{errors.New("err1"), nil},
			maxRetries:    3,
			expectedError: false,
			expectedCalls: 2,
		},
		{
			name:          "all attempts fail",
			errors:        []error{errors.New("err1"), errors.New("err2"), errors.New("err3")},
			maxRetries:    3,
			expectedError: true,
			expectedCalls: 3,
		},
		{
			name:          "success on last attempt",
			errors:        []error{errors.New("err1"), errors.New("err2"), nil},
			maxRetries:    3,
			expectedError: false,
			expectedCalls: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := atomic.Int32{}
			operation := func() error {
				idx := int(callCount.Load())
				callCount.Add(1)
				if idx < len(tt.errors) {
					return tt.errors[idx]
				}
				return errors.New("unexpected call")
			}

			err := RetryOperation(operation, tt.maxRetries, 1*time.Millisecond)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, int32(tt.expectedCalls), callCount.Load()) //nolint:gosec // expectedCalls es siempre pequeño (1-3)
		})
	}
}

// TestRetryOperation_LastErrorWrapped verifica que el último error se preserva en el wrap
func TestRetryOperation_LastErrorWrapped(t *testing.T) {
	callCount := atomic.Int32{}
	lastErr := errors.New("final specific error")

	operation := func() error {
		count := callCount.Add(1)
		if count < 3 {
			return errors.New("temporary error")
		}
		return lastErr
	}

	err := RetryOperation(operation, 3, 1*time.Millisecond)

	assert.Error(t, err)
	assert.ErrorIs(t, err, lastErr, "El último error debe estar wrapeado en el resultado")
}

// =============================================================================
// Tests para ExecSQLFile (solo la parte de lectura de archivo, sin DB)
// =============================================================================

// TestExecSQLFile_FileNotFound_Unit verifica error cuando archivo no existe (sin DB)
func TestExecSQLFile_FileNotFound_Unit(t *testing.T) {
	ctx := context.Background()

	// Pasar nil como db - no debería llegar a usarse porque el archivo no existe
	err := ExecSQLFile(ctx, nil, "/path/to/nonexistent/file.sql")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error leyendo archivo SQL")
}

// TestExecSQLFile_ValidFilePath_Unit verifica que lee el archivo correctamente
func TestExecSQLFile_ValidFilePath_Unit(t *testing.T) {
	tmpDir := t.TempDir()
	sqlFile := filepath.Join(tmpDir, "test.sql")

	// Crear archivo SQL de prueba
	err := os.WriteFile(sqlFile, []byte("SELECT 1;"), 0600)
	require.NoError(t, err)

	// Sin una DB real, la ejecución fallará pero debería pasar la lectura del archivo.
	// Pasamos nil como DB para verificar que el panic/error viene de la ejecución, no de la lectura.
	// Nota: sql.DB nil causará un panic en ExecContext, lo cual confirma que la lectura fue exitosa.
	assert.Panics(t, func() {
		_ = ExecSQLFile(context.Background(), nil, sqlFile) //nolint:errcheck // panic is expected; error is irrelevant
	}, "Debe hacer panic en ExecContext con nil DB (confirma que la lectura del archivo fue exitosa)")
}
