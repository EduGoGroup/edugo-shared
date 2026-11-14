# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/bootstrap/rabbitmq` se registran aqui.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Migración a la banda `0.900.x` (ADR 0010).
- `go.mod`: **eliminado** el `replace github.com/EduGoGroup/edugo-shared/bootstrap => ../`; el módulo requiere el tag publicado `edugo-shared/bootstrap v0.900.0` (HU-E1.2). Sin cambios de API.
- Dependencias directas: `github.com/rabbitmq/amqp091-go` `v1.10.0` → `v1.16.0`.
- Lint 2.14.0: `//nolint:errcheck` con motivo en el cierre best-effort de la ruta de error. Sin cambio de comportamiento.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.101.0] - 2026-04-02

### Added

- Factory RabbitMQ con connection timeout y QoS configurado.
- `NewFactory()` constructor para crear instancias de la factory.
- `CreateConnection(ctx, RabbitMQConfig) (*amqp.Connection, error)` con timeout de 10 segundos.
- `CreateChannel(*amqp.Connection) (*amqp.Channel, error)` con QoS prefetch count 10.
- `DeclareQueue(*amqp.Channel, string) (amqp.Queue, error)` con TTL 1h, max priority 10, lazy mode.
- `Close(*amqp.Channel, *amqp.Connection) error` con manejo seguro de `amqp.ErrClosed`.
- Targets Makefile: build, test, check, lint, fmt, vet, tidy, deps, release.

### Dependencies

- `github.com/EduGoGroup/edugo-shared/bootstrap` v0.101.0
- `github.com/rabbitmq/amqp091-go` v1.10.0
