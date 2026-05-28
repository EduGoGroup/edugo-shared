# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/bootstrap/mongodb` se registran aqui.

## [Unreleased]

## [0.101.0] - 2026-04-02

### Added

- Factory MongoDB con pool configurado y graceful degradation.
- `NewFactory()` constructor para crear instancias de la factory.
- `CreateConnection(ctx, MongoDBConfig) (*mongo.Client, error)` con pool size 100/10, timeouts y ping automatico.
- `GetDatabase(*mongo.Client, string) *mongo.Database` para obtener databases especificas.
- `Ping(ctx, *mongo.Client) error` con timeout de 5 segundos y read preference primaria.
- `Close(ctx, *mongo.Client) error` con timeout de 10 segundos.
- Connection timeout de 10 segundos por defecto.
- Targets Makefile: build, test, check, lint, fmt, vet, tidy, deps, release.

### Dependencies

- `github.com/EduGoGroup/edugo-shared/bootstrap` v0.101.0
- `go.mongodb.org/mongo-driver/v2` v2.5.0
