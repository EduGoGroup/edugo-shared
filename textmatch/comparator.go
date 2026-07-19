package textmatch

import (
	"context"
	"fmt"
)

// Outcome es el veredicto de una estrategia sobre UN par (expected, candidate).
type Outcome int

// Los valores de Outcome gobiernan el escalado de la cascada (D-045.2).
const (
	OutcomeNoMatch   Outcome = iota // negativo → escala a la siguiente estrategia
	OutcomeMatch                    // positivo → corta la cascada
	OutcomeUncertain                // incierto → escala a la siguiente estrategia
)

// DefaultFuzzyThreshold es el umbral de similitud por defecto del fuzzy (0.85),
// reusa el del fuzzyMatch histórico de learning. Conservador a propósito
// (D-045.4): estricto con lo incorrecto, tolerante con typos de lo correcto.
const DefaultFuzzyThreshold = 0.85

// Result es el resultado de una comparación de un par.
type Result struct {
	Outcome    Outcome
	Confidence float64 // 0..1
	Evidence   string  // por qué (legible: para el profesor y para depurar)
	Strategy   string  // qué estrategia lo produjo (procedencia)
}

// Strategy es una estrategia de comparación de un par (D-045.2). Contrato mínimo
// y estable (ISP: dos métodos). El ctx + error acomodan estrategias cancelables y
// transitorias (p. ej. LLM); las deterministas nunca devuelven error.
type Strategy interface {
	Name() string
	Compare(ctx context.Context, expected, candidate string) (Result, error)
}

// Comparator es lo que consume el SetMatcher para cada celda. Cascade lo
// implementa; también lo implementa una Strategy suelta (misma firma de Compare).
type Comparator interface {
	Compare(ctx context.Context, expected, candidate string) (Result, error)
}

// Exact es la estrategia de igualdad de strings normalizados (D-045.4).
type Exact struct{}

// Name identifica la estrategia (procedencia en Result.Strategy).
func (Exact) Name() string { return "exact" }

// Compare devuelve Match/1.0 si expected y candidate normalizan al mismo texto.
func (Exact) Compare(_ context.Context, expected, candidate string) (Result, error) {
	if Normalize(expected) == Normalize(candidate) {
		return Result{Outcome: OutcomeMatch, Confidence: 1.0, Evidence: "iguales tras normalizar", Strategy: "exact"}, nil
	}
	return Result{Outcome: OutcomeNoMatch, Confidence: 0.0, Evidence: "distintos tras normalizar", Strategy: "exact"}, nil
}

// Fuzzy es la estrategia ortográfica: distancia de edición normalizada por runas
// sobre los textos normalizados; sim = 1 - dist/maxLen (D-045.4). Es el escalón
// del medio que faltaba entre el match exacto y el juicio del LLM.
type Fuzzy struct {
	// Threshold es la similitud mínima para Match (0..1). Si es <= 0 se usa
	// DefaultFuzzyThreshold, así el literal Fuzzy{} también es seguro.
	Threshold float64
}

// NewFuzzy construye un Fuzzy; threshold <= 0 cae al default 0.85.
func NewFuzzy(threshold float64) Fuzzy {
	if threshold <= 0 {
		threshold = DefaultFuzzyThreshold
	}
	return Fuzzy{Threshold: threshold}
}

// Name identifica la estrategia (procedencia en Result.Strategy).
func (Fuzzy) Name() string { return "fuzzy" }

// Compare produce Match si la similitud alcanza el umbral (Confidence = sim), o
// NoMatch en caso contrario (Confidence = sim igualmente, para depurar). No emite
// OutcomeUncertain: la banda incierta queda reservada para estrategias futuras.
func (f Fuzzy) Compare(_ context.Context, expected, candidate string) (Result, error) {
	threshold := f.Threshold
	if threshold <= 0 {
		threshold = DefaultFuzzyThreshold
	}
	e, c := Normalize(expected), Normalize(candidate)
	if e == c {
		return Result{Outcome: OutcomeMatch, Confidence: 1.0, Evidence: "iguales tras normalizar", Strategy: "fuzzy"}, nil
	}
	maxLen := len([]rune(e))
	if n := len([]rune(c)); n > maxLen {
		maxLen = n
	}
	if maxLen == 0 {
		return Result{Outcome: OutcomeMatch, Confidence: 1.0, Evidence: "ambos vacíos", Strategy: "fuzzy"}, nil
	}
	sim := 1.0 - float64(EditDistance(e, c))/float64(maxLen)
	outcome := OutcomeNoMatch
	if sim >= threshold {
		outcome = OutcomeMatch
	}
	return Result{
		Outcome:    outcome,
		Confidence: sim,
		Evidence:   fmt.Sprintf("similitud %.3f (umbral %.3f)", sim, threshold),
		Strategy:   "fuzzy",
	}, nil
}

// Cascade orquesta una lista ordenada de estrategias (barata→cara) con escalado
// explícito (D-045.3): positivo corta; incierto/negativo escala a la siguiente;
// un error se propaga (transitorio del LLM → el caller reintenta el intento). Si
// se agota sin positivo devuelve el último Result (negativo/incierto): qué
// significa ese no-match lo decide el caller (red del profesor), no el motor.
// Implementa Comparator.
type Cascade struct {
	strategies []Strategy
}

// NewCascade construye una Cascade con las estrategias en orden barata→cara.
func NewCascade(strategies ...Strategy) *Cascade {
	return &Cascade{strategies: strategies}
}

// Compare recorre las estrategias hasta el primer Match; ver Cascade.
func (c *Cascade) Compare(ctx context.Context, expected, candidate string) (Result, error) {
	var last Result
	for _, s := range c.strategies {
		r, err := s.Compare(ctx, expected, candidate)
		if err != nil {
			return Result{}, err
		}
		if r.Outcome == OutcomeMatch {
			return r, nil
		}
		last = r
	}
	return last, nil
}
