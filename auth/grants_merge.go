package auth

// MergeGrantChain funde una cadena de Grants (uno por nivel de la
// jerarquía de roles) en un único Grants plano. El orden esperado es
// "ancestro primero, hijo después" pero el resultado es invariante al
// orden: se unen todos los allow y todos los deny y se deduplican los
// patterns repetidos preservando la primera aparición.
//
// El aplanado es deliberado: las APIs y el KMP siguen recibiendo la
// misma estructura plana Grants{Allow, Deny} de hoy. La precedencia
// deny-wins NO se aplica aquí (no se descartan allows que algún deny
// cubra): eso es responsabilidad del matcher (EvaluateGrants), que es
// set-based y evalúa por request. Aplanar aquí mantiene intacta esa
// semántica — un deny del hijo gana sobre un allow del ancestro al
// evaluar, sin importar el nivel del que provengan.
func MergeGrantChain(chain []Grants) Grants {
	out := Grants{Allow: []string{}, Deny: []string{}}
	seenAllow := make(map[string]struct{})
	seenDeny := make(map[string]struct{})
	for _, g := range chain {
		for _, p := range g.Allow {
			if _, ok := seenAllow[p]; ok {
				continue
			}
			seenAllow[p] = struct{}{}
			out.Allow = append(out.Allow, p)
		}
		for _, p := range g.Deny {
			if _, ok := seenDeny[p]; ok {
				continue
			}
			seenDeny[p] = struct{}{}
			out.Deny = append(out.Deny, p)
		}
	}
	return out
}
