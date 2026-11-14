# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/health` se registran aqui.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencias directas: `github.com/rabbitmq/amqp091-go` `v1.10.0` → `v1.16.0`, `github.com/stretchr/testify` `v1.11.1` → `v1.12.1`, `go.mongodb.org/mongo-driver/v2` `v2.5.0` → `v2.10.0`.
- 5 dependencias indirectas suben con `go mod tidy` (ver `go.mod`).
- Lint 2.14.0: aserciones de tipo comprobadas en tests y comentario del bloque de estados. Sin cambio de comportamiento.
- El módulo entra a la banda `v0.900.x` (ADR 0010); el tag anterior `v0.1.x` queda como procedencia.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.102.0] - 2026-04-03

### Added

- Framework extensible de health checks con estados Healthy/Unhealthy/Degraded.
- `Checker` que agrega multiples checks y retorna status consolidado.
- `PostgreSQLCheck` con metadata de pool (open_connections, in_use, idle).
- `MongoDBCheck` con ping y response_time_ms.
- `RabbitMQCheck` con timeout/context y goroutine para evitar bloqueos.
- Validacion de inputs en constructores (panic en nil, default 5s timeout).
- Fix: PostgreSQL degraded check ignora MaxOpenConnections==0 (unlimited pool).
