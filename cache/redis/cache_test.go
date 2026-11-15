package redis

import (
	"context"
	"strings"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func newUnavailableRedisClient() *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:         "127.0.0.1:1",
		DialTimeout:  20 * time.Millisecond,
		ReadTimeout:  20 * time.Millisecond,
		WriteTimeout: 20 * time.Millisecond,
		PoolTimeout:  20 * time.Millisecond,
		MaxRetries:   0,
	})
}

func TestCacheService_Get_ReturnsRedisError(t *testing.T) {
	client := newUnavailableRedisClient()
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("closing redis client: %v", err)
		}
	})

	service := NewCacheService(client)

	var out map[string]string
	err := service.Get(context.Background(), "user:1", &out)
	if err == nil {
		t.Fatal("expected redis error on Get")
	}
}

func TestCacheService_Set_MarshalError(t *testing.T) {
	client := newUnavailableRedisClient()
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("closing redis client: %v", err)
		}
	})

	service := NewCacheService(client)

	err := service.Set(context.Background(), "key", make(chan int), time.Minute)
	if err == nil {
		t.Fatal("expected marshal error")
	}
	if !strings.Contains(err.Error(), "marshaling cache value") {
		t.Fatalf("expected marshaling error, got: %v", err)
	}
}

func TestCacheService_Set_ReturnsRedisError(t *testing.T) {
	client := newUnavailableRedisClient()
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("closing redis client: %v", err)
		}
	})

	service := NewCacheService(client)

	err := service.Set(context.Background(), "key", map[string]string{"name": "alice"}, time.Minute)
	if err == nil {
		t.Fatal("expected redis error on Set")
	}
}

func TestCacheService_Delete_NoKeys(t *testing.T) {
	client := newUnavailableRedisClient()
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("closing redis client: %v", err)
		}
	})

	service := NewCacheService(client)

	if err := service.Delete(context.Background()); err != nil {
		t.Fatalf("expected nil error when deleting no keys, got: %v", err)
	}
}

func TestCacheService_Delete_WithKeys_ReturnsRedisError(t *testing.T) {
	client := newUnavailableRedisClient()
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("closing redis client: %v", err)
		}
	})

	service := NewCacheService(client)

	err := service.Delete(context.Background(), "key:1", "key:2")
	if err == nil {
		t.Fatal("expected redis error on Delete with keys")
	}
}

func TestCacheService_DeleteByPattern_ScanError(t *testing.T) {
	client := newUnavailableRedisClient()
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("closing redis client: %v", err)
		}
	})

	service := NewCacheService(client)

	err := service.DeleteByPattern(context.Background(), "user:*")
	if err == nil {
		t.Fatal("expected scan error")
	}
	if !strings.Contains(err.Error(), "scanning keys") {
		t.Fatalf("expected scanning keys error, got: %v", err)
	}
}
