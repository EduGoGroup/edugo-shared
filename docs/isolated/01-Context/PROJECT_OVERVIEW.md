# PROJECT OVERVIEW - SHARED

## Información General

**Proyecto:** EduGo Shared Library  
**Tipo:** Librería Go (módulo reutilizable)  
**Lenguaje:** Go 1.21+  
**Especificación de Origen:** spec-04-shared  
**Estado:** En Desarrollo (Sprint 1/4)

---

## Propósito del Proyecto

SHARED es una librería Go centralizada que proporciona módulos reutilizables para todos los otros proyectos de EduGo (api-mobile, api-admin, worker). Implementa patrones comunes y evita duplicación de código.

### Responsabilidades Principales
- Logger centralizado (structured logging)
- Database connections (PostgreSQL, MongoDB)
- Authentication (JWT validation)
- Messaging (RabbitMQ client)
- Context management (timeouts, user info)
- Error handling estandarizado
- Configuration management
- Health checks

---

## Estructura de Módulos

```
github.com/EduGoGroup/edugo-shared/
├── logger/                  # Structured logging
├── database/                # DB connections
│   ├── postgres/
│   └── mongo/
├── auth/                    # JWT, authentication
├── messaging/               # RabbitMQ client
├── models/                  # Shared data structures
├── context/                 # Context utilities
├── errors/                  # Error types
├── health/                  # Health checks
├── config/                  # Configuration
└── VERSION                  # Versión semántica
```

---

## Módulos Detallados

### 1. Logger Module

**Propósito:** Logging estructurado centralizado

```go
import "github.com/EduGoGroup/edugo-shared/logger"

// Inicializar
logger.Init(logger.Config{
    Level:  "info",
    Format: "json",
    Output: "stdout",
})

// Usar
logger.Info("Evento importante", map[string]interface{}{
    "user_id": 42,
    "action": "create_evaluation",
})

logger.Error("Error ocurrido", map[string]interface{}{
    "error": err.Error(),
    "stack": debug.Stack(),
})
```

**Características:**
- JSON structured output
- Niveles: debug, info, warn, error
- Request tracing (request_id)
- User context (user_id, school_id)
- Stack traces en errores

---

### 2. Database Module

#### PostgreSQL Connection

```go
import "github.com/EduGoGroup/edugo-shared/database"

// Inicializar
database.Init(database.Config{
    Host:     "localhost",
    Port:     5432,
    User:     "edugo_user",
    Password: "password",
    Database: "edugo_mobile",
})

// Usar
db := database.GetDB()
var user User
db.First(&user, id)
```

**Características:**
- Connection pooling
- GORM integration
- Automatic migrations support
- Transaction support
- Prepared statements

#### MongoDB Connection

```go
import "github.com/EduGoGroup/edugo-shared/database"

// Inicializar MongoDB
mongo := database.InitMongo(
    "mongodb://localhost:27017",
    "edugo_assessments",
)

// Usar
client := mongo.Client()
collection := client.Database("edugo").Collection("results")
```

---

### 3. Auth Module

**Propósito:** JWT validation y context extraction

```go
import "github.com/EduGoGroup/edugo-shared/auth"

// En middleware Gin
func AuthMiddleware(c *gin.Context) {
    claims, err := auth.ValidateToken(c.GetHeader("Authorization"))
    if err != nil {
        c.JSON(401, gin.H{"error": "Unauthorized"})
        c.Abort()
        return
    }
    
    // Inyectar en contexto
    c.Set("user_id", claims.UserID)
    c.Set("school_id", claims.SchoolID)
    c.Set("roles", claims.Roles)
    c.Next()
}

// En handlers
func GetEvaluation(c *gin.Context) {
    userID := c.GetInt64("user_id")
    schoolID := c.GetInt64("school_id")
    // ... lógica
}
```

**Estructura de Claims:**
```go
type Claims struct {
    UserID   int64    `json:"user_id"`
    SchoolID int64    `json:"school_id"`
    Email    string   `json:"email"`
    Roles    []string `json:"roles"`
    ExpiresAt int64   `json:"exp"`
}
```

---

### 4. Messaging Module

**Propósito:** RabbitMQ client abstraction

```go
import "github.com/EduGoGroup/edugo-shared/messaging"

// Publicador
publisher := messaging.NewPublisher()
err := publisher.Publish(
    "assessment.requests",        // exchange
    "worker.assessment.requests", // routing key
    payload,                       // []byte
)

// Suscriptor
subscriber := messaging.NewSubscriber()
messages := subscriber.Subscribe(
    "assessment.responses",        // exchange
    "api-mobile.assessment.responses", // queue
)

for msg := range messages {
    // Procesar mensaje
    msg.Ack(false)  // Acknowledge
}
```

**Características:**
- Connection pooling
- Auto-reconnect
- Dead letter queue support
- Manual acknowledgement
- Retry logic

---

### 5. Models Module

**Propósito:** Estructuras de datos compartidas

```go
import "github.com/EduGoGroup/edugo-shared/models"

// User
type User struct {
    ID        int64  `json:"id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    SchoolID  int64  `json:"school_id"`
}

// School
type School struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    City string `json:"city"`
}

// AcademicUnit
type AcademicUnit struct {
    ID       int64  `json:"id"`
    ParentID *int64 `json:"parent_id"`
    Name     string `json:"name"`
    Type     string `json:"type"` // faculty, department, etc
}
```

---

### 6. Context Module

**Propósito:** Manejo de contexto y timeouts globales

```go
import "github.com/EduGoGroup/edugo-shared/context"

// Crear contexto con timeout
ctx := context.WithTimeout(
    context.Background(),
    30*time.Second,  // timeout global
)

// Usar en operaciones
ch := make(chan interface{})
go func() {
    result := expensiveOperation()
    ch <- result
}()

select {
case result := <-ch:
    return result
case <-ctx.Done():
    return errors.New("timeout exceeded")
}
```

---

### 7. Errors Module

**Propósito:** Error types estandarizados

```go
import "github.com/EduGoGroup/edugo-shared/errors"

// Crear error
err := errors.NewNotFound("Evaluation", "evaluation_id", id)

// Usar en handlers
if err != nil {
    httpStatus, response := errors.ToHTTP(err)
    c.JSON(httpStatus, response)
}

// Response:
{
  "error": "not_found",
  "message": "Evaluation with id 123 not found",
  "code": "EVALUATION_NOT_FOUND"
}
```

**Error Types:**
- NotFound (404)
- BadRequest (400)
- Unauthorized (401)
- Forbidden (403)
- InternalError (500)

---

### 8. Health Module

**Propósito:** Health checks y readiness

```go
import "github.com/EduGoGroup/edugo-shared/health"

// Registrar checkers
health.Register("postgres", &DatabaseChecker{})
health.Register("mongodb", &MongoChecker{})
health.Register("rabbitmq", &RabbitMQChecker{})

// Usar en Gin
router.GET("/health", health.Handler())
router.GET("/ready", health.ReadinessHandler())

// Response:
{
  "status": "ok",
  "checks": {
    "postgres": "ok",
    "mongodb": "ok",
    "rabbitmq": "ok"
  },
  "timestamp": "2025-11-15T10:30:00Z"
}
```

---

## Versionamiento

**Estrategia:** Semantic Versioning (MAJOR.MINOR.PATCH)

### Versión Actual: v1.3.0

```
v1.3.0
│  └─ PATCH: Bug fixes (compatible)
├─ MINOR: Nuevos features (backward compatible)
└─ MAJOR: Breaking changes (nueva versión)
```

### Cambios por Versión

| Versión | Cambios |
|---------|---------|
| v1.0.0 | Initial release |
| v1.1.0 | + Auth module |
| v1.2.0 | + Messaging module |
| v1.3.0 | + Health checks |
| v1.4.0 | PostgreSQL improvements (en desarrollo) |

---

## Dependencias de SHARED

```go
// Logger
github.com/sirupsen/logrus
github.com/fatih/color

// Database
gorm.io/gorm
gorm.io/driver/postgres
go.mongodb.org/mongo-driver

// Auth
github.com/golang-jwt/jwt/v5

// Messaging
github.com/streadway/amqp

// Configuration
github.com/spf13/viper
```

---

## Uso en Otros Proyectos

### En api-mobile/go.mod
```
require github.com/EduGoGroup/edugo-shared v1.3.0
```

### En api-mobile/main.go
```go
import (
    "github.com/EduGoGroup/edugo-shared/logger"
    "github.com/EduGoGroup/edugo-shared/database"
    "github.com/EduGoGroup/edugo-shared/auth"
)

func main() {
    // Usar módulos de SHARED
    logger.Init(...)
    database.Init(...)
    
    // Resto de la app
}
```

---

## Compilación y Release

### Compilación Local
```bash
go mod download
go mod tidy
go build ./...
```

### Tests
```bash
go test ./...
go test -cover ./...
```

### Release (desde rama dev)
```bash
# 1. Mergear cambios a dev
# 2. Tag versión
git tag v1.3.0

# 3. Push tag
git push origin v1.3.0

# 4. Otros proyectos actualizan go.mod
# En api-mobile:
go get -u github.com/EduGoGroup/edugo-shared@v1.3.0
```

---

## Sprint Planning (4 Sprints)

| Sprint | Funcionalidad | Duración |
|--------|---------------|----------|
| 1 | Setup + Logger + Database | 2 semanas |
| 2 | Auth + Messaging | 2 semanas |
| 3 | Models + Context + Errors | 2 semanas |
| 4 | Health + Optimizaciones | 2 semanas |

---

## Contacto y Referencias

- **Repositorio GitHub:** https://github.com/EduGoGroup/edugo-shared
- **Especificación Completa:** docs/ESTADO_PROYECTO.md (repo análisis)
- **Documentación Técnica:** Este directorio (01-Context/)
