package config

import (
	"testing"
	"time"
)

type testDefaultsConfig struct {
	Environment string `mapstructure:"environment" default:"development"`
	Server      testDefaultsServer `mapstructure:"server"`
	NoTag       string
}

type testDefaultsServer struct {
	Port         int           `mapstructure:"port"          default:"8080"`
	Host         string        `mapstructure:"host"          default:"0.0.0.0"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"  default:"15s"`
	Debug        bool          `mapstructure:"debug"         default:"true"`
	SwaggerHost  string        `mapstructure:"swagger_host"`
}

func TestExtractDefaults_Basic(t *testing.T) {
	t.Parallel()

	defaults := ExtractDefaults(testDefaultsConfig{})

	tests := []struct {
		key  string
		want any
	}{
		{"environment", "development"},
		{"server.port", int64(8080)},
		{"server.host", "0.0.0.0"},
		{"server.read_timeout", 15 * time.Second},
		{"server.debug", true},
	}

	for _, tt := range tests {
		val, ok := defaults[tt.key]
		if !ok {
			t.Errorf("key %q not found in defaults", tt.key)
			continue
		}
		if val != tt.want {
			t.Errorf("defaults[%q] = %v (%T), want %v (%T)", tt.key, val, val, tt.want, tt.want)
		}
	}
}

func TestExtractDefaults_ZeroValueForMissingDefault(t *testing.T) {
	t.Parallel()

	defaults := ExtractDefaults(testDefaultsConfig{})

	// swagger_host no tiene default tag, debe registrarse con zero value
	val, ok := defaults["server.swagger_host"]
	if !ok {
		t.Fatal("key 'server.swagger_host' should be registered even without default tag")
	}
	if val != "" {
		t.Errorf("server.swagger_host = %v, want empty string", val)
	}
}

func TestExtractDefaults_SkipsFieldsWithoutMapstructure(t *testing.T) {
	t.Parallel()

	defaults := ExtractDefaults(testDefaultsConfig{})

	// NoTag no tiene mapstructure tag, no debe aparecer
	for key := range defaults {
		if key == "NoTag" || key == "notag" {
			t.Errorf("field without mapstructure tag should not appear in defaults, found %q", key)
		}
	}
}

func TestExtractDefaults_Pointer(t *testing.T) {
	t.Parallel()

	defaults := ExtractDefaults(&testDefaultsConfig{})

	if _, ok := defaults["environment"]; !ok {
		t.Error("ExtractDefaults should work with pointer to struct")
	}
}

func TestExtractDefaults_NestedKeyPaths(t *testing.T) {
	t.Parallel()

	type inner struct {
		Password string `mapstructure:"password"`
	}
	type mid struct {
		DB inner `mapstructure:"db"`
	}
	type root struct {
		Database mid `mapstructure:"database"`
	}

	defaults := ExtractDefaults(root{})

	if _, ok := defaults["database.db.password"]; !ok {
		t.Error("expected deeply nested key 'database.db.password' in defaults")
	}
}
