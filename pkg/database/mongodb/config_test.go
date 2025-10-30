package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("retorna configuración con valores por defecto correctos", func(t *testing.T) {
		config := DefaultConfig()

		assert.Equal(t, "mongodb://localhost:27017", config.URI)
		assert.Equal(t, "test", config.Database)
		assert.Equal(t, 10*time.Second, config.Timeout)
		assert.Equal(t, uint64(100), config.MaxPoolSize)
		assert.Equal(t, uint64(10), config.MinPoolSize)
	})

	t.Run("pool size defaults son razonables", func(t *testing.T) {
		config := DefaultConfig()

		assert.Greater(t, config.MaxPoolSize, config.MinPoolSize,
			"MaxPoolSize debe ser mayor que MinPoolSize")
		assert.GreaterOrEqual(t, config.MaxPoolSize, uint64(50),
			"MaxPoolSize por defecto debe soportar concurrencia moderada")
	})
}

func TestConfigCustomization(t *testing.T) {
	t.Run("permite personalizar todos los campos", func(t *testing.T) {
		config := Config{
			URI:         "mongodb://user:pass@custom-host:27017/admin",
			Database:    "custom-db",
			Timeout:     30 * time.Second,
			MaxPoolSize: 200,
			MinPoolSize: 20,
		}

		assert.Equal(t, "mongodb://user:pass@custom-host:27017/admin", config.URI)
		assert.Equal(t, "custom-db", config.Database)
		assert.Equal(t, 30*time.Second, config.Timeout)
		assert.Equal(t, uint64(200), config.MaxPoolSize)
		assert.Equal(t, uint64(20), config.MinPoolSize)
	})

	t.Run("puede modificar config por defecto", func(t *testing.T) {
		config := DefaultConfig()

		// Modificar valores
		config.URI = "mongodb://prod-host:27017"
		config.Database = "prod-db"
		config.Timeout = 20 * time.Second
		config.MaxPoolSize = 150

		assert.Equal(t, "mongodb://prod-host:27017", config.URI)
		assert.Equal(t, "prod-db", config.Database)
		assert.Equal(t, 20*time.Second, config.Timeout)
		assert.Equal(t, uint64(150), config.MaxPoolSize)

		// Valores no modificados mantienen defaults
		assert.Equal(t, uint64(10), config.MinPoolSize)
	})
}

func TestConfigURIFormats(t *testing.T) {
	testCases := []struct {
		name        string
		uri         string
		description string
	}{
		{
			name:        "URI simple localhost",
			uri:         "mongodb://localhost:27017",
			description: "Conexión local sin autenticación",
		},
		{
			name:        "URI con autenticación",
			uri:         "mongodb://user:password@localhost:27017/admin",
			description: "Conexión con usuario y password",
		},
		{
			name:        "URI con autenticación y authSource",
			uri:         "mongodb://user:password@localhost:27017/mydb?authSource=admin",
			description: "Especifica base de datos de autenticación",
		},
		{
			name:        "URI con múltiples hosts (replica set)",
			uri:         "mongodb://host1:27017,host2:27017,host3:27017/mydb",
			description: "Replica set con múltiples nodos",
		},
		{
			name:        "URI con opciones de conexión",
			uri:         "mongodb://localhost:27017/mydb?maxPoolSize=50&minPoolSize=5",
			description: "URI con parámetros de pool",
		},
		{
			name:        "URI de MongoDB Atlas",
			uri:         "mongodb+srv://user:password@cluster.mongodb.net/mydb",
			description: "Conexión a MongoDB Atlas (SRV record)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := Config{
				URI:         tc.uri,
				Database:    "test",
				Timeout:     10 * time.Second,
				MaxPoolSize: 100,
				MinPoolSize: 10,
			}

			assert.Equal(t, tc.uri, config.URI)
			assert.NotEmpty(t, config.URI, "URI no debe estar vacía")
		})
	}
}

func TestConfigConnectionPool(t *testing.T) {
	t.Run("MaxPoolSize debe ser mayor o igual que MinPoolSize", func(t *testing.T) {
		config := DefaultConfig()

		assert.GreaterOrEqual(t, config.MaxPoolSize, config.MinPoolSize,
			"MaxPoolSize debe ser mayor o igual que MinPoolSize")
	})

	t.Run("permite configurar pool pequeño para desarrollo", func(t *testing.T) {
		config := Config{
			URI:         "mongodb://localhost:27017",
			Database:    "dev",
			Timeout:     5 * time.Second,
			MaxPoolSize: 10,
			MinPoolSize: 2,
		}

		assert.Equal(t, uint64(10), config.MaxPoolSize)
		assert.Equal(t, uint64(2), config.MinPoolSize)
	})

	t.Run("permite configurar pool grande para producción", func(t *testing.T) {
		config := Config{
			URI:         "mongodb://prod-host:27017",
			Database:    "prod",
			Timeout:     30 * time.Second,
			MaxPoolSize: 500,
			MinPoolSize: 50,
		}

		assert.Equal(t, uint64(500), config.MaxPoolSize)
		assert.Equal(t, uint64(50), config.MinPoolSize)
	})

	t.Run("pool con tamaño mínimo = máximo (pool fijo)", func(t *testing.T) {
		config := Config{
			URI:         "mongodb://localhost:27017",
			Database:    "test",
			Timeout:     10 * time.Second,
			MaxPoolSize: 25,
			MinPoolSize: 25, // Pool de tamaño fijo
		}

		assert.Equal(t, config.MaxPoolSize, config.MinPoolSize,
			"Pool de tamaño fijo: min = max")
	})
}

func TestConfigTimeouts(t *testing.T) {
	t.Run("permite configurar timeouts cortos", func(t *testing.T) {
		config := DefaultConfig()
		config.Timeout = 3 * time.Second

		assert.Equal(t, 3*time.Second, config.Timeout)
	})

	t.Run("permite configurar timeouts largos", func(t *testing.T) {
		config := DefaultConfig()
		config.Timeout = 60 * time.Second

		assert.Equal(t, 60*time.Second, config.Timeout)
	})

	t.Run("timeout por defecto es razonable", func(t *testing.T) {
		config := DefaultConfig()

		assert.Equal(t, 10*time.Second, config.Timeout,
			"Timeout por defecto debe ser razonable (10 segundos)")
		assert.GreaterOrEqual(t, config.Timeout, 5*time.Second,
			"Timeout no debe ser demasiado corto")
		assert.LessOrEqual(t, config.Timeout, 30*time.Second,
			"Timeout no debe ser excesivamente largo")
	})
}

func TestConfigForDifferentEnvironments(t *testing.T) {
	t.Run("configuración para desarrollo local", func(t *testing.T) {
		devConfig := Config{
			URI:         "mongodb://localhost:27017",
			Database:    "dev_db",
			Timeout:     5 * time.Second,
			MaxPoolSize: 20,
			MinPoolSize: 5,
		}

		assert.Contains(t, devConfig.URI, "localhost")
		assert.Equal(t, uint64(20), devConfig.MaxPoolSize, "Desarrollo usa pool más pequeño")
	})

	t.Run("configuración para producción con replica set", func(t *testing.T) {
		prodConfig := Config{
			URI:         "mongodb://user:pass@host1:27017,host2:27017,host3:27017/prod?replicaSet=rs0&authSource=admin",
			Database:    "prod_db",
			Timeout:     30 * time.Second,
			MaxPoolSize: 300,
			MinPoolSize: 50,
		}

		assert.NotContains(t, prodConfig.URI, "localhost")
		assert.Contains(t, prodConfig.URI, "replicaSet")
		assert.GreaterOrEqual(t, prodConfig.MaxPoolSize, uint64(100),
			"Producción necesita pool más grande")
	})

	t.Run("configuración para MongoDB Atlas", func(t *testing.T) {
		atlasConfig := Config{
			URI:         "mongodb+srv://user:pass@cluster0.example.mongodb.net/mydb?retryWrites=true&w=majority",
			Database:    "atlas_db",
			Timeout:     20 * time.Second,
			MaxPoolSize: 150,
			MinPoolSize: 30,
		}

		assert.Contains(t, atlasConfig.URI, "mongodb+srv")
		assert.Contains(t, atlasConfig.URI, "mongodb.net")
		assert.Equal(t, uint64(150), atlasConfig.MaxPoolSize)
	})

	t.Run("configuración para testing", func(t *testing.T) {
		testConfig := Config{
			URI:         "mongodb://localhost:27018", // Puerto diferente
			Database:    "test_db",
			Timeout:     3 * time.Second,
			MaxPoolSize: 10,
			MinPoolSize: 1,
		}

		assert.Equal(t, "mongodb://localhost:27018", testConfig.URI)
		assert.Contains(t, testConfig.URI, "27018", "Testing usa puerto diferente")
		assert.Equal(t, uint64(10), testConfig.MaxPoolSize, "Testing usa pool pequeño")
	})
}

func TestConfigValidation(t *testing.T) {
	t.Run("configuración válida mínima", func(t *testing.T) {
		config := Config{
			URI:      "mongodb://localhost:27017",
			Database: "mydb",
		}

		assert.NotEmpty(t, config.URI, "URI es requerida")
		assert.NotEmpty(t, config.Database, "Database es requerida")
	})

	t.Run("config sin timeout usa default de 0", func(t *testing.T) {
		config := Config{
			URI:      "mongodb://localhost:27017",
			Database: "mydb",
			// Timeout no especificado (valor zero)
		}

		assert.Equal(t, time.Duration(0), config.Timeout,
			"Timeout no especificado debe ser 0 (usará default del driver)")
	})

	t.Run("config sin pool sizes usa default de 0", func(t *testing.T) {
		config := Config{
			URI:      "mongodb://localhost:27017",
			Database: "mydb",
			// Pool sizes no especificados
		}

		assert.Equal(t, uint64(0), config.MaxPoolSize,
			"MaxPoolSize no especificado debe ser 0 (usará default del driver)")
		assert.Equal(t, uint64(0), config.MinPoolSize,
			"MinPoolSize no especificado debe ser 0 (usará default del driver)")
	})
}
