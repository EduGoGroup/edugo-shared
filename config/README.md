# Config

Carga configuración YAML con Viper, la mezcla con variables de entorno y la valida con tags struct.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/config
```

El módulo se versiona y se consume de forma independiente gracias a su `go.mod` propio.

## Quick Start

### Definir el struct de configuración

Cada servicio define su propio struct. No hay struct base impuesta.

```go
type Config struct {
    Environment string         `mapstructure:"environment" validate:"required,oneof=local dev qa prod"`
    ServiceName string         `mapstructure:"service_name" validate:"required"`
    Server      ServerConfig   `mapstructure:"server"`
    Database    DatabaseConfig `mapstructure:"database"`
}
```

### Cargar configuración

```go
loader := config.NewLoader(
    config.WithConfigPath("./config"),
    config.WithEnvPrefix("APP"),
    config.WithEnvironmentOverride(os.Getenv("APP_ENV")), // fusiona config-dev.yaml si APP_ENV=dev
    config.WithEnvFiles(".env"),
)

cfg := &Config{}
if err := loader.Load(cfg); err != nil {
    return err
}
```

### Validar configuración

```go
validator := config.NewValidator()
if err := validator.Validate(cfg); err != nil {
    var valErr *config.ValidationError
    if errors.As(err, &valErr) {
        for _, fe := range valErr.Errors {
            fmt.Printf("%s: %s\n", fe.Field, fe.Message)
        }
    }
    return err
}
```

### Secrets sin prefijo (ExplicitBindings)

Útil para variables de entorno que no siguen el prefijo del servicio, como secretos inyectados por el orquestador:

```go
loader := config.NewLoader(
    config.WithConfigPath("./config"),
    config.WithEnvPrefix("MYAPP"),
    config.WithExplicitBindings(map[string]string{
        "database.password": "POSTGRES_PASSWORD",
        "auth.jwt_secret":   "JWT_SECRET",
    }),
)
```

## Componentes

- **Loader**: Carga configurable con Viper. El struct destino lo define el servicio.
- **Validator**: Validación con tags struct (`required`, `min`, `max`, `oneof`, `email`, `url`, etc.).
- **ValidationError**: Error estructurado con lista de `FieldError` (campo, tag, valor, mensaje).

## Opciones del Loader

| Opción | Descripción |
|--------|-------------|
| `WithConfigPath(path)` | Directorio donde buscar el archivo. Se puede usar varias veces. |
| `WithConfigName(name)` | Nombre del archivo sin extensión. Default: `"config"`. |
| `WithConfigType(type)` | Formato: `"yaml"`, `"json"`, `"toml"`. Default: `"yaml"`. |
| `WithEnvPrefix(prefix)` | Prefijo para variables de entorno (ej. `"APP"` → `APP_SERVER_PORT`). |
| `WithEnvironmentOverride(env)` | Fusiona `config-{env}.yaml` sobre el archivo base. |
| `WithExplicitBindings(map)` | Vincula keys de Viper a env vars específicas sin prefijo. |
| `WithDefaults(map)` | Valores por defecto en memoria antes de leer cualquier archivo. |
| `WithEnvFiles(files...)` | Archivos `.env` a cargar antes de que Viper actúe. |

## Métodos del Loader

| Método | Descripción |
|--------|-------------|
| `Load(cfg any) error` | Carga archivo + env vars. Tolera ausencia del archivo. |
| `LoadFromFile(cfg any) error` | Carga solo desde archivo. Falla si el archivo no existe. |
| `Get(key string) any` | Valor por key tras la última carga. |
| `GetString(key string) string` | Igual, tipado como string. |
| `GetInt(key string) int` | Igual, tipado como int. |
| `GetBool(key string) bool` | Igual, tipado como bool. |

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build      # Compilar módulo
make test       # Ejecutar tests
make test-race  # Tests con race detector
make check      # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- **Sin struct base impuesta**: El módulo provee el mecanismo de carga, no el esquema. Cada servicio define sus propios tipos.
- **Sin singleton global**: Cada `Load()` y `LoadFromFile()` crea una instancia local de Viper → seguro para tests paralelos.
- **Load vs LoadFromFile**: `Load()` tolera archivo ausente (continúa con env vars y defaults). `LoadFromFile()` exige que el archivo exista.
- **AutomaticEnv y keys conocidas**: `AutomaticEnv()` solo resuelve keys que Viper conoce (via archivo o defaults). Para env vars sin archivo usa `WithExplicitBindings`.
- **Sin secretos en código**: Passwords, tokens y API keys deben inyectarse via variables de entorno o `WithExplicitBindings`.
