# Changelog

Todos los cambios notables a este proyecto serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

---

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
