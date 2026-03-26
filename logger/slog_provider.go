package logger

import (
	"log/slog"
	"os"
)

// SlogConfig holds configuration for creating a slog.Logger.
type SlogConfig struct {
	// Level is the minimum log level: "debug", "info", "warn", "error". Default: "info".
	Level string
	// Format is the output format: "json" or "text". Default: "json".
	Format string
	// Service is the service name added to every log entry.
	Service string
	// Env is the environment name (dev, staging, production).
	Env string
	// Version is the application version.
	Version string
}

// NewSlogProvider creates a *slog.Logger using stdlib slog handlers.
// JSON format for production (structured, Datadog-compatible), text for development.
//
// The returned logger includes service, env, and version as base fields.
// For backwards compatibility with the Logger interface, wrap with NewSlogAdapter.
func NewSlogProvider(cfg SlogConfig) *slog.Logger {
	level := parseSlogLevel(cfg.Level)
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handler slog.Handler
	if cfg.Format == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	l := slog.New(handler)

	// Add base fields if provided
	if cfg.Service != "" {
		l = l.With(slog.String(FieldService, cfg.Service))
	}
	if cfg.Env != "" {
		l = l.With(slog.String(FieldEnvironment, cfg.Env))
	}
	if cfg.Version != "" {
		l = l.With(slog.String(FieldVersion, cfg.Version))
	}

	return l
}

func parseSlogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
