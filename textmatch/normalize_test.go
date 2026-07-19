package textmatch

import (
	"reflect"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"minúsculas", "FaceBook", "facebook"},
		{"tildes fuera", "Café Perú", "cafe peru"},
		{"diéresis fuera", "Pingüino", "pinguino"},
		{"ñ preservada minúscula", "año", "año"},
		{"ñ preservada mayúscula", "AÑO", "año"},
		{"ano sin ñ intacto", "ano", "ano"},
		{"colapsa espacios y trim", "  hola   mundo  ", "hola mundo"},
		{"combinada", "  El Niño  ", "el niño"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.in); got != tt.want {
				t.Errorf("Normalize(%q) = %q, quiero %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestNormalizeAnioNoEsAno es la regresión canónica D-045.7: la «ñ» se preserva,
// así «año» y «ano» NO colapsan al mismo texto normalizado.
func TestNormalizeAnioNoEsAno(t *testing.T) {
	if Normalize("año") == Normalize("ano") {
		t.Fatal("«año» y «ano» no deben normalizar igual: la ñ es letra, no tilde")
	}
}

func TestSplitTokens(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "descarta conector y y relleno queda como tokens",
			in:   "whastapp instalgram y el famoso facebook",
			want: []string{"whastapp", "instalgram", "el", "famoso", "facebook"},
		},
		{
			name: "puntuación como frontera y conector y",
			in:   "los países son: ecuador, venezuela y colombia",
			want: []string{"los", "paises", "son", "ecuador", "venezuela", "colombia"},
		},
		{
			name: "conector e descartado",
			in:   "café e historia",
			want: []string{"cafe", "historia"},
		},
		{
			name: "barras y comas",
			in:   "a,b|c",
			want: []string{"a", "b", "c"},
		},
		{
			name: "vacío",
			in:   "   ",
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitTokens(tt.in)
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitTokens(%q) = %v, quiero %v", tt.in, got, tt.want)
			}
		})
	}
}
