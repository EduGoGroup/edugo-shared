# Bootstrap RabbitMQ — Documentacion tecnica

## Descripcion general

Sub-modulo que implementa la creacion de conexiones RabbitMQ usando amqp091-go con timeout, QoS y declaracion de colas con configuracion por defecto.

## Componentes principales

### Factory

```go
type Factory struct {
    connectionTimeout time.Duration // default: 10s
}
```

### CreateConnection

```go
func (f *Factory) CreateConnection(ctx context.Context, cfg bootstrap.RabbitMQConfig) (*amqp.Connection, error)
```

Conecta a RabbitMQ con timeout. Usa goroutine interna + select para manejar:
- Conexion exitosa
- Error de conexion
- Timeout (10 segundos)
- Cancelacion de contexto

### CreateChannel

```go
func (f *Factory) CreateChannel(conn *amqp.Connection) (*amqp.Channel, error)
```

Crea un canal con QoS:
- Prefetch count: 10
- Prefetch size: 0
- Global: false

Si falla el QoS, cierra el canal antes de retornar error.

### DeclareQueue

```go
func (f *Factory) DeclareQueue(channel *amqp.Channel, queueName string) (amqp.Queue, error)
```

Declara una cola con:
- Durable: true
- TTL: 1 hora (x-message-ttl)
- Max priority: 10 (x-max-priority)
- Queue mode: lazy (x-queue-mode)
- Dead letter exchange: vacio por defecto

### Close

```go
func (f *Factory) Close(channel *amqp.Channel, conn *amqp.Connection) error
```

Cierra canal y conexion. Ignora `amqp.ErrClosed` para evitar errores en shutdown.

## Flujos comunes

### Crear conexion, canal y cola

```go
factory := rabbitbootstrap.NewFactory()
conn, err := factory.CreateConnection(ctx, bootstrap.RabbitMQConfig{URL: rabbitURL})
ch, err := factory.CreateChannel(conn)
queue, err := factory.DeclareQueue(ch, "events")
defer factory.Close(ch, conn)
```

## Dependencias

### Internas
- `github.com/EduGoGroup/edugo-shared/bootstrap` — RabbitMQConfig

### Externas
- `github.com/rabbitmq/amqp091-go` — Cliente AMQP 0.9.1

## Notas de diseño

- **Timeout con goroutine**: `amqp.Dial` no acepta contexto, asi que usamos goroutine + select para timeout.
- **QoS prefetch 10**: Balance entre throughput y uso de memoria. El consumidor puede ajustar post-creacion.
- **Lazy mode por defecto**: Reduce uso de RAM al almacenar mensajes en disco. Recomendado para colas de larga duracion.
- **Ningun consumidor actual usa esta factory**: Worker y Mobile usan `rabbit.Connect()` del modulo messaging/rabbit. Esta factory existe para futura adopcion o uso directo de AMQP.
