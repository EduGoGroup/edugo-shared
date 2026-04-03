# Bootstrap MongoDB

Factory para crear conexiones MongoDB con pool configurado y verificacion automatica de conectividad.

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/bootstrap/mongodb
```

## Uso rapido

```go
import (
    "github.com/EduGoGroup/edugo-shared/bootstrap"
    mongobootstrap "github.com/EduGoGroup/edugo-shared/bootstrap/mongodb"
)

// Crear conexion
mongoFactory := mongobootstrap.NewFactory()
client, err := mongoFactory.CreateConnection(ctx, bootstrap.MongoDBConfig{
    URI:      "mongodb://localhost:27017",
    Database: "mydb",
})

// Obtener database
db := mongoFactory.GetDatabase(client, "mydb")

// Graceful degradation (patron usado por Mobile)
if err != nil {
    log.Warn("MongoDB unavailable, continuing without it", "error", err)
}
```

## API Publica

- `NewFactory() *Factory` — Crea una nueva factory de MongoDB.
- `CreateConnection(ctx, MongoDBConfig) (*mongo.Client, error)` — Conexion con pool y ping automatico.
- `GetDatabase(*mongo.Client, string) *mongo.Database` — Obtiene una database especifica.
- `Ping(ctx, *mongo.Client) error` — Verifica conectividad (5s timeout, read preference primaria).
- `Close(ctx, *mongo.Client) error` — Cierra la conexion (10s timeout).

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

- `github.com/EduGoGroup/edugo-shared/bootstrap` — Config structs
- `go.mongodb.org/mongo-driver/v2` — Driver MongoDB
