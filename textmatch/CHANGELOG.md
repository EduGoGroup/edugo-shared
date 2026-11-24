# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/textmatch` se registran aquí.

## [Unreleased]

## [0.900.2] - 2026-11-23

### Added
- **Reglas de lengua de la respuesta corta** (ADR 0049, HU-INF.4 del MVP-INF): el contrato `Rule` (nombre estable, valor por defecto, forma canónica de un texto) y el opcional `TargetGuard`; el registro con sus cinco reglas en orden fijo (`punctuation`, `accents`, `number_words`, `leading_article` y `gender_number`, esta apagada por defecto), con `Rules`, `RuleNames`, `LookupRule`, `ResolveRules` y el error tipado `ErrUnknownRule`.
- `AnswerComparator` (`NewAnswerComparator`): el comparador que arma una pregunta con su mapa de reglas. Implementa `Strategy` y da `OutcomeMatch` (casa), `OutcomeNoMatch` (no se da por buena: una regla apagada lo exige, o un número que no vale lo mismo) u `OutcomeUncertain` (no opina: decide `Fuzzy`). Un número, con o sin unidad de una tabla cerrada, se compara con aritmética exacta y sin cercanía.
- `NormalizeBase`: minúsculas, NFC y espacios colapsados, sin quitar tildes. Con la regla `accents` da exactamente `Normalize`.

### Design Notes
- `Normalize`, `Fuzzy`, `Cascade`, `SetMatcher`, `EditDistance` y el umbral 0,85 no cambian: el worker no nota nada.
- `AnswerComparator` no va dentro de una `Cascade` delante de `Fuzzy`: la cascada escala también en `OutcomeNoMatch`.

## [0.900.1] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencias directas: `golang.org/x/text` `v0.38.0` → `v0.44.0`.

## [0.900.0] - 2026-07-17

### Added

- Módulo nuevo `textmatch`: motor de comparación de respuestas en cascada, puro y reutilizable (ADR 0035, plan 045 F0/F1). Sin consumidores todavía; queda como fundación para learning, worker y el pipeline material→evaluación (planes 043/044).
- **Nivel 1 (comparación de un par):** `Outcome` (`OutcomeNoMatch`/`OutcomeMatch`/`OutcomeUncertain`), `Result`, `Strategy` y `Comparator`; estrategias deterministas `Exact` y `Fuzzy{Threshold}` (con `NewFuzzy`, default `DefaultFuzzyThreshold = 0.85`); orquestador `Cascade` (`NewCascade`) con escalado explícito (positivo corta, incierto/negativo escala, error se propaga). La estrategia LLM NO vive aquí: se inyecta desde su infra (DIP).
- **Nivel 2 (match de conjunto):** `SetMatcher` (`NewSetMatcher`) con `Policy` (`PolicyStrict`/`PolicyLenient`), `MatchReport`, `Candidate` y `GenerateCandidates` (tokens + n-gramas contiguos). Dos entradas: `Match` (candidatos atómicos) y `MatchAnswer` (helper de alto nivel que deriva tokens y arma candidatos; para Strict, token base sin consumir = foráneo).
- **Normalización canónica:** `Normalize` (minúsculas, sin tildes/diéresis, **preserva la «ñ»**, colapsa espacios) y `SplitTokens` (frontera no-alfanumérica unicode, descarta conectores «y»/«e»); `EditDistance` por runas.
- Suite de tests unitarios (tabla-driven) con los casos reales del research, verde con race detector.

### Design Notes

- **Distancia Damerau-Levenshtein (OSA)** en vez de Levenshtein puro: la transposición de runas adyacentes cuenta como una sola edición, lo que rescata `whastapp`≈`whatsapp` (similitud 0.875 ≥ 0.85) sin bajar el umbral conservador de 0.85.
- La «ñ» se preserva como letra propia (`año` ≠ `ano`), unificando la divergencia histórica entre learning (la borraba) y worker (la preservaba a mano).
- Dependencia externa única: `golang.org/x/text`. Sin acoplamiento a LLM/red.
- Módulo de nivel 0: no depende de otros módulos de `edugo-shared`.
