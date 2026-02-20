package screenconfig

import "strings"

// ExtractResourceKeys extrae las claves unicas de recursos desde un slice de permisos
// en formato "resource:action".
//
// Los permisos malformados (sin ":" o vacios) son ignorados silenciosamente.
//
// Ejemplo:
//
//	Input:  []string{"users:read", "users:write", "materials:create"}
//	Output: []string{"users", "materials"}
func ExtractResourceKeys(permissions []string) []string {
	keySet := make(map[string]bool)
	for _, perm := range permissions {
		parts := strings.SplitN(perm, ":", 2)
		if len(parts) >= 2 {
			keySet[parts[0]] = true
		}
	}

	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	return keys
}

// HasPermission verifica si un permiso especifico existe en el slice de permisos.
//
// Ejemplo:
//
//	perms := []string{"materials:read", "materials:write"}
//	HasPermission(perms, "materials:read")   // true
//	HasPermission(perms, "materials:delete") // false
func HasPermission(permissions []string, required string) bool {
	for _, p := range permissions {
		if p == required {
			return true
		}
	}
	return false
}
