# Textmatch

Motor de comparación de respuestas en **cascada**, puro y reutilizable (ADR 0035, plan 045). Consolida en un solo lugar —neutral, sin infraestructura ni LLM— las utilidades de similitud y normalización que antes vivían duplicadas y divergentes entre `edugo-api-learning` y `edugo-worker`. Lo consumen (o consumirán) learning, worker y el pipeline material→evaluación (planes 043/044) sin reinventar la comparación.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/textmatch
```

Dependencia externa única: `golang.org/x/text` (normalización Unicode).

## Quick Start

### Nivel 1 — comparar UN par (Comparator / Cascade)

Una cascada ordena estrategias baratas→caras. Positivo corta; incierto/negativo escala a la siguiente; un `error` (transitorio, típico del LLM) se propaga.

```go
cmp := textmatch.NewCascade(textmatch.Exact{}, textmatch.NewFuzzy(0.85))

r, err := cmp.Compare(ctx, "whatsapp", "whastapp")
// r.Outcome == textmatch.OutcomeMatch, r.Confidence ~= 0.875, r.Strategy == "fuzzy"
```

La estrategia LLM **no vive aquí** (depende de infra/red): se define donde vive el provider y se inyecta como una `Strategy` más en la cascada (DIP).

### Nivel 2 — matchear un CONJUNTO (SetMatcher + Policy)

`SetMatcher` matchea el conjunto de ítems del alumno contra el conjunto esperado, con una **política de completitud** que es decisión de negocio, ortogonal a cómo se compara un par.

```go
m := textmatch.NewSetMatcher(cmp, textmatch.PolicyLenient)

rep, err := m.MatchAnswer(ctx,
    []string{"facebook", "instagram", "whatsapp"},
    "whastapp instalgram y el famoso facebook")
// rep.Complete == true: los typos se rescatan sin LLM; el relleno "el famoso" se ignora.
```

- `MatchAnswer(expected, studentAnswer)` — helper de alto nivel: deriva tokens y arma candidatos (tokens sueltos + n-gramas contiguos hasta el esperado más largo, así los ítems multi-palabra como "costa rica" casan). Para `PolicyStrict`, un **token base sin consumir** es foráneo e invalida el match.
- `Match(expected, candidates)` — nivel bajo para callers que ya tienen sus unidades discretas; un candidato sobrante es foráneo.

### Políticas

- `PolicyStrict` (aula/learning): todos los esperados cubiertos **y** ningún token foráneo (un ítem extra no reconocido invalida).
- `PolicyLenient` (worker triturado): todos los esperados cubiertos; los sobrantes se ignoran.

## Componentes principales

- **Normalize / SplitTokens**: contrato canónico de normalización (minúsculas, sin tildes/diéresis, **preserva la «ñ»**, colapsa espacios) y tokenización (frontera no-alfanumérica unicode, descarta conectores «y»/«e»).
- **EditDistance**: distancia de edición por runas.
- **Strategy / Comparator / Cascade**: contrato de comparación de un par y su orquestador con escalado explícito.
- **Exact / Fuzzy**: estrategias deterministas puras (Fuzzy = el escalón ortográfico que faltaba entre el match exacto y el juicio del LLM).
- **SetMatcher / Policy / GenerateCandidates / MatchReport**: match de conjunto con generación de candidatos y política de estrictez.

## Operación local

```bash
make build      # Verificar que el módulo compila
make test       # Ejecutar tests
make test-race  # Tests con race detector
make check      # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- **Distancia = Damerau-Levenshtein (OSA), no Levenshtein puro.** Cuenta la transposición de dos runas adyacentes como **una** edición. Es lo que rescata el caso canónico `whastapp`≈`whatsapp` (un intercambio s↔t): con Levenshtein puro esa distancia sería 2 (similitud 0.75 < 0.85 → no casaría), contradiciendo el objetivo del plan; con transposición es 1 (similitud 0.875 ≥ 0.85) **manteniendo el umbral 0.85 conservador**.
- **La «ñ» se preserva** (es letra, no tilde): «año» ≠ «ano». Unifica la divergencia histórica (learning borraba la ñ, worker la preservaba a mano).
- Dependencia externa única (`golang.org/x/text`); sin acoplamiento a LLM/red (la estrategia LLM se inyecta desde su infra).
- El motor **no decide negocio**: qué significa un no-match (ítem ausente → red del profesor) lo interpreta el caller.

## Documentación

- [Changelog](CHANGELOG.md)
- ADR 0035 · plan 045 (`docs/plans/045-motor-comparacion-cascada/`).
