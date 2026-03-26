# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/logger` deben registrarse aqui.

## [Unreleased]

## [0.51.0] - 2026-03-26

### Added

- `SlogProvider`: factory para crear `*slog.Logger` con backend JSON/text y campos base (service, env, version).
- `NewSlogProviderFromEnv()`: factory que lee configuracion de variables de entorno.
- `SlogAdapter`: adapter que implementa `Logger` delegando a `*slog.Logger` para migracion gradual.
- `NewContext`/`FromContext`/`L`: helpers para propagar `*slog.Logger` via `context.Context`.
- Helpers tipados `slog.Attr`: `WithRequestID`, `WithUserID`, `WithCorrelationID`, `WithError`, `WithDuration`, `WithComponent`, `WithSchoolID`, `WithRole`, `WithResource`, `WithResourceID`, `WithAction`, `WithIP`.
- Constantes nuevas: `FieldRole`, `FieldSchoolID`, `FieldBytes`.
- Nil guards en `NewSlogAdapter` y `WithError` para prevenir panics.
- Benchmarks para provider, adapter, context y helpers.

