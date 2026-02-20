package screenconfig

import (
	"sort"
	"testing"
)

func TestExtractResourceKeys_BasicParsing(t *testing.T) {
	permissions := []string{"users:read", "materials:create", "assessments:delete"}

	keys := ExtractResourceKeys(permissions)
	sort.Strings(keys)

	expected := []string{"assessments", "materials", "users"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected key %q at index %d, got %q", expected[i], i, k)
		}
	}
}

func TestExtractResourceKeys_Dedup(t *testing.T) {
	permissions := []string{"users:read", "users:create", "users:delete"}

	keys := ExtractResourceKeys(permissions)

	if len(keys) != 1 {
		t.Fatalf("expected 1 unique key, got %d", len(keys))
	}
	if keys[0] != "users" {
		t.Errorf("expected 'users', got %q", keys[0])
	}
}

func TestExtractResourceKeys_Empty(t *testing.T) {
	keys := ExtractResourceKeys([]string{})

	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}

	keys = ExtractResourceKeys(nil)
	if len(keys) != 0 {
		t.Errorf("expected 0 keys for nil, got %d", len(keys))
	}
}

func TestExtractResourceKeys_Malformed(t *testing.T) {
	permissions := []string{"nocolon", "valid:read", "", "also:valid"}

	keys := ExtractResourceKeys(permissions)
	sort.Strings(keys)

	// "nocolon" and "" should be skipped (no colon)
	expected := []string{"also", "valid"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d: %v", len(expected), len(keys), keys)
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected key %q at index %d, got %q", expected[i], i, k)
		}
	}
}

func TestHasPermission_Found(t *testing.T) {
	perms := []string{"materials:read", "materials:write", "assessments:read"}

	if !HasPermission(perms, "materials:read") {
		t.Error("expected to find 'materials:read'")
	}
	if !HasPermission(perms, "assessments:read") {
		t.Error("expected to find 'assessments:read'")
	}
}

func TestHasPermission_NotFound(t *testing.T) {
	perms := []string{"materials:read"}

	if HasPermission(perms, "assessments:read") {
		t.Error("should not find 'assessments:read'")
	}
	if HasPermission(perms, "materials:write") {
		t.Error("should not find 'materials:write'")
	}
}

func TestHasPermission_EmptyPerms(t *testing.T) {
	if HasPermission([]string{}, "materials:read") {
		t.Error("empty slice should return false")
	}
	if HasPermission(nil, "materials:read") {
		t.Error("nil slice should return false")
	}
}
