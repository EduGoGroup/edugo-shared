package textmatch

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// nSentinel es un rune del área de uso privado Unicode que sustituye a la «ñ»
// mientras se descomponen y quitan los diacríticos. Como no es una marca (Mn) ni
// tiene descomposición, sobrevive intacto a la cadena NFD→remove(Mn)→NFC y se
// restaura después. Así la «ñ» se preserva como letra propia (D-045.7: «año»≠«ano»)
// aunque el pipeline de tildes, sin protección, la borraría (NFD la parte en
// «n» + tilde combinante U+0303, que es Mn).
const nSentinel = "\uE000"

// connectorTokens son las palabras conectoras que separan ítems en la respuesta
// del alumno pero no forman parte de ningún ítem («ecuador y colombia» → dos
// ítems). Se descartan tras el split; nunca casan un ítem. Portado de
// edugo-api-learning/internal/core/domain/scoring_prep.go.
var connectorTokens = map[string]struct{}{"y": {}, "e": {}}

// nonAlnum casa runs de caracteres que no son letra ni dígito (unicode). Toda la
// puntuación y los separadores del contrato («,» «|» «;») caen aquí y actúan como
// frontera de token, así "ecuador, venezuela" no pega "ecuador," ≠ "ecuador".
// Portado del splitPattern más correcto (scoring_prep.go), no del regex del worker
// que produce fragmentos pegados.
var nonAlnum = regexp.MustCompile(`[^\p{L}\p{N}]+`)

// diacriticStripper quita las marcas diacríticas (tildes, diéresis) descomponiendo,
// removiendo las marcas (categoría Unicode Mn) y recomponiendo.
var diacriticStripper = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

// Normalize es el contrato de normalización canónico del ecosistema (D-045.7):
// minúsculas, sin tildes/diéresis, PRESERVA la «ñ» (es letra, no tilde), colapsa
// los espacios internos y hace trim. Es una función pura. Reemplaza al
// removeDiacritics de learning (que borraba la «ñ») y al accentMap/normalize del
// worker (que la preservaba a mano); unifica la divergencia histórica de la «ñ».
func Normalize(s string) string {
	s = strings.ToLower(s)
	// Compone primero (n + tilde combinante → «ñ» precompuesta) para poder
	// protegerla aunque venga descompuesta desde el origen.
	s = norm.NFC.String(s)
	s = strings.ReplaceAll(s, "ñ", nSentinel)
	// La cadena NFD→remove(Mn)→NFC no falla para texto UTF-8 válido; ante un error
	// improbable se conserva la entrada sin descomponer (degradación segura).
	if stripped, _, err := transform.String(diacriticStripper, s); err == nil {
		s = stripped
	}
	s = strings.ReplaceAll(s, nSentinel, "ñ")
	// strings.Fields colapsa cualquier run de espacios y descarta los extremos.
	return strings.Join(strings.Fields(s), " ")
}

// SplitTokens descompone un texto en tokens normalizados para el match de listas
// (D-045.7). Es pura. Reglas: aplica Normalize, usa como frontera todo carácter
// no alfanumérico unicode (incluida la puntuación), descarta las conectoras «y»/«e»
// (separan ítems, no son ítems) y colapsa los tokens vacíos. Portado de
// scoring_prep.go:SplitAndNormalizeAnswer.
func SplitTokens(s string) []string {
	normalized := Normalize(s)
	raw := nonAlnum.Split(normalized, -1)
	out := make([]string, 0, len(raw))
	for _, tok := range raw {
		if tok == "" {
			continue
		}
		if _, isConnector := connectorTokens[tok]; isConnector {
			continue
		}
		out = append(out, tok)
	}
	return out
}
