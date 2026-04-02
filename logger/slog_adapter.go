package logger

import (
	"log/slog"
	"os"
)

// SlogAdapter implementa la interfaz Logger delegando a un *slog.Logger.
// Esto permite una migración gradual: el código existente usa la interfaz Logger
// mientras que el código nuevo puede usar *slog.Logger directamente.
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlogAdapter envuelve un *slog.Logger para satisfacer la interfaz Logger.
// Si se pasa nil, se usa slog.Default() como fallback para evitar panics.
//
//	slogLogger := logger.NewSlogProvider(cfg)
//	appLogger := logger.NewSlogAdapter(slogLogger)
func NewSlogAdapter(l *slog.Logger) Logger {
	if l == nil {
		l = slog.Default()
	}
	return &SlogAdapter{logger: l}
}

// SlogLogger retorna el *slog.Logger subyacente para uso directo.
// Útil cuando necesitas pasar el slog.Logger a middleware o contexto.
func (a *SlogAdapter) SlogLogger() *slog.Logger {
	if a == nil || a.logger == nil {
		return slog.Default()
	}
	return a.logger
}

// Debug registra un mensaje de nivel debug.
func (a *SlogAdapter) Debug(msg string, fields ...any) {
	a.logger.Debug(msg, fields...)
}

// Info registra un mensaje de nivel info.
func (a *SlogAdapter) Info(msg string, fields ...any) {
	a.logger.Info(msg, fields...)
}

// Warn registra un mensaje de nivel warning.
func (a *SlogAdapter) Warn(msg string, fields ...any) {
	a.logger.Warn(msg, fields...)
}

// Error registra un mensaje de nivel error.
func (a *SlogAdapter) Error(msg string, fields ...any) {
	a.logger.Error(msg, fields...)
}

// Fatal registra un mensaje de nivel fatal y termina la aplicación.
func (a *SlogAdapter) Fatal(msg string, fields ...any) {
	a.logger.Error(msg, fields...)
	os.Exit(1)
}

// With agrega campos contextuales y retorna un nuevo Logger inmutable.
func (a *SlogAdapter) With(fields ...any) Logger {
	return &SlogAdapter{logger: a.logger.With(fields...)}
}

// Sync es un no-op para compatibilidad con la interfaz Logger.
func (a *SlogAdapter) Sync() error {
	return nil
}
