# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/bootstrap/postgres` se registran aqui.

## [Unreleased]

## [0.101.0] - 2026-04-02

### Added

- Factory PostgreSQL con soporte GORM, pgx SimpleProtocol y SearchPath configurable.
- `NewFactory()` constructor para crear instancias de la factory.
- `CreateRawConnection(ctx, PostgreSQLConfig) (*sql.DB, error)` para conexiones SQL nativas via pgx.
- `CreateGORMConnection(ctx, PostgreSQLConfig, ...GORMOption) (*gorm.DB, error)` para conexiones GORM completas.
- `Ping(ctx, *gorm.DB) error` para verificar conectividad.
- `Close(*gorm.DB) error` para cerrar conexiones.
- Soporte para `pgx.QueryExecModeSimpleProtocol` (PgBouncer/Neon) habilitado por defecto.
- Configuracion de `search_path` via `PostgreSQLConfig.SearchPath`.
- Pool configurable: `MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime`, `ConnMaxIdleTime` con defaults sensibles.
- Functional options: `WithGORMLogger()`, `WithSimpleProtocol()`, `WithPrepareStmt()`.
- Targets Makefile: build, test, check, lint, fmt, vet, tidy, deps, release.

### Dependencies

- `github.com/EduGoGroup/edugo-shared/bootstrap` v0.101.0
- `github.com/jackc/pgx/v5` v5.9.1
- `gorm.io/driver/postgres` v1.6.0
- `gorm.io/gorm` v1.31.1
