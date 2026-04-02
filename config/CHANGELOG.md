# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/config` se registran aquí.

## [0.100.0] - 2026-04-02

### Changed

- Removed trivial struct field tests from `base_test.go`, kept `TestConnectionString`

### Added

- **Loader**: Carga configurable con Viper (opciones: path, nombre, tipo, prefijo de variables de entorno).
- **LoaderOption**: Constructor de opciones (WithPath, WithName, WithType, WithEnvPrefix).
- **Load method**: Cargar desde archivo tolerando ausencia, aplicar sobrescritura de variables de entorno.
- **LoadFromFile method**: Cargar desde archivo específico exigiendo su existencia.
- **BaseConfig**: Estructura base con campos: Server, Logger, Database, MongoDB, Bootstrap.
- **Validator**: Validación de configuración usando struct tags (required, min, max, email, etc.).
- **ValidationError**: Error estructurado con errores por campo (map[string][]string).
- Suite completa de tests unitarios con race detector.
- Documentación técnica detallada en docs/README.md con flujos comunes y ejemplos.
- Makefile con targets: build, test, test-race, check, lint, fmt, vet, tidy, deps, release.

### Design Notes

- Separación de responsabilidades: Loader carga, Validator valida.
- Tolerancia configurada: `Load()` es lenient (tolera ausencia), `LoadFromFile()` es estricto (exige existencia).
- Configuración base no universal: BaseConfig define shape estándar; servicios pueden extenderla.
- Variables de entorno: Mecanismo simple de override con SetEnvKeyReplacer (APP_DATABASE_HOST → database.host).
- Sin secretos: Módulo solo carga datos; secretos (API keys, passwords) resueltos externamente.
