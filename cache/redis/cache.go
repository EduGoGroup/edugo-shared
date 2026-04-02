package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// CacheService provides a generic cache interface with JSON serialization.
type CacheService interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

type redisCacheService struct {
	client *goredis.Client
}

// NewCacheService creates a new Redis-backed CacheService.
func NewCacheService(client *goredis.Client) CacheService {
	return &redisCacheService{client: client}
}

func (s *redisCacheService) Get(ctx context.Context, key string, dest any) error {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (s *redisCacheService) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshaling cache value: %w", err)
	}
	return s.client.Set(ctx, key, data, ttl).Err()
}

func (s *redisCacheService) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return s.client.Del(ctx, keys...).Err()
}

func (s *redisCacheService) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := s.client.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("scanning keys: %w", err)
	}
	if len(keys) > 0 {
		return s.client.Del(ctx, keys...).Err()
	}
	return nil
}
