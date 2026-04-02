# Config

Carga configuración YAML con Viper, la mezcla con variables de entorno y la valida con tags struct.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/config
```

El módulo se versiona y se consume de forma independiente gracias a su `go.mod` propio.

## Quick Start

### Cargar configuración con Loader

```go
// Crear loader con opciones
loader := config.NewLoader(
    config.WithPath("./config"),
    config.WithName("app"),
    config.WithType("yaml"),
    config.WithEnvPrefix("APP"),
)

// Cargar archivo (tolera ausencia)
cfg := &config.BaseConfig{}
err := loader.Load(cfg)
if err != nil {
    return err
}
```

### Validar configuración

```go
// Usar Validator para reglas struct
validator := config.NewValidator()
err := validator.Validate(cfg)
if err != nil {
    appErr := err.(config.ValidationError)
    for field, messages := range appErr.Errors {
        fmt.Printf("%s: %v\n", field, messages)
    }
}
```

### Sobrescribir con variables de entorno

```go
// Las variables de entorno sobrescriben valores del YAML
// APP_DATABASE_HOST=prod.db.com → BaseConfig.Database.Host
//
// El loader automáticamente:
// 1. Lee config.yaml
// 2. Aplica variables de entorno (AutomaticEnv)
// 3. Usa SetEnvKeyReplacer para mapear APP_DB_MAX_CONN → DatabaseConfig.MaxConn
```

## Componentes principales

- **Loader**: Carga configurable con Viper (path, nombre, tipo, prefijo de entorno)
- **BaseConfig**: Estructura base con Server, Logger, Database, MongoDB, Bootstrap
- **Validator**: Validación struct con tags (required, min, max, email, etc.)
- **ValidationError**: Error estructurado con errores por campo

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build    # Compilar módulo
make test     # Ejecutar tests
make test-race # Tests con race detector
make check    # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- **Flujo**: Load (leer archivo + variables de entorno) → Unmarshal → Validate
- **Tolerancia**: `Load()` tolera ausencia de archivo; `LoadFromFile()` exige que exista
- **Configuración base, no universal**: BaseConfig define shape estándar, no es obligatoria para todos los servicios
- **Variables de entorno**: Sobrescriben valores YAML mediante SetEnvKeyReplacer (APP_DATABASE_HOST → Database.Host)
