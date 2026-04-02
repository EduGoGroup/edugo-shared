# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/database/postgres` deben registrarse aqui.

## [Unreleased]

## [0.100.0] - 2026-04-02

### Changed

- Removed trivial `TestDefaultConstants_Unit` from `connection_unit_test.go`
- Added `//go:build integration` tag to `connection_test.go`

### Added

- **PostgreSQL Connection**: Módulo de bajo nivel para conectar, validar y manejar transacciones en PostgreSQL
- **Dual API Support**: Soporte nativo (sql.DB) y GORM para máxima flexibilidad
- **Config Structure**: Configuración centralizada con host, user, password, database, SSL, pool sizing y timeouts
- **DefaultConfig()**: Valores por defecto seguros para desarrollo local (localhost:5432, postgres user, disable SSL, 25 conexiones max, 5 idle)
- **Connect()**: Establece conexión sql.DB nativa con validación de conectividad contra PostgreSQL
- **ConnectGORM()**: Establece conexión GORM con misma configuración de pool que sql.DB
- **HealthCheck()**: Verifica estado de conexión con timeout configurable (5s default)
- **GetStats()**: Retorna estadísticas del pool (OpenConnections, InUse, Idle, WaitCount, WaitDuration, etc)
- **Close()**: Cierra conexión de forma ordenada con validación nil-guard
- **WithTransaction()**: Ejecuta código dentro de transacción con rollback/commit automático
- **WithTransactionIsolation()**: Transacción con nivel de aislamiento específico (ReadCommitted, RepeatableRead, Serializable)
- **Timeouts configurables**: DefaultTimeout (10s para operaciones), DefaultConnectTimeout (10s), DefaultHealthCheckTimeout (5s)
- **Pool Management**: Configuración de MaxConnections (default 25), MaxIdleConnections (default 5), MaxLifetime (default 5 min)
- **Schema Search Path**: Soporte para múltiples schemas con SearchPath configurable
- **SSL Flexible**: Soporte SSLMode (disable, require, verify-ca, verify-full) para dev y producción
- **Error Handling**: Errores contextuales con wrapping para operaciones de conexión, ping, transacción
- **Panic Protection**: Transacciones con rollback automático en caso de panic (defer protection)
- **Concurrency Safe**: Todas las operaciones thread-safe, *sql.DB y *gorm.DB compartibles entre goroutines
- **Defensive Defaults**: Nil-guard en Close, timeouts en todas las operaciones, validación temprana en Connect
- **Suite completa de tests unitarios e integración** sin dependencias externas para unitarios
- **Documentación técnica detallada** en docs/README.md con componentes, flujos comunes, ciclo de vida y patrones
- **Makefile** con targets: fmt, vet, lint, test, build, check

### Design Notes

- **Dual API**: Expone sql.DB nativo y GORM, usuario decide qué abstracción usar según necesidades
- **No elige por ti**: No fuerza una opción ORM, soporta ambas con igual configuración de pool
- **Bajo nivel**: Módulo de conexión solamente, no abstrae repositorios o query builders
- **Agnóstico**: No depende de otros módulos EduGo, funciona autocontendido
- **Pool configurable**: MaxConnections, MaxIdleConnections, MaxLifetime adaptables a carga esperada
- **Validación temprana**: Connect valida inmediatamente con Ping, no lazy connection
- **Transacciones seguras**: WithTransaction rollback automático en error o panic con defer protection
- **Aislamiento flexible**: WithTransactionIsolation permite control fino de concurrencia (dirty reads vs serializable)
- **Monitoring**: GetStats permite observabilidad del pool (conexiones activas, waiting, closed stats)
- **Search path**: Soporta múltiples schemas PostgreSQL en SearchPath configurable
- **SSL flexible**: SSLMode cubre desarrollo (disable) a producción (verify-full)
- **Timeouts defensivos**: ConnectTimeout y HealthCheckTimeout evitan bloqueos indefinidos
- **Concurrency**: pool thread-safe, seguro compartir *sql.DB entre goroutines
