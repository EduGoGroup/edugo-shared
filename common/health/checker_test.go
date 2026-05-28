package health

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckerFunc(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		checker := CheckerFunc(func(ctx context.Context) error {
			return nil
		})
		assert.NoError(t, checker.Check(context.Background()))
	})

	t.Run("unhealthy", func(t *testing.T) {
		checker := CheckerFunc(func(ctx context.Context) error {
			return errors.New("connection refused")
		})
		assert.Error(t, checker.Check(context.Background()))
	})
}

func TestCheckWithTimeout(t *testing.T) {
	t.Run("passes within timeout", func(t *testing.T) {
		checker := CheckerFunc(func(ctx context.Context) error {
			return nil
		})
		err := CheckWithTimeout(context.Background(), checker, 1*time.Second)
		assert.NoError(t, err)
	})

	t.Run("uses default timeout when zero", func(t *testing.T) {
		checker := CheckerFunc(func(ctx context.Context) error {
			deadline, ok := ctx.Deadline()
			assert.True(t, ok)
			assert.WithinDuration(t, time.Now().Add(DefaultTimeout), deadline, 100*time.Millisecond)
			return nil
		})
		err := CheckWithTimeout(context.Background(), checker, 0)
		assert.NoError(t, err)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		checker := CheckerFunc(func(ctx context.Context) error {
			return ctx.Err()
		})
		err := CheckWithTimeout(ctx, checker, 1*time.Second)
		assert.Error(t, err)
	})
}

func TestDefaultTimeout(t *testing.T) {
	assert.Equal(t, 5*time.Second, DefaultTimeout)
}
