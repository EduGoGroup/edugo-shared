package bootstrap

import (
	logger "github.com/EduGoGroup/edugo-shared/logger"
)

// mergeFactories combina las factories base con las factories mock.
//
// Parámetros:
//   - base: Factories base
//   - mocks: Factories mock que sobrescriben las base
//
// Retorna las factories combinadas.
func mergeFactories(base *Factories, mocks *MockFactories) *Factories {
	result := &Factories{}
	if base != nil {
		*result = *base
	}
	if mocks != nil {
		if mocks.Logger != nil {
			result.Logger = mocks.Logger
		}
		if mocks.PostgreSQL != nil {
			result.PostgreSQL = mocks.PostgreSQL
		}
		if mocks.MongoDB != nil {
			result.MongoDB = mocks.MongoDB
		}
		if mocks.RabbitMQ != nil {
			result.RabbitMQ = mocks.RabbitMQ
		}
		if mocks.S3 != nil {
			result.S3 = mocks.S3
		}
	}
	return result
}

// isRequired verifica si un recurso está en la lista de recursos requeridos.
//
// Parámetros:
//   - resource: Nombre del recurso a verificar
//   - opts: Opciones de bootstrap
//
// Retorna true si el recurso es requerido.
func isRequired(resource string, opts *BootstrapOptions) bool {
	for _, r := range opts.RequiredResources {
		if r == resource {
			return true
		}
	}
	return false
}

// logWarning registra un mensaje de advertencia si el logger está disponible.
//
// Parámetros:
//   - log: Logger a usar
//   - msg: Mensaje de advertencia
//   - err: Error asociado
func logWarning(log logger.Logger, msg string, err error) {
	if log != nil {
		log.With(logger.FieldError, err).Warn(msg)
	}
}
