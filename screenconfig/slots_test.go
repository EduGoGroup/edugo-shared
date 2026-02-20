package screenconfig

import (
	"encoding/json"
	"testing"
)

func TestResolveSlots_BasicReplacement(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title", "subtitle": "static text"}`)
	slotData := json.RawMessage(`{"page_title": "My Page Title"}`)

	result := ResolveSlots(definition, slotData)

	var resultMap map[string]interface{}
	if err := json.Unmarshal(result, &resultMap); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if resultMap["title"] != "My Page Title" {
		t.Errorf("expected title 'My Page Title', got %v", resultMap["title"])
	}
	if resultMap["subtitle"] != "static text" {
		t.Errorf("expected subtitle 'static text', got %v", resultMap["subtitle"])
	}
}

func TestResolveSlots_EmptySlotData(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title"}`)

	result := ResolveSlots(definition, json.RawMessage(`{}`))

	if string(result) != string(definition) {
		t.Errorf("expected unchanged definition, got %s", string(result))
	}
}

func TestResolveSlots_NullSlotData(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title"}`)

	result := ResolveSlots(definition, json.RawMessage(`null`))

	if string(result) != string(definition) {
		t.Errorf("expected unchanged definition, got %s", string(result))
	}
}

func TestResolveSlots_NilSlotData(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title"}`)

	result := ResolveSlots(definition, nil)

	if string(result) != string(definition) {
		t.Errorf("expected unchanged definition, got %s", string(result))
	}
}

func TestResolveSlots_NestedStructure(t *testing.T) {
	definition := json.RawMessage(`{
		"nav": {
			"topBar": {
				"title": "slot:header_title"
			}
		},
		"items": [
			{"label": "slot:item_label"},
			{"label": "fixed"}
		]
	}`)
	slotData := json.RawMessage(`{
		"header_title": "Dashboard",
		"item_label": "Home"
	}`)

	result := ResolveSlots(definition, slotData)

	resultStr := string(result)
	if contains(resultStr, "slot:header_title") {
		t.Error("slot:header_title should have been resolved")
	}
	if contains(resultStr, "slot:item_label") {
		t.Error("slot:item_label should have been resolved")
	}
	if !contains(resultStr, "Dashboard") {
		t.Error("expected 'Dashboard' in result")
	}
	if !contains(resultStr, "Home") {
		t.Error("expected 'Home' in result")
	}
	if !contains(resultStr, "fixed") {
		t.Error("expected 'fixed' to remain unchanged")
	}
}

func TestResolveSlots_UnknownSlotKey_KeepsOriginal(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:unknown_key"}`)
	slotData := json.RawMessage(`{"other_key": "value"}`)

	result := ResolveSlots(definition, slotData)

	var resultMap map[string]interface{}
	if err := json.Unmarshal(result, &resultMap); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if resultMap["title"] != "slot:unknown_key" {
		t.Errorf("expected unresolved slot reference, got %v", resultMap["title"])
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
