# FASE 3: Refactoring Estructural

> **Prioridad**: MEDIA  
> **Duración estimada**: 3-4 días  
> **Prerrequisitos**: Fase 1 y 2 completadas  
> **Rama**: `fase-3-refactoring-estructural`  
> **Objetivo**: Mejorar estructura del código para mejor mantenibilidad

---

## Flujo de Trabajo de Esta Fase

### 1. Inicio de la Fase

```bash
# Asegurarse de estar en dev actualizado
git checkout dev
git pull origin dev

# Crear rama de la fase
git checkout -b fase-3-refactoring-estructural

# Verificar estado inicial
make build
make test-all-modules
```

### 2. Durante la Fase

- Ejecutar cada paso en orden
- Commit atómico después de cada paso completado
- Verificar que tests pasen después de cada cambio
- Esta fase es más delicada: verificar regresiones constantemente

### 3. Fin de la Fase

```bash
# Push de la rama
git push origin fase-3-refactoring-estructural

# Crear PR en GitHub hacia dev
# - Título: "refactor: Fase 3 - Refactoring Estructural"
# - Descripción: Lista de refactorizaciones realizadas

# Esperar revisión de GitHub Copilot
# - DESCARTAR: Comentarios de traducción inglés/español
# - CORREGIR: Problemas importantes detectados
# - DOCUMENTAR: Lo que queda como deuda técnica futura

# Esperar pipelines (máx 10 min, revisar cada 1 min)
# - Si hay errores: Corregir (regla de 3 intentos)
# - Todos los errores se corrigen (propios o heredados)

# Merge cuando todo esté verde
```

---

## Resumen de la Fase

Esta fase reestructura el código para:
- Reducir tamaño de archivos grandes
- Eliminar duplicación de código
- Mejorar type safety
- Unificar implementaciones duplicadas

---

## Paso 3.1: Dividir bootstrap.go en Archivos Más Pequeños

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo** | `bootstrap/bootstrap.go` |
| **Líneas Actuales** | 623 |
| **Líneas Objetivo** | < 150 por archivo |

### Estructura Objetivo

```
bootstrap/
├── bootstrap.go           # Solo función principal Bootstrap() (~100 líneas)
├── init_logger.go         # initLogger (~50 líneas)
├── init_postgresql.go     # initPostgreSQL (~80 líneas)
├── init_mongodb.go        # initMongoDB (~80 líneas)
├── init_rabbitmq.go       # initRabbitMQ (~80 líneas)
├── init_s3.go             # initS3 (~80 líneas)
├── health_check.go        # performHealthChecks (~50 líneas)
├── config_extractors.go   # Funciones extract*Config (~100 líneas)
├── cleanup_registrars.go  # Funciones register*Cleanup (~80 líneas)
├── interfaces.go          # (ya existe)
├── options.go             # (ya existe)
└── resources.go           # (ya existe)
```

### Pasos de Implementación

#### Paso 3.1.1: Crear init_logger.go

```go
package bootstrap

import (
    "context"
)

// initLogger inicializa el logger para la aplicación.
//
// Parámetros:
//   - ctx: Contexto para cancelación
//   - factories: Fábricas disponibles
//   - config: Configuración de la aplicación
//   - resources: Recursos a inicializar
//
// Retorna error si la inicialización falla.
func initLogger(
    ctx context.Context,
    factories *Factories,
    config interface{},
    resources *Resources,
) error {
    if factories.Logger == nil {
        return nil
    }

    env, version := extractEnvAndVersion(config)
    
    logger, err := factories.Logger.CreateLogger(ctx, env, version)
    if err != nil {
        return err
    }

    resources.Logger = logger
    return nil
}
```

#### Paso 3.1.2: Crear init_postgresql.go

```go
package bootstrap

import (
    "context"
    "fmt"
)

// initPostgreSQL inicializa la conexión a PostgreSQL.
//
// Parámetros:
//   - ctx: Contexto para cancelación
//   - factories: Fábricas disponibles
//   - config: Configuración de la aplicación
//   - resources: Recursos a inicializar
//   - lm: Manager de lifecycle para cleanup
//   - options: Opciones de bootstrap
//
// Retorna error si el recurso es requerido y falla la inicialización.
func initPostgreSQL(
    ctx context.Context,
    factories *Factories,
    config interface{},
    resources *Resources,
    lm *lifecycle.Manager,
    options *bootstrapOptions,
) error {
    if factories.PostgreSQL == nil {
        if options.IsResourceRequired("postgresql") {
            return fmt.Errorf("PostgreSQL factory is required but not provided")
        }
        return nil
    }

    pgConfig, err := extractPostgreSQLConfig(config)
    if err != nil {
        if options.IsResourceRequired("postgresql") {
            return fmt.Errorf("failed to extract PostgreSQL config: %w", err)
        }
        logOptionalResourceSkipped(resources.Logger, "postgresql", err)
        return nil
    }

    db, err := factories.PostgreSQL.CreateConnection(ctx, pgConfig)
    if err != nil {
        if options.IsResourceRequired("postgresql") {
            return fmt.Errorf("failed to create PostgreSQL connection: %w", err)
        }
        logOptionalResourceSkipped(resources.Logger, "postgresql", err)
        return nil
    }

    resources.PostgreSQL = db
    registerPostgreSQLCleanup(lm, factories.PostgreSQL, db)

    return nil
}
```

#### Paso 3.1.3: Crear archivos restantes

Continuar con el mismo patrón para:
- `init_mongodb.go`
- `init_rabbitmq.go`
- `init_s3.go`
- `health_check.go`
- `config_extractors.go`
- `cleanup_registrars.go`

#### Paso 3.1.4: Simplificar bootstrap.go principal

```go
package bootstrap

import (
    "context"
    "fmt"
)

// Bootstrap inicializa todos los recursos de infraestructura de la aplicación.
//
// Parámetros:
//   - ctx: Contexto para cancelación y timeouts
//   - config: Struct de configuración con campos para cada recurso
//   - factories: Fábricas para crear recursos
//   - lm: Manager de lifecycle para cleanup ordenado
//   - opts: Opciones adicionales de configuración
//
// Retorna los recursos inicializados o error si falla algún recurso requerido.
func Bootstrap(
    ctx context.Context,
    config interface{},
    factories *Factories,
    lm *lifecycle.Manager,
    opts ...Option,
) (*Resources, error) {
    options := newBootstrapOptions(opts...)
    resources := &Resources{}

    // 1. Initialize Logger (always first)
    if err := initLogger(ctx, factories, config, resources); err != nil {
        return nil, fmt.Errorf("logger initialization failed: %w", err)
    }

    // 2. Initialize PostgreSQL
    if err := initPostgreSQL(ctx, factories, config, resources, lm, options); err != nil {
        return nil, err
    }

    // 3. Initialize MongoDB
    if err := initMongoDB(ctx, factories, config, resources, lm, options); err != nil {
        return nil, err
    }

    // 4. Initialize RabbitMQ
    if err := initRabbitMQ(ctx, factories, config, resources, lm, options); err != nil {
        return nil, err
    }

    // 5. Initialize S3
    if err := initS3(ctx, factories, config, resources, lm, options); err != nil {
        return nil, err
    }

    // 6. Health checks
    if !options.skipHealthCheck {
        if err := performHealthChecks(ctx, factories, resources, options); err != nil {
            return nil, err
        }
    }

    logBootstrapComplete(resources.Logger)
    return resources, nil
}
```

### Criterios de Éxito

- [ ] `bootstrap.go` tiene menos de 150 líneas
- [ ] Cada archivo init tiene menos de 100 líneas
- [ ] Todos los tests siguen pasando
- [ ] `go build ./...` compila sin errores

### Commit

```bash
git add bootstrap/
git commit -m "refactor(bootstrap): dividir bootstrap.go en archivos más pequeños

- Separa inicialización de cada recurso en archivo propio
- Extrae funciones de extracción de config
- Extrae funciones de cleanup
- Mejora legibilidad y mantenibilidad"
```

### Actualización de Documentación

Actualizar `documents/ARCHITECTURE.md` para reflejar la nueva estructura de archivos del módulo bootstrap.

---

## Paso 3.2: Crear extractConfigField Genérico

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo** | `bootstrap/config_extractors.go` |
| **Duplicación Actual** | ~320 líneas |
| **Objetivo** | < 50 líneas |

### Implementación con Generics (Go 1.18+)

```go
package bootstrap

import (
    "fmt"
    "reflect"
)

// extractConfigField es un helper genérico para extraer configuración de un struct.
//
// Parámetros:
//   - config: Struct de configuración (valor o puntero)
//   - fieldName: Nombre del campo a extraer
//
// Retorna el valor del campo del tipo T o error si no se encuentra.
func extractConfigField[T any](config interface{}, fieldName string) (T, error) {
    var zero T
    
    // Intentar type assertion directo primero
    if typedConfig, ok := config.(T); ok {
        return typedConfig, nil
    }
    
    // Usar reflection para extraer el campo
    v := reflect.ValueOf(config)
    if v.Kind() == reflect.Ptr {
        if v.IsNil() {
            return zero, fmt.Errorf("config is nil")
        }
        v = v.Elem()
    }
    
    if v.Kind() != reflect.Struct {
        return zero, fmt.Errorf("config must be a struct, got %T", config)
    }
    
    field := v.FieldByName(fieldName)
    if !field.IsValid() {
        return zero, fmt.Errorf("%s field not found in config", fieldName)
    }
    
    // Intentar convertir el campo al tipo deseado
    fieldInterface := field.Interface()
    if typedField, ok := fieldInterface.(T); ok {
        return typedField, nil
    }
    
    // Si el campo es un puntero, intentar desreferenciarlo
    if field.Kind() == reflect.Ptr && !field.IsNil() {
        if typedField, ok := field.Elem().Interface().(T); ok {
            return typedField, nil
        }
    }
    
    return zero, fmt.Errorf("%s field is not of expected type, got %T", fieldName, fieldInterface)
}

// Funciones simplificadas usando el helper genérico
func extractPostgreSQLConfig(config interface{}) (PostgreSQLConfig, error) {
    return extractConfigField[PostgreSQLConfig](config, "PostgreSQL")
}

func extractMongoDBConfig(config interface{}) (MongoDBConfig, error) {
    return extractConfigField[MongoDBConfig](config, "MongoDB")
}

func extractRabbitMQConfig(config interface{}) (RabbitMQConfig, error) {
    return extractConfigField[RabbitMQConfig](config, "RabbitMQ")
}

func extractS3Config(config interface{}) (S3Config, error) {
    return extractConfigField[S3Config](config, "S3")
}
```

### Tests para el Helper Genérico

```go
func TestExtractConfigField(t *testing.T) {
    type TestConfig struct {
        Name  string
        Value int
    }
    
    type ParentConfig struct {
        Test TestConfig
    }
    
    tests := []struct {
        name      string
        config    interface{}
        fieldName string
        wantErr   bool
    }{
        {
            name:      "direct struct",
            config:    ParentConfig{Test: TestConfig{Name: "test", Value: 42}},
            fieldName: "Test",
            wantErr:   false,
        },
        {
            name:      "pointer to struct",
            config:    &ParentConfig{Test: TestConfig{Name: "test", Value: 42}},
            fieldName: "Test",
            wantErr:   false,
        },
        {
            name:      "field not found",
            config:    ParentConfig{},
            fieldName: "NonExistent",
            wantErr:   true,
        },
        {
            name:      "nil pointer",
            config:    (*ParentConfig)(nil),
            fieldName: "Test",
            wantErr:   true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := extractConfigField[TestConfig](tt.config, tt.fieldName)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotZero(t, result)
            }
        })
    }
}
```

### Criterios de Éxito

- [ ] Duplicación reducida de ~320 a ~50 líneas
- [ ] Tests del helper pasan
- [ ] Todas las funciones extract siguen funcionando
- [ ] `go build ./...` compila

### Commit

```bash
git add bootstrap/config_extractors.go bootstrap/config_extractors_test.go
git commit -m "refactor(bootstrap): crear extractConfigField genérico

- Reduce duplicación de ~320 a ~50 líneas
- Usa generics de Go 1.18+
- Mantiene API existente de funciones extract*Config
- Agrega tests para el helper genérico"
```

---

## Paso 3.3: Unificar MessagePublisher

### Información General

| Campo | Valor |
|-------|-------|
| **Implementaciones Actuales** | 2 (bootstrap y messaging/rabbit) |
| **Objetivo** | 1 implementación |

### Implementación

#### Paso 3.3.1: Modificar bootstrap para usar messaging/rabbit

En `bootstrap/init_rabbitmq.go`:

```go
import "github.com/EduGoGroup/edugo-shared/messaging/rabbit"

func initRabbitMQ(
    ctx context.Context,
    factories *Factories,
    config interface{},
    resources *Resources,
    lm *lifecycle.Manager,
    options *bootstrapOptions,
) error {
    // ... código existente de conexión ...
    
    // Usar implementación de messaging/rabbit
    publisher := rabbit.NewPublisher(conn)
    resources.MessagePublisher = &publisherAdapter{publisher}
    
    // ... resto del código ...
}

// publisherAdapter adapta rabbit.RabbitMQPublisher a la interface MessagePublisher
type publisherAdapter struct {
    publisher *rabbit.RabbitMQPublisher
}

func (a *publisherAdapter) Publish(ctx context.Context, queueName string, body []byte) error {
    return a.publisher.Publish(ctx, "", queueName, body)
}

func (a *publisherAdapter) PublishWithPriority(ctx context.Context, queueName string, body []byte, priority uint8) error {
    return a.publisher.PublishWithPriority(ctx, "", queueName, body, priority)
}

func (a *publisherAdapter) Close() error {
    return a.publisher.Close()
}
```

#### Paso 3.3.2: Eliminar implementación duplicada

Eliminar `defaultMessagePublisher` de `bootstrap/resource_implementations.go`.

### Criterios de Éxito

- [ ] Solo existe una implementación real de publisher
- [ ] Tests de bootstrap siguen pasando
- [ ] Tests de messaging siguen pasando

### Commit

```bash
git add bootstrap/init_rabbitmq.go bootstrap/resource_implementations.go
git commit -m "refactor(bootstrap): unificar MessagePublisher con messaging/rabbit

- Elimina implementación duplicada en bootstrap
- Usa adapter para mantener interface de bootstrap
- Centraliza lógica de publicación en un solo lugar"
```

---

## Paso 3.4: Corregir Tipo de presignClient

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo** | `bootstrap/resource_implementations.go` |
| **Línea** | 73 |
| **Cambio** | `interface{}` → `*s3.PresignClient` |

### Implementación

```go
import "github.com/aws/aws-sdk-go-v2/service/s3"

type defaultStorageClient struct {
    client        *s3.Client
    presignClient *s3.PresignClient  // Tipo específico en lugar de interface{}
    bucket        string
}
```

### Cambios en Factory

En `bootstrap/factory_s3.go`:

```go
func (f *defaultS3Factory) CreateStorageClient(
    client *s3.Client,
    bucket string,
) StorageClient {
    presignClient := s3.NewPresignClient(client)
    return &defaultStorageClient{
        client:        client,
        presignClient: presignClient,
        bucket:        bucket,
    }
}
```

### Actualizar GetPresignedURL

Ya no necesita type assertion:

```go
func (c *defaultStorageClient) GetPresignedURL(ctx context.Context, key string, expirationMinutes int) (string, error) {
    if c.presignClient == nil {
        return "", fmt.Errorf("presign client not initialized")
    }
    
    request, err := c.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(c.bucket),
        Key:    aws.String(key),
    }, func(opts *s3.PresignOptions) {
        opts.Expires = time.Duration(expirationMinutes) * time.Minute
    })
    if err != nil {
        return "", fmt.Errorf("failed to generate presigned URL: %w", err)
    }
    
    return request.URL, nil
}
```

### Criterios de Éxito

- [ ] No hay `interface{}` para presignClient
- [ ] Compile-time type safety
- [ ] GetPresignedURL no necesita type assertion
- [ ] Tests pasan

### Commit

```bash
git add bootstrap/resource_implementations.go bootstrap/factory_s3.go
git commit -m "refactor(bootstrap): cambiar presignClient de interface{} a *s3.PresignClient

- Mejora type safety con tipo específico
- Elimina necesidad de type assertion en runtime
- Detecta errores de tipo en compile-time"
```

---

## Paso 3.5: Agregar Control de Goroutine en Consumer

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo** | `messaging/rabbit/consumer.go` |
| **Problema** | Goroutine sin tracking ni error reporting |

### Implementación

```go
package rabbit

import (
    "context"
    "fmt"
    "sync"

    amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConsumer consume mensajes de una cola RabbitMQ con control de goroutines.
type RabbitMQConsumer struct {
    conn     *Connection
    config   ConsumerConfig
    wg       sync.WaitGroup
    errChan  chan error
    stopCh   chan struct{}
    stopOnce sync.Once
    mu       sync.Mutex
    running  bool
}

// NewConsumer crea un nuevo consumer con las opciones especificadas.
func NewConsumer(conn *Connection, config ConsumerConfig) *RabbitMQConsumer {
    return &RabbitMQConsumer{
        conn:    conn,
        config:  config,
        errChan: make(chan error, 1),
        stopCh:  make(chan struct{}),
    }
}

// Consume inicia el consumo de mensajes de la cola especificada.
func (c *RabbitMQConsumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
    c.mu.Lock()
    if c.running {
        c.mu.Unlock()
        return fmt.Errorf("consumer already running")
    }
    c.running = true
    c.mu.Unlock()

    msgs, err := c.conn.GetChannel().Consume(
        queueName,
        c.config.Name,
        c.config.AutoAck,
        c.config.Exclusive,
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        c.mu.Lock()
        c.running = false
        c.mu.Unlock()
        return fmt.Errorf("failed to start consuming: %w", err)
    }

    c.wg.Add(1)
    go func() {
        defer c.wg.Done()
        defer func() {
            c.mu.Lock()
            c.running = false
            c.mu.Unlock()
        }()

        for {
            select {
            case <-ctx.Done():
                return
            case <-c.stopCh:
                return
            case msg, ok := <-msgs:
                if !ok {
                    select {
                    case c.errChan <- fmt.Errorf("message channel closed unexpectedly"):
                    default:
                    }
                    return
                }
                c.processMessage(ctx, queueName, msg, handler)
            }
        }
    }()

    return nil
}

// Wait bloquea hasta que el consumer se detenga.
func (c *RabbitMQConsumer) Wait() error {
    c.wg.Wait()
    select {
    case err := <-c.errChan:
        return err
    default:
        return nil
    }
}

// Stop detiene el consumer de forma graceful.
func (c *RabbitMQConsumer) Stop() {
    c.stopOnce.Do(func() {
        close(c.stopCh)
    })
}

// Errors retorna un canal para recibir errores del consumer.
func (c *RabbitMQConsumer) Errors() <-chan error {
    return c.errChan
}

// IsRunning retorna true si el consumer está actualmente ejecutándose.
func (c *RabbitMQConsumer) IsRunning() bool {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.running
}

// Close detiene y limpia el consumer.
func (c *RabbitMQConsumer) Close() error {
    c.Stop()
    c.wg.Wait()
    return nil
}
```

### Criterios de Éxito

- [ ] Consumer tiene WaitGroup para esperar terminación
- [ ] Errores se reportan a través de canal
- [ ] Stop graceful implementado
- [ ] Tests pasan

### Commit

```bash
git add messaging/rabbit/consumer.go
git commit -m "refactor(messaging): agregar control de goroutine en Consumer

- Agrega WaitGroup para esperar terminación
- Implementa canal de errores para reportar problemas
- Agrega Stop() para terminación graceful
- Agrega IsRunning() para verificar estado"
```

---

## Verificación Final de Fase 3

### Antes de Crear el PR

```bash
# Ejecutar todos los tests
make test-all-modules

# Verificar que no hay duplicación
# (revisar manualmente que las funciones extract usan el helper)

# Build limpio
make build

# Linter
make lint
```

### Crear Pull Request

```bash
# Push de la rama
git push origin fase-3-refactoring-estructural

# En GitHub:
# 1. Crear PR hacia dev
# 2. Título: "refactor: Fase 3 - Refactoring Estructural"
# 3. Descripción con lista de refactorizaciones
```

### Revisión de GitHub Copilot

| Tipo de Comentario | Acción |
|-------------------|--------|
| Traducción inglés/español | DESCARTAR |
| Problema de diseño | CORREGIR |
| Mejora de rendimiento | DOCUMENTAR como deuda futura |

### Esperar Pipelines

```bash
# Revisar cada minuto durante máximo 10 minutos
# Si hay errores:
#   1. Analizar causa
#   2. Corregir (máx 3 intentos)
#   3. Push y esperar nuevamente
```

### Criterios de Éxito de Fase

- [ ] bootstrap.go < 150 líneas
- [ ] Duplicación reducida significativamente
- [ ] Una sola implementación de MessagePublisher
- [ ] Type safety mejorado
- [ ] Todos los tests pasan
- [ ] PR aprobado
- [ ] Pipelines verdes
- [ ] Merge a dev completado

---

## Resumen de la Fase 3

| Paso | Descripción | Reducción |
|------|-------------|-----------|
| 3.1 | Dividir bootstrap.go | 623 → 6 archivos ~100 c/u |
| 3.2 | extractConfigField genérico | 320 → 50 líneas |
| 3.3 | Unificar MessagePublisher | 2 → 1 implementación |
| 3.4 | Corregir tipo presignClient | Type safety mejorado |
| 3.5 | Control de goroutine | Mejor manejo de lifecycle |

---

## Siguiente Fase

Después de completar esta fase y hacer merge a dev, continuar con:
→ [FASE-4: Mejoras de Calidad](./FASE-4_MEJORAS_CALIDAD.md)
