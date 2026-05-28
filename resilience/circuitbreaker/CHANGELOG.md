# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/resilience/circuitbreaker` se registran aqui.

## [Unreleased]

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
