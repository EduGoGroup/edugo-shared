package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
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
	// OtelLevel es el nivel mínimo para el exporter OTel/Loki.
	// Se resuelve independiente de Level (stdout) para evitar inundar
	// backends de observabilidad cuando se sube DEBUG en local.
	// Si está vacío, el call-site decide vía ResolveOtelLevel.
	OtelLevel string
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

// NewSlogProviderFromEnv crea un *slog.Logger leyendo la configuración de variables de entorno:
//   - LOGGING_LEVEL: debug, info, warn, error (default: info)
//   - LOGGING_FORMAT: json, text (default: json)
//   - SERVICE_NAME: nombre del servicio
//   - APP_ENV: entorno (dev, staging, production)
//   - APP_VERSION: versión de la aplicación
//   - OTEL_LOG_LEVEL: nivel mínimo para el exporter OTel/Loki (opcional, ver ResolveOtelLevel).
func NewSlogProviderFromEnv() *slog.Logger {
	return NewSlogProvider(SlogConfig{
		Level:     strings.ToLower(getEnvOrDefault("LOGGING_LEVEL", "info")),
		Format:    strings.ToLower(getEnvOrDefault("LOGGING_FORMAT", "json")),
		Service:   getEnvOrDefault("SERVICE_NAME", ""),
		Env:       getEnvOrDefault("APP_ENV", ""),
		Version:   getEnvOrDefault("APP_VERSION", ""),
		OtelLevel: strings.ToLower(getEnvOrDefault("OTEL_LOG_LEVEL", "")),
	})
}

// ResolveOtelLevel resuelve el nivel mínimo del exporter OTel/Loki según la
// política DA-MPH-5: el level del exporter está desacoplado del level del
// slog raíz (stdout) para que un dev pueda subir DEBUG en local sin inundar
// Loki.
//
// Cascada de resolución:
//  1. Si envVarValue (OTEL_LOG_LEVEL) es válido → se usa.
//  2. Si está vacío → fallback estricto a info para cualquier deploymentEnv
//     (local/staging/prod). Para volver a ver DEBUG en Loki, el dev debe
//     setear OTEL_LOG_LEVEL=debug explícitamente.
//  3. Si envVarValue es inválido → log warning a stderr y se usa el fallback.
//
// La comparación es case-insensitive: "DEBUG", "Debug", "debug" valen igual.
func ResolveOtelLevel(envVarValue, deploymentEnv string) slog.Level {
	v := strings.ToLower(strings.TrimSpace(envVarValue))
	if v == "" {
		return resolveOtelFallback(deploymentEnv)
	}
	switch v {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		fmt.Fprintf(os.Stderr,
			"logger: OTEL_LOG_LEVEL=%q inválido (esperado debug|info|warn|error); "+
				"usando fallback estricto info (DA-MPH-5). APP_ENV=%q (ignorado).\n",
			envVarValue, deploymentEnv)
		return resolveOtelFallback(deploymentEnv)
	}
}

func resolveOtelFallback(deploymentEnv string) slog.Level {
	// DA-MPH-5: fallback estricto info para cualquier deploymentEnv.
	// El parámetro se mantiene en la firma para extensibilidad futura
	// (ej. si en cloud se decide forzar warn como suelo).
	_ = deploymentEnv
	return slog.LevelInfo
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
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
