package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/common/config"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_GET_ENV"
	defaultValue := "default"

	t.Run("existing variable", func(t *testing.T) {
		expected := "value"
		os.Setenv(key, expected)
		defer os.Unsetenv(key)

		if got := config.GetEnv(key, defaultValue); got != expected {
			t.Errorf("GetEnv() = %v, want %v", got, expected)
		}
	})

	t.Run("non-existing variable", func(t *testing.T) {
		os.Unsetenv(key)
		if got := config.GetEnv(key, defaultValue); got != defaultValue {
			t.Errorf("GetEnv() = %v, want %v", got, defaultValue)
		}
	})

	t.Run("empty variable", func(t *testing.T) {
		os.Setenv(key, "")
		defer os.Unsetenv(key)
		if got := config.GetEnv(key, defaultValue); got != defaultValue {
			t.Errorf("GetEnv() = %v, want %v", got, defaultValue)
		}
	})
}

func TestGetEnvRequired(t *testing.T) {
	key := "TEST_GET_ENV_REQUIRED"

	t.Run("existing variable", func(t *testing.T) {
		expected := "value"
		os.Setenv(key, expected)
		defer os.Unsetenv(key)

		if got := config.GetEnvRequired(key); got != expected {
			t.Errorf("GetEnvRequired() = %v, want %v", got, expected)
		}
	})

	t.Run("non-existing variable", func(t *testing.T) {
		os.Unsetenv(key)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("GetEnvRequired() did not panic for missing variable")
			}
		}()
		config.GetEnvRequired(key)
	})

	t.Run("empty variable", func(t *testing.T) {
		os.Setenv(key, "")
		defer os.Unsetenv(key)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("GetEnvRequired() did not panic for empty variable")
			}
		}()
		config.GetEnvRequired(key)
	})
}

func TestGetEnvInt(t *testing.T) {
	key := "TEST_GET_ENV_INT"
	defaultValue := 42

	t.Run("valid integer", func(t *testing.T) {
		os.Setenv(key, "100")
		defer os.Unsetenv(key)
		if got := config.GetEnvInt(key, defaultValue); got != 100 {
			t.Errorf("GetEnvInt() = %v, want %v", got, 100)
		}
	})

	t.Run("invalid integer", func(t *testing.T) {
		os.Setenv(key, "not-an-int")
		defer os.Unsetenv(key)
		if got := config.GetEnvInt(key, defaultValue); got != defaultValue {
			t.Errorf("GetEnvInt() = %v, want %v", got, defaultValue)
		}
	})

	t.Run("non-existing variable", func(t *testing.T) {
		os.Unsetenv(key)
		if got := config.GetEnvInt(key, defaultValue); got != defaultValue {
			t.Errorf("GetEnvInt() = %v, want %v", got, defaultValue)
		}
	})
}

func TestGetEnvBool(t *testing.T) {
	key := "TEST_GET_ENV_BOOL"
	defaultValue := true

	t.Run("valid bool true", func(t *testing.T) {
		os.Setenv(key, "true")
		defer os.Unsetenv(key)
		if got := config.GetEnvBool(key, false); got != true {
			t.Errorf("GetEnvBool() = %v, want %v", got, true)
		}
	})

	t.Run("valid bool false", func(t *testing.T) {
		os.Setenv(key, "false")
		defer os.Unsetenv(key)
		if got := config.GetEnvBool(key, true); got != false {
			t.Errorf("GetEnvBool() = %v, want %v", got, false)
		}
	})

	t.Run("invalid bool", func(t *testing.T) {
		os.Setenv(key, "not-a-bool")
		defer os.Unsetenv(key)
		if got := config.GetEnvBool(key, defaultValue); got != defaultValue {
			t.Errorf("GetEnvBool() = %v, want %v", got, defaultValue)
		}
	})

	t.Run("non-existing variable", func(t *testing.T) {
		os.Unsetenv(key)
		if got := config.GetEnvBool(key, defaultValue); got != defaultValue {
			t.Errorf("GetEnvBool() = %v, want %v", got, defaultValue)
		}
	})
}

func TestGetEnvDuration(t *testing.T) {
	key := "TEST_GET_ENV_DURATION"
	defaultValue := 5 * time.Minute

	t.Run("valid duration", func(t *testing.T) {
		os.Setenv(key, "10s")
		defer os.Unsetenv(key)
		expected, _ := time.ParseDuration("10s")
		if got := config.GetEnvDuration(key, defaultValue); got != expected {
			t.Errorf("GetEnvDuration() = %v, want %v", got, expected)
		}
	})

	t.Run("invalid duration", func(t *testing.T) {
		os.Setenv(key, "not-a-duration")
		defer os.Unsetenv(key)
		if got := config.GetEnvDuration(key, defaultValue); got != defaultValue {
			t.Errorf("GetEnvDuration() = %v, want %v", got, defaultValue)
		}
	})

	t.Run("non-existing variable", func(t *testing.T) {
		os.Unsetenv(key)
		if got := config.GetEnvDuration(key, defaultValue); got != defaultValue {
			t.Errorf("GetEnvDuration() = %v, want %v", got, defaultValue)
		}
	})
}

func TestMustGetEnv(t *testing.T) {
	key := "TEST_MUST_GET_ENV"

	t.Run("existing variable", func(t *testing.T) {
		expected := "value"
		os.Setenv(key, expected)
		defer os.Unsetenv(key)

		if got := config.MustGetEnv(key); got != expected {
			t.Errorf("MustGetEnv() = %v, want %v", got, expected)
		}
	})

	t.Run("non-existing variable", func(t *testing.T) {
		os.Unsetenv(key)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustGetEnv() did not panic for missing variable")
			}
		}()
		config.MustGetEnv(key)
	})
}

func TestSetUnsetLookupEnv(t *testing.T) {
	key := "TEST_ENV_UTILS"
	value := "some-value"

	// Test SetEnv
	err := config.SetEnv(key, value)
	if err != nil {
		t.Errorf("SetEnv() error = %v", err)
	}

	// Test LookupEnv
	got, exists := config.LookupEnv(key)
	if !exists || got != value {
		t.Errorf("LookupEnv() = %v, %v, want %v, %v", got, exists, value, true)
	}

	// Test UnsetEnv
	err = config.UnsetEnv(key)
	if err != nil {
		t.Errorf("UnsetEnv() error = %v", err)
	}

	got, exists = config.LookupEnv(key)
	if exists {
		t.Errorf("LookupEnv() after UnsetEnv: exists = %v, want false", exists)
	}
}

func TestEnvironmentCheckers(t *testing.T) {
	key := "APP_ENV"
	originalValue, originalExists := os.LookupEnv(key)
	defer func() {
		if originalExists {
			os.Setenv(key, originalValue)
		} else {
			os.Unsetenv(key)
		}
	}()

	tests := []struct {
		env         string
		isDev       bool
		isProd      bool
		isStage     bool
		expectedEnv string
	}{
		{"development", true, false, false, "development"},
		{"dev", true, false, false, "dev"},
		{"local", true, false, false, "local"},
		{"production", false, true, false, "production"},
		{"prod", false, true, false, "prod"},
		{"staging", false, false, true, "staging"},
		{"stage", false, false, true, "stage"},
		{"qa", false, false, true, "qa"},
		{"unknown", false, false, false, "unknown"},
		{"", true, false, false, "development"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			if tt.env != "" {
				os.Setenv(key, tt.env)
			} else {
				os.Unsetenv(key)
			}

			if got := config.GetEnvironment(); got != tt.expectedEnv {
				t.Errorf("GetEnvironment() = %v, want %v", got, tt.expectedEnv)
			}
			if got := config.IsDevelopment(); got != tt.isDev {
				t.Errorf("IsDevelopment() = %v, want %v", got, tt.isDev)
			}
			if got := config.IsProduction(); got != tt.isProd {
				t.Errorf("IsProduction() = %v, want %v", got, tt.isProd)
			}
			if got := config.IsStaging(); got != tt.isStage {
				t.Errorf("IsStaging() = %v, want %v", got, tt.isStage)
			}
		})
	}
}
