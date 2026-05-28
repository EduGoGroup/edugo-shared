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
	// Limitar attempt para evitar overflow (max 30 = ~5.7 años con base 5s)
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 30 {
		attempt = 30
	}
	return c.RetryDelay * time.Duration(1<<uint(attempt)) //nolint:gosec // attempt está validado arriba
}
