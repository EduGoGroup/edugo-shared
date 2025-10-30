package postgres

import "time"

// Config contiene la configuración para conectarse a PostgreSQL
type Config struct {
	// Host del servidor PostgreSQL
	Host string

	// Port del servidor PostgreSQL (por defecto 5432)
	Port int

	// User para autenticación
	User string

	// Password para autenticación
	Password string

	// Database nombre de la base de datos
	Database string

	// MaxConnections número máximo de conexiones en el pool
	MaxConnections int

	// MaxIdleConnections número máximo de conexiones idle
	MaxIdleConnections int

	// MaxLifetime tiempo máximo de vida de una conexión
	MaxLifetime time.Duration

	// SSLMode modo SSL: disable, require, verify-ca, verify-full
	SSLMode string

	// ConnectTimeout timeout para establecer conexión
	ConnectTimeout time.Duration
}

// DefaultConfig retorna una configuración con valores por defecto
func DefaultConfig() Config {
	return Config{
		Host:               "localhost",
		Port:               5432,
		User:               "postgres",
		Database:           "postgres",
		MaxConnections:     25,
		MaxIdleConnections: 5,
		MaxLifetime:        5 * time.Minute,
		SSLMode:            "disable",
		ConnectTimeout:     10 * time.Second,
	}
}
