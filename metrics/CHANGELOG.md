# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/metrics` deben registrarse aqui.

## [Unreleased]

## [0.2.0] - 2026-03-26

### Added

- Modulo completo de metricas con fachada `Metrics` y `NoopRecorder` por defecto.
- Metricas de autenticacion: login, token refresh (con histograma de duracion), rate limit, permission checks.
- Metricas de operaciones de negocio: CRUD con duracion y status.
- Metricas HTTP: request duration y response status.
- Metricas de base de datos: query duration y status.
- Metricas de messaging: message processing, DLQ routing, circuit breaker.
- Interfaz `Recorder` extensible para futuros backends (Prometheus, Datadog, OTel).
- Makefile con targets estandar (fmt, vet, lint, test, build).

