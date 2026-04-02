package config

import (
	"testing"
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
