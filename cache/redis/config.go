package redis

// RedisConfig holds Redis connection configuration.
//
//nolint:revive // Keep exported name for backward compatibility in this shared package API.
type RedisConfig struct {
	// URL is the full Redis connection URL.
	// Supports rediss:// for TLS (e.g., Upstash).
	URL string
}
