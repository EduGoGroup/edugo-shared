# Config — Documentación técnica

Carga configuración YAML con Viper, la mezcla con variables de entorno y la valida con tags struct.

## Propósito

Proveer un mecanismo de carga y validación declarativa de configuración. El módulo define **cómo** se carga, no **qué** se carga: cada servicio define su propio struct de configuración.

## Componentes principales

### Loader — Carga configurable

Instancia local de Viper por llamada. No usa singleton global.

**Constructor:**
```go
loader := config.NewLoader(opts ...config.LoaderOption) *config.Loader
```

**Opciones disponibles:**

```go
config.WithConfigPath("./config")          // directorio de búsqueda (apilable)
config.WithConfigName("config")            // nombre del archivo sin extensión
config.WithConfigType("yaml")              // formato: yaml, json, toml
config.WithEnvPrefix("APP")               // prefijo para AutomaticEnv
config.WithEnvironmentOverride("dev")      // fusiona config-dev.yaml sobre config.yaml
config.WithEnvFiles(".env", "../.env")     // archivos .env a cargar antes de Viper
config.WithDefaults(map[string]interface{}{"server.port": 8080})
config.WithExplicitBindings(map[string]string{"database.password": "POSTGRES_PASSWORD"})
```

**Métodos de carga:**

```go
// Tolera ausencia del archivo; continúa con env vars y defaults
err := loader.Load(cfg)

// Exige que el archivo exista; falla si no lo encuentra
err := loader.LoadFromFile(cfg)
```

Ambos métodos son simétricos: aplican defaults, explicitBindings, envPrefix y AutomaticEnv.

**Acceso por key tras la carga:**

```go
loader.Get("server.port")         // any
loader.GetString("log.level")     // string
loader.GetInt("server.port")      // int
loader.GetBool("feature.enabled") // bool
```

Los métodos `Get*` retornan el zero value si se llaman antes de `Load()` o `LoadFromFile()`.

### Validator — Validación estructurada

```go
v := config.NewValidator()

// Valida un struct completo con tags validate:"..."
err := v.Validate(cfg)

// Valida un campo puntual
err := v.ValidateField("user@example.com", "email")
```

**ValidationError:**

```go
type ValidationError struct {
    Errors []FieldError
}

type FieldError struct {
    Field   string // nombre del campo
    Tag     string // tag que falló (required, min, oneof, ...)
    Value   any    // valor que se intentó validar
    Message string // mensaje legible
}
```

Tags soportados con mensaje legible: `required`, `min`, `max`, `oneof`, `email`, `url`. Cualquier otro tag retorna un mensaje genérico.

## Flujos comunes

### 1. Carga estándar con archivo YAML y env vars

```go
// internal/config/config.go

type Config struct {
    Environment string         `mapstructure:"app_env"      validate:"required,oneof=local dev qa prod"`
    ServiceName string         `mapstructure:"service_name" validate:"required"`
    Server      ServerConfig   `mapstructure:"server"`
    Database    DatabaseConfig `mapstructure:"database"`
}

func Load() (*Config, error) {
    loader := sharedconfig.NewLoader(
        sharedconfig.WithConfigPath("./config"),
        sharedconfig.WithEnvPrefix(""),
        sharedconfig.WithEnvFiles(".env", "../.env"),
        sharedconfig.WithEnvironmentOverride(os.Getenv("APP_ENV")),
    )

    cfg := &Config{}
    if err := loader.Load(cfg); err != nil {
        return nil, fmt.Errorf("loading config: %w", err)
    }

    v := sharedconfig.NewValidator()
    if err := v.Validate(cfg); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }

    return cfg, nil
}
```

### 2. Secrets sin prefijo con ExplicitBindings

Para secrets inyectados directamente como variables de entorno (sin prefijo del servicio):

```go
loader := sharedconfig.NewLoader(
    sharedconfig.WithConfigPath("./config"),
    sharedconfig.WithEnvPrefix("MYSERVICE"),
    sharedconfig.WithExplicitBindings(map[string]string{
        "database.postgres.password": "POSTGRES_PASSWORD",
        "messaging.rabbitmq.url":     "RABBITMQ_URL",
        "auth.jwt_secret":            "JWT_SECRET",
    }),
)
```

Con esto, `POSTGRES_PASSWORD` se mapea directamente a `database.postgres.password` independientemente del prefijo `MYSERVICE`.

### 3. Valores por defecto en memoria

```go
loader := sharedconfig.NewLoader(
    sharedconfig.WithConfigPath("./config"),
    sharedconfig.WithDefaults(map[string]interface{}{
        "server.port":         8080,
        "server.read_timeout": "30s",
        "log.level":           "info",
        "log.format":          "json",
    }),
)
```

Los defaults tienen la prioridad más baja: archivo YAML y env vars los sobrescriben.

### 4. Fusionar archivo de entorno

```go
// Carga config.yaml, luego fusiona config-prod.yaml encima
loader := sharedconfig.NewLoader(
    sharedconfig.WithConfigPath("./config"),
    sharedconfig.WithEnvironmentOverride("prod"),
)
```

Si `config-prod.yaml` no existe, se ignora silenciosamente.

### 5. Cargar solo desde variables de entorno (sin archivo)

`AutomaticEnv()` solo resuelve keys que Viper ya conoce. Sin archivo ni defaults, `Unmarshal` no sabe qué buscar. Usar `WithExplicitBindings` para este caso:

```go
loader := sharedconfig.NewLoader(
    sharedconfig.WithExplicitBindings(map[string]string{
        "environment":        "APP_ENV",
        "database.host":      "DB_HOST",
        "database.password":  "DB_PASSWORD",
    }),
)

cfg := &Config{}
_ = loader.Load(cfg) // tolera ausencia de archivo
```

### 6. Manejar errores de validación

```go
err := validator.Validate(cfg)

var valErr *config.ValidationError
if errors.As(err, &valErr) {
    for _, fe := range valErr.Errors {
        fmt.Printf("campo=%s tag=%s valor=%v mensaje=%s\n",
            fe.Field, fe.Tag, fe.Value, fe.Message)
    }
}
```

## Arquitectura

```
NewLoader(opciones)
    ↓
Load(cfg) / LoadFromFile(cfg)
    ├─ godotenv.Load(.env files)        ← solo Load()
    ├─ viper.New()                      ← instancia local, sin singleton
    ├─ SetConfigType / SetConfigName
    ├─ AddConfigPath (× n rutas)
    ├─ SetDefault (× n defaults)
    ├─ SetEnvPrefix + AutomaticEnv
    ├─ BindEnv (× n explicitBindings)
    ├─ ReadInConfig                     ← tolera ausencia en Load()
    ├─ MergeInConfig (env override)     ← opcional
    ├─ Unmarshal → struct
    └─ guarda instancia en l.viper
    ↓
Validator.Validate(cfg)
    ├─ validate.Struct(cfg)
    ├─ Construye []FieldError
    └─ Retorna *ValidationError
    ↓
cfg listo para usar
```

## Dependencias

- **Internas**: Ninguna (módulo independiente)
- **Externas**:
  - `github.com/spf13/viper` — Carga y parsing YAML/JSON/TOML
  - `github.com/go-playground/validator/v10` — Validación struct
  - `github.com/joho/godotenv` — Carga de archivos `.env`

## Testing

```bash
make test          # Tests básicos
make test-race     # Tests con race detector (verifica ausencia de race conditions)
make check         # Tests + linting + format
```

Dado que cada `Load()` usa `viper.New()`, los tests pueden correr en paralelo sin `viper.Reset()`.

## Notas de diseño

- **Sin struct base impuesta**: El módulo provee el mecanismo, no el esquema. Cada servicio define sus tipos según sus necesidades reales.
- **Sin singleton global**: Cada llamada crea su propia instancia de Viper → predecible, testeable, sin estado compartido.
- **Load vs LoadFromFile**: `Load()` es tolerante (archivo opcional, continúa con env vars). `LoadFromFile()` es estricto (archivo obligatorio).
- **ExplicitBindings para secrets**: Secrets inyectados sin prefijo del servicio deben vincularse explícitamente. `AutomaticEnv()` solo funciona para keys ya conocidas por Viper.
- **Sin secretos hardcodeados**: Passwords, tokens y API keys se inyectan via variables de entorno o `WithExplicitBindings`, nunca en YAML versionado.
