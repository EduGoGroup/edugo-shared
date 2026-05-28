# Bootstrap

Tipos compartidos para la inicializacion de recursos de infraestructura. Cada tecnologia tiene su propio sub-modulo.

## Arquitectura

```
bootstrap/                    # Modulo raiz: configs, options, errores (0 deps externas)
├── bootstrap/postgres/       # Factory PostgreSQL + GORM
├── bootstrap/mongodb/        # Factory MongoDB
├── bootstrap/rabbitmq/       # Factory RabbitMQ
└── bootstrap/s3/             # Factory S3
```

## Instalacion

```bash
# Solo config structs y options
go get github.com/EduGoGroup/edugo-shared/bootstrap

# PostgreSQL + GORM (lo mas comun)
go get github.com/EduGoGroup/edugo-shared/bootstrap/postgres

# MongoDB
go get github.com/EduGoGroup/edugo-shared/bootstrap/mongodb
```

## Uso rapido

```go
import (
    "github.com/EduGoGroup/edugo-shared/bootstrap"
    pgbootstrap "github.com/EduGoGroup/edugo-shared/bootstrap/postgres"
)

pgFactory := pgbootstrap.NewFactory()
db, err := pgFactory.CreateGORMConnection(ctx, bootstrap.PostgreSQLConfig{
    Host:       "localhost",
    Port:       5432,
    User:       "app",
    Password:   "secret",
    Database:   "mydb",
    SSLMode:    "disable",
    SearchPath: "auth,iam,public",
},
    bootstrap.WithGORMLogger(gormLogger),
)
```

## Modulo raiz — API Publica

### Config structs
- `PostgreSQLConfig` — Host, Port, User, Password, Database, SSLMode, SearchPath, pool config
- `MongoDBConfig` — URI, Database
- `RabbitMQConfig` — URL
- `S3Config` — Bucket, Region, AccessKeyID, SecretAccessKey, Endpoint, ForcePathStyle

### GORM options
- `WithGORMLogger(logger any)` — Logger para GORM
- `WithSimpleProtocol(bool)` — SimpleProtocol para PgBouncer/Neon (default: true)
- `WithPrepareStmt(bool)` — PrepareStmt para GORM (default: false)

### Interfaces
- `LifecycleManager` — Registro de cleanup functions para graceful shutdown

### Errores
- `ErrMissingFactory` — Factory requerida no encontrada
- `ErrConnectionFailed` — Fallo de conexion a un recurso

## Sub-modulos

| Sub-modulo | Funcion principal | Consumidores |
|-----------|------------------|-------------|
| [postgres](postgres/) | `CreateGORMConnection()` | IAM, Admin, Mobile, Worker |
| [mongodb](mongodb/) | `CreateConnection()` | Mobile, Worker |
| [rabbitmq](rabbitmq/) | `CreateConnection()` | (futura adopcion) |
| [s3](s3/) | `CreateClient()` | (futura adopcion) |

## Navegacion

- [Changelog](CHANGELOG.md)

## Comandos disponibles

```bash
make build     # Compilar el modulo
make test      # Ejecutar tests
make check     # Lint y validacion
```
