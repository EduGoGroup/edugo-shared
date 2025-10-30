package mongodb

import "time"

// Config contiene la configuración para conectarse a MongoDB
type Config struct {
	// URI de conexión a MongoDB
	// Formato: mongodb://[username:password@]host[:port][/[defaultauthdb][?options]]
	URI string

	// Database nombre de la base de datos a usar
	Database string

	// Timeout para operaciones (connect, read, write)
	Timeout time.Duration

	// MaxPoolSize número máximo de conexiones en el pool
	MaxPoolSize uint64

	// MinPoolSize número mínimo de conexiones en el pool
	MinPoolSize uint64
}

// DefaultConfig retorna una configuración con valores por defecto
func DefaultConfig() Config {
	return Config{
		URI:         "mongodb://localhost:27017",
		Database:    "test",
		Timeout:     10 * time.Second,
		MaxPoolSize: 100,
		MinPoolSize: 10,
	}
}
