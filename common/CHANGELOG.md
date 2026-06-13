# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/common` se registran aquí.

## [Unreleased]

## [0.900.3] - 2026-06-13

### Added
- Nuevo package `common/timeutil`: helpers de tiempo para el estándar de fechas UTC (bug 0001, MP-05 F1). `NowUTC()`, `FormatISO()`/`ParseISO()` (instantes en UTC con sufijo `Z`) y `FormatDate()`/`ParseDate()` (fechas puras `YYYY-MM-DD` sin zona). Tests round-trip verdes. Cambio aditivo.

### Removed
- **BREAKING** — Eliminadas las 12 constantes de permisos de features muertas en `types/enum/permission.go` (defs + entradas en el mapa `AllPermissions`): `PermissionColors*`, `PermissionSchedules*` y `PermissionCalendarEvents*` (CRUD `color`/`schedule`/`calendar_event` podados en platform, MP-01 F3). Cualquier consumidor que referencie estas constantes deja de compilar.

## [0.900.2] - 2026-06-11

### Added
- Tipo `Scope` y constante `ScopeNotificationsDispatch` (`notifications.dispatch`) en `types/enum/scope.go`, con catálogo cerrado `AllScopes` y método `IsValid()`. Es el scope M2M que autoriza a un cliente de servicio a invocar el Notification Gateway (`POST /api/v1/internal/notifications/dispatch`); única fuente de verdad del scope, análoga a `Permission*` (plan 020 N5, D14/D17).

## [0.900.1] - 2026-06-07

### Added
- Permisos del recurso `academic.grades_detail` (modo DETALLADO de notas, N4 / ADR 0020) en `types/enum/permission.go` y en `AllPermissions`:
  - `academic.grades_detail.create` (enum `PermissionGradesDetailCreate`).
  - `academic.grades_detail.read` (enum `PermissionGradesDetailRead`).
  - `academic.grades_detail.update` (enum `PermissionGradesDetailUpdate`).
  - `academic.grades_detail.delete` (enum `PermissionGradesDetailDelete`).

  Gestionan los componentes de nota (`academic.grade_item`) y habilitan el desglose transparente en "Mis Notas". El modo BÁSICO usa solo `academic.grades`; el modo DETALLADO los otorga vía grant condicional según el perfil de la escuela (`academic.schools.grade_profile`).

## [0.900.0] - 2026-06-06

### Added
- Permiso `academic.my_grades.read:own` (enum `PermissionMyGradesReadOwn`) para que el alumno consulte sus propias notas vía `GET /me/grades` (N3 F4 'Mis Notas').

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.100.0] - 2026-04-02

### Removed

- `AllSystemRoles()`, `AllSystemRolesStrings()` from `types/enum/role.go` (dead code, no external consumers)
- `AllMaterialStatuses()`, `AllProgressStatuses()`, `AllProcessingStatuses()` from `types/enum/status.go` (dead code)
- `AllEventTypes()` from `types/enum/event.go` (dead code)
- `AllAssessmentTypes()` from `types/enum/assessment.go` (dead code)
- `AllPermissionsSlice()` from `types/enum/permission.go` (dead code)
- Corresponding test functions for all removed functions

### Added

- **config**: Resolución de variables de entorno (GetEnv, GetEnvInt, GetEnvBool) y detección de ambiente (dev/staging/prod).
- **errors**: AppError tipado con constructores (ValidationError, UnauthorizedError, ForbiddenError, NotFoundError, ConflictError, InternalError) y mapeo automático a status HTTP.
- **validator**: Validador fluido (NewValidator) con helpers: RequireNotEmpty, RequireLength, RequireEmail, Require; agregación de múltiples errores.
- **types**: UUID generation (NewUUID), parsing (ParseUUID) y validación (IsValidUUID).
- **types/enum**: Enumeraciones de dominio: roles (Admin, SuperAdmin, Teacher, Student, etc.), permisos (UserRead, UserWrite, SchoolRead, etc.), estados.
- Suite completa de tests unitarios con race detector.
- Benchmarks para config, errors, validator y UUID.
- Documentación en README.md y docs/README.md.
- Makefile con targets: build, test, test-race, check, lint, fmt, vet, tidy, deps, release.

### Design Notes

- Módulo base sin dependencias circulares a otros módulos de edugo-shared.
- Cada subpaquete (config, errors, validator, types) puede consumirse independientemente.
- Enumeraciones de dominio centralizadas para evitar duplicación transversal.
- AppError mapea automáticamente a status HTTP para facilitar respuestas API.
- Validator acumula múltiples errores en una sola validación para mejor UX.
