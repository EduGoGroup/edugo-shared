package screenconfig

import (
	"encoding/json"
	"strings"
)

// ResolveSlots replaces "slot:xxx" references in a template definition
// with corresponding values from slotData. Returns definition unchanged if slotData is empty/null.
func ResolveSlots(definition json.RawMessage, slotData json.RawMessage) json.RawMessage {
	if len(slotData) == 0 || string(slotData) == "null" || string(slotData) == "{}" {
		return definition
	}

	var slots map[string]interface{}
	if err := json.Unmarshal(slotData, &slots); err != nil {
		return definition
	}

	if len(slots) == 0 {
		return definition
	}

	var defMap interface{}
	if err := json.Unmarshal(definition, &defMap); err != nil {
		return definition
	}

	resolved := resolveValue(defMap, slots)

	result, err := json.Marshal(resolved)
	if err != nil {
		return definition
	}

	return result
}

// resolveValue recursively resolves slot:xxx references in a parsed JSON value.
func resolveValue(value interface{}, slots map[string]interface{}) interface{} {
	switch v := value.(type) {
	case string:
		if strings.HasPrefix(v, "slot:") {
			slotKey := strings.TrimPrefix(v, "slot:")
			if slotValue, ok := slots[slotKey]; ok {
				return slotValue
			}
		}
		return v
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			result[key] = resolveValue(val, slots)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = resolveValue(val, slots)
		}
		return result
	default:
		return v
	}
}
