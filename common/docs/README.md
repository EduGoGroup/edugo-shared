# Common — Documentación técnica

Módulo base del repositorio: primitivos reutilizables, resolución de configuración, manejo de errores, validación y tipos compartidos.

## Propósito

Centralizar contratos y primitivos compartidos con dependencias mínimas, siendo la base sobre la cual se construyen otros módulos.

## Componentes principales

### common/config — Configuración de entorno

Resolución de variables de entorno con fallbacks y detección de ambiente.

**Funciones principales:**
- `GetEnv(key, defaultValue string) string` — Resuelve variable de entorno con fallback
- `GetEnvInt(key string, defaultValue int) int` — Resuelve como entero
- `GetEnvironment() string` — Retorna "dev", "staging" o "prod"
- `GetEnvBool(key string, defaultValue bool) bool` — Resuelve como booleano

### common/errors — Errores tipados

Define `AppError` con constructores tipados y mapeo automático a status HTTP.

**Constructores:**
- `NewValidationError(msg string) error` — 400 Bad Request
- `NewUnauthorizedError(msg string) error` — 401 Unauthorized
- `NewForbiddenError(msg string) error` — 403 Forbidden
- `NewNotFoundError(msg string) error` — 404 Not Found
- `NewConflictError(msg string) error` — 409 Conflict
- `NewInternalError(msg string, err error) error` — 500 Internal Server Error

**Interfaz:**
```go
type AppError interface {
    Error() string
    HTTPStatus() int
    Details() map[string]interface{}
}
```

### common/validator — Validación de datos

Agregación de múltiples errores de validación con helpers comunes.

**Métodos principales:**
- `NewValidator() *Validator`
- `RequireNotEmpty(field, value string) *Validator` — Validar campo no vacío
- `RequireLength(field string, value string, minLen, maxLen int) *Validator` — Rango de longitud
- `RequireEmail(field, email string) *Validator` — Formato de email
- `Require(condition bool, msg string) *Validator` — Validación personalizada
- `Valid() bool` — Verificar si todas las validaciones pasaron
- `Error() error` — Retornar error con todos los problemas

### common/types — Tipos compartidos

**UUID:**
- `NewUUID() string` — Generar UUID v4
- `ParseUUID(s string) (string, error)` — Parsear y validar UUID
- `IsValidUUID(s string) bool` — Verificar si es UUID válido

### common/types/enum — Enumeraciones de dominio

Constantes de roles, permisos, estados y tipos de evento compartidos en toda la aplicación.

**Roles:**
- `RoleAdmin`, `RoleSuperAdmin`, `RoleTeacher`, `RoleStudent`, etc.

**Permisos:**
- `PermissionUserRead`, `PermissionUserWrite`, `PermissionSchoolRead`, etc.

**Estados:**
- Estados de usuario, ciclo escolar, eventos, etc.

## Flujos comunes

### 1. Cargar configuración al inicializar

```go
func loadConfig() {
    dbHost := common.GetEnv("DB_HOST", "localhost")
    dbPort := common.GetEnvInt("DB_PORT", 5432)
    environment := common.GetEnvironment()

    if environment == "prod" {
        log.Println("Running in production")
    }
}
```

### 2. Validar entrada de usuario

```go
func validateUser(user *User) error {
    v := common.NewValidator()
    v.RequireNotEmpty("email", user.Email)
    v.RequireEmail("email", user.Email)
    v.RequireNotEmpty("password", user.Password)
    v.RequireLength("password", user.Password, 8, 72)

    if !v.Valid() {
        return v.Error() // AppError con todos los problemas
    }
    return nil
}
```

### 3. Manejar errores con mapeo HTTP

```go
func handler(w http.ResponseWriter, r *http.Request) {
    user, err := getUser(r.Context())

    if err != nil {
        appErr, ok := err.(common.AppError)
        if ok {
            w.WriteHeader(appErr.HTTPStatus())
            json.NewEncoder(w).Encode(appErr.Details())
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        return
    }

    // Éxito...
}
```

### 4. Usar tipos compartidos (UUID, Enums)

```go
// Generar ID único
userID := common.NewUUID()

// Asignar rol desde enum
user.Role = common.RoleTeacher

// Validar UUID recibido
id, err := common.ParseUUID(requestID)
if err != nil {
    return common.NewValidationError("invalid user id")
}
```

## Arquitectura

No existe un único package raíz; el consumo ocurre via subpaquetes específicos:

```
common
├── config/      → Configuración y entorno
├── errors/      → Errores tipados con mapeo HTTP
├── validator/   → Validación y agregación de errores
├── types/       → UUID y tipos compartidos
└── types/enum/  → Enumeraciones de dominio
```

Cada subpaquete tiene responsabilidad única y puede ser consumido independientemente.

## Dependencias

- **Internas**: Ninguna (módulo base, no depende de otros módulos de edugo-shared)
- **Externas**: `github.com/google/uuid` (mínimo necesario)

## Testing

Suite de tests comprensiva:
- Validación de errores tipados y mapeo HTTP
- Helpers de validación comunes
- UUID generation, parsing y validación
- Enumeraciones y constantes de dominio

Ejecutar:
```bash
make test          # Tests básicos
make test-race     # Tests con race detector
make check         # Tests + linting + format
```

## Notas de diseño

- **Modulo fundacional**: Es la base sobre la que otros módulos dependen; cambios aquí impactan toda la aplicación
- **Bajo acoplamiento**: Cada subpaquete puede usarse independientemente sin importar `common` completo
- **Contrato estable**: Los tipos y errores aquí son relativamente inmutables para evitar cambios masivos en consumidores
- **Sin dependencias circulares**: Por diseño, `common` no importa ningún otro módulo de edugo-shared
- **Enums centralizados**: Roles, permisos y estados viven aquí para evitar duplicación transversal
