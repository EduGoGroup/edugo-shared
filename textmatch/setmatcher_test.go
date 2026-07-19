package textmatch

import (
	"context"
	"reflect"
	"testing"
)

// cascade857 es el comparador determinista de los dos carriles: exacto + fuzzy 0.85.
func cascade085() *Cascade { return NewCascade(Exact{}, NewFuzzy(0.85)) }

func TestMatchAnswerCasoRealResearch(t *testing.T) {
	ctx := context.Background()
	expected := []string{"facebook", "instagram", "whatsapp"}
	student := "whastapp instalgram y el famoso facebook"

	t.Run("Lenient completa sin LLM", func(t *testing.T) {
		m := NewSetMatcher(cascade085(), PolicyLenient)
		rep, err := m.MatchAnswer(ctx, expected, student)
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		for i, c := range rep.Covered {
			if !c {
				t.Errorf("esperado %q no cubierto", expected[i])
			}
		}
		if !rep.Complete {
			t.Errorf("Lenient debía completar (typos rescatados, relleno ignorado); rep=%+v", rep)
		}
	})

	t.Run("Strict no completa por tokens foráneos", func(t *testing.T) {
		m := NewSetMatcher(cascade085(), PolicyStrict)
		rep, err := m.MatchAnswer(ctx, expected, student)
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		for i, c := range rep.Covered {
			if !c {
				t.Errorf("esperado %q no cubierto", expected[i])
			}
		}
		if rep.Complete {
			t.Error("Strict NO debía completar: quedan tokens foráneos 'el','famoso'")
		}
		if len(rep.Leftover) != 2 {
			t.Errorf("esperaba 2 tokens foráneos ('el','famoso'), Leftover=%v", rep.Leftover)
		}
	})
}

func TestMatchAnswerMultiPalabra(t *testing.T) {
	ctx := context.Background()
	expected := []string{"costa rica", "panama"}
	m := NewSetMatcher(cascade085(), PolicyStrict)
	rep, err := m.MatchAnswer(ctx, expected, "costa rica y panama")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !rep.Complete {
		t.Errorf("ítem multi-palabra 'costa rica' debía casar vía n-grama; rep=%+v", rep)
	}
}

func TestMatchAnswerItemForaneoStrict(t *testing.T) {
	ctx := context.Background()
	expected := []string{"ecuador", "venezuela", "colombia"}
	student := "ecuador venezuela colombia y peru"

	strict := NewSetMatcher(cascade085(), PolicyStrict)
	repS, err := strict.MatchAnswer(ctx, expected, student)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repS.Complete {
		t.Error("Strict NO debía completar: 'peru' es un ítem extra foráneo")
	}

	lenient := NewSetMatcher(cascade085(), PolicyLenient)
	repL, err := lenient.MatchAnswer(ctx, expected, student)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !repL.Complete {
		t.Error("Lenient SÍ debía completar: el sobrante 'peru' se ignora")
	}
}

func TestMatchAnswerFaltanteReal(t *testing.T) {
	ctx := context.Background()
	expected := []string{"facebook", "instagram", "whatsapp"}
	// Falta whatsapp de verdad: ni exacto ni fuzzy lo rescatan.
	for _, policy := range []Policy{PolicyStrict, PolicyLenient} {
		m := NewSetMatcher(cascade085(), policy)
		rep, err := m.MatchAnswer(ctx, expected, "facebook instagram")
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if rep.Complete {
			t.Errorf("policy %v: no debía completar con un esperado ausente; rep=%+v", policy, rep)
		}
	}
}

// TestMatchAnswerForaneosVsDuplicados cubre la semántica de foráneos de Strict:
// un sobrante es foráneo solo si no corresponde a ningún esperado (regla 042: los
// duplicados/variantes no penalizan; un ítem extra real sí).
func TestMatchAnswerForaneosVsDuplicados(t *testing.T) {
	ctx := context.Background()
	expected := []string{"ecuador", "venezuela", "colombia"}
	tests := []struct {
		name         string
		expected     []string
		student      string
		wantComplete bool
	}{
		{
			name:         "Caso 1: relleno foráneo (el/famoso) invalida",
			expected:     []string{"facebook", "instagram", "whatsapp"},
			student:      "whastapp instalgram y el famoso facebook",
			wantComplete: false,
		},
		{
			name:         "ítem extra real (peru) invalida",
			expected:     expected,
			student:      "ecuador venezuela colombia y peru",
			wantComplete: false,
		},
		{
			name:         "duplicado exacto no penaliza (restaura 042)",
			expected:     expected,
			student:      "ecuador ecuador venezuela colombia",
			wantComplete: true,
		},
		{
			name:         "duplicado con typo casa por fuzzy y no penaliza",
			expected:     expected,
			student:      "ecuador ecuadr venezuela colombia",
			wantComplete: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewSetMatcher(cascade085(), PolicyStrict)
			rep, err := m.MatchAnswer(ctx, tt.expected, tt.student)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			for i, c := range rep.Covered {
				if !c {
					t.Errorf("esperado %q no cubierto", tt.expected[i])
				}
			}
			if rep.Complete != tt.wantComplete {
				t.Errorf("Strict Complete = %v, quiero %v (rep=%+v)", rep.Complete, tt.wantComplete, rep)
			}
		})
	}
}

// TestMatchAnswerDuplicadoLenient confirma que Lenient no toca foráneos: los
// duplicados y sobrantes se ignoran igual.
func TestMatchAnswerDuplicadoLenient(t *testing.T) {
	ctx := context.Background()
	expected := []string{"ecuador", "venezuela", "colombia"}
	m := NewSetMatcher(cascade085(), PolicyLenient)
	rep, err := m.MatchAnswer(ctx, expected, "ecuador ecuador venezuela colombia y peru")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !rep.Complete {
		t.Errorf("Lenient debía completar (duplicado y sobrante ignorados); rep=%+v", rep)
	}
}

// TestMatchAtómico cubre el nivel bajo (candidatos discretos, foráneo = sobrante).
func TestMatchAtómico(t *testing.T) {
	ctx := context.Background()
	expected := []string{"facebook", "instagram"}
	candidates := []string{"facebook", "instalgram", "extra"}

	strict := NewSetMatcher(cascade085(), PolicyStrict)
	repS, err := strict.Match(ctx, expected, candidates)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repS.Complete {
		t.Error("Strict NO debía completar: 'extra' es candidato foráneo")
	}
	if !reflect.DeepEqual(repS.Leftover, []int{2}) {
		t.Errorf("esperaba Leftover=[2] (índice de 'extra'), obtuve %v", repS.Leftover)
	}

	lenient := NewSetMatcher(cascade085(), PolicyLenient)
	repL, err := lenient.Match(ctx, expected, candidates)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !repL.Complete {
		t.Error("Lenient SÍ debía completar: el candidato 'extra' se ignora")
	}
}

// TestMatchAtómicoDuplicado confirma que en el nivel atómico un candidato sobrante
// que matchea un esperado (duplicado) no penaliza en Strict.
func TestMatchAtómicoDuplicado(t *testing.T) {
	ctx := context.Background()
	expected := []string{"facebook", "instagram"}
	candidates := []string{"facebook", "instagram", "facebook"} // el tercero es duplicado

	strict := NewSetMatcher(cascade085(), PolicyStrict)
	rep, err := strict.Match(ctx, expected, candidates)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !rep.Complete {
		t.Errorf("Strict debía completar: el sobrante 'facebook' es duplicado de un esperado; rep=%+v", rep)
	}
}

func TestGenerateCandidates(t *testing.T) {
	got := GenerateCandidates([]string{"costa", "rica", "panama"}, 2)
	want := []Candidate{
		{Text: "costa", Start: 0, End: 1},
		{Text: "costa rica", Start: 0, End: 2},
		{Text: "rica", Start: 1, End: 2},
		{Text: "rica panama", Start: 1, End: 3},
		{Text: "panama", Start: 2, End: 3},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GenerateCandidates = %+v, quiero %+v", got, want)
	}
}
