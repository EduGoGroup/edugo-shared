# Testing — Documentación técnica

Infraestructura de testing basada en Testcontainers para PostgreSQL, MongoDB y RabbitMQ.

## Propósito

Proporcionar infraestructura de testing que:
- Gestiona ciclo de vida de containers (startup y cleanup)
- Expone acceso a bases de datos y brokers para integración
- Reutiliza containers entre múltiples tests (singleton)
- Proporciona utilities para limpiar estado entre tests
- Maneja health checks y retries automáticos

## Componentes principales

### ConfigBuilder — Construcción fluida

Builder que configura qué containers crear.

**Constructor:**
```go
func NewConfig() *ConfigBuilder
```

**Métodos:**
- `WithPostgres(enabled bool) *ConfigBuilder` — Habilitar/deshabilitar PostgreSQL
- `WithMongoDB(enabled bool) *ConfigBuilder` — Habilitar/deshabilitar MongoDB
- `WithRabbitMQ(enabled bool) *ConfigBuilder` — Habilitar/deshabilitar RabbitMQ
- `WithNetwork(name string) *ConfigBuilder` — Especificar red Docker
- `GetManager(ctx context.Context) (*Manager, error)` — Crear/obtener singleton

**Ejemplo:**
```go
config := containers.NewConfig().
    WithPostgres(true).
    WithMongoDB(true).
    WithRabbitMQ(false)

mgr, err := config.GetManager(context.Background())
```

### Manager — Gestor singleton

Centraliza acceso y vida útil de todos los containers.

**Métodos principales:**
- `PostgreSQL() *PostgresContainer` — Obtener container PostgreSQL
- `MongoDB() *MongoDBContainer` — Obtener container MongoDB
- `RabbitMQ() *RabbitMQContainer` — Obtener container RabbitMQ
- `TruncatePostgres(ctx, tables []string) error` — Truncar tablas
- `DropMongoDB(ctx, collections []string) error` — Dropear colecciones
- `PurgeRabbitMQ(ctx, queues []string) error` — Purgar colas
- `Cleanup(ctx context.Context) error` — Parar todos los containers

**Características:**
- Singleton pattern: única instancia por proceso
- Thread-safe: protección interna con mutex
- Lazy initialization: containers se crean solo si están habilitados

### PostgresContainer — Wrapper PostgreSQL

Expone acceso a base de datos PostgreSQL con utilidades.

**Métodos:**
- `DB() *gorm.DB` — Obtener conexión GORM
- `Port() int` — Obtener puerto mapeado
- `Host() string` — Obtener host (localhost)
- `User() string` — Obtener usuario de conexión
- `Password() string` — Obtener password
- `Database() string` — Obtener nombre de base de datos
- `WaitForHealthy(ctx context.Context) error` — Esperar a que esté listo
- `ExecSQL(ctx context.Context, sql string) error` — Ejecutar SQL raw
- `ExecSQLFile(ctx context.Context, filePath string) error` — Ejecutar script SQL
- `Stop(ctx context.Context) error` — Parar container

**Ejemplo:**
```go
pgContainer := mgr.PostgreSQL()
db := pgContainer.DB()

// Ejecutar operación
var count int64
db.Model(&User{}).Count(&count)

// Ejecutar SQL raw
err := pgContainer.ExecSQL(ctx, "CREATE TABLE IF NOT EXISTS logs (...)")
```

### MongoDBContainer — Wrapper MongoDB

Expone acceso a MongoDB con utilidades.

**Métodos:**
- `Client() *mongo.Client` — Obtener cliente MongoDB
- `Database(name string) *mongo.Database` — Obtener base de datos
- `Port() int` — Obtener puerto mapeado
- `Host() string` — Obtener host
- `WaitForHealthy(ctx context.Context) error` — Esperar a que esté listo
- `Stop(ctx context.Context) error` — Parar container

**Ejemplo:**
```go
mongoContainer := mgr.MongoDB()
client := mongoContainer.Client()

// Acceder a colección
collection := mongoContainer.Database("myapp").Collection("users")

// Operación
result, err := collection.InsertOne(ctx, bson.M{"name": "John"})
```

### RabbitMQContainer — Wrapper RabbitMQ

Expone acceso a RabbitMQ con utilidades.

**Métodos:**
- `Connection() *amqp.Connection` — Obtener conexión AMQP
- `Port() int` — Obtener puerto mapeado
- `Host() string` — Obtener host
- `User() string` — Obtener usuario
- `Password() string` — Obtener password
- `WaitForHealthy(ctx context.Context) error` — Esperar a que esté listo
- `Stop(ctx context.Context) error` — Parar container

**Ejemplo:**
```go
rmqContainer := mgr.RabbitMQ()
conn := rmqContainer.Connection()

// Crear canal
ch, err := conn.Channel()
defer ch.Close()

// Declarar queue
q, err := ch.QueueDeclare("my_queue", true, false, false, false, nil)
```

## Flujos comunes

### 1. Configurar containers para suite de tests

```go
import (
    "testing"
    "github.com/EduGoGroup/edugo-shared/testing/containers"
)

var mgr *containers.Manager

func TestMain(m *testing.M) {
    ctx := context.Background()

    // Configurar containers necesarios
    config := containers.NewConfig().
        WithPostgres(true).
        WithMongoDB(true).
        WithRabbitMQ(false)

    var err error
    mgr, err = config.GetManager(ctx)
    if err != nil {
        log.Fatalf("failed to get manager: %v", err)
    }

    // Ejecutar tests
    code := m.Run()

    // Limpiar
    mgr.Cleanup(ctx)

    os.Exit(code)
}
```

### 2. Limpiar estado PostgreSQL entre tests

```go
func setupPostgresTest(t *testing.T) {
    ctx := context.Background()

    // Truncar tablas relevantes
    tables := []string{"users", "schools", "memberships"}
    if err := mgr.TruncatePostgres(ctx, tables); err != nil {
        t.Fatalf("failed to truncate: %v", err)
    }

    // Ejecutar migrations si es necesario
    db := mgr.PostgreSQL().DB()
    // ... ejecutar migrations ...
}

func TestCreateUser(t *testing.T) {
    setupPostgresTest(t)

    db := mgr.PostgreSQL().DB()
    user := User{Name: "John"}

    if err := db.Create(&user).Error; err != nil {
        t.Fatalf("failed to create user: %v", err)
    }
}
```

### 3. Usar múltiples containers en test

```go
func TestCrossDatabaseIntegration(t *testing.T) {
    ctx := context.Background()

    // PostgreSQL para datos persistentes
    pgDB := mgr.PostgreSQL().DB()
    user := User{Name: "Alice"}
    pgDB.Create(&user)

    // MongoDB para cache/profiles
    mongoCollection := mgr.MongoDB().Database("app").Collection("profiles")
    mongoCollection.InsertOne(ctx, bson.M{
        "user_id": user.ID,
        "cached":  true,
    })

    // RabbitMQ para eventos
    ch, _ := mgr.RabbitMQ().Connection().Channel()
    q, _ := ch.QueueDeclare("events", true, false, false, false, nil)
    ch.Publish("", q.Name, false, false, amqp.Publishing{
        Body: []byte(fmt.Sprintf("user_created:%d", user.ID)),
    })

    // Assertions
    assert.NotZero(t, user.ID)
}
```

### 4. Esperar a que containers estén listos con retries

```go
func TestWithRetries(t *testing.T) {
    ctx := context.Background()

    // Esperar a que PostgreSQL esté completamente listo
    if err := mgr.PostgreSQL().WaitForHealthy(ctx); err != nil {
        t.Fatalf("postgres failed to become healthy: %v", err)
    }

    // RetryOperation para operación que podría fallar al inicio
    result := containers.RetryOperation(
        func() (interface{}, error) {
            db := mgr.PostgreSQL().DB()
            var count int64
            return count, db.Model(&User{}).Count(&count).Error
        },
        5,                          // máximo 5 intentos
        100*time.Millisecond,       // 100ms entre intentos
    )

    if result.Err != nil {
        t.Fatalf("operation failed after retries: %v", result.Err)
    }
}
```

## Arquitectura

Flujo del ciclo de vida:

```
1. NewConfig()
   ↓
2. WithPostgres(true)
   WithMongoDB(true)
   WithRabbitMQ(false)
   ↓
3. GetManager(ctx)
   ├─ Crear containers habilitados
   ├─ Mapear puertos
   └─ Retornar singleton Manager
   ↓
4. [Tests ejecutándose]
   ├─ mgr.PostgreSQL().DB()
   ├─ mgr.MongoDB().Client()
   ├─ Operaciones en containers
   └─ Cleanup entre tests
   ↓
5. mgr.Cleanup(ctx)
   ├─ Stop PostgreSQL
   ├─ Stop MongoDB
   ├─ Stop RabbitMQ
   └─ Cleanup Docker resources
```

Diagrama de componentes:

```
ConfigBuilder
      ↓
GetManager (singleton)
      ↓
Manager
  ├─ PostgresContainer
  │  ├─ gorm.DB
  │  ├─ WaitForHealthy
  │  └─ TruncatePostgres
  ├─ MongoDBContainer
  │  ├─ mongo.Client
  │  ├─ WaitForHealthy
  │  └─ DropMongoDB
  └─ RabbitMQContainer
     ├─ amqp.Connection
     ├─ WaitForHealthy
     └─ PurgeRabbitMQ
```

## Dependencias

- **Internas**: Ninguna (módulo autónomo)
- **Externas**:
  - `github.com/testcontainers/testcontainers-go` (Testcontainers core)
  - `github.com/testcontainers/testcontainers-go/wait` (Wait strategies)
  - `gorm.io/gorm` (para PostgresContainer)
  - `go.mongodb.org/mongo-go-driver` (para MongoDBContainer)
  - `github.com/rabbitmq/amqp091-go` (para RabbitMQContainer)
  - `github.com/docker/docker` (indirectamente, via testcontainers)

## Testing

Suite de tests completa:

- Creación y obtención de Manager singleton
- Configuración de containers habilitados/deshabilitados
- Acceso a conexiones (GORM, MongoDB client, AMQP connection)
- Health checks y retries
- Cleanup de estado (truncate, drop, purge)
- Ejecución de SQL desde archivo
- Tests unitarios (sin Docker) e integración (con Docker)

Ejecutar:
```bash
make test          # Tests unitarios (sin Docker)
make test-race     # Tests con race detector (sin Docker)
make check         # Tests + linting + format
make test-all      # Tests con integración Docker (requiere Docker)
```

## Notas de diseño

- **Singleton Manager**: Centraliza vida útil de containers para abaratar suites largas (containers se crean una sola vez)
- **ConfigBuilder fluido**: API clara y composable para configurar containers
- **Cleanup tipado**: Helpers específicos por backend (TruncatePostgres, DropMongoDB, PurgeRabbitMQ)
- **Sin framework específico**: Funciona con testing, testify, ginkgo, etc.
- **Health checks**: Retries automáticos para esperar a que containers estén listos
- **Docker agnóstico**: Testcontainers abstrae detalles de Docker (funciona con Docker/Podman)
- **Lazy initialization**: Containers se crean solo si están habilitados
- **Thread-safe**: Protección interna para acceso concurrente
