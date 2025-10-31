// Package postgres provides PostgreSQL connection utilities, transaction management,
// and configuration for the EduGo shared library.
package postgres

import "time"

const (
	// DefaultPort is the default PostgreSQL port
	DefaultPort = 5432
	// DefaultMaxConnections is the default maximum number of connections
	DefaultMaxConnections = 25
	// DefaultMaxIdleConnections is the default maximum number of idle connections
	DefaultMaxIdleConnections = 5
	// DefaultMaxLifetime is the default maximum lifetime for connections
	DefaultMaxLifetime = 5 * time.Minute
	// DefaultConnectTimeout is the default connection timeout
	DefaultConnectTimeout = 10 * time.Second
)

// Config contiene la configuración para conectarse a PostgreSQL
type Config struct {
	// Host del servidor PostgreSQL
	Host string

	// User para autenticación
	User string

	// Password para autenticación
	Password string

	// Database nombre de la base de datos
	Database string

	// SSLMode modo SSL: disable, require, verify-ca, verify-full
	SSLMode string

	// MaxLifetime tiempo máximo de vida de una conexión
	MaxLifetime time.Duration

	// ConnectTimeout timeout para establecer conexión
	ConnectTimeout time.Duration

	// Port del servidor PostgreSQL (por defecto 5432)
	Port int

	// MaxConnections número máximo de conexiones en el pool
	MaxConnections int

	// MaxIdleConnections número máximo de conexiones idle
	MaxIdleConnections int
}

// DefaultConfig retorna una configuración con valores por defecto
func DefaultConfig() Config {
	return Config{
		Host:               "localhost",
		Port:               DefaultPort,
		User:               "postgres",
		Database:           "postgres",
		MaxConnections:     DefaultMaxConnections,
		MaxIdleConnections: DefaultMaxIdleConnections,
		MaxLifetime:        DefaultMaxLifetime,
		SSLMode:            "disable",
		ConnectTimeout:     DefaultConnectTimeout,
	}
}
