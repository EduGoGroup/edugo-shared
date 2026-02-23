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

func TestConnectRedis_InvalidURL(t *testing.T) {
	_, err := ConnectRedis(RedisConfig{URL: "not-a-valid-redis-url"})
	if err == nil {
		t.Fatal("expected error for invalid redis URL")
	}
	if !strings.Contains(err.Error(), "parsing redis URL") {
		t.Fatalf("expected parsing redis URL error, got: %v", err)
	}
}

func TestConnectRedis_PingFailure(t *testing.T) {
	client, err := ConnectRedis(RedisConfig{URL: "redis://127.0.0.1:1/0"})
	if err == nil {
		if client != nil {
			_ = client.Close()
		}
		t.Fatal("expected ping failure error")
	}
	if !strings.Contains(err.Error(), "pinging redis") {
		t.Fatalf("expected pinging redis error, got: %v", err)
	}
}

func TestCacheService_Get_ReturnsRedisError(t *testing.T) {
	client := newUnavailableRedisClient()
	t.Cleanup(func() {
		_ = client.Close()
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
		_ = client.Close()
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
		_ = client.Close()
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
		_ = client.Close()
	})

	service := NewCacheService(client)

	if err := service.Delete(context.Background()); err != nil {
		t.Fatalf("expected nil error when deleting no keys, got: %v", err)
	}
}

func TestCacheService_Delete_WithKeys_ReturnsRedisError(t *testing.T) {
	client := newUnavailableRedisClient()
	t.Cleanup(func() {
		_ = client.Close()
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
		_ = client.Close()
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
