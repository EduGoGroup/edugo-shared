package textmatch

import (
	"errors"
	"fmt"
)

// Nombres de las reglas de lengua de la respuesta corta (ADR 0049, punto 2). El
// nombre es parte del contrato: es la clave de answer_spec.rules y no cambia.
const (
	RulePunctuation    = "punctuation"
	RuleAccents        = "accents"
	RuleNumberWords    = "number_words"
	RuleLeadingArticle = "leading_article"
	RuleGenderNumber   = "gender_number"
)

// ErrUnknownRule es el error de pedir una regla que el registro no conoce.
var ErrUnknownRule = errors.New("textmatch: unknown rule")

// Rule es el contrato de una regla de lengua (ADR 0049, punto 2): un nombre
// estable, un valor por defecto y la forma canónica de UN texto. No conoce la otra
// respuesta; por eso las reglas se encadenan. Encendida, perdona la diferencia que
// borra; apagada, la exige.
//
// Canonical recibe un texto ya pasado por NormalizeBase y por las reglas que van
// antes en el orden, y devuelve un texto con los espacios colapsados. Es pura e
// idempotente: aplicarla dos veces da lo mismo que una.
type Rule interface {
	Name() string
	DefaultEnabled() bool
	Canonical(s string) string
}

// TargetGuard es opcional. Una regla que lo implementa queda apagada al comparar
// contra un texto esperado que ya trae lo que ella perdonaría: si la docente lo
// escribió, cuenta. target llega en la forma canónica de las reglas encendidas que
// van antes. Es lo que hace `leading_article` (ADR 0049, punto 3).
type TargetGuard interface {
	GuardsTarget(target string) bool
}

// registry es el registro de reglas, en el orden fijo en que se encadenan, que es
// parte del contrato: puntuación, tilde, cifras, artículo, género.
var registry = []Rule{
	punctuationRule{},
	accentsRule{},
	numberWordsRule{},
	leadingArticleRule{},
	genderNumberRule{},
}

// Rules devuelve las reglas registradas, en su orden. La lista es una copia.
func Rules() []Rule {
	return append([]Rule(nil), registry...)
}

// RuleNames devuelve los nombres registrados, en el orden de las reglas.
func RuleNames() []string {
	names := make([]string, len(registry))
	for i, r := range registry {
		names[i] = r.Name()
	}
	return names
}

// LookupRule devuelve la regla con ese nombre, o un error que envuelve
// ErrUnknownRule si el registro no la conoce.
func LookupRule(name string) (Rule, error) {
	for _, r := range registry {
		if r.Name() == name {
			return r, nil
		}
	}
	return nil, fmt.Errorf("%w: %q", ErrUnknownRule, name)
}

// ResolveRules aplica a los valores por defecto lo que la pregunta cambió
// (answer_spec.rules: solo lo que la docente tocó) y devuelve el estado de cada
// regla registrada. Un nombre que no está en el registro es un error que envuelve
// ErrUnknownRule.
func ResolveRules(overrides map[string]bool) (map[string]bool, error) {
	state := make(map[string]bool, len(registry))
	for _, r := range registry {
		state[r.Name()] = r.DefaultEnabled()
	}
	for name, on := range overrides {
		if _, err := LookupRule(name); err != nil {
			return nil, err
		}
		state[name] = on
	}
	return state, nil
}
