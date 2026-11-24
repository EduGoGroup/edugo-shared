package textmatch

import (
	"math/big"
	"testing"
)

func TestNumberWordsRule(t *testing.T) {
	r := numberWordsRule{}
	cases := map[string]string{
		"once":              "11",
		"veintitrés":        "23",
		"veintitres":        "23",
		"dieciseis":         "16",
		"dieciséis":         "16",
		"cien":              "100",
		"cero":              "0",
		"uno":               "1",
		"treinta y uno":     "31",
		"noventa y nueve":   "99",
		"treinta y cero":    "30 y 0",
		"ciento uno":        "ciento uno",
		"dos mil":           "dos mil",
		"un gato":           "un gato",
		"una célula":        "una célula",
		"los tres cerditos": "los 3 cerditos",
		"cinco centímetros": "5 centímetros",
		"doce meses":        "12 meses",
	}
	for in, want := range cases {
		if got := canonicalOf(r, in); got != want {
			t.Errorf("number_words(%q) = %q, quiero %q", in, got, want)
		}
	}
}

func TestParseNumber(t *testing.T) {
	cases := []struct {
		in   string
		want string // RatString; "" si no es un número
	}{
		{"3,5", "7/2"},
		{"3.5", "7/2"},
		{"1/2", "1/2"},
		{"0,5", "1/2"},
		{"1.000.000", "1000000"},
		{"8.848", "8848"},
		{"1 000 000", "1000000"},
		{"1.000,5", ""},
		{"1.0", "1"},
		{"-5", "-5"},
		{"1234567", "1234567"},
		{"1/0", ""},
		{"1.2.3", ""},
		{"12 34", ""},
		{"1,2,3", ""},
	}
	for _, c := range cases {
		got, ok := parseNumber(c.in)
		switch {
		case c.want == "" && ok:
			t.Errorf("parseNumber(%q) = %s, quiero que no sea un número", c.in, got.RatString())
		case c.want != "" && (!ok || got.RatString() != c.want):
			t.Errorf("parseNumber(%q) = %v (%v), quiero %s", c.in, got, ok, c.want)
		}
	}
}

func TestParseQuantityUnits(t *testing.T) {
	c, err := NewAnswerComparator(nil)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		in   string
		val  int64
		unit string
		ok   bool
	}{
		{"5cm", 5, "cm", true},
		{"5 cm", 5, "cm", true},
		{"5 centímetros", 5, "cm", true},
		{"5 centimetros", 5, "cm", true},
		{"cinco centímetros", 5, "cm", true},
		{"120 Km", 120, "km", true},
		{"2 l", 2, "l", true},
		{"100°C", 100, "°c", true},
		{"100 grados Celsius", 100, "°c", true},
		{"100 grados centígrados", 100, "°c", true},
		{"50 %", 50, "%", true},
		{"12", 12, "", true},
		{"100 grados", 0, "", false},
		{"12 meses", 0, "", false},
		{"100 °F", 0, "", false},
		{"ciento uno", 0, "", false},
	}
	for _, tc := range cases {
		form, _ := c.canonical(tc.in, nil)
		q, ok := parseQuantity(form, c.units)
		if ok != tc.ok {
			t.Errorf("parseQuantity(%q→%q) ok=%v, quiero %v", tc.in, form, ok, tc.ok)
			continue
		}
		if ok && (q.unit != tc.unit || q.value.Cmp(big.NewRat(tc.val, 1)) != 0) {
			t.Errorf("parseQuantity(%q) = %s %q, quiero %d %q", tc.in, q.value.RatString(), q.unit, tc.val, tc.unit)
		}
	}
}

// Un número se compara como número y la unidad no es opcional.
func TestQuantityEquality(t *testing.T) {
	c, err := NewAnswerComparator(nil)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		a, b string
		same bool
	}{
		{"5 cm", "5cm", true},
		{"5 cm", "5 centímetros", true},
		{"5 m", "5 cm", false},
		{"1000 g", "1000 kg", false},
		{"12", "12 kg", false},
		{"1234567", "1234568", false},
		{"1/2", "0,5", true},
		{"3,5", "3.5", true},
	}
	for _, tc := range cases {
		fa, _ := c.canonical(tc.a, nil)
		fb, _ := c.canonical(tc.b, nil)
		qa, okA := parseQuantity(fa, c.units)
		qb, okB := parseQuantity(fb, c.units)
		if !okA || !okB {
			t.Errorf("%q o %q no se leen como número", tc.a, tc.b)
			continue
		}
		if qa.equal(qb) != tc.same {
			t.Errorf("%q y %q: iguales=%v, quiero %v", tc.a, tc.b, qa.equal(qb), tc.same)
		}
	}
}

// La tabla de unidades no choca con ninguna combinación de reglas: dos unidades
// distintas nunca comparten forma canónica.
func TestUnitTableHasNoCollisions(t *testing.T) {
	names := RuleNames()
	for mask := 0; mask < 1<<len(names); mask++ {
		overrides := map[string]bool{}
		for i, n := range names {
			overrides[n] = mask&(1<<i) != 0
		}
		c, err := NewAnswerComparator(overrides)
		if err != nil {
			t.Fatal(err)
		}
		seen := map[string]string{}
		for unit, aliases := range unitAliases {
			for _, alias := range aliases {
				for _, a := range []string{alias, Normalize(alias)} {
					form, _ := c.canonical(a, nil)
					if prev, ok := seen[form]; ok && prev != unit {
						t.Errorf("reglas %v: %q es a la vez %q y %q", overrides, form, prev, unit)
					}
					seen[form] = unit
				}
			}
		}
	}
}
