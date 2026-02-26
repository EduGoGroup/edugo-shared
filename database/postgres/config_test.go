package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, DefaultPort, cfg.Port)
	assert.Equal(t, "postgres", cfg.User)
	assert.Equal(t, "postgres", cfg.Database)
	assert.Equal(t, DefaultMaxConnections, cfg.MaxConnections)
	assert.Equal(t, DefaultMaxIdleConnections, cfg.MaxIdleConnections)
	assert.Equal(t, DefaultMaxLifetime, cfg.MaxLifetime)
	assert.Equal(t, "disable", cfg.SSLMode)
	assert.Equal(t, DefaultConnectTimeout, cfg.ConnectTimeout)
	assert.Contains(t, cfg.SearchPath, "public")
}

func TestConfig_Validation(t *testing.T) {
	// Add config validation tests if validation logic exists or create a method for it
	// Assuming no validation method exists yet, this is mostly checking struct integrity
}

// Since ConnectGORM is not refactored yet, we can't unit test the DSN building easily
// without duplicating the logic.
// However, we can add a test that ensures Config struct can hold all necessary values.

func TestConfig_CustomValues(t *testing.T) {
	cfg := Config{
		Host:           "127.0.0.1",
		Port:           5433,
		User:           "custom_user",
		Password:       "secure",
		Database:       "custom_db",
		SSLMode:        "require",
		MaxConnections: 100,
		ConnectTimeout: 5 * time.Second,
		SearchPath:     "myschema,public",
	}

	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, 5433, cfg.Port)
	assert.Equal(t, "custom_user", cfg.User)
	assert.Equal(t, "secure", cfg.Password)
	assert.Equal(t, "custom_db", cfg.Database)
	assert.Equal(t, "require", cfg.SSLMode)
	assert.Equal(t, 100, cfg.MaxConnections)
	assert.Equal(t, 5*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, "myschema,public", cfg.SearchPath)
}
