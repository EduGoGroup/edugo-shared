# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/resilience/circuitbreaker` se registran aqui.

## [Unreleased]

## [0.102.0] - 2026-04-03

### Added

- Circuit breaker con 3 estados: Closed, Open, HalfOpen.
- Configuracion flexible: MaxFailures, Timeout, MaxRequests, SuccessThreshold.
- `MetricsHook` interface para inyeccion opcional de metricas.
- `Execute(ctx, fn)` para ejecutar operaciones protegidas.
- `DefaultConfig(name)` para configuracion sensata por defecto.
- Errores sentinel: `ErrCircuitOpen`, `ErrTooManyRequests`.
