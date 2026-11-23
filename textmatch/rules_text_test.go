package textmatch

import "testing"

// canonicalOf pasa un texto por NormalizeBase y una sola regla.
func canonicalOf(r Rule, s string) string { return r.Canonical(NormalizeBase(s)) }

type ruleCase struct {
	a, b string
	same bool
}

func checkRule(t *testing.T, r Rule, cases []ruleCase) {
	t.Helper()
	for _, c := range cases {
		ca, cb := canonicalOf(r, c.a), canonicalOf(r, c.b)
		if (ca == cb) != c.same {
			t.Errorf("%s: %q→%q y %q→%q, iguales=%v, quiero %v", r.Name(), c.a, ca, c.b, cb, ca == cb, c.same)
		}
	}
}

func TestPunctuationRule(t *testing.T) {
	checkRule(t, punctuationRule{}, []ruleCase{
		{"Lima.", "lima", true},
		{"¡París!", "parís", true},
		{"París...", "parís", true},
		{`"Roma"`, "roma", true},
		{"(agua)", "agua", true},
		{"Simón Bolívar!!", "simón bolívar", true},
		{"10.", "10", true},
		{"1/2", "12", false},
		{"3,5", "35", false},
		{"1.0", "10", false},
		{"1.000.000", "1000000", false},
	})
	keep := map[string]string{
		"1/2":     "1/2",
		"3,5.":    "3,5",
		"-5 °C":   "-5 °c",
		"50 %":    "50 %",
		"(-5)":    "-5",
		"a-5":     "a 5",
		"¿Qué?":   "qué",
		"8.848":   "8.848",
		"100°C":   "100°c",
		"O'Brien": "o brien",
	}
	for in, want := range keep {
		if got := canonicalOf(punctuationRule{}, in); got != want {
			t.Errorf("punctuation(%q) = %q, quiero %q", in, got, want)
		}
	}
}

func TestAccentsRule(t *testing.T) {
	checkRule(t, accentsRule{}, []ruleCase{
		{"oxigeno", "oxígeno", true},
		{"canté", "cante", true},
		{"pingüino", "pinguino", true},
		{"año", "ano", false},
	})
}

// La base nueva más la tilde da exactamente Normalize: lo que hoy vale sigue valiendo.
func TestNormalizeBasePlusAccentsIsNormalize(t *testing.T) {
	for _, s := range []string{
		"Lima", "  Ciudad   de México ", "oxígeno", "Bogotá", "BOGOTÁ", "córazon", "año", "AÑO",
		"ño", "pingüino", "canté", "¡París!", "1.000.000", "100 °C", "Cristóbal Colón.", "",
		"Leonardo da Vinci", "Venezuela; Ecuador; Colombia.", "Ñandú", "éxito",
	} {
		if got, want := (accentsRule{}).Canonical(NormalizeBase(s)), Normalize(s); got != want {
			t.Errorf("accents(NormalizeBase(%q)) = %q, Normalize = %q", s, got, want)
		}
	}
}

func TestNormalizeBaseKeepsAccents(t *testing.T) {
	if got := NormalizeBase("  Canté   ÁRBOL "); got != "canté árbol" {
		t.Errorf("NormalizeBase = %q", got)
	}
}

func TestLeadingArticleRule(t *testing.T) {
	r := leadingArticleRule{}
	checkRule(t, r, []ruleCase{
		{"la fotosíntesis", "fotosíntesis", true},
		{"el Sol", "sol", true},
		{"los pulmones", "pulmones", true},
		{"unos gatos", "gatos", true},
		{"una célula", "célula", true},
		{"la", "", false},
		{"Ciudad de la Plata", "ciudad de plata", false},
	})
	if r.GuardsTarget(NormalizeBase("fotosíntesis")) {
		t.Error("una esperada sin artículo no apaga la regla")
	}
	for _, target := range []string{"la capital", "el pacífico", "el salvador"} {
		if !r.GuardsTarget(target) {
			t.Errorf("la esperada %q empieza por artículo: debe apagar la regla", target)
		}
	}
}

func TestGenderNumberRule(t *testing.T) {
	checkRule(t, genderNumberRule{}, []ruleCase{
		{"blanca", "blanco", true},
		{"aves", "ave", true},
		{"verdes", "verde", true},
		{"flores", "flor", true},
		{"lapices", "lapiz", true},
		{"lápices", "lápiz", true},
		{"pulmones", "pulmon", true},
		{"cachorros", "cachorro", true},
		{"maestra", "maestro", true},
		{"ramo", "rama", true},
		{"casas", "casa", true},
		{"sol", "luna", false},
		{"blanco", "negro", false},
	})
}

// Aplicar una regla dos veces da lo mismo que una.
func TestRulesAreIdempotent(t *testing.T) {
	inputs := []string{
		"Lima.", "¡París!", "1.0.", "-5 °C", "(-5)", "1/2.", "oxígeno", "canté",
		"la la casa", "el Sol", "los", "treinta y uno", "ciento uno", "veintitrés",
		"mes", "meses", "casas", "Amazonas", "lápices", "aves", "flores", "blanca",
		"mitosis", "pulmón", "5 cm", "100 grados Celsius", "",
	}
	for _, r := range Rules() {
		for _, in := range inputs {
			once := canonicalOf(r, in)
			if twice := r.Canonical(once); twice != once {
				t.Errorf("%s no es idempotente con %q: %q y luego %q", r.Name(), in, once, twice)
			}
		}
	}
}
