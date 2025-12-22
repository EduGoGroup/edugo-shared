# FASE 4: Mejoras de Calidad

> **Prioridad**: MEDIA  
> **Duración estimada**: 2 días  
> **Prerrequisitos**: Fase 3 completada  
> **Rama**: `fase-4-mejoras-calidad`  
> **Objetivo**: Limpiar código y mejorar consistencia

---

## Flujo de Trabajo de Esta Fase

### 1. Inicio de la Fase

```bash
# Asegurarse de estar en dev actualizado
git checkout dev
git pull origin dev

# Crear rama de la fase
git checkout -b fase-4-mejoras-calidad

# Verificar estado inicial
make build
make test-all-modules
```

### 2. Durante la Fase

- Ejecutar cada paso en orden
- Commit atómico después de cada paso completado
- Verificar que tests pasen después de cada cambio

### 3. Fin de la Fase

```bash
# Push de la rama
git push origin fase-4-mejoras-calidad

# Crear PR en GitHub hacia dev
# - Título: "chore: Fase 4 - Mejoras de Calidad"
# - Descripción: Lista de mejoras realizadas

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

Esta fase se enfoca en mejoras de calidad de código:
- Limpieza de código muerto
- Documentación de APIs
- Consistencia en logging
- Migración a interfaces abstractas

---

## Paso 4.1: Limpiar Imports Comentados

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo** | `messaging/rabbit/consumer.go` |
| **Línea** | 7 |
| **Tipo** | Limpieza |

### Acción

Eliminar cualquier import comentado del código:

```go
// ANTES
import (
    "context"
    "encoding/json"
    "fmt"
    // amqp "github.com/rabbitmq/amqp091-go" // No usado actualmente
)

// DESPUÉS
import (
    "context"
    "encoding/json"
    "fmt"
)
```

### Búsqueda de Otros Imports Comentados

```bash
# Buscar imports comentados en todo el proyecto
grep -rn "// .*import" --include="*.go" .
grep -rn "//.*github.com" --include="*.go" . | grep -v "_test.go"
```

### Criterios de Éxito

- [ ] No hay imports comentados
- [ ] `go build ./...` compila
- [ ] `goimports` no reporta cambios

### Commit

```bash
git add messaging/rabbit/consumer.go
git commit -m "chore(messaging): eliminar imports comentados

- Limpia código muerto
- Mejora legibilidad"
```

---

## Paso 4.2: Documentar API de Containers

### Información General

| Campo | Valor |
|-------|-------|
| **Ubicación** | `testing/containers/` |
| **Problema** | API no documentada |

### Crear Archivo de Documentación

Crear `testing/containers/README.md`:

```markdown
# Testing Containers

Infraestructura de testing con testcontainers para PostgreSQL, MongoDB y RabbitMQ.

## Uso Básico

```go
func TestExample(t *testing.T) {
    ctx := context.Background()
    
    config := containers.NewConfig().
        WithPostgreSQL(nil).
        WithMongoDB(nil).
        Build()
    
    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    defer manager.Cleanup(ctx)
    
    // PostgreSQL
    pgConfig, _ := manager.PostgreSQL().Config(ctx)
    
    // MongoDB
    mongoURI, _ := manager.MongoDB().ConnectionString(ctx)
    
    // RabbitMQ
    rmqURL, _ := manager.RabbitMQ().ConnectionString(ctx)
}
```

## API

### Manager

```go
// GetManager obtiene o crea el singleton del manager
func GetManager(t *testing.T, config *Config) (*Manager, error)

// Cleanup libera todos los containers
func (m *Manager) Cleanup(ctx context.Context) error
```

### PostgreSQL Container

```go
// Config retorna la configuración de conexión
func (c *PostgreSQLContainer) Config(ctx context.Context) (*PostgreSQLConfig, error)

// ConnectionString retorna el DSN de conexión
func (c *PostgreSQLContainer) ConnectionString(ctx context.Context) (string, error)
```

### MongoDB Container

```go
// ConnectionString retorna la URI de conexión
func (c *MongoDBContainer) ConnectionString(ctx context.Context) (string, error)
```

### RabbitMQ Container

```go
// ConnectionString retorna la URL AMQP
func (c *RabbitMQContainer) ConnectionString(ctx context.Context) (string, error)
```

## Configuración

### Configurar container específico

```go
config := containers.NewConfig().
    WithPostgreSQL(&containers.PostgreSQLOptions{
        Database: "mydb",
        User:     "myuser",
    }).
    Build()
```

### Timeouts

Por defecto, los containers tienen un timeout de 30 segundos para iniciar.
Puedes configurarlo:

```go
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

pgConfig, err := manager.PostgreSQL().Config(ctx)
```
```

### Criterios de Éxito

- [ ] README.md creado en testing/containers/
- [ ] API documentada claramente
- [ ] Ejemplos de uso incluidos

### Commit

```bash
git add testing/containers/README.md
git commit -m "docs(testing): documentar API de containers

- Agrega README con API pública
- Documenta uso básico con ejemplos
- Incluye configuración y timeouts"
```

### Actualización de Documentación Principal

Actualizar `documents/TESTING.md` para referenciar la nueva documentación de containers.

---

## Paso 4.3: Migrar a logger.Logger Interface

### Información General

| Campo | Valor |
|-------|-------|
| **Problema** | bootstrap usa `*logrus.Logger` directamente |
| **Objetivo** | Usar interface `logger.Logger` para flexibilidad |

### Cambios en bootstrap/interfaces.go

```go
import "github.com/EduGoGroup/edugo-shared/logger"

// Resources contiene todos los recursos inicializados por Bootstrap.
type Resources struct {
    Logger logger.Logger  // Interface abstracta en lugar de implementación específica
    // ...otros campos...
}
```

### Actualizar LoggerFactory

```go
// LoggerFactory crea instancias de logger.
type LoggerFactory interface {
    CreateLogger(ctx context.Context, env, version string) (logger.Logger, error)
}
```

### Beneficios

- Permite cambiar implementación de logger (logrus → zap → slog)
- Facilita testing con mocks
- Desacopla bootstrap de implementación específica

### Criterios de Éxito

- [ ] Resources.Logger es de tipo `logger.Logger`
- [ ] No hay referencias directas a logrus en bootstrap interfaces
- [ ] Tests siguen pasando

### Commit

```bash
git add bootstrap/interfaces.go bootstrap/init_logger.go
git commit -m "refactor(bootstrap): migrar a logger.Logger interface

- Cambia Resources.Logger de *logrus.Logger a logger.Logger
- Actualiza LoggerFactory para retornar interface
- Desacopla bootstrap de implementación específica"
```

---

## Paso 4.4: Crear Constantes para Campos de Log

### Información General

| Campo | Valor |
|-------|-------|
| **Problema** | Campos de log inconsistentes (user_id vs userId) |
| **Objetivo** | Constantes centralizadas para campos comunes |

### Crear Archivo de Constantes

Crear `logger/fields.go`:

```go
package logger

// Campos estándar para logging estructurado.
// Usar estas constantes garantiza consistencia en todos los logs.
const (
    // Identificadores
    FieldUserID     = "user_id"
    FieldRequestID  = "request_id"
    FieldSessionID  = "session_id"
    FieldTraceID    = "trace_id"
    
    // Servicio y operación
    FieldService    = "service"
    FieldOperation  = "operation"
    FieldMethod     = "method"
    FieldPath       = "path"
    
    // Recursos
    FieldResourceID   = "resource_id"
    FieldResourceType = "resource_type"
    FieldQueue        = "queue"
    FieldExchange     = "exchange"
    FieldDatabase     = "database"
    FieldCollection   = "collection"
    FieldTable        = "table"
    
    // Métricas
    FieldDurationMs = "duration_ms"
    FieldStatusCode = "status_code"
    FieldBytesIn    = "bytes_in"
    FieldBytesOut   = "bytes_out"
    
    // Errores
    FieldError        = "error"
    FieldErrorCode    = "error_code"
    FieldErrorMessage = "error_message"
    FieldStackTrace   = "stack_trace"
    
    // Mensajería
    FieldDeliveryTag = "delivery_tag"
    FieldMessageID   = "message_id"
    FieldPriority    = "priority"
    FieldRetryCount  = "retry_count"
)
```

### Uso en Código

```go
import "github.com/EduGoGroup/edugo-shared/logger"

// Usar constantes para campos de log
log.Info("user created", logger.FieldUserID, userID)
log.Error("failed to process message",
    logger.FieldError, err,
    logger.FieldDeliveryTag, msg.DeliveryTag,
    logger.FieldQueue, queueName,
)
```

### Criterios de Éxito

- [ ] Archivo `logger/fields.go` creado
- [ ] Campos comunes definidos como constantes
- [ ] Al menos un uso actualizado para validar

### Commit

```bash
git add logger/fields.go
git commit -m "feat(logger): agregar constantes para campos de log

- Define campos estándar para logging estructurado
- Asegura consistencia en nombres de campos
- Facilita búsqueda y análisis de logs"
```

### Actualización de Documentación

Actualizar `documents/SERVICES.md` o crear sección en documentación para explicar el uso de las constantes de logging.

---

## Verificación Final de Fase 4

### Antes de Crear el PR

```bash
# Verificar imports
goimports -l .

# Verificar documentación existe
ls testing/containers/README.md

# Build y tests
make build
make test-all-modules

# Linter
make lint
```

### Crear Pull Request

```bash
# Push de la rama
git push origin fase-4-mejoras-calidad

# En GitHub:
# 1. Crear PR hacia dev
# 2. Título: "chore: Fase 4 - Mejoras de Calidad"
# 3. Descripción con lista de mejoras
```

### Revisión de GitHub Copilot

| Tipo de Comentario | Acción |
|-------------------|--------|
| Traducción inglés/español | DESCARTAR |
| Sugerencia de mejora | EVALUAR y posiblemente DOCUMENTAR |
| Error detectado | CORREGIR |

### Esperar Pipelines

```bash
# Revisar cada minuto durante máximo 10 minutos
# Si hay errores:
#   1. Analizar causa
#   2. Corregir (máx 3 intentos)
#   3. Push y esperar nuevamente
```

### Criterios de Éxito de Fase

- [ ] No hay imports comentados
- [ ] Documentación de containers creada
- [ ] Interface de logger abstracta
- [ ] Campos de log estandarizados
- [ ] Todos los tests pasan
- [ ] PR aprobado
- [ ] Pipelines verdes
- [ ] Merge a dev completado

---

## Resumen de la Fase 4

| Paso | Descripción | Commit |
|------|-------------|--------|
| 4.1 | Limpiar imports comentados | `chore(messaging): eliminar imports comentados` |
| 4.2 | Documentar API containers | `docs(testing): documentar API de containers` |
| 4.3 | Migrar a logger.Logger | `refactor(bootstrap): migrar a logger.Logger` |
| 4.4 | Crear constantes de log | `feat(logger): agregar constantes de log` |

---

## Siguiente Fase

Después de completar esta fase y hacer merge a dev, continuar con:
→ [FASE-5: Deuda Técnica](./FASE-5_DEUDA_TECNICA.md)
