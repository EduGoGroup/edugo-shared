# Bootstrap — Documentacion tecnica

## Descripcion general

Modulo raiz que define tipos compartidos para la inicializacion de recursos de infraestructura. Contiene config structs, functional options para GORM, interface de lifecycle y tipos de error. **No tiene dependencias externas**.

Las implementaciones concretas viven en sub-modulos:
- `bootstrap/postgres` — PostgreSQL + GORM
- `bootstrap/mongodb` — MongoDB
- `bootstrap/rabbitmq` — RabbitMQ
- `bootstrap/s3` — S3

## Arquitectura

```
bootstrap/                     # Configs, options, errores (0 deps)
├── config.go                  # PostgreSQLConfig, MongoDBConfig, RabbitMQConfig, S3Config
├── gorm_options.go            # GORMOption, WithGORMLogger, WithSimpleProtocol
├── lifecycle.go               # LifecycleManager interface
├── errors.go                  # ErrMissingFactory, ErrConnectionFailed
│
├── postgres/                  # Factory PostgreSQL + GORM
│   └── factory.go             # NewFactory, CreateGORMConnection, CreateRawConnection
├── mongodb/                   # Factory MongoDB
│   └── factory.go             # NewFactory, CreateConnection, GetDatabase
├── rabbitmq/                  # Factory RabbitMQ
│   └── factory.go             # NewFactory, CreateConnection, CreateChannel
└── s3/                        # Factory S3
    └── factory.go             # NewFactory, CreateClient, CreatePresignClient
```

## Componentes principales

### PostgreSQLConfig

```go
type PostgreSQLConfig struct {
    Host, User, Password, Database, SSLMode string
    Port                                     int
    SearchPath                               string        // "auth,iam,public"
    MaxOpenConns, MaxIdleConns               int           // defaults: 25, 5
    ConnMaxLifetime, ConnMaxIdleTime         time.Duration // defaults: 1h, 10m
}
```

### GORMOption

```go
type GORMOption func(*GORMOptions)

type GORMOptions struct {
    Logger         any  // acepta gorm.logger.Interface
    SimpleProtocol bool // default: true (PgBouncer/Neon)
    PrepareStmt    bool // default: false
}
```

Funciones disponibles:
- `WithGORMLogger(logger any)` — Configura logger GORM
- `WithSimpleProtocol(bool)` — Controla pgx SimpleProtocol
- `WithPrepareStmt(bool)` — Controla GORM PrepareStmt
- `ApplyGORMOptions(...GORMOption) GORMOptions` — Aplica options sobre defaults

### LifecycleManager

```go
type LifecycleManager interface {
    RegisterCleanup(name string, cleanup func() error)
}
```

Interface tipada para registro de cleanup functions. Reemplaza el parametro `lifecycleManager any` del modulo anterior.

### Errores

```go
type ErrMissingFactory struct{ Resource string }
type ErrConnectionFailed struct{ Resource string; Err error }
```

## Flujos comunes

### 1. API con GORM (IAM, Admin)

```go
pgFactory := pgbootstrap.NewFactory()
db, err := pgFactory.CreateGORMConnection(ctx, bootstrap.PostgreSQLConfig{
    Host: cfg.DB.Host, Port: cfg.DB.Port,
    User: cfg.DB.User, Password: cfg.DB.Password,
    Database: cfg.DB.Database, SSLMode: cfg.DB.SSLMode,
    SearchPath: "auth,iam,academic,public",
    MaxOpenConns: cfg.DB.MaxOpenConns,
},
    bootstrap.WithGORMLogger(customLogger),
)
```

### 2. Worker con raw SQL

```go
pgFactory := pgbootstrap.NewFactory()
sqlDB, err := pgFactory.CreateRawConnection(ctx, bootstrap.PostgreSQLConfig{...})
```

### 3. Mobile con PostgreSQL + MongoDB (graceful degradation)

```go
db, _ := pgFactory.CreateGORMConnection(ctx, bootstrap.PostgreSQLConfig{...})
mongoClient, err := mongoFactory.CreateConnection(ctx, bootstrap.MongoDBConfig{...})
if err != nil {
    log.Warn("MongoDB unavailable, continuing without it")
}
```

## Dependencias

### Modulo raiz
Ninguna dependencia externa. Solo stdlib de Go.

### Sub-modulos
Cada sub-modulo solo importa lo que necesita:
- postgres: pgx, gorm
- mongodb: mongo-driver
- rabbitmq: amqp091-go
- s3: aws-sdk-go-v2

## Notas de diseño

- **0 dependencias en raiz**: Configs, options y errores no necesitan librerias externas. Esto elimina el problema del "god module" donde IAM descargaba MongoDB/RabbitMQ/Docker solo por importar bootstrap.
- **GORMOption en raiz**: Para que consumidores configuren sin importar el sub-modulo directamente.
- **SimpleProtocol por defecto**: Todos los entornos usan PgBouncer o Neon.
- **SearchPath configurable**: Cada servicio tiene su propio set de schemas.
- **Pool desde config**: En vez de hardcoded 25/5, los valores vienen de la configuracion del servicio.

## Operacion

```bash
make build     # Compilar
make test      # Tests unitarios
make test-all  # Tests incluyendo integracion (requiere Docker)
make check     # Lint + vet + test + build
```
