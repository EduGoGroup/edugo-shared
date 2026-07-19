package textmatch

import "testing"

func TestEditDistance(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{"iguales", "facebook", "facebook", 0},
		{"vacío contra algo", "", "abc", 3},
		{"algo contra vacío", "abc", "", 3},
		{"clásico kitten/sitting", "kitten", "sitting", 3},
		{"inserción instalgram/instagram", "instalgram", "instagram", 1},
		// Transposición adyacente s<->t: Damerau-Levenshtein (OSA) la cuenta 1,
		// no 2 como el Levenshtein puro. Es la desviación clave del plan 045.
		{"transposición whastapp/whatsapp", "whastapp", "whatsapp", 1},
		{"por runas año/ano", "año", "ano", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EditDistance(tt.a, tt.b); got != tt.want {
				t.Errorf("EditDistance(%q,%q) = %d, quiero %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestEditDistanceSimétrica confirma que la distancia no depende del orden.
func TestEditDistanceSimétrica(t *testing.T) {
	pairs := [][2]string{{"whastapp", "whatsapp"}, {"instalgram", "instagram"}, {"costa", "rica"}}
	for _, p := range pairs {
		if EditDistance(p[0], p[1]) != EditDistance(p[1], p[0]) {
			t.Errorf("EditDistance no simétrica para %q/%q", p[0], p[1])
		}
	}
}
