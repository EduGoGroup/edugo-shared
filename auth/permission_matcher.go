package auth

import "strings"

// PermissionMatches reporta si `pattern` cubre el permiso `request`.
// Es el mirror Go exacto de iam.permission_matches() (Postgres), y debe
// coincidir 1:1 con el PermissionMatcher Kotlin del KMP.
//
// Reglas (D1 path-based + extensión wildcard-first):
//
//	`*`                 → cubre cualquier request
//	`pattern == request`→ exacto
//	`prefix.*`          → cubre `prefix` y `prefix.<lo-que-sea>` (subárbol)
//	`*.suffix`          → cubre cualquier request cuyo último tramo sea
//	                      `.suffix` (ej. `*.create` matchea `users.create`
//	                      y `academic.units.create`)
//	`prefix.*.suffix`   → cubre cualquier request que empiece con
//	                      `prefix.`, tenga al menos un segmento intermedio
//	                      y termine con `.suffix`
//	sufijo `:own`       → semánticamente distinto, no se mezcla con
//	                      el mismo pattern sin `:own`
//
// La gramática válida de pattern y request vive en
// enum.PathPermissionRegex.
func PermissionMatches(pattern, request string) bool {
	if pattern == "*" {
		return true
	}
	if pattern == request {
		return true
	}
	// prefix.*  (subárbol). Debe evaluarse antes que prefix.*.suffix
	// para preservar la semántica histórica cuando el pattern termina
	// literalmente con `.*`.
	if strings.HasSuffix(pattern, ".*") {
		prefix := pattern[:len(pattern)-2]
		return request == prefix || strings.HasPrefix(request, prefix+".")
	}
	// *.suffix  → cualquier request `<algo>.suffix`
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:] // conserva el punto, ej ".create"
		return len(request) > len(suffix) && strings.HasSuffix(request, suffix)
	}
	// prefix.*.suffix  → request startsWith `prefix.` + algo + `.suffix`
	if i := strings.Index(pattern, ".*."); i > 0 {
		head := pattern[:i+1]      // `prefix.`
		tail := pattern[i+2:]      // `.suffix`
		// Validar que no haya otro `*` (los patterns soportados son
		// los listados arriba; ya descartamos `prefix.*`).
		if strings.Contains(head, "*") || strings.Contains(tail, "*") {
			return false
		}
		if !strings.HasPrefix(request, head) || !strings.HasSuffix(request, tail) {
			return false
		}
		// Si head y tail solapan en request, no hay segmento intermedio.
		if len(request) <= len(head)+len(tail) {
			return false
		}
		// Exigir al menos un segmento intermedio entre head y tail.
		middle := request[len(head) : len(request)-len(tail)]
		return !strings.HasPrefix(middle, ".") && !strings.HasSuffix(middle, ".")
	}
	return false
}

// EvaluateGrants aplica deny precedence sobre `g` para `request`:
// si algún deny matchea → false. Si no, devuelve true sólo si algún
// allow matchea. Default deny (sin allow → false).
func EvaluateGrants(g Grants, request string) bool {
	for _, d := range g.Deny {
		if PermissionMatches(d, request) {
			return false
		}
	}
	for _, a := range g.Allow {
		if PermissionMatches(a, request) {
			return true
		}
	}
	return false
}
