# FASE 2: Restauración de Tests

> **Prioridad**: ALTA  
> **Duración estimada**: 2-3 días  
> **Prerrequisitos**: Fase 1 completada  
> **Rama**: `fase-2-restauracion-tests`  
> **Objetivo**: Restaurar todos los tests de integración deshabilitados

---

## Flujo de Trabajo de Esta Fase

### 1. Inicio de la Fase

```bash
# Asegurarse de estar en dev actualizado
git checkout dev
git pull origin dev

# Crear rama de la fase
git checkout -b fase-2-restauracion-tests

# Verificar estado inicial
make build
make test-all-modules
```

### 2. Durante la Fase

- Ejecutar cada paso en orden
- Commit atómico después de cada paso completado
- Verificar que tests pasen después de cada cambio

### 3. Fin de la Fase

```bash
# Push de la rama
git push origin fase-2-restauracion-tests

# Crear PR en GitHub hacia dev
# - Título: "test: Fase 2 - Restauración de Tests"
# - Descripción: Lista de tests restaurados

# Esperar revisión de GitHub Copilot
# - DESCARTAR: Comentarios de traducción inglés/español
# - CORREGIR: Problemas importantes detectados
# - DOCUMENTAR: Lo que queda como deuda técnica futura

# Esperar pipelines (máx 10 min, revisar cada 1 min)
# - Si hay errores: Corregir (regla de 3 intentos)
# - Todos los errores se corrigen (propios o heredados)

# Merge cuando todo esté verde
```

---

## Resumen de la Fase

Esta fase restaura los archivos de test que fueron renombrados a `.skip` debido a cambios en la API de containers.

### Archivos a Restaurar

| Archivo | Líneas | Estado |
|---------|--------|--------|
| `factory_mongodb_integration_test.go.skip` | ~529 | ⏳ Pendiente |
| `factory_postgresql_integration_test.go.skip` | ~500 | ⏳ Pendiente |
| `factory_rabbitmq_integration_test.go.skip` | ~450 | ⏳ Pendiente |

---

## API de Containers (Referencia)

La API de containers para tests de integración:

### Manager

```go
// GetManager obtiene o crea el singleton del manager
func GetManager(t *testing.T, config *Config) (*Manager, error)

// Métodos del Manager
func (m *Manager) PostgreSQL() *PostgreSQLContainer
func (m *Manager) MongoDB() *MongoDBContainer
func (m *Manager) RabbitMQ() *RabbitMQContainer
func (m *Manager) Cleanup(ctx context.Context) error
```

### PostgreSQL Container

```go
// Config retorna la configuración de conexión
func (c *PostgreSQLContainer) Config(ctx context.Context) (*PostgreSQLConfig, error)

// ConnectionString retorna el DSN de conexión
func (c *PostgreSQLContainer) ConnectionString(ctx context.Context) (string, error)
```

### MongoDB Container

```go
// ConnectionString retorna la URI de conexión
func (c *MongoDBContainer) ConnectionString(ctx context.Context) (string, error)
```

### RabbitMQ Container

```go
// ConnectionString retorna la URL AMQP
func (c *RabbitMQContainer) ConnectionString(ctx context.Context) (string, error)
```

---

## Paso 2.1: Restaurar Tests de Integración MongoDB

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo Original** | `bootstrap/factory_mongodb_integration_test.go.skip` |
| **Archivo Destino** | `bootstrap/factory_mongodb_integration_test.go` |
| **Líneas** | ~529 |

### Pasos de Implementación

#### Paso 2.1.1: Crear nuevo archivo de tests

Crear `bootstrap/factory_mongodb_integration_test.go`:

```go
//go:build integration
// +build integration

package bootstrap

import (
    "context"
    "testing"
    "time"

    "github.com/EduGoGroup/edugo-shared/testing/containers"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMongoDBFactory_CreateConnection_Success(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Setup container
    config := containers.NewConfig().
        WithMongoDB(nil).
        Build()

    manager, err := containers.GetManager(t, config)
    require.NoError(t, err, "Failed to get container manager")
    
    // Obtener MongoDB container
    mongoContainer := manager.MongoDB()
    require.NotNil(t, mongoContainer, "MongoDB container is nil")
    
    // Obtener connection string
    mongoURI, err := mongoContainer.ConnectionString(ctx)
    require.NoError(t, err, "Failed to get MongoDB connection string")
    require.NotEmpty(t, mongoURI, "MongoDB URI is empty")
    
    // Crear factory y probar conexión
    factory := NewDefaultMongoDBFactory()
    
    mongoConfig := MongoDBConfig{
        URI:         mongoURI,
        Database:    "test_integration_db",
        MaxPoolSize: 10,
        MinPoolSize: 1,
        Timeout:     10 * time.Second,
    }
    
    client, err := factory.CreateConnection(ctx, mongoConfig)
    require.NoError(t, err, "Failed to create MongoDB connection")
    require.NotNil(t, client, "MongoDB client is nil")
    
    // Cleanup
    defer func() {
        closeErr := factory.Close(ctx, client)
        assert.NoError(t, closeErr, "Failed to close MongoDB connection")
    }()
    
    // Verificar que podemos hacer ping
    err = factory.Ping(ctx, client)
    assert.NoError(t, err, "MongoDB ping failed")
}

func TestMongoDBFactory_GetDatabase(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Setup
    config := containers.NewConfig().
        WithMongoDB(nil).
        Build()

    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    
    mongoURI, err := manager.MongoDB().ConnectionString(ctx)
    require.NoError(t, err)
    
    factory := NewDefaultMongoDBFactory()
    client, err := factory.CreateConnection(ctx, MongoDBConfig{
        URI:      mongoURI,
        Database: "test_db",
        Timeout:  10 * time.Second,
    })
    require.NoError(t, err)
    defer factory.Close(ctx, client)
    
    // Test GetDatabase
    db := factory.GetDatabase(client, "test_database")
    assert.NotNil(t, db, "Database should not be nil")
    assert.Equal(t, "test_database", db.Name())
}

func TestMongoDBFactory_Integration_FullWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    // Setup
    config := containers.NewConfig().
        WithMongoDB(nil).
        Build()

    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    
    mongoURI, err := manager.MongoDB().ConnectionString(ctx)
    require.NoError(t, err)
    
    factory := NewDefaultMongoDBFactory()
    
    // 1. Create connection
    client, err := factory.CreateConnection(ctx, MongoDBConfig{
        URI:         mongoURI,
        Database:    "integration_test",
        MaxPoolSize: 10,
        MinPoolSize: 1,
        Timeout:     10 * time.Second,
    })
    require.NoError(t, err)
    defer factory.Close(ctx, client)
    
    // 2. Get database
    db := factory.GetDatabase(client, "integration_test")
    require.NotNil(t, db)
    
    // 3. Create collection and insert document
    collection := db.Collection("test_collection")
    _, err = collection.InsertOne(ctx, map[string]interface{}{
        "name":      "test_document",
        "timestamp": time.Now(),
    })
    assert.NoError(t, err)
    
    // 4. Query document
    var result map[string]interface{}
    err = collection.FindOne(ctx, map[string]interface{}{
        "name": "test_document",
    }).Decode(&result)
    assert.NoError(t, err)
    assert.Equal(t, "test_document", result["name"])
    
    // 5. Cleanup
    err = collection.Drop(ctx)
    assert.NoError(t, err)
}
```

#### Paso 2.1.2: Eliminar archivo .skip

```bash
rm bootstrap/factory_mongodb_integration_test.go.skip
```

#### Paso 2.1.3: Ejecutar tests

```bash
cd bootstrap && go test -v -tags=integration -run MongoDB
```

### Criterios de Éxito

- [ ] Archivo `.skip` eliminado
- [ ] Nuevo archivo `.go` creado
- [ ] Tests pasan con `go test -tags=integration`
- [ ] No hay errores de compilación

### Commit

```bash
git add bootstrap/factory_mongodb_integration_test.go
git rm bootstrap/factory_mongodb_integration_test.go.skip
git commit -m "test(bootstrap): restaurar tests de integración MongoDB

- Actualiza uso de containers API
- Agrega manejo de contexto en todas las llamadas
- Tests funcionando con testcontainers"
```

---

## Paso 2.2: Restaurar Tests de Integración PostgreSQL

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo Original** | `bootstrap/factory_postgresql_integration_test.go.skip` |
| **Archivo Destino** | `bootstrap/factory_postgresql_integration_test.go` |

### Pasos de Implementación

#### Paso 2.2.1: Crear nuevo archivo de tests

```go
//go:build integration
// +build integration

package bootstrap

import (
    "context"
    "testing"
    "time"

    "github.com/EduGoGroup/edugo-shared/testing/containers"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPostgreSQLFactory_CreateConnection_Success(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Setup container
    config := containers.NewConfig().
        WithPostgreSQL(nil).
        Build()

    manager, err := containers.GetManager(t, config)
    require.NoError(t, err, "Failed to get container manager")
    
    // Obtener PostgreSQL container
    pgContainer := manager.PostgreSQL()
    require.NotNil(t, pgContainer, "PostgreSQL container is nil")
    
    // Obtener configuración
    pgConfig, err := pgContainer.Config(ctx)
    require.NoError(t, err, "Failed to get PostgreSQL config")
    
    // Crear factory y probar
    factory := NewDefaultPostgreSQLFactory()
    
    db, err := factory.CreateConnection(ctx, PostgreSQLConfig{
        Host:               pgConfig.Host,
        Port:               pgConfig.Port,
        User:               pgConfig.User,
        Password:           pgConfig.Password,
        Database:           pgConfig.Database,
        SSLMode:            "disable",
        MaxConnections:     10,
        MaxIdleConnections: 5,
        MaxLifetime:        time.Minute * 5,
    })
    require.NoError(t, err, "Failed to create PostgreSQL connection")
    require.NotNil(t, db, "PostgreSQL DB is nil")
    
    defer factory.Close(db)
    
    // Verificar ping
    err = factory.Ping(ctx, db)
    assert.NoError(t, err, "PostgreSQL ping failed")
}

func TestPostgreSQLFactory_Integration_FullWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    config := containers.NewConfig().
        WithPostgreSQL(nil).
        Build()

    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    
    pgContainer := manager.PostgreSQL()
    pgConfig, err := pgContainer.Config(ctx)
    require.NoError(t, err)
    
    factory := NewDefaultPostgreSQLFactory()
    
    // 1. Create connection
    db, err := factory.CreateConnection(ctx, PostgreSQLConfig{
        Host:               pgConfig.Host,
        Port:               pgConfig.Port,
        User:               pgConfig.User,
        Password:           pgConfig.Password,
        Database:           pgConfig.Database,
        SSLMode:            "disable",
        MaxConnections:     10,
        MaxIdleConnections: 5,
        MaxLifetime:        time.Minute * 5,
    })
    require.NoError(t, err)
    defer factory.Close(db)
    
    // 2. Create table
    err = db.Exec(`CREATE TABLE IF NOT EXISTS test_table (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`).Error
    require.NoError(t, err)
    
    // 3. Insert data
    err = db.Exec(`INSERT INTO test_table (name) VALUES (?)`, "test_record").Error
    assert.NoError(t, err)
    
    // 4. Query data
    var count int64
    err = db.Raw(`SELECT COUNT(*) FROM test_table WHERE name = ?`, "test_record").Scan(&count).Error
    assert.NoError(t, err)
    assert.Equal(t, int64(1), count)
    
    // 5. Cleanup
    err = db.Exec(`DROP TABLE IF EXISTS test_table`).Error
    assert.NoError(t, err)
}
```

#### Paso 2.2.2: Eliminar archivo .skip y ejecutar tests

```bash
rm bootstrap/factory_postgresql_integration_test.go.skip
cd bootstrap && go test -v -tags=integration -run PostgreSQL
```

### Commit

```bash
git add bootstrap/factory_postgresql_integration_test.go
git rm bootstrap/factory_postgresql_integration_test.go.skip
git commit -m "test(bootstrap): restaurar tests de integración PostgreSQL

- Actualiza uso de containers API
- Agrega manejo de contexto
- Tests funcionando con testcontainers"
```

---

## Paso 2.3: Restaurar Tests de Integración RabbitMQ

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo Original** | `bootstrap/factory_rabbitmq_integration_test.go.skip` |
| **Archivo Destino** | `bootstrap/factory_rabbitmq_integration_test.go` |

### Pasos de Implementación

#### Paso 2.3.1: Crear nuevo archivo de tests

```go
//go:build integration
// +build integration

package bootstrap

import (
    "context"
    "testing"
    "time"

    "github.com/EduGoGroup/edugo-shared/testing/containers"
    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestRabbitMQFactory_CreateConnection_Success(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    config := containers.NewConfig().
        WithRabbitMQ(nil).
        Build()

    manager, err := containers.GetManager(t, config)
    require.NoError(t, err, "Failed to get container manager")
    
    rmqContainer := manager.RabbitMQ()
    require.NotNil(t, rmqContainer, "RabbitMQ container is nil")
    
    // Obtener URL
    rmqURL, err := rmqContainer.ConnectionString(ctx)
    require.NoError(t, err, "Failed to get RabbitMQ connection string")
    require.NotEmpty(t, rmqURL, "RabbitMQ URL is empty")
    
    factory := NewDefaultRabbitMQFactory()
    
    conn, err := factory.CreateConnection(ctx, RabbitMQConfig{
        URL: rmqURL,
    })
    require.NoError(t, err, "Failed to create RabbitMQ connection")
    require.NotNil(t, conn, "RabbitMQ connection is nil")
    
    defer factory.Close(nil, conn)
}

func TestRabbitMQFactory_Integration_FullWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    config := containers.NewConfig().
        WithRabbitMQ(nil).
        Build()

    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    
    rmqURL, err := manager.RabbitMQ().ConnectionString(ctx)
    require.NoError(t, err)
    
    factory := NewDefaultRabbitMQFactory()
    
    // 1. Create connection
    conn, err := factory.CreateConnection(ctx, RabbitMQConfig{URL: rmqURL})
    require.NoError(t, err)
    defer factory.Close(nil, conn)
    
    // 2. Create channel
    channel, err := factory.CreateChannel(conn)
    require.NoError(t, err)
    defer channel.Close()
    
    // 3. Declare queue
    queueName := "integration_test_queue"
    queue, err := factory.DeclareQueue(channel, queueName)
    require.NoError(t, err)
    require.Equal(t, queueName, queue.Name)
    
    // 4. Publish message
    testMessage := []byte(`{"test": "message"}`)
    err = channel.Publish(
        "",        // exchange
        queueName, // routing key
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        testMessage,
        },
    )
    assert.NoError(t, err)
    
    // 5. Consume message
    msgs, err := channel.Consume(
        queueName,
        "",    // consumer
        true,  // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    require.NoError(t, err)
    
    select {
    case msg := <-msgs:
        assert.Equal(t, testMessage, msg.Body)
    case <-time.After(5 * time.Second):
        t.Fatal("Timeout waiting for message")
    }
    
    // 6. Cleanup - delete queue
    _, err = channel.QueueDelete(queueName, false, false, false)
    assert.NoError(t, err)
}
```

#### Paso 2.3.2: Eliminar archivo .skip y ejecutar tests

```bash
rm bootstrap/factory_rabbitmq_integration_test.go.skip
cd bootstrap && go test -v -tags=integration -run RabbitMQ
```

### Commit

```bash
git add bootstrap/factory_rabbitmq_integration_test.go
git rm bootstrap/factory_rabbitmq_integration_test.go.skip
git commit -m "test(bootstrap): restaurar tests de integración RabbitMQ

- Actualiza uso de containers API
- Agrega manejo de contexto
- Tests de workflow completo (connect, channel, queue, pub/sub)"
```

---

## Paso 2.4: Verificar Coverage

### Verificación de Coverage

```bash
# Coverage de todo el módulo bootstrap
cd bootstrap
go test -coverprofile=coverage.out -tags=integration ./...
go tool cover -func=coverage.out

# Ver reporte HTML
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # macOS
```

### Objetivo de Coverage

| Módulo | Antes | Objetivo |
|--------|-------|----------|
| bootstrap | ~65% | >= 80% |

### Criterios de Éxito

- [ ] Coverage de bootstrap >= 80%
- [ ] No hay archivos `.skip` restantes
- [ ] Todos los tests de integración pasan

---

## Verificación Final de Fase 2

### Antes de Crear el PR

```bash
# Verificar que no hay archivos .skip
find . -name "*.skip" -type f
# Debería retornar vacío

# Ejecutar todos los tests de integración
cd bootstrap && go test -v -tags=integration ./...

# Verificar coverage
make test-coverage

# Build limpio
make build
```

### Crear Pull Request

```bash
# Push de la rama
git push origin fase-2-restauracion-tests

# En GitHub:
# 1. Crear PR hacia dev
# 2. Título: "test: Fase 2 - Restauración de Tests"
# 3. Descripción con lista de tests restaurados
```

### Revisión de GitHub Copilot

| Tipo de Comentario | Acción |
|-------------------|--------|
| Traducción inglés/español | DESCARTAR |
| Error de lógica en tests | CORREGIR |
| Sugerencia de mejora menor | DOCUMENTAR como deuda futura |

### Esperar Pipelines

```bash
# Revisar cada minuto durante máximo 10 minutos
# Los tests de integración pueden tardar más
# Si hay errores:
#   1. Analizar causa
#   2. Corregir (máx 3 intentos)
#   3. Push y esperar nuevamente
```

### Criterios de Éxito de Fase

- [ ] Todos los archivos `.skip` eliminados
- [ ] Todos los tests de integración pasan
- [ ] Coverage >= 80%
- [ ] `make build` compila sin errores
- [ ] PR aprobado
- [ ] Pipelines verdes
- [ ] Merge a dev completado

---

## Resumen de la Fase 2

| Paso | Descripción | Commit |
|------|-------------|--------|
| 2.1 | Restaurar tests MongoDB | `test(bootstrap): restaurar tests MongoDB` |
| 2.2 | Restaurar tests PostgreSQL | `test(bootstrap): restaurar tests PostgreSQL` |
| 2.3 | Restaurar tests RabbitMQ | `test(bootstrap): restaurar tests RabbitMQ` |
| 2.4 | Verificar coverage >= 80% | N/A |

---

## Siguiente Fase

Después de completar esta fase y hacer merge a dev, continuar con:
→ [FASE-3: Refactoring Estructural](./FASE-3_REFACTORING.md)
