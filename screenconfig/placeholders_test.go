package screenconfig_test

import (
	"testing"

	"github.com/EduGoGroup/edugo-shared/screenconfig"
)

// --- ResourcePrefixFromPermission ---

func TestResourcePrefixFromPermission_Standard(t *testing.T) {
	got := screenconfig.ResourcePrefixFromPermission("content.assessments.read")
	if got != "content.assessments" {
		t.Errorf("got %q, want %q", got, "content.assessments")
	}
}

func TestResourcePrefixFromPermission_TwoSegments(t *testing.T) {
	got := screenconfig.ResourcePrefixFromPermission("content.assessments.write")
	if got != "content.assessments" {
		t.Errorf("got %q, want %q", got, "content.assessments")
	}
}

func TestResourcePrefixFromPermission_EmptyOrSingleToken(t *testing.T) {
	// Table-driven sobre los edge cases que el codigo trata como "sin
	// segmentos separables": vacio, sin punto, y leading-dot (idx==0
	// dispara la guarda idx <= 0).
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"single token no dot", "single", ""},
		{"leading dot", ".leading", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := screenconfig.ResourcePrefixFromPermission(tc.input)
			if got != tc.want {
				t.Errorf("ResourcePrefixFromPermission(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// --- ExpandResourcePlaceholders ---

func TestExpandResourcePlaceholders_PlaceholderPresent(t *testing.T) {
	in := []map[string]any{
		{"id": "save", "permission": "$resource$.read"},
	}
	got := screenconfig.ExpandResourcePlaceholders(in, "content.assessments")
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1", len(got))
	}
	if got[0]["permission"] != "content.assessments.read" {
		t.Errorf("permission = %v, want content.assessments.read", got[0]["permission"])
	}
	// id no contiene placeholder → se conserva.
	if got[0]["id"] != "save" {
		t.Errorf("id mutated: %v", got[0]["id"])
	}
}

func TestExpandResourcePlaceholders_NoPlaceholder(t *testing.T) {
	in := []map[string]any{
		{"id": "save", "permission": "content.x.read", "label": "Guardar"},
	}
	got := screenconfig.ExpandResourcePlaceholders(in, "content.assessments")
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1", len(got))
	}
	if got[0]["permission"] != "content.x.read" {
		t.Errorf("permission mutated: %v", got[0]["permission"])
	}
	if got[0]["label"] != "Guardar" {
		t.Errorf("label mutated: %v", got[0]["label"])
	}
}

func TestExpandResourcePlaceholders_EmptyPrefix(t *testing.T) {
	// prefix vacio → fail loud: devuelve actions sin tocar, dejando el
	// placeholder visible.
	in := []map[string]any{
		{"id": "save", "permission": "$resource$.read"},
	}
	got := screenconfig.ExpandResourcePlaceholders(in, "")
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1", len(got))
	}
	if got[0]["permission"] != "$resource$.read" {
		t.Errorf("permission expected unchanged ($resource$.read), got %v", got[0]["permission"])
	}
}

func TestExpandResourcePlaceholders_MultipleFields(t *testing.T) {
	// Action con multiples string-fields: solo el campo con $resource$
	// cambia; los demas se preservan tal cual.
	in := []map[string]any{
		{
			"id":         "save_new",
			"permission": "$resource$.create",
			"label":      "Nuevo",
			"icon":       "add",
			"scope":      "form-submit",
		},
	}
	got := screenconfig.ExpandResourcePlaceholders(in, "content.assessments")
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1", len(got))
	}
	a := got[0]
	if a["permission"] != "content.assessments.create" {
		t.Errorf("permission = %v, want content.assessments.create", a["permission"])
	}
	if a["id"] != "save_new" {
		t.Errorf("id mutated: %v", a["id"])
	}
	if a["label"] != "Nuevo" {
		t.Errorf("label mutated: %v", a["label"])
	}
	if a["icon"] != "add" {
		t.Errorf("icon mutated: %v", a["icon"])
	}
	if a["scope"] != "form-submit" {
		t.Errorf("scope mutated: %v", a["scope"])
	}
}
