package logger

import (
	"github.com/sirupsen/logrus"
)

// logrusLogger es la implementación privada de la interface Logger usando
// la librería Logrus como backend. Mantiene una referencia al logger base
// y una entry que acumula campos contextuales a través de llamadas a With().
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

// With agrega campos contextuales al logger y retorna un nuevo logger con esos campos.
//
// Los campos se pasan en pares clave-valor:
//
//	logger.With("user_id", 123, "action", "login")
//
// Las claves deben ser strings. Si hay un número impar de argumentos,
// el último se ignora silenciosamente.
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

// convertToLogrusFields convierte los campos variádicos a logrus.Fields.
//
// Espera que los argumentos se pasen en forma de pares clave-valor:
//
//	convertToLogrusFields("clave1", valor1, "clave2", valor2, ...)
//
// Comportamiento:
//   - Las claves deben ser de tipo string. Si una clave no es string,
//     se ignora ese par completo (no se agrega ninguna entrada al resultado).
//   - Si se recibe un número impar de argumentos, el último argumento queda
//     sin pareja y se ignora silenciosamente.
//   - Los valores se almacenan tal cual se reciben, sin transformaciones
//     adicionales.
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
