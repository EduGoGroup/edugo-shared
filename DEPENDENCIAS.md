# Dependencias del Módulo Shared

**Versión:** 0.1.0
**Última actualización:** 30 de Octubre, 2025
**Módulo Go:** `github.com/edugo/shared`

---

## 📦 Paquetes Disponibles

El módulo `shared` proporciona 8 paquetes reutilizables para los microservicios de EduGo:

### 1. `pkg/auth` - Autenticación JWT

**Descripción:** Manejo de tokens JWT para autenticación de usuarios y administradores.

**Archivos:**
- `jwt.go` - JWTManager para generación y validación de tokens

**Funcionalidad principal:**
- `NewJWTManager(secret, issuer string) *JWTManager` - Crear gestor de tokens
- `GenerateToken(userID uuid.UUID, role enum.Role) (string, error)` - Generar token JWT
- `ValidateToken(tokenString string) (*Claims, error)` - Validar y extraer claims

**Dependencias externas:**
- `github.com/golang-jwt/jwt/v5 v5.3.0`

**Usado por:**
- ✅ api-mobile
- ✅ api-administracion
- ❌ worker (no requiere autenticación)

---

### 2. `pkg/config` - Configuración

**Descripción:** Carga de variables de entorno y configuración de aplicación.

**Archivos:**
- `env.go` - Helpers para cargar variables de entorno

**Funcionalidad principal:**
- Loaders de variables de entorno con valores por defecto
- Validación de configuración requerida

**Dependencias externas:**
- Ninguna (solo standard library)

**Usado por:**
- ✅ api-mobile
- ✅ api-administracion
- ✅ worker

---

### 3. `pkg/database` - Conexiones a Bases de Datos

**Descripción:** Gestión de conexiones a PostgreSQL y MongoDB con pools de conexiones.

#### 3.1 `pkg/database/postgres`

**Archivos:**
- `config.go` - Configuración de conexión a PostgreSQL
- `connection.go` - Gestión del pool de conexiones
- `transaction.go` - Helpers para transacciones

**Funcionalidad principal:**
- `NewPostgresConnection(config Config) (*sql.DB, error)` - Crear conexión con pool
- `NewTransaction(db *sql.DB) (*sql.Tx, error)` - Iniciar transacción
- Health checks y reconexión automática

**Dependencias externas:**
- `github.com/lib/pq v1.10.9` - Driver PostgreSQL

**Usado por:**
- ✅ api-mobile (almacena usuarios, materiales, progreso)
- ✅ api-administracion (almacena usuarios, escuelas, unidades, asignaturas)
- ✅ worker (actualiza estados de procesamiento)

#### 3.2 `pkg/database/mongodb`

**Archivos:**
- `config.go` - Configuración de conexión a MongoDB
- `connection.go` - Gestión de cliente MongoDB

**Funcionalidad principal:**
- `NewMongoConnection(config Config) (*mongo.Client, error)` - Crear cliente
- `GetDatabase(client *mongo.Client, dbName string) *mongo.Database` - Obtener BD
- Reconexión automática y timeouts

**Dependencias externas:**
- `go.mongodb.org/mongo-driver v1.17.6`

**Usado por:**
- ✅ api-mobile (almacena resúmenes y evaluaciones procesadas)
- ❌ api-administracion (no usa MongoDB)
- ✅ worker (guarda resultados de procesamiento de PDFs)

---

### 4. `pkg/errors` - Manejo de Errores Personalizado

**Descripción:** Errores tipados para respuestas HTTP consistentes.

**Archivos:**
- `errors.go` - Tipos de errores personalizados

**Tipos de errores:**
- `NotFoundError` - Recurso no encontrado (404)
- `ValidationError` - Error de validación (400)
- `UnauthorizedError` - No autenticado (401)
- `ForbiddenError` - Sin permisos (403)
- `InternalError` - Error interno del servidor (500)
- `ConflictError` - Conflicto de recursos (409)

**Dependencias externas:**
- Ninguna

**Usado por:**
- ✅ api-mobile (manejo de errores HTTP)
- ✅ api-administracion (manejo de errores HTTP)
- ✅ worker (logging de errores)

---

### 5. `pkg/logger` - Sistema de Logging

**Descripción:** Logging estructurado con Zap para todos los servicios.

**Archivos:**
- `logger.go` - Interfaz Logger
- `zap_logger.go` - Implementación con Uber Zap

**Funcionalidad principal:**
- `NewZapLogger(level, format string) (Logger, error)` - Crear logger
- Niveles: Debug, Info, Warn, Error, Fatal
- Formatos: JSON (producción), Text (desarrollo)
- Logging estructurado con campos adicionales

**Dependencias externas:**
- `go.uber.org/zap v1.27.0`

**Usado por:**
- ✅ api-mobile (todos los logs)
- ✅ api-administracion (todos los logs)
- ✅ worker (todos los logs)

---

### 6. `pkg/messaging` - Message Queue (RabbitMQ)

**Descripción:** Cliente RabbitMQ para publicación y consumo de mensajes.

**Archivos:**
- `config.go` - Configuración de RabbitMQ
- `connection.go` - Gestión de conexión AMQP
- `publisher.go` - Publicador de eventos
- `consumer.go` - Consumidor de eventos

**Funcionalidad principal:**
- `NewRabbitMQPublisher(config Config) (*Publisher, error)` - Crear publicador
- `Publish(event Event) error` - Publicar evento a queue
- `NewRabbitMQConsumer(config Config) (*Consumer, error)` - Crear consumidor
- `Consume(queueName string, handler func(Event)) error` - Consumir eventos
- Reconexión automática
- Dead Letter Queue (DLQ) para mensajes fallidos

**Dependencias externas:**
- `github.com/rabbitmq/amqp091-go v1.10.0`

**Usado por:**
- ✅ api-mobile (publica eventos: MaterialUploaded, AssessmentAttempt, etc.)
- ❌ api-administracion (no usa messaging actualmente)
- ✅ worker (consume eventos de la queue)

---

### 7. `pkg/types` - Tipos Compartidos

**Descripción:** Tipos de datos y enumeraciones compartidas entre servicios.

**Archivos:**
- `uuid.go` - Tipo UUID personalizado con serialización JSON

#### 7.1 `pkg/types/enum`

**Enumeraciones disponibles:**

**`role.go` - Roles de usuarios:**
- `RoleTeacher` - Profesor
- `RoleStudent` - Estudiante
- `RoleGuardian` - Tutor/Padre
- `RoleAdmin` - Administrador

**`status.go` - Estados generales:**
- `StatusPublished` - Publicado
- `StatusDraft` - Borrador
- `StatusArchived` - Archivado
- `StatusDeleted` - Eliminado

**`assessment.go` - Estados de evaluaciones:**
- `AssessmentPending` - Pendiente
- `AssessmentProcessing` - Procesando
- `AssessmentCompleted` - Completado
- `AssessmentFailed` - Fallido

**`event.go` - Tipos de eventos RabbitMQ:**
- `EventMaterialUploaded` - Material subido
- `EventAssessmentAttempt` - Intento de evaluación
- `EventMaterialDeleted` - Material eliminado
- `EventMaterialReprocess` - Reprocesar material
- `EventStudentEnrolled` - Estudiante inscrito

**Dependencias externas:**
- `github.com/google/uuid v1.6.0`

**Usado por:**
- ✅ api-mobile (todos los tipos)
- ✅ api-administracion (todos los tipos)
- ✅ worker (eventos y estados)

---

### 8. `pkg/validator` - Validaciones

**Descripción:** Funciones de validación comunes para entrada de datos.

**Archivos:**
- `validator.go` - Validadores de datos

**Funcionalidad principal:**
- `ValidateEmail(email string) error` - Validar formato de email
- `ValidateUUID(id string) error` - Validar formato UUID
- `ValidateRequired(value interface{}) error` - Validar campo requerido
- `ValidateStringLength(value string, min, max int) error` - Validar longitud

**Dependencias externas:**
- Ninguna (usa regex y standard library)

**Usado por:**
- ✅ api-mobile (validación de requests)
- ✅ api-administracion (validación de requests)
- ❌ worker (no valida input de usuarios)

---

## 📊 Matriz de Dependencias por Servicio

| Paquete | api-mobile | api-administracion | worker |
|---------|------------|-------------------|--------|
| `pkg/auth` | ✅ Sí | ✅ Sí | ❌ No |
| `pkg/config` | ✅ Sí | ✅ Sí | ✅ Sí |
| `pkg/database/postgres` | ✅ Sí | ✅ Sí | ✅ Sí |
| `pkg/database/mongodb` | ✅ Sí | ❌ No | ✅ Sí |
| `pkg/errors` | ✅ Sí | ✅ Sí | ✅ Sí |
| `pkg/logger` | ✅ Sí | ✅ Sí | ✅ Sí |
| `pkg/messaging` | ✅ Sí | ❌ No | ✅ Sí |
| `pkg/types` | ✅ Sí | ✅ Sí | ✅ Sí |
| `pkg/validator` | ✅ Sí | ✅ Sí | ❌ No |

---

## 🔗 Dependencias Externas de `shared`

Todas las dependencias externas definidas en `shared/go.mod`:

```go
require (
    github.com/golang-jwt/jwt/v5 v5.3.0         // JWT para autenticación
    github.com/google/uuid v1.6.0               // Manejo de UUIDs
    github.com/lib/pq v1.10.9                   // Driver PostgreSQL
    github.com/rabbitmq/amqp091-go v1.10.0      // Cliente RabbitMQ
    go.mongodb.org/mongo-driver v1.17.6         // Driver MongoDB
    go.uber.org/zap v1.27.0                     // Logger estructurado
)
```

**Nota:** Todas estas dependencias son estables y ampliamente usadas en producción.

---

## 📝 Uso de `shared` en los Servicios

### Patrón de Import Actual (Monorepo)

Cada servicio tiene en su `go.mod`:

```go
module github.com/edugo/api-mobile

require (
    github.com/edugo/shared v0.0.0-00010101000000-000000000000
)

replace github.com/edugo/shared => ../../shared
```

### Patrón de Import Futuro (Multi-repo)

Después de la separación, los servicios usarán:

```go
module github.com/edugo/api-mobile

require (
    github.com/edugo/edugo-shared v0.1.0
)

// ¡Ya no hay replace!
```

---

## 🔄 Flujo de Datos entre Servicios

```
┌─────────────────┐
│   api-mobile    │
│                 │
│ Usa:            │
│ - auth (JWT)    │
│ - database/     │
│   postgres      │
│ - messaging     │
│   (Publisher)   │
│ - logger        │
│ - types/enum    │
└────────┬────────┘
         │
         │ Publica eventos
         ↓
┌─────────────────┐
│    RabbitMQ     │
└────────┬────────┘
         │
         │ Consume eventos
         ↓
┌─────────────────┐
│     worker      │
│                 │
│ Usa:            │
│ - database/     │
│   postgres,     │
│   mongodb       │
│ - messaging     │
│   (Consumer)    │
│ - logger        │
│ - types/enum    │
└─────────────────┘
```

---

## 🚀 Versionamiento

El módulo `shared` seguirá [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** (1.x.x): Cambios incompatibles en la API
- **MINOR** (x.1.x): Nueva funcionalidad compatible
- **PATCH** (x.x.1): Bug fixes compatibles

**Versión actual:** `v0.1.0` (pre-release antes de separación)

**Próxima versión:** `v0.1.0` (primera versión estable después de separación)

---

## 📚 Recursos Adicionales

- **README principal:** [shared/README.md](README.md)
- **Changelog:** [shared/CHANGELOG.md](CHANGELOG.md)
- **Guía de uso:** `/GUIA_USO_SHARED.md`

---

**Última actualización:** 30 de Octubre, 2025
**Mantenedor:** Equipo EduGo
