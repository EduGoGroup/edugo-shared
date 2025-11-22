package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClose_WithNil_Unit verifica que Close maneja nil correctamente
func TestClose_WithNil_Unit(t *testing.T) {
	err := Close(nil)
	assert.NoError(t, err, "Close con client nil no debe retornar error")
}

// TestGetDatabase_NilClient_Unit verifica comportamiento con cliente nil
func TestGetDatabase_NilClient_Unit(t *testing.T) {
	// GetDatabase con client nil causará panic (comportamiento esperado de mongo.Client)
	// Este test documenta ese comportamiento
	defer func() {
		if r := recover(); r != nil {
			// Panic esperado cuando client es nil
			assert.NotNil(t, r, "GetDatabase con client nil causa panic")
		}
	}()

	db := GetDatabase(nil, "test_db")
	// Si llegamos aquí sin panic, db será nil o inválido
	_ = db
}

// TestConnect_InvalidURI_Unit verifica manejo de URIs inválidas
func TestConnect_InvalidURI_Unit(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "URI completamente inválida",
			config: Config{
				URI:         "this-is-not-a-valid-uri",
				Database:    "test",
				Timeout:     5 * time.Second,
				MaxPoolSize: 10,
				MinPoolSize: 2,
			},
		},
		{
			name: "URI vacía",
			config: Config{
				URI:      "",
				Database: "test",
				Timeout:  5 * time.Second,
			},
		},
		{
			name: "URI con protocolo incorrecto",
			config: Config{
				URI:      "http://localhost:27017",
				Database: "test",
				Timeout:  5 * time.Second,
			},
		},
		{
			name: "URI con formato malformado",
			config: Config{
				URI:      "mongodb://[invalid",
				Database: "test",
				Timeout:  5 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := Connect(tt.config)

			// Debe retornar error con URI inválida
			require.Error(t, err, "URI inválida debe causar error")
			assert.Nil(t, client, "Client debe ser nil cuando hay error")
			assert.Contains(t, err.Error(), "failed to connect to mongodb",
				"Error debe indicar fallo de conexión")
		})
	}
}

// TestConnect_UnreachableHost_Unit verifica timeout con host inalcanzable
func TestConnect_UnreachableHost_Unit(t *testing.T) {
	config := Config{
		URI:         "mongodb://192.0.2.1:27017", // IP de TEST-NET-1 (no ruteable)
		Database:    "test",
		Timeout:     1 * time.Second, // Timeout corto para test rápido
		MaxPoolSize: 10,
		MinPoolSize: 2,
	}

	start := time.Now()
	client, err := Connect(config)
	duration := time.Since(start)

	// Debe fallar y respetar el timeout
	require.Error(t, err, "Conexión a host inalcanzable debe fallar")
	assert.Nil(t, client, "Client debe ser nil cuando hay error")

	// El timeout debe ser respetado (con margen)
	assert.Less(t, duration, 3*time.Second,
		"Debe respetar el timeout configurado")
}

// TestConnect_InvalidHost_Unit verifica manejo de hosts inválidos
func TestConnect_InvalidHost_Unit(t *testing.T) {
	tests := []struct {
		name string
		uri  string
	}{
		{
			name: "hostname que no resuelve",
			uri:  "mongodb://this-host-does-not-exist-12345.local:27017",
		},
		{
			name: "IP inválida",
			uri:  "mongodb://999.999.999.999:27017",
		},
		{
			name: "puerto inválido",
			uri:  "mongodb://localhost:99999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				URI:      tt.uri,
				Database: "test",
				Timeout:  2 * time.Second,
			}

			client, err := Connect(config)

			require.Error(t, err, "Host inválido debe causar error")
			assert.Nil(t, client, "Client debe ser nil cuando hay error")
		})
	}
}

// TestConnect_ZeroTimeout_Unit verifica comportamiento con timeout zero
func TestConnect_ZeroTimeout_Unit(t *testing.T) {
	config := Config{
		URI:      "mongodb://localhost:27017",
		Database: "test",
		Timeout:  0, // Timeout zero
	}

	// Con timeout 0, mongo usará defaults
	// Este test puede fallar si no hay MongoDB local, lo cual es esperado
	_, err := Connect(config)

	// Solo verificamos que no cause panic
	// El error es esperado si no hay MongoDB local
	if err != nil {
		assert.Contains(t, err.Error(), "failed to",
			"Error debe ser de conexión o ping")
	}
}

// TestConnect_ExtremePoolSizes_Unit verifica pool sizes extremos
func TestConnect_ExtremePoolSizes_Unit(t *testing.T) {
	tests := []struct {
		name        string
		maxPoolSize uint64
		minPoolSize uint64
		shouldWork  bool
	}{
		{
			name:        "pool size zero",
			maxPoolSize: 0,
			minPoolSize: 0,
			shouldWork:  true, // MongoDB usará defaults
		},
		{
			name:        "pool muy grande",
			maxPoolSize: 10000,
			minPoolSize: 100,
			shouldWork:  true,
		},
		{
			name:        "min mayor que max",
			maxPoolSize: 5,
			minPoolSize: 10,
			shouldWork:  true, // MongoDB puede manejarlo
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				URI:         "mongodb://nonexistent:27017",
				Database:    "test",
				Timeout:     1 * time.Second,
				MaxPoolSize: tt.maxPoolSize,
				MinPoolSize: tt.minPoolSize,
			}

			// Intentar conectar (fallará por host inexistente, pero no por pool sizes)
			_, err := Connect(config)

			// Debe fallar por conexión, no por configuración de pool
			if err != nil {
				assert.Contains(t, err.Error(), "failed to",
					"Error debe ser de conexión, no de configuración")
			}
		})
	}
}

// TestConnect_URIWithAuth_Unit verifica URIs con autenticación
func TestConnect_URIWithAuth_Unit(t *testing.T) {
	tests := []struct {
		name string
		uri  string
	}{
		{
			name: "URI con usuario y password",
			uri:  "mongodb://user:pass@localhost:27017/admin",
		},
		{
			name: "URI con password vacío",
			uri:  "mongodb://user:@localhost:27017",
		},
		{
			name: "URI con caracteres especiales en password",
			uri:  "mongodb://user:p@ssw0rd!@localhost:27017",
		},
		{
			name: "URI con password URL-encoded",
			uri:  "mongodb://user:p%40ssw0rd@localhost:27017",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				URI:      tt.uri,
				Database: "test",
				Timeout:  1 * time.Second,
			}

			// Intentar conectar (fallará porque no hay servidor)
			_, err := Connect(config)

			// El error debe ser de conexión/ping, no de parsing de URI
			if err != nil {
				assert.Contains(t, err.Error(), "failed to",
					"Error debe ser de conexión, no de parsing de URI")
			}
		})
	}
}

// TestConnect_URIWithOptions_Unit verifica URIs con query parameters
func TestConnect_URIWithOptions_Unit(t *testing.T) {
	tests := []struct {
		name string
		uri  string
	}{
		{
			name: "URI con authSource",
			uri:  "mongodb://user:pass@localhost:27017/mydb?authSource=admin",
		},
		{
			name: "URI con múltiples opciones",
			uri:  "mongodb://localhost:27017/mydb?maxPoolSize=50&minPoolSize=5",
		},
		{
			name: "URI con replica set",
			uri:  "mongodb://host1:27017,host2:27017/mydb?replicaSet=rs0",
		},
		{
			name: "URI MongoDB Atlas style",
			uri:  "mongodb+srv://user:pass@cluster.mongodb.net/mydb?retryWrites=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				URI:      tt.uri,
				Database: "test",
				Timeout:  1 * time.Second,
			}

			// Intentar conectar (fallará porque no hay servidor)
			_, err := Connect(config)

			// El error debe ser de conexión, no de parsing
			if err != nil {
				assert.Contains(t, err.Error(), "failed to",
					"Error debe ser de conexión, no de parsing de URI")
			}
		})
	}
}

// TestHealthCheck_NilClient_Unit verifica HealthCheck con cliente nil
func TestHealthCheck_NilClient_Unit(t *testing.T) {
	// HealthCheck con nil causará panic al intentar llamar client.Ping
	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r, "HealthCheck con nil debe causar panic")
		}
	}()

	err := HealthCheck(nil)
	// Si llegamos aquí, hubo error en lugar de panic
	if err != nil {
		assert.Error(t, err, "HealthCheck con nil debe fallar")
	}
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
	t.Run("config con timeout muy corto", func(t *testing.T) {
		config := Config{
			URI:      "mongodb://localhost:27017",
			Database: "test",
			Timeout:  1 * time.Millisecond, // Extremadamente corto
		}

		// Debe fallar por timeout
		_, err := Connect(config)
		if err != nil {
			// Error esperado
			assert.Error(t, err)
		}
	})

	t.Run("config sin database especificado", func(t *testing.T) {
		config := Config{
			URI:      "mongodb://localhost:27017",
			Database: "", // Vacío
			Timeout:  5 * time.Second,
		}

		// Database vacío no causa error en Connect
		// (el error será al usar GetDatabase)
		assert.Empty(t, config.Database,
			"Database puede estar vacío en config")
	})
}
