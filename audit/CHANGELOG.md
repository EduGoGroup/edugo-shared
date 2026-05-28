# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/audit` se registran aquí.

## [0.100.0] - 2026-04-01

### Added

- Contrato base `AuditLogger` para persistencia agnóstica de eventos.
- Estructura `AuditEvent` centralizada con campos de actor, acción, recurso, request y contexto.
- Funciones `AuditOption` declarativas: `WithChanges`, `WithSeverity`, `WithCategory`, `WithMetadata`, `WithPermission`, `WithError`.
- Implementación `NoopAuditLogger` para tests y desarrollo local.
- Constantes de severidad: `SeverityInfo`, `SeverityWarning`, `SeverityCritical`.
- Constantes de categoría: `CategoryAuth`, `CategoryData`, `CategoryConfig`, `CategoryAdmin`.
- Suite completa de tests (12 casos) validando contrato e implementación noop.
- Documentación en README.md y docs/README.md.
- Makefile con targets: build, test, check, lint, fmt, vet, tidy, deps, release.

### Design Notes

- Interfaz `AuditLogger` permite múltiples adaptadores (PostgreSQL, Kafka, archivos, etc.).
- Sin dependencias externas, solo stdlib Go.
- Versión v0.100.0 marca estabilización del contrato base.
