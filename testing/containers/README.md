# Testing Containers

Infraestructura de testing con testcontainers para PostgreSQL, MongoDB y RabbitMQ.

Este paquete proporciona containers de Docker listos para usar en tests de integración, implementando el patrón Singleton para reutilizar containers entre tests y mejorar el rendimiento.

## Características

- ✅ **Singleton Pattern**: Los containers se crean una sola vez y se reutilizan
- ✅ **Configuración Flexible**: Habilita solo los containers que necesitas
- ✅ **Cleanup Automático**: Gestión automática del ciclo de vida de containers
- ✅ **Métodos de Limpieza**: Limpia datos entre tests sin recrear containers
- ✅ **Thread-Safe**: Seguro para uso concurrente

## Uso Básico

### Configuración Simple

```go
func TestExample(t *testing.T) {
    ctx := context.Background()
    
    // Crear configuración con todos los containers
    config := containers.NewConfig().
        WithPostgreSQL(nil).
        WithMongoDB(nil).
        WithRabbitMQ(nil).
        Build()
    
    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    defer manager.Cleanup(ctx)
    
    // Usar los containers...
}
```

### PostgreSQL

```go
// Obtener configuración de conexión
pgConfig, err := manager.PostgreSQL().Config(ctx)
require.NoError(t, err)

fmt.Println("Host:", pgConfig.Host)
fmt.Println("Port:", pgConfig.Port)
fmt.Println("Database:", pgConfig.Database)
fmt.Println("User:", pgConfig.User)
fmt.Println("Password:", pgConfig.Password)

// O obtener connection string directamente
connString, err := manager.PostgreSQL().ConnectionString(ctx)
require.NoError(t, err)
// connString: "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable"
```

### MongoDB

```go
// Obtener URI de conexión
mongoURI, err := manager.MongoDB().ConnectionString(ctx)
require.NoError(t, err)
// mongoURI: "mongodb://localhost:27017/testdb"
```

### RabbitMQ

```go
// Obtener URL AMQP
rmqURL, err := manager.RabbitMQ().ConnectionString(ctx)
require.NoError(t, err)
// rmqURL: "amqp://guest:guest@localhost:5672/"
```

## API Completa

### Manager

#### GetManager

```go
func GetManager(t *testing.T, config *Config) (*Manager, error)
```

Obtiene o crea el singleton del manager. Si es la primera llamada, crea todos los containers según la configuración. Las llamadas subsiguientes retornan el mismo manager.

**Parámetros:**
- `t`: *testing.T para logging (puede ser nil)
- `config`: Configuración de containers a crear

**Retorna:**
- `*Manager`: Manager con containers listos
- `error`: Error si falla la creación de algún container

#### Cleanup

```go
func (m *Manager) Cleanup(ctx context.Context) error
```

Limpia y termina todos los containers creados. Debe llamarse al final de los tests.

**Recomendación:** Usar con `defer` inmediatamente después de `GetManager`.

### Container Getters

#### PostgreSQL

```go
func (m *Manager) PostgreSQL() *PostgresContainer
```

Retorna el container de PostgreSQL. Retorna `nil` si PostgreSQL no fue habilitado en la configuración.

#### MongoDB

```go
func (m *Manager) MongoDB() *MongoDBContainer
```

Retorna el container de MongoDB. Retorna `nil` si MongoDB no fue habilitado en la configuración.

#### RabbitMQ

```go
func (m *Manager) RabbitMQ() *RabbitMQContainer
```

Retorna el container de RabbitMQ. Retorna `nil` si RabbitMQ no fue habilitado en la configuración.

### PostgreSQL Container

#### Config

```go
func (c *PostgresContainer) Config(ctx context.Context) (*PostgreSQLConfig, error)
```

Retorna la configuración completa de conexión a PostgreSQL.

**PostgreSQLConfig contiene:**
- `Host string`: Hostname del container
- `Port int`: Puerto expuesto
- `Database string`: Nombre de la base de datos
- `User string`: Usuario
- `Password string`: Contraseña
- `SSLMode string`: Modo SSL (típicamente "disable" en tests)

#### ConnectionString

```go
func (c *PostgresContainer) ConnectionString(ctx context.Context) (string, error)
```

Retorna el DSN (Data Source Name) para conectarse a PostgreSQL.

Formato: `postgres://user:password@host:port/database?sslmode=disable`

#### Truncate

```go
func (c *PostgresContainer) Truncate(ctx context.Context, tables ...string) error
```

Trunca las tablas especificadas. Útil para limpiar datos entre tests sin recrear el container.

```go
// Limpiar tablas específicas
err := manager.PostgreSQL().Truncate(ctx, "users", "orders", "products")
```

### MongoDB Container

#### ConnectionString

```go
func (c *MongoDBContainer) ConnectionString(ctx context.Context) (string, error)
```

Retorna la URI de conexión para MongoDB.

Formato: `mongodb://host:port/database`

#### DropAllCollections

```go
func (c *MongoDBContainer) DropAllCollections(ctx context.Context) error
```

Elimina todas las colecciones de la base de datos. Útil para limpiar entre tests.

```go
// Limpiar todas las colecciones
err := manager.MongoDB().DropAllCollections(ctx)
```

### RabbitMQ Container

#### ConnectionString

```go
func (c *RabbitMQContainer) ConnectionString(ctx context.Context) (string, error)
```

Retorna la URL AMQP de conexión para RabbitMQ.

Formato: `amqp://user:password@host:port/`

#### PurgeAll

```go
func (c *RabbitMQContainer) PurgeAll(ctx context.Context) error
```

Elimina todas las colas y exchanges. Útil para limpiar entre tests.

```go
// Limpiar todas las colas y exchanges
err := manager.RabbitMQ().PurgeAll(ctx)
```

## Configuración Avanzada

### Configurar Container Específico

Puedes personalizar cada container con opciones específicas:

```go
config := containers.NewConfig().
    WithPostgreSQL(&containers.PostgreSQLOptions{
        Database: "mydb",
        User:     "myuser",
        Password: "mypassword",
    }).
    WithMongoDB(&containers.MongoDBOptions{
        Database: "testmongo",
    }).
    WithRabbitMQ(&containers.RabbitMQOptions{
        User:     "admin",
        Password: "secret",
    }).
    Build()
```

### Habilitar Solo Containers Necesarios

No es necesario crear todos los containers si no los usas:

```go
// Solo PostgreSQL
config := containers.NewConfig().
    WithPostgreSQL(nil).
    Build()

// Solo MongoDB y RabbitMQ
config := containers.NewConfig().
    WithMongoDB(nil).
    WithRabbitMQ(nil).
    Build()
```

## Métodos de Limpieza del Manager

El Manager proporciona métodos convenientes para limpiar datos sin acceder directamente a los containers:

### CleanPostgreSQL

```go
func (m *Manager) CleanPostgreSQL(ctx context.Context, tables ...string) error
```

Trunca tablas de PostgreSQL.

```go
err := manager.CleanPostgreSQL(ctx, "users", "sessions")
```

### CleanMongoDB

```go
func (m *Manager) CleanMongoDB(ctx context.Context) error
```

Elimina todas las colecciones de MongoDB.

```go
err := manager.CleanMongoDB(ctx)
```

### PurgeRabbitMQ

```go
func (m *Manager) PurgeRabbitMQ(ctx context.Context) error
```

Purga todas las colas y exchanges de RabbitMQ.

```go
err := manager.PurgeRabbitMQ(ctx)
```

## Timeouts

Por defecto, los containers tienen un timeout de **30 segundos** para iniciar.

Puedes configurar timeouts personalizados usando el contexto:

```go
// Timeout de 60 segundos para inicialización
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

pgConfig, err := manager.PostgreSQL().Config(ctx)
if err != nil {
    // Timeout o error de conexión
}
```

## Patrón de Uso Recomendado

### En TestMain

```go
func TestMain(m *testing.M) {
    ctx := context.Background()
    
    config := containers.NewConfig().
        WithPostgreSQL(nil).
        WithMongoDB(nil).
        WithRabbitMQ(nil).
        Build()
    
    manager, err := containers.GetManager(nil, config)
    if err != nil {
        log.Fatal("Failed to create containers:", err)
    }
    
    // Ejecutar tests
    code := m.Run()
    
    // Cleanup al final
    if err := manager.Cleanup(ctx); err != nil {
        log.Printf("Cleanup error: %v", err)
    }
    
    os.Exit(code)
}
```

### En Tests Individuales

```go
func TestDatabaseOperation(t *testing.T) {
    ctx := context.Background()
    
    // Obtener manager (reutiliza containers si ya existen)
    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    
    // Limpiar datos antes del test
    err = manager.CleanPostgreSQL(ctx, "users")
    require.NoError(t, err)
    
    // Tu test...
}
```

## Troubleshooting

### Container no inicia

Si un container tarda mucho en iniciar:
1. Verifica que Docker Desktop esté corriendo
2. Aumenta el timeout del contexto
3. Revisa los logs del container

### Puerto ya en uso

Los containers usan puertos aleatorios para evitar conflictos. Si aún así hay problemas:
1. Limpia containers huérfanos: `docker ps -a`
2. Reinicia Docker Desktop

### Tests lentos

Si los tests son lentos:
1. Asegúrate de usar el patrón Singleton (GetManager reutiliza containers)
2. Usa métodos de limpieza en lugar de recrear containers
3. Considera ejecutar tests con `-short` flag

## Ejemplos Completos

### Test con PostgreSQL

```go
func TestUserRepository(t *testing.T) {
    ctx := context.Background()
    
    config := containers.NewConfig().WithPostgreSQL(nil).Build()
    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    
    // Obtener connection string
    connStr, err := manager.PostgreSQL().ConnectionString(ctx)
    require.NoError(t, err)
    
    // Conectar a la base de datos
    db, err := sql.Open("postgres", connStr)
    require.NoError(t, err)
    defer db.Close()
    
    // Crear tabla
    _, err = db.ExecContext(ctx, `
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100)
        )
    `)
    require.NoError(t, err)
    
    // Tu test...
    
    // Limpiar al final
    err = manager.CleanPostgreSQL(ctx, "users")
    require.NoError(t, err)
}
```

### Test con MongoDB

```go
func TestMongoRepository(t *testing.T) {
    ctx := context.Background()
    
    config := containers.NewConfig().WithMongoDB(nil).Build()
    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    
    // Obtener URI
    mongoURI, err := manager.MongoDB().ConnectionString(ctx)
    require.NoError(t, err)
    
    // Conectar a MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    require.NoError(t, err)
    defer client.Disconnect(ctx)
    
    // Tu test...
    
    // Limpiar al final
    err = manager.CleanMongoDB(ctx)
    require.NoError(t, err)
}
```

### Test con RabbitMQ

```go
func TestMessagePublisher(t *testing.T) {
    ctx := context.Background()
    
    config := containers.NewConfig().WithRabbitMQ(nil).Build()
    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    
    // Obtener URL
    rmqURL, err := manager.RabbitMQ().ConnectionString(ctx)
    require.NoError(t, err)
    
    // Conectar a RabbitMQ
    conn, err := amqp.Dial(rmqURL)
    require.NoError(t, err)
    defer conn.Close()
    
    // Tu test...
    
    // Limpiar al final
    err = manager.PurgeRabbitMQ(ctx)
    require.NoError(t, err)
}
```

## Referencias

- [Testcontainers Go](https://golang.testcontainers.org/)
- [Testing en Go](https://go.dev/doc/tutorial/add-a-test)
- [Documentación del Proyecto](../../documents/TESTING.md)
