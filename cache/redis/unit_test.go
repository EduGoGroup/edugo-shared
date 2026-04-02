package redis

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// startMiniRedis spins up an in-memory Redis and returns the client + cleanup.
func startMiniRedis(t *testing.T) (*goredis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("starting miniredis: %v", err)
	}
	client := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	t.Cleanup(func() {
		_ = client.Close() //nolint:errcheck,gosec // best-effort cleanup
		mr.Close()
	})
	return client, mr
}

// mustSet is a test helper that calls Set and fails the test on error.
func mustSet(ctx context.Context, t *testing.T, svc CacheService, key string, value any, ttl time.Duration) {
	t.Helper()
	if err := svc.Set(ctx, key, value, ttl); err != nil {
		t.Fatalf("Set(%q): %v", key, err)
	}
}

// ---------------------------------------------------------------------------
// RedisConfig tests
// ---------------------------------------------------------------------------

func TestRedisConfig_ZeroValue(t *testing.T) {
	var cfg RedisConfig
	if cfg.URL != "" {
		t.Fatalf("expected empty URL, got %q", cfg.URL)
	}
}

func TestRedisConfig_WithURL(t *testing.T) {
	cfg := RedisConfig{URL: "redis://localhost:6379/0"}
	if cfg.URL != "redis://localhost:6379/0" {
		t.Fatalf("unexpected URL: %q", cfg.URL)
	}
}

// ---------------------------------------------------------------------------
// ConnectRedis tests (no real server needed for error paths)
// ---------------------------------------------------------------------------

func TestConnectRedis_InvalidURL(t *testing.T) {
	_, err := ConnectRedis(RedisConfig{URL: "not-a-url"})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
	if !strings.Contains(err.Error(), "parsing redis URL") {
		t.Fatalf("expected 'parsing redis URL' in error, got: %v", err)
	}
}

func TestConnectRedis_EmptyURL(t *testing.T) {
	_, err := ConnectRedis(RedisConfig{URL: ""})
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestConnectRedis_PingFailure(t *testing.T) {
	// Use a valid URL format but unreachable address so ParseURL succeeds
	// but Ping fails.
	_, err := ConnectRedis(RedisConfig{URL: "redis://127.0.0.1:1/0"})
	if err == nil {
		t.Fatal("expected ping failure")
	}
	if !strings.Contains(err.Error(), "pinging redis") {
		t.Fatalf("expected 'pinging redis' in error, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// NewCacheService
// ---------------------------------------------------------------------------

func TestNewCacheService_ReturnsNonNil(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)
	if svc == nil {
		t.Fatal("expected non-nil CacheService")
	}
}

// ---------------------------------------------------------------------------
// CacheService.Set tests
// ---------------------------------------------------------------------------

func TestSet_SimpleString(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	if err := svc.Set(ctx, "greeting", "hello", time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// Verify raw data stored is valid JSON.
	raw, err := client.Get(ctx, "greeting").Result()
	if err != nil {
		t.Fatalf("raw Get: %v", err)
	}
	if raw != `"hello"` {
		t.Fatalf("expected JSON string, got %q", raw)
	}
}

func TestSet_Struct(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	ctx := context.Background()
	user := User{Name: "Alice", Email: "alice@example.com"}
	if err := svc.Set(ctx, "user:1", user, 5*time.Minute); err != nil {
		t.Fatalf("Set struct: %v", err)
	}
}

func TestSet_Map(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	m := map[string]int{"a": 1, "b": 2}
	if err := svc.Set(ctx, "counts", m, time.Hour); err != nil {
		t.Fatalf("Set map: %v", err)
	}
}

func TestSet_NilValue(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	if err := svc.Set(ctx, "nil-key", nil, time.Minute); err != nil {
		t.Fatalf("Set nil: %v", err)
	}
}

func TestSet_MarshalError(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	// Channels cannot be marshalled to JSON.
	err := svc.Set(ctx, "bad", make(chan int), time.Minute)
	if err == nil {
		t.Fatal("expected marshal error for channel")
	}
	if !strings.Contains(err.Error(), "marshaling cache value") {
		t.Fatalf("expected 'marshaling cache value' in error, got: %v", err)
	}
}

func TestSet_MarshalError_FuncValue(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	err := svc.Set(ctx, "fn", func() {}, time.Minute)
	if err == nil {
		t.Fatal("expected marshal error for func")
	}
	if !strings.Contains(err.Error(), "marshaling cache value") {
		t.Fatalf("expected 'marshaling cache value', got: %v", err)
	}
}

func TestSet_ZeroTTL(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	// TTL=0 means no expiration in Redis.
	if err := svc.Set(ctx, "forever", "val", 0); err != nil {
		t.Fatalf("Set with zero TTL: %v", err)
	}
}

// ---------------------------------------------------------------------------
// CacheService.Get tests
// ---------------------------------------------------------------------------

func TestGet_String(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "k", "value123", time.Minute)

	var got string
	if err := svc.Get(ctx, "k", &got); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "value123" {
		t.Fatalf("expected 'value123', got %q", got)
	}
}

func TestGet_Struct(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	ctx := context.Background()
	mustSet(ctx, t, svc, "item:42", Item{ID: 42, Name: "Widget"}, time.Minute)

	var got Item
	if err := svc.Get(ctx, "item:42", &got); err != nil {
		t.Fatalf("Get struct: %v", err)
	}
	if got.ID != 42 || got.Name != "Widget" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestGet_Map(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	expected := map[string]float64{"x": 1.5, "y": 2.5}
	mustSet(ctx, t, svc, "coords", expected, time.Minute)

	var got map[string]float64
	if err := svc.Get(ctx, "coords", &got); err != nil {
		t.Fatalf("Get map: %v", err)
	}
	if got["x"] != 1.5 || got["y"] != 2.5 {
		t.Fatalf("unexpected map: %v", got)
	}
}

func TestGet_KeyNotFound(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	var out string
	err := svc.Get(ctx, "nonexistent", &out)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestGet_UnmarshalError(t *testing.T) {
	client, mr := startMiniRedis(t)
	svc := NewCacheService(client)

	// Write raw invalid JSON directly via miniredis.
	if err := mr.Set("bad-json", "not-valid-json{{{"); err != nil {
		t.Fatalf("mr.Set: %v", err)
	}

	ctx := context.Background()
	var out map[string]string
	err := svc.Get(ctx, "bad-json", &out)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
	// The error message should relate to JSON parsing.
	errMsg := err.Error()
	if !strings.Contains(errMsg, "invalid") && !strings.Contains(errMsg, "json") && !strings.Contains(errMsg, "unmarshal") {
		t.Fatalf("expected JSON-related error, got: %v", err)
	}
}

func TestGet_WrongType_Unmarshal(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	// Store a string, try to read as struct.
	mustSet(ctx, t, svc, "str", "hello", time.Minute)

	var out struct{ Name string }
	err := svc.Get(ctx, "str", &out)
	if err == nil {
		t.Fatal("expected unmarshal error when types mismatch")
	}
}

func TestGet_Slice(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	input := []string{"a", "b", "c"}
	mustSet(ctx, t, svc, "list", input, time.Minute)

	var got []string
	if err := svc.Get(ctx, "list", &got); err != nil {
		t.Fatalf("Get slice: %v", err)
	}
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("unexpected slice: %v", got)
	}
}

func TestGet_Integer(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "count", 42, time.Minute)

	var got int
	if err := svc.Get(ctx, "count", &got); err != nil {
		t.Fatalf("Get int: %v", err)
	}
	// JSON numbers deserialize to float64 by default, but into int it should work.
	if got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
}

func TestGet_Boolean(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "flag", true, time.Minute)

	var got bool
	if err := svc.Get(ctx, "flag", &got); err != nil {
		t.Fatalf("Get bool: %v", err)
	}
	if !got {
		t.Fatal("expected true")
	}
}

func TestGet_NilDest(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "k", "val", time.Minute)

	// Passing a non-pointer should cause unmarshal error.
	err := svc.Get(ctx, "k", "not-a-pointer")
	if err == nil {
		t.Fatal("expected error when dest is not a pointer")
	}
}

// ---------------------------------------------------------------------------
// CacheService.Delete tests
// ---------------------------------------------------------------------------

func TestDelete_NoKeys(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	if err := svc.Delete(ctx); err != nil {
		t.Fatalf("Delete with no keys should not error: %v", err)
	}
}

func TestDelete_SingleKey(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "del-me", "value", time.Minute)

	if err := svc.Delete(ctx, "del-me"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	var out string
	if err := svc.Get(ctx, "del-me", &out); err == nil {
		t.Fatal("expected key to be deleted")
	}
}

func TestDelete_MultipleKeys(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "a", 1, time.Minute)
	mustSet(ctx, t, svc, "b", 2, time.Minute)
	mustSet(ctx, t, svc, "c", 3, time.Minute)

	if err := svc.Delete(ctx, "a", "b"); err != nil {
		t.Fatalf("Delete multiple: %v", err)
	}

	// a and b should be gone, c should remain.
	var v int
	if err := svc.Get(ctx, "a", &v); err == nil {
		t.Fatal("key 'a' should have been deleted")
	}
	if err := svc.Get(ctx, "b", &v); err == nil {
		t.Fatal("key 'b' should have been deleted")
	}
	if err := svc.Get(ctx, "c", &v); err != nil {
		t.Fatalf("key 'c' should still exist: %v", err)
	}
	if v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}
}

func TestDelete_NonexistentKey(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	// Deleting a key that doesn't exist should not error.
	if err := svc.Delete(ctx, "ghost"); err != nil {
		t.Fatalf("Delete nonexistent: %v", err)
	}
}

// ---------------------------------------------------------------------------
// CacheService.DeleteByPattern tests
// ---------------------------------------------------------------------------

func TestDeleteByPattern_MatchingKeys(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "user:1", "alice", time.Minute)
	mustSet(ctx, t, svc, "user:2", "bob", time.Minute)
	mustSet(ctx, t, svc, "session:1", "s1", time.Minute)

	if err := svc.DeleteByPattern(ctx, "user:*"); err != nil {
		t.Fatalf("DeleteByPattern: %v", err)
	}

	// user keys should be deleted.
	var out string
	if err := svc.Get(ctx, "user:1", &out); err == nil {
		t.Fatal("user:1 should have been deleted")
	}
	if err := svc.Get(ctx, "user:2", &out); err == nil {
		t.Fatal("user:2 should have been deleted")
	}
	// session key should remain.
	if err := svc.Get(ctx, "session:1", &out); err != nil {
		t.Fatalf("session:1 should remain: %v", err)
	}
}

func TestDeleteByPattern_NoMatches(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "x", 1, time.Minute)

	// Pattern that matches nothing should succeed silently.
	if err := svc.DeleteByPattern(ctx, "nonexistent:*"); err != nil {
		t.Fatalf("DeleteByPattern no match: %v", err)
	}

	// Original key should still be there.
	var v int
	if err := svc.Get(ctx, "x", &v); err != nil {
		t.Fatalf("key 'x' should remain: %v", err)
	}
}

func TestDeleteByPattern_AllKeys(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "item:1", "a", time.Minute)
	mustSet(ctx, t, svc, "item:2", "b", time.Minute)

	if err := svc.DeleteByPattern(ctx, "item:*"); err != nil {
		t.Fatalf("DeleteByPattern all: %v", err)
	}

	var out string
	if err := svc.Get(ctx, "item:1", &out); err == nil {
		t.Fatal("item:1 should be deleted")
	}
	if err := svc.Get(ctx, "item:2", &out); err == nil {
		t.Fatal("item:2 should be deleted")
	}
}

// ---------------------------------------------------------------------------
// Round-trip tests (Set then Get for various types)
// ---------------------------------------------------------------------------

func TestRoundTrip_NestedStruct(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	ctx := context.Background()
	input := Person{
		Name: "Charlie",
		Age:  30,
		Address: Address{
			City:    "Bogota",
			Country: "CO",
		},
	}
	if err := svc.Set(ctx, "person:3", input, time.Minute); err != nil {
		t.Fatalf("Set nested: %v", err)
	}

	var got Person
	if err := svc.Get(ctx, "person:3", &got); err != nil {
		t.Fatalf("Get nested: %v", err)
	}
	if got.Name != "Charlie" || got.Age != 30 || got.Address.City != "Bogota" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestRoundTrip_EmptyString(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	if err := svc.Set(ctx, "empty", "", time.Minute); err != nil {
		t.Fatalf("Set empty string: %v", err)
	}

	var got string
	if err := svc.Get(ctx, "empty", &got); err != nil {
		t.Fatalf("Get empty string: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestRoundTrip_LargeSlice(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	input := make([]int, 1000)
	for i := range input {
		input[i] = i
	}
	if err := svc.Set(ctx, "big", input, time.Minute); err != nil {
		t.Fatalf("Set large slice: %v", err)
	}

	var got []int
	if err := svc.Get(ctx, "big", &got); err != nil {
		t.Fatalf("Get large slice: %v", err)
	}
	if len(got) != 1000 || got[999] != 999 {
		t.Fatalf("unexpected slice length %d or last element %d", len(got), got[999])
	}
}

func TestRoundTrip_NullJSON(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	// nil marshals to "null" in JSON.
	if err := svc.Set(ctx, "null-val", nil, time.Minute); err != nil {
		t.Fatalf("Set nil: %v", err)
	}

	var got *string
	if err := svc.Get(ctx, "null-val", &got); err != nil {
		t.Fatalf("Get null: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// Context cancellation
// ---------------------------------------------------------------------------

func TestGet_CanceledContext(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	var out string
	err := svc.Get(ctx, "k", &out)
	if err == nil {
		t.Fatal("expected error with canceled context")
	}
}

func TestSet_CanceledContext(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := svc.Set(ctx, "k", "v", time.Minute)
	if err == nil {
		t.Fatal("expected error with canceled context")
	}
}

// ---------------------------------------------------------------------------
// CacheService interface compliance
// ---------------------------------------------------------------------------

func TestCacheService_InterfaceCompliance(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)
	if svc == nil {
		t.Fatal("CacheService should not be nil")
	}
	// Exercise a method to confirm the interface is usable.
	ctx := context.Background()
	mustSet(ctx, t, svc, "iface-test", "ok", time.Second)
}

// ---------------------------------------------------------------------------
// Overwrite behavior
// ---------------------------------------------------------------------------

func TestSet_OverwriteExistingKey(t *testing.T) {
	client, _ := startMiniRedis(t)
	svc := NewCacheService(client)

	ctx := context.Background()
	mustSet(ctx, t, svc, "k", "old", time.Minute)
	mustSet(ctx, t, svc, "k", "new", time.Minute)

	var got string
	if err := svc.Get(ctx, "k", &got); err != nil {
		t.Fatalf("Get after overwrite: %v", err)
	}
	if got != "new" {
		t.Fatalf("expected 'new', got %q", got)
	}
}
