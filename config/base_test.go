package config

import (
	"testing"
	"time"
)

func TestDatabaseConfig_ConnectionString(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	actual := cfg.ConnectionString()

	if actual != expected {
		t.Errorf("ConnectionString() = %v, want %v", actual, expected)
	}
}

func TestDatabaseConfig_ConnectionStringWithDB(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Database: "defaultdb",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=customdb sslmode=disable"
	actual := cfg.ConnectionStringWithDB("customdb")

	if actual != expected {
		t.Errorf("ConnectionStringWithDB() = %v, want %v", actual, expected)
	}
}

func TestBaseConfig_DefaultValues(t *testing.T) {
	cfg := BaseConfig{
		Environment: "local",
		ServiceName: "test-service",
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "pass",
			Database: "db",
			SSLMode:  "disable",
		},
		MongoDB: MongoDBConfig{
			URI:      "mongodb://localhost:27017",
			Database: "testdb",
		},
		Logger: LoggerConfig{
			Level:  "info",
			Format: "json",
		},
	}

	if cfg.Environment != "local" {
		t.Errorf("Environment = %v, want local", cfg.Environment)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %v, want 8080", cfg.Server.Port)
	}
	if cfg.Logger.Level != "info" {
		t.Errorf("Logger.Level = %v, want info", cfg.Logger.Level)
	}
}
