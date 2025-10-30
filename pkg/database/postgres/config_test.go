package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("retorna configuración con valores por defecto correctos", func(t *testing.T) {
		config := DefaultConfig()

		assert.Equal(t, "localhost", config.Host)
		assert.Equal(t, 5432, config.Port)
		assert.Equal(t, "postgres", config.User)
		assert.Equal(t, "postgres", config.Database)
		assert.Equal(t, 25, config.MaxConnections)
		assert.Equal(t, 5, config.MaxIdleConnections)
		assert.Equal(t, 5*time.Minute, config.MaxLifetime)
		assert.Equal(t, "disable", config.SSLMode)
		assert.Equal(t, 10*time.Second, config.ConnectTimeout)
	})

	t.Run("password vacío por defecto", func(t *testing.T) {
		config := DefaultConfig()

		assert.Empty(t, config.Password, "Password debe estar vacío por defecto")
	})
}

func TestConfigCustomization(t *testing.T) {
	t.Run("permite personalizar todos los campos", func(t *testing.T) {
		config := Config{
			Host:               "custom-host",
			Port:               5433,
			User:               "custom-user",
			Password:           "custom-password",
			Database:           "custom-db",
			MaxConnections:     50,
			MaxIdleConnections: 10,
			MaxLifetime:        10 * time.Minute,
			SSLMode:            "require",
			ConnectTimeout:     30 * time.Second,
		}

		assert.Equal(t, "custom-host", config.Host)
		assert.Equal(t, 5433, config.Port)
		assert.Equal(t, "custom-user", config.User)
		assert.Equal(t, "custom-password", config.Password)
		assert.Equal(t, "custom-db", config.Database)
		assert.Equal(t, 50, config.MaxConnections)
		assert.Equal(t, 10, config.MaxIdleConnections)
		assert.Equal(t, 10*time.Minute, config.MaxLifetime)
		assert.Equal(t, "require", config.SSLMode)
		assert.Equal(t, 30*time.Second, config.ConnectTimeout)
	})

	t.Run("puede modificar config por defecto", func(t *testing.T) {
		config := DefaultConfig()

		// Modificar valores
		config.Host = "production-host"
		config.Port = 5433
		config.User = "prod-user"
		config.Password = "secure-password"
		config.Database = "prod-db"
		config.MaxConnections = 100
		config.SSLMode = "verify-full"

		assert.Equal(t, "production-host", config.Host)
		assert.Equal(t, 5433, config.Port)
		assert.Equal(t, "prod-user", config.User)
		assert.Equal(t, "secure-password", config.Password)
		assert.Equal(t, "prod-db", config.Database)
		assert.Equal(t, 100, config.MaxConnections)
		assert.Equal(t, "verify-full", config.SSLMode)

		// Valores no modificados mantienen defaults
		assert.Equal(t, 5, config.MaxIdleConnections)
		assert.Equal(t, 5*time.Minute, config.MaxLifetime)
	})
}

func TestConfigSSLModes(t *testing.T) {
	validSSLModes := []string{
		"disable",
		"require",
		"verify-ca",
		"verify-full",
	}

	for _, sslMode := range validSSLModes {
		t.Run("permite SSL mode: "+sslMode, func(t *testing.T) {
			config := DefaultConfig()
			config.SSLMode = sslMode

			assert.Equal(t, sslMode, config.SSLMode)
		})
	}
}

func TestConfigConnectionPool(t *testing.T) {
	t.Run("MaxConnections debe ser mayor que MaxIdleConnections", func(t *testing.T) {
		config := DefaultConfig()

		assert.Greater(t, config.MaxConnections, config.MaxIdleConnections,
			"MaxConnections debe ser mayor que MaxIdleConnections para un pool eficiente")
	})

	t.Run("permite configurar pool pequeño", func(t *testing.T) {
		config := Config{
			Host:               "localhost",
			Port:               5432,
			User:               "test",
			Database:           "test",
			MaxConnections:     5,
			MaxIdleConnections: 2,
			MaxLifetime:        1 * time.Minute,
			SSLMode:            "disable",
			ConnectTimeout:     5 * time.Second,
		}

		assert.Equal(t, 5, config.MaxConnections)
		assert.Equal(t, 2, config.MaxIdleConnections)
	})

	t.Run("permite configurar pool grande para alta concurrencia", func(t *testing.T) {
		config := Config{
			Host:               "localhost",
			Port:               5432,
			User:               "test",
			Database:           "test",
			MaxConnections:     200,
			MaxIdleConnections: 50,
			MaxLifetime:        30 * time.Minute,
			SSLMode:            "disable",
			ConnectTimeout:     15 * time.Second,
		}

		assert.Equal(t, 200, config.MaxConnections)
		assert.Equal(t, 50, config.MaxIdleConnections)
		assert.Equal(t, 30*time.Minute, config.MaxLifetime)
	})
}

func TestConfigTimeouts(t *testing.T) {
	t.Run("permite configurar timeouts cortos", func(t *testing.T) {
		config := DefaultConfig()
		config.ConnectTimeout = 3 * time.Second

		assert.Equal(t, 3*time.Second, config.ConnectTimeout)
	})

	t.Run("permite configurar timeouts largos", func(t *testing.T) {
		config := DefaultConfig()
		config.ConnectTimeout = 60 * time.Second
		config.MaxLifetime = 1 * time.Hour

		assert.Equal(t, 60*time.Second, config.ConnectTimeout)
		assert.Equal(t, 1*time.Hour, config.MaxLifetime)
	})

	t.Run("timeout por defecto es razonable", func(t *testing.T) {
		config := DefaultConfig()

		assert.Equal(t, 10*time.Second, config.ConnectTimeout,
			"Timeout por defecto debe ser razonable (10 segundos)")
		assert.LessOrEqual(t, config.ConnectTimeout, 30*time.Second,
			"Timeout por defecto no debe ser excesivamente largo")
	})
}

func TestConfigForDifferentEnvironments(t *testing.T) {
	t.Run("configuración para desarrollo local", func(t *testing.T) {
		devConfig := Config{
			Host:               "localhost",
			Port:               5432,
			User:               "dev_user",
			Password:           "dev_password",
			Database:           "dev_db",
			MaxConnections:     10,
			MaxIdleConnections: 2,
			MaxLifetime:        5 * time.Minute,
			SSLMode:            "disable",
			ConnectTimeout:     5 * time.Second,
		}

		assert.Equal(t, "localhost", devConfig.Host)
		assert.Equal(t, "disable", devConfig.SSLMode)
		assert.Equal(t, 10, devConfig.MaxConnections, "Desarrollo usa menos conexiones")
	})

	t.Run("configuración para producción", func(t *testing.T) {
		prodConfig := Config{
			Host:               "prod-db.example.com",
			Port:               5432,
			User:               "prod_user",
			Password:           "secure_prod_password",
			Database:           "prod_db",
			MaxConnections:     100,
			MaxIdleConnections: 25,
			MaxLifetime:        30 * time.Minute,
			SSLMode:            "verify-full",
			ConnectTimeout:     15 * time.Second,
		}

		assert.NotEqual(t, "localhost", prodConfig.Host)
		assert.Equal(t, "verify-full", prodConfig.SSLMode, "Producción debe usar SSL")
		assert.GreaterOrEqual(t, prodConfig.MaxConnections, 50,
			"Producción necesita más conexiones")
	})

	t.Run("configuración para testing", func(t *testing.T) {
		testConfig := Config{
			Host:               "localhost",
			Port:               5433, // Puerto diferente
			User:               "test_user",
			Password:           "test_password",
			Database:           "test_db",
			MaxConnections:     5,
			MaxIdleConnections: 1,
			MaxLifetime:        1 * time.Minute,
			SSLMode:            "disable",
			ConnectTimeout:     3 * time.Second,
		}

		assert.Equal(t, "localhost", testConfig.Host)
		assert.Equal(t, 5433, testConfig.Port, "Testing usa puerto diferente")
		assert.Equal(t, 5, testConfig.MaxConnections, "Testing usa pocas conexiones")
	})
}
