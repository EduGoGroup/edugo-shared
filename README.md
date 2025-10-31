# EduGo Shared Library

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/EduGoGroup/edugo-shared)](https://github.com/EduGoGroup/edugo-shared/releases)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Professional Go shared library with modular architecture and reusable components for EduGo projects.

## 🏗️ Arquitectura Modular (v2.0.5+)

Este proyecto utiliza **módulos Go independientes** para optimizar dependencias y permitir instalación selectiva.

### ✨ Beneficios

- ✅ **Dependencias selectivas**: Solo descarga las librerías que necesitas
- ✅ **Binarios optimizados**: Menor tamaño del ejecutable final
- ✅ **Menor superficie de ataque**: Menos dependencias = menos vulnerabilidades
- ✅ **Builds más rápidos**: Menos código que compilar
- ✅ **Testing modular**: Tests aislados por módulo

### 📦 Módulos Disponibles

```
edugo-shared/
├── common/                    # Errors, Types, Validator, Config
│   └── go.mod                # Deps: google/uuid (liviano)
├── logger/                    # Logging con Zap
│   └── go.mod                # Deps: go.uber.org/zap
├── auth/                      # JWT Authentication
│   └── go.mod                # Deps: jwt, uuid, common
├── messaging/
│   └── rabbit/                # RabbitMQ helpers
│       └── go.mod            # Deps: rabbitmq/amqp091-go
└── database/
    ├── postgres/              # PostgreSQL utilities
    │   └── go.mod            # Deps: lib/pq
    └── mongodb/               # MongoDB utilities
        └── go.mod            # Deps: mongo-driver
```

---

## 📦 Instalación

### Módulo Common (Errors, Types, Validator, Config)

**El más liviano - Sin dependencias externas pesadas**

```bash
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
```

**Incluye:**
- 🚨 Manejo de errores estructurado (`errors`)
- 🏷️ Types compartidos: UUID, Enums (`types`)
- ✅ Validaciones comunes (`validator`)
- ⚙️ Configuration loaders (`config`)

**Dependencias:** Solo `google/uuid` (liviana)

---

### Módulo Logger

```bash
go get github.com/EduGoGroup/edugo-shared/logger@v2.0.5
```

**Incluye:**
- 📝 Interface Logger
- 📊 Implementación con Uber Zap
- 🎨 Formatos: JSON, Console con colores

**Dependencias:** `go.uber.org/zap`

---

### Módulo Auth (JWT)

```bash
go get github.com/EduGoGroup/edugo-shared/auth@v2.0.5
```

**Incluye:**
- 🔐 Generación de tokens JWT
- ✅ Validación de tokens
- 🔄 Refresh tokens
- 👥 Soporte para roles (admin, teacher, student, guardian)

**Dependencias:** `golang-jwt/jwt`, `google/uuid`, `common`

---

### Módulo RabbitMQ

```bash
go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@v2.0.5
```

**Incluye:**
- 📨 Publisher interface
- 📥 Consumer interface
- 🔌 Connection management
- ⚙️ Configuration helpers

**Dependencias:** `rabbitmq/amqp091-go`

---

### Módulo PostgreSQL

```bash
go get github.com/EduGoGroup/edugo-shared/database/postgres@v2.0.5
```

**Incluye:**
- 🗄️ Connection pooling
- 🔒 Transaction support
- 🏥 Health checks
- ⚙️ Configuration utilities

**Dependencias:** `lib/pq`

---

### Módulo MongoDB

```bash
go get github.com/EduGoGroup/edugo-shared/database/mongodb@v2.0.5
```

**Incluye:**
- 🗄️ Client configuration
- 🔄 Replica set support
- ☁️ MongoDB Atlas support
- ⚙️ Connection pooling

**Dependencias:** `mongo-driver`

---

## 🚀 Quick Start

### Ejemplo 1: Solo Errores y Validación (Ultra Liviano)

```bash
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
```

```go
import (
    "github.com/EduGoGroup/edugo-shared/common/errors"
    "github.com/EduGoGroup/edugo-shared/common/validator"
)

// Manejo de errores
err := errors.NewValidationError("email inválido")
err.WithField("email", userEmail)

// Validación
v := validator.New()
v.Email(email, "email")
v.Required(name, "name")

if v.HasErrors() {
    return v.GetError()
}
```

**Resultado:** `go.mod` con CERO dependencias externas pesadas ✅

---

### Ejemplo 2: Autenticación JWT

```bash
go get github.com/EduGoGroup/edugo-shared/auth@v2.0.5
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
```

```go
import (
    "github.com/EduGoGroup/edugo-shared/auth"
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// Crear JWT Manager
jwtManager := auth.NewJWTManager("secret-key", "edugo-api")

// Generar token
token, err := jwtManager.GenerateToken(
    userID,
    email,
    enum.SystemRoleTeacher,
    24*time.Hour,
)

// Validar token
claims, err := jwtManager.ValidateToken(token)
```

**Resultado:** Solo 3 dependencias (jwt, uuid, common) ✅

---

### Ejemplo 3: API Completa con Postgres y Logger

```bash
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
go get github.com/EduGoGroup/edugo-shared/logger@v2.0.5
go get github.com/EduGoGroup/edugo-shared/auth@v2.0.5
go get github.com/EduGoGroup/edugo-shared/database/postgres@v2.0.5
```

```go
import (
    "github.com/EduGoGroup/edugo-shared/common/errors"
    "github.com/EduGoGroup/edugo-shared/logger"
    "github.com/EduGoGroup/edugo-shared/auth"
    "github.com/EduGoGroup/edugo-shared/database/postgres"
)

// Logger
log := logger.NewZapLogger("info", "json")

// Database
db, err := postgres.Connect(postgres.Config{
    Host:     "localhost",
    Port:     5432,
    Database: "edugo",
    User:     "user",
    Password: "pass",
})

// Auth
jwtManager := auth.NewJWTManager(secretKey, "api")
```

---

## 📚 Documentación por Módulo

### Common

- **Errors**: Códigos de error estandarizados con HTTP status codes
- **Types**: UUID wrapper, Enums (Role, Status, Events, Assessment)
- **Validator**: Email, UUID, URL, nombres, rangos numéricos
- **Config**: Helpers para leer variables de entorno

### Logger

- **Niveles**: Debug, Info, Warn, Error, Fatal
- **Formatos**: JSON (producción), Console con colores (desarrollo)
- **Features**: Structured logging, context fields, caller info

### Auth

- **JWT Manager**: Generación y validación de tokens
- **Claims**: UserID, Email, Role, timestamps estándar
- **Refresh**: Soporte para refresh tokens
- **Roles**: Admin, Teacher, Student, Guardian

### Messaging/Rabbit

- **Connection**: Connection pooling y reconexión automática
- **Publisher**: Publicar mensajes con routing keys
- **Consumer**: Consumir mensajes con prefetch configurable
- **Config**: Exchange y Queue declaration helpers

### Database/Postgres

- **Connection**: Pool con configuración avanzada
- **Transactions**: Begin, Commit, Rollback helpers
- **Health**: Health check endpoint support
- **SSL**: Soporte para SSL modes (disable, require, verify-ca, verify-full)

### Database/MongoDB

- **Client**: Configuración de cliente MongoDB
- **Replica Sets**: Soporte completo
- **Atlas**: Compatible con `mongodb+srv://`
- **Pooling**: Connection pool configurable

---

## 🔄 Migración desde v2.0.1

Ver [UPGRADE_GUIDE.md](UPGRADE_GUIDE.md) para instrucciones detalladas.

### Resumen Rápido

**ANTES (v2.0.1):**
```go
import "github.com/EduGoGroup/edugo-shared/v2/pkg/errors"
import "github.com/EduGoGroup/edugo-shared/v2/pkg/auth"
```

```bash
go get github.com/EduGoGroup/edugo-shared/v2@v2.0.1
# Descarga: RabbitMQ, JWT, Zap, Postgres, Mongo (TODO)
```

**DESPUÉS (v2.0.5):**
```go
import "github.com/EduGoGroup/edugo-shared/common/errors"
import "github.com/EduGoGroup/edugo-shared/auth"
```

```bash
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
go get github.com/EduGoGroup/edugo-shared/auth@v2.0.5
# Descarga: Solo JWT + UUID + common (selectivo) ✅
```

---

## 🛠️ Desarrollo

### Comandos Make Disponibles

```bash
make help                # Ver todos los comandos
make test-all-modules    # Tests de todos los módulos
make build-all-modules   # Build de todos los módulos
make lint-all-modules    # Linter en todos los módulos
make coverage-all-modules # Coverage de todos los módulos
```

### Tests por Módulo

```bash
cd common && go test ./...
cd logger && go test ./...
cd auth && go test ./...
cd messaging/rabbit && go test ./...
cd database/postgres && go test ./...
cd database/mongodb && go test ./...
```

---

## 📊 Comparación de Dependencias

| Caso de Uso | v2.0.1 (Monolítico) | v2.0.5 (Modular) | Ahorro |
|-------------|---------------------|------------------|--------|
| Solo errores y types | 15+ deps | 1 dep (uuid) | ~93% |
| Errors + Auth | 15+ deps | 3 deps | ~80% |
| API completa (Postgres + Auth + Logger) | 15+ deps | ~8 deps | ~47% |
| API + RabbitMQ + Mongo | 15+ deps | ~12 deps | ~20% |

---

## 📋 Versionamiento

Este proyecto sigue [Semantic Versioning](https://semver.org/):

- **MAJOR**: Cambios incompatibles en la API
- **MINOR**: Nueva funcionalidad compatible
- **PATCH**: Corrección de bugs

### Historial de Versiones

- **v2.0.5** (2025-10-31): Arquitectura modular completa con 6 módulos independientes
- **v2.0.0** (2025-10-31): Separación inicial de databases en sub-módulos
- **v1.0.0** (2025-10-31): Primera versión estable monolítica

Ver [CHANGELOG.md](CHANGELOG.md) para detalles completos.

---

## 🤝 Contribuir

1. Asegurarse que el cambio sea **realmente compartido** (usado por 2+ proyectos)
2. Documentar públicamente con comentarios Go
3. Agregar tests unitarios (coverage mínimo 80%)
4. Usar interfaces para flexibilidad
5. Mantener módulos independientes

### Agregar Nueva Funcionalidad

```bash
# 1. Determinar módulo correcto
# 2. Implementar con tests
cd <module> && go test ./...

# 3. Verificar que no rompe otros módulos
make test-all-modules

# 4. Commit siguiendo conventional commits
git commit -m "feat(module): descripción"
```

---

## 📄 License

MIT License - EduGo Project

---

## 📞 Soporte

- **Issues**: [GitHub Issues](https://github.com/EduGoGroup/edugo-shared/issues)
- **Docs**: [pkg.go.dev](https://pkg.go.dev/github.com/EduGoGroup/edugo-shared)
- **Changelog**: [CHANGELOG.md](CHANGELOG.md)
- **Migration Guide**: [UPGRADE_GUIDE.md](UPGRADE_GUIDE.md)
