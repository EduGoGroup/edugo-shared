package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLoader_DefaultValues(t *testing.T) {
	loader := NewLoader()

	if len(loader.configPaths) == 0 || loader.configPaths[0] != "./config" {
		t.Errorf("configPaths = %v, want [./config]", loader.configPaths)
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

	found := false
	for _, p := range loader.configPaths {
		if p == "/custom/path" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("configPaths = %v, want it to contain /custom/path", loader.configPaths)
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

	found := false
	for _, p := range loader.configPaths {
		if p == "/test/path" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("configPaths = %v, want /test/path included", loader.configPaths)
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

	if err := os.WriteFile(configFile, []byte(yamlContent), 0600); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	loader := NewLoader(
		WithConfigPath(tmpDir),
		WithConfigName("config"),
	)

	var cfg struct {
		Environment string
		ServiceName string `mapstructure:"service_name"`
		Server      struct {
			Port int
		}
		Database struct {
			Host string
		}
		Logger struct {
			Level string
		}
	}
	err := loader.LoadFromFile(&cfg)

	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

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

	var cfg struct{}
	err := loader.LoadFromFile(&cfg)

	if err == nil {
		t.Error("Expected error when config file not found, got nil")
	}
}

func TestLoader_Load_WithEnvVars(t *testing.T) {
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

	if err := os.WriteFile(configFile, []byte(yamlContent), 0600); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	t.Setenv("APP_ENVIRONMENT", "prod")
	t.Setenv("APP_SERVER_PORT", "9090")

	loader := NewLoader(
		WithConfigPath(tmpDir),
		WithConfigName("config"),
		WithEnvPrefix("APP"),
	)

	var cfg struct {
		Environment string
		ServiceName string `mapstructure:"service_name"`
		Server      struct {
			Port int
		}
	}
	err := loader.Load(&cfg)

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Environment != "prod" {
		t.Errorf("Environment = %v, want prod (from env)", cfg.Environment)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %v, want 9090 (from env)", cfg.Server.Port)
	}
	if cfg.ServiceName != "test-service" {
		t.Errorf("ServiceName = %v, want test-service (from file)", cfg.ServiceName)
	}
}

func TestLoader_Load_FileNotFoundContinuesWithEnv(t *testing.T) {
	// AutomaticEnv solo resuelve keys que Viper ya conoce.
	// Sin archivo de config ni defaults, se deben usar WithExplicitBindings
	// para que Viper pueda leer env vars al hacer Unmarshal.
	t.Setenv("MY_ENVIRONMENT", "qa")
	t.Setenv("MY_SERVICE_NAME", "env-service")

	loader := NewLoader(
		WithConfigPath("/nonexistent"),
		WithConfigName("config"),
		WithExplicitBindings(map[string]string{
			"environment":  "MY_ENVIRONMENT",
			"service_name": "MY_SERVICE_NAME",
		}),
	)

	var cfg struct {
		Environment string `mapstructure:"environment"`
		ServiceName string `mapstructure:"service_name"`
	}

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
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
string_key: test_value
int_key: 42
bool_key: true
`
	if err := os.WriteFile(configFile, []byte(yamlContent), 0600); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	loader := NewLoader(
		WithConfigPath(tmpDir),
		WithConfigName("config"),
	)

	var cfg map[string]any
	if err := loader.Load(&cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

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

func TestLoader_GetMethods_BeforeLoad(t *testing.T) {
	loader := NewLoader()

	if val := loader.Get("key"); val != nil {
		t.Errorf("Get before Load = %v, want nil", val)
	}
	if val := loader.GetString("key"); val != "" {
		t.Errorf("GetString before Load = %v, want empty", val)
	}
	if val := loader.GetInt("key"); val != 0 {
		t.Errorf("GetInt before Load = %v, want 0", val)
	}
	if val := loader.GetBool("key"); val != false {
		t.Errorf("GetBool before Load = %v, want false", val)
	}
}

func TestLoader_MultipleOptions(t *testing.T) {
	loader := NewLoader(
		WithConfigPath("/path1"),
		WithConfigName("name1"),
		WithConfigPath("/path2"),
		WithConfigType("json"),
	)

	if len(loader.configPaths) < 2 {
		t.Fatalf("Expected at least 2 config paths, got %d", len(loader.configPaths))
	}

	found1, found2 := false, false
	for _, p := range loader.configPaths {
		if p == "/path1" {
			found1 = true
		}
		if p == "/path2" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("configPaths = %v, want /path1 and /path2", loader.configPaths)
	}

	if loader.configName != "name1" {
		t.Errorf("configName = %v, want name1", loader.configName)
	}
	if loader.configType != "json" {
		t.Errorf("configType = %v, want json", loader.configType)
	}
}

func TestLoader_WithDefaults(t *testing.T) {
	loader := NewLoader(
		WithConfigPath("/nonexistent"),
		WithDefaults(map[string]interface{}{
			"server.port": 8080,
			"log.level":   "info",
		}),
	)

	var cfg struct {
		Server struct {
			Port int `mapstructure:"port"`
		} `mapstructure:"server"`
		Log struct {
			Level string `mapstructure:"level"`
		} `mapstructure:"log"`
	}

	if err := loader.Load(&cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %v, want 8080 (from defaults)", cfg.Server.Port)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level = %v, want info (from defaults)", cfg.Log.Level)
	}
}

func TestLoader_WithExplicitBindings(t *testing.T) {
	t.Setenv("MY_SECRET_PASSWORD", "supersecret")

	loader := NewLoader(
		WithConfigPath("/nonexistent"),
		WithExplicitBindings(map[string]string{
			"database.password": "MY_SECRET_PASSWORD",
		}),
	)

	var cfg struct {
		Database struct {
			Password string `mapstructure:"password"`
		} `mapstructure:"database"`
	}

	if err := loader.Load(&cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Database.Password != "supersecret" {
		t.Errorf("Database.Password = %v, want supersecret (from explicit binding)", cfg.Database.Password)
	}
}

func TestLoader_LoadFromFile_WithDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `server:
  host: myhost
`
	if err := os.WriteFile(configFile, []byte(yamlContent), 0600); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	loader := NewLoader(
		WithConfigPath(tmpDir),
		WithDefaults(map[string]interface{}{
			"server.port": 9000,
		}),
	)

	var cfg struct {
		Server struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"server"`
	}

	if err := loader.LoadFromFile(&cfg); err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	if cfg.Server.Host != "myhost" {
		t.Errorf("Server.Host = %v, want myhost", cfg.Server.Host)
	}
	if cfg.Server.Port != 9000 {
		t.Errorf("Server.Port = %v, want 9000 (from defaults)", cfg.Server.Port)
	}
}

func TestLoader_ParallelSafety(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `environment: local`
	if err := os.WriteFile(configFile, []byte(yamlContent), 0600); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	t.Run("parallel1", func(t *testing.T) {
		t.Parallel()
		loader := NewLoader(WithConfigPath(tmpDir), WithEnvPrefix("P1"))
		var cfg struct{ Environment string }
		if err := loader.Load(&cfg); err != nil {
			t.Errorf("Load failed: %v", err)
		}
	})

	t.Run("parallel2", func(t *testing.T) {
		t.Parallel()
		loader := NewLoader(WithConfigPath(tmpDir), WithEnvPrefix("P2"))
		var cfg struct{ Environment string }
		if err := loader.Load(&cfg); err != nil {
			t.Errorf("Load failed: %v", err)
		}
	})

	t.Run("parallel3", func(t *testing.T) {
		t.Parallel()
		loader := NewLoader(WithConfigPath(tmpDir), WithEnvPrefix("P3"))
		var cfg struct{ Environment string }
		if err := loader.Load(&cfg); err != nil {
			t.Errorf("Load failed: %v", err)
		}
	})
}
