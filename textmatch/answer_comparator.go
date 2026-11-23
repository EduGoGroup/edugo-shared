package textmatch

import (
	"context"
	"fmt"
	"strings"
)

// AnswerComparator es el comparador que arma una pregunta corta a partir de su
// mapa de reglas (ADR 0049, puntos 3 a 6). Implementa Strategy y da uno de tres
// resultados contra cada texto esperado:
//
//   - OutcomeMatch («casa»): las formas canónicas con las reglas encendidas son
//     iguales y no vacías; o, si el esperado es un número, los dos valen lo mismo
//     con la misma unidad.
//   - OutcomeNoMatch («no se da por buena»): solo difieren en lo que perdonaría una
//     regla apagada, que la pregunta exige; o el esperado es un número y no vale lo
//     mismo; o la respuesta está vacía. No pasa a la cercanía.
//   - OutcomeUncertain («no opina»): en cualquier otro caso. Decide quien llama,
//     con la cercanía de siempre (Fuzzy) sobre el texto crudo.
//
// No va dentro de una Cascade delante de Fuzzy: la cascada escala también en
// OutcomeNoMatch, y eso deshace lo que la regla apagada exige. Quien lo usa
// compone a mano: Match cierra, NoMatch no se da por buena, Uncertain va a Fuzzy.
type AnswerComparator struct {
	enabled map[string]bool
	units   unitTable
}

// NewAnswerComparator arma el comparador de una pregunta. overrides es lo que la
// docente cambió (answer_spec.rules); lo demás vale su valor por defecto. Un
// nombre que no está en el registro es un error que envuelve ErrUnknownRule.
func NewAnswerComparator(overrides map[string]bool) (*AnswerComparator, error) {
	enabled, err := ResolveRules(overrides)
	if err != nil {
		return nil, err
	}
	c := &AnswerComparator{enabled: enabled}
	c.units = newUnitTable(func(s string) string {
		form, _ := c.canonical(s, nil)
		return form
	})
	return c, nil
}

// Name identifica la estrategia (procedencia en Result.Strategy).
func (*AnswerComparator) Name() string { return "rules" }

// Enabled dice si la regla está encendida en esta pregunta.
func (c *AnswerComparator) Enabled(name string) bool { return c.enabled[name] }

// Compare compara la respuesta del alumno (candidate) con UN texto esperado: la
// respuesta esperada o una alternativa, cada una por su lado (ADR 0049, punto 5).
// Nunca devuelve error.
func (c *AnswerComparator) Compare(_ context.Context, expected, candidate string) (Result, error) {
	targetForm, guarded := c.canonical(expected, nil)
	candidateForm, _ := c.canonical(candidate, guarded)

	if target, ok := parseQuantity(targetForm, c.units); ok {
		return c.compareNumbers(target, candidateForm), nil
	}

	candidateAll := canonicalAll(candidate)
	switch {
	case strings.TrimSpace(targetForm) == "":
		return c.result(OutcomeUncertain, 0, "la esperada queda vacía tras las reglas"), nil
	case candidateForm == "" || candidateAll == "":
		return c.result(OutcomeNoMatch, 0, "respuesta vacía"), nil
	case candidateForm == targetForm:
		return c.result(OutcomeMatch, 1, "iguales con las reglas encendidas: "+strings.Join(c.on(guarded), ", ")), nil
	case candidateAll == canonicalAll(expected):
		return c.result(OutcomeNoMatch, 0, "solo difieren en lo que la pregunta exige: "+strings.Join(c.off(guarded), ", ")), nil
	default:
		return c.result(OutcomeUncertain, 0, "las reglas no deciden"), nil
	}
}

// compareNumbers decide sin cercanía (ADR 0049, punto 6): la unidad no es opcional.
func (c *AnswerComparator) compareNumbers(target quantity, candidateForm string) Result {
	got, ok := parseQuantity(candidateForm, c.units)
	if !ok {
		return c.result(OutcomeNoMatch, 0, "número: la respuesta no es un número con una unidad conocida")
	}
	if !got.equal(target) {
		return c.result(OutcomeNoMatch, 0, fmt.Sprintf("número: %s %s no es %s %s",
			got.value.RatString(), got.unit, target.value.RatString(), target.unit))
	}
	return c.result(OutcomeMatch, 1, "número: mismo valor y misma unidad")
}

// canonical pasa un texto por las reglas encendidas, en su orden. Con guarded nil
// es el texto esperado: devuelve además las reglas que su forma apaga
// (TargetGuard). Con guarded, salta esas reglas: la respuesta del alumno se
// compara con las mismas reglas que el esperado.
func (c *AnswerComparator) canonical(s string, guarded map[string]bool) (string, map[string]bool) {
	isTarget := guarded == nil
	if isTarget {
		guarded = map[string]bool{}
	}
	form := NormalizeBase(s)
	for _, r := range registry {
		if !c.enabled[r.Name()] || guarded[r.Name()] {
			continue
		}
		if g, ok := r.(TargetGuard); ok && isTarget && g.GuardsTarget(form) {
			guarded[r.Name()] = true
			continue
		}
		form = r.Canonical(form)
	}
	return form, guarded
}

// canonicalAll pasa un texto por todas las reglas, también las apagadas.
func canonicalAll(s string) string {
	form := NormalizeBase(s)
	for _, r := range registry {
		form = r.Canonical(form)
	}
	return form
}

// on y off listan, en su orden, las reglas encendidas y las apagadas (o apagadas
// por el esperado) de una comparación.
func (c *AnswerComparator) on(guarded map[string]bool) []string {
	var names []string
	for _, r := range registry {
		if c.enabled[r.Name()] && !guarded[r.Name()] {
			names = append(names, r.Name())
		}
	}
	return names
}

func (c *AnswerComparator) off(guarded map[string]bool) []string {
	var names []string
	for _, r := range registry {
		if !c.enabled[r.Name()] || guarded[r.Name()] {
			names = append(names, r.Name())
		}
	}
	return names
}

func (c *AnswerComparator) result(outcome Outcome, confidence float64, evidence string) Result {
	return Result{Outcome: outcome, Confidence: confidence, Evidence: evidence, Strategy: c.Name()}
}
