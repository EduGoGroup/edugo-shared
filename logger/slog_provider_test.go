package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSlogProvider_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	cfg := SlogConfig{
		Level:   "info",
		Format:  "json",
		Service: "test-service",
		Env:     "test",
		Version: "1.0.0",
	}

	// Create logger writing to buffer for testing
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	l := slog.New(handler).With(
		slog.String(FieldService, cfg.Service),
		slog.String(FieldEnvironment, cfg.Env),
		slog.String(FieldVersion, cfg.Version),
	)

	l.Info("test message", "key", "value")

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "test message", entry["msg"])
	assert.Equal(t, "INFO", entry["level"])
	assert.Equal(t, "test-service", entry[FieldService])
	assert.Equal(t, "test", entry[FieldEnvironment])
	assert.Equal(t, "1.0.0", entry[FieldVersion])
	assert.Equal(t, "value", entry["key"])
}

func TestNewSlogProvider_Defaults(t *testing.T) {
	l := NewSlogProvider(SlogConfig{})
	assert.NotNil(t, l)
}

func TestNewSlogProvider_AllLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "unknown"}
	for _, level := range levels {
		l := NewSlogProvider(SlogConfig{Level: level})
		assert.NotNil(t, l, "level: %s", level)
	}
}

func TestNewSlogProvider_TextFormat(t *testing.T) {
	l := NewSlogProvider(SlogConfig{Format: "text"})
	assert.NotNil(t, l)
}

func TestParseSlogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"unknown", slog.LevelInfo},
		{"", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, parseSlogLevel(tt.input))
		})
	}
}
