package bootstrap

import "time"

// =============================================================================
// CONFIGURATION STRUCTS
// =============================================================================
//
// Structs de configuracion para cada tecnologia de infraestructura.
// Usados tanto por el modulo raiz como por los sub-modulos (bootstrap/postgres, etc.).
// Solo dependen de la stdlib — sin imports de librerias externas.
// =============================================================================

// PostgreSQLConfig define la configuracion para conexion PostgreSQL.
type PostgreSQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string //nolint:gosec // G117: Password is a required database config field, not a hardcoded secret
	Database string
	SSLMode  string

	// SearchPath configura el search_path de PostgreSQL.
	// Ejemplo: "auth,iam,academic,ui_config,public"
	// Si esta vacio, se usa "public" por defecto.
	SearchPath string

	// Pool configuration
	MaxOpenConns    int           // default: 25
	MaxIdleConns    int           // default: 5
	ConnMaxLifetime time.Duration // default: 1h
	ConnMaxIdleTime time.Duration // default: 10m
}

// MongoDBConfig define la configuracion para conexion MongoDB.
type MongoDBConfig struct {
	URI      string
	Database string
}

// RabbitMQConfig define la configuracion para conexion RabbitMQ.
type RabbitMQConfig struct {
	URL string
}

// S3Config define la configuracion para cliente AWS S3.
type S3Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string // Para LocalStack o endpoints custom
	ForcePathStyle  bool   // Para LocalStack compatibility
}
