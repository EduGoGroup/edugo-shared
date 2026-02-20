package screenconfig

import "sort"

// MenuNode is the generic input for building the tree. Both APIs convert their domain types to MenuNode.
type MenuNode struct {
	ID          string
	Key         string
	DisplayName string
	Icon        string
	ParentID    string            // empty if top-level
	SortOrder   int
	Scope       string
	Permissions []string          // optional
	Screens     map[string]string // optional: screenType -> screenKey
}

// MenuTreeItem is a node in the built tree.
type MenuTreeItem struct {
	Key         string            `json:"key"`
	DisplayName string            `json:"displayName"`
	Icon        string            `json:"icon,omitempty"`
	Scope       string            `json:"scope,omitempty"`
	SortOrder   int               `json:"sortOrder"`
	ScreenKey   string            `json:"screenKey,omitempty"`
	Permissions []string          `json:"permissions,omitempty"`
	Screens     map[string]string `json:"screens,omitempty"`
	Children    []MenuTreeItem    `json:"children,omitempty"`
}

// BuildMenuTree builds a hierarchical tree from a flat list of MenuNodes.
// visibleKeys: if non-nil, only includes nodes whose Key is in the set.
// screenMap: if non-nil, maps resourceKey -> default screenKey.
func BuildMenuTree(nodes []MenuNode, visibleKeys map[string]bool, screenMap map[string]string) []MenuTreeItem {
	// Group nodes by parentID
	childrenOf := make(map[string][]MenuNode)
	for _, n := range nodes {
		childrenOf[n.ParentID] = append(childrenOf[n.ParentID], n)
	}

	return buildSubTree(childrenOf, "", visibleKeys, screenMap)
}

func buildSubTree(childrenOf map[string][]MenuNode, parentID string, visibleKeys map[string]bool, screenMap map[string]string) []MenuTreeItem {
	nodes := childrenOf[parentID]
	var items []MenuTreeItem

	for _, n := range nodes {
		// Filter by visibleKeys if provided
		if visibleKeys != nil && !visibleKeys[n.Key] {
			continue
		}

		item := MenuTreeItem{
			Key:         n.Key,
			DisplayName: n.DisplayName,
			Icon:        n.Icon,
			Scope:       n.Scope,
			SortOrder:   n.SortOrder,
		}

		if len(n.Permissions) > 0 {
			item.Permissions = n.Permissions
		}

		if len(n.Screens) > 0 {
			item.Screens = n.Screens
		}

		// Apply default screen key from screenMap
		if screenMap != nil {
			if sk, ok := screenMap[n.Key]; ok {
				item.ScreenKey = sk
			}
		}

		// Recursively build children
		children := buildSubTree(childrenOf, n.ID, visibleKeys, screenMap)
		if len(children) > 0 {
			item.Children = children
		}

		items = append(items, item)
	}

	// Sort by SortOrder
	sort.Slice(items, func(i, j int) bool {
		return items[i].SortOrder < items[j].SortOrder
	})

	return items
}
