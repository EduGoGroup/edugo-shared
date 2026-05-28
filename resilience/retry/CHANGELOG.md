# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/resilience/retry` se registran aqui.

## [Unreleased]

## [0.102.0] - 2026-04-03

### Added

- `WithRetry(ctx, cfg, operation)` con backoff exponencial configurable.
- `ErrorClassifier` inyectable para distinguir errores permanentes vs transitorios.
- `Logger` opcional (nil-safe) para observabilidad de reintentos.
- `IsContextError()` helper para detectar cancelacion de contexto.
- `DefaultConfig()` con valores sensatos (3 retries, 500ms initial, 10s max).
