# Testing Containers

Este directorio contiene el package `containers` del modulo [`testing`](../README.md). Su documentacion principal ahora vive en:

- [README del modulo](../README.md)
- [Documentacion tecnica del modulo](../docs/README.md)
- [Changelog del modulo](../CHANGELOG.md)

## Alcance de este package

- Builder de configuracion para PostgreSQL, MongoDB y RabbitMQ.
- `Manager` singleton para reutilizar containers entre tests.
- Wrappers por backend con helpers de conexion y limpieza.
- Helpers transversales como `ExecSQLFile`, `WaitForHealthy` y `RetryOperation`.

## Cuando leer este archivo

- Cuando necesites ubicar rapidamente el package real a importar: `github.com/EduGoGroup/edugo-shared/testing/containers`.
- Cuando quieras entrar directo a los wrappers de containers sin pasar por la documentacion general del repositorio.

## Nota

La explicacion completa de procesos, arquitectura y operacion del modulo se centralizo en `testing/docs/README.md` para evitar duplicacion y enlaces rotos.
