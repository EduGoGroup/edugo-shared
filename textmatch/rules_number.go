package textmatch

import (
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// numberWords son los números en letra del 0 al 29, las decenas y «cien», con y
// sin tilde (la regla `accents` va antes y puede estar apagada). «un» y «una» no
// están: son artículos (ficha HU-INF.4, decisión 6).
var numberWords = map[string]int{
	"cero": 0, "uno": 1, "dos": 2, "tres": 3, "cuatro": 4, "cinco": 5, "seis": 6,
	"siete": 7, "ocho": 8, "nueve": 9, "diez": 10, "once": 11, "doce": 12,
	"trece": 13, "catorce": 14, "quince": 15,
	"dieciseis": 16, "dieciséis": 16, "diecisiete": 17, "dieciocho": 18, "diecinueve": 19,
	"veinte": 20, "veintiuno": 21, "veintidos": 22, "veintidós": 22,
	"veintitres": 23, "veintitrés": 23, "veinticuatro": 24, "veinticinco": 25,
	"veintiseis": 26, "veintiséis": 26, "veintisiete": 27, "veintiocho": 28, "veintinueve": 29,
	"treinta": 30, "cuarenta": 40, "cincuenta": 50, "sesenta": 60,
	"setenta": 70, "ochenta": 80, "noventa": 90, "cien": 100,
}

// largeNumberWords marcan un número de más de cien en letra. Una palabra de
// numberWords pegada a una de ellas no se convierte: «ciento uno» no es «ciento 1».
var largeNumberWords = map[string]struct{}{
	"ciento": {}, "cientos": {}, "doscientos": {}, "trescientos": {}, "cuatrocientos": {},
	"quinientos": {}, "seiscientos": {}, "setecientos": {}, "ochocientos": {},
	"novecientos": {}, "mil": {}, "millon": {}, "millón": {}, "millones": {},
}

// numberWordsRule perdona la cifra y la letra del 0 al 100: «once» y «11»,
// «treinta y uno» y «31». Lleva la letra a la cifra; fuera de ese rango no toca nada.
type numberWordsRule struct{}

func (numberWordsRule) Name() string         { return RuleNumberWords }
func (numberWordsRule) DefaultEnabled() bool { return true }

func (numberWordsRule) Canonical(s string) string {
	tokens := strings.Fields(s)
	out := make([]string, 0, len(tokens))
	for i := 0; i < len(tokens); {
		value, used := readNumberWords(tokens, i)
		if used == 0 {
			out = append(out, tokens[i])
			i++
			continue
		}
		out = append(out, strconv.Itoa(value))
		i += used
	}
	return strings.Join(out, " ")
}

// readNumberWords lee un número en letra que empieza en tokens[i] y devuelve su
// valor y cuántos tokens usa (0 si no hay número o si va pegado a uno mayor).
func readNumberWords(tokens []string, i int) (int, int) {
	value, ok := numberWords[tokens[i]]
	if !ok {
		return 0, 0
	}
	used := 1
	if value >= 30 && value <= 90 && i+2 < len(tokens) && tokens[i+1] == "y" {
		if unit, ok := numberWords[tokens[i+2]]; ok && unit >= 1 && unit <= 9 {
			value += unit
			used = 3
		}
	}
	if isLargeNumberWord(tokens, i-1) || isLargeNumberWord(tokens, i+used) {
		return 0, 0
	}
	return value, used
}

func isLargeNumberWord(tokens []string, i int) bool {
	if i < 0 || i >= len(tokens) {
		return false
	}
	_, ok := largeNumberWords[tokens[i]]
	return ok
}

// unitAliases es la tabla cerrada de unidades (ficha HU-INF.4, decisión 2): cada
// unidad con su símbolo y su nombre en singular y plural. «grados» a secas no
// entra: también son ángulos.
var unitAliases = map[string][]string{
	"mm":  {"mm", "milímetro", "milímetros"},
	"cm":  {"cm", "centímetro", "centímetros"},
	"m":   {"m", "metro", "metros"},
	"km":  {"km", "kilómetro", "kilómetros"},
	"mg":  {"mg", "miligramo", "miligramos"},
	"g":   {"g", "gramo", "gramos"},
	"kg":  {"kg", "kilogramo", "kilogramos"},
	"ml":  {"ml", "mililitro", "mililitros"},
	"l":   {"l", "litro", "litros"},
	"s":   {"s", "segundo", "segundos"},
	"min": {"min", "minuto", "minutos"},
	"h":   {"h", "hora", "horas"},
	"°c": {
		"°c", "ºc", "° c", "º c", "grado celsius", "grados celsius",
		"grado centígrado", "grados centígrados",
	},
	"%": {"%", "por ciento"},
}

// quantityPattern separa la cifra de lo que la sigue: «5cm», «100 °c», «1 000 000».
var quantityPattern = regexp.MustCompile(`^(-?\d+(?:[.,/ ]\d+)*)\s*(.*)$`)

var (
	thousandsDot      = regexp.MustCompile(`^\d{1,3}(\.\d{3})+$`)
	thousandsDotComma = regexp.MustCompile(`^\d{1,3}(\.\d{3})+,\d+$`)
	thousandsSpace    = regexp.MustCompile(`^\d{1,3}( \d{3})+$`)
)

// quantity es un número con su unidad ("" si no lleva).
type quantity struct {
	value *big.Rat
	unit  string
}

func (q quantity) equal(o quantity) bool {
	return q.unit == o.unit && q.value.Cmp(o.value) == 0
}

// unitTable resuelve un nombre de unidad, en la forma canónica de las reglas
// encendidas, a su unidad.
type unitTable map[string]string

// newUnitTable pasa cada alias, con su tilde y sin ella, por la misma cadena de
// reglas que la respuesta, así «centímetros» resuelve con cualquier combinación.
func newUnitTable(canonical func(string) string) unitTable {
	table := make(unitTable)
	for unit, aliases := range unitAliases {
		for _, alias := range aliases {
			table[canonical(alias)] = unit
			table[canonical(Normalize(alias))] = unit
		}
	}
	return table
}

// parseQuantity lee un texto canónico como número con unidad opcional. ok es
// false si no es un número o si lo que sigue a la cifra no es una unidad de la tabla.
func parseQuantity(s string, units unitTable) (quantity, bool) {
	m := quantityPattern.FindStringSubmatch(s)
	if m == nil {
		return quantity{}, false
	}
	value, ok := parseNumber(m[1])
	if !ok {
		return quantity{}, false
	}
	rest := strings.TrimSpace(m[2])
	if rest == "" {
		return quantity{value: value}, true
	}
	unit, ok := units[rest]
	if !ok {
		return quantity{}, false
	}
	return quantity{value: value, unit: unit}, true
}

// parseNumber lee la cifra (ficha HU-INF.4, decisión 3): la coma y el punto valen
// como decimal; un punto, o un espacio, con grupos exactos de tres cifras separa
// los miles («1.000.000», «8.848», «1 000 000»); una fracción vale por su valor.
// Compara con aritmética exacta: «1/2» es «0,5».
func parseNumber(s string) (*big.Rat, bool) {
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign, s = "-", s[1:]
	}
	var literal string
	switch {
	case strings.Contains(s, " "):
		if !thousandsSpace.MatchString(s) {
			return nil, false
		}
		literal = strings.ReplaceAll(s, " ", "")
	case thousandsDot.MatchString(s):
		literal = strings.ReplaceAll(s, ".", "")
	case thousandsDotComma.MatchString(s):
		literal = strings.Replace(strings.ReplaceAll(s, ".", ""), ",", ".", 1)
	case strings.Count(s, ".")+strings.Count(s, ",")+strings.Count(s, "/") > 1:
		return nil, false
	case strings.Contains(s, "/"):
		// big.Rat rechaza el denominador cero.
		literal = s
	default:
		literal = strings.Replace(s, ",", ".", 1)
	}
	r, ok := new(big.Rat).SetString(sign + literal)
	return r, ok
}
