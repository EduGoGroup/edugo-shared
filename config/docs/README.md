# Config — Documentación técnica

Carga configuración YAML con Viper, la mezcla con variables de entorno y la valida con tags struct.

## Propósito

Proporcionar carga y validación declarativa de configuración usando YAML con soporte para sobrescritura via variables de entorno.

## Componentes principales

### Loader — Carga configurable

Constructor con opciones para path, nombre de archivo, tipo (YAML/JSON) y prefijo de variables de entorno.

**Métodos principales:**
- `NewLoader(options ...LoaderOption) *Loader` — Constructor con opciones
- `Load(v interface{}) error` — Cargar desde archivo (tolera ausencia) y sobrescribir con env vars
- `LoadFromFile(path string, v interface{}) error` — Cargar desde archivo específico (exige existencia)

**Opciones (LoaderOption):**
- `WithPath(path string)` — Directorio donde buscar archivo
- `WithName(name string)` — Nombre del archivo sin extensión (ej. "app" → "app.yaml")
- `WithType(typ string)` — Tipo: "yaml", "json", "toml", etc.
- `WithEnvPrefix(prefix string)` — Prefijo para variables de entorno (ej. "APP_")

### BaseConfig — Estructura de configuración

Estructura base que define configuración común para servicios: servidor, logger, base de datos y bootstrap.

**Campos principales:**
- `Server` — Puerto, host, timeouts HTTP
- `Logger` — Nivel, formato (JSON/text), salida
- `Database` — PostgreSQL: host, puerto, usuario, contraseña, base de datos, pool
- `MongoDB` — (opcional) host, puerto, base de datos
- `Bootstrap` — Configuración de inicialización específica del servicio

### Validator — Validación estructurada

Valida configuración usando tags struct (`required`, `min`, `max`, `email`, etc.).

**Métodos:**
- `NewValidator() *Validator`
- `Validate(v interface{}) error` — Retorna ValidationError si hay problemas

**ValidationError:**
```go
type ValidationError struct {
    Errors map[string][]string // campo -> lista de errores
}
```

## Flujos comunes

### 1. Cargar y validar configuración al inicializar

```go
func initConfig() (*config.BaseConfig, error) {
    // Crear loader
    loader := config.NewLoader(
        config.WithPath("./config"),
        config.WithName("app"),
        config.WithType("yaml"),
        config.WithEnvPrefix("APP"),
    )

    // Cargar configuración
    cfg := &config.BaseConfig{}
    if err := loader.Load(cfg); err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    // Validar
    validator := config.NewValidator()
    if err := validator.Validate(cfg); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }

    return cfg, nil
}
```

### 2. Sobrescribir archivo YAML con variables de entorno

```yaml
# config/app.yaml
server:
  port: 8080
  host: localhost
database:
  host: localhost
  port: 5432
  database: edugo
```

```bash
# Al ejecutar:
export APP_SERVER_PORT=9000
export APP_DATABASE_HOST=prod.db.com

# El Loader:
# 1. Lee config/app.yaml
# 2. Aplica AutomaticEnv() para buscar APP_* variables
# 3. Sobrescribe: server.port = 9000, database.host = prod.db.com
```

### 3. Manejar errores de validación

```go
validator := config.NewValidator()
err := validator.Validate(cfg)

if validationErr, ok := err.(config.ValidationError); ok {
    for field, messages := range validationErr.Errors {
        fmt.Printf("Field %q has errors: %v\n", field, messages)
        // Output: Field "Database.Password" has errors: ["required"]
    }
}
```

### 4. Cargar desde archivo específico

```go
// LoadFromFile exige que el archivo exista
err := loader.LoadFromFile("./config/prod.yaml", cfg)
if err != nil {
    // El archivo no existe o hay error de parsing
    return err
}
```

## Arquitectura

Flujo de carga y validación:

```
1. NewLoader(opciones)
   ↓
2. Load(cfg)
   ├─ Lee archivo YAML
   ├─ Aplica variables de entorno
   ├─ Unmarshal a struct
   └─ Retorna error si falla parsing
   ↓
3. Validator.Validate(cfg)
   ├─ Valida tags struct (required, min, max, etc.)
   └─ Retorna ValidationError si hay problemas
   ↓
4. cfg listo para usar
```

## Dependencias

- **Internas**: Ninguna (módulo independiente)
- **Externas**:
  - `github.com/spf13/viper` — Carga y parsing YAML/JSON
  - `github.com/go-playground/validator/v10` — Validación struct

## Testing

Suite de tests completa:
- Carga de archivos YAML/JSON
- Sobrescritura de variables de entorno
- Validación de estructura
- Manejo de errores
- Tolerancia a ausencia de archivo en Load

Ejecutar:
```bash
make test          # Tests básicos
make test-race     # Tests con race detector
make check         # Tests + linting + format
```

## Notas de diseño

- **Separación de responsabilidades**: Loader carga, Validator valida
- **Tolerancia configurada**: `Load()` es lenient (tolera ausencia), `LoadFromFile()` es estricto
- **Configuración base, no universal**: BaseConfig define shape estándar; servicios pueden extenderla
- **Variables de entorno**: Mecanismo de override simple: APP_DATABASE_HOST → database.host
- **Sin secretos**: Este módulo solo carga datos; secretos (API keys, passwords) resueltos externamente
