package screenconfig

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractResourceKeys_BasicParsing(t *testing.T) {
	permissions := []string{"users:read", "materials:create", "assessments:delete"}

	keys := ExtractResourceKeys(permissions)
	sort.Strings(keys)

	require.Len(t, keys, 3)
	assert.Equal(t, []string{"assessments", "materials", "users"}, keys)
}

func TestExtractResourceKeys_Dedup(t *testing.T) {
	permissions := []string{"users:read", "users:create", "users:delete"}

	keys := ExtractResourceKeys(permissions)

	require.Len(t, keys, 1)
	assert.Equal(t, "users", keys[0])
}

func TestExtractResourceKeys_Empty(t *testing.T) {
	assert.Empty(t, ExtractResourceKeys([]string{}))
	assert.Empty(t, ExtractResourceKeys(nil))
}

func TestExtractResourceKeys_Malformed(t *testing.T) {
	permissions := []string{"nocolon", "valid:read", "", "also:valid"}

	keys := ExtractResourceKeys(permissions)
	sort.Strings(keys)

	assert.Equal(t, []string{"also", "valid"}, keys)
}

func TestHasPermission_Found(t *testing.T) {
	perms := []string{"materials:read", "materials:write", "assessments:read"}

	assert.True(t, HasPermission(perms, "materials:read"))
	assert.True(t, HasPermission(perms, "assessments:read"))
}

func TestHasPermission_NotFound(t *testing.T) {
	perms := []string{"materials:read"}

	assert.False(t, HasPermission(perms, "assessments:read"))
	assert.False(t, HasPermission(perms, "materials:write"))
}

func TestHasPermission_EmptyPerms(t *testing.T) {
	assert.False(t, HasPermission([]string{}, "materials:read"))
	assert.False(t, HasPermission(nil, "materials:read"))
}
