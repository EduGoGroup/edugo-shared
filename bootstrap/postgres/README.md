# Bootstrap PostgreSQL

Factory para crear conexiones PostgreSQL con GORM, SimpleProtocol (PgBouncer/Neon) y SearchPath configurable.

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/bootstrap/postgres
```

## Uso rapido

```go
import (
    "github.com/EduGoGroup/edugo-shared/bootstrap"
    pgbootstrap "github.com/EduGoGroup/edugo-shared/bootstrap/postgres"
)

// Conexion GORM con SearchPath y SimpleProtocol
pgFactory := pgbootstrap.NewFactory()
db, err := pgFactory.CreateGORMConnection(ctx, bootstrap.PostgreSQLConfig{
    Host:       "localhost",
    Port:       5432,
    User:       "app",
    Password:   "secret",
    Database:   "mydb",
    SSLMode:    "disable",
    SearchPath: "auth,iam,public",
    MaxOpenConns: 25,
},
    bootstrap.WithGORMLogger(gormLogger),
)

// Conexion raw *sql.DB
sqlDB, err := pgFactory.CreateRawConnection(ctx, bootstrap.PostgreSQLConfig{...})
```

## API Publica

- `NewFactory() *Factory` — Crea una nueva factory de PostgreSQL.
- `CreateRawConnection(ctx, PostgreSQLConfig) (*sql.DB, error)` — Conexion SQL nativa via pgx con SimpleProtocol y SearchPath.
- `CreateGORMConnection(ctx, PostgreSQLConfig, ...GORMOption) (*gorm.DB, error)` — Conexion GORM completa.
- `Ping(ctx, *gorm.DB) error` — Verifica conectividad.
- `Close(*gorm.DB) error` — Cierra la conexion.

## Navegacion

- [Documentacion tecnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Comandos disponibles

```bash
make build     # Compilar el modulo
make test      # Ejecutar tests
make check     # Lint y validacion
```

## Dependencias

- `github.com/EduGoGroup/edugo-shared/bootstrap` — Config structs y GORM options
- `github.com/jackc/pgx/v5` — Driver PostgreSQL
- `gorm.io/gorm` — ORM
- `gorm.io/driver/postgres` — GORM PostgreSQL dialector
