# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/repository` se registran aquí.

## [0.100.0] - 2026-04-02

### Added

- **ListFilters**: Estructura para filtros seguros con búsqueda ILIKE, paginación y validación de campos contra SQL injection.
- **ApplySearch**: Método que aplica búsqueda segura escapando caracteres especiales (%, _) en patrones SQL.
- **ApplyPagination**: Método que aplica LIMIT y OFFSET para paginación.
- **GetOffset**: Método para obtener desplazamiento (útil en cálculo de páginas).
- **UserRepository**: Interfaz con operaciones CRUD (Create, FindByID, FindByEmail, ExistsByEmail, Update, Delete, List).
- **SchoolRepository**: Interfaz con operaciones CRUD sobre entidades School.
- **MembershipRepository**: Interfaz con operaciones CRUD sobre entidades Membership (relaciones usuario-escuela).
- **MembershipAdminRepository**: Extensión de MembershipRepository con FindBySchool para consultas administrativas.
- **NewPostgresUserRepository**: Constructor que retorna UserRepository implementado con GORM.
- **NewPostgresSchoolRepository**: Constructor que retorna SchoolRepository implementado con GORM.
- **NewPostgresMembershipRepository**: Constructor que retorna MembershipRepository implementado con GORM.
- **NewPostgresMembershipAdminRepository**: Constructor que retorna MembershipAdminRepository implementado con GORM.
- **ErrNotFound**: Error tipado para operaciones que no encuentran registros.
- **Context propagation**: Todas las operaciones requieren context.Context para trazabilidad.
- **Field validation**: Validación de nombres de campo con regex `^[a-zA-Z_][a-zA-Z0-9_]*$` para prevenir SQL injection.
- **SQL injection prevention**: Escaping de patrones de búsqueda y validación de field names.
- Suite completa de tests unitarios con cobertura de CRUD, búsqueda segura y paginación.
- Documentación técnica detallada en docs/README.md con arquitectura, componentes y flujos comunes.
- Makefile con targets: build, test, test-race, check, lint, fmt, vet, tidy, deps, release.

### Design Notes

- Sin lógica de negocio: módulo proporciona CRUD genérico sin reglas específicas del dominio.
- Seguridad en búsqueda: ListFilters previene SQL injection validando campos y escapando patrones.
- Múltiples adaptadores: cada entidad tiene su propio repositorio intercambiable.
- Context-aware: propagación de contexto para trazabilidad en todas las operaciones.
- Errores tipados: ErrNotFound permite manejo explícito de registros no encontrados.
- GORM agnóstico: adaptadores usan *gorm.DB internamente sin exponerlo en interfaces públicas.

## [0.7.3] - 2026-03-30
### Changed
- Actualizada dependencia `edugo-infrastructure/postgres`

## [0.7.2] - 2026-03-28
### Changed
- Actualizada dependencia `edugo-infrastructure/postgres`

## [0.7.1] - 2026-03-28
### Changed
- Actualizada dependencia `edugo-infrastructure/postgres`

## [0.7.0] - 2026-03-27
### Changed
- Actualizada dependencia `edugo-infrastructure/postgres` a `v0.71.0`.
- Soporte para nuevas entidades de evaluación en los repositorios base.

## [0.4.6] - 2026-03-23

### Added

- Baseline de documentación con `README.md` y `docs/README.md`.
- Nueva interfaz `MembershipAdminRepository` que extiende `MembershipRepository` con `FindBySchool`.
- Nuevo constructor `NewPostgresMembershipAdminRepository` que retorna `MembershipAdminRepository`.

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-infrastructure/postgres` de `v0.65.0` a `v0.66.0`.
