# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/audit/postgres` se registran aquí.

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
