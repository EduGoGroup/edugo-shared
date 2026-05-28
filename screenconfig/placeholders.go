package screenconfig

import "strings"

// ResourcePrefixFromPermission extrae el prefijo de recurso a partir
// de una permission del tipo "content.assessments.read" → "content.assessments".
// Si el permission no tiene segmentos separables (o esta vacio),
// devuelve string vacio: los placeholders quedaran sin sustituir, lo
// que es preferible a inyectar basura.
func ResourcePrefixFromPermission(permission string) string {
	if permission == "" {
		return ""
	}
	idx := strings.LastIndex(permission, ".")
	if idx <= 0 {
		return ""
	}
	return permission[:idx]
}

// ExpandResourcePlaceholders sustituye "$resource$" en cualquier valor
// string de los defaults por el prefix dado. Solo reemplaza si el
// prefix no es vacio; en otro caso devuelve el default sin tocar (el
// placeholder queda visible — falla loud para detectar misconfig).
func ExpandResourcePlaceholders(actions []map[string]any, prefix string) []map[string]any {
	if prefix == "" {
		return actions
	}
	out := make([]map[string]any, 0, len(actions))
	for _, a := range actions {
		cp := make(map[string]any, len(a))
		for k, v := range a {
			if s, ok := v.(string); ok && strings.Contains(s, "$resource$") {
				cp[k] = strings.ReplaceAll(s, "$resource$", prefix)
			} else {
				cp[k] = v
			}
		}
		out = append(out, cp)
	}
	return out
}
