package logger

import (
	"github.com/sirupsen/logrus"
)

// logrusLogger es la implementación de Logger usando Logrus
type logrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

// NewLogrusLogger crea un nuevo logger usando Logrus
func NewLogrusLogger(logger *logrus.Logger) Logger {
	return &logrusLogger{
		logger: logger,
		entry:  logrus.NewEntry(logger),
	}
}

// Debug registra un mensaje de nivel debug
func (l *logrusLogger) Debug(msg string, fields ...interface{}) {
	l.entry.WithFields(convertToLogrusFields(fields...)).Debug(msg)
}

// Info registra un mensaje de nivel info
func (l *logrusLogger) Info(msg string, fields ...interface{}) {
	l.entry.WithFields(convertToLogrusFields(fields...)).Info(msg)
}

// Warn registra un mensaje de nivel warning
func (l *logrusLogger) Warn(msg string, fields ...interface{}) {
	l.entry.WithFields(convertToLogrusFields(fields...)).Warn(msg)
}

// Error registra un mensaje de nivel error
func (l *logrusLogger) Error(msg string, fields ...interface{}) {
	l.entry.WithFields(convertToLogrusFields(fields...)).Error(msg)
}

// Fatal registra un mensaje de nivel fatal y termina la aplicación
func (l *logrusLogger) Fatal(msg string, fields ...interface{}) {
	l.entry.WithFields(convertToLogrusFields(fields...)).Fatal(msg)
}

// With agrega campos contextuales al logger
func (l *logrusLogger) With(fields ...interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(convertToLogrusFields(fields...)),
	}
}

// Sync sincroniza el buffer del logger (no-op para logrus)
func (l *logrusLogger) Sync() error {
	return nil
}

// convertToLogrusFields convierte los campos variádicos a logrus.Fields
func convertToLogrusFields(fields ...interface{}) logrus.Fields {
	result := logrus.Fields{}

	// Los campos vienen en pares clave-valor
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		result[key] = fields[i+1]
	}

	return result
}
