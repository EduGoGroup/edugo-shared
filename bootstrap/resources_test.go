package bootstrap

import (
	"testing"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func TestResources_HasLogger(t *testing.T) {
	t.Run("with_logger", func(t *testing.T) {
		r := &Resources{
			Logger: logrus.New(),
		}

		if !r.HasLogger() {
			t.Error("HasLogger() should return true when logger is set")
		}
	})

	t.Run("without_logger", func(t *testing.T) {
		r := &Resources{}

		if r.HasLogger() {
			t.Error("HasLogger() should return false when logger is nil")
		}
	})
}

func TestResources_HasPostgreSQL(t *testing.T) {
	t.Run("with_postgresql", func(t *testing.T) {
		r := &Resources{
			PostgreSQL: &gorm.DB{},
		}

		if !r.HasPostgreSQL() {
			t.Error("HasPostgreSQL() should return true when PostgreSQL is set")
		}
	})

	t.Run("without_postgresql", func(t *testing.T) {
		r := &Resources{}

		if r.HasPostgreSQL() {
			t.Error("HasPostgreSQL() should return false when PostgreSQL is nil")
		}
	})
}

func TestResources_HasMongoDB(t *testing.T) {
	t.Run("with_both_client_and_database", func(t *testing.T) {
		r := &Resources{
			MongoDB:       &mongo.Client{},
			MongoDatabase: &mongo.Database{},
		}

		if !r.HasMongoDB() {
			t.Error("HasMongoDB() should return true when both client and database are set")
		}
	})

	t.Run("with_only_client", func(t *testing.T) {
		r := &Resources{
			MongoDB: &mongo.Client{},
		}

		if r.HasMongoDB() {
			t.Error("HasMongoDB() should return false when database is nil")
		}
	})

	t.Run("with_only_database", func(t *testing.T) {
		r := &Resources{
			MongoDatabase: &mongo.Database{},
		}

		if r.HasMongoDB() {
			t.Error("HasMongoDB() should return false when client is nil")
		}
	})

	t.Run("with_neither", func(t *testing.T) {
		r := &Resources{}

		if r.HasMongoDB() {
			t.Error("HasMongoDB() should return false when both are nil")
		}
	})
}

func TestResources_HasMessagePublisher(t *testing.T) {
	t.Run("with_message_publisher", func(t *testing.T) {
		r := &Resources{
			MessagePublisher: &mockMessagePublisher{},
		}

		if !r.HasMessagePublisher() {
			t.Error("HasMessagePublisher() should return true when publisher is set")
		}
	})

	t.Run("without_message_publisher", func(t *testing.T) {
		r := &Resources{}

		if r.HasMessagePublisher() {
			t.Error("HasMessagePublisher() should return false when publisher is nil")
		}
	})
}

func TestResources_HasStorageClient(t *testing.T) {
	t.Run("with_storage_client", func(t *testing.T) {
		r := &Resources{
			StorageClient: &mockStorageClient{},
		}

		if !r.HasStorageClient() {
			t.Error("HasStorageClient() should return true when client is set")
		}
	})

	t.Run("without_storage_client", func(t *testing.T) {
		r := &Resources{}

		if r.HasStorageClient() {
			t.Error("HasStorageClient() should return false when client is nil")
		}
	})
}

func TestResources_AllResourcesPresent(t *testing.T) {
	r := &Resources{
		Logger:           logrus.New(),
		PostgreSQL:       &gorm.DB{},
		MongoDB:          &mongo.Client{},
		MongoDatabase:    &mongo.Database{},
		MessagePublisher: &mockMessagePublisher{},
		StorageClient:    &mockStorageClient{},
	}

	if !r.HasLogger() {
		t.Error("Logger should be present")
	}
	if !r.HasPostgreSQL() {
		t.Error("PostgreSQL should be present")
	}
	if !r.HasMongoDB() {
		t.Error("MongoDB should be present")
	}
	if !r.HasMessagePublisher() {
		t.Error("MessagePublisher should be present")
	}
	if !r.HasStorageClient() {
		t.Error("StorageClient should be present")
	}
}

func TestResources_NoResourcesPresent(t *testing.T) {
	r := &Resources{}

	if r.HasLogger() {
		t.Error("Logger should not be present")
	}
	if r.HasPostgreSQL() {
		t.Error("PostgreSQL should not be present")
	}
	if r.HasMongoDB() {
		t.Error("MongoDB should not be present")
	}
	if r.HasMessagePublisher() {
		t.Error("MessagePublisher should not be present")
	}
	if r.HasStorageClient() {
		t.Error("StorageClient should not be present")
	}
}

func TestResources_PartialConfiguration(t *testing.T) {
	// Test a typical configuration with only some resources
	r := &Resources{
		Logger:     logrus.New(),
		PostgreSQL: &gorm.DB{},
	}

	if !r.HasLogger() {
		t.Error("Logger should be present")
	}
	if !r.HasPostgreSQL() {
		t.Error("PostgreSQL should be present")
	}
	if r.HasMongoDB() {
		t.Error("MongoDB should not be present")
	}
	if r.HasMessagePublisher() {
		t.Error("MessagePublisher should not be present")
	}
	if r.HasStorageClient() {
		t.Error("StorageClient should not be present")
	}
}

// Mock implementations for testing

type mockMessagePublisher struct{}

func (m *mockMessagePublisher) Publish(exchange, routingKey string, body []byte) error {
	return nil
}

type mockStorageClient struct{}

func (m *mockStorageClient) Upload(bucket, key string, data []byte) error {
	return nil
}
