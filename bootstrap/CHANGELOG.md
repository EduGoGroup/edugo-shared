# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/bootstrap` se registran aqui.

## [Unreleased]

## [0.101.0] - 2026-04-02

### Changed

- **BREAKING**: Refactoring completo del modulo monolitico a sub-modulos por tecnologia.
- Modulo raiz reducido a configs, options, lifecycle y errores — 0 dependencias externas.
- Factory interfaces movidas a sus respectivos sub-modulos.
- Eliminado `Bootstrap()` orquestador (ningún consumidor lo usaba).
- Eliminado `Resources` struct, `Factories` struct, `BootstrapOptions`.
- Eliminadas interfaces `MessagePublisher`, `StorageClient`, `DatabaseClient`, `HealthChecker`, `LoggerFactory`.
- Eliminadas factory implementations del modulo raiz (movidas a sub-modulos).
- Eliminados `init_*.go`, `cleanup_registrars.go`, `config_extractors.go`, `helpers.go`, `health_check.go`.

### Added

- `config.go` — `PostgreSQLConfig` con `SearchPath`, pool config, `ConnMaxIdleTime`.
- `gorm_options.go` — `GORMOption`, `WithGORMLogger()`, `WithSimpleProtocol()`, `WithPrepareStmt()`.
- `lifecycle.go` — `LifecycleManager` interface tipada (reemplaza `any`).
- `errors.go` — `ErrMissingFactory`, `ErrConnectionFailed`.
- Sub-modulo `bootstrap/postgres` — Factory PostgreSQL + GORM con pgx SimpleProtocol y SearchPath.
- Sub-modulo `bootstrap/mongodb` — Factory MongoDB con pool configurado.
- Sub-modulo `bootstrap/rabbitmq` — Factory RabbitMQ con timeout y QoS.
- Sub-modulo `bootstrap/s3` — Factory S3 con soporte LocalStack.

### Removed

- 20+ archivos del modulo raiz (orquestador, factories, init functions, tests).
- 12 dependencias directas del go.mod raiz (AWS SDK, MongoDB, RabbitMQ, GORM, logrus, testcontainers).
- ~87 dependencias indirectas (Docker SDK, OpenTelemetry, etc.).

## [0.100.0] - 2026-04-02

### Changed

- Consolidated `cleanup_test.go`, `options_test.go`, `resources_test.go` into `bootstrap_test.go`
- Added `//go:build integration` tag to `factory_*_integration_test.go` and `bootstrap_integration_test.go`
