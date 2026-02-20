package screenconfig

import "sort"

// MenuNode es la entrada generica para construir el arbol.
// Ambas APIs convierten sus tipos de dominio a MenuNode antes de llamar a BuildMenuTree.
type MenuNode struct {
	ID          string
	Key         string
	DisplayName string
	Icon        string
	ParentID    string            // vacio si top-level
	SortOrder   int
	Scope       string
	Permissions []string          // opcional
	Screens     map[string]string // opcional: screenType -> screenKey
}

// MenuTreeItem es un nodo en el arbol construido.
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

// BuildMenuTree construye un arbol jerarquico a partir de una lista plana de MenuNodes.
// Los nodos se organizan por su campo ParentID, y se ordenan por SortOrder.
//
// Parametros:
//   - nodes: Lista plana de nodos de menu a procesar
//   - visibleKeys: Si no es nil, solo incluye nodos cuya Key este en el conjunto
//   - screenMap: Si no es nil, mapea resourceKey a screenKey por defecto
//
// Retorna un slice de MenuTreeItem ordenado por SortOrder.
// Si no hay nodos validos, retorna un slice vacio (no nil).
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
	items := make([]MenuTreeItem, 0)

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
