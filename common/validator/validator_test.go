package validator_test

import (
	"testing"

	"github.com/EduGoGroup/edugo-shared/common/validator"
)

func TestNew(t *testing.T) {
	v := validator.New()
	if v == nil {
		t.Fatal("New() returned nil")
	}
	if v.HasErrors() {
		t.Error("New validator should have no errors")
	}
}

func TestAddError(t *testing.T) {
	v := validator.New()
	v.AddError("error 1")
	v.AddError("error 2")

	if !v.HasErrors() {
		t.Error("Should have errors")
	}

	errors := v.GetErrors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}
}

func TestAddErrorf(t *testing.T) {
	v := validator.New()
	v.AddErrorf("error %d: %s", 1, "test")

	errors := v.GetErrors()
	if len(errors) != 1 {
		t.Fatal("Expected 1 error")
	}
	if errors[0] != "error 1: test" {
		t.Errorf("Unexpected error message: %s", errors[0])
	}
}

func TestGetError(t *testing.T) {
	t.Run("with_errors", func(t *testing.T) {
		v := validator.New()
		v.AddError("error 1")
		v.AddError("error 2")

		err := v.GetError()
		if err == nil {
			t.Error("Should return error when has errors")
		}
	})

	t.Run("without_errors", func(t *testing.T) {
		v := validator.New()
		err := v.GetError()
		if err != nil {
			t.Error("Should return nil when no errors")
		}
	})
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		hasError bool
	}{
		{"valid_value", "test", false},
		{"empty_string", "", true},
		{"whitespace_only", "   ", true},
		{"tabs_and_spaces", " \t ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			v.Required(tt.value, "field")

			if v.HasErrors() != tt.hasError {
				t.Errorf("Expected hasError=%v for value '%s'", tt.hasError, tt.value)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	v := validator.New()
	v.MinLength("ab", 3, "field")

	if !v.HasErrors() {
		t.Error("Should have error for string shorter than min length")
	}

	v2 := validator.New()
	v2.MinLength("abc", 3, "field")

	if v2.HasErrors() {
		t.Error("Should not have error for string >= min length")
	}
}

func TestMaxLength(t *testing.T) {
	v := validator.New()
	v.MaxLength("abcde", 3, "field")

	if !v.HasErrors() {
		t.Error("Should have error for string longer than max length")
	}

	v2 := validator.New()
	v2.MaxLength("abc", 5, "field")

	if v2.HasErrors() {
		t.Error("Should not have error for string <= max length")
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		hasError bool
	}{
		{"valid_email", "test@example.com", false},
		{"valid_with_subdomain", "user@mail.example.com", false},
		{"valid_with_plus", "user+tag@example.com", false},
		{"invalid_no_at", "testexample.com", true},
		{"invalid_no_domain", "test@", true},
		{"invalid_no_tld", "test@example", true},
		{"empty_string", "", false}, // Empty is allowed, use Required for mandatory
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			v.Email(tt.email, "email")

			if v.HasErrors() != tt.hasError {
				t.Errorf("Expected hasError=%v for email '%s'", tt.hasError, tt.email)
			}
		})
	}
}

func TestUUID(t *testing.T) {
	tests := []struct {
		name     string
		uuid     string
		hasError bool
	}{
		{"valid_uuid", "123e4567-e89b-12d3-a456-426614174000", false},
		{"invalid_uuid", "not-a-uuid", true},
		{"invalid_format", "123e4567-e89b-12d3", true},
		{"empty_string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			v.UUID(tt.uuid, "id")

			if v.HasErrors() != tt.hasError {
				t.Errorf("Expected hasError=%v for uuid '%s'", tt.hasError, tt.uuid)
			}
		})
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		hasError bool
	}{
		{"valid_http", "http://example.com", false},
		{"valid_https", "https://example.com", false},
		{"valid_with_path", "https://example.com/path/to/resource", false},
		{"invalid_no_protocol", "example.com", true},
		{"invalid_ftp", "ftp://example.com", true},
		{"empty_string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			v.URL(tt.url, "url")

			if v.HasErrors() != tt.hasError {
				t.Errorf("Expected hasError=%v for url '%s'", tt.hasError, tt.url)
			}
		})
	}
}

func TestInSlice(t *testing.T) {
	allowed := []string{"admin", "user", "guest"}

	t.Run("valid_value", func(t *testing.T) {
		v := validator.New()
		v.InSlice("admin", allowed, "role")
		if v.HasErrors() {
			t.Error("Should not have error for valid value")
		}
	})

	t.Run("invalid_value", func(t *testing.T) {
		v := validator.New()
		v.InSlice("superuser", allowed, "role")
		if !v.HasErrors() {
			t.Error("Should have error for invalid value")
		}
	})

	t.Run("empty_value", func(t *testing.T) {
		v := validator.New()
		v.InSlice("", allowed, "role")
		if v.HasErrors() {
			t.Error("Should not have error for empty value")
		}
	})
}

func TestMinValue(t *testing.T) {
	v := validator.New()
	v.MinValue(5, 10, "age")

	if !v.HasErrors() {
		t.Error("Should have error for value < min")
	}

	v2 := validator.New()
	v2.MinValue(10, 10, "age")

	if v2.HasErrors() {
		t.Error("Should not have error for value >= min")
	}
}

func TestMaxValue(t *testing.T) {
	v := validator.New()
	v.MaxValue(15, 10, "age")

	if !v.HasErrors() {
		t.Error("Should have error for value > max")
	}

	v2 := validator.New()
	v2.MaxValue(10, 10, "age")

	if v2.HasErrors() {
		t.Error("Should not have error for value <= max")
	}
}

func TestRange(t *testing.T) {
	t.Run("below_range", func(t *testing.T) {
		v := validator.New()
		v.Range(5, 10, 20, "score")
		if !v.HasErrors() {
			t.Error("Should have error for value below range")
		}
	})

	t.Run("above_range", func(t *testing.T) {
		v := validator.New()
		v.Range(25, 10, 20, "score")
		if !v.HasErrors() {
			t.Error("Should have error for value above range")
		}
	})

	t.Run("within_range", func(t *testing.T) {
		v := validator.New()
		v.Range(15, 10, 20, "score")
		if v.HasErrors() {
			t.Error("Should not have error for value within range")
		}
	})

	t.Run("at_min_boundary", func(t *testing.T) {
		v := validator.New()
		v.Range(10, 10, 20, "score")
		if v.HasErrors() {
			t.Error("Should not have error at min boundary")
		}
	})

	t.Run("at_max_boundary", func(t *testing.T) {
		v := validator.New()
		v.Range(20, 10, 20, "score")
		if v.HasErrors() {
			t.Error("Should not have error at max boundary")
		}
	})
}

func TestName(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		hasError bool
	}{
		{"valid_name", "John Doe", false},
		{"valid_with_accent", "José García", false},
		{"valid_with_apostrophe", "O'Brien", false},
		{"valid_with_hyphen", "Mary-Jane", false},
		{"invalid_with_numbers", "John123", true},
		{"invalid_with_symbols", "John@Doe", true},
		{"empty_string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			v.Name(tt.value, "name")

			if v.HasErrors() != tt.hasError {
				t.Errorf("Expected hasError=%v for name '%s'", tt.hasError, tt.value)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user+tag@example.co.uk", true},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
		{"a@b.co", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := validator.IsValidEmail(tt.email)
			if result != tt.valid {
				t.Errorf("IsValidEmail('%s') = %v, want %v", tt.email, result, tt.valid)
			}
		})
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		uuid  string
		valid bool
	}{
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"not-a-uuid", false},
		{"", false},
		{"123", false},
	}

	for _, tt := range tests {
		t.Run(tt.uuid, func(t *testing.T) {
			result := validator.IsValidUUID(tt.uuid)
			if result != tt.valid {
				t.Errorf("IsValidUUID('%s') = %v, want %v", tt.uuid, result, tt.valid)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		url   string
		valid bool
	}{
		{"http://example.com", true},
		{"https://example.com/path", true},
		{"https://subdomain.example.com", true},
		{"example.com", false},
		{"ftp://example.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := validator.IsValidURL(tt.url)
			if result != tt.valid {
				t.Errorf("IsValidURL('%s') = %v, want %v", tt.url, result, tt.valid)
			}
		})
	}
}

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"John Doe", true},
		{"Mary-Jane", true},
		{"O'Brien", true},
		{"José", true},
		{"John123", false},
		{"", false},
		{"a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.IsValidName(tt.name)
			if result != tt.valid {
				t.Errorf("IsValidName('%s') = %v, want %v", tt.name, result, tt.valid)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		value string
		empty bool
	}{
		{"", true},
		{"  ", true},
		{"\t", true},
		{" \n ", true},
		{"text", false},
		{" text ", false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			result := validator.IsEmpty(tt.value)
			if result != tt.empty {
				t.Errorf("IsEmpty('%s') = %v, want %v", tt.value, result, tt.empty)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{" TEST ", "test"},
		{"Test", "test"},
		{"  multiple   spaces  ", "multiple   spaces"},
		{"\tTAB\t", "tab"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := validator.Normalize(tt.input)
			if result != tt.expected {
				t.Errorf("Normalize('%s') = '%s', want '%s'", tt.input, result, tt.expected)
			}
		})
	}
}
