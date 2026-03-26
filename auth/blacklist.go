package auth

import (
	"context"
	"sync"
	"time"
)

// TokenBlacklist checks if a token has been revoked.
type TokenBlacklist interface {
	// Revoke adds a token's JTI to the blacklist until its expiration time.
	Revoke(jti string, expiresAt time.Time)
	// IsRevoked returns true if the token's JTI has been revoked.
	IsRevoked(jti string) bool
}

// InMemoryBlacklist implements TokenBlacklist using sync.Map with TTL cleanup.
type InMemoryBlacklist struct {
	store sync.Map // jti -> expiresAt (time.Time)
	done  chan struct{}
}

// NewInMemoryBlacklist creates a blacklist with a background cleanup goroutine.
// Cancel the context to stop the cleanup goroutine.
func NewInMemoryBlacklist(ctx context.Context) *InMemoryBlacklist {
	b := &InMemoryBlacklist{
		done: make(chan struct{}),
	}
	go b.cleanupLoop(ctx)
	return b
}

// Revoke adds a token's JTI to the blacklist until its expiration time.
func (b *InMemoryBlacklist) Revoke(jti string, expiresAt time.Time) {
	b.store.Store(jti, expiresAt)
}

// IsRevoked returns true if the token's JTI has been revoked.
func (b *InMemoryBlacklist) IsRevoked(jti string) bool {
	_, ok := b.store.Load(jti)
	return ok
}

func (b *InMemoryBlacklist) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	defer close(b.done)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			b.store.Range(func(key, value any) bool {
				if exp, ok := value.(time.Time); ok && now.After(exp) {
					b.store.Delete(key)
				}
				return true
			})
		}
	}
}

// NoOpBlacklist is a blacklist that never revokes anything. Useful for services
// that don't need token revocation or for testing.
type NoOpBlacklist struct{}

// Revoke is a no-op; the token is not stored.
func (n *NoOpBlacklist) Revoke(_ string, _ time.Time) {}

// IsRevoked always returns false.
func (n *NoOpBlacklist) IsRevoked(_ string) bool { return false }
