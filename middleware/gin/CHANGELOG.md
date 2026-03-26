# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/middleware/gin` deben registrarse aqui.

## [Unreleased]

## [0.52.0] - 2026-03-26

### Added

- `RequestLogging`: middleware de logging estructurado por request con request_id, correlation_id, method, path, client_ip.
- `PostAuthLogging`: middleware post-JWT que enriquece el logger del contexto con user_id, role y school_id.
- `GetLogger`/`GetRequestID`: helpers para extraer logger y request_id del gin.Context.
- `setLogger`: helper interno para inyectar logger en gin.Context y context.Context.
- Nil guard en `RequestLogging` cuando `baseLogger` es nil (fallback a `slog.Default()`).
- Fallback a `c.Request.URL.Path` cuando `c.FullPath()` retorna vacio (rutas no registradas / 404).
- Log level automatico segun status code: 5xx=Error, 4xx=Warn, 2xx/3xx=Info.

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-shared/logger` de `v0.50.1` a `v0.51.0`.
- Eliminado replace directive local para logger (ahora usa version publicada).

## [0.51.6] - 2026-03-23

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-shared/repository` de `v0.4.5` a `v0.4.6`.
- Actualización de dependencia indirecta `github.com/EduGoGroup/edugo-infrastructure/postgres` de `v0.65.0` a `v0.66.0`.

## [0.51.4] - 2026-03-17

### Added

- Baseline de documentacion de fase 1 con `README.md`, `docs/README.md` y organizacion local por modulo.

