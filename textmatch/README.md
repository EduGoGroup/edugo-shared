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

### Reglas de lengua de la respuesta corta (AnswerComparator)

La pregunta corta arma su comparador con su mapa de reglas (`answer_spec.rules`, ADR 0049). Cada regla es una implementación del contrato `Rule`: un nombre estable, un valor por defecto y la forma canónica de **un** texto. Se encadenan en un orden fijo, que es parte del contrato:

| Regla | Qué perdona | Nace |
|---|---|---|
| `punctuation` | Los signos: «Lima.» y «lima». Conserva la coma, el punto y la barra entre cifras, el signo menos, «%» y el grado | encendida |
| `accents` | La tilde: «oxigeno» y «oxígeno» (la «ñ» se conserva) | encendida |
| `number_words` | Cifra y letra del 0 al 100: «once» y «11». «un» y «una» no son números; «ciento uno» no se convierte | encendida |
| `leading_article` | El artículo inicial: «la fotosíntesis» y «fotosíntesis». Si la esperada lo lleva, se exige | encendida |
| `gender_number` | El género y el número: «blanca» y «blanco», «aves» y «ave» | **apagada** |

```go
cmp, err := textmatch.NewAnswerComparator(map[string]bool{"accents": false}) // solo lo que la docente cambió
r, _ := cmp.Compare(ctx, "canté", "cante")
// r.Outcome == textmatch.OutcomeNoMatch: solo difieren en lo que la regla apagada exige
```

- **`OutcomeMatch`** (casa): las formas canónicas con las reglas encendidas son iguales y no vacías.
- **`OutcomeNoMatch`** (no se da por buena): solo difieren en lo que perdonaría una regla apagada; o la esperada es un número (con o sin unidad de la tabla cerrada) y no vale lo mismo; o la respuesta está vacía. **No pasa a la cercanía.**
- **`OutcomeUncertain`** (no opina): decide quien llama, con `Fuzzy` sobre el texto crudo.

Cada texto esperado (la respuesta y cada alternativa) se compara por su lado: una alternativa «12 kg» se compara como número; una «la capital», exige su artículo. **No lo metas en una `Cascade` delante de `Fuzzy`**: la cascada escala también en `NoMatch`. `Normalize`, `Fuzzy`, `Cascade` y `SetMatcher` no cambian.

## Componentes principales

- **Normalize / SplitTokens**: contrato canónico de normalización (minúsculas, sin tildes/diéresis, **preserva la «ñ»**, colapsa espacios) y tokenización (frontera no-alfanumérica unicode, descarta conectores «y»/«e»).
- **EditDistance**: distancia de edición por runas.
- **Strategy / Comparator / Cascade**: contrato de comparación de un par y su orquestador con escalado explícito.
- **Exact / Fuzzy**: estrategias deterministas puras (Fuzzy = el escalón ortográfico que faltaba entre el match exacto y el juicio del LLM).
- **SetMatcher / Policy / GenerateCandidates / MatchReport**: match de conjunto con generación de candidatos y política de estrictez.
- **Rule / TargetGuard / Rules / RuleNames / LookupRule / ResolveRules**: el contrato de regla de lengua y su registro (nombres, valores por defecto, orden). Una regla desconocida da `ErrUnknownRule`.
- **AnswerComparator / NormalizeBase**: el comparador que arma una pregunta corta y la forma base sin quitar tildes de la que parten las reglas.

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
- ADR 0049 (las reglas de la respuesta corta) · ficha HU-INF.4 del MVP-INF, en el repo raíz del ecosistema.
