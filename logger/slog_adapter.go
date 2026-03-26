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
// Usar junto con NewSlogProvider:
//
//	slogLogger := logger.NewSlogProvider(cfg)
//	appLogger := logger.NewSlogAdapter(slogLogger)
func NewSlogAdapter(l *slog.Logger) Logger {
	return &SlogAdapter{logger: l}
}

// SlogLogger retorna el *slog.Logger subyacente para uso directo.
// Útil cuando necesitas pasar el slog.Logger a middleware o contexto.
func (a *SlogAdapter) SlogLogger() *slog.Logger {
	return a.logger
}

func (a *SlogAdapter) Debug(msg string, fields ...interface{}) {
	a.logger.Debug(msg, fields...)
}

func (a *SlogAdapter) Info(msg string, fields ...interface{}) {
	a.logger.Info(msg, fields...)
}

func (a *SlogAdapter) Warn(msg string, fields ...interface{}) {
	a.logger.Warn(msg, fields...)
}

func (a *SlogAdapter) Error(msg string, fields ...interface{}) {
	a.logger.Error(msg, fields...)
}

func (a *SlogAdapter) Fatal(msg string, fields ...interface{}) {
	a.logger.Error(msg, fields...)
	os.Exit(1)
}

func (a *SlogAdapter) With(fields ...interface{}) Logger {
	return &SlogAdapter{logger: a.logger.With(fields...)}
}

func (a *SlogAdapter) Sync() error {
	return nil
}
