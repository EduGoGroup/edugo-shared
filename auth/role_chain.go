package auth

// ResolveRoleChain devuelve la cadena de herencia de un rol: el propio
// rol (rootID) seguido de todos sus ancestros a lo largo de
// parent_role_id (ADR-6), de más cercano a más lejano. Es agnóstica al
// almacenamiento: `parentOf` resuelve el padre de un rol y reporta
// ok=false (o parent="") cuando el rol es canónico (sin padre).
//
// Soporta profundidad N aunque hoy el seed solo use profundidad 1.
// Incluye guard de ciclos: si un rol ya fue visitado, corta la cadena
// sin error — los grants ya acumulados son válidos y la precedencia
// deny-wins del matcher se mantiene (defensa ante datos corruptos; el
// schema no impide un ciclo A→B→A).
//
// Es la versión pura del CTE recursivo de login (edugo-api-identity
// fetchGrantsForUser/resolveRoleChain): misma semántica, expresada como
// caminata sin acoplarse a la BD, para reusarla desde el admin Go.
func ResolveRoleChain(rootID string, parentOf func(id string) (parent string, ok bool, err error)) ([]string, error) {
	chain := make([]string, 0, 4)
	visited := make(map[string]struct{}, 4)
	current := rootID
	for current != "" {
		if _, seen := visited[current]; seen {
			// Ciclo detectado: corta sin propagar error.
			break
		}
		visited[current] = struct{}{}
		chain = append(chain, current)

		parent, ok, err := parentOf(current)
		if err != nil {
			return nil, err
		}
		if !ok || parent == "" {
			break
		}
		current = parent
	}
	return chain, nil
}
