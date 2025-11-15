package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestNewLoader_DefaultValues(t *testing.T) {
	loader := NewLoader()

	if loader.configPath != "./config" {
		t.Errorf("configPath = %v, want ./config", loader.configPath)
	}
	if loader.configName != "config" {
		t.Errorf("configName = %v, want config", loader.configName)
	}
	if loader.configType != "yaml" {
		t.Errorf("configType = %v, want yaml", loader.configType)
	}
	if loader.envPrefix != "" {
		t.Errorf("envPrefix = %v, want empty string", loader.envPrefix)
	}
}

func TestNewLoader_WithOptions(t *testing.T) {
	loader := NewLoader(
		WithConfigPath("/custom/path"),
		WithConfigName("custom"),
		WithConfigType("json"),
		WithEnvPrefix("APP"),
	)

	if loader.configPath != "/custom/path" {
		t.Errorf("configPath = %v, want /custom/path", loader.configPath)
	}
	if loader.configName != "custom" {
		t.Errorf("configName = %v, want custom", loader.configName)
	}
	if loader.configType != "json" {
		t.Errorf("configType = %v, want json", loader.configType)
	}
	if loader.envPrefix != "APP" {
		t.Errorf("envPrefix = %v, want APP", loader.envPrefix)
	}
}

func TestWithConfigPath(t *testing.T) {
	loader := NewLoader(WithConfigPath("/test/path"))

	if loader.configPath != "/test/path" {
		t.Errorf("configPath = %v, want /test/path", loader.configPath)
	}
}

func TestWithConfigName(t *testing.T) {
	loader := NewLoader(WithConfigName("test-config"))

	if loader.configName != "test-config" {
		t.Errorf("configName = %v, want test-config", loader.configName)
	}
}

func TestWithConfigType(t *testing.T) {
	loader := NewLoader(WithConfigType("toml"))

	if loader.configType != "toml" {
		t.Errorf("configType = %v, want toml", loader.configType)
	}
}

func TestWithEnvPrefix(t *testing.T) {
	loader := NewLoader(WithEnvPrefix("TEST"))

	if loader.envPrefix != "TEST" {
		t.Errorf("envPrefix = %v, want TEST", loader.envPrefix)
	}
}

func TestLoader_LoadFromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
environment: local
service_name: test-service
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 60s
database:
  host: localhost
  port: 5432
  user: testuser
  password: testpass
  database: testdb
  ssl_mode: disable
mongodb:
  uri: mongodb://localhost:27017
  database: testdb
logger:
  level: info
  format: json
`

	if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading
	loader := NewLoader(
		WithConfigPath(tmpDir),
		WithConfigName("config"),
	)

	var cfg BaseConfig
	err := loader.LoadFromFile(&cfg)

	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	// Verify loaded values
	if cfg.Environment != "local" {
		t.Errorf("Environment = %v, want local", cfg.Environment)
	}
	if cfg.ServiceName != "test-service" {
		t.Errorf("ServiceName = %v, want test-service", cfg.ServiceName)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %v, want 8080", cfg.Server.Port)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Database.Host = %v, want localhost", cfg.Database.Host)
	}
	if cfg.Logger.Level != "info" {
		t.Errorf("Logger.Level = %v, want info", cfg.Logger.Level)
	}
}

func TestLoader_LoadFromFile_FileNotFound(t *testing.T) {
	loader := NewLoader(
		WithConfigPath("/nonexistent/path"),
		WithConfigName("config"),
	)

	var cfg BaseConfig
	err := loader.LoadFromFile(&cfg)

	if err == nil {
		t.Error("Expected error when config file not found, got nil")
	}
}

func TestLoader_Load_WithEnvVars(t *testing.T) {
	// Reset viper for this test
	viper.Reset()

	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
environment: local
service_name: test-service
server:
  port: 8080
database:
  host: localhost
  port: 5432
  user: testuser
  password: testpass
  database: testdb
  ssl_mode: disable
mongodb:
  uri: mongodb://localhost:27017
  database: testdb
logger:
  level: info
  format: json
`

	if err := os.WriteFile(configFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Set environment variables
	os.Setenv("APP_ENVIRONMENT", "prod")
	os.Setenv("APP_SERVER_PORT", "9090")
	defer func() {
		os.Unsetenv("APP_ENVIRONMENT")
		os.Unsetenv("APP_SERVER_PORT")
		viper.Reset()
	}()

	// Test loading with env vars
	loader := NewLoader(
		WithConfigPath(tmpDir),
		WithConfigName("config"),
		WithEnvPrefix("APP"),
	)

	var cfg BaseConfig
	err := loader.Load(&cfg)

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Env vars should override file values
	if cfg.Environment != "prod" {
		t.Errorf("Environment = %v, want prod (from env)", cfg.Environment)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %v, want 9090 (from env)", cfg.Server.Port)
	}
	// Values not overridden should come from file
	if cfg.ServiceName != "test-service" {
		t.Errorf("ServiceName = %v, want test-service (from file)", cfg.ServiceName)
	}
}

func TestLoader_Load_FileNotFoundContinuesWithEnv(t *testing.T) {
	// Reset viper for this test
	viper.Reset()

	// Set environment variables - Viper usa mayúsculas para todo el key
	os.Setenv("APP_ENVIRONMENT", "qa")
	os.Setenv("APP_SERVICE_NAME", "env-service")
	defer func() {
		os.Unsetenv("APP_ENVIRONMENT")
		os.Unsetenv("APP_SERVICE_NAME")
		viper.Reset()
	}()

	// Load with non-existent file should still work with env vars
	loader := NewLoader(
		WithConfigPath("/nonexistent"),
		WithConfigName("config"),
		WithEnvPrefix("APP"),
	)

	var cfg struct {
		Environment string `mapstructure:"environment"`
		ServiceName string `mapstructure:"service_name"`
	}

	// Bind env vars before loading (needed for viper to recognize them)
	viper.BindEnv("environment", "APP_ENVIRONMENT")
	viper.BindEnv("service_name", "APP_SERVICE_NAME")

	err := loader.Load(&cfg)

	if err != nil {
		t.Fatalf("Load should not fail when file not found but env vars present: %v", err)
	}

	if cfg.Environment != "qa" {
		t.Errorf("Environment = %v, want qa", cfg.Environment)
	}
	if cfg.ServiceName != "env-service" {
		t.Errorf("ServiceName = %v, want env-service", cfg.ServiceName)
	}
}

func TestLoader_GetMethods(t *testing.T) {
	// Reset viper and set test values
	viper.Reset()
	defer viper.Reset()

	viper.Set("string_key", "test_value")
	viper.Set("int_key", 42)
	viper.Set("bool_key", true)

	loader := NewLoader()

	t.Run("Get", func(t *testing.T) {
		val := loader.Get("string_key")
		if val != "test_value" {
			t.Errorf("Get('string_key') = %v, want test_value", val)
		}
	})

	t.Run("GetString", func(t *testing.T) {
		val := loader.GetString("string_key")
		if val != "test_value" {
			t.Errorf("GetString('string_key') = %v, want test_value", val)
		}
	})

	t.Run("GetInt", func(t *testing.T) {
		val := loader.GetInt("int_key")
		if val != 42 {
			t.Errorf("GetInt('int_key') = %v, want 42", val)
		}
	})

	t.Run("GetBool", func(t *testing.T) {
		val := loader.GetBool("bool_key")
		if !val {
			t.Error("GetBool('bool_key') = false, want true")
		}
	})

	t.Run("Get_NonExistentKey", func(t *testing.T) {
		val := loader.Get("nonexistent")
		if val != nil {
			t.Errorf("Get('nonexistent') = %v, want nil", val)
		}
	})
}

func TestLoader_MultipleOptions(t *testing.T) {
	loader := NewLoader(
		WithConfigPath("/path1"),
		WithConfigName("name1"),
		WithConfigPath("/path2"), // Should override previous
		WithConfigType("json"),
	)

	// Last option should win
	if loader.configPath != "/path2" {
		t.Errorf("configPath = %v, want /path2", loader.configPath)
	}
	if loader.configName != "name1" {
		t.Errorf("configName = %v, want name1", loader.configName)
	}
	if loader.configType != "json" {
		t.Errorf("configType = %v, want json", loader.configType)
	}
}
