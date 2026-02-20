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

	var resultMap map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &resultMap))

	assert.NotContains(t, resultMap, "platformOverrides", "platformOverrides should be removed from result")

	zones, ok := resultMap["zones"].([]interface{})
	require.True(t, ok)
	require.Len(t, zones, 1)

	zone := zones[0].(map[string]interface{})
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

	var resultMap map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &resultMap))

	zones := resultMap["zones"].([]interface{})
	zone := zones[0].(map[string]interface{})
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

	var resultMap map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &resultMap))

	zones := resultMap["zones"].([]interface{})
	zone := zones[0].(map[string]interface{})
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
