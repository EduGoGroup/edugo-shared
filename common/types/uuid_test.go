package types_test

import (
	"encoding/json"
	"testing"

	"github.com/EduGoGroup/edugo-shared/common/types"
)

func TestNewUUID(t *testing.T) {
	id := types.NewUUID()

	if id.IsZero() {
		t.Error("NewUUID should generate non-zero UUID")
	}

	// Generate another to ensure they're different
	id2 := types.NewUUID()
	if id.String() == id2.String() {
		t.Error("NewUUID should generate unique UUIDs")
	}
}

func TestParseUUID(t *testing.T) {
	t.Run("valid_uuid", func(t *testing.T) {
		uuidStr := "123e4567-e89b-12d3-a456-426614174000"
		id, err := types.ParseUUID(uuidStr)

		if err != nil {
			t.Fatalf("ParseUUID failed: %v", err)
		}

		if id.String() != uuidStr {
			t.Errorf("Expected %s, got %s", uuidStr, id.String())
		}
	})

	t.Run("invalid_uuid", func(t *testing.T) {
		_, err := types.ParseUUID("invalid-uuid")

		if err == nil {
			t.Error("ParseUUID should return error for invalid UUID")
		}
	})

	t.Run("empty_string", func(t *testing.T) {
		_, err := types.ParseUUID("")

		if err == nil {
			t.Error("ParseUUID should return error for empty string")
		}
	})
}

func TestMustParseUUID(t *testing.T) {
	t.Run("valid_uuid", func(t *testing.T) {
		uuidStr := "123e4567-e89b-12d3-a456-426614174000"
		id := types.MustParseUUID(uuidStr)

		if id.String() != uuidStr {
			t.Errorf("Expected %s, got %s", uuidStr, id.String())
		}
	})

	t.Run("invalid_uuid_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustParseUUID should panic on invalid UUID")
			}
		}()

		types.MustParseUUID("invalid-uuid")
	})
}

func TestUUID_String(t *testing.T) {
	uuidStr := "123e4567-e89b-12d3-a456-426614174000"
	id := types.MustParseUUID(uuidStr)

	if id.String() != uuidStr {
		t.Errorf("String() = %s, want %s", id.String(), uuidStr)
	}
}

func TestUUID_IsZero(t *testing.T) {
	t.Run("zero_uuid", func(t *testing.T) {
		var id types.UUID

		if !id.IsZero() {
			t.Error("Default UUID should be zero")
		}
	})

	t.Run("non_zero_uuid", func(t *testing.T) {
		id := types.NewUUID()

		if id.IsZero() {
			t.Error("Generated UUID should not be zero")
		}
	})
}

func TestUUID_MarshalJSON(t *testing.T) {
	id := types.MustParseUUID("123e4567-e89b-12d3-a456-426614174000")

	data, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected := `"123e4567-e89b-12d3-a456-426614174000"`
	if string(data) != expected {
		t.Errorf("MarshalJSON() = %s, want %s", string(data), expected)
	}
}

func TestUUID_UnmarshalJSON(t *testing.T) {
	t.Run("valid_json", func(t *testing.T) {
		jsonData := []byte(`"123e4567-e89b-12d3-a456-426614174000"`)
		var id types.UUID

		err := json.Unmarshal(jsonData, &id)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}

		if id.String() != "123e4567-e89b-12d3-a456-426614174000" {
			t.Errorf("Unexpected UUID: %s", id.String())
		}
	})

	t.Run("invalid_json_too_short", func(t *testing.T) {
		jsonData := []byte(`"x"`)
		var id types.UUID

		err := json.Unmarshal(jsonData, &id)
		if err == nil {
			t.Error("UnmarshalJSON should error on invalid JSON")
		}
	})

	t.Run("invalid_uuid_format", func(t *testing.T) {
		jsonData := []byte(`"invalid-uuid-format"`)
		var id types.UUID

		err := json.Unmarshal(jsonData, &id)
		if err == nil {
			t.Error("UnmarshalJSON should error on invalid UUID format")
		}
	})

	t.Run("empty_json", func(t *testing.T) {
		jsonData := []byte(`""`)
		var id types.UUID

		err := json.Unmarshal(jsonData, &id)
		if err == nil {
			t.Error("UnmarshalJSON should error on empty string")
		}
	})
}

func TestUUID_Value(t *testing.T) {
	id := types.MustParseUUID("123e4567-e89b-12d3-a456-426614174000")

	value, err := id.Value()
	if err != nil {
		t.Fatalf("Value() failed: %v", err)
	}

	strValue, ok := value.(string)
	if !ok {
		t.Fatal("Value() should return string")
	}

	if strValue != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("Value() = %s, want 123e4567-e89b-12d3-a456-426614174000", strValue)
	}
}

func TestUUID_Scan(t *testing.T) {
	t.Run("scan_string", func(t *testing.T) {
		var id types.UUID
		err := id.Scan("123e4567-e89b-12d3-a456-426614174000")

		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}

		if id.String() != "123e4567-e89b-12d3-a456-426614174000" {
			t.Errorf("Scan string: got %s", id.String())
		}
	})

	t.Run("scan_bytes", func(t *testing.T) {
		var id types.UUID
		err := id.Scan([]byte("123e4567-e89b-12d3-a456-426614174000"))

		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}

		if id.String() != "123e4567-e89b-12d3-a456-426614174000" {
			t.Errorf("Scan bytes: got %s", id.String())
		}
	})

	t.Run("scan_nil", func(t *testing.T) {
		var id types.UUID
		err := id.Scan(nil)

		if err != nil {
			t.Fatalf("Scan nil failed: %v", err)
		}

		if !id.IsZero() {
			t.Error("Scanning nil should result in zero UUID")
		}
	})

	t.Run("scan_invalid_type", func(t *testing.T) {
		var id types.UUID
		err := id.Scan(123)

		if err == nil {
			t.Error("Scan should error on invalid type")
		}
	})

	t.Run("scan_invalid_uuid", func(t *testing.T) {
		var id types.UUID
		err := id.Scan("invalid-uuid")

		if err == nil {
			t.Error("Scan should error on invalid UUID string")
		}
	})
}

func TestUUID_RoundTrip(t *testing.T) {
	// Test that we can marshal and unmarshal without losing data
	original := types.NewUUID()

	// JSON round trip
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded types.UUID
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if original.String() != decoded.String() {
		t.Errorf("Round trip failed: %s != %s", original.String(), decoded.String())
	}
}

func TestUUID_InStruct(t *testing.T) {
	type TestStruct struct {
		ID   types.UUID `json:"id"`
		Name string     `json:"name"`
	}

	// Create struct with UUID
	original := TestStruct{
		ID:   types.NewUUID(),
		Name: "test",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal back
	var decoded TestStruct
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if original.ID.String() != decoded.ID.String() {
		t.Error("UUID in struct not preserved during JSON round trip")
	}
	if original.Name != decoded.Name {
		t.Error("Name in struct not preserved during JSON round trip")
	}
}
