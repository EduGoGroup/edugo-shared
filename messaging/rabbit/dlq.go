package rabbit

import (
	"time"
)

// DLQConfig configura el Dead Letter Queue
type DLQConfig struct {
	Enabled               bool
	MaxRetries            int           // Default: 3
	RetryDelay            time.Duration // Default: 5s
	DLXExchange           string        // Dead Letter Exchange
	DLXRoutingKey         string        // Routing key para DLQ
	UseExponentialBackoff bool          // Default: true
}

// DefaultDLQConfig retorna configuración por defecto
func DefaultDLQConfig() DLQConfig {
	return DLQConfig{
		Enabled:               true,
		MaxRetries:            3,
		RetryDelay:            5 * time.Second,
		DLXExchange:           "dlx",
		DLXRoutingKey:         "dlq",
		UseExponentialBackoff: true,
	}
}

// CalculateBackoff calcula el delay con exponential backoff
func (c *DLQConfig) CalculateBackoff(attempt int) time.Duration {
	if !c.UseExponentialBackoff {
		return c.RetryDelay
	}
	// Exponential: 5s, 10s, 20s, 40s...
	return c.RetryDelay * time.Duration(1<<uint(attempt))
}
