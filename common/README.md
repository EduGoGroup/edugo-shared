# Common

Módulo base del repositorio: primitivos reutilizables, resolución de configuración, manejo de errores, validación y tipos compartidos.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/common
```

El módulo se distribuye via subpaquetes: `common/config`, `common/errors`, `common/validator`, `common/types` y `common/types/enum`.

## Quick Start

### Configuración de entorno

```go
// Resolver variable de entorno con fallback
dbHost := common.GetEnv("DB_HOST", "localhost")
maxConnections := common.GetEnvInt("DB_MAX_CONN", 10)

// Detectar ambiente actual
env := common.GetEnvironment() // "dev", "staging", "prod"
```

### Manejo de errores tipados

```go
// Errores de validación
if user.Email == "" {
    return common.NewValidationError("email is required")
}

// Errores de autenticación
if !tokenValid {
    return common.NewUnauthorizedError("invalid token")
}

// Mapeo automático a status HTTP
statusCode := error.HTTPStatus() // 400 para validation, 401 para unauthorized, etc.
```

### Validación de datos

```go
// Acumular múltiples errores de validación
v := common.NewValidator()
v.RequireNotEmpty("email", user.Email)
v.RequireNotEmpty("password", user.Password)
v.RequireLength("password", user.Password, 8, 72)
if !v.Valid() {
    return v.Error() // Retorna error con todos los problemas
}
```

### Tipos compartidos (UUID, Enums)

```go
// Generar UUID
id := common.NewUUID() // string

// Enums de dominio (roles, permisos, estados)
role := common.RoleAdmin
permission := common.PermissionUserRead
```

## Componentes principales

- **config**: Resolución de variables de entorno y detección de ambiente
- **errors**: Errores tipados (Validation, Unauthorized, NotFound, etc.) con mapeo a HTTP
- **validator**: Agregación de errores de validación y helpers de validación comunes
- **types**: UUID y tipos compartidos
- **types/enum**: Enums de dominio (roles, permisos, estados, tipos de evento)

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

- No existe un package raíz único; consumir mediante subpaquetes específicos (`common/errors`, `common/config`, etc.)
- Dependencias mínimas por diseño: solo `github.com/google/uuid` como dependencia externa
- Módulo fundacional: otros módulos dependen de estos contratos
- Arquitectura: bajo acoplamiento, alta reutilización, contrato estable
