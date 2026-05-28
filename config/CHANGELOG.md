# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/config` se registran aquí.

## [0.101.0] - 2026-04-03

### Breaking Changes

- **Eliminado `BaseConfig`** y todos los tipos asociados (`ServerConfig`, `DatabaseConfig`, `MongoDBConfig`, `LoggerConfig`, `BootstrapConfig`, `OptionalResourcesConfig`). Cada servicio define su propio struct de configuración. El módulo provee el mecanismo de carga, no el esquema.

### Changed

- **`Load()`** ya no usa la instancia global de Viper. Crea una instancia local con `viper.New()` → seguro para tests paralelos, sin contaminación entre llamadas.
- **`LoadFromFile()`** ahora es completamente simétrico con `Load()`: aplica `defaults`, `explicitBindings`, `envPrefix` y `AutomaticEnv()`. Antes ignoraba estas opciones silenciosamente.
- **`Get()`, `GetString()`, `GetInt()`, `GetBool()`** usan la instancia de Viper guardada tras la última llamada a `Load()` o `LoadFromFile()`. Antes leían del singleton global, por lo que no funcionaban correctamente con `LoadFromFile()`.
- Tests reescritos: eliminados todos los `viper.Reset()` globales, uso de `t.Setenv()` para cleanup automático.

### Added

- Tests nuevos: `WithDefaults`, `WithExplicitBindings`, `LoadFromFile_WithDefaults`, `GetMethods_BeforeLoad`, `ParallelSafety`.

### Fixed

- `go.work`: corregida ruta de `edugo-api-identity` de `../EduApi/edugo-api-identity` → `./edugo-api-identity`.

### Notes

- `AutomaticEnv()` de Viper solo resuelve keys que ya conoce (via archivo o defaults). Para cargar configuración **únicamente desde variables de entorno** sin archivo YAML, usar `WithExplicitBindings`. Los servicios que usan archivo YAML + `WithEnvPrefix` no se ven afectados.

---

## [0.100.0] - 2026-04-02

### Changed

- Removed trivial struct field tests from `base_test.go`, kept `TestConnectionString`

### Added

- **Loader**: Carga configurable con Viper (opciones: path, nombre, tipo, prefijo de variables de entorno).
- **LoaderOption**: Constructor de opciones (WithConfigPath, WithConfigName, WithConfigType, WithEnvPrefix, WithEnvironmentOverride, WithExplicitBindings, WithDefaults, WithEnvFiles).
- **Load method**: Cargar desde archivo tolerando ausencia, aplicar sobrescritura de variables de entorno.
- **LoadFromFile method**: Cargar desde archivo específico exigiendo su existencia.
- **BaseConfig**: Estructura base con campos: Server, Logger, Database, MongoDB, Bootstrap.
- **Validator**: Validación de configuración usando struct tags (required, min, max, email, etc.).
- **ValidationError**: Error estructurado con errores por campo.
- Suite completa de tests unitarios con race detector.
- Documentación técnica detallada en docs/README.md con flujos comunes y ejemplos.
- Makefile con targets: build, test, test-race, check, lint, fmt, vet, tidy, deps, release.

### Design Notes

- Separación de responsabilidades: Loader carga, Validator valida.
- Tolerancia configurada: `Load()` es lenient (tolera ausencia), `LoadFromFile()` es estricto (exige existencia).
- Configuración base no universal: BaseConfig define shape estándar; servicios pueden extenderla.
- Variables de entorno: Mecanismo simple de override con SetEnvKeyReplacer (APP_DATABASE_HOST → database.host).
- Sin secretos: Módulo solo carga datos; secretos (API keys, passwords) resueltos externamente.
