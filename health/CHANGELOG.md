# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/health` se registran aqui.

## [Unreleased]

## [0.102.0] - 2026-04-03

### Added

- Framework extensible de health checks con estados Healthy/Unhealthy/Degraded.
- `Checker` que agrega multiples checks y retorna status consolidado.
- `PostgreSQLCheck` con metadata de pool (open_connections, in_use, idle).
- `MongoDBCheck` con ping y response_time_ms.
- `RabbitMQCheck` con timeout/context y goroutine para evitar bloqueos.
- Validacion de inputs en constructores (panic en nil, default 5s timeout).
- Fix: PostgreSQL degraded check ignora MaxOpenConnections==0 (unlimited pool).
