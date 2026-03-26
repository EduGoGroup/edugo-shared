package logger

import (
	"log/slog"
	"os"
)

// SlogConfig contiene la configuración para crear un slog.Logger.
type SlogConfig struct {
	// Level es el nivel mínimo de log: "debug", "info", "warn", "error". Por defecto: "info".
	Level string
	// Format es el formato de salida: "json" o "text". Por defecto: "json".
	Format string
	// Service es el nombre del servicio agregado a cada entrada de log.
	Service string
	// Env es el nombre del entorno (dev, staging, production).
	Env string
	// Version es la versión de la aplicación.
	Version string
}

// NewSlogProvider crea un *slog.Logger usando handlers estándar de slog.
// Formato JSON para producción (estructurado, compatible con Datadog), texto para desarrollo.
//
// El logger retornado incluye service, env y version como campos base.
// Para compatibilidad con la interfaz Logger, envuelve con NewSlogAdapter.
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

	// Agregar campos base si fueron proporcionados
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
