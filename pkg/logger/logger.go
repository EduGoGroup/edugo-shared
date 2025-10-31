// Package logger provides structured logging interfaces and implementations
// for the EduGo shared library, including Zap logger support.
package logger

// Logger define la interfaz para logging estructurado
// This interface allows multiple implementations (Zap, Logrus, etc.)
type Logger interface {
	// Debug registra un mensaje de nivel debug
	Debug(msg string, fields ...interface{})

	// Info registra un mensaje de nivel info
	Info(msg string, fields ...interface{})

	// Warn registra un mensaje de nivel warning
	Warn(msg string, fields ...interface{})

	// Error registra un mensaje de nivel error
	Error(msg string, fields ...interface{})

	// Fatal registra un mensaje de nivel fatal y termina la aplicación
	Fatal(msg string, fields ...interface{})

	// With agrega campos contextuales al logger y retorna un nuevo logger
	// Útil para agregar información de contexto que se incluirá en todos los logs
	// Ejemplo: logger.With("user_id", "123").Info("operación exitosa")
	With(fields ...interface{}) Logger

	// Sync sincroniza el buffer del logger (útil antes de terminar la aplicación)
	Sync() error
}

// Fields es un mapa de campos adicionales para logging estructurado
type Fields map[string]interface{}
