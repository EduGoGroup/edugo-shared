package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSlogProvider_JSONFormat(t *testing.T) {
	cfg := SlogConfig{
		Level:   "info",
		Format:  "json",
		Service: "test-service",
		Env:     "test",
		Version: "1.0.0",
	}

	// Llamar a NewSlogProvider directamente y verificar que retorna non-nil
	l := NewSlogProvider(cfg)
	require.NotNil(t, l)

	// Verificar que el logger funciona sin panic en todos los niveles
	assert.NotPanics(t, func() {
		l.Info("test message", "key", "value")
		l.Debug("debug message")
		l.Warn("warn message")
		l.Error("error message")
	})

	// Verificar que el logger respeta el nivel configurado:
	// Con nivel "info", Debug no debe estar habilitado
	assert.False(t, l.Handler().Enabled(context.Background(), slog.LevelDebug),
		"con nivel info, debug no debe estar habilitado")
	assert.True(t, l.Handler().Enabled(context.Background(), slog.LevelInfo),
		"con nivel info, info debe estar habilitado")
	assert.True(t, l.Handler().Enabled(context.Background(), slog.LevelWarn),
		"con nivel info, warn debe estar habilitado")
	assert.True(t, l.Handler().Enabled(context.Background(), slog.LevelError),
		"con nivel info, error debe estar habilitado")
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
