package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestConnectGORM(t *testing.T) {
	// Pruebas básicas sobre un DSN erróneo o inexistente
	t.Run("Fallo de conexión por configuración inválida", func(t *testing.T) {
		cfg := &Config{
			Host:     "host-invalido-que-no-existe.local",
			Port:     5432,
			User:     "user",
			Password: "password",
			Database: "db",
			SSLMode:  "disable",
		}

		db, err := ConnectGORM(cfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("Conexión con configuraciones de timeout y search path", func(t *testing.T) {
		cfg := &Config{
			Host:           "localhost",
			Port:           5432,
			User:           "user",
			Password:       "password",
			Database:       "db",
			SSLMode:        "disable",
			SearchPath:     "public",
			ConnectTimeout: 10 * time.Second,
		}

		// Como no hay una BD real escuchando, esperamos que falle,
		// pero validamos que no lance pánico y maneje el error al construir el DSN
		db, err := ConnectGORM(cfg, logger.Default)
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}
