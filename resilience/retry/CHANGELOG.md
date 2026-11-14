# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/resilience/retry` se registran aqui.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencias directas: `github.com/stretchr/testify` `v1.11.1` → `v1.12.1`.
- Lint 2.14.0: `//nolint:errcheck` en un descarte deliberado de test.
- El módulo entra a la banda `v0.900.x` (ADR 0010); el tag anterior `v0.1.x` queda como procedencia.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.102.0] - 2026-04-03

### Added

- `WithRetry(ctx, cfg, operation)` con backoff exponencial configurable.
- `ErrorClassifier` inyectable para distinguir errores permanentes vs transitorios.
- `Logger` opcional (nil-safe) para observabilidad de reintentos.
- `IsContextError()` helper para detectar cancelacion de contexto.
- `DefaultConfig()` con valores sensatos (3 retries, 500ms initial, 10s max).
