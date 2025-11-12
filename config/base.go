package config

import (
	"fmt"
	"time"
)

// BaseConfig contiene configuración común a todos los servicios
type BaseConfig struct {
	Environment string              `mapstructure:"environment" validate:"required,oneof=local dev qa prod"`
	ServiceName string              `mapstructure:"service_name" validate:"required"`
	Server      ServerConfig        `mapstructure:"server" validate:"required"`
	Database    DatabaseConfig      `mapstructure:"database" validate:"required"`
	MongoDB     MongoDBConfig       `mapstructure:"mongodb" validate:"required"`
	Logger      LoggerConfig        `mapstructure:"logger" validate:"required"`
	Bootstrap   BootstrapConfig     `mapstructure:"bootstrap"`
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Port         int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" validate:"required"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" validate:"required"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" validate:"required"`
	Host         string        `mapstructure:"host"`
}

// DatabaseConfig configuración de PostgreSQL
type DatabaseConfig struct {
	Host            string        `mapstructure:"host" validate:"required"`
	Port            int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	User            string        `mapstructure:"user" validate:"required"`
	Password        string        `mapstructure:"password" validate:"required"`
	Database        string        `mapstructure:"database" validate:"required"`
	SSLMode         string        `mapstructure:"ssl_mode" validate:"required,oneof=disable require verify-ca verify-full"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" validate:"min=1"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" validate:"min=1"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

// MongoDBConfig configuración de MongoDB
type MongoDBConfig struct {
	URI            string        `mapstructure:"uri" validate:"required"`
	Database       string        `mapstructure:"database" validate:"required"`
	MaxPoolSize    uint64        `mapstructure:"max_pool_size"`
	MinPoolSize    uint64        `mapstructure:"min_pool_size"`
	ConnectTimeout time.Duration `mapstructure:"connect_timeout"`
}

// LoggerConfig configuración del logger
type LoggerConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Format string `mapstructure:"format" validate:"required,oneof=json console"`
}

// BootstrapConfig configuración de recursos opcionales
type BootstrapConfig struct {
	OptionalResources OptionalResourcesConfig `mapstructure:"optional_resources"`
}

// OptionalResourcesConfig recursos que pueden estar deshabilitados
type OptionalResourcesConfig struct {
	RabbitMQ bool `mapstructure:"rabbitmq"`
	S3       bool `mapstructure:"s3"`
}

// ConnectionString construye el DSN para PostgreSQL
func (d *DatabaseConfig) ConnectionString() string {
	return d.ConnectionStringWithDB(d.Database)
}

// ConnectionStringWithDB construye el DSN para PostgreSQL con database específica
func (d *DatabaseConfig) ConnectionStringWithDB(database string) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, database, d.SSLMode)
}
