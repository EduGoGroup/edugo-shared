package screenconfig

import (
	"encoding/json"
	"testing"
)

func TestApplyPlatformOverrides_MatchDesktop(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [
			{"id": "list_content", "distribution": "stacked"}
		],
		"platformOverrides": {
			"desktop": {
				"zones": {
					"list_content": {"distribution": "grid", "columns": 3}
				}
			}
		}
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	var resultMap map[string]interface{}
	if err := json.Unmarshal(result, &resultMap); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	// platformOverrides should be removed
	if _, exists := resultMap["platformOverrides"]; exists {
		t.Error("platformOverrides should be removed from result")
	}

	// zone should have overrides applied
	zones, ok := resultMap["zones"].([]interface{})
	if !ok || len(zones) == 0 {
		t.Fatal("expected zones array")
	}
	zone := zones[0].(map[string]interface{})
	if zone["distribution"] != "grid" {
		t.Errorf("expected distribution 'grid', got %v", zone["distribution"])
	}
	if zone["columns"] != float64(3) {
		t.Errorf("expected columns 3, got %v", zone["columns"])
	}
}

func TestApplyPlatformOverrides_FallbackIosToMobile(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [
			{"id": "header", "height": 60}
		],
		"platformOverrides": {
			"mobile": {
				"zones": {
					"header": {"height": 44}
				}
			}
		}
	}`)

	// ios falls back to mobile
	result := ApplyPlatformOverrides(definition, "ios")

	var resultMap map[string]interface{}
	if err := json.Unmarshal(result, &resultMap); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	zones := resultMap["zones"].([]interface{})
	zone := zones[0].(map[string]interface{})
	if zone["height"] != float64(44) {
		t.Errorf("expected height 44 (mobile fallback), got %v", zone["height"])
	}
}

func TestApplyPlatformOverrides_NoPlatformMatch(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [
			{"id": "list_content", "distribution": "stacked"}
		],
		"platformOverrides": {
			"desktop": {
				"zones": {
					"list_content": {"distribution": "grid"}
				}
			}
		}
	}`)

	// "web" has no override and no fallback
	result := ApplyPlatformOverrides(definition, "web")

	var resultMap map[string]interface{}
	if err := json.Unmarshal(result, &resultMap); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	zones := resultMap["zones"].([]interface{})
	zone := zones[0].(map[string]interface{})
	if zone["distribution"] != "stacked" {
		t.Errorf("expected distribution unchanged 'stacked', got %v", zone["distribution"])
	}
}

func TestApplyPlatformOverrides_NoOverridesKey(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [
			{"id": "list_content", "distribution": "stacked"}
		]
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	// Without platformOverrides, the template should be identical
	var origMap, resultMap map[string]interface{}
	json.Unmarshal(definition, &origMap)
	json.Unmarshal(result, &resultMap)

	origZones := origMap["zones"].([]interface{})
	resultZones := resultMap["zones"].([]interface{})
	origZone := origZones[0].(map[string]interface{})
	resultZone := resultZones[0].(map[string]interface{})

	if origZone["distribution"] != resultZone["distribution"] {
		t.Error("definition should be unchanged without platformOverrides")
	}
}

func TestApplyPlatformOverrides_InvalidJSON(t *testing.T) {
	definition := json.RawMessage(`{invalid json}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	if string(result) != string(definition) {
		t.Error("invalid JSON should return definition unchanged")
	}
}
