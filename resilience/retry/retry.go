package retry

import (
	"context"
	"errors"
	"time"
)

// ErrorType clasifica el tipo de error para decidir si reintentar.
type ErrorType int

const (
	// ErrorTypePermanent indica un error que no se puede resolver con reintentos.
	ErrorTypePermanent ErrorType = iota
	// ErrorTypeTransient indica un error temporal que puede resolverse con reintentos.
	ErrorTypeTransient
)

// ErrorClassifier determina si un error es transitorio o permanente.
// Si retorna ErrorTypePermanent, WithRetry no reintenta.
// Si es nil, todos los errores se tratan como transitorios.
type ErrorClassifier func(error) ErrorType

// Logger interfaz simple para logging de reintentos.
type Logger interface {
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

// Config configura el comportamiento de retry.
type Config struct {
	MaxRetries      int
	InitialBackoff  time.Duration
	MaxBackoff      time.Duration
	BackoffMultiple float64
	Logger          Logger
	Classifier      ErrorClassifier
}

// DefaultConfig retorna la configuracion por defecto.
func DefaultConfig() Config {
	return Config{
		MaxRetries:      3,
		InitialBackoff:  500 * time.Millisecond,
		MaxBackoff:      10 * time.Second,
		BackoffMultiple: 2.0,
	}
}

// WithRetry ejecuta una operacion con logica de reintento y backoff exponencial.
// Si Classifier es nil, todos los errores se tratan como transitorios.
func WithRetry(ctx context.Context, cfg Config, operation func() error) error {
	var lastErr error
	backoff := cfg.InitialBackoff

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			logWarn(cfg.Logger, "operacion cancelada por contexto",
				"attempt", attempt,
				"lastError", lastErr,
			)
			return ctx.Err()
		default:
		}

		err := operation()
		if err == nil {
			if attempt > 0 {
				logInfo(cfg.Logger, "operacion exitosa despues de reintentos", "attempts", attempt+1)
			}
			return nil
		}

		lastErr = err

		if cfg.Classifier != nil && cfg.Classifier(err) == ErrorTypePermanent {
			logWarn(cfg.Logger, "error permanente detectado, no se reintentara",
				"error", err,
				"attempt", attempt+1,
			)
			return err
		}

		if attempt == cfg.MaxRetries {
			logError(cfg.Logger, "maximo de reintentos alcanzado",
				"error", err,
				"attempts", attempt+1,
				"maxRetries", cfg.MaxRetries,
			)
			return err
		}

		logWarn(cfg.Logger, "error transitorio, reintentando",
			"error", err,
			"attempt", attempt+1,
			"backoff", backoff,
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}

		backoff = time.Duration(float64(backoff) * cfg.BackoffMultiple)
		if backoff > cfg.MaxBackoff {
			backoff = cfg.MaxBackoff
		}
	}

	return lastErr
}

// IsContextError verifica si un error es por cancelacion de contexto.
func IsContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func logInfo(l Logger, msg string, kv ...any) {
	if l != nil {
		l.Info(msg, kv...)
	}
}

func logWarn(l Logger, msg string, kv ...any) {
	if l != nil {
		l.Warn(msg, kv...)
	}
}

func logError(l Logger, msg string, kv ...any) {
	if l != nil {
		l.Error(msg, kv...)
	}
}
