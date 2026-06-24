# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/audit` se registran aquí.

## [0.900.0] - 2026-06-24

### Changed
- Migración a la banda `0.900.x` (estándar de versionado del ecosistema; ver `../AGENTS.md` §Versionado).
  **Rescate** del tag `v0.1.0`, que fue re-tagueado (contenido cambiado bajo el mismo número) y dejó a los
  consumidores con `checksum mismatch` (`SECURITY ERROR`) en CI/cloud. Mismo contrato y código que el
  `v0.1.0` actual; solo cambia el número a un tag nuevo e inmutable. Sin cambios de API.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

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
