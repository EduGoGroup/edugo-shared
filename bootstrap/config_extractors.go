package bootstrap

import (
	"fmt"
	"reflect"
)

// extractConfigField es un helper genérico para extraer configuración de un struct.
//
// Parámetros:
//   - config: Struct de configuración (valor o puntero)
//   - fieldName: Nombre del campo a extraer
//
// Retorna el valor del campo del tipo T o error si no se encuentra.
func extractConfigField[T any](config interface{}, fieldName string) (T, error) {
	var zero T

	// Intentar type assertion directo primero
	if typedConfig, ok := config.(T); ok {
		return typedConfig, nil
	}

	// Usar reflection para extraer el campo
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return zero, fmt.Errorf("config is nil")
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return zero, fmt.Errorf("config must be a struct, got %T", config)
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return zero, fmt.Errorf("%s field not found in config", fieldName)
	}

	// Intentar convertir el campo al tipo deseado
	fieldInterface := field.Interface()
	if typedField, ok := fieldInterface.(T); ok {
		return typedField, nil
	}

	// Si el campo es un puntero, intentar desreferenciarlo
	if field.Kind() == reflect.Ptr && !field.IsNil() {
		if typedField, ok := field.Elem().Interface().(T); ok {
			return typedField, nil
		}
	}

	return zero, fmt.Errorf("%s field is not of expected type, got %T", fieldName, fieldInterface)
}

// extractPostgreSQLConfig extrae configuración de PostgreSQL usando el helper genérico.
//
// Parámetros:
//   - config: Configuración de la aplicación (puede ser struct o puntero)
//
// Retorna la configuración de PostgreSQL o error si no se encuentra.
func extractPostgreSQLConfig(config interface{}) (PostgreSQLConfig, error) {
	return extractConfigField[PostgreSQLConfig](config, "PostgreSQL")
}

// extractMongoDBConfig extrae configuración de MongoDB usando el helper genérico.
//
// Parámetros:
//   - config: Configuración de la aplicación (puede ser struct o puntero)
//
// Retorna la configuración de MongoDB o error si no se encuentra.
func extractMongoDBConfig(config interface{}) (MongoDBConfig, error) {
	return extractConfigField[MongoDBConfig](config, "MongoDB")
}

// extractRabbitMQConfig extrae configuración de RabbitMQ usando el helper genérico.
//
// Parámetros:
//   - config: Configuración de la aplicación (puede ser struct o puntero)
//
// Retorna la configuración de RabbitMQ o error si no se encuentra.
func extractRabbitMQConfig(config interface{}) (RabbitMQConfig, error) {
	return extractConfigField[RabbitMQConfig](config, "RabbitMQ")
}

// extractS3Config extrae configuración de S3 usando el helper genérico.
//
// Parámetros:
//   - config: Configuración de la aplicación (puede ser struct o puntero)
//
// Retorna la configuración de S3 o error si no se encuentra.
func extractS3Config(config interface{}) (S3Config, error) {
	return extractConfigField[S3Config](config, "S3")
}

// extractEnvAndVersion extrae los campos Environment y Version de una configuración.
//
// Busca campos llamados "Environment" y "Version" en el struct proporcionado.
// Si no los encuentra o el config es nil, retorna valores por defecto.
//
// Parámetros:
//   - config: Struct de configuración (puede ser valor o puntero)
//
// Retorna:
//   - environment: Valor del campo Environment o "unknown"
//   - version: Valor del campo Version o "0.0.0"
func extractEnvAndVersion(config interface{}) (string, string) {
	if config == nil {
		return "unknown", "0.0.0"
	}

	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "unknown", "0.0.0"
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "unknown", "0.0.0"
	}

	// Buscar campo Environment
	env := "unknown"
	envField := v.FieldByName("Environment")
	if envField.IsValid() && envField.Kind() == reflect.String {
		env = envField.String()
		if env == "" {
			env = "unknown"
		}
	}

	// Buscar campo Version
	version := "0.0.0"
	versionField := v.FieldByName("Version")
	if versionField.IsValid() && versionField.Kind() == reflect.String {
		ver := versionField.String()
		if ver != "" {
			version = ver
		}
	}

	return env, version
}
