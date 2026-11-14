# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/lifecycle/shutdown` se registran aqui.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencias directas: `github.com/stretchr/testify` `v1.11.1` → `v1.12.1`.
- Lint 2.14.0: `else { if }` colapsado a `else if`. Sin cambio de comportamiento.
- El módulo entra a la banda `v0.900.x` (ADR 0010); el tag anterior `v0.1.x` queda como procedencia.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.102.0] - 2026-04-03

### Added

- `GracefulShutdown` orquestador con ejecucion LIFO de tareas de cleanup.
- `Register(name, fn)` para registrar tareas (valida nil).
- `Shutdown(ctx)` con timeout y `errors.Join` para errores compuestos.
- `WaitForSignal()` para escuchar SIGTERM/SIGINT.
- `Logger` interface simple e inyectable.
