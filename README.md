# EduGo Shared Package

Módulo compartido con utilidades y componentes reutilizables para los proyectos de EduGo.

## Estructura

```
shared/
├── pkg/
│   ├── logger/           # Logging interface y implementación con Zap
│   ├── database/         # Helpers de conexión a bases de datos
│   │   ├── postgres/     # PostgreSQL connection pool
│   │   └── mongodb/      # MongoDB connection
│   ├── messaging/        # RabbitMQ helpers (publisher, consumer)
│   ├── errors/           # Error handling personalizado
│   ├── validator/        # Validaciones comunes
│   ├── auth/             # JWT helpers y autenticación
│   ├── config/           # Configuration loaders
│   └── types/            # Tipos compartidos (UUID, Timestamp, Enums)
│       └── enum/         # Enumeraciones (Role, Status, etc.)
└── go.mod
```

## Uso en Proyectos

### 1. Agregar como dependencia local

En el `go.mod` de tu proyecto:

```go
require (
    github.com/edugo/shared v0.1.0
)

replace github.com/edugo/shared => ../../shared
```

### 2. Importar en código

```go
import (
    "github.com/edugo/shared/pkg/logger"
    "github.com/edugo/shared/pkg/database/postgres"
    "github.com/edugo/shared/pkg/auth"
)
```

### 3. Actualizar dependencias

```bash
go mod tidy
```

## Paquetes Disponibles

### Logger

Interface de logging con implementación Zap:

```go
logger := logger.NewZapLogger("info", "json")
logger.Info("mensaje", "key", "value")
logger.Error("error", "error", err)
```

### Database - PostgreSQL

Helper para conexión a PostgreSQL:

```go
db, err := postgres.Connect(postgres.Config{
    Host:     "localhost",
    Port:     5432,
    Database: "edugo",
    User:     "user",
    Password: "pass",
})
```

### Database - MongoDB

Helper para conexión a MongoDB:

```go
client, err := mongodb.Connect(mongodb.Config{
    URI:      "mongodb://localhost:27017",
    Database: "edugo",
})
```

### Messaging - RabbitMQ

Publisher y Consumer interfaces:

```go
publisher := messaging.NewPublisher(conn)
publisher.Publish(ctx, "exchange", "routing.key", payload)
```

### Errors

Errores personalizados con códigos:

```go
err := errors.NewNotFoundError("user not found")
err := errors.NewValidationError("invalid email")
err := errors.NewInternalError("database connection failed")
```

### Validator

Validaciones comunes:

```go
validator.IsValidEmail("test@example.com")
validator.IsValidUUID("123e4567-e89b-12d3-a456-426614174000")
```

### Auth - JWT

Generación y validación de JWT:

```go
token, err := auth.GenerateToken(userID, role, expiresIn)
claims, err := auth.ValidateToken(token)
```

### Types

Tipos compartidos:

```go
import "github.com/edugo/shared/pkg/types/enum"

role := enum.SystemRoleTeacher
status := enum.MaterialStatusPublished
```

## Versionamiento

Este paquete sigue [Semantic Versioning](https://semver.org/):

- **MAJOR**: Cambios incompatibles en la API
- **MINOR**: Nueva funcionalidad compatible con versiones anteriores
- **PATCH**: Corrección de bugs compatibles

## Desarrollo

### Agregar nueva funcionalidad

1. Crear nuevo paquete en `pkg/`
2. Implementar con interfaces cuando sea posible
3. Agregar tests unitarios
4. Actualizar este README
5. Hacer commit siguiendo conventional commits

### Tests

```bash
cd shared
go test ./...
```

### Formato y Lint

```bash
go fmt ./...
go vet ./...
golangci-lint run
```

## Contribuir

Al agregar nuevo código a `shared`:

1. Asegurarse que sea **realmente compartido** (usado por 2+ proyectos)
2. Documentar públicamente con comentarios Go
3. Agregar tests unitarios (coverage mínimo 80%)
4. Usar interfaces para flexibilidad
5. Evitar dependencias externas pesadas

## Licencia

Propietario - EduGo Project
