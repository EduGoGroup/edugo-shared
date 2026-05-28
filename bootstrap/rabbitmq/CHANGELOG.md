# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/bootstrap/rabbitmq` se registran aqui.

## [Unreleased]

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
