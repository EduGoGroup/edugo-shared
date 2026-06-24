# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/audit/postgres` se registran aquí.

## [0.900.0] - 2026-06-24

### Changed
- Migración a la banda `0.900.x` (**rescate** del tag `v0.1.0` re-tagueado; ver `../../AGENTS.md` §Versionado).
- `go.mod`: `require audit v0.900.0` y **eliminado** el `replace github.com/EduGoGroup/edugo-shared/audit => ../`
  (un `replace` no debe viajar en un módulo publicado: se ignora aguas abajo y ensucia el hash del go.mod).
  Sin cambios de API.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.100.0] - 2026-04-01

### Added

- Adaptador PostgreSQL para persistencia de eventos de auditoría usando GORM.
- Implementación de `audit.AuditLogger` con normalización automática de eventos.
- Constructor `NewPostgresAuditLogger(db, serviceName)` para inicializar el logger.
- Método `Log(ctx, event)` para persistencia directa de eventos.
- Método `LogFromGin(c, action, resourceType, resourceID, opts...)` para integración con Gin.
- Aplicación automática de defaults (Severity, Category, Actor).
- Conversión interna de `AuditEvent` a modelo GORM `AuditEventDB`.
- Soporte para campos opcionales (IP, UserAgent, ResourceID, etc.).
- Serialización JSON automática para `Changes` y `Metadata`.
- Estructura modular con API pública en root e implementación privada en `internal/`.
- Documentación completa en README.md y docs/README.md.
- Targets Makefile: build, test, check, lint, fmt, vet, tidy, deps, release.

### Dependencies

- `github.com/EduGoGroup/edugo-shared/audit` v0.1.0
- `gorm.io/gorm` v1.25.0+
- `github.com/gin-gonic/gin` v1.9.0+ (opcional)
