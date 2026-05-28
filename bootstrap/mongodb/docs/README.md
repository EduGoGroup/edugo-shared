# Bootstrap MongoDB — Documentacion tecnica

## Descripcion general

Sub-modulo que implementa la creacion de conexiones MongoDB usando mongo-driver v2 con pool configurado, timeouts y verificacion automatica de conectividad.

## Componentes principales

### Factory

```go
type Factory struct {
    connectionTimeout time.Duration // default: 10s
}
```

### CreateConnection

```go
func (f *Factory) CreateConnection(ctx context.Context, cfg bootstrap.MongoDBConfig) (*mongo.Client, error)
```

Configuracion del pool:
- MaxPoolSize: 100
- MinPoolSize: 10
- MaxConnIdleTime: 30 minutos
- ServerSelectionTimeout: 5 segundos
- ConnectTimeout: 10 segundos

Ejecuta Ping con `readpref.Primary()` antes de retornar. Si el ping falla, desconecta el cliente y retorna error.

### GetDatabase

```go
func (f *Factory) GetDatabase(client *mongo.Client, dbName string) *mongo.Database
```

Wrapper simple sobre `client.Database(dbName)`.

## Flujos comunes

### 1. Conexion con graceful degradation (Mobile pattern)

```go
if cfg.Database.MongoDB.URI != "" {
    mongoFactory := mongobootstrap.NewFactory()
    mongoClient, err := mongoFactory.CreateConnection(ctx, bootstrap.MongoDBConfig{
        URI:      cfg.Database.MongoDB.URI,
        Database: cfg.Database.MongoDB.Database,
    })
    if err != nil {
        log.Warn("MongoDB unavailable, continuing without MongoDB", "error", err)
        mongoClient = nil
    } else {
        mongoDB = mongoFactory.GetDatabase(mongoClient, cfg.Database.MongoDB.Database)
    }
}
```

### 2. Conexion requerida (Worker pattern)

```go
mongoFactory := mongobootstrap.NewFactory()
client, err := mongoFactory.CreateConnection(ctx, bootstrap.MongoDBConfig{
    URI:      cfg.Database.MongoDB.URI,
    Database: cfg.Database.MongoDB.Database,
})
if err != nil {
    return fmt.Errorf("failed to connect to MongoDB: %w", err)
}
defer mongoFactory.Close(ctx, client)
```

## Dependencias

### Internas
- `github.com/EduGoGroup/edugo-shared/bootstrap` — MongoDBConfig

### Externas
- `go.mongodb.org/mongo-driver/v2` — Driver MongoDB oficial

## Notas de diseño

- **Pool conservador**: 100 max / 10 min es suficiente para la mayoria de workloads. El driver gestiona el scaling automaticamente.
- **Ping obligatorio**: CreateConnection siempre verifica conectividad. Si MongoDB no esta disponible, el consumidor decide si continuar (graceful degradation) o fallar.
- **Timeout de 10s**: Suficiente para cold starts en entornos cloud. Evita bloquear indefinidamente.
