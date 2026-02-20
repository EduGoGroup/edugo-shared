package screenconfig

import "encoding/json"

// ApplyPlatformOverrides applies platform-specific zone overrides to a template definition.
// Uses the fallback chain (e.g., ios -> mobile -> none). Returns definition with overrides applied
// and the "platformOverrides" key removed.
func ApplyPlatformOverrides(definition json.RawMessage, platform string) json.RawMessage {
	var defMap map[string]interface{}
	if err := json.Unmarshal(definition, &defMap); err != nil {
		return definition
	}

	overrides, ok := defMap["platformOverrides"]
	if !ok {
		return definition
	}

	overridesMap, ok := overrides.(map[string]interface{})
	if !ok {
		return definition
	}

	// Resolve override with fallback chain (e.g., ios -> mobile -> no override)
	resolvedKey, found := ResolvePlatformOverrideKey(Platform(platform), overridesMap)
	if !found {
		return definition
	}

	platformOverride, ok := overridesMap[resolvedKey]
	if !ok {
		return definition
	}

	platformMap, ok := platformOverride.(map[string]interface{})
	if !ok {
		return definition
	}

	// Apply zone overrides
	applyZoneOverrides(defMap, platformMap)

	// Remove platformOverrides from final result
	delete(defMap, "platformOverrides")

	result, err := json.Marshal(defMap)
	if err != nil {
		return definition
	}

	return result
}

// applyZoneOverrides merges zone-level overrides from platformMap into defMap["zones"].
func applyZoneOverrides(defMap, platformMap map[string]interface{}) {
	zonesMap, ok := toStringMap(platformMap["zones"])
	if !ok {
		return
	}

	zonesArr, ok := defMap["zones"].([]interface{})
	if !ok {
		return
	}

	for i, zone := range zonesArr {
		zoneMap, ok := zone.(map[string]interface{})
		if !ok {
			continue
		}
		zoneID, _ := zoneMap["id"].(string)
		overrideMap, ok := toStringMap(zonesMap[zoneID])
		if !ok {
			continue
		}
		for k, v := range overrideMap {
			zoneMap[k] = v
		}
		zonesArr[i] = zoneMap
	}
}

func toStringMap(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok
}
