# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/repository` deben registrarse aqui.

## [Unreleased]

### Added

- Baseline de documentacion de fase 1 con `README.md`, `docs/README.md` y organizacion local por modulo.
- Nuevo método `FindBySchool` en la interfaz `MembershipRepository` y su implementación en `postgresMembershipRepository`, permitiendo consultar membresías filtradas por `schoolID` con soporte de filtros, búsqueda y paginación.

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-infrastructure/postgres` de `v0.65.0` a `v0.66.0`.
