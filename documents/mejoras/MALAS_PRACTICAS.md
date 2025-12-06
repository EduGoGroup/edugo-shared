# Malas Prácticas de Código

> Código que funciona pero viola best practices y debería ser corregido.

---

## Issue #1: Ignorar Errores de Ack/Nack en Consumer

### Ubicación
```
messaging/rabbit/consumer.go:67-70
```

### Código Actual
```go
// Manejar acknowledgment si no es auto-ack
if !c.config.AutoAck {
    if err != nil {
        // Nack con requeue si hubo error
        _ = msg.Nack(false, true) // Ignore nack errors
    } else {
        // Ack si fue exitoso
        _ = msg.Ack(false) // Ignore ack errors
    }
}
```

### Problema
- Los errores de `Ack()` y `Nack()` se ignoran con `_ =`
- Si el Ack falla, el mensaje podría procesarse dos veces
- Si el Nack falla, el mensaje podría perderse
- No hay logging ni métricas de estos fallos
- Dificulta debugging de problemas de mensajería

### Impacto
- **Alto en producción**: Mensajes pueden perderse o duplicarse silenciosamente
- **Debugging difícil**: No hay forma de saber si hay problemas de acknowledgment

### Solución Sugerida
```go
// Manejar acknowledgment si no es auto-ack
if !c.config.AutoAck {
    if err != nil {
        // Nack con requeue si hubo error
        if nackErr := msg.Nack(false, true); nackErr != nil {
            c.logger.Error("failed to nack message",
                "error", nackErr,
                "original_error", err,
                "delivery_tag", msg.DeliveryTag,
            )
            // Considerar: métricas, alertas, circuit breaker
        }
    } else {
        // Ack si fue exitoso
        if ackErr := msg.Ack(false); ackErr != nil {
            c.logger.Error("failed to ack message",
                "error", ackErr,
                "delivery_tag", msg.DeliveryTag,
            )
            // El mensaje podría reprocessarse - registrar para idempotencia
        }
    }
}
```

### Cambios Adicionales
1. Agregar logger al `RabbitMQConsumer`
2. Agregar métricas de Ack/Nack fallidos
3. Considerar retry con backoff para Ack fallidos

### Prioridad: **ALTA**

---

## Issue #2: Import Comentado sin Eliminar

### Ubicación
```
messaging/rabbit/consumer.go:7
```

### Código Actual
```go
import (
    "context"
    "encoding/json"
    "fmt"
    // amqp "github.com/rabbitmq/amqp091-go" // No usado actualmente
)
```

### Problema
- Import comentado indica código incompleto o removido
- Ensucian el código y confunden
- Debería eliminarse o usarse

### Impacto
- **Bajo**: Solo afecta legibilidad

### Solución
Eliminar la línea comentada:
```go
import (
    "context"
    "encoding/json"
    "fmt"
)
```

### Prioridad: **BAJA**

---

## Issue #3: Uso de `interface{}` en lugar de Tipos Específicos

### Ubicación
```
bootstrap/resource_implementations.go:72-74
```

### Código Actual
```go
type defaultStorageClient struct {
    client        *s3.Client
    presignClient interface{}  // <-- Debería ser *s3.PresignClient
    bucket        string
}
```

### Problema
- `interface{}` pierde type safety
- Requiere type assertion en runtime
- No hay validación en compile-time
- Dificulta el autocompletado del IDE

### Impacto
- **Medio**: Bugs potenciales en runtime que podrían detectarse en compile-time

### Solución
```go
import "github.com/aws/aws-sdk-go-v2/service/s3"

type defaultStorageClient struct {
    client        *s3.Client
    presignClient *s3.PresignClient
    bucket        string
}
```

### Cambios Requeridos
1. Cambiar tipo de `presignClient`
2. Actualizar `factory_s3.go` para retornar tipo correcto
3. Actualizar `S3Factory` interface si es necesario

### Prioridad: **MEDIA**

---

## Issue #4: Error Handling en Exists() Silencia Todos los Errores

### Ubicación
```
bootstrap/resource_implementations.go:133-145
```

### Código Actual
```go
func (c *defaultStorageClient) Exists(ctx context.Context, key string) (bool, error) {
    _, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
        Bucket: aws.String(c.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        // Si es error "not found", retornar false sin error
        return false, nil  // <-- Silencia TODOS los errores
    }
    return true, nil
}
```

### Problema
- Cualquier error (network, permisos, etc.) se trata como "no existe"
- Un problema de red haría que todos los archivos "no existan"
- Imposible diferenciar entre "no existe" y "error de acceso"

### Impacto
- **Alto**: Bugs silenciosos en producción

### Solución
```go
import (
    "errors"
    "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *defaultStorageClient) Exists(ctx context.Context, key string) (bool, error) {
    _, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
        Bucket: aws.String(c.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        // Verificar si es específicamente "not found"
        var notFound *types.NotFound
        if errors.As(err, &notFound) {
            return false, nil
        }
        // Cualquier otro error se propaga
        return false, fmt.Errorf("failed to check object existence: %w", err)
    }
    return true, nil
}
```

### Prioridad: **ALTA**

---

## Issue #5: Goroutine sin Control en Consumer

### Ubicación
```
messaging/rabbit/consumer.go:49-75
```

### Código Actual
```go
func (c *RabbitMQConsumer) Consume(...) error {
    msgs, err := c.conn.GetChannel().Consume(...)
    if err != nil {
        return err
    }

    // Procesar mensajes en un loop
    go func() {  // <-- Goroutine sin tracking
        for {
            select {
            case <-ctx.Done():
                return
            case msg, ok := <-msgs:
                // ...
            }
        }
    }()

    return nil  // <-- Retorna antes de que la goroutine termine
}
```

### Problema
- La goroutine no tiene forma de reportar errores después de iniciar
- No hay WaitGroup para esperar su terminación
- El caller no sabe si el consumer sigue funcionando
- Memory leaks potenciales si el consumer falla silenciosamente

### Impacto
- **Medio-Alto**: Problemas difíciles de diagnosticar en producción

### Solución
```go
type RabbitMQConsumer struct {
    conn     *Connection
    config   ConsumerConfig
    wg       sync.WaitGroup
    errChan  chan error
    stopOnce sync.Once
}

func (c *RabbitMQConsumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
    msgs, err := c.conn.GetChannel().Consume(...)
    if err != nil {
        return err
    }

    c.errChan = make(chan error, 1)
    c.wg.Add(1)
    
    go func() {
        defer c.wg.Done()
        for {
            select {
            case <-ctx.Done():
                return
            case msg, ok := <-msgs:
                if !ok {
                    c.errChan <- fmt.Errorf("message channel closed unexpectedly")
                    return
                }
                // procesar...
            }
        }
    }()

    return nil
}

func (c *RabbitMQConsumer) Wait() error {
    c.wg.Wait()
    select {
    case err := <-c.errChan:
        return err
    default:
        return nil
    }
}

func (c *RabbitMQConsumer) Errors() <-chan error {
    return c.errChan
}
```

### Prioridad: **MEDIA**

---

## Resumen de Malas Prácticas

| Archivo | Línea | Problema | Prioridad |
|---------|-------|----------|-----------|
| `consumer.go` | 67-70 | Ignorar errores Ack/Nack | Alta |
| `consumer.go` | 7 | Import comentado | Baja |
| `resource_implementations.go` | 73 | `interface{}` en lugar de tipo | Media |
| `resource_implementations.go` | 133-145 | Silenciar errores en Exists | Alta |
| `consumer.go` | 49-75 | Goroutine sin control | Media |

---

## Checklist de Corrección

- [ ] Manejar errores de Ack/Nack correctamente
- [ ] Eliminar import comentado
- [ ] Cambiar `interface{}` a `*s3.PresignClient`
- [ ] Detectar tipo de error específico en `Exists()`
- [ ] Agregar control de goroutine en Consumer
- [ ] Agregar tests para nuevos casos de error
