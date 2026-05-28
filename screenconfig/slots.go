package screenconfig

import (
	"encoding/json"
	"strings"
)

// ResolveSlots reemplaza las referencias "slot:xxx" en una definicion de plantilla
// con los valores correspondientes de slotData.
//
// Parametros:
//   - definition: JSON que puede contener referencias "slot:xxx" en strings
//   - slotData: JSON con los valores para reemplazar (formato: {"xxx": valor})
//
// Retorna la definicion con los slots resueltos. Si slotData es vacio, null, o no es
// valido, retorna la definicion sin cambios. Las referencias a slots no encontrados
// se mantienen sin resolver.
//
// Ejemplo:
//
//	definition: {"title": "slot:page_title"}
//	slotData:   {"page_title": "Mi Pagina"}
//	resultado:  {"title": "Mi Pagina"}
func ResolveSlots(definition json.RawMessage, slotData json.RawMessage) json.RawMessage {
	if len(slotData) == 0 || string(slotData) == "null" || string(slotData) == "{}" {
		return definition
	}

	var slots map[string]any
	if err := json.Unmarshal(slotData, &slots); err != nil {
		return definition
	}

	if len(slots) == 0 {
		return definition
	}

	var defMap any
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

// resolveValue resuelve recursivamente referencias slot:xxx en un valor JSON parseado.
func resolveValue(value any, slots map[string]any) any {
	switch v := value.(type) {
	case string:
		if after, ok := strings.CutPrefix(v, "slot:"); ok {
			slotKey := after
			if slotValue, ok := slots[slotKey]; ok {
				return slotValue
			}
		}
		return v
	case map[string]any:
		result := make(map[string]any, len(v))
		for key, val := range v {
			result[key] = resolveValue(val, slots)
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, val := range v {
			result[i] = resolveValue(val, slots)
		}
		return result
	default:
		return v
	}
}
