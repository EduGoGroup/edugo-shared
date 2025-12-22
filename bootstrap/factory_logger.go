package bootstrap

import (
	"context"
	"os"

	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/sirupsen/logrus"
)

// =============================================================================
// LOGGER FACTORY IMPLEMENTATION
// =============================================================================

// DefaultLoggerFactory implementa LoggerFactory usando logrus
type DefaultLoggerFactory struct{}

// NewDefaultLoggerFactory crea una nueva instancia de DefaultLoggerFactory
func NewDefaultLoggerFactory() *DefaultLoggerFactory {
	return &DefaultLoggerFactory{}
}

// CreateLogger crea un logger configurado según el entorno
func (f *DefaultLoggerFactory) CreateLogger(ctx context.Context, env string, version string) (logger.Logger, error) {
	logrusLogger := logrus.New()

	// Configurar output
	logrusLogger.SetOutput(os.Stdout)

	// Configurar formato según entorno
	if env == "production" || env == "prod" {
		logrusLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		logrusLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Configurar nivel de log según entorno
	switch env {
	case "production", "prod":
		logrusLogger.SetLevel(logrus.InfoLevel)
	case "qa", "staging":
		logrusLogger.SetLevel(logrus.InfoLevel)
	case "development", "dev":
		logrusLogger.SetLevel(logrus.DebugLevel)
	case "local":
		logrusLogger.SetLevel(logrus.TraceLevel)
	default:
		logrusLogger.SetLevel(logrus.InfoLevel)
	}

	// Crear wrapper con fields globales
	wrappedLogger := logger.NewLogrusLogger(logrusLogger).With(
		"version", version,
		"env", env,
	)

	return wrappedLogger, nil
}

// Verificar que DefaultLoggerFactory implementa LoggerFactory
var _ LoggerFactory = (*DefaultLoggerFactory)(nil)
