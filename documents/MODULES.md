# Módulos del Sistema

## Visión General de Módulos

```
edugo-shared/
├── auth/              → Autenticación y seguridad
├── bootstrap/         → Inicialización de aplicaciones
├── common/            → Código compartido
│   ├── errors/        → Manejo de errores
│   ├── types/         → Tipos del dominio
│   └── validator/     → Validación de datos
├── config/            → Gestión de configuración
├── database/          → Conexiones a bases de datos
│   ├── mongodb/       → Cliente MongoDB
│   └── postgres/      → Cliente PostgreSQL
├── lifecycle/         → Gestión de ciclo de vida
├── logger/            → Logging estructurado
├── messaging/         → Sistema de mensajería
│   └── rabbit/        → Cliente RabbitMQ
├── middleware/        → Middlewares HTTP
│   └── gin/           → Middlewares para Gin
└── testing/           → Infraestructura de testing
    └── containers/    → Testcontainers
```

---

## 1. Auth Module

**Ubicación:** `auth/`  
**Import:** `github.com/EduGoGroup/edugo-shared/auth`

### Propósito
Proporciona autenticación basada en JWT y manejo seguro de contraseñas.

### Componentes

#### JWTManager
Gestiona la generación y validación de tokens JWT.

```go
// Crear manager
jwtManager := auth.NewJWTManager("secret-key", "edugo-issuer")

// Generar token
token, err := jwtManager.GenerateToken(
    "user-uuid",           // userID
    "user@edugo.com",      // email
    enum.SystemRoleStudent, // role
    time.Hour * 24,        // expiration
)

// Validar token
claims, err := jwtManager.ValidateToken(tokenString)
// claims.UserID, claims.Email, claims.Role

// Refresh token
newToken, err := jwtManager.RefreshToken(tokenString, time.Hour * 24)
```

#### Claims Structure
```go
type Claims struct {
    UserID string          `json:"user_id"`
    Email  string          `json:"email"`
    Role   enum.SystemRole `json:"role"`
    jwt.RegisteredClaims
}
```

#### Password Hashing
```go
// Hash password
hashedPassword, err := auth.HashPassword("plain-password")

// Verify password
valid := auth.VerifyPassword(hashedPassword, "plain-password")
```

#### Refresh Tokens
```go
// Generar refresh token
refreshToken := auth.GenerateRefreshToken()

// Validar formato
valid := auth.ValidateRefreshToken(refreshToken)
```

### Dependencias
- `github.com/golang-jwt/jwt/v5`
- `github.com/google/uuid`
- `golang.org/x/crypto/bcrypt`

---

## 2. Bootstrap Module

**Ubicación:** `bootstrap/`  
**Import:** `github.com/EduGoGroup/edugo-shared/bootstrap`

### Propósito
Inicialización centralizada de todos los recursos de infraestructura de una aplicación.

### Componentes Principales

#### Bootstrap Function
```go
resources, err := bootstrap.Bootstrap(
    ctx,              // context
    config,           // configuración
    factories,        // factories de recursos
    lifecycleManager, // gestor de ciclo de vida
    options...,       // opciones adicionales
)
```

#### Factories
```go
type Factories struct {
    Logger     LoggerFactory
    PostgreSQL PostgreSQLFactory
    MongoDB    MongoDBFactory
    RabbitMQ   RabbitMQFactory
    S3         S3Factory
}
```

#### Resources (Output)
```go
type Resources struct {
    Logger           *logrus.Logger
    PostgreSQL       *gorm.DB
    MongoDB          *mongo.Client
    MongoDatabase    *mongo.Database
    MessagePublisher MessagePublisher
    StorageClient    StorageClient
}
```

#### Options
```go
// Recursos requeridos (falla si no pueden inicializar)
bootstrap.WithRequiredResources("logger", "postgresql", "mongodb")

// Recursos opcionales (continúa si fallan)
bootstrap.WithOptionalResources("rabbitmq", "s3")

// Saltar health checks
bootstrap.WithSkipHealthCheck()

// Mock factories para testing
bootstrap.WithMockFactories(mockFactories)
```

### Secuencia de Inicialización
1. Logger (siempre primero)
2. PostgreSQL
3. MongoDB
4. RabbitMQ
5. S3
6. Health Checks

---

## 3. Common Module

**Ubicación:** `common/`  
**Import:** `github.com/EduGoGroup/edugo-shared/common/...`

### 3.1 Errors

```go
import "github.com/EduGoGroup/edugo-shared/common/errors"

// Crear error
err := errors.NewValidationError("email inválido")

// Con campos adicionales
err := errors.NewNotFoundError("user").
    WithField("user_id", "123").
    WithDetails("User was deleted")

// Wrap error interno
err := errors.NewDatabaseError("insert", originalErr)
```

#### Códigos de Error Disponibles
| Código | HTTP Status | Descripción |
|--------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Error de validación |
| `INVALID_INPUT` | 400 | Input malformado |
| `NOT_FOUND` | 404 | Recurso no encontrado |
| `ALREADY_EXISTS` | 409 | Recurso ya existe |
| `CONFLICT` | 409 | Conflicto de estado |
| `UNAUTHORIZED` | 401 | No autenticado |
| `FORBIDDEN` | 403 | Sin permisos |
| `INVALID_TOKEN` | 401 | Token inválido |
| `TOKEN_EXPIRED` | 401 | Token expirado |
| `BUSINESS_RULE_VIOLATION` | 422 | Regla de negocio |
| `INTERNAL_ERROR` | 500 | Error interno |
| `DATABASE_ERROR` | 500 | Error de BD |
| `EXTERNAL_SERVICE_ERROR` | 500 | Servicio externo |
| `TIMEOUT` | 408 | Timeout |
| `RATE_LIMIT_EXCEEDED` | 429 | Rate limit |

### 3.2 Types/Enum

```go
import "github.com/EduGoGroup/edugo-shared/common/types/enum"

// Roles del sistema
role := enum.SystemRoleStudent  // admin, teacher, student, guardian
if role.IsValid() { ... }

// Status de materiales
status := enum.MaterialStatusPublished  // draft, published, archived

// Status de progreso
progress := enum.ProgressStatusInProgress  // not_started, in_progress, completed

// Status de procesamiento
processing := enum.ProcessingStatusPending  // pending, processing, completed, failed
```

### 3.3 Types/UUID

```go
import "github.com/EduGoGroup/edugo-shared/common/types"

// Generar nuevo UUID
id := types.NewUUID()

// Validar UUID
valid := types.IsValidUUID("550e8400-e29b-41d4-a716-446655440000")
```

---

## 4. Config Module

**Ubicación:** `config/`  
**Import:** `github.com/EduGoGroup/edugo-shared/config`

### Propósito
Carga y validación de configuración desde archivos YAML.

### BaseConfig Structure
```go
type BaseConfig struct {
    Environment string          // local, dev, qa, prod
    ServiceName string          // Nombre del servicio
    Server      ServerConfig    // Configuración HTTP
    Database    DatabaseConfig  // PostgreSQL config
    MongoDB     MongoDBConfig   // MongoDB config
    Logger      LoggerConfig    // Logger config
    Bootstrap   BootstrapConfig // Recursos opcionales
}
```

### Uso
```go
import "github.com/EduGoGroup/edugo-shared/config"

// Cargar configuración
cfg, err := config.Load("config.yaml")

// Validar configuración
if err := config.Validate(cfg); err != nil {
    log.Fatal(err)
}
```

---

## 5. Database Module

**Ubicación:** `database/`

### 5.1 PostgreSQL

```go
import "github.com/EduGoGroup/edugo-shared/database/postgres"

// Conectar
db, err := postgres.Connect(cfg)

// Health check
err := postgres.HealthCheck(db)

// Estadísticas del pool
stats := postgres.GetStats(db)

// Cerrar
postgres.Close(db)

// Transacciones
err := postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
    // operaciones...
    return nil
})
```

### 5.2 MongoDB

```go
import "github.com/EduGoGroup/edugo-shared/database/mongodb"

// Conectar
client, err := mongodb.Connect(cfg)

// Obtener database
db := mongodb.GetDatabase(client, "dbname")

// Health check
err := mongodb.HealthCheck(client)

// Cerrar
mongodb.Close(client)
```

---

## 6. Lifecycle Module

**Ubicación:** `lifecycle/`  
**Import:** `github.com/EduGoGroup/edugo-shared/lifecycle`

### Propósito
Gestión ordenada de startup y shutdown de recursos.

### Uso
```go
import "github.com/EduGoGroup/edugo-shared/lifecycle"

// Crear manager
manager := lifecycle.NewManager(logger)

// Registrar recurso con startup y cleanup
manager.Register("database",
    func(ctx context.Context) error { return db.Connect() }, // startup
    func() error { return db.Close() },                      // cleanup
)

// Registrar solo cleanup
manager.RegisterSimple("cache", func() error {
    return cache.Close()
})

// Ejecutar startup (en orden de registro)
err := manager.Startup(ctx)

// Ejecutar cleanup (en orden inverso - LIFO)
err := manager.Cleanup()
```

---

## 7. Logger Module

**Ubicación:** `logger/`  
**Import:** `github.com/EduGoGroup/edugo-shared/logger`

### Interface
```go
type Logger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    Fatal(msg string, fields ...interface{})
    With(fields ...interface{}) Logger
    Sync() error
}
```

### Uso (Zap Implementation)
```go
import "github.com/EduGoGroup/edugo-shared/logger"

// Crear logger
log, err := logger.NewZapLogger("info", "json")

// Logging
log.Info("operación exitosa", "user_id", "123", "action", "login")

// Con contexto
userLog := log.With("user_id", "123")
userLog.Info("operación")

// Sync antes de cerrar
defer log.Sync()
```

---

## 8. Messaging Module

**Ubicación:** `messaging/rabbit/`  
**Import:** `github.com/EduGoGroup/edugo-shared/messaging/rabbit`

### Connection
```go
import "github.com/EduGoGroup/edugo-shared/messaging/rabbit"

// Crear conexión
conn, err := rabbit.NewConnection(cfg)
defer conn.Close()
```

### Publisher
```go
// Crear publisher
pub := rabbit.NewPublisher(conn)

// Publicar mensaje
err := pub.Publish(ctx, "exchange", "routing.key", message)

// Con prioridad
err := pub.PublishWithPriority(ctx, "exchange", "routing.key", message, 5)
```

### Consumer
```go
// Crear consumer
consumer := rabbit.NewConsumer(conn, rabbit.ConsumerConfig{
    Name:      "my-consumer",
    AutoAck:   false,
    Exclusive: false,
})

// Consumir mensajes
err := consumer.Consume(ctx, "queue-name", func(ctx context.Context, body []byte) error {
    var msg MyMessage
    if err := rabbit.UnmarshalMessage(body, &msg); err != nil {
        return err
    }
    // procesar mensaje...
    return nil
})
```

### Dead Letter Queue (DLQ)
```go
// Configurar DLQ automático
dlqConfig := rabbit.DLQConfig{
    QueueName:    "my-queue",
    DLQName:      "my-queue.dlq",
    MaxRetries:   3,
    RetryDelay:   time.Second * 5,
}
```

---

## 9. Middleware Module

**Ubicación:** `middleware/gin/`  
**Import:** `github.com/EduGoGroup/edugo-shared/middleware/gin`

### JWT Auth Middleware
```go
import ginmw "github.com/EduGoGroup/edugo-shared/middleware/gin"

router := gin.New()

// Aplicar middleware
router.Use(ginmw.JWTAuthMiddleware(jwtManager))

// En handlers, obtener datos del token
router.GET("/profile", func(c *gin.Context) {
    userID, _ := ginmw.GetUserID(c)
    email, _ := ginmw.GetEmail(c)
    role, _ := ginmw.GetRole(c)
    
    // O todos los claims
    claims, _ := ginmw.GetClaims(c)
})
```

### Context Keys
```go
const (
    ContextKeyUserID = "user_id"
    ContextKeyEmail  = "email"
    ContextKeyRole   = "role"
    ContextKeyClaims = "jwt_claims"
)
```

---

## 10. Testing Module

**Ubicación:** `testing/containers/`  
**Import:** `github.com/EduGoGroup/edugo-shared/testing/containers`

### Setup
```go
import "github.com/EduGoGroup/edugo-shared/testing/containers"

func TestMain(m *testing.M) {
    config := containers.NewConfig().
        WithPostgreSQL(nil).
        WithMongoDB(nil).
        WithRabbitMQ(nil).
        Build()
    
    manager, err := containers.GetManager(nil, config)
    if err != nil {
        panic(err)
    }
    defer manager.Cleanup(context.Background())
    
    os.Exit(m.Run())
}
```

### Uso en Tests
```go
func TestDatabase(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // PostgreSQL
    db := manager.PostgreSQL().DB()
    
    // MongoDB
    mongoDb := manager.MongoDB().Database()
    
    // RabbitMQ
    ch, _ := manager.RabbitMQ().Channel()
    
    // Cleanup
    ctx := context.Background()
    manager.CleanPostgreSQL(ctx, "users", "orders")
    manager.CleanMongoDB(ctx)
}
```

---

## Dependencias entre Módulos

```
                    ┌──────────┐
                    │  common  │
                    │ (errors, │
                    │  types)  │
                    └────┬─────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
    ┌─────────┐    ┌──────────┐    ┌──────────┐
    │  auth   │    │  logger  │    │  config  │
    └────┬────┘    └────┬─────┘    └────┬─────┘
         │              │               │
         └──────────────┼───────────────┘
                        │
                        ▼
                  ┌───────────┐
                  │ lifecycle │
                  └─────┬─────┘
                        │
         ┌──────────────┼──────────────┐
         │              │              │
         ▼              ▼              ▼
    ┌──────────┐  ┌───────────┐  ┌────────────┐
    │ database │  │ messaging │  │ middleware │
    └──────────┘  └───────────┘  └────────────┘
         │              │              │
         └──────────────┼──────────────┘
                        │
                        ▼
                  ┌───────────┐
                  │ bootstrap │
                  └───────────┘
```
