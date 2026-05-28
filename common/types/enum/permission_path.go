package enum

import "regexp"

// PathPermissionRegex valida el formato path-based de un Permission o
// pattern de grant:
//
//	`*`                  → super wildcard
//	`prefix.*`           → wildcard de subárbol (1 a 3 niveles)
//	`*.suffix`           → wildcard leading (matchea cualquier path cuyo
//	                       último segmento sea `suffix`)
//	`prefix.*.suffix`    → wildcard medio (matchea `prefix.<algo>.suffix`,
//	                       con uno o más segmentos intermedios)
//	`a` / `a.b` / `a.b.c`→ exacto (1-3 segmentos)
//	cualquiera + `:own`  → sufijo de ownership (D1)
//
// El regex acepta patterns (con wildcards) además de strings exactos —
// `IsValid()` sobre Permission, en cambio, verifica pertenencia al
// catálogo cerrado AllPermissions.
var PathPermissionRegex = regexp.MustCompile(
	`^(` +
		`\*` +
		`|[a-z_]+(\.[a-z_]+){0,3}(\.\*)?` +
		`|\*\.[a-z_]+` +
		`|[a-z_]+\.\*\.[a-z_]+` +
		`)(:own)?$`,
)

// IsPathFormat reporta si s respeta la gramática path-based.
func IsPathFormat(s string) bool {
	return PathPermissionRegex.MatchString(s)
}
