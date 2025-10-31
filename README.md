# EduGo Shared Library

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/EduGoGroup/edugo-shared)](https://github.com/EduGoGroup/edugo-shared/releases)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/coverage-87.2%25-brightgreen)](https://github.com/EduGoGroup/edugo-shared)

Professional Go shared library with utilities and reusable components for EduGo projects.

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

## 📦 Installation

### Latest Stable Version (Recommended)

```bash
go get github.com/EduGoGroup/edugo-shared@v1.0.0
```

### Latest Version

```bash
go get github.com/EduGoGroup/edugo-shared@latest
```

### Specific Version

```bash
# List available versions
go list -m -versions github.com/EduGoGroup/edugo-shared

# Install specific version
go get github.com/EduGoGroup/edugo-shared@v1.0.0
```

## 🚀 Quick Start
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

### Comandos Make Disponibles

Este proyecto incluye un Makefile con comandos útiles para desarrollo:

```bash
# Ver todos los comandos disponibles
make help

# Configurar entorno de desarrollo
make setup

# Comandos básicos
make build          # Compilar proyecto
make test           # Ejecutar tests
make fmt            # Formatear código
make lint           # Ejecutar linter
make vet            # Análisis estático

# Tests avanzados
make test-race      # Tests con detección de race conditions
make test-coverage  # Tests con reporte de cobertura
make benchmark      # Ejecutar benchmarks

# Verificaciones completas
make check-all      # Todas las verificaciones
make ci             # Pipeline CI completo
make pre-commit     # Verificaciones rápidas antes de commit

# Herramientas
make install-tools  # Instalar herramientas de desarrollo
make docs-serve     # Servir documentación en localhost:6060
make security       # Análisis de seguridad

# Utilidades
make clean          # Limpiar archivos generados
make deps           # Actualizar dependencias
make version        # Mostrar versiones
```

### Agregar nueva funcionalidad

1. Crear nuevo paquete en `pkg/`
2. Implementar con interfaces cuando sea posible
3. Agregar tests unitarios
4. Actualizar este README
5. Ejecutar `make pre-commit` antes del commit
6. Hacer commit siguiendo conventional commits

### Tests

```bash
# Tests básicos
make test

# Tests con cobertura
make test-coverage

# Tests completos (con race detection)
make test-race
```

### Formato y Lint

```bash
# Formatear código
make fmt

# Análisis estático
make vet

# Linter completo
make lint

# Todo junto
make check-all
```

## Contribuir

Al agregar nuevo código a `shared`:

1. Asegurarse que sea **realmente compartido** (usado por 2+ proyectos)
2. Documentar públicamente con comentarios Go
3. Agregar tests unitarios (coverage mínimo 80%)
4. Usar interfaces para flexibilidad
5. Evitar dependencias externas pesadas

## 📋 Versioning

This project follows [Semantic Versioning](https://semver.org/). 

### Version History

- **v1.0.0** (2025-10-31): First stable release with complete feature set
- **v0.1.0**: Initial development version

### Compatibility Promise

Starting from v1.0.0, we guarantee:

- ✅ **Backward compatibility** for all PATCH and MINOR releases
- ✅ **API stability** - no breaking changes without major version bump
- ✅ **Clear migration guides** for any major version changes

### Upgrade Guide

See [UPGRADE_GUIDE.md](UPGRADE_GUIDE.md) for detailed instructions on updating your projects.

## 📚 Documentation

- **[CHANGELOG.md](CHANGELOG.md)**: Detailed change history
- **[UPGRADE_GUIDE.md](UPGRADE_GUIDE.md)**: Migration instructions
- **API Documentation**: Available via `go doc` or [pkg.go.dev](https://pkg.go.dev/github.com/EduGoGroup/edugo-shared)

## 🤝 Contributing

1. Follow semantic versioning principles
2. Update CHANGELOG.md for all changes
3. Maintain backward compatibility
4. Add comprehensive tests
5. Update documentation

## 📄 License

MIT License - EduGo Project
