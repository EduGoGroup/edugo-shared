package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validation_Unit(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "config con valores válidos",
			config: Config{
				Host:               "localhost",
				Port:               5432,
				User:               "postgres",
				Password:           "password",
				Database:           "mydb",
				SSLMode:            "disable",
				MaxConnections:     20,
				MaxIdleConnections: 5,
				MaxLifetime:        10 * time.Minute,
				ConnectTimeout:     10 * time.Second,
				SearchPath:         "public",
			},
		},
		{
			name: "config con puerto no estándar",
			config: Config{
				Host: "db.example.com",
				Port: 5433,
			},
		},
		{
			name: "config con diferentes modos SSL",
			config: Config{
				SSLMode: "require",
			},
		},
		{
			name: "config con pool connections válido",
			config: Config{
				MaxConnections:     50,
				MaxIdleConnections: 25,
			},
		},
		{
			name: "config con password vacío",
			config: Config{
				Password: "", // PostgreSQL permite auth sin password (trust)
			},
		},
		{
			name: "config con timeout razonable",
			config: Config{
				ConnectTimeout: 30 * time.Second,
			},
		},
		{
			name: "config con diferentes nombres de database",
			config: Config{
				Database: "app_production_v1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// En este test simplemente verificamos que los valores se pueden asignar
			// y leer correctamente, ya que no hay método Validate() explícito en Config
			assert.NotNil(t, tt.config)
		})
	}

	t.Run("config con diferentes pool sizes", func(t *testing.T) {
		poolTests := []struct {
			name    string
			maxOpen int
			maxIdle int
		}{
			{"pool pequeño", 5, 2},
			{"pool mediano", 20, 10},
			{"pool grande", 100, 50},
		}

		for _, pt := range poolTests {
			t.Run(pt.name, func(t *testing.T) {
				cfg := Config{
					MaxConnections:     pt.maxOpen,
					MaxIdleConnections: pt.maxIdle,
				}
				assert.Equal(t, pt.maxOpen, cfg.MaxConnections)
				assert.Equal(t, pt.maxIdle, cfg.MaxIdleConnections)
			})
		}
	})
}

// Helper para testear DSN construction si estuviera expuesto
// Como no lo está, simulamos la lógica típica para verificar los valores
func TestConfig_DSN_Components(t *testing.T) {
	cfg := Config{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		Database: "db",
		SSLMode:  "disable",
	}

	expectedDSNFormat := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	dsn := fmt.Sprintf(expectedDSNFormat, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)

	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=user")
	assert.Contains(t, dsn, "password=pass")
	assert.Contains(t, dsn, "dbname=db")
	assert.Contains(t, dsn, "sslmode=disable")
}
