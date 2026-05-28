# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/lifecycle/shutdown` se registran aqui.

## [Unreleased]

## [0.102.0] - 2026-04-03

### Added

- `GracefulShutdown` orquestador con ejecucion LIFO de tareas de cleanup.
- `Register(name, fn)` para registrar tareas (valida nil).
- `Shutdown(ctx)` con timeout y `errors.Join` para errores compuestos.
- `WaitForSignal()` para escuchar SIGTERM/SIGINT.
- `Logger` interface simple e inyectable.
