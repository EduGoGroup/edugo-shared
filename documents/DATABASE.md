# Diagrama de Base de Datos y Configuración

## Sistemas de Base de Datos

El proyecto utiliza un enfoque **polyglot persistence** con dos sistemas de base de datos:

| Sistema | Uso Principal | Driver |
|---------|---------------|--------|
| **PostgreSQL** | Datos relacionales transaccionales | GORM + lib/pq |
| **MongoDB** | Documentos y datos flexibles | mongo-driver |

---

## PostgreSQL

### Diagrama de Conexión

```
┌──────────────────────────────────────────────────────────────────┐
│                      POSTGRESQL CONNECTION                        │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌─────────────┐                      ┌─────────────────────┐    │
│  │  Service    │                      │    PostgreSQL       │    │
│  │  (GORM)     │                      │      Server         │    │
│  └──────┬──────┘                      │                     │    │
│         │                             │  ┌───────────────┐  │    │
│         │  1. Connect with DSN        │  │   Database    │  │    │
│         │ ──────────────────────────► │  │  edugo_test   │  │    │
│         │                             │  └───────────────┘  │    │
│         │  2. Connection Pool         │                     │    │
│         │ ◄────────────────────────── │  Port: 5432         │    │
│         │                             │  SSL: configurable  │    │
│  ┌──────┴──────┐                      └─────────────────────┘    │
│  │   Pool      │                                                  │
│  │ ┌────┬────┐ │                                                  │
│  │ │conn│conn│ │  MaxOpenConns: configurable                     │
│  │ ├────┼────┤ │  MaxIdleConns: configurable                     │
│  │ │conn│conn│ │  ConnMaxLifetime: configurable                  │
│  │ └────┴────┘ │                                                  │
│  └─────────────┘                                                  │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

### Configuración

```go
// database/postgres/config.go
type Config struct {
    Host               string        // Host del servidor (default: localhost)
    Port               int           // Puerto (default: 5432)
    User               string        // Usuario de conexión
    Password           string        // Contraseña
    Database           string        // Nombre de la base de datos
    SSLMode            string        // disable | require | verify-ca | verify-full
    MaxConnections     int           // Conexiones máximas abiertas
    MaxIdleConnections int           // Conexiones idle máximas
    MaxLifetime        time.Duration // Tiempo máximo de vida de conexión
    ConnectTimeout     time.Duration // Timeout de conexión inicial
}
```

### DSN Format

```
host={host} port={port} user={user} password={password} dbname={database} sslmode={sslmode} connect_timeout={seconds}
```

### Ejemplo de Uso

```go
import "github.com/EduGoGroup/edugo-shared/database/postgres"

// Crear configuración
cfg := postgres.NewConfig().
    WithHost("localhost").
    WithPort(5432).
    WithUser("edugo_user").
    WithPassword("secret").
    WithDatabase("edugo_db").
    WithSSLMode("disable").
    Build()

// Conectar
db, err := postgres.Connect(cfg)
if err != nil {
    log.Fatal(err)
}
defer postgres.Close(db)

// Health check
if err := postgres.HealthCheck(db); err != nil {
    log.Fatal("Database unhealthy:", err)
}
```

### Transacciones

```go
import "github.com/EduGoGroup/edugo-shared/database/postgres"

// Ejecutar en transacción
err := postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
    _, err := tx.Exec("INSERT INTO users (name) VALUES ($1)", "Alice")
    if err != nil {
        return err // Rollback automático
    }
    return nil // Commit automático
})
```

---

## MongoDB

### Diagrama de Conexión

```
┌──────────────────────────────────────────────────────────────────┐
│                       MONGODB CONNECTION                          │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌─────────────┐                      ┌─────────────────────┐    │
│  │  Service    │                      │      MongoDB        │    │
│  │  (Driver)   │                      │      Server         │    │
│  └──────┬──────┘                      │                     │    │
│         │                             │  ┌───────────────┐  │    │
│         │  1. Connect with URI        │  │   Database    │  │    │
│         │ ──────────────────────────► │  │  edugo_test   │  │    │
│         │                             │  │               │  │    │
│         │  2. Client Pool             │  │ ┌───────────┐ │  │    │
│         │ ◄────────────────────────── │  │ │Collection │ │  │    │
│         │                             │  │ │Collection │ │  │    │
│  ┌──────┴──────┐                      │  │ └───────────┘ │  │    │
│  │   Pool      │                      │  └───────────────┘  │    │
│  │ ┌────┬────┐ │                      │                     │    │
│  │ │conn│conn│ │  MaxPoolSize: 100    │  Port: 27017        │    │
│  │ ├────┼────┤ │  MinPoolSize: 10     │                     │    │
│  │ │conn│conn│ │                      └─────────────────────┘    │
│  │ └────┴────┘ │                                                  │
│  └─────────────┘                                                  │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

### Configuración

```go
// database/mongodb/config.go
type Config struct {
    URI         string        // MongoDB connection URI
    Database    string        // Nombre de la base de datos
    MaxPoolSize uint64        // Tamaño máximo del pool (default: 100)
    MinPoolSize uint64        // Tamaño mínimo del pool (default: 10)
    Timeout     time.Duration // Timeout de conexión (default: 30s)
}
```

### URI Format

```
mongodb://[user:password@]host:port/[database][?options]
```

### Ejemplo de Uso

```go
import "github.com/EduGoGroup/edugo-shared/database/mongodb"

// Crear configuración
cfg := mongodb.Config{
    URI:         "mongodb://localhost:27017",
    Database:    "edugo_db",
    MaxPoolSize: 100,
    MinPoolSize: 10,
    Timeout:     30 * time.Second,
}

// Conectar
client, err := mongodb.Connect(cfg)
if err != nil {
    log.Fatal(err)
}
defer mongodb.Close(client)

// Obtener database
db := mongodb.GetDatabase(client, cfg.Database)

// Obtener collection
users := db.Collection("users")

// Insertar documento
result, err := users.InsertOne(ctx, bson.M{
    "name":  "Alice",
    "email": "alice@edugo.com",
})
```

---

## Esquema de Configuración Centralizada

```yaml
# config.yaml ejemplo
database:
  host: localhost
  port: 5432
  user: edugo_user
  password: secret
  database: edugo_db
  ssl_mode: disable
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

mongodb:
  uri: mongodb://localhost:27017
  database: edugo_db
  max_pool_size: 100
  min_pool_size: 10
  connect_timeout: 30s
```

---

## Diagrama de Entidades Conceptual

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        ENTIDADES DEL DOMINIO EDUGO                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  POSTGRESQL (Datos Transaccionales)                                          │
│  ═══════════════════════════════════                                         │
│                                                                              │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐                │
│  │    Users     │     │  Materials   │     │   Progress   │                │
│  ├──────────────┤     ├──────────────┤     ├──────────────┤                │
│  │ id (UUID)    │────►│ id (UUID)    │────►│ id (UUID)    │                │
│  │ email        │     │ title        │     │ user_id (FK) │                │
│  │ role         │     │ type         │     │ material_id  │                │
│  │ created_at   │     │ status       │     │ status       │                │
│  └──────────────┘     │ created_by   │     │ percentage   │                │
│         │             └──────────────┘     │ completed_at │                │
│         │                    │             └──────────────┘                │
│         │                    │                                              │
│         ▼                    ▼                                              │
│  ┌──────────────┐     ┌──────────────┐                                     │
│  │  Enrollments │     │ Assessments  │                                     │
│  ├──────────────┤     ├──────────────┤                                     │
│  │ user_id (FK) │     │ material_id  │                                     │
│  │ course_id    │     │ type         │                                     │
│  │ enrolled_at  │     │ questions    │                                     │
│  └──────────────┘     └──────────────┘                                     │
│                                                                              │
│  MONGODB (Documentos Flexibles)                                              │
│  ══════════════════════════════                                              │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────┐              │
│  │                     audit_logs                            │              │
│  ├──────────────────────────────────────────────────────────┤              │
│  │  {                                                        │              │
│  │    "_id": ObjectId,                                       │              │
│  │    "user_id": "uuid-string",                              │              │
│  │    "action": "material.created",                          │              │
│  │    "resource_type": "material",                           │              │
│  │    "resource_id": "uuid-string",                          │              │
│  │    "metadata": { ... },                                   │              │
│  │    "timestamp": ISODate                                   │              │
│  │  }                                                        │              │
│  └──────────────────────────────────────────────────────────┘              │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────┐              │
│  │                 material_content                          │              │
│  ├──────────────────────────────────────────────────────────┤              │
│  │  {                                                        │              │
│  │    "_id": ObjectId,                                       │              │
│  │    "material_id": "uuid-string",                          │              │
│  │    "content": { rich_text_data },                         │              │
│  │    "attachments": [ { url, type, size } ],                │              │
│  │    "version": 1,                                          │              │
│  │    "updated_at": ISODate                                  │              │
│  │  }                                                        │              │
│  └──────────────────────────────────────────────────────────┘              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Health Checks

### PostgreSQL
```go
// Ejecuta: SELECT 1
err := db.PingContext(ctx)
```

### MongoDB
```go
// Ejecuta: ping command
err := client.Ping(ctx, readpref.Primary())
```

---

## Timeouts por Defecto

| Operación | Timeout |
|-----------|---------|
| PostgreSQL Connect | Configurable |
| PostgreSQL Health Check | 5 segundos |
| MongoDB Connect | 30 segundos |
| MongoDB Health Check | 5 segundos |
| MongoDB Disconnect | 10 segundos |
