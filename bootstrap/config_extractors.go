package bootstrap

import (
	"fmt"
	"reflect"
)

// extractPostgreSQLConfig extrae configuración de PostgreSQL usando reflection.
//
// Parámetros:
//   - config: Configuración de la aplicación (puede ser struct o puntero)
//
// Retorna la configuración de PostgreSQL o error si no se encuentra.
func extractPostgreSQLConfig(config interface{}) (PostgreSQLConfig, error) {
	// Intentar type assertion directo primero
	if pgConfig, ok := config.(PostgreSQLConfig); ok {
		return pgConfig, nil
	}

	// Usar reflection para extraer campo PostgreSQL
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return PostgreSQLConfig{}, fmt.Errorf("config must be a struct, got %T", config)
	}

	pgField := v.FieldByName("PostgreSQL")
	if !pgField.IsValid() {
		return PostgreSQLConfig{}, fmt.Errorf("PostgreSQL field not found in config")
	}

	if pgConfig, ok := pgField.Interface().(PostgreSQLConfig); ok {
		return pgConfig, nil
	}

	return PostgreSQLConfig{}, fmt.Errorf("PostgreSQL field is not of type PostgreSQLConfig")
}

// extractMongoDBConfig extrae configuración de MongoDB usando reflection.
//
// Parámetros:
//   - config: Configuración de la aplicación (puede ser struct o puntero)
//
// Retorna la configuración de MongoDB o error si no se encuentra.
func extractMongoDBConfig(config interface{}) (MongoDBConfig, error) {
	// Intentar type assertion directo
	if mongoConfig, ok := config.(MongoDBConfig); ok {
		return mongoConfig, nil
	}

	// Usar reflection
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return MongoDBConfig{}, fmt.Errorf("config must be a struct, got %T", config)
	}

	mongoField := v.FieldByName("MongoDB")
	if !mongoField.IsValid() {
		return MongoDBConfig{}, fmt.Errorf("MongoDB field not found in config")
	}

	if mongoConfig, ok := mongoField.Interface().(MongoDBConfig); ok {
		return mongoConfig, nil
	}

	return MongoDBConfig{}, fmt.Errorf("MongoDB field is not of type MongoDBConfig")
}

// extractRabbitMQConfig extrae configuración de RabbitMQ usando reflection.
//
// Parámetros:
//   - config: Configuración de la aplicación (puede ser struct o puntero)
//
// Retorna la configuración de RabbitMQ o error si no se encuentra.
func extractRabbitMQConfig(config interface{}) (RabbitMQConfig, error) {
	// Intentar type assertion directo
	if rabbitConfig, ok := config.(RabbitMQConfig); ok {
		return rabbitConfig, nil
	}

	// Usar reflection
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return RabbitMQConfig{}, fmt.Errorf("config must be a struct, got %T", config)
	}

	rabbitField := v.FieldByName("RabbitMQ")
	if !rabbitField.IsValid() {
		return RabbitMQConfig{}, fmt.Errorf("RabbitMQ field not found in config")
	}

	if rabbitConfig, ok := rabbitField.Interface().(RabbitMQConfig); ok {
		return rabbitConfig, nil
	}

	return RabbitMQConfig{}, fmt.Errorf("RabbitMQ field is not of type RabbitMQConfig")
}

// extractS3Config extrae configuración de S3 usando reflection.
//
// Parámetros:
//   - config: Configuración de la aplicación (puede ser struct o puntero)
//
// Retorna la configuración de S3 o error si no se encuentra.
func extractS3Config(config interface{}) (S3Config, error) {
	// Intentar type assertion directo
	if s3Config, ok := config.(S3Config); ok {
		return s3Config, nil
	}

	// Usar reflection
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return S3Config{}, fmt.Errorf("config must be a struct, got %T", config)
	}

	s3Field := v.FieldByName("S3")
	if !s3Field.IsValid() {
		return S3Config{}, fmt.Errorf("S3 field not found in config")
	}

	if s3Config, ok := s3Field.Interface().(S3Config); ok {
		return s3Config, nil
	}

	return S3Config{}, fmt.Errorf("S3 field is not of type S3Config")
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
