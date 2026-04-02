package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryBlacklist_RevokeAndIsRevoked(t *testing.T) {
	ctx := t.Context()

	bl := NewInMemoryBlacklist(ctx)

	assert.False(t, bl.IsRevoked("jti-1"))

	bl.Revoke("jti-1", time.Now().Add(1*time.Hour))

	assert.True(t, bl.IsRevoked("jti-1"))
	assert.False(t, bl.IsRevoked("jti-2"))
}

func TestInMemoryBlacklist_MultipleRevocations(t *testing.T) {
	ctx := t.Context()

	bl := NewInMemoryBlacklist(ctx)

	bl.Revoke("jti-1", time.Now().Add(1*time.Hour))
	bl.Revoke("jti-2", time.Now().Add(1*time.Hour))
	bl.Revoke("jti-3", time.Now().Add(1*time.Hour))

	assert.True(t, bl.IsRevoked("jti-1"))
	assert.True(t, bl.IsRevoked("jti-2"))
	assert.True(t, bl.IsRevoked("jti-3"))
	assert.False(t, bl.IsRevoked("jti-4"))
}

func TestNoOpBlacklist(t *testing.T) {
	bl := &NoOpBlacklist{}

	bl.Revoke("jti-1", time.Now().Add(1*time.Hour))
	assert.False(t, bl.IsRevoked("jti-1"))
}

func TestInMemoryBlacklist_CancelStopsCleanup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	bl := NewInMemoryBlacklist(ctx)
	require.NotNil(t, bl)

	cancel()
	// Wait for cleanup goroutine to finish
	<-bl.done
}
