# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/resilience/ratelimiter` se registran aqui.

## [Unreleased]

## [0.102.0] - 2026-04-03

### Added

- Token bucket rate limiter con `Allow()` (no-blocking) y `Wait(ctx)` (blocking).
- `MultiRateLimiter` para throttling por entidad con creacion lazy de limiters.
- Timer reutilizable en `Wait()` con calculo exacto de duracion de espera.
- Soporte de burst capacity configurable.
