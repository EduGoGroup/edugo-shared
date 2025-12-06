# TECH STACK - SHARED

## Resumen Ejecutivo

| Layer | Tecnología | Versión | Propósito |
|-------|-----------|---------|----------|
| **Language** | Go | 1.21+ | Compilación de módulos |
| **Module System** | Go Modules | Native | Versionamiento |
| **Logging** | Logrus | v1.9+ | Structured logging |
| **Auth** | golang-jwt | v5+ | JWT handling |
| **Database** | GORM | v1.25+ | ORM abstraction |
| **PostgreSQL** | go-sql-driver | v1.5+ | PostgreSQL connection |
| **MongoDB** | mongo-go-driver | v1.12+ | MongoDB connection |
| **Messaging** | streadway/amqp | v1.0+ | RabbitMQ client |
| **Config** | Viper | v1.16+ | Configuration |

---

## Stack por Módulo

### Logger Module
```go
// Logrus - Structured logging
import "github.com/sirupsen/logrus"

type Logger interface {
    Info(msg string, fields map[string]interface{})
    Error(msg string, fields map[string]interface{})
    Debug(msg string, fields map[string]interface{})
    Warn(msg string, fields map[string]interface{})
}

// JSON output format
{
  "level": "info",
  "msg": "User created",
  "timestamp": "2025-11-15T10:30:00Z",
  "user_id": 42,
  "service": "api-mobile"
}
```

**Performance:**
- Async logging para no bloquear
- Buffering automático
- Rotación de files (opcional)

---

### Auth Module
```go
// golang-jwt/jwt - JWT handling
import "github.com/golang-jwt/jwt/v5"

type Claims struct {
    UserID   int64    `json:"user_id"`
    SchoolID int64    `json:"school_id"`
    Email    string   `json:"email"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

// Crear token
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, _ := token.SignedString(secretKey)

// Validar token
claims := &Claims{}
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    return secretKey, nil
})
```

**Features:**
- HS256 (HMAC)
- RS256 (RSA, para key rotation)
- Custom claims
- Token expiration
- Refresh tokens (optional)

---

### Database - PostgreSQL

```go
// GORM + PostgreSQL driver
import "gorm.io/gorm"
import "gorm.io/driver/postgres"

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    SkipDefaultTransaction: true,
    NowFunc: func() time.Time {
        return time.Now().UTC()
    },
})

// Connection pooling
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

**Characteristics:**
- GORM abstraction
- Transaction support
- Query logging
- Prepared statements
- CTEs for recursion
- Window functions

---

### Database - MongoDB

```go
// mongo-go-driver
import "go.mongodb.org/mongo-driver/mongo"

client, err := mongo.Connect(context.Background(), options.Client().
    ApplyURI(mongoURI))

db := client.Database("edugo")
collection := db.Collection("evaluation_results")

// CRUD
result, _ := collection.InsertOne(ctx, document)
_ = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
```

**Features:**
- Connection pooling
- Aggregation pipeline
- TTL indexes
- Full-text search
- Transactions (v1.1+)

---

### Messaging - RabbitMQ

```go
// streadway/amqp
import "github.com/streadway/amqp"

conn, err := amqp.Dial("amqp://user:pass@localhost:5672/")
ch, err := conn.Channel()

// Declare
ch.ExchangeDeclare("assessment.requests", "direct", true, false, false, false, nil)
ch.QueueDeclare("worker.assessment.requests", true, false, false, false, nil)

// Publish
ch.Publish("assessment.requests", "worker.assessment.requests", false, false, amqp.Publishing{
    ContentType: "application/json",
    Body:        []byte(payload),
})

// Consume
msgs, _ := ch.Consume("api-mobile.responses", "", false, false, false, false, nil)
for msg := range msgs {
    msg.Ack(false)
}
```

**Features:**
- Connection pooling
- Auto-reconnect
- QoS management
- Dead letter queues
- Message acknowledgement

---

### Configuration - Viper

```go
// spf13/viper
import "github.com/spf13/viper"

viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.ReadInConfig()

// Environment variables
viper.AutomaticEnv()
viper.BindEnv("db.host", "DB_HOST")

// Get values
host := viper.GetString("db.host")
port := viper.GetInt("db.port")
timeout := viper.GetDuration("api.timeout")
```

**File Support:**
- YAML
- JSON
- TOML
- HCL
- Environment variables

---

## Arquitectura de Módulos

```
github.com/EduGoGroup/edugo-shared/

├── logger/
│   └── Using: logrus
│   └── Exports: Logger interface
│
├── database/
│   ├── postgres/
│   │   └── Using: GORM, PostgreSQL driver
│   │   └── Exports: *gorm.DB
│   │
│   └── mongo/
│       └── Using: mongo-go-driver
│       └── Exports: *mongo.Client
│
├── auth/
│   └── Using: golang-jwt
│   └── Exports: Claims, ValidateToken()
│
├── messaging/
│   └── Using: streadway/amqp
│   └── Exports: Publisher, Subscriber interfaces
│
├── models/
│   └── No external deps
│   └── Exports: User, School, etc structs
│
├── context/
│   └── No external deps
│   └── Exports: WithTimeout, WithUser, etc
│
├── errors/
│   └── No external deps
│   └── Exports: NotFound, BadRequest, etc
│
└── health/
    └── Using: Otros módulos
    └── Exports: HealthChecker interface
```

---

## Dependency Injection Pattern

SHARED proporciona construcción explícita de dependencias:

```go
// En api-mobile
func main() {
    // Paso 1: Logger
    log := logger.NewLogger(config.Logger)
    
    // Paso 2: Database
    pgDB := database.InitPostgres(config.Postgres)
    mongoClient := database.InitMongo(config.Mongo)
    
    // Paso 3: Auth middleware
    authMiddleware := auth.NewValidator(config.Auth)
    
    // Paso 4: Messaging
    publisher := messaging.NewPublisher(config.RabbitMQ)
    subscriber := messaging.NewSubscriber(config.RabbitMQ)
    
    // Paso 5: Services (usando deps inyectadas)
    evaluationService := services.NewEvaluationService(
        pgDB,
        mongoClient.Database("edugo"),
        publisher,
        log,
    )
    
    // Paso 6: Handlers
    handler := handlers.NewEvaluationHandler(evaluationService)
}
```

---

## Testing Strategy

### Unit Tests (sin deps externas)
```go
// errors/ module
func TestNotFoundError(t *testing.T) {
    err := errors.NewNotFound("User", "id", 42)
    assert.Equal(t, "not_found", err.Code)
}
```

### Integration Tests (con mocks)
```go
// database/ module
func TestPostgresConnection(t *testing.T) {
    // Usar PostgreSQL real en Docker
    db := testDatabase.NewPostgresDB()
    defer db.Close()
    
    var count int64
    db.Model(&User{}).Count(&count)
    assert.Equal(t, int64(0), count)
}
```

### Cross-Project Tests
```go
// Verificar que api-mobile funciona con SHARED
go test -run TestAPIMobileWithShared ./...
```

---

## Performance Characteristics

### Logging
- ~1-2 μs per log entry
- Async si se configura
- Buffering automático

### JWT
- ~100 μs para crear token
- ~500 μs para validar token
- Cacheable

### Database Connections
- ~10ms para primer query
- ~1-5ms para queries subsecuentes
- Connection pooling automático

### RabbitMQ
- ~5ms para publish
- ~2-5ms para consume
- Retry automático

---

## Comparativa: Go vs Alternativas

| Aspecto | Go | Python | Node.js |
|--------|----|---------| --------|
| Startup time | <100ms | 500ms+ | 100ms+ |
| Memory footprint | 10-50MB | 100MB+ | 50MB+ |
| Concurrency | Goroutines | Threads | Async/Await |
| Type safety | Strong | Weak | Weak |
| Compilation | Fast | Interpreted | Interpreted |

**Por qué Go para SHARED:**
- Single binary fácil de distribuir
- Performance crítica para librería base
- Type safety crucial para compatibilidad
- Excelente para aplicaciones backend
