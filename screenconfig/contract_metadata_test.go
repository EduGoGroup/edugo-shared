package screenconfig

import (
	"reflect"
	"testing"
)

func strPtr(s string) *string { return &s }

func TestExtractContractMetadata(t *testing.T) {
	tests := []struct {
		name               string
		slotData           map[string]any
		requiredPermission string
		want               *ContractMetadata
	}{
		{
			name: "completa: todos los campos presentes",
			slotData: map[string]any{
				"api_prefix":      "platform",
				"api_base_path":   "/api/v1/colors",
				"resource":        "colors",
				"form_screen_key": "colors-form",
				"list_screen_key": "colors-list",
				"transforms": map[string]any{
					"submit": "identity",
				},
			},
			requiredPermission: "irrelevant.value.here",
			want: &ContractMetadata{
				APIPrefix:     "platform",
				BasePath:      "/api/v1/colors",
				Resource:      "colors",
				FormScreenKey: strPtr("colors-form"),
				ListScreenKey: strPtr("colors-list"),
				ParentIDParam: nil,
				Transforms: map[string]any{
					"submit": "identity",
				},
			},
		},
		{
			name: "parcial: resource derivado de permission, defaults aplicados",
			slotData: map[string]any{
				"api_prefix": "platform",
			},
			requiredPermission: "platform.colors.read",
			want: &ContractMetadata{
				APIPrefix:     "platform",
				BasePath:      "/api/v1/colors",
				Resource:      "colors",
				FormScreenKey: nil,
				ListScreenKey: nil,
				ParentIDParam: nil,
				Transforms:    map[string]any{},
			},
		},
		{
			name: "sin api_prefix: devuelve nil aunque el resto venga",
			slotData: map[string]any{
				"resource":        "colors",
				"api_base_path":   "/api/v1/colors",
				"form_screen_key": "colors-form",
			},
			requiredPermission: "platform.colors.read",
			want:               nil,
		},
		{
			name: "nested: parent_id_param presente, basePath calculado",
			slotData: map[string]any{
				"api_prefix":      "academic",
				"resource":        "questions",
				"parent_id_param": "assessmentId",
			},
			requiredPermission: "academic.questions.read",
			want: &ContractMetadata{
				APIPrefix:     "academic",
				BasePath:      "/api/v1/questions",
				Resource:      "questions",
				FormScreenKey: nil,
				ListScreenKey: nil,
				ParentIDParam: strPtr("assessmentId"),
				Transforms:    map[string]any{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtractContractMetadata(tc.slotData, tc.requiredPermission)
			if !equalContractMetadata(got, tc.want) {
				t.Fatalf("ExtractContractMetadata() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestParseResourceFromPermission(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"platform", ""},
		{"platform.colors", ""},
		{"platform.colors.read", "colors"},
		{"academic.questions.write", "questions"},
		{"identity.users.roles.assign", "users"},
	}
	for _, c := range cases {
		if got := parseResourceFromPermission(c.in); got != c.want {
			t.Errorf("parseResourceFromPermission(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// equalContractMetadata compara dos *ContractMetadata considerando los
// punteros a string y el mapa transforms. Devuelve true tambien cuando ambos
// son nil.
func equalContractMetadata(a, b *ContractMetadata) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.APIPrefix != b.APIPrefix ||
		a.BasePath != b.BasePath ||
		a.Resource != b.Resource {
		return false
	}
	if !equalStringPtr(a.FormScreenKey, b.FormScreenKey) ||
		!equalStringPtr(a.ListScreenKey, b.ListScreenKey) ||
		!equalStringPtr(a.ParentIDParam, b.ParentIDParam) {
		return false
	}
	return reflect.DeepEqual(a.Transforms, b.Transforms)
}

func equalStringPtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
