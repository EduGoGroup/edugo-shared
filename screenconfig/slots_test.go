package screenconfig

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveSlots_BasicReplacement(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title", "subtitle": "static text"}`)
	slotData := json.RawMessage(`{"page_title": "My Page Title"}`)

	result := ResolveSlots(definition, slotData)

	var resultMap map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &resultMap))

	assert.Equal(t, "My Page Title", resultMap["title"])
	assert.Equal(t, "static text", resultMap["subtitle"])
}

func TestResolveSlots_EmptySlotData(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title"}`)

	result := ResolveSlots(definition, json.RawMessage(`{}`))

	assert.Equal(t, string(definition), string(result))
}

func TestResolveSlots_NullSlotData(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title"}`)

	result := ResolveSlots(definition, json.RawMessage(`null`))

	assert.Equal(t, string(definition), string(result))
}

func TestResolveSlots_NilSlotData(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title"}`)

	result := ResolveSlots(definition, nil)

	assert.Equal(t, string(definition), string(result))
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

	assert.False(t, strings.Contains(resultStr, "slot:header_title"), "slot:header_title should have been resolved")
	assert.False(t, strings.Contains(resultStr, "slot:item_label"), "slot:item_label should have been resolved")
	assert.Contains(t, resultStr, "Dashboard")
	assert.Contains(t, resultStr, "Home")
	assert.Contains(t, resultStr, "fixed")
}

func TestResolveSlots_UnknownSlotKey_KeepsOriginal(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:unknown_key"}`)
	slotData := json.RawMessage(`{"other_key": "value"}`)

	result := ResolveSlots(definition, slotData)

	var resultMap map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &resultMap))

	assert.Equal(t, "slot:unknown_key", resultMap["title"])
}
