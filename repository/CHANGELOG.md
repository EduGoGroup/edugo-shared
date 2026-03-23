# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/repository` deben registrarse aqui.

## [Unreleased]

## [0.4.6] - 2026-03-23

### Added

- Baseline de documentacion de fase 1 con `README.md`, `docs/README.md` y organizacion local por modulo.
- Nueva interfaz `MembershipAdminRepository` que extiende `MembershipRepository` con `FindBySchool`, permitiendo consultar membresías filtradas por `schoolID` con soporte de filtros, búsqueda y paginación. Solo consumidores con necesidades de administración deben referenciar esta interfaz.
- Nuevo constructor `NewPostgresMembershipAdminRepository` que retorna `MembershipAdminRepository`.

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-infrastructure/postgres` de `v0.65.0` a `v0.66.0`.
