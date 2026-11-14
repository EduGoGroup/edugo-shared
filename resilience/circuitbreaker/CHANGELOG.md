# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/resilience/circuitbreaker` se registran aqui.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencias directas: `github.com/stretchr/testify` `v1.11.1` → `v1.12.1`.
- Lint 2.14.0: comentario del bloque de estados y `//nolint:errcheck` en descartes deliberados de tests.
- El módulo entra a la banda `v0.900.x` (ADR 0010); el tag anterior `v0.1.x` queda como procedencia.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.102.0] - 2026-04-03

### Added

- Circuit breaker con 3 estados: Closed, Open, HalfOpen.
- Configuracion flexible: MaxFailures, Timeout, MaxRequests, SuccessThreshold.
- `MetricsHook` interface para inyeccion opcional de metricas.
- `Execute(ctx, fn)` para ejecutar operaciones protegidas.
- `DefaultConfig(name)` para configuracion sensata por defecto.
- Errores sentinel: `ErrCircuitOpen`, `ErrTooManyRequests`.
