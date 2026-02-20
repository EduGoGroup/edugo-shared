package screenconfig

import (
	"testing"
)

func TestBuildMenuTree_FlatList(t *testing.T) {
	nodes := []MenuNode{
		{ID: "1", Key: "dashboard", DisplayName: "Dashboard", SortOrder: 0, Scope: "system"},
		{ID: "2", Key: "materials", DisplayName: "Materials", SortOrder: 1, Scope: "school"},
		{ID: "3", Key: "settings", DisplayName: "Settings", SortOrder: 2, Scope: "system"},
	}

	items := BuildMenuTree(nodes, nil, nil)

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[0].Key != "dashboard" {
		t.Errorf("expected first item 'dashboard', got %q", items[0].Key)
	}
	if items[1].Key != "materials" {
		t.Errorf("expected second item 'materials', got %q", items[1].Key)
	}
	if items[2].Key != "settings" {
		t.Errorf("expected third item 'settings', got %q", items[2].Key)
	}
}

func TestBuildMenuTree_Hierarchical(t *testing.T) {
	nodes := []MenuNode{
		{ID: "1", Key: "admin", DisplayName: "Admin", SortOrder: 0},
		{ID: "2", Key: "users", DisplayName: "Users", ParentID: "1", SortOrder: 0},
		{ID: "3", Key: "roles", DisplayName: "Roles", ParentID: "1", SortOrder: 1},
	}

	items := BuildMenuTree(nodes, nil, nil)

	if len(items) != 1 {
		t.Fatalf("expected 1 top-level item, got %d", len(items))
	}
	if items[0].Key != "admin" {
		t.Errorf("expected 'admin', got %q", items[0].Key)
	}
	if len(items[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(items[0].Children))
	}
	if items[0].Children[0].Key != "users" {
		t.Errorf("expected first child 'users', got %q", items[0].Children[0].Key)
	}
	if items[0].Children[1].Key != "roles" {
		t.Errorf("expected second child 'roles', got %q", items[0].Children[1].Key)
	}
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

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Key != "dashboard" {
		t.Errorf("expected 'dashboard', got %q", items[0].Key)
	}
	if items[1].Key != "settings" {
		t.Errorf("expected 'settings', got %q", items[1].Key)
	}
}

func TestBuildMenuTree_SortOrder(t *testing.T) {
	nodes := []MenuNode{
		{ID: "1", Key: "c-item", DisplayName: "C", SortOrder: 3},
		{ID: "2", Key: "a-item", DisplayName: "A", SortOrder: 1},
		{ID: "3", Key: "b-item", DisplayName: "B", SortOrder: 2},
	}

	items := BuildMenuTree(nodes, nil, nil)

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[0].Key != "a-item" {
		t.Errorf("expected first 'a-item', got %q", items[0].Key)
	}
	if items[1].Key != "b-item" {
		t.Errorf("expected second 'b-item', got %q", items[1].Key)
	}
	if items[2].Key != "c-item" {
		t.Errorf("expected third 'c-item', got %q", items[2].Key)
	}
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

	if items[0].ScreenKey != "dashboard-teacher" {
		t.Errorf("expected screenKey 'dashboard-teacher', got %q", items[0].ScreenKey)
	}
	if items[1].ScreenKey != "materials-list" {
		t.Errorf("expected screenKey 'materials-list', got %q", items[1].ScreenKey)
	}
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

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if len(items[0].Permissions) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(items[0].Permissions))
	}
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

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Screens["list"] != "materials-list" {
		t.Errorf("expected screen 'materials-list' for type 'list', got %q", items[0].Screens["list"])
	}
	if items[0].Screens["detail"] != "material-detail" {
		t.Errorf("expected screen 'material-detail' for type 'detail', got %q", items[0].Screens["detail"])
	}
}

func TestBuildMenuTree_Empty(t *testing.T) {
	items := BuildMenuTree(nil, nil, nil)

	if items != nil {
		t.Errorf("expected nil for empty input, got %v", items)
	}

	items = BuildMenuTree([]MenuNode{}, nil, nil)
	if items != nil {
		t.Errorf("expected nil for empty slice, got %v", items)
	}
}
