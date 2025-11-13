# Testing Module - EduGo Shared

Módulo de testing infrastructure que proporciona containers reutilizables para tests de integración usando [testcontainers-go](https://github.com/testcontainers/testcontainers-go).

## 🎯 Características

- ✅ **Manager Singleton**: Los containers se crean una sola vez y se reutilizan entre tests
- ✅ **Builder Pattern**: Configuración flexible e intuitiva
- ✅ **3 Containers Soportados**: PostgreSQL, MongoDB, RabbitMQ
- ✅ **Cleanup Automático**: Gestión de recursos y limpieza entre tests
- ✅ **Performance**: Setup inicial ~17s, tests subsiguientes instantáneos

## 📦 Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/testing@latest
```

## 🚀 Uso Rápido

### Ejemplo Básico: PostgreSQL

```go
package myapp_test

import (
    "context"
    "os"
    "testing"
    
    "github.com/EduGoGroup/edugo-shared/testing/containers"
)

func TestMain(m *testing.M) {
    // Configurar containers
    config := containers.NewConfig().
        WithPostgreSQL(nil). // nil usa defaults
        Build()
    
    // Obtener manager (singleton)
    manager, err := containers.GetManager(nil, config)
    if err != nil {
        panic(err)
    }
    
    // Cleanup al finalizar
    defer func() {
        ctx := context.Background()
        if err := manager.Cleanup(ctx); err != nil {
            panic(err)
        }
    }()
    
    // Ejecutar tests
    os.Exit(m.Run())
}

func TestDatabaseOperation(t *testing.T) {
    // Reutilizar manager
    manager, _ := containers.GetManager(t, nil)
    
    // Acceder a la base de datos
    db := manager.PostgreSQL().DB()
    
    // Crear tabla
    _, err := db.Exec(`CREATE TABLE users (id SERIAL, name TEXT)`)
    if err != nil {
        t.Fatal(err)
    }
    
    // Insertar datos
    _, err = db.Exec(`INSERT INTO users (name) VALUES ('Alice')`)
    if err != nil {
        t.Fatal(err)
    }
    
    // Verificar
    var count int
    db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
    if count != 1 {
        t.Errorf("Esperado 1 usuario, obtenido %d", count)
    }
    
    // Limpiar para el siguiente test
    ctx := context.Background()
    manager.CleanPostgreSQL(ctx, "users")
}
```

### Ejemplo: Todos los Containers

```go
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

func TestFullStack(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // PostgreSQL
    db := manager.PostgreSQL().DB()
    db.Exec(`CREATE TABLE events (id SERIAL, data TEXT)`)
    
    // MongoDB
    mongoDb := manager.MongoDB().Database()
    coll := mongoDb.Collection("logs")
    coll.InsertOne(context.Background(), bson.M{"level": "info"})
    
    // RabbitMQ
    ch, _ := manager.RabbitMQ().Channel()
    ch.QueueDeclare("tasks", false, false, false, false, nil)
    
    // Tests...
}
```

## 📖 Configuración Avanzada

### PostgreSQL con Scripts de Inicialización

```go
config := containers.NewConfig().
    WithPostgreSQL(&containers.PostgresConfig{
        Image:    "postgres:15-alpine",
        Database: "my_test_db",
        Username: "test_user",
        Password: "test_pass",
        InitScripts: []string{
            "../../migrations/001_create_schema.sql",
            "../../migrations/002_seed_data.sql",
        },
    }).
    Build()
```

### MongoDB con Autenticación

```go
config := containers.NewConfig().
    WithMongoDB(&containers.MongoConfig{
        Image:    "mongo:7.0",
        Database: "test_db",
        Username: "admin",
        Password: "secret",
    }).
    Build()
```

### RabbitMQ Personalizado

```go
config := containers.NewConfig().
    WithRabbitMQ(&containers.RabbitConfig{
        Image:    "rabbitmq:3.12-management-alpine",
        Username: "guest",
        Password: "guest",
    }).
    Build()
```

## 🧹 Limpieza Entre Tests

### Truncate de Tablas PostgreSQL

```go
func TestSomething(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // ... test logic ...
    
    // Limpiar tablas al final
    ctx := context.Background()
    manager.CleanPostgreSQL(ctx, "users", "orders", "products")
}
```

### Drop de Colecciones MongoDB

```go
func TestMongo(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // ... test logic ...
    
    // Eliminar todas las colecciones
    ctx := context.Background()
    manager.CleanMongoDB(ctx)
}
```

### Purgar Colas RabbitMQ

```go
func TestRabbit(t *testing.T) {
    manager, _ := containers.GetManager(t, nil)
    
    // ... test logic ...
    
    // Purgar cola específica
    manager.RabbitMQ().PurgeQueue("my_queue")
}
```

## 🔍 API Reference

### Manager

```go
// Obtener container específico
manager.PostgreSQL() *PostgresContainer
manager.MongoDB() *MongoDBContainer
manager.RabbitMQ() *RabbitMQContainer

// Limpieza
manager.Cleanup(ctx) error
manager.CleanPostgreSQL(ctx, ...tables) error
manager.CleanMongoDB(ctx) error
manager.PurgeRabbitMQ(ctx) error
```

### PostgresContainer

```go
pg := manager.PostgreSQL()

pg.DB() *sql.DB
pg.ConnectionString(ctx) (string, error)
pg.Truncate(ctx, ...tables) error
pg.Terminate(ctx) error
```

### MongoDBContainer

```go
mongo := manager.MongoDB()

mongo.Client() *mongo.Client
mongo.Database() *mongo.Database
mongo.ConnectionString(ctx) (string, error)
mongo.DropAllCollections(ctx) error
mongo.DropCollections(ctx, ...names) error
mongo.Terminate(ctx) error
```

### RabbitMQContainer

```go
rabbit := manager.RabbitMQ()

rabbit.Connection() *amqp.Connection
rabbit.Channel() (*amqp.Channel, error)
rabbit.ConnectionString(ctx) (string, error)
rabbit.PurgeQueue(name) error
rabbit.DeleteQueue(name) error
rabbit.Terminate(ctx) error
```

## ⚙️ Configuración por Defecto

| Container | Imagen | Database | Usuario | Password |
|-----------|--------|----------|---------|----------|
| **PostgreSQL** | `postgres:15-alpine` | `edugo_test` | `edugo_user` | `edugo_pass` |
| **MongoDB** | `mongo:7.0` | `edugo_test` | - | - |
| **RabbitMQ** | `rabbitmq:3.12-alpine` | - | `edugo_user` | `edugo_pass` |

## 📊 Performance

- **Primera ejecución**: ~17s (crear 3 containers)
- **Tests subsiguientes**: <1s (reutiliza containers)
- **Cleanup entre tests**: ~100ms (truncate/drop)
- **Memoria**: ~800MB (todos los containers)

## 🐛 Troubleshooting

### Error: "Docker not running"
```bash
# Verificar que Docker está corriendo
docker ps
```

### Error: "Port already in use"
```bash
# Los containers usan puertos aleatorios, este error no debería ocurrir
# Si ocurre, limpiar containers huérfanos:
docker ps -a | grep testcontainers | awk '{print $1}' | xargs docker rm -f
```

### Tests lentos
```bash
# Usar -short para skipear tests de integración
go test -short ./...
```

## 🔗 Recursos

- [testcontainers-go Documentation](https://golang.testcontainers.org/)
- [Repositorio EduGo](https://github.com/EduGoGroup/edugo-shared)

## 📝 Licencia

Uso interno EduGo - Todos los derechos reservados
