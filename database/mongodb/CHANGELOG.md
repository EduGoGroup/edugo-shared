# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/database/mongodb` deben registrarse aqui.

## [Unreleased]

## [0.100.0] - 2026-04-02

### Added

- **MongoDB Connection**: Módulo de bajo nivel para conectar, validar y cerrar conexiones a MongoDB
- **Config Structure**: Configuración centralizada con URI, database, timeouts y pool sizing (MaxPoolSize, MinPoolSize)
- **DefaultConfig()**: Valores por defecto seguros para desarrollo local (localhost:27017, test database, 10s timeout, 10-100 pool)
- **Connect()**: Establece conexión a MongoDB con validación de conectividad contra read preference Primary
- **GetDatabase()**: Obtiene instancia de database del cliente conectado
- **HealthCheck()**: Verifica estado de conexión con timeout configurable (5s default)
- **Close()**: Cierra conexión de forma ordenada con drain de operaciones in-flight
- **Timeouts configurables**: DefaultTimeout (10s), DefaultHealthCheckTimeout (5s), DefaultDisconnectTimeout (10s)
- **Pool Management**: Configuración de MaxPoolSize (default 100) y MinPoolSize (default 10)
- **Error Handling**: Errores contextuales con wrapping ("failed to connect to mongodb: %w", "failed to ping mongodb: %w")
- **Concurrency Safe**: Todas las operaciones thread-safe, Client y Database compartibles
- **Connection Validation**: Ping en Connect valida que al menos una primary sea accesible
- **Defensive Defaults**: Nil-guard en Close, timeout en todas las operaciones
- **Suite completa de tests unitarios e integración** sin dependencias externas para unitarios
- **Documentación técnica detallada** en docs/README.md con componentes, flujos comunes, ciclo de vida y patrones
- **Makefile** con targets: fmt, vet, lint, test, build, check

### Design Notes

- **Bajo nivel**: Módulo de conexión solamente, no abstrae colecciones o operaciones de datos
- **Agnóstico**: No depende de other EduGo modules, funciona de forma autocontendida
- **Pool configurable**: MaxPoolSize y MinPoolSize adaptables a carga esperada de aplicación
- **Validación temprana**: Connect valida conectividad inmediatamente (Ping), no lazy
- **Timeouts defensivos**: Evitan bloqueos indefinidos en Connect, HealthCheck, Close
- **Health checking**: HealthCheck usa Ping con timeout más corto (5s vs 10s de operaciones)
- **Concurrency**: Driver Go MongoDB es thread-safe, seguro compartir Client entre goroutines
- **Lazy database**: Database existe solo en operaciones reales, GetDatabase no valida existencia

## [0.1.0] - [Fecha anterior]

### Added

- Baseline de documentacion de fase 1 con `README.md`, `docs/README.md` y organizacion local por modulo.
