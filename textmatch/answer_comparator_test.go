package textmatch

import (
	"context"
	"errors"
	"testing"
)

type comparatorCase struct {
	expected, candidate string
	want                Outcome
}

func runComparator(t *testing.T, overrides map[string]bool, cases []comparatorCase) {
	t.Helper()
	c, err := NewAnswerComparator(overrides)
	if err != nil {
		t.Fatalf("NewAnswerComparator(%v): %v", overrides, err)
	}
	for _, tc := range cases {
		r, err := c.Compare(context.Background(), tc.expected, tc.candidate)
		if err != nil {
			t.Fatalf("Compare(%q, %q): %v", tc.expected, tc.candidate, err)
		}
		if r.Outcome != tc.want {
			t.Errorf("reglas %v: %q / %q = %v (%s), quiero %v",
				overrides, tc.expected, tc.candidate, r.Outcome, r.Evidence, tc.want)
		}
		if r.Strategy != "rules" || r.Evidence == "" {
			t.Errorf("%q / %q: falta la procedencia (%q, %q)", tc.expected, tc.candidate, r.Strategy, r.Evidence)
		}
	}
}

func TestAnswerComparatorDefaults(t *testing.T) {
	runComparator(t, nil, []comparatorCase{
		// Casa.
		{"fotosíntesis", "La fotosíntesis.", OutcomeMatch},
		{"Lima", "Lima.", OutcomeMatch},
		{"oxígeno", "oxigeno", OutcomeMatch},
		{"Sol", "El Sol.", OutcomeMatch},
		{"3", "Tres.", OutcomeMatch},
		{"ocho", "8", OutcomeMatch},
		{"100 °C", "100 grados Celsius", OutcomeMatch},
		{"5 cm", "5 centímetros", OutcomeMatch},
		{"1000000", "1.000.000", OutcomeMatch},
		// No se da por buena: una regla apagada lo exige.
		{"maestra", "maestro", OutcomeNoMatch},
		{"hermana", "hermano", OutcomeNoMatch},
		{"la capital", "el capital", OutcomeNoMatch},
		{"El Salvador", "Salvador", OutcomeNoMatch},
		{"el Pacífico", "Pacífico", OutcomeNoMatch},
		// Un número decide sin cercanía.
		{"1234567", "1234568", OutcomeNoMatch},
		{"5 cm", "5 m", OutcomeNoMatch},
		{"1000 g", "1000 kg", OutcomeNoMatch},
		{"12", "12 kg", OutcomeNoMatch},
		{"24", "veintitrés", OutcomeNoMatch},
		{"100 °C", "100 °F", OutcomeNoMatch},
		// No opina: la cercanía de siempre decide.
		{"oxígeno", "oxijeno", OutcomeUncertain},
		{"Austria", "Australia", OutcomeUncertain},
		{"auto", "coche", OutcomeUncertain},
		// Una respuesta vacía no casa con nada.
		{"Lima", "", OutcomeNoMatch},
		{"Sol", ".", OutcomeNoMatch},
		{"3", "-", OutcomeNoMatch},
		{"París", "   ", OutcomeNoMatch},
	})
}

func TestAnswerComparatorAccentsOff(t *testing.T) {
	runComparator(t, map[string]bool{RuleAccents: false}, []comparatorCase{
		{"canté", "cante", OutcomeNoMatch},
		{"habló", "hablo", OutcomeNoMatch},
		{"canté", "Canté.", OutcomeMatch},
		{"5 centímetros", "5 centimetros", OutcomeMatch},
	})
}

func TestAnswerComparatorGenderOn(t *testing.T) {
	runComparator(t, map[string]bool{RuleGenderNumber: true}, []comparatorCase{
		{"blanco", "blanca", OutcomeMatch},
		{"blanco", "Blanca.", OutcomeMatch},
		{"pulmones", "el pulmón", OutcomeMatch},
		{"aves", "ave", OutcomeMatch},
		{"blanco", "negro", OutcomeUncertain},
		{"5 cm", "5 centímetros", OutcomeMatch},
	})
}

func TestAnswerComparatorRulesOff(t *testing.T) {
	runComparator(t, map[string]bool{RulePunctuation: false}, []comparatorCase{
		{"Lima", "Lima.", OutcomeNoMatch},
		{"10", "10.", OutcomeNoMatch},
	})
	runComparator(t, map[string]bool{RuleNumberWords: false}, []comparatorCase{
		{"3", "tres", OutcomeNoMatch},
		{"ocho", "8", OutcomeNoMatch},
	})
	runComparator(t, map[string]bool{RuleLeadingArticle: false}, []comparatorCase{
		{"Sol", "el Sol", OutcomeNoMatch},
	})
}

func TestAnswerComparatorUnknownRule(t *testing.T) {
	if _, err := NewAnswerComparator(map[string]bool{"plural": true}); !errors.Is(err, ErrUnknownRule) {
		t.Fatalf("una regla desconocida da %v, quiero ErrUnknownRule", err)
	}
}

// Las sondas de fuera del banco que lista el informe de las reglas, en su
// dirección (esperada / alumno) y con el resultado que dice allí. «Pendiente» es no
// darla por buena: NoMatch, o Uncertain cuando la cercanía tampoco llega
// (TestReportProbesStayPending).
func TestAnswerComparatorReportProbes(t *testing.T) {
	runComparator(t, nil, []comparatorCase{
		{"oxígeno", "oxígena", OutcomeNoMatch},
		{"neuronas", "neurona", OutcomeNoMatch},
		{"la circunferencia", "el circunferencia", OutcomeNoMatch},
		{"12", "12 kg", OutcomeNoMatch},
		{"3,5", "3.5", OutcomeMatch},
		{"1/2", "0,5", OutcomeMatch},
		{"veintiuno", "21", OutcomeMatch},
		{"cien", "100", OutcomeMatch},
		{"ciento uno", "101", OutcomeUncertain},
		{"101", "ciento uno", OutcomeNoMatch},
		{"8848", "8.848", OutcomeMatch},
	})
	runComparator(t, map[string]bool{RuleGenderNumber: true}, []comparatorCase{
		{"oxígeno", "oxígena", OutcomeMatch},
		{"neuronas", "neurona", OutcomeMatch},
		{"la circunferencia", "el circunferencia", OutcomeNoMatch},
		{"12", "12 kg", OutcomeNoMatch},
		{"3,5", "3.5", OutcomeMatch},
		{"1/2", "0,5", OutcomeMatch},
		{"veintiuno", "21", OutcomeMatch},
		{"cien", "100", OutcomeMatch},
		{"8848", "8.848", OutcomeMatch},
		{"ciento uno", "101", OutcomeUncertain},
		{"101", "ciento uno", OutcomeNoMatch},
	})
}

// Las sondas que el comparador no opina tampoco las da por buenas la cercanía.
func TestReportProbesStayPending(t *testing.T) {
	r, err := NewFuzzy(0).Compare(context.Background(), "ciento uno", "101")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outcome == OutcomeMatch {
		t.Errorf("la cercanía da por buena «101» por «ciento uno» (%s)", r.Evidence)
	}
}

// Con una regla apagada, la igualdad gana al vacío: si la docente exige la
// puntuación, «?» por «?» casa.
func TestAnswerComparatorEqualBeatsEmpty(t *testing.T) {
	runComparator(t, map[string]bool{RulePunctuation: false}, []comparatorCase{
		{"?", "?", OutcomeMatch},
		{"Lima", "...", OutcomeNoMatch},
	})
}

// Cada texto esperado se mira por su lado: una alternativa numérica se compara
// como número y una que empieza por artículo lo exige.
func TestAnswerComparatorLooksAtEachTarget(t *testing.T) {
	runComparator(t, nil, []comparatorCase{
		{"12 kg", "12", OutcomeNoMatch},
		{"12 kg", "doce kilogramos", OutcomeMatch},
		{"la capital", "capital", OutcomeNoMatch},
		{"capital", "la capital", OutcomeMatch},
	})
}
