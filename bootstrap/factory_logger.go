package bootstrap

import (
	"context"
	"os"

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
func (f *DefaultLoggerFactory) CreateLogger(ctx context.Context, env string, version string) (*logrus.Logger, error) {
	logger := logrus.New()

	// Configurar output
	logger.SetOutput(os.Stdout)

	// Configurar formato según entorno
	if env == "production" || env == "prod" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Configurar nivel de log según entorno
	switch env {
	case "production", "prod":
		logger.SetLevel(logrus.InfoLevel)
	case "qa", "staging":
		logger.SetLevel(logrus.InfoLevel)
	case "development", "dev":
		logger.SetLevel(logrus.DebugLevel)
	case "local":
		logger.SetLevel(logrus.TraceLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Agregar fields globales
	logger.WithFields(logrus.Fields{
		"version": version,
		"env":     env,
	})

	return logger, nil
}

// Verificar que DefaultLoggerFactory implementa LoggerFactory
var _ LoggerFactory = (*DefaultLoggerFactory)(nil)
