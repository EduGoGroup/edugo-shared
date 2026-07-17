package textmatch

import (
	"context"
	"errors"
	"testing"
)

func TestExactCompare(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name           string
		expected, cand string
		wantOutcome    Outcome
	}{
		{"igual tras normalizar", "Facebook", "facebook", OutcomeMatch},
		{"tildes igualan", "Perú", "peru", OutcomeMatch},
		{"ñ distingue", "año", "ano", OutcomeNoMatch},
		{"typo no casa exacto", "whatsapp", "whastapp", OutcomeNoMatch},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Exact{}.Compare(ctx, tt.expected, tt.cand)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if r.Outcome != tt.wantOutcome {
				t.Errorf("Exact(%q,%q) outcome = %v, quiero %v", tt.expected, tt.cand, r.Outcome, tt.wantOutcome)
			}
		})
	}
}

func TestFuzzyCompare(t *testing.T) {
	ctx := context.Background()
	f := NewFuzzy(0.85)
	tests := []struct {
		name           string
		expected, cand string
		wantOutcome    Outcome
	}{
		{"typo transpuesto whatsapp", "whatsapp", "whastapp", OutcomeMatch},   // sim 0.875
		{"typo inserción instagram", "instagram", "instalgram", OutcomeMatch}, // sim 0.9
		{"dos ítems distintos", "facebook", "instagram", OutcomeNoMatch},
		{"ñ no se cuela bajo umbral", "año", "ano", OutcomeNoMatch}, // sim 0.667
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := f.Compare(ctx, tt.expected, tt.cand)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if r.Outcome != tt.wantOutcome {
				t.Errorf("Fuzzy(%q,%q) outcome = %v (conf %.3f), quiero %v", tt.expected, tt.cand, r.Outcome, r.Confidence, tt.wantOutcome)
			}
		})
	}
}

// TestNewFuzzyDefault confirma que threshold <= 0 cae al default 0.85 (rescata el
// typo transpuesto).
func TestNewFuzzyDefault(t *testing.T) {
	ctx := context.Background()
	for _, f := range []Fuzzy{NewFuzzy(0), NewFuzzy(-1), {}} {
		r, err := f.Compare(ctx, "whatsapp", "whastapp")
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if r.Outcome != OutcomeMatch {
			t.Errorf("Fuzzy con threshold no positivo debería usar 0.85 y casar el typo; outcome = %v", r.Outcome)
		}
	}
}

// errStrategy es un stub que registra si fue invocada y devuelve un error, para
// probar el escalado y la propagación de error de la Cascade.
type errStrategy struct{ called *bool }

func (errStrategy) Name() string { return "err-stub" }
func (s errStrategy) Compare(_ context.Context, _, _ string) (Result, error) {
	*s.called = true
	return Result{}, errors.New("fallo transitorio")
}

func TestCascade(t *testing.T) {
	ctx := context.Background()

	t.Run("positivo corta sin llamar a la siguiente", func(t *testing.T) {
		called := false
		c := NewCascade(Exact{}, errStrategy{called: &called})
		r, err := c.Compare(ctx, "facebook", "facebook")
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if r.Outcome != OutcomeMatch || r.Strategy != "exact" {
			t.Errorf("esperaba Match de exact, obtuve %+v", r)
		}
		if called {
			t.Error("la estrategia siguiente no debía invocarse tras un positivo")
		}
	})

	t.Run("negativo escala a la siguiente", func(t *testing.T) {
		c := NewCascade(Exact{}, NewFuzzy(0.85))
		r, err := c.Compare(ctx, "whatsapp", "whastapp")
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if r.Outcome != OutcomeMatch || r.Strategy != "fuzzy" {
			t.Errorf("esperaba Match de fuzzy tras escalar, obtuve %+v", r)
		}
	})

	t.Run("error se propaga", func(t *testing.T) {
		called := false
		c := NewCascade(errStrategy{called: &called})
		_, err := c.Compare(ctx, "a", "b")
		if err == nil {
			t.Fatal("esperaba que el error se propagara")
		}
	})

	t.Run("agotada devuelve el último negativo", func(t *testing.T) {
		c := NewCascade(Exact{}, NewFuzzy(0.99))
		r, err := c.Compare(ctx, "facebook", "instagram")
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if r.Outcome == OutcomeMatch {
			t.Errorf("esperaba no-match final, obtuve %+v", r)
		}
	})
}
