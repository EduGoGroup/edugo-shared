# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/testing` se registran aquí.

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
