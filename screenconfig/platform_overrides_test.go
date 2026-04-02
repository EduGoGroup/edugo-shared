package screenconfig

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	var resultMap map[string]any
	require.NoError(t, json.Unmarshal(result, &resultMap))

	assert.NotContains(t, resultMap, "platformOverrides", "platformOverrides should be removed from result")

	zones, ok := resultMap["zones"].([]any)
	require.True(t, ok)
	require.Len(t, zones, 1)

	zone, ok2 := zones[0].(map[string]any)
	require.True(t, ok2)
	assert.Equal(t, "grid", zone["distribution"])
	assert.Equal(t, float64(3), zone["columns"])
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

	result := ApplyPlatformOverrides(definition, "ios")

	var resultMap map[string]any
	require.NoError(t, json.Unmarshal(result, &resultMap))

	zones, ok := resultMap["zones"].([]any)
	require.True(t, ok)
	zone, ok2 := zones[0].(map[string]any)
	require.True(t, ok2)
	assert.Equal(t, float64(44), zone["height"], "ios should fallback to mobile override")
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

	result := ApplyPlatformOverrides(definition, "web")

	var resultMap map[string]any
	require.NoError(t, json.Unmarshal(result, &resultMap))

	zones, ok := resultMap["zones"].([]any)
	require.True(t, ok)
	zone, ok2 := zones[0].(map[string]any)
	require.True(t, ok2)
	assert.Equal(t, "stacked", zone["distribution"], "should remain unchanged without matching override")
}

func TestApplyPlatformOverrides_NoOverridesKey(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [
			{"id": "list_content", "distribution": "stacked"}
		]
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	assert.JSONEq(t, string(definition), string(result))
}

func TestApplyPlatformOverrides_InvalidJSON(t *testing.T) {
	definition := json.RawMessage(`{invalid json}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	assert.Equal(t, string(definition), string(result), "invalid JSON should return definition unchanged")
}

func TestApplyPlatformOverrides_OverridesNotAMap(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [{"id": "header"}],
		"platformOverrides": "not a map"
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	assert.JSONEq(t, string(definition), string(result), "non-map platformOverrides should return definition unchanged")
}

func TestApplyPlatformOverrides_PlatformValueNotAMap(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [{"id": "header"}],
		"platformOverrides": {
			"desktop": "not a map"
		}
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	assert.JSONEq(t, string(definition), string(result), "non-map platform value should return definition unchanged")
}

func TestApplyPlatformOverrides_ZoneOverridesNotAMap(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [{"id": "header", "height": 60}],
		"platformOverrides": {
			"desktop": {
				"zones": "not a map"
			}
		}
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	var resultMap map[string]any
	require.NoError(t, json.Unmarshal(result, &resultMap))
	assert.NotContains(t, resultMap, "platformOverrides", "platformOverrides should still be removed")
}

func TestApplyPlatformOverrides_DefMapZonesNotAnArray(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": "not an array",
		"platformOverrides": {
			"desktop": {
				"zones": {"header": {"height": 44}}
			}
		}
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	var resultMap map[string]any
	require.NoError(t, json.Unmarshal(result, &resultMap))
	assert.NotContains(t, resultMap, "platformOverrides")
	assert.Equal(t, "not an array", resultMap["zones"])
}

func TestApplyPlatformOverrides_ZoneItemNotAMap(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": ["not a map", {"id": "header", "height": 60}],
		"platformOverrides": {
			"desktop": {
				"zones": {"header": {"height": 44}}
			}
		}
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	var resultMap map[string]any
	require.NoError(t, json.Unmarshal(result, &resultMap))
	zones := resultMap["zones"].([]any)
	assert.Equal(t, "not a map", zones[0], "non-map zone items should be skipped")
	zone1 := zones[1].(map[string]any)
	assert.Equal(t, float64(44), zone1["height"], "valid zone should still get override applied")
}

func TestApplyPlatformOverrides_ZoneWithoutStringID(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [{"id": 123, "height": 60}, {"id": "header", "height": 60}],
		"platformOverrides": {
			"desktop": {
				"zones": {"header": {"height": 44}}
			}
		}
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	var resultMap map[string]any
	require.NoError(t, json.Unmarshal(result, &resultMap))
	zones := resultMap["zones"].([]any)
	zone0 := zones[0].(map[string]any)
	assert.Equal(t, float64(60), zone0["height"], "zone with non-string id should be skipped")
	zone1 := zones[1].(map[string]any)
	assert.Equal(t, float64(44), zone1["height"], "zone with string id should get override")
}

func TestApplyPlatformOverrides_ZoneIDNotInOverrides(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [{"id": "footer", "height": 40}],
		"platformOverrides": {
			"desktop": {
				"zones": {"header": {"height": 44}}
			}
		}
	}`)

	result := ApplyPlatformOverrides(definition, "desktop")

	var resultMap map[string]any
	require.NoError(t, json.Unmarshal(result, &resultMap))
	zones := resultMap["zones"].([]any)
	zone := zones[0].(map[string]any)
	assert.Equal(t, float64(40), zone["height"], "zone with no matching override should keep original values")
}
