package rabbit

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestGetRetryCount(t *testing.T) {
	tests := []struct {
		name     string
		headers  amqp.Table
		expected int
	}{
		{
			name:     "no retry header",
			headers:  amqp.Table{},
			expected: 0,
		},
		{
			name: "retry count as int32",
			headers: amqp.Table{
				"x-retry-count": int32(3),
			},
			expected: 3,
		},
		{
			name: "retry count as int64",
			headers: amqp.Table{
				"x-retry-count": int64(5),
			},
			expected: 5,
		},
		{
			name: "retry count as int",
			headers: amqp.Table{
				"x-retry-count": int(2),
			},
			expected: 2,
		},
		{
			name: "other headers present",
			headers: amqp.Table{
				"x-retry-count": int32(4),
				"other-header":  "value",
				"timestamp":     int64(1234567890),
			},
			expected: 4,
		},
		{
			name: "nil headers",
			headers: nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := getRetryCount(tt.headers)
			assert.Equal(t, tt.expected, count)
		})
	}
}

func TestCloneHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  amqp.Table
		validate func(t *testing.T, original, cloned amqp.Table)
	}{
		{
			name:    "nil headers",
			headers: nil,
			validate: func(t *testing.T, original, cloned amqp.Table) {
				assert.NotNil(t, cloned)
				assert.Empty(t, cloned)
			},
		},
		{
			name:    "empty headers",
			headers: amqp.Table{},
			validate: func(t *testing.T, original, cloned amqp.Table) {
				assert.NotNil(t, cloned)
				assert.Empty(t, cloned)
			},
		},
		{
			name: "single header",
			headers: amqp.Table{
				"key": "value",
			},
			validate: func(t *testing.T, original, cloned amqp.Table) {
				assert.Equal(t, original["key"], cloned["key"])
				assert.Equal(t, len(original), len(cloned))
			},
		},
		{
			name: "multiple headers",
			headers: amqp.Table{
				"key1":        "value1",
				"key2":        int32(42),
				"key3":        int64(123),
				"retry-count": int32(3),
			},
			validate: func(t *testing.T, original, cloned amqp.Table) {
				assert.Equal(t, len(original), len(cloned))
				for key, value := range original {
					assert.Equal(t, value, cloned[key])
				}
			},
		},
		{
			name: "mutation test - cloned is independent",
			headers: amqp.Table{
				"original-key": "original-value",
			},
			validate: func(t *testing.T, original, cloned amqp.Table) {
				// Modificar el clon
				cloned["new-key"] = "new-value"
				cloned["original-key"] = "modified-value"
				
				// Original no debe cambiar
				assert.NotContains(t, original, "new-key")
				assert.Equal(t, "original-value", original["original-key"])
				
				// Clon debe tener los cambios
				assert.Equal(t, "modified-value", cloned["original-key"])
				assert.Equal(t, "new-value", cloned["new-key"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloned := cloneHeaders(tt.headers)
			tt.validate(t, tt.headers, cloned)
		})
	}
}

func TestCloneHeaders_DeepCopy(t *testing.T) {
	original := amqp.Table{
		"counter": int32(1),
		"name":    "test",
	}

	cloned := cloneHeaders(original)

	// Verificar que es una copia
	assert.Equal(t, original["counter"], cloned["counter"])
	assert.Equal(t, original["name"], cloned["name"])

	// Modificar original
	original["counter"] = int32(99)
	original["new"] = "added"

	// Clon no debe cambiar
	assert.Equal(t, int32(1), cloned["counter"])
	assert.NotContains(t, cloned, "new")
}

func TestDLQConfig_Enabled(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{"enabled", true},
		{"disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DLQConfig{
				Enabled: tt.enabled,
			}

			assert.Equal(t, tt.enabled, config.Enabled)
		})
	}
}

func TestDLQConfig_MaxRetries(t *testing.T) {
	tests := []struct {
		name       string
		maxRetries int
	}{
		{"zero retries", 0},
		{"one retry", 1},
		{"default retries", 3},
		{"many retries", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DLQConfig{
				MaxRetries: tt.maxRetries,
			}

			assert.Equal(t, tt.maxRetries, config.MaxRetries)
		})
	}
}
