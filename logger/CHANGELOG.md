# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/logger` se registran aquí.

## [0.100.0] - 2026-04-02

### Added

- **Logger interface**: Contrato común (Debug, Info, Warn, Error, Fatal, With, Sync).
- **ZapLogger**: Implementación con go.uber.org/zap (JSON o console format).
- **LogrusLogger**: Adaptador para github.com/sirupsen/logrus existentes.
- **SlogProvider**: Factory moderna para crear `*slog.Logger` con configuración desde ENV.
- **SlogAdapter**: Implementa Logger delegando a `*slog.Logger` para migración gradual.
- **Context helpers**: NewContext, FromContext, L para propagación de logger via context.Context.
- **Fields helpers**: Constructores tipados para atributos comunes (WithUserID, WithRequestID, WithDuration, etc.).
- **Log levels**: Soporte para debug, info, warn, error, fatal con configuración por string.
- **Campos variadicos**: Key/value pairs para minimizar acoplamiento.
- **With() returns new logger**: Composición inmutable de loggers contextuales.
- **Sync support**: Sincronización de buffers (crítico antes de exit).
- **Thread-safe**: Implementaciones seguras para concurrencia.
- Suite completa de tests unitarios con race detector.
- Documentación técnica detallada con flujos comunes.
- Makefile con targets: build, test, test-race, check, lint, fmt, vet, tidy, deps, release.

### Design Notes

- Sin políticas de logging: módulo proporciona abstracción e implementaciones, no reglas.
- Múltiples backends intercambiables: Zap, Logrus y Slog son opciones equivalentes según necesidad.
- Campos variadicos: minimiza acoplamiento entre consumidores e implementaciones.
- Migración gradual: SlogAdapter permite usar slog en código nuevo mientras se migra.
- Context-aware: propagación de logger via context.Context es práctica moderna recomendada.

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

