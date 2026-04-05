package config

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

type testServerConfig struct {
	Port int    `validate:"required,min=1,max=65535"`
	Host string `validate:"required"`
}

type testAppConfig struct {
	Environment string           `validate:"required,oneof=local dev qa prod"`
	ServiceName string           `validate:"required"`
	Server      testServerConfig `validate:"required"`
}

func TestValidator_Validate_Success(t *testing.T) {
	v := NewValidator()

	cfg := testAppConfig{
		Environment: "local",
		ServiceName: "test-service",
		Server: testServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	}

	err := v.Validate(&cfg)
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestValidator_Validate_MissingRequired(t *testing.T) {
	v := NewValidator()

	cfg := testAppConfig{
		Environment: "local",
		// ServiceName missing - should fail
	}

	err := v.Validate(&cfg)
	if err == nil {
		t.Error("Validate() error = nil, want error for missing ServiceName")
	}

	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		if len(validationErr.Errors) == 0 {
			t.Error("ValidationError.Errors is empty, expected errors")
		}

		// Verifica que el path incluye la ruta completa sin el struct raíz
		found := false
		for _, fe := range validationErr.Errors {
			if fe.Field == "ServiceName" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected field path 'ServiceName' in validation errors")
		}
	}
}

func TestValidator_Validate_InvalidEnvironment(t *testing.T) {
	v := NewValidator()

	cfg := testAppConfig{
		Environment: "staging", // not in oneof=local dev qa prod
		ServiceName: "test",
		Server: testServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	}

	err := v.Validate(&cfg)
	if err == nil {
		t.Error("Validate() error = nil, want error for invalid environment")
	}
}

func TestValidator_Validate_PortOutOfRange(t *testing.T) {
	v := NewValidator()

	cfg := testAppConfig{
		Environment: "local",
		ServiceName: "test",
		Server: testServerConfig{
			Port: 99999, // exceeds max=65535
			Host: "localhost",
		},
	}

	err := v.Validate(&cfg)
	if err == nil {
		t.Error("Validate() error = nil, want error for port out of range")
	}
}

func TestValidationError_Error(t *testing.T) {
	v := NewValidator()

	cfg := testAppConfig{
		Environment: "invalid",
		// ServiceName missing
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

func TestValidator_ValidateField(t *testing.T) {
	v := NewValidator()

	t.Run("valid email", func(t *testing.T) {
		err := v.ValidateField("user@example.com", "email")
		if err != nil {
			t.Errorf("ValidateField() error = %v, want nil", err)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		err := v.ValidateField("not-an-email", "email")
		if err == nil {
			t.Error("ValidateField() error = nil, want error for invalid email")
		}
	})

	t.Run("required non-empty", func(t *testing.T) {
		err := v.ValidateField("somevalue", "required")
		if err != nil {
			t.Errorf("ValidateField() error = %v, want nil", err)
		}
	})

	t.Run("required empty fails", func(t *testing.T) {
		err := v.ValidateField("", "required")
		if err == nil {
			t.Error("ValidateField() error = nil, want error for empty required field")
		}
	})
}

func TestValidator_Validate_NestedFieldPath(t *testing.T) {
	v := NewValidator()

	cfg := testAppConfig{
		Environment: "local",
		ServiceName: "test",
		Server: testServerConfig{
			Port: 0, // fails required
			Host: "", // fails required
		},
	}

	err := v.Validate(&cfg)
	if err == nil {
		t.Fatal("expected validation error for nested fields")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}

	paths := make(map[string]bool)
	for _, fe := range validationErr.Errors {
		paths[fe.Field] = true
	}

	// Debe incluir el path completo: Server.Port, Server.Host
	if !paths["Server.Port"] {
		t.Errorf("expected field path 'Server.Port', got paths: %v", paths)
	}
	if !paths["Server.Host"] {
		t.Errorf("expected field path 'Server.Host', got paths: %v", paths)
	}
}

func TestValidator_RegisterValidation(t *testing.T) {
	v := NewValidator()

	// Registra un validador custom que verifica que el string empiece con "svc-"
	err := v.RegisterValidation("svc_prefix", func(fl validator.FieldLevel) bool {
		return strings.HasPrefix(fl.Field().String(), "svc-")
	})
	if err != nil {
		t.Fatalf("RegisterValidation failed: %v", err)
	}

	type svcConfig struct {
		Name string `validate:"required,svc_prefix"`
	}

	t.Run("valid", func(t *testing.T) {
		if err := v.Validate(&svcConfig{Name: "svc-identity"}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		if err := v.Validate(&svcConfig{Name: "identity"}); err == nil {
			t.Error("expected error for missing svc- prefix")
		}
	})
}

func TestValidationError_MultipleErrors(t *testing.T) {
	v := NewValidator()

	// cfg con múltiples campos inválidos
	cfg := testAppConfig{
		Environment: "bad-env",
		// ServiceName missing
		// Server missing (Port=0, Host="")
	}

	err := v.Validate(&cfg)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}

	if len(validationErr.Errors) < 2 {
		t.Errorf("expected at least 2 field errors, got %d", len(validationErr.Errors))
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error() returned empty string")
	}
}
