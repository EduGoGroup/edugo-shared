package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNextBackoff(t *testing.T) {
	tests := []struct {
		name     string
		current  time.Duration
		maxDelay time.Duration
		want     time.Duration
	}{
		{
			name:     "doubles delay",
			current:  1 * time.Second,
			maxDelay: 60 * time.Second,
			want:     2 * time.Second,
		},
		{
			name:     "caps at maxDelay",
			current:  40 * time.Second,
			maxDelay: 60 * time.Second,
			want:     60 * time.Second,
		},
		{
			name:     "already at max",
			current:  60 * time.Second,
			maxDelay: 60 * time.Second,
			want:     60 * time.Second,
		},
		{
			name:     "small values",
			current:  100 * time.Millisecond,
			maxDelay: 5 * time.Second,
			want:     200 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextBackoff(tt.current, tt.maxDelay)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExponentialDelay(t *testing.T) {
	base := 5 * time.Second

	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{"attempt 0", 0, 5 * time.Second},
		{"attempt 1", 1, 10 * time.Second},
		{"attempt 2", 2, 20 * time.Second},
		{"attempt 3", 3, 40 * time.Second},
		{"negative attempt", -1, 5 * time.Second},
		{"attempt capped at 30", 50, base * time.Duration(1<<30)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExponentialDelay(base, tt.attempt)
			assert.Equal(t, tt.want, got)
		})
	}
}
