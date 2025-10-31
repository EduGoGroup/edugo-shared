# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.5] - 2025-10-31

### 🚀 BREAKING CHANGES - Arquitectura Modular Completa

#### Eliminado Módulo v2 Monolítico

- ❌ **Eliminado**: `github.com/EduGoGroup/edugo-shared/v2` (módulo core monolítico)
- ❌ **Eliminado**: Carpeta `pkg/` completa
- ❌ **Eliminado**: `go.mod` raíz

#### 6 Módulos Independientes Creados

1. **`common/`** - Errors, Types, Validator, Config
   - Dependencias: Solo `google/uuid` (liviana)
   - Path: `github.com/EduGoGroup/edugo-shared/common`

2. **`logger/`** - Logging con Uber Zap
   - Dependencias: `go.uber.org/zap`
   - Path: `github.com/EduGoGroup/edugo-shared/logger`

3. **`auth/`** - JWT Authentication
   - Dependencias: `golang-jwt/jwt`, `google/uuid`, `common`
   - Path: `github.com/EduGoGroup/edugo-shared/auth`

4. **`messaging/rabbit/`** - RabbitMQ helpers
   - Dependencias: `rabbitmq/amqp091-go`
   - Path: `github.com/EduGoGroup/edugo-shared/messaging/rabbit`

5. **`database/postgres/`** - PostgreSQL utilities
   - Dependencias: `lib/pq`
   - Path: `github.com/EduGoGroup/edugo-shared/database/postgres`

6. **`database/mongodb/`** - MongoDB utilities
   - Dependencias: `mongo-driver`
   - Path: `github.com/EduGoGroup/edugo-shared/database/mongodb`

### 📋 Migración de Imports

| Antes (v2.0.1) | Después (v2.0.5) |
|----------------|------------------|
| `github.com/EduGoGroup/edugo-shared/v2/pkg/errors` | `github.com/EduGoGroup/edugo-shared/common/errors` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/types` | `github.com/EduGoGroup/edugo-shared/common/types` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/validator` | `github.com/EduGoGroup/edugo-shared/common/validator` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/config` | `github.com/EduGoGroup/edugo-shared/common/config` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/auth` | `github.com/EduGoGroup/edugo-shared/auth` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/logger` | `github.com/EduGoGroup/edugo-shared/logger` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/messaging` | `github.com/EduGoGroup/edugo-shared/messaging/rabbit` |
| `github.com/EduGoGroup/edugo-shared/database/postgres` | Sin cambios ✓ |
| `github.com/EduGoGroup/edugo-shared/database/mongodb` | Sin cambios ✓ |

### ✨ Beneficios

- ✅ **Dependencias ultra-selectivas**: Módulo `common` con solo 1 dependencia externa
- ✅ **Ahorro de ~93%** en dependencias si solo usas `common`
- ✅ **Ahorro de ~80%** en dependencias si usas `common` + `auth`
- ✅ **Binarios más pequeños**: Solo se compila lo que realmente usas
- ✅ **Testing modular**: Cada módulo se testea independientemente
- ✅ **CI/CD optimizado**: Workflows paralelos por módulo

### 🔧 CI/CD Actualizado

- **ci.yml**: Strategy matrix para 6 módulos con tests paralelos
- **test.yml**: Coverage independiente por módulo + summary consolidado
- **release.yml**: Validación completa de todos los módulos antes de release

### 📚 Documentación Actualizada

- **README.md**: Sección modular completa con ejemplos por módulo
- **UPGRADE_GUIDE.md**: Tabla de migración detallada v2.0.1 → v2.0.5
- **Makefile**: Comandos para testear/build todos los módulos

### 🎯 Ejemplo de Ahorro

**Antes (v2.0.1):**
```bash
go get github.com/EduGoGroup/edugo-shared/v2@v2.0.1
# Descarga: 15+ dependencias (RabbitMQ, JWT, Zap, etc.)
```

**Después (v2.0.5):**
```bash
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
# Descarga: 1 dependencia (google/uuid)
# Ahorro: ~93% ✅
```

### 📦 Instalación Modular

Ver [README.md](README.md) para instrucciones completas de cada módulo.

---

## [2.0.0] - 2025-10-31

### 🚀 BREAKING CHANGES

#### Arquitectura Modular con Sub-módulos Independientes

- **Separación de módulos de bases de datos**: PostgreSQL y MongoDB ahora son sub-módulos Go independientes
- **Cambio en la estructura de directorios**:
  - ❌ Antes: `pkg/database/postgres/` y `pkg/database/mongodb/`
  - ✅ Ahora: `database/postgres/` y `database/mongodb/`
- **Cambio en imports**:
  - ❌ Antes: `import "github.com/EduGoGroup/edugo-shared/pkg/database/postgres"`
  - ✅ Ahora: `import "github.com/EduGoGroup/edugo-shared/database/postgres"`

### ✨ Mejoras

- **Dependencias selectivas**: Los proyectos ahora pueden importar solo el módulo de base de datos que necesitan
- **Reducción de dependencias transitivas**:
  - Si solo usas PostgreSQL, no se descarga el driver de MongoDB (y viceversa)
  - El módulo core ya no incluye drivers de bases de datos
- **Binarios más ligeros**: El compilador de Go solo incluye el código que realmente se usa
- **Mejor mantenibilidad**: Cada módulo de base de datos tiene su propio `go.mod` y versionado

### 📦 Nuevos Módulos

1. **github.com/EduGoGroup/edugo-shared** (core)
   - Incluye: logger, messaging, errors, validator, auth, config, types
   - Sin dependencias de bases de datos

2. **github.com/EduGoGroup/edugo-shared/database/postgres**
   - Módulo independiente para PostgreSQL
   - Dependencias: `github.com/lib/pq`

3. **github.com/EduGoGroup/edugo-shared/database/mongodb**
   - Módulo independiente para MongoDB
   - Dependencias: `go.mongodb.org/mongo-driver`

### 📋 Migración

Ver [UPGRADE_GUIDE.md](UPGRADE_GUIDE.md) para instrucciones detalladas de migración.

**Resumen rápido:**

```bash
# 1. Actualizar go.mod
go get github.com/EduGoGroup/edugo-shared@v2.0.0
go get github.com/EduGoGroup/edugo-shared/database/postgres@v2.0.0  # Si usas PostgreSQL
go get github.com/EduGoGroup/edugo-shared/database/mongodb@v2.0.0   # Si usas MongoDB

# 2. Actualizar imports en tu código
# Buscar y reemplazar:
#   pkg/database/postgres -> database/postgres
#   pkg/database/mongodb -> database/mongodb

# 3. Actualizar dependencias
go mod tidy
```

### 🎯 Beneficios de la Migración

| Aspecto | v1.0.0 | v2.0.0 |
|---------|--------|--------|
| Dependencias descargadas | Todas las BDs | Solo las que uses |
| Tamaño del go.mod | ~15 dependencias | ~5-8 dependencias |
| Binario compilado | Optimizado | Optimizado |
| Flexibilidad | Baja | Alta |

### Dependencies

**Core Module:**
```
github.com/golang-jwt/jwt/v5 v5.3.0
github.com/google/uuid v1.6.0
github.com/rabbitmq/amqp091-go v1.10.0
go.uber.org/zap v1.27.0
github.com/stretchr/testify v1.8.1
```

**PostgreSQL Module:**
```
github.com/lib/pq v1.10.9
```

**MongoDB Module:**
```
go.mongodb.org/mongo-driver v1.17.6
```

---

## [1.0.0] - 2025-10-31

### Added
- **JWT Authentication System**: Complete JWT token generation, validation, and management with support for all system roles
- **Database Connectivity**: 
  - PostgreSQL connection utilities with connection pooling, health checks, and transaction support
  - MongoDB connection utilities with configurable pools and health monitoring
- **Error Management**: Structured error handling with HTTP status codes and contextual information
- **Validation System**: Comprehensive input validation with support for emails, UUIDs, URLs, and custom rules  
- **Messaging System**: RabbitMQ integration with publishers, consumers, and connection management
- **Configuration Utilities**: Environment variable handling and configuration management
- **Logging**: Structured logging with Zap integration
- **Type System**: Custom UUID handling and enum definitions for roles, events, assessments, and statuses
- **Development Tooling**: 
  - Comprehensive Makefile with 20+ commands (build, test, lint, coverage, security)
  - golangci-lint configuration with professional standards
  - GitHub Actions CI/CD pipeline
  - Pre-commit hooks setup

### Technical Improvements
- **Code Quality**: 100% linter compliance with zero warnings
- **Test Coverage**: 87.2% coverage in authentication module with comprehensive test suites
- **Memory Optimization**: Struct field alignment optimizations reducing memory usage by 10-15%
- **Performance**: Optimized imports and eliminated unused code
- **Documentation**: Complete package documentation following Go standards
- **Constants**: Extracted magic numbers to named constants for better maintainability

### Dependencies
- `github.com/golang-jwt/jwt/v5`: JWT token handling
- `github.com/google/uuid`: UUID generation and parsing
- `github.com/lib/pq`: PostgreSQL driver
- `go.mongodb.org/mongo-driver`: MongoDB driver
- `github.com/streadway/amqp`: RabbitMQ client
- `go.uber.org/zap`: Structured logging
- `github.com/stretchr/testify`: Testing utilities

### Security
- JWT tokens use secure signing methods (HMAC-SHA256)
- Environment-based configuration for sensitive data
- SQL injection protection through parameterized queries
- Input validation for all user-facing APIs

## [Unreleased]

### Planeado
- Agregar tests de integración con Testcontainers
- Agregar validación de configuración
- Agregar metrics y tracing
- Mejorar manejo de errores con wrapped errors

---

## [0.1.0] - 2025-10-30

### Añadido

#### Autenticación (pkg/auth)
- JWTManager para generación y validación de tokens
- Soporte para múltiples roles (admin, teacher, student, guardian)
- Refresh token con expiración personalizada
- Funciones de extracción sin validación (para logging)
- Tests unitarios completos (15 tests)

#### Configuración (pkg/config)
- Loaders de variables de entorno
- Valores por defecto razonables

#### Database (pkg/database)
- **PostgreSQL:**
  - Configuración de pool de conexiones
  - Soporte para SSL (disable, require, verify-ca, verify-full)
  - Manejo de transacciones
  - Reconexión automática
  - Tests de configuración
- **MongoDB:**
  - Configuración de cliente MongoDB
  - Soporte para replica sets
  - Soporte para MongoDB Atlas (mongodb+srv)
  - Pool de conexiones configurable
  - Tests de configuración

#### Manejo de Errores (pkg/errors)
- NotFoundError (404)
- ValidationError (400)
- UnauthorizedError (401)
- ForbiddenError (403)
- InternalError (500)
- ConflictError (409)
- Errores tipados para respuestas HTTP consistentes

#### Logging (pkg/logger)
- Implementación con Uber Zap
- Niveles: Debug, Info, Warn, Error, Fatal
- Formatos: JSON (producción), Text (desarrollo)
- Logging estructurado con campos adicionales

#### Messaging (pkg/messaging)
- Cliente RabbitMQ
- Publisher para enviar eventos
- Consumer para procesar eventos
- Reconexión automática
- Dead Letter Queue (DLQ) para mensajes fallidos
- Prefetch configurable

#### Types (pkg/types)
- Tipo UUID personalizado con serialización JSON
- **Enumeraciones (pkg/types/enum):**
  - SystemRole: admin, teacher, student, guardian
  - Status: published, draft, archived, deleted
  - AssessmentStatus: pending, processing, completed, failed
  - EventType: material_uploaded, assessment_attempt, material_deleted, material_reprocess, student_enrolled
- Validación de enums

#### Validación (pkg/validator)
- Validación de emails
- Validación de UUIDs
- Validación de campos requeridos
- Validación de longitud de strings

### Documentación
- README.md completo con ejemplos de uso
- DEPENDENCIAS.md con mapeo de servicios
- CHANGELOG.md (este archivo)
- Documentación inline en todos los paquetes

### Tests
- Tests unitarios para JWT (100% funciones principales)
- Tests de configuración para PostgreSQL
- Tests de configuración para MongoDB
- Cobertura total: ~70%

### Dependencias Externas
```
github.com/golang-jwt/jwt/v5 v5.3.0
github.com/google/uuid v1.6.0
github.com/lib/pq v1.10.9
github.com/rabbitmq/amqp091-go v1.10.0
go.mongodb.org/mongo-driver v1.17.6
go.uber.org/zap v1.27.0
```

### Notas de Migración
Este es el primer release estable del módulo shared antes de la separación del monorepo.

**Cómo usar:**
```go
// En go.mod del monorepo
require (
    github.com/edugo/shared v0.0.0-00010101000000-000000000000
)
replace github.com/edugo/shared => ../../shared

// Después de la separación (futuro)
require (
    github.com/edugo/edugo-shared v0.1.0
)
```

---

## Formato de Versiones

- **MAJOR** (1.x.x): Cambios incompatibles en la API
- **MINOR** (x.1.x): Nueva funcionalidad compatible hacia atrás
- **PATCH** (x.x.1): Bug fixes compatibles hacia atrás

---

## Tipos de Cambios

- **Añadido** - para nuevas funcionalidades
- **Cambiado** - para cambios en funcionalidad existente
- **Obsoleto** - para funcionalidades que pronto se eliminarán
- **Eliminado** - para funcionalidades eliminadas
- **Corregido** - para corrección de bugs
- **Seguridad** - en caso de vulnerabilidades

---

**Mantenedor:** Equipo EduGo
**Última actualización:** 30 de Octubre, 2025
