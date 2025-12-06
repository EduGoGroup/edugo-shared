# Guía de Testing

## Visión General

edugo-shared proporciona infraestructura de testing robusta basada en testcontainers-go para tests de integración con servicios reales.

---

## Estructura de Testing

```
testing/
├── containers/
│   ├── manager.go          # Manager singleton de containers
│   ├── config.go           # Configuración de containers
│   ├── postgres.go         # Container PostgreSQL
│   ├── mongodb.go          # Container MongoDB
│   └── rabbitmq.go         # Container RabbitMQ
└── README.md               # Documentación del módulo
```

---

## Configuración Rápida

### Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/testing@latest
```

### TestMain Básico

```go
package myapp_test

import (
    "context"
    "os"
    "testing"
    
    "github.com/EduGoGroup/edugo-shared/testing/containers"
)

func TestMain(m *testing.M) {
    // 1. Configurar containers necesarios
    config := containers.NewConfig().
        WithPostgreSQL(nil).  // nil = usar defaults
        WithMongoDB(nil).
        Build()
    
    // 2. Obtener manager (crea containers si no existen)
    manager, err := containers.GetManager(nil, config)
    if err != nil {
        panic(err)
    }
    
    // 3. Cleanup al finalizar
    defer func() {
        ctx := context.Background()
        if err := manager.Cleanup(ctx); err != nil {
            panic(err)
        }
    }()
    
    // 4. Ejecutar tests
    os.Exit(m.Run())
}
```

---

## Containers Disponibles

### PostgreSQL

```go
// Configuración por defecto
config := containers.NewConfig().
    WithPostgreSQL(nil).
    Build()

// Configuración personalizada
config := containers.NewConfig().
    WithPostgreSQL(&containers.PostgresConfig{
        Image:    "postgres:15-alpine",
        Database: "my_test_db",
        Username: "test_user",
        Password: "test_pass",
        InitScripts: []string{
            "./migrations/001_schema.sql",
            "./fixtures/seed_data.sql",
        },
    }).
    Build()
```

#### Valores por Defecto PostgreSQL

| Campo | Valor |
|-------|-------|
| Image | `postgres:15-alpine` |
| Database | `edugo_test` |
| Username | `edugo_user` |
| Password | `edugo_pass` |

#### Uso en Tests

```go
func TestUserRepository(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // Obtener conexión
    db := manager.PostgreSQL().DB()
    
    // Crear tabla de prueba
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL
        )
    `)
    if err != nil {
        t.Fatal(err)
    }
    
    // Insertar datos de prueba
    _, err = db.Exec(`INSERT INTO users (name, email) VALUES ($1, $2)`, "Alice", "alice@test.com")
    if err != nil {
        t.Fatal(err)
    }
    
    // Test assertions...
    
    // Limpiar al final
    ctx := context.Background()
    manager.CleanPostgreSQL(ctx, "users")
}
```

### MongoDB

```go
// Configuración por defecto
config := containers.NewConfig().
    WithMongoDB(nil).
    Build()

// Configuración personalizada
config := containers.NewConfig().
    WithMongoDB(&containers.MongoConfig{
        Image:    "mongo:7.0",
        Database: "my_test_db",
        Username: "admin",
        Password: "secret",
    }).
    Build()
```

#### Valores por Defecto MongoDB

| Campo | Valor |
|-------|-------|
| Image | `mongo:7.0` |
| Database | `edugo_test` |
| Username | - |
| Password | - |

#### Uso en Tests

```go
func TestAuditLogRepository(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // Obtener database
    db := manager.MongoDB().Database()
    
    // Obtener collection
    logs := db.Collection("audit_logs")
    
    // Insertar documento
    ctx := context.Background()
    _, err := logs.InsertOne(ctx, bson.M{
        "action":    "user.created",
        "user_id":   "123",
        "timestamp": time.Now(),
    })
    if err != nil {
        t.Fatal(err)
    }
    
    // Test assertions...
    
    // Limpiar al final
    manager.CleanMongoDB(ctx)
}
```

### RabbitMQ

```go
// Configuración por defecto
config := containers.NewConfig().
    WithRabbitMQ(nil).
    Build()

// Configuración personalizada
config := containers.NewConfig().
    WithRabbitMQ(&containers.RabbitConfig{
        Image:    "rabbitmq:3.12-management-alpine",
        Username: "guest",
        Password: "guest",
    }).
    Build()
```

#### Valores por Defecto RabbitMQ

| Campo | Valor |
|-------|-------|
| Image | `rabbitmq:3.12-alpine` |
| Username | `edugo_user` |
| Password | `edugo_pass` |

#### Uso en Tests

```go
func TestMessagePublisher(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // Obtener channel
    ch, err := manager.RabbitMQ().Channel()
    if err != nil {
        t.Fatal(err)
    }
    
    // Declarar cola
    q, err := ch.QueueDeclare(
        "test_queue",
        false, // durable
        true,  // auto-delete
        false, // exclusive
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        t.Fatal(err)
    }
    
    // Publicar mensaje
    err = ch.Publish(
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte("test message"),
        },
    )
    if err != nil {
        t.Fatal(err)
    }
    
    // Test assertions...
    
    // Limpiar
    manager.RabbitMQ().PurgeQueue("test_queue")
}
```

---

## API del Manager

### Obtener Containers

```go
manager, _ := containers.GetManager(t, nil)

// Acceder a containers específicos
pg := manager.PostgreSQL()      // *PostgresContainer
mongo := manager.MongoDB()      // *MongoDBContainer
rabbit := manager.RabbitMQ()    // *RabbitMQContainer
```

### PostgresContainer

```go
pg := manager.PostgreSQL()

// Conexión *sql.DB
db := pg.DB()

// String de conexión
connStr, err := pg.ConnectionString(ctx)

// Truncar tablas
err := pg.Truncate(ctx, "users", "orders", "products")

// Terminar container
err := pg.Terminate(ctx)
```

### MongoDBContainer

```go
mongo := manager.MongoDB()

// Cliente *mongo.Client
client := mongo.Client()

// Database *mongo.Database
db := mongo.Database()

// String de conexión
connStr, err := mongo.ConnectionString(ctx)

// Eliminar todas las colecciones
err := mongo.DropAllCollections(ctx)

// Eliminar colecciones específicas
err := mongo.DropCollections(ctx, "logs", "events")

// Terminar container
err := mongo.Terminate(ctx)
```

### RabbitMQContainer

```go
rabbit := manager.RabbitMQ()

// Conexión *amqp.Connection
conn := rabbit.Connection()

// Nuevo canal
ch, err := rabbit.Channel()

// String de conexión
connStr, err := rabbit.ConnectionString(ctx)

// Purgar cola (vaciar mensajes)
err := rabbit.PurgeQueue("my_queue")

// Eliminar cola
err := rabbit.DeleteQueue("my_queue")

// Terminar container
err := rabbit.Terminate(ctx)
```

---

## Limpieza Entre Tests

### PostgreSQL

```go
func TestA(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    db := manager.PostgreSQL().DB()
    
    // Setup y test...
    
    // Limpiar tablas afectadas
    ctx := context.Background()
    manager.CleanPostgreSQL(ctx, "users", "orders")
}

func TestB(t *testing.T) {
    // Comienza con tablas vacías
    manager, _ := containers.GetManager(t, nil)
    // ...
}
```

### MongoDB

```go
func TestA(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // Setup y test...
    
    // Eliminar todas las colecciones
    ctx := context.Background()
    manager.CleanMongoDB(ctx)
}
```

### RabbitMQ

```go
func TestA(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // Setup y test...
    
    // Purgar colas
    ctx := context.Background()
    manager.PurgeRabbitMQ(ctx)
}
```

---

## Patrón Singleton

El manager usa patrón singleton para reutilizar containers entre tests:

```go
// Primera llamada: crea containers (~17s)
manager1, _ := containers.GetManager(nil, config)

// Llamadas subsiguientes: retorna mismo manager (instantáneo)
manager2, _ := containers.GetManager(t, nil)  // config es opcional

// manager1 == manager2 (misma instancia)
```

### Beneficios
- **Performance**: Setup inicial una sola vez
- **Recursos**: Un solo set de containers
- **Simplicidad**: No necesitas pasar referencias

---

## Performance

| Operación | Tiempo |
|-----------|--------|
| Primera ejecución (3 containers) | ~17s |
| Tests subsiguientes | <1s |
| Cleanup entre tests (truncate/drop) | ~100ms |
| Memoria total | ~800MB |

---

## Comandos Make para Testing

```bash
# Tests unitarios de todos los módulos
make test-all-modules

# Tests con race detection
make test-race-all-modules

# Tests con cobertura
make coverage-all-modules

# Tests cortos (skip integración)
make test-short

# Verificación completa
make check-all-modules
```

---

## Skip Tests de Integración

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Test de integración...
}
```

```bash
# Ejecutar solo tests cortos
go test -short ./...
```

---

## Fixtures y Seeds

### Scripts SQL

```go
config := containers.NewConfig().
    WithPostgreSQL(&containers.PostgresConfig{
        InitScripts: []string{
            "./testdata/schema.sql",
            "./testdata/fixtures.sql",
        },
    }).
    Build()
```

### fixtures.sql

```sql
-- testdata/fixtures.sql
INSERT INTO users (id, name, email, role) VALUES
    ('uuid-1', 'Alice', 'alice@test.com', 'admin'),
    ('uuid-2', 'Bob', 'bob@test.com', 'student'),
    ('uuid-3', 'Carol', 'carol@test.com', 'teacher');
```

### Fixtures Programáticos

```go
func setupFixtures(t *testing.T, db *sql.DB) {
    t.Helper()
    
    fixtures := []struct {
        name  string
        email string
    }{
        {"Alice", "alice@test.com"},
        {"Bob", "bob@test.com"},
    }
    
    for _, f := range fixtures {
        _, err := db.Exec(`INSERT INTO users (name, email) VALUES ($1, $2)`, f.name, f.email)
        if err != nil {
            t.Fatalf("Failed to insert fixture: %v", err)
        }
    }
}
```

---

## Troubleshooting

### Error: "Docker not running"

```bash
# Verificar Docker
docker ps

# Iniciar Docker (macOS)
open -a Docker
```

### Error: "Port already in use"

Los containers usan puertos aleatorios, pero si hay conflictos:

```bash
# Limpiar containers huérfanos
docker ps -a | grep testcontainers | awk '{print $1}' | xargs docker rm -f
```

### Tests Lentos

```bash
# Usar modo short para desarrollo
go test -short ./...

# Ver qué containers están activos
docker ps --filter "name=testcontainers"
```

### Memoria Insuficiente

Reducir containers activos simultáneamente:

```go
// Usar solo los containers necesarios
config := containers.NewConfig().
    WithPostgreSQL(nil).  // Solo PostgreSQL
    Build()
```

---

## Ejemplos Completos

### Test de Repositorio

```go
package repository_test

import (
    "context"
    "os"
    "testing"
    
    "github.com/EduGoGroup/edugo-shared/testing/containers"
)

var manager *containers.Manager

func TestMain(m *testing.M) {
    config := containers.NewConfig().
        WithPostgreSQL(nil).
        Build()
    
    var err error
    manager, err = containers.GetManager(nil, config)
    if err != nil {
        panic(err)
    }
    
    defer manager.Cleanup(context.Background())
    
    os.Exit(m.Run())
}

func TestUserRepository_Create(t *testing.T) {
    db := manager.PostgreSQL().DB()
    
    // Setup
    _, _ = db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL, name TEXT, email TEXT UNIQUE)`)
    defer manager.CleanPostgreSQL(context.Background(), "users")
    
    // Test
    repo := NewUserRepository(db)
    user, err := repo.Create(context.Background(), "Alice", "alice@test.com")
    
    // Assertions
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("Expected name 'Alice', got '%s'", user.Name)
    }
}

func TestUserRepository_FindByEmail(t *testing.T) {
    db := manager.PostgreSQL().DB()
    
    // Setup
    _, _ = db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL, name TEXT, email TEXT UNIQUE)`)
    _, _ = db.Exec(`INSERT INTO users (name, email) VALUES ('Bob', 'bob@test.com')`)
    defer manager.CleanPostgreSQL(context.Background(), "users")
    
    // Test
    repo := NewUserRepository(db)
    user, err := repo.FindByEmail(context.Background(), "bob@test.com")
    
    // Assertions
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if user.Name != "Bob" {
        t.Errorf("Expected name 'Bob', got '%s'", user.Name)
    }
}
```
