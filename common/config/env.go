// Package config provides environment variable loading and configuration
// utilities for the EduGo shared library.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// GetEnv obtiene una variable de entorno o retorna un valor por defecto
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvRequired obtiene una variable de entorno requerida, panic si no existe
func GetEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

// GetEnvInt obtiene una variable de entorno como int o retorna un valor por defecto
func GetEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetEnvBool obtiene una variable de entorno como bool o retorna un valor por defecto
func GetEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetEnvDuration obtiene una variable de entorno como time.Duration o retorna un valor por defecto
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// MustGetEnv obtiene una variable de entorno requerida, panic si no existe o está vacía
func MustGetEnv(key string) string {
	return GetEnvRequired(key)
}

// SetEnv establece una variable de entorno
func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// UnsetEnv elimina una variable de entorno
func UnsetEnv(key string) error {
	return os.Unsetenv(key)
}

// LookupEnv verifica si una variable de entorno existe
func LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// GetEnvironment retorna el ambiente actual (development, staging, production)
func GetEnvironment() string {
	return GetEnv("APP_ENV", "development")
}

// IsDevelopment verifica si estamos en ambiente de desarrollo
func IsDevelopment() bool {
	env := GetEnvironment()
	return env == "development" || env == "dev" || env == "local"
}

// IsProduction verifica si estamos en ambiente de producción
func IsProduction() bool {
	env := GetEnvironment()
	return env == "production" || env == "prod"
}

// IsStaging verifica si estamos en ambiente de staging
func IsStaging() bool {
	env := GetEnvironment()
	return env == "staging" || env == "stage" || env == "qa"
}
