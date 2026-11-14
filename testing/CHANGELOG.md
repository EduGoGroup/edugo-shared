# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/testing` se registran aquí.

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencias directas: `github.com/rabbitmq/amqp091-go` `v1.10.0` → `v1.16.0`, `github.com/stretchr/testify` `v1.11.1` → `v1.12.1`, `github.com/testcontainers/testcontainers-go` `v0.42.0` → `v0.45.0`, `github.com/testcontainers/testcontainers-go/modules/mongodb` `v0.42.0` → `v0.45.0`, `github.com/testcontainers/testcontainers-go/modules/postgres` `v0.42.0` → `v0.45.0`, `github.com/testcontainers/testcontainers-go/modules/rabbitmq` `v0.42.0` → `v0.45.0`, `go.mongodb.org/mongo-driver/v2` `v2.5.0` → `v2.10.0`.
- 24 dependencias indirectas suben con `go mod tidy` (ver `go.mod`).
- El módulo entra a la banda `v0.900.x` (ADR 0010); el tag anterior `v0.1.x` queda como procedencia.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.100.0] - 2026-04-02

### Removed

- Deleted `helpers_unit_test.go` (invalid tests passing db=nil)

### Changed

- Added `//go:build integration` tag to `postgres_test.go`, `mongodb_test.go`, `rabbitmq_test.go`, `helpers_test.go`

### Added

- **ConfigBuilder**: Constructor fluido para configurar qué containers habilitar.
- **WithPostgres/WithMongoDB/WithRabbitMQ**: Métodos para habilitar/deshabilitar backends.
- **WithNetwork**: Método para especificar red Docker personalizada.
- **Manager**: Singleton que centraliza acceso y vida útil de todos los containers.
- **PostgresContainer**: Wrapper para PostgreSQL con GORM, health checks y utilities.
- **MongoDBContainer**: Wrapper para MongoDB con acceso client y health checks.
- **RabbitMQContainer**: Wrapper para RabbitMQ con acceso AMQP connection.
- **TruncatePostgres**: Truncar tablas PostgreSQL para limpiar estado entre tests.
- **DropMongoDB**: Dropear colecciones MongoDB para limpiar estado entre tests.
- **PurgeRabbitMQ**: Purgar colas RabbitMQ para limpiar estado entre tests.
- **WaitForHealthy**: Esperar a que containers estén listos con health checks.
- **RetryOperation**: Ejecutar operación con reintentos automáticos.
- **ExecSQL/ExecSQLFile**: Ejecutar SQL raw o desde archivo en PostgreSQL.
- **Cleanup**: Parar todos los containers y limpiar recursos Docker.
- Thread-safety: Protección interna con mutex para acceso concurrente.
- Singleton pattern: única instancia de Manager por proceso.
- Lazy initialization: containers se crean solo si están habilitados.
- Suite completa de tests unitarios (sin Docker) e integración (con Docker).
- Documentación técnica detallada en docs/README.md con componentes, flujos comunes y arquitectura.
- Makefile con targets: build, test, test-race, check, test-all.

### Design Notes

- Singleton Manager: centraliza vida útil de containers para abaratar suites largas.
- ConfigBuilder fluido: API clara para habilitar solo containers necesarios.
- Cleanup tipado: helpers específicos por backend (TruncatePostgres, DropMongoDB, PurgeRabbitMQ).
- Sin framework específico: funciona con testing, testify, ginkgo, etc.
- Health checks: retries automáticos para esperar a que containers estén listos.
- Docker agnóstico: Testcontainers abstrae detalles de Docker (compatible con Podman).

## [0.1.0] - 2026-03-26

### Added

- Baseline de documentación de fase 1 con `README.md` y `docs/README.md`.
