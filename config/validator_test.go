package config

import (
	"testing"
	"time"
)

func TestValidator_Validate_Success(t *testing.T) {
	v := NewValidator()

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
			Host:            "localhost",
			Port:            5432,
			User:            "user",
			Password:        "pass",
			Database:        "db",
			SSLMode:         "disable",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
			ConnMaxIdleTime: 2 * time.Minute,
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

	err := v.Validate(&cfg)
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestValidator_Validate_MissingRequired(t *testing.T) {
	v := NewValidator()

	cfg := BaseConfig{
		Environment: "local",
		// ServiceName missing - should fail
	}

	err := v.Validate(&cfg)
	if err == nil {
		t.Error("Validate() error = nil, want error for missing ServiceName")
	}

	if validationErr, ok := err.(*ValidationError); ok {
		if len(validationErr.Errors) == 0 {
			t.Error("ValidationError.Errors is empty, expected errors")
		}
	}
}

func TestValidator_Validate_InvalidEnvironment(t *testing.T) {
	v := NewValidator()

	cfg := BaseConfig{
		Environment: "invalid",
		ServiceName: "test",
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: DatabaseConfig{
			Host:         "localhost",
			Port:         5432,
			User:         "user",
			Password:     "pass",
			Database:     "db",
			SSLMode:      "disable",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
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

	err := v.Validate(&cfg)
	if err == nil {
		t.Error("Validate() error = nil, want error for invalid environment")
	}
}

func TestValidationError_Error(t *testing.T) {
	v := NewValidator()

	cfg := BaseConfig{
		Environment: "invalid",
	}

	err := v.Validate(&cfg)
	if err == nil {
		t.Fatal("expected validation error")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error() returned empty string")
	}
}
