package redis

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	// URL is the full Redis connection URL.
	// Supports rediss:// for TLS (e.g., Upstash).
	URL string
}
