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

func TestResolveOtelLevel(t *testing.T) {
	tests := []struct {
		name        string
		envVar      string
		appEnv      string
		expected    slog.Level
		description string
	}{
		// 1) OTEL_LOG_LEVEL explícito gana siempre, en cualquier entorno.
		{"explicit_debug_local", "debug", "local", slog.LevelDebug, "DEBUG explícito en local debe pasar"},
		{"explicit_debug_staging", "debug", "staging", slog.LevelDebug, "DEBUG explícito en staging debe pasar"},
		{"explicit_debug_prod", "debug", "prod", slog.LevelDebug, "DEBUG explícito en prod debe pasar"},
		{"explicit_info", "info", "local", slog.LevelInfo, ""},
		{"explicit_warn", "warn", "prod", slog.LevelWarn, ""},
		{"explicit_error", "error", "staging", slog.LevelError, ""},

		// 2) Vacío → fallback estricto info para todo entorno (DA-MPH-5).
		{"empty_local_strict", "", "local", slog.LevelInfo, "local sin var debe filtrar a info (estricta)"},
		{"empty_staging", "", "staging", slog.LevelInfo, ""},
		{"empty_prod", "", "prod", slog.LevelInfo, ""},
		{"empty_env_empty", "", "", slog.LevelInfo, ""},

		// 3) Inválido → fallback a APP_ENV (info), sin panic.
		{"typo_warning", "warning", "local", slog.LevelInfo, "typo debe caer al fallback sin panic"},
		{"typo_random", "verbose", "prod", slog.LevelInfo, ""},

		// 4) Case-insensitive: mayúsculas y espacios.
		{"uppercase_debug", "DEBUG", "local", slog.LevelDebug, ""},
		{"mixed_case_warn", "Warn", "prod", slog.LevelWarn, ""},
		{"padded_info", "  info  ", "staging", slog.LevelInfo, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveOtelLevel(tt.envVar, tt.appEnv)
			assert.Equal(t, tt.expected, got, tt.description)
		})
	}
}
