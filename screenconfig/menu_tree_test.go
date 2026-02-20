package screenconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildMenuTree_FlatList(t *testing.T) {
	nodes := []MenuNode{
		{ID: "1", Key: "dashboard", DisplayName: "Dashboard", SortOrder: 0, Scope: "system"},
		{ID: "2", Key: "materials", DisplayName: "Materials", SortOrder: 1, Scope: "school"},
		{ID: "3", Key: "settings", DisplayName: "Settings", SortOrder: 2, Scope: "system"},
	}

	items := BuildMenuTree(nodes, nil, nil)

	require.Len(t, items, 3)
	assert.Equal(t, "dashboard", items[0].Key)
	assert.Equal(t, "materials", items[1].Key)
	assert.Equal(t, "settings", items[2].Key)
}

func TestBuildMenuTree_Hierarchical(t *testing.T) {
	nodes := []MenuNode{
		{ID: "1", Key: "admin", DisplayName: "Admin", SortOrder: 0},
		{ID: "2", Key: "users", DisplayName: "Users", ParentID: "1", SortOrder: 0},
		{ID: "3", Key: "roles", DisplayName: "Roles", ParentID: "1", SortOrder: 1},
	}

	items := BuildMenuTree(nodes, nil, nil)

	require.Len(t, items, 1)
	assert.Equal(t, "admin", items[0].Key)
	require.Len(t, items[0].Children, 2)
	assert.Equal(t, "users", items[0].Children[0].Key)
	assert.Equal(t, "roles", items[0].Children[1].Key)
}

func TestBuildMenuTree_FilterByVisibleKeys(t *testing.T) {
	nodes := []MenuNode{
		{ID: "1", Key: "dashboard", DisplayName: "Dashboard", SortOrder: 0},
		{ID: "2", Key: "materials", DisplayName: "Materials", SortOrder: 1},
		{ID: "3", Key: "settings", DisplayName: "Settings", SortOrder: 2},
	}

	visibleKeys := map[string]bool{
		"dashboard": true,
		"settings":  true,
	}

	items := BuildMenuTree(nodes, visibleKeys, nil)

	require.Len(t, items, 2)
	assert.Equal(t, "dashboard", items[0].Key)
	assert.Equal(t, "settings", items[1].Key)
}

func TestBuildMenuTree_SortOrder(t *testing.T) {
	nodes := []MenuNode{
		{ID: "1", Key: "c-item", DisplayName: "C", SortOrder: 3},
		{ID: "2", Key: "a-item", DisplayName: "A", SortOrder: 1},
		{ID: "3", Key: "b-item", DisplayName: "B", SortOrder: 2},
	}

	items := BuildMenuTree(nodes, nil, nil)

	require.Len(t, items, 3)
	assert.Equal(t, "a-item", items[0].Key)
	assert.Equal(t, "b-item", items[1].Key)
	assert.Equal(t, "c-item", items[2].Key)
}

func TestBuildMenuTree_ScreenMapping(t *testing.T) {
	nodes := []MenuNode{
		{ID: "1", Key: "dashboard", DisplayName: "Dashboard", SortOrder: 0},
		{ID: "2", Key: "materials", DisplayName: "Materials", SortOrder: 1},
	}

	screenMap := map[string]string{
		"dashboard": "dashboard-teacher",
		"materials": "materials-list",
	}

	items := BuildMenuTree(nodes, nil, screenMap)

	assert.Equal(t, "dashboard-teacher", items[0].ScreenKey)
	assert.Equal(t, "materials-list", items[1].ScreenKey)
}

func TestBuildMenuTree_WithPermissions(t *testing.T) {
	nodes := []MenuNode{
		{
			ID:          "1",
			Key:         "users",
			DisplayName: "Users",
			SortOrder:   0,
			Permissions: []string{"users:read", "users:create"},
		},
	}

	items := BuildMenuTree(nodes, nil, nil)

	require.Len(t, items, 1)
	assert.Len(t, items[0].Permissions, 2)
}

func TestBuildMenuTree_WithScreens(t *testing.T) {
	nodes := []MenuNode{
		{
			ID:          "1",
			Key:         "materials",
			DisplayName: "Materials",
			SortOrder:   0,
			Screens:     map[string]string{"list": "materials-list", "detail": "material-detail"},
		},
	}

	items := BuildMenuTree(nodes, nil, nil)

	require.Len(t, items, 1)
	assert.Equal(t, "materials-list", items[0].Screens["list"])
	assert.Equal(t, "material-detail", items[0].Screens["detail"])
}

func TestBuildMenuTree_Empty(t *testing.T) {
	items := BuildMenuTree(nil, nil, nil)
	assert.Empty(t, items)
	assert.NotNil(t, items, "should return empty slice not nil for JSON serialization")

	items = BuildMenuTree([]MenuNode{}, nil, nil)
	assert.Empty(t, items)
	assert.NotNil(t, items)
}
