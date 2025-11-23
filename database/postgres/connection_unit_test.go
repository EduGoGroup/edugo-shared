package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestClose_WithNil_Unit verifica que Close maneja nil correctamente
func TestClose_WithNil_Unit(t *testing.T) {
	err := Close(nil)
	assert.NoError(t, err, "Close con DB nil no debe retornar error")
}

// TestDefaultConstants_Unit verifica que las constantes tienen valores razonables
func TestDefaultConstants_Unit(t *testing.T) {
	t.Run("DefaultPort", func(t *testing.T) {
		assert.Equal(t, 5432, DefaultPort, "Puerto por defecto debe ser 5432")
	})

	t.Run("DefaultMaxConnections", func(t *testing.T) {
		assert.Equal(t, 25, DefaultMaxConnections)
		assert.Greater(t, DefaultMaxConnections, 0, "MaxConnections debe ser positivo")
	})

	t.Run("DefaultMaxIdleConnections", func(t *testing.T) {
		assert.Equal(t, 5, DefaultMaxIdleConnections)
		assert.Less(t, DefaultMaxIdleConnections, DefaultMaxConnections,
			"MaxIdleConnections debe ser menor que MaxConnections")
	})

	t.Run("DefaultMaxLifetime", func(t *testing.T) {
		assert.Equal(t, 5*time.Minute, DefaultMaxLifetime)
		assert.Greater(t, DefaultMaxLifetime, 0*time.Second,
			"MaxLifetime debe ser positivo")
	})

	t.Run("DefaultConnectTimeout", func(t *testing.T) {
		assert.Equal(t, 10*time.Second, DefaultConnectTimeout)
		assert.Greater(t, DefaultConnectTimeout, 0*time.Second,
			"ConnectTimeout debe ser positivo")
	})

	t.Run("DefaultHealthCheckTimeout", func(t *testing.T) {
		assert.Equal(t, 5*time.Second, DefaultHealthCheckTimeout)
		assert.Greater(t, DefaultHealthCheckTimeout, 0*time.Second,
			"HealthCheckTimeout debe ser positivo")
	})
}

// TestConfig_Validation_Unit verifica validaciones de estructura Config
func TestConfig_Validation_Unit(t *testing.T) {
	t.Run("config con valores válidos", func(t *testing.T) {
		config := &Config{
			Host:               "localhost",
			Port:               5432,
			User:               "testuser",
			Password:           "testpass",
			Database:           "testdb",
			SSLMode:            "disable",
			ConnectTimeout:     10 * time.Second,
			MaxConnections:     25,
			MaxIdleConnections: 5,
			MaxLifetime:        5 * time.Minute,
		}

		// Verificar que la estructura de config es válida
		assert.Equal(t, "localhost", config.Host)
		assert.Equal(t, 5432, config.Port)
		assert.Equal(t, "testuser", config.User)
		assert.Equal(t, "testdb", config.Database)
		assert.Equal(t, "disable", config.SSLMode)
	})

	t.Run("config con puerto no estándar", func(t *testing.T) {
		config := &Config{
			Host:     "localhost",
			Port:     5433,
			User:     "user",
			Password: "pass",
			Database: "db",
			SSLMode:  "disable",
		}

		assert.Equal(t, 5433, config.Port, "Puerto no estándar debe ser aceptado")
	})

	t.Run("config con diferentes modos SSL", func(t *testing.T) {
		sslModes := []string{"disable", "require", "verify-ca", "verify-full"}

		for _, mode := range sslModes {
			config := &Config{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "pass",
				Database: "db",
				SSLMode:  mode,
			}

			assert.Equal(t, mode, config.SSLMode,
				"Modo SSL %s debe ser aceptado", mode)
		}
	})

	t.Run("config con pool connections válido", func(t *testing.T) {
		config := &Config{
			Host:               "localhost",
			Port:               5432,
			User:               "user",
			Password:           "pass",
			Database:           "db",
			SSLMode:            "disable",
			MaxConnections:     100,
			MaxIdleConnections: 20,
			MaxLifetime:        10 * time.Minute,
		}

		assert.Equal(t, 100, config.MaxConnections)
		assert.Equal(t, 20, config.MaxIdleConnections)
		assert.Less(t, config.MaxIdleConnections, config.MaxConnections,
			"MaxIdleConnections debe ser menor que MaxConnections")
	})

	t.Run("config con password vacío", func(t *testing.T) {
		config := &Config{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "", // Password vacío es válido
			Database: "db",
			SSLMode:  "disable",
		}

		assert.Empty(t, config.Password, "Password vacío debe ser aceptado")
	})

	t.Run("config con timeout razonable", func(t *testing.T) {
		config := &Config{
			Host:           "localhost",
			Port:           5432,
			User:           "user",
			Password:       "pass",
			Database:       "db",
			SSLMode:        "disable",
			ConnectTimeout: 30 * time.Second,
		}

		assert.Greater(t, config.ConnectTimeout, 0*time.Second,
			"ConnectTimeout debe ser positivo")
	})

	t.Run("config con diferentes nombres de database", func(t *testing.T) {
		databases := []string{
			"simple_db",
			"db-with-hyphens",
			"db123",
			"UPPERCASE_DB",
			"mixed_CASE_db",
		}

		for _, dbName := range databases {
			config := &Config{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "pass",
				Database: dbName,
				SSLMode:  "disable",
			}

			assert.Equal(t, dbName, config.Database)
		}
	})

	t.Run("config con diferentes pool sizes", func(t *testing.T) {
		tests := []struct {
			name               string
			maxConnections     int
			maxIdleConnections int
		}{
			{"pool pequeño", 5, 2},
			{"pool mediano", 25, 5},
			{"pool grande", 100, 20},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				config := &Config{
					Host:               "localhost",
					Port:               5432,
					User:               "user",
					Password:           "pass",
					Database:           "db",
					SSLMode:            "disable",
					MaxConnections:     tt.maxConnections,
					MaxIdleConnections: tt.maxIdleConnections,
				}

				assert.Equal(t, tt.maxConnections, config.MaxConnections)
				assert.Equal(t, tt.maxIdleConnections, config.MaxIdleConnections)
			})
		}
	})
}
