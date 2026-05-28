# Bootstrap PostgreSQL — Documentacion tecnica

## Descripcion general

Sub-modulo que implementa la creacion de conexiones PostgreSQL usando pgx v5 como driver y GORM como ORM. Diseñado para ser compatible con PgBouncer (transaction mode) y Neon via SimpleProtocol.

## Componentes principales

### Factory

```go
type Factory struct{}
```

Struct sin estado que implementa los metodos de creacion de conexiones.

### CreateRawConnection

```go
func (f *Factory) CreateRawConnection(ctx context.Context, cfg bootstrap.PostgreSQLConfig) (*sql.DB, error)
```

Crea una conexion `*sql.DB` usando `pgx.ParseConfig` + `stdlib.OpenDB`. Configura:
- `search_path` en `RuntimeParams` si `cfg.SearchPath` no esta vacio
- Connection pool: MaxOpenConns (default 25), MaxIdleConns (default 5), ConnMaxLifetime (default 1h), ConnMaxIdleTime (default 10m)
- Ejecuta Ping para verificar conectividad antes de retornar

### CreateGORMConnection

```go
func (f *Factory) CreateGORMConnection(ctx context.Context, cfg bootstrap.PostgreSQLConfig, opts ...bootstrap.GORMOption) (*gorm.DB, error)
```

Crea una conexion `*gorm.DB` sobre pgx. Pasos:
1. Construye pgx config con DSN y SearchPath
2. Si `SimpleProtocol` esta habilitado (default: true), configura `pgx.QueryExecModeSimpleProtocol`
3. Abre `*sql.DB` via `stdlib.OpenDB` y configura pool
4. Ejecuta Ping
5. Abre GORM con `postgres.Config{Conn: sqlDB}` y las opciones aplicadas

### GORMOption (del modulo raiz bootstrap)

Opciones funcionales definidas en el modulo raiz para evitar dependencia circular:

- `WithGORMLogger(logger any)` — Acepta `gorm.logger.Interface`
- `WithSimpleProtocol(bool)` — Default: true (para PgBouncer/Neon)
- `WithPrepareStmt(bool)` — Default: false (incompatible con SimpleProtocol)

## PostgreSQLConfig (del modulo raiz bootstrap)

```go
type PostgreSQLConfig struct {
    Host, User, Password, Database, SSLMode string
    Port                                     int
    SearchPath                               string        // "auth,iam,public"
    MaxOpenConns, MaxIdleConns               int
    ConnMaxLifetime, ConnMaxIdleTime         time.Duration
}
```

## Flujos comunes

### 1. API con GORM y SearchPath

```go
db, err := pgFactory.CreateGORMConnection(ctx, bootstrap.PostgreSQLConfig{
    Host: cfg.DB.Host, Port: cfg.DB.Port,
    User: cfg.DB.User, Password: cfg.DB.Password,
    Database: cfg.DB.Database, SSLMode: cfg.DB.SSLMode,
    SearchPath: "auth,iam,academic,public",
    MaxOpenConns: cfg.DB.MaxOpenConns,
    MaxIdleConns: cfg.DB.MaxIdleConns,
    ConnMaxLifetime: time.Hour,
},
    bootstrap.WithGORMLogger(customLogger),
)
```

### 2. Worker con raw SQL

```go
sqlDB, err := pgFactory.CreateRawConnection(ctx, bootstrap.PostgreSQLConfig{
    Host: cfg.DB.Host, Port: cfg.DB.Port,
    User: cfg.DB.User, Password: cfg.DB.Password,
    Database: cfg.DB.Database, SSLMode: cfg.DB.SSLMode,
})
```

## Arquitectura

```
Consumer (IAM/Admin/Mobile/Worker)
    |
    v
pgbootstrap.NewFactory()
    |
    +-- CreateGORMConnection()
    |       |
    |       +-- pgx.ParseConfig(dsn)
    |       +-- SimpleProtocol + SearchPath
    |       +-- stdlib.OpenDB()
    |       +-- Pool config
    |       +-- Ping
    |       +-- gorm.Open(postgres.Config{Conn: sqlDB})
    |
    +-- CreateRawConnection()
            |
            +-- pgx.ParseConfig(dsn)
            +-- SearchPath
            +-- stdlib.OpenDB()
            +-- Pool config
            +-- Ping
```

## Dependencias

### Internas
- `github.com/EduGoGroup/edugo-shared/bootstrap` — Config structs, GORMOption

### Externas
- `github.com/jackc/pgx/v5` — Driver PostgreSQL con SimpleProtocol
- `gorm.io/gorm` — ORM
- `gorm.io/driver/postgres` — GORM dialector

## Notas de diseño

- **SimpleProtocol por defecto**: PgBouncer en transaction mode y Neon no soportan prepared statements. pgx v5 los cachea internamente, asi que `QueryExecModeSimpleProtocol` los deshabilita a nivel driver.
- **SearchPath en RuntimeParams**: Mas robusto que incluirlo en el DSN, ya que pgx lo maneja como parametro de conexion.
- **Pool defaults conservadores**: 25 open / 5 idle / 1h lifetime. Los consumidores pueden override via config.
- **No PrepareStmt**: Incompatible con SimpleProtocol. Se deshabilita en GORM config.
