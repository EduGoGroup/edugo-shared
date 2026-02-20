package screenconfig

import (
	"testing"
)

func TestValidatePattern_ValidPatterns(t *testing.T) {
	validPatterns := []string{
		"login", "form", "list", "dashboard", "settings",
		"detail", "search", "profile", "modal", "notification",
		"onboarding", "empty-state",
	}
	for _, p := range validPatterns {
		if err := ValidatePattern(p); err != nil {
			t.Errorf("expected pattern %q to be valid but got error: %v", p, err)
		}
	}
}

func TestValidatePattern_InvalidPatterns(t *testing.T) {
	invalidPatterns := []string{"", "unknown", "LOGIN", "foo-bar"}
	for _, p := range invalidPatterns {
		if err := ValidatePattern(p); err == nil {
			t.Errorf("expected pattern %q to be invalid but got no error", p)
		}
	}
}

func TestValidateActionType_ValidTypes(t *testing.T) {
	validTypes := []string{
		"NAVIGATE", "NAVIGATE_BACK", "API_CALL", "SUBMIT_FORM",
		"REFRESH", "CONFIRM", "LOGOUT",
	}
	for _, a := range validTypes {
		if err := ValidateActionType(a); err != nil {
			t.Errorf("expected action type %q to be valid but got error: %v", a, err)
		}
	}
}

func TestValidateActionType_InvalidTypes(t *testing.T) {
	invalidTypes := []string{"", "navigate", "UNKNOWN", "DELETE"}
	for _, a := range invalidTypes {
		if err := ValidateActionType(a); err == nil {
			t.Errorf("expected action type %q to be invalid but got no error", a)
		}
	}
}

func TestValidateScreenType_ValidTypes(t *testing.T) {
	validTypes := []string{
		"list", "detail", "create", "edit", "dashboard", "settings",
	}
	for _, st := range validTypes {
		if err := ValidateScreenType(st); err != nil {
			t.Errorf("expected screen type %q to be valid but got error: %v", st, err)
		}
	}
}

func TestValidateScreenType_InvalidTypes(t *testing.T) {
	invalidTypes := []string{"", "LIST", "unknown", "view"}
	for _, st := range invalidTypes {
		if err := ValidateScreenType(st); err == nil {
			t.Errorf("expected screen type %q to be invalid but got no error", st)
		}
	}
}

func TestValidatePlatform_ValidPlatforms(t *testing.T) {
	validPlatforms := []string{"ios", "android", "mobile", "desktop", "web"}
	for _, p := range validPlatforms {
		if err := ValidatePlatform(p); err != nil {
			t.Errorf("expected platform %q to be valid but got error: %v", p, err)
		}
	}
}

func TestValidatePlatform_InvalidPlatforms(t *testing.T) {
	invalidPlatforms := []string{"", "iOS", "ANDROID", "windows", "linux"}
	for _, p := range invalidPlatforms {
		if err := ValidatePlatform(p); err == nil {
			t.Errorf("expected platform %q to be invalid but got no error", p)
		}
	}
}

func TestResolvePlatformOverrideKey(t *testing.T) {
	overrides := map[string]interface{}{
		"mobile":  nil,
		"desktop": nil,
	}

	// ios deberia hacer fallback a mobile
	key, ok := ResolvePlatformOverrideKey(PlatformIOS, overrides)
	if !ok || key != "mobile" {
		t.Errorf("expected ios to fallback to mobile, got key=%q ok=%v", key, ok)
	}

	// android deberia hacer fallback a mobile
	key, ok = ResolvePlatformOverrideKey(PlatformAndroid, overrides)
	if !ok || key != "mobile" {
		t.Errorf("expected android to fallback to mobile, got key=%q ok=%v", key, ok)
	}

	// desktop deberia encontrar directamente
	key, ok = ResolvePlatformOverrideKey(PlatformDesktop, overrides)
	if !ok || key != "desktop" {
		t.Errorf("expected desktop to match directly, got key=%q ok=%v", key, ok)
	}

	// web no existe en overrides y no tiene fallback
	key, ok = ResolvePlatformOverrideKey(PlatformWeb, overrides)
	if ok {
		t.Errorf("expected web to not match, got key=%q ok=%v", key, ok)
	}
}

func TestResolvePlatformOverrideKey_SpecificOverIDE(t *testing.T) {
	// Cuando existe override especifico de ios, debe preferirlo sobre mobile
	overrides := map[string]interface{}{
		"ios":    nil,
		"mobile": nil,
	}

	key, ok := ResolvePlatformOverrideKey(PlatformIOS, overrides)
	if !ok || key != "ios" {
		t.Errorf("expected ios to match directly when available, got key=%q ok=%v", key, ok)
	}
}

func TestValidateTemplateDefinition_Valid(t *testing.T) {
	definition := []byte(`{
		"zones": [
			{
				"id": "header",
				"type": "fixed",
				"slots": [
					{"id": "title", "controlType": "text"},
					{"id": "subtitle", "controlType": "text"}
				]
			}
		]
	}`)

	if err := ValidateTemplateDefinition(definition); err != nil {
		t.Errorf("expected valid definition but got error: %v", err)
	}
}

func TestValidateTemplateDefinition_InvalidJSON(t *testing.T) {
	definition := []byte(`not json`)
	if err := ValidateTemplateDefinition(definition); err == nil {
		t.Error("expected error for invalid JSON but got nil")
	}
}

func TestValidateTemplateDefinition_NoZones(t *testing.T) {
	definition := []byte(`{"zones": []}`)
	if err := ValidateTemplateDefinition(definition); err == nil {
		t.Error("expected error for empty zones but got nil")
	}
}

func TestValidateTemplateDefinition_ZoneMissingID(t *testing.T) {
	definition := []byte(`{
		"zones": [{"type": "fixed", "slots": []}]
	}`)
	if err := ValidateTemplateDefinition(definition); err == nil {
		t.Error("expected error for zone missing id but got nil")
	}
}

func TestValidateTemplateDefinition_ZoneMissingType(t *testing.T) {
	definition := []byte(`{
		"zones": [{"id": "header", "slots": []}]
	}`)
	if err := ValidateTemplateDefinition(definition); err == nil {
		t.Error("expected error for zone missing type but got nil")
	}
}

func TestValidateTemplateDefinition_SlotMissingID(t *testing.T) {
	definition := []byte(`{
		"zones": [
			{
				"id": "header",
				"type": "fixed",
				"slots": [{"controlType": "text"}]
			}
		]
	}`)
	if err := ValidateTemplateDefinition(definition); err == nil {
		t.Error("expected error for slot missing id but got nil")
	}
}

func TestValidateTemplateDefinition_SlotMissingControlType(t *testing.T) {
	definition := []byte(`{
		"zones": [
			{
				"id": "header",
				"type": "fixed",
				"slots": [{"id": "title"}]
			}
		]
	}`)
	if err := ValidateTemplateDefinition(definition); err == nil {
		t.Error("expected error for slot missing controlType but got nil")
	}
}
