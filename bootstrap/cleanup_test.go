package bootstrap

import (
	"testing"

	"github.com/EduGoGroup/edugo-shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

// Mock lifecycle manager
type mockLifecycleManager struct {
	registered map[string]func() error
}

func (m *mockLifecycleManager) RegisterSimple(name string, cleanup func() error) {
	m.registered[name] = cleanup
}

func TestRegisterPostgreSQLCleanup(t *testing.T) {
	t.Run("with nil lifecycle manager", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())

		assert.NotPanics(t, func() {
			registerPostgreSQLCleanup(nil, nil, nil, logger)
		})
	})

	t.Run("with lifecycle manager that doesn't implement RegisterSimple", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := struct{}{}

		assert.NotPanics(t, func() {
			registerPostgreSQLCleanup(lifecycleManager, nil, nil, logger)
		})
	})

	t.Run("with nil factory", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}

		registerPostgreSQLCleanup(lifecycleManager, nil, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with nil db", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockPostgreSQLFactory{}

		registerPostgreSQLCleanup(lifecycleManager, factory, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with wrong type of db", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockPostgreSQLFactory{}
		wrongDB := "not a *gorm.DB"

		registerPostgreSQLCleanup(lifecycleManager, factory, wrongDB, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("successful registration", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockPostgreSQLFactory{}
		db := &gorm.DB{}

		registerPostgreSQLCleanup(lifecycleManager, factory, db, logger)

		assert.Contains(t, lifecycleManager.registered, "postgresql")
		assert.NotNil(t, lifecycleManager.registered["postgresql"])
	})
}

func TestRegisterMongoDBCleanup(t *testing.T) {
	t.Run("with nil lifecycle manager", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())

		assert.NotPanics(t, func() {
			registerMongoDBCleanup(nil, nil, nil, logger)
		})
	})

	t.Run("with nil factory", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}

		registerMongoDBCleanup(lifecycleManager, nil, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with nil client", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockMongoDBFactory{}

		registerMongoDBCleanup(lifecycleManager, factory, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with wrong type of client", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockMongoDBFactory{}
		wrongClient := "not a *mongo.Client"

		registerMongoDBCleanup(lifecycleManager, factory, wrongClient, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("successful registration", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockMongoDBFactory{}
		client := &mongo.Client{}

		registerMongoDBCleanup(lifecycleManager, factory, client, logger)

		assert.Contains(t, lifecycleManager.registered, "mongodb")
		assert.NotNil(t, lifecycleManager.registered["mongodb"])
	})
}

func TestRegisterRabbitMQCleanup(t *testing.T) {
	t.Run("with nil lifecycle manager", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())

		assert.NotPanics(t, func() {
			registerRabbitMQCleanup(nil, nil, nil, nil, logger)
		})
	})

	t.Run("with nil factory", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}

		registerRabbitMQCleanup(lifecycleManager, nil, nil, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with nil channel", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockRabbitMQFactory{}

		registerRabbitMQCleanup(lifecycleManager, factory, nil, nil, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with wrong type of channel", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockRabbitMQFactory{}
		wrongChannel := "not a *amqp.Channel"
		conn := &amqp.Connection{}

		registerRabbitMQCleanup(lifecycleManager, factory, wrongChannel, conn, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("with wrong type of connection", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockRabbitMQFactory{}
		channel := &amqp.Channel{}
		wrongConn := "not a *amqp.Connection"

		registerRabbitMQCleanup(lifecycleManager, factory, channel, wrongConn, logger)

		assert.Empty(t, lifecycleManager.registered)
	})

	t.Run("successful registration", func(t *testing.T) {
		logger := logger.NewLogrusLogger(logrus.New())
		lifecycleManager := &mockLifecycleManager{registered: make(map[string]func() error)}
		factory := &mockRabbitMQFactory{}
		channel := &amqp.Channel{}
		conn := &amqp.Connection{}

		registerRabbitMQCleanup(lifecycleManager, factory, channel, conn, logger)

		assert.Contains(t, lifecycleManager.registered, "rabbitmq")
		assert.NotNil(t, lifecycleManager.registered["rabbitmq"])
	})
}

func TestExtractPostgreSQLConfig(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		_, err := extractPostgreSQLConfig(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with non-struct config", func(t *testing.T) {
		_, err := extractPostgreSQLConfig("not a struct")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with struct without PostgreSQL field", func(t *testing.T) {
		config := struct {
			OtherField string
		}{
			OtherField: "value",
		}

		_, err := extractPostgreSQLConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PostgreSQL field not found")
	})
}

func TestExtractMongoDBConfig(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		_, err := extractMongoDBConfig(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with non-struct config", func(t *testing.T) {
		_, err := extractMongoDBConfig(123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})
}

func TestExtractRabbitMQConfig(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		_, err := extractRabbitMQConfig(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with non-struct config", func(t *testing.T) {
		_, err := extractRabbitMQConfig([]string{"not", "a", "struct"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})
}

func TestExtractS3Config(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		_, err := extractS3Config(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})

	t.Run("with non-struct config", func(t *testing.T) {
		_, err := extractS3Config(make(map[string]string))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config must be a struct")
	})
}
