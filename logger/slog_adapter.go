package logger

import (
	"log/slog"
	"os"
)

// SlogAdapter implements the Logger interface by delegating to a *slog.Logger.
// This allows gradual migration: existing code uses Logger interface while
// new code can use *slog.Logger directly.
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlogAdapter wraps a *slog.Logger to satisfy the Logger interface.
// Use together with NewSlogProvider:
//
//	slogLogger := logger.NewSlogProvider(cfg)
//	appLogger := logger.NewSlogAdapter(slogLogger)
func NewSlogAdapter(l *slog.Logger) Logger {
	return &SlogAdapter{logger: l}
}

// SlogLogger returns the underlying *slog.Logger for direct use.
// Useful when you need to pass the slog.Logger to middleware or context.
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
