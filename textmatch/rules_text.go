package textmatch

import (
	"strings"
	"unicode"
)

// unitMarks son los signos que no son puntuación sino unidad: «%» y el grado. La
// regla `punctuation` los conserva para que el comparador numérico los lea.
var unitMarks = map[rune]struct{}{'%': {}, '°': {}, 'º': {}}

// numberSeparators son los signos que, entre dos cifras, son parte del número
// («3,5», «1.000.000», «1/2») y no puntuación.
var numberSeparators = map[rune]struct{}{'.': {}, ',': {}, '/': {}}

// punctuationRule perdona los signos de puntuación: «Lima.» y «lima», «¡París!» y
// «parís». Conserva lo que es parte de un número o de su unidad: la coma, el punto
// y la barra entre dos cifras, el signo menos delante de una cifra, «%» y el grado.
type punctuationRule struct{}

func (punctuationRule) Name() string         { return RulePunctuation }
func (punctuationRule) DefaultEnabled() bool { return true }

func (punctuationRule) Canonical(s string) string {
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s))
	for i, r := range rs {
		if keepRune(rs, i, b.String()) {
			b.WriteRune(r)
			continue
		}
		b.WriteRune(' ')
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// keepRune decide si el rune i sobrevive a la regla de puntuación. out es lo
// escrito hasta ahora (para saber si el signo menos abre un token).
func keepRune(rs []rune, i int, out string) bool {
	r := rs[i]
	switch {
	case unicode.IsLetter(r), unicode.IsDigit(r), unicode.IsMark(r), unicode.IsSpace(r):
		return true
	}
	if _, ok := unitMarks[r]; ok {
		return true
	}
	nextIsDigit := i+1 < len(rs) && unicode.IsDigit(rs[i+1])
	if _, ok := numberSeparators[r]; ok {
		return i > 0 && unicode.IsDigit(rs[i-1]) && nextIsDigit
	}
	if r == '-' || r == '−' {
		opensToken := out == "" || strings.HasSuffix(out, " ")
		return opensToken && nextIsDigit
	}
	return false
}

// accentsRule perdona la tilde y la diéresis: «oxigeno» y «oxígeno». Es lo que
// hasta hoy hacía Normalize para todas las preguntas, ahora con nombre: la «ñ» se
// conserva (es letra, no tilde). Apagada, la tilde cuenta: «canté» no es «cante».
type accentsRule struct{}

func (accentsRule) Name() string         { return RuleAccents }
func (accentsRule) DefaultEnabled() bool { return true }

// Canonical delega en Normalize, que sobre un texto ya en minúsculas solo quita
// las marcas y conserva la «ñ».
func (accentsRule) Canonical(s string) string { return Normalize(s) }

// leadingArticles son los artículos que la regla `leading_article` quita al
// principio de una respuesta.
var leadingArticles = map[string]struct{}{
	"el": {}, "la": {}, "los": {}, "las": {},
	"un": {}, "una": {}, "unos": {}, "unas": {},
}

// leadingArticleRule perdona el artículo inicial: «la fotosíntesis» y
// «fotosíntesis». Una respuesta que es solo un artículo se queda como está.
// Implementa TargetGuard: si la respuesta esperada empieza por artículo, la docente
// lo escribió y se exige («la capital» no es «el capital»).
type leadingArticleRule struct{}

func (leadingArticleRule) Name() string         { return RuleLeadingArticle }
func (leadingArticleRule) DefaultEnabled() bool { return true }

func (leadingArticleRule) Canonical(s string) string {
	tokens := strings.Fields(s)
	for len(tokens) > 1 {
		if _, ok := leadingArticles[tokens[0]]; !ok {
			break
		}
		tokens = tokens[1:]
	}
	return strings.Join(tokens, " ")
}

// GuardsTarget es true si la esperada empieza por un artículo que la regla quitaría.
func (r leadingArticleRule) GuardsTarget(target string) bool {
	return r.Canonical(target) != strings.Join(strings.Fields(target), " ")
}

// genderNumberRule perdona el género y el número (ADR 0049, «Lo que decide»):
// plural en «-s» y «-es» («-ces» a «-z») y género en «-o» y «-a». Cada palabra de
// tres letras o más se lleva a una clave común: se quita el plural, se quita la
// vocal final («-a», «-o» o «-e», la que deja el plural en «-es») y se añade «-o».
// Así «blanca» y «blanco», «aves» y «ave», «lápices» y «lápiz» comparten clave, y
// la clave es un punto fijo (termina en «-o»): aplicarla dos veces da lo mismo.
// Nace apagada: encendida acepta también «ramo» por «rama», y eso lo decide la
// docente en cada pregunta.
type genderNumberRule struct{}

func (genderNumberRule) Name() string         { return RuleGenderNumber }
func (genderNumberRule) DefaultEnabled() bool { return false }

func (genderNumberRule) Canonical(s string) string {
	tokens := strings.Fields(s)
	for i, t := range tokens {
		tokens[i] = genderNumberKey(t)
	}
	return strings.Join(tokens, " ")
}

// genderNumberKey da la clave de una palabra; las que no son solo letras o tienen
// menos de tres, y aquellas cuya raíz quedaría de menos de dos letras, no cambian.
func genderNumberKey(word string) string {
	rs := []rune(word)
	if len(rs) < 3 {
		return word
	}
	for _, r := range rs {
		if !unicode.IsLetter(r) {
			return word
		}
	}
	stem := string(rs)
	switch {
	case strings.HasSuffix(stem, "ces"):
		stem = strings.TrimSuffix(stem, "ces") + "z"
	case strings.HasSuffix(stem, "es"):
		stem = strings.TrimSuffix(stem, "es")
	case strings.HasSuffix(stem, "s"):
		stem = strings.TrimSuffix(stem, "s")
	}
	if n := len(stem); n > 0 && strings.ContainsRune("aoe", rune(stem[n-1])) {
		stem = stem[:n-1]
	}
	if len([]rune(stem)) < 2 {
		return word
	}
	return stem + "o"
}
