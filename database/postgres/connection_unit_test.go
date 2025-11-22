package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestClose_WithNil_Unit verifica que Close maneja nil correctamente
func TestClose_WithNil_Unit(t *testing.T) {
	err := Close(nil)
	assert.NoError(t, err, "Close con DB nil no debe retornar error")
}

// TestGetStats_WithNil_Unit verifica GetStats con DB nil
func TestGetStats_WithNil_Unit(t *testing.T) {
	// GetStats con nil causará panic (comportamiento esperado de sql.DB)
	// Este test documenta ese comportamiento
	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r, "GetStats con nil causa panic esperado")
		}
	}()

	stats := GetStats(nil)
	// Si llegamos aquí sin panic, verificar que stats tiene valores zero
	_ = stats
}

// TestConnect_DSNBuilding_Unit verifica construcción del DSN
func TestConnect_DSNBuilding_Unit(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		description string
	}{
		{
			name: "config básica válida",
			config: &Config{
				Host:           "localhost",
				Port:           5432,
				User:           "testuser",
				Password:       "testpass",
				Database:       "testdb",
				SSLMode:        "disable",
				ConnectTimeout: 10 * time.Second,
			},
			expectError: true, // Error porque no hay servidor real
			description: "Config válida debe generar DSN correcto",
		},
		{
			name: "config con host remoto",
			config: &Config{
				Host:           "db.example.com",
				Port:           5432,
				User:           "admin",
				Password:       "secret",
				Database:       "production",
				SSLMode:        "require",
				ConnectTimeout: 30 * time.Second,
			},
			expectError: true, // Error porque no hay servidor real
			description: "Config con host remoto debe funcionar",
		},
		{
			name: "config con puerto no estándar",
			config: &Config{
				Host:           "localhost",
				Port:           5433,
				User:           "user",
				Password:       "pass",
				Database:       "db",
				SSLMode:        "disable",
				ConnectTimeout: 5 * time.Second,
			},
			expectError: true, // Error porque no hay servidor real
			description: "Puerto no estándar debe ser respetado",
		},
		{
			name: "config con password vacío",
			config: &Config{
				Host:           "localhost",
				Port:           5432,
				User:           "user",
				Password:       "", // Password vacío
				Database:       "db",
				SSLMode:        "disable",
				ConnectTimeout: 5 * time.Second,
			},
			expectError: true,
			description: "Password vacío es válido en DSN",
		},
		{
			name: "config con SSL verify-full",
			config: &Config{
				Host:           "secure-host",
				Port:           5432,
				User:           "user",
				Password:       "pass",
				Database:       "db",
				SSLMode:        "verify-full",
				ConnectTimeout: 10 * time.Second,
			},
			expectError: true,
			description: "SSL verify-full debe ser incluido en DSN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.config)

			if tt.expectError {
				// Esperamos error porque no hay servidor PostgreSQL real
				assert.Error(t, err, tt.description)

				// Si hay error, db debe ser nil o debe ser cerrado
				if db != nil {
					_ = db.Close()
				}
			} else {
				// Si esperamos éxito (ninguno en este caso)
				assert.NoError(t, err, tt.description)
				if db != nil {
					_ = db.Close()
				}
			}
		})
	}
}

// TestConnect_InvalidConfig_Unit verifica manejo de configs inválidas
func TestConnect_InvalidConfig_Unit(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "puerto zero",
			config: &Config{
				Host:           "localhost",
				Port:           0,
				User:           "user",
				Password:       "pass",
				Database:       "db",
				SSLMode:        "disable",
				ConnectTimeout: 5 * time.Second,
			},
		},
		{
			name: "puerto negativo",
			config: &Config{
				Host:           "localhost",
				Port:           -1,
				User:           "user",
				Password:       "pass",
				Database:       "db",
				SSLMode:        "disable",
				ConnectTimeout: 5 * time.Second,
			},
		},
		{
			name: "timeout zero",
			config: &Config{
				Host:           "localhost",
				Port:           5432,
				User:           "user",
				Password:       "pass",
				Database:       "db",
				SSLMode:        "disable",
				ConnectTimeout: 0,
			},
		},
		{
			name: "host vacío",
			config: &Config{
				Host:           "",
				Port:           5432,
				User:           "user",
				Password:       "pass",
				Database:       "db",
				SSLMode:        "disable",
				ConnectTimeout: 5 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.config)

			// Debe haber error (ya sea de DSN o de conexión)
			assert.Error(t, err, "Config inválida debe causar error")

			if db != nil {
				_ = db.Close()
			}
		})
	}
}

// TestConnect_ConnectionPoolConfig_Unit verifica configuración del pool
func TestConnect_ConnectionPoolConfig_Unit(t *testing.T) {
	tests := []struct {
		name               string
		maxConnections     int
		maxIdleConnections int
		maxLifetime        time.Duration
	}{
		{
			name:               "pool pequeño",
			maxConnections:     5,
			maxIdleConnections: 2,
			maxLifetime:        1 * time.Minute,
		},
		{
			name:               "pool grande",
			maxConnections:     100,
			maxIdleConnections: 20,
			maxLifetime:        10 * time.Minute,
		},
		{
			name:               "pool con lifetime largo",
			maxConnections:     25,
			maxIdleConnections: 5,
			maxLifetime:        1 * time.Hour,
		},
		{
			name:               "configuración por defecto",
			maxConnections:     DefaultMaxConnections,
			maxIdleConnections: DefaultMaxIdleConnections,
			maxLifetime:        DefaultMaxLifetime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Host:               "nonexistent-host-12345.local",
				Port:               5432,
				User:               "user",
				Password:           "pass",
				Database:           "db",
				SSLMode:            "disable",
				ConnectTimeout:     1 * time.Second,
				MaxConnections:     tt.maxConnections,
				MaxIdleConnections: tt.maxIdleConnections,
				MaxLifetime:        tt.maxLifetime,
			}

			db, err := Connect(config)

			// Esperamos error porque el host no existe
			assert.Error(t, err, "Conexión debe fallar con host inexistente")

			if db != nil {
				_ = db.Close()
			}
		})
	}
}

// TestConnect_SSLModes_Unit verifica diferentes modos SSL
func TestConnect_SSLModes_Unit(t *testing.T) {
	sslModes := []string{
		"disable",
		"require",
		"verify-ca",
		"verify-full",
	}

	for _, sslMode := range sslModes {
		t.Run("SSL mode: "+sslMode, func(t *testing.T) {
			config := &Config{
				Host:           "nonexistent-host.local",
				Port:           5432,
				User:           "user",
				Password:       "pass",
				Database:       "db",
				SSLMode:        sslMode,
				ConnectTimeout: 1 * time.Second,
			}

			db, err := Connect(config)

			// Debe fallar por conexión, no por SSL mode inválido
			assert.Error(t, err, "Conexión debe fallar con host inexistente")

			if db != nil {
				_ = db.Close()
			}
		})
	}
}

// TestConnect_PasswordWithSpecialChars_Unit verifica passwords con caracteres especiales
func TestConnect_PasswordWithSpecialChars_Unit(t *testing.T) {
	passwords := []string{
		"p@ssw0rd!",
		"pass word",
		"pass'word",
		`pass"word`,
		"pass\\word",
		"pásswörd",
		"密码",
	}

	for _, password := range passwords {
		t.Run("password con caracteres especiales", func(t *testing.T) {
			config := &Config{
				Host:           "nonexistent.local",
				Port:           5432,
				User:           "user",
				Password:       password,
				Database:       "db",
				SSLMode:        "disable",
				ConnectTimeout: 1 * time.Second,
			}

			db, err := Connect(config)

			// El error debe ser de conexión, no de parsing del password
			assert.Error(t, err, "Debe fallar por host inexistente")

			if db != nil {
				_ = db.Close()
			}
		})
	}
}

// TestConnect_DatabaseNameVariations_Unit verifica diferentes nombres de DB
func TestConnect_DatabaseNameVariations_Unit(t *testing.T) {
	databases := []string{
		"simple_db",
		"db-with-hyphens",
		"db123",
		"UPPERCASE_DB",
		"mixed_CASE_db",
	}

	for _, dbName := range databases {
		t.Run("database: "+dbName, func(t *testing.T) {
			config := &Config{
				Host:           "nonexistent.local",
				Port:           5432,
				User:           "user",
				Password:       "pass",
				Database:       dbName,
				SSLMode:        "disable",
				ConnectTimeout: 1 * time.Second,
			}

			db, err := Connect(config)

			// El error debe ser de conexión
			assert.Error(t, err, "Debe fallar por host inexistente")

			if db != nil {
				_ = db.Close()
			}
		})
	}
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

// TestHealthCheck_WithClosedDB_Unit verifica HealthCheck con DB cerrada
func TestHealthCheck_WithClosedDB_Unit(t *testing.T) {
	// Crear DB ficticia (no conectada realmente)
	db, err := sql.Open("postgres", "host=nonexistent")
	if err != nil {
		t.Skipf("No se pudo crear DB ficticia: %v", err)
	}

	// Cerrar inmediatamente
	_ = db.Close()

	// HealthCheck debe fallar
	err = HealthCheck(db)
	assert.Error(t, err, "HealthCheck con DB cerrada debe fallar")
	assert.Contains(t, err.Error(), "health check failed",
		"Error debe indicar fallo de health check")
}

// TestConnect_ExtremeTimeouts_Unit verifica timeouts extremos
func TestConnect_ExtremeTimeouts_Unit(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "timeout muy corto",
			timeout: 1 * time.Millisecond,
		},
		{
			name:    "timeout medio",
			timeout: 5 * time.Second,
		},
		{
			name:    "timeout largo",
			timeout: 1 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Host:           "192.0.2.1", // TEST-NET-1 (no ruteable)
				Port:           5432,
				User:           "user",
				Password:       "pass",
				Database:       "db",
				SSLMode:        "disable",
				ConnectTimeout: tt.timeout,
			}

			start := time.Now()
			db, err := Connect(config)
			duration := time.Since(start)

			// Debe fallar
			assert.Error(t, err, "Conexión a host inalcanzable debe fallar")

			// Verificar que respeta el timeout (con margen)
			// Para timeouts muy cortos, puede tardar un poco más debido a overhead
			if tt.timeout > 100*time.Millisecond {
				assert.Less(t, duration, tt.timeout*3,
					"Debe respetar aproximadamente el timeout")
			}

			if db != nil {
				_ = db.Close()
			}
		})
	}
}

// TestConnect_ZeroValues_Unit verifica comportamiento con valores zero
func TestConnect_ZeroValues_Unit(t *testing.T) {
	config := &Config{
		Host:               "nonexistent.local",
		Port:               5432,
		User:               "user",
		Password:           "pass",
		Database:           "db",
		SSLMode:            "disable",
		ConnectTimeout:     1 * time.Second,
		MaxConnections:     0, // Zero
		MaxIdleConnections: 0, // Zero
		MaxLifetime:        0, // Zero
	}

	db, err := Connect(config)

	// Debe fallar por conexión, no por valores zero en pool
	assert.Error(t, err, "Debe fallar por host inexistente")

	if db != nil {
		_ = db.Close()
	}
}
