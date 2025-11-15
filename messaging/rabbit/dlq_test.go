package rabbit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/messaging/rabbit"
)

func TestDLQConfig_CalculateBackoff(t *testing.T) {
	config := rabbit.DLQConfig{
		RetryDelay:            5 * time.Second,
		UseExponentialBackoff: true,
	}

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{0, 5 * time.Second},  // 5s * 2^0 = 5s
		{1, 10 * time.Second}, // 5s * 2^1 = 10s
		{2, 20 * time.Second}, // 5s * 2^2 = 20s
		{3, 40 * time.Second}, // 5s * 2^3 = 40s
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("attempt_%d", tt.attempt), func(t *testing.T) {
			got := config.CalculateBackoff(tt.attempt)
			if got != tt.want {
				t.Errorf("CalculateBackoff(%d) = %v, want %v", tt.attempt, got, tt.want)
			}
		})
	}
}

func TestDLQConfig_LinearBackoff(t *testing.T) {
	config := rabbit.DLQConfig{
		RetryDelay:            5 * time.Second,
		UseExponentialBackoff: false,
	}

	// Sin exponential backoff, siempre retorna RetryDelay
	for attempt := 0; attempt < 5; attempt++ {
		got := config.CalculateBackoff(attempt)
		if got != 5*time.Second {
			t.Errorf("CalculateBackoff(%d) = %v, want 5s", attempt, got)
		}
	}
}

func TestDefaultDLQConfig(t *testing.T) {
	config := rabbit.DefaultDLQConfig()

	tests := []struct {
		name     string
		got      interface{}
		want     interface{}
		errField string
	}{
		{"Enabled", config.Enabled, true, "Enabled should be true"},
		{"MaxRetries", config.MaxRetries, 3, "MaxRetries should be 3"},
		{"RetryDelay", config.RetryDelay, 5 * time.Second, "RetryDelay should be 5s"},
		{"DLXExchange", config.DLXExchange, "dlx", "DLXExchange should be 'dlx'"},
		{"DLXRoutingKey", config.DLXRoutingKey, "dlq", "DLXRoutingKey should be 'dlq'"},
		{"UseExponentialBackoff", config.UseExponentialBackoff, true, "UseExponentialBackoff should be true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s: got %v, want %v", tt.errField, tt.got, tt.want)
			}
		})
	}
}
