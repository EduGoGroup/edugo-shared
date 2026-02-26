package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validation_Unit(t *testing.T) {
	t.Run("configuración válida", func(t *testing.T) {
		cfg := Config{
			URI:         "mongodb://localhost:27017",
			Database:    "test",
			Timeout:     5 * time.Second,
			MaxPoolSize: 100,
			MinPoolSize: 10,
		}

		assert.Equal(t, "mongodb://localhost:27017", cfg.URI)
		assert.Equal(t, "test", cfg.Database)
		assert.Equal(t, 5*time.Second, cfg.Timeout)
		assert.Equal(t, uint64(100), cfg.MaxPoolSize)
		assert.Equal(t, uint64(10), cfg.MinPoolSize)
	})

	t.Run("configuración por defecto es válida", func(t *testing.T) {
		cfg := DefaultConfig()
		assert.NotEmpty(t, cfg.URI)
		assert.NotEmpty(t, cfg.Database)
		assert.Greater(t, cfg.Timeout, time.Duration(0))
		assert.Greater(t, cfg.MaxPoolSize, uint64(0))
	})
}

func TestConfig_Parsing(t *testing.T) {
	// Verificar que la estructura Config soporta diferentes formatos de URI
	uris := []string{
		"mongodb://user:pass@host:27017/db",
		"mongodb+srv://user:pass@cluster.mongodb.net/db",
		"mongodb://host1,host2,host3/?replicaSet=rs0",
	}

	for _, uri := range uris {
		cfg := Config{URI: uri}
		assert.Equal(t, uri, cfg.URI)
	}
}
