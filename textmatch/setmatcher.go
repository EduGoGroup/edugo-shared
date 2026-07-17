package textmatch

import (
	"context"
	"strings"
)

// Policy captura la ESTRICTEZ de un match de conjunto: una decisión de negocio,
// ortogonal a cómo se compara un par (D-045.6).
type Policy int

const (
	// PolicyStrict (aula/learning): todos los esperados cubiertos Y ningún
	// sobrante FORÁNEO. Un sobrante es foráneo solo si NO corresponde a ningún
	// esperado vía el comparador (exact/fuzzy); así un ítem extra no reconocido
	// ("perú") invalida el match, pero un DUPLICADO o VARIANTE de algo esperado
	// ("ecuador" repetido, o "ecuadr" con typo) NO penaliza. Generaliza con fuzzy
	// la membresía de tokens de MatchesShortAnswerListPrep (regla firmada 042: los
	// duplicados de un ítem esperado no penalizan).
	PolicyStrict Policy = iota
	// PolicyLenient (worker triturado): todos los esperados cubiertos; los
	// sobrantes se ignoran (relleno de cortesía como "el famoso" no penaliza).
	// Es el Grade de hoy.
	PolicyLenient
)

// Candidate es un candidato de match con su rango sobre los tokens base del
// alumno. El rango [Start, End) permite a MatchAnswer saber qué tokens consume un
// n-grama, para el chequeo de foráneos de PolicyStrict.
type Candidate struct {
	Text  string
	Start int // índice de token base inicial (inclusive)
	End   int // índice de token base final (exclusivo)
}

// GenerateCandidates produce, a partir de los tokens base del alumno, los tokens
// sueltos más los n-gramas contiguos hasta maxLen tokens (D-045.6). Esto arregla
// el split pegado del worker: "whastapp instalgram y el famoso facebook" ofrece
// tokens sueltos y cada ítem encuentra su casi-match, y los ítems multi-palabra
// ("costa rica") encuentran su n-grama. Orden determinista: por inicio ascendente
// y, dentro de cada inicio, por longitud ascendente (el primer Match gana, así el
// n-grama del largo justo se prefiere sin desperdiciar tokens). Exportada para
// que 043/044 la reusen.
func GenerateCandidates(tokens []string, maxLen int) []Candidate {
	if maxLen < 1 {
		maxLen = 1
	}
	out := make([]Candidate, 0, len(tokens)*maxLen)
	for start := 0; start < len(tokens); start++ {
		for l := 1; l <= maxLen && start+l <= len(tokens); l++ {
			out = append(out, Candidate{
				Text:  strings.Join(tokens[start:start+l], " "),
				Start: start,
				End:   start + l,
			})
		}
	}
	return out
}

// MatchReport es el resultado de un match de conjunto (D-045.6).
type MatchReport struct {
	Covered  []bool // por ítem esperado: ¿cubierto?
	UsedBy   []int  // por ítem esperado: índice que lo cubrió, o -1 (ver cada método)
	Leftover []int  // índices no consumidos (candidatos o tokens; ver cada método)
	Complete bool   // según la Policy
}

// SetMatcher matchea un conjunto de candidatos contra un conjunto de esperados,
// usando un Comparator (Nivel 1) por celda y una Policy de completitud (D-045.6).
type SetMatcher struct {
	cmp    Comparator
	policy Policy
}

// NewSetMatcher construye un SetMatcher con el comparador (típicamente una Cascade)
// y la política de estrictez.
func NewSetMatcher(cmp Comparator, policy Policy) *SetMatcher {
	return &SetMatcher{cmp: cmp, policy: policy}
}

// Match asigna cada esperado al mejor candidato NO usado que dé OutcomeMatch vía el
// Comparator (greedy, marca usado, empate = primer candidato para determinismo).
// Trata cada candidato como una unidad ATÓMICA: no genera n-gramas. Es el nivel
// bajo para callers que ya tienen sus unidades discretas (p. ej. los fragmentos
// del worker). En el MatchReport: UsedBy[e] = índice de candidato o -1; Leftover =
// índices de candidatos no consumidos. Para PolicyStrict, un candidato sobrante es
// foráneo solo si NO matchea ningún ESPERADO vía el comparador (un duplicado de un
// esperado ya cubierto no penaliza); un foráneo real invalida el match. A
// diferencia de MatchAnswer, la referencia del chequeo de foráneos son los ítems
// esperados completos (no sus tokens), porque aquí el candidato es la unidad. Un
// error del Comparator se propaga.
func (m *SetMatcher) Match(ctx context.Context, expected, candidates []string) (MatchReport, error) {
	used := make([]bool, len(candidates))
	rep := newReport(len(expected))
	for e, exp := range expected {
		for j, cand := range candidates {
			if used[j] {
				continue
			}
			r, err := m.cmp.Compare(ctx, exp, cand)
			if err != nil {
				return MatchReport{}, err
			}
			if r.Outcome == OutcomeMatch {
				rep.Covered[e] = true
				rep.UsedBy[e] = j
				used[j] = true
				break
			}
		}
	}
	rep.Leftover = leftoverIndices(used)
	complete, err := m.completeUnderPolicy(ctx, rep.Covered, textsAt(candidates, rep.Leftover), expected)
	if err != nil {
		return MatchReport{}, err
	}
	rep.Complete = complete
	return rep, nil
}

// MatchAnswer es el helper de alto nivel: toma la respuesta CRUDA del alumno,
// deriva sus tokens base, arma los candidatos (tokens + n-gramas contiguos hasta
// la longitud del esperado más largo) y hace el match greedy. Es la forma que
// reusan learning y worker sin duplicar el armado de candidatos. A diferencia de
// Match, los índices del MatchReport son de TOKEN base: UsedBy[e] = índice del
// token inicial del span que cubrió al esperado (o -1); Leftover = índices de
// tokens base no consumidos. Para PolicyStrict, un token sobrante es foráneo solo
// si NO matchea ningún TOKEN de ningún esperado vía el comparador (reproduce la
// membresía de tokens de MatchesShortAnswerListPrep generalizada con fuzzy:
// duplicados/variantes con typo de algo esperado no penalizan; un token que no
// corresponde a nada —"el", "peru"— sí invalida). Un error del Comparator se propaga.
func (m *SetMatcher) MatchAnswer(ctx context.Context, expected []string, studentAnswer string) (MatchReport, error) {
	tokens := SplitTokens(studentAnswer)
	maxLen := 1
	for _, e := range expected {
		if n := len(SplitTokens(e)); n > maxLen {
			maxLen = n
		}
	}
	candidates := GenerateCandidates(tokens, maxLen)
	used := make([]bool, len(tokens))
	rep := newReport(len(expected))
	for e, exp := range expected {
		for _, c := range candidates {
			if spanUsed(used, c) {
				continue
			}
			r, err := m.cmp.Compare(ctx, exp, c.Text)
			if err != nil {
				return MatchReport{}, err
			}
			if r.Outcome == OutcomeMatch {
				rep.Covered[e] = true
				rep.UsedBy[e] = c.Start
				for k := c.Start; k < c.End; k++ {
					used[k] = true
				}
				break
			}
		}
	}
	rep.Leftover = leftoverIndices(used)
	// El chequeo de foráneos de Strict compara cada token sobrante contra los
	// TOKENS de los esperados (no los ítems completos): así "costa" sobrante
	// corresponde al token de "costa rica", y un duplicado casa su propio token.
	complete, err := m.completeUnderPolicy(ctx, rep.Covered, textsAt(tokens, rep.Leftover), expectedTokenSet(expected))
	if err != nil {
		return MatchReport{}, err
	}
	rep.Complete = complete
	return rep, nil
}

// newReport inicializa un MatchReport con UsedBy en -1 (nadie cubre aún).
func newReport(n int) MatchReport {
	rep := MatchReport{Covered: make([]bool, n), UsedBy: make([]int, n)}
	for e := range rep.UsedBy {
		rep.UsedBy[e] = -1
	}
	return rep
}

// spanUsed indica si algún token del rango del candidato ya fue consumido.
func spanUsed(used []bool, c Candidate) bool {
	for k := c.Start; k < c.End; k++ {
		if used[k] {
			return true
		}
	}
	return false
}

// leftoverIndices devuelve los índices marcados como no usados.
func leftoverIndices(used []bool) []int {
	var out []int
	for i, u := range used {
		if !u {
			out = append(out, i)
		}
	}
	return out
}

// completeUnderPolicy aplica la Policy. Lenient: basta con todos los esperados
// cubiertos (sobrantes ignorados). Strict: además ningún sobrante puede ser
// FORÁNEO; un sobrante (texto en leftoverTexts) es foráneo si no matchea NINGUNA
// referencia vía el comparador. Así los duplicados/variantes de algo esperado no
// penalizan y solo un ítem extra real invalida el match. references son los ítems
// esperados (en Match) o sus tokens (en MatchAnswer). Un error del comparador se
// propaga.
func (m *SetMatcher) completeUnderPolicy(ctx context.Context, covered []bool, leftoverTexts, references []string) (bool, error) {
	if !allCovered(covered) {
		return false, nil
	}
	if m.policy != PolicyStrict {
		return true, nil
	}
	for _, text := range leftoverTexts {
		matched, err := m.matchesAny(ctx, text, references)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil // sobrante que no corresponde a ningún esperado: foráneo
		}
	}
	return true, nil
}

// matchesAny reporta si target da OutcomeMatch contra alguna referencia vía el
// comparador (orden de argumentos: la referencia es el "expected", target el
// "candidate"). Corta en el primer match; un error se propaga.
func (m *SetMatcher) matchesAny(ctx context.Context, target string, references []string) (bool, error) {
	for _, ref := range references {
		r, err := m.cmp.Compare(ctx, ref, target)
		if err != nil {
			return false, err
		}
		if r.Outcome == OutcomeMatch {
			return true, nil
		}
	}
	return false, nil
}

// allCovered reporta si todos los esperados quedaron cubiertos.
func allCovered(covered []bool) bool {
	for _, c := range covered {
		if !c {
			return false
		}
	}
	return true
}

// textsAt proyecta los textos de all en los índices dados.
func textsAt(all []string, idx []int) []string {
	out := make([]string, len(idx))
	for i, j := range idx {
		out[i] = all[j]
	}
	return out
}

// expectedTokenSet es la unión (sin duplicados, orden de aparición) de los tokens
// de todos los ítems esperados. Es la referencia del chequeo de foráneos por token
// de MatchAnswer; alimenta la membresía generalizada con fuzzy.
func expectedTokenSet(expected []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, e := range expected {
		for _, t := range SplitTokens(e) {
			if _, ok := seen[t]; ok {
				continue
			}
			seen[t] = struct{}{}
			out = append(out, t)
		}
	}
	return out
}
