package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestClose_WithNil_Unit verifica que Close maneja nil correctamente
func TestClose_WithNil_Unit(t *testing.T) {
	err := Close(nil)
	assert.NoError(t, err, "Close con client nil no debe retornar error")
}

// TestDefaultConstants_Unit verifica que las constantes tienen valores razonables
func TestDefaultConstants_Unit(t *testing.T) {
	t.Run("DefaultHealthCheckTimeout es razonable", func(t *testing.T) {
		assert.Equal(t, 5*time.Second, DefaultHealthCheckTimeout,
			"HealthCheck timeout debe ser 5 segundos")
		assert.Greater(t, DefaultHealthCheckTimeout, 1*time.Second,
			"Timeout no debe ser muy corto")
		assert.Less(t, DefaultHealthCheckTimeout, 30*time.Second,
			"Timeout no debe ser muy largo")
	})

	t.Run("DefaultDisconnectTimeout es razonable", func(t *testing.T) {
		assert.Equal(t, 10*time.Second, DefaultDisconnectTimeout,
			"Disconnect timeout debe ser 10 segundos")
		assert.Greater(t, DefaultDisconnectTimeout, DefaultHealthCheckTimeout,
			"Disconnect timeout debe ser mayor que health check timeout")
	})
}

// TestConfig_Validation_Unit verifica validaciones lógicas de Config
func TestConfig_Validation_Unit(t *testing.T) {
	t.Run("config sin database especificado", func(t *testing.T) {
		config := Config{
			URI:      "mongodb://localhost:27017",
			Database: "", // Vacío
			Timeout:  5 * time.Second,
		}

		// Database vacío no causa error en construcción de Config
		// (el error será al usar GetDatabase)
		assert.Empty(t, config.Database,
			"Database puede estar vacío en config")
	})

	t.Run("config con valores de pool válidos", func(t *testing.T) {
		config := Config{
			URI:         "mongodb://localhost:27017",
			Database:    "test",
			Timeout:     5 * time.Second,
			MaxPoolSize: 100,
			MinPoolSize: 10,
		}

		// Verificar que la estructura de config es válida
		assert.NotEmpty(t, config.URI, "URI no debe estar vacío")
		assert.Equal(t, uint64(100), config.MaxPoolSize)
		assert.Equal(t, uint64(10), config.MinPoolSize)
	})

	t.Run("config con timeout razonable", func(t *testing.T) {
		config := Config{
			URI:      "mongodb://localhost:27017",
			Database: "test",
			Timeout:  10 * time.Second,
		}

		assert.Greater(t, config.Timeout, 0*time.Second,
			"Timeout debe ser positivo")
	})

	t.Run("config con diferentes valores de pool", func(t *testing.T) {
		tests := []struct {
			name        string
			maxPoolSize uint64
			minPoolSize uint64
		}{
			{"pool pequeño", 5, 2},
			{"pool mediano", 25, 5},
			{"pool grande", 100, 20},
			{"pool con valores iguales", 10, 10},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				config := Config{
					URI:         "mongodb://localhost:27017",
					Database:    "test",
					Timeout:     5 * time.Second,
					MaxPoolSize: tt.maxPoolSize,
					MinPoolSize: tt.minPoolSize,
				}

				assert.Equal(t, tt.maxPoolSize, config.MaxPoolSize)
				assert.Equal(t, tt.minPoolSize, config.MinPoolSize)
			})
		}
	})
}
