# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/common` se registran aquí.

## [0.100.0] - 2026-04-01

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
