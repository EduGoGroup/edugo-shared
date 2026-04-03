# Bootstrap RabbitMQ

Factory para crear conexiones RabbitMQ con timeout, QoS configurado y declaracion de colas.

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/bootstrap/rabbitmq
```

## Uso rapido

```go
import (
    "github.com/EduGoGroup/edugo-shared/bootstrap"
    rabbitbootstrap "github.com/EduGoGroup/edugo-shared/bootstrap/rabbitmq"
)

factory := rabbitbootstrap.NewFactory()

conn, err := factory.CreateConnection(ctx, bootstrap.RabbitMQConfig{
    URL: "amqp://guest:guest@localhost:5672/",
})

ch, err := factory.CreateChannel(conn)
queue, err := factory.DeclareQueue(ch, "my-queue")
defer factory.Close(ch, conn)
```

## API Publica

- `NewFactory() *Factory` — Crea una nueva factory de RabbitMQ.
- `CreateConnection(ctx, RabbitMQConfig) (*amqp.Connection, error)` — Conexion con timeout de 10s.
- `CreateChannel(*amqp.Connection) (*amqp.Channel, error)` — Canal con QoS prefetch 10.
- `DeclareQueue(*amqp.Channel, string) (amqp.Queue, error)` — Cola durable con TTL 1h, max priority 10, lazy mode.
- `Close(*amqp.Channel, *amqp.Connection) error` — Cierra canal y conexion.

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
- `github.com/rabbitmq/amqp091-go` — Cliente AMQP 0.9.1
