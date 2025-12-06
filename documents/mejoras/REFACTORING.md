# Refactoring Necesario

> Código que funciona pero necesita reestructuración para mejorar mantenibilidad, testabilidad o extensibilidad.

---

## Issue #1: Bootstrap.go Demasiado Grande

### Ubicación
```
bootstrap/bootstrap.go (623 líneas)
```

### Problema
- Un solo archivo con 623 líneas
- Mezcla múltiples responsabilidades:
  - Función principal Bootstrap
  - Inicialización de Logger
  - Inicialización de PostgreSQL
  - Inicialización de MongoDB
  - Inicialización de RabbitMQ
  - Inicialización de S3
  - Health checks
  - Funciones de extracción de config
  - Funciones de registro de cleanup
- Difícil de testear unitariamente
- Difícil de navegar y entender

### Impacto
- **Medio**: Mantenibilidad reducida

### Solución Sugerida
Dividir en múltiples archivos:

```
bootstrap/
├── bootstrap.go           # Solo función principal Bootstrap()
├── init_logger.go         # initLogger y helpers
├── init_postgresql.go     # initPostgreSQL y helpers
├── init_mongodb.go        # initMongoDB y helpers
├── init_rabbitmq.go       # initRabbitMQ y helpers
├── init_s3.go             # initS3 y helpers
├── health_check.go        # performHealthChecks
├── config_extractors.go   # Todas las funciones extract*Config
├── cleanup_registrars.go  # Todas las funciones register*Cleanup
├── interfaces.go          # (ya existe)
├── options.go             # (ya existe)
└── resources.go           # (ya existe)
```

### Beneficios
- Archivos más pequeños (~100-150 líneas cada uno)
- Cada archivo con una responsabilidad clara
- Más fácil de testear independientemente
- Mejor navegación en IDE

### Prioridad: **MEDIA**

---

## Issue #2: Duplicación en Funciones de Extracción de Config

### Ubicación
```
bootstrap/bootstrap.go:438-552
```

### Código Actual
```go
// extractPostgreSQLConfig extrae configuración de PostgreSQL usando reflection
func extractPostgreSQLConfig(config interface{}) (PostgreSQLConfig, error) {
    // Intentar type assertion directo primero
    if pgConfig, ok := config.(PostgreSQLConfig); ok {
        return pgConfig, nil
    }

    // Usar reflection para extraer campo PostgreSQL
    v := reflect.ValueOf(config)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    if v.Kind() != reflect.Struct {
        return PostgreSQLConfig{}, fmt.Errorf("config must be a struct, got %T", config)
    }

    pgField := v.FieldByName("PostgreSQL")
    if !pgField.IsValid() {
        return PostgreSQLConfig{}, fmt.Errorf("PostgreSQL field not found in config")
    }
    // ...
}

// extractMongoDBConfig - MISMO PATRÓN
// extractRabbitMQConfig - MISMO PATRÓN
// extractS3Config - MISMO PATRÓN
```

### Problema
- 4 funciones con el mismo patrón de reflection
- Código duplicado: ~80 líneas × 4 = 320 líneas
- Si cambia el patrón, hay que cambiar 4 lugares
- Propenso a errores de copia-paste

### Solución Sugerida
```go
// extractConfigField es un helper genérico para extraer configuración
func extractConfigField[T any](config interface{}, fieldName string) (T, error) {
    var zero T
    
    // Intentar type assertion directo
    if typedConfig, ok := config.(T); ok {
        return typedConfig, nil
    }

    // Usar reflection
    v := reflect.ValueOf(config)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    if v.Kind() != reflect.Struct {
        return zero, fmt.Errorf("config must be a struct, got %T", config)
    }

    field := v.FieldByName(fieldName)
    if !field.IsValid() {
        return zero, fmt.Errorf("%s field not found in config", fieldName)
    }

    if typedField, ok := field.Interface().(T); ok {
        return typedField, nil
    }

    return zero, fmt.Errorf("%s field is not of expected type", fieldName)
}

// Uso simplificado
func extractPostgreSQLConfig(config interface{}) (PostgreSQLConfig, error) {
    return extractConfigField[PostgreSQLConfig](config, "PostgreSQL")
}

func extractMongoDBConfig(config interface{}) (MongoDBConfig, error) {
    return extractConfigField[MongoDBConfig](config, "MongoDB")
}
```

### Beneficios
- De ~320 líneas a ~50 líneas
- Un solo lugar para mantener la lógica
- Más fácil de testear
- Usa generics de Go 1.18+

### Prioridad: **MEDIA**

---

## Issue #3: MessagePublisher en Bootstrap vs Messaging Module

### Ubicación
```
bootstrap/resource_implementations.go:18-64
bootstrap/interfaces.go:118-128
messaging/rabbit/publisher.go:12-83
```

### Problema
- Hay DOS implementaciones de MessagePublisher:
  1. `defaultMessagePublisher` en bootstrap (simple)
  2. `RabbitMQPublisher` en messaging/rabbit (completo)
- La de bootstrap es más limitada
- Confusión sobre cuál usar
- Duplicación de código y conceptos

### Comparación

| Feature | bootstrap | messaging/rabbit |
|---------|-----------|-----------------|
| Publish básico | ✅ | ✅ |
| Publish con prioridad | ✅ | ✅ |
| Serialización JSON | ❌ | ✅ |
| Exchange routing | ❌ | ✅ |
| Configuración | Mínima | Completa |

### Solución Sugerida

**Opción A: Eliminar implementación de bootstrap**
```go
// bootstrap/bootstrap.go
func initRabbitMQ(...) error {
    // ...
    
    // Usar implementación de messaging/rabbit directamente
    resources.MessagePublisher = rabbit.NewPublisher(conn)
    
    // ...
}
```

**Opción B: Abstraer con adapter**
```go
// bootstrap/adapters.go
type bootstrapPublisherAdapter struct {
    rabbitPublisher *rabbit.RabbitMQPublisher
}

func (a *bootstrapPublisherAdapter) Publish(ctx context.Context, queueName string, body []byte) error {
    return a.rabbitPublisher.Publish(ctx, "", queueName, body)
}
```

### Recomendación
Opción A: Eliminar la implementación duplicada de bootstrap y usar la de messaging/rabbit.

### Prioridad: **MEDIA**

---

## Issue #4: Testing Containers API Inconsistente

### Ubicación
```
testing/containers/
```

### Problema
- API de containers cambió y rompió tests de integración
- No hay versioning de la API interna
- Cambios breaking sin migración clara

### Evidencia
```go
// Antes (tests que ahora fallan)
mongo := manager.MongoDB()
uri := mongo.URI  // Acceso directo a campo

// Después (API actual)
mongo := manager.MongoDB()
uri, err := mongo.ConnectionString(ctx)  // Método con error
```

### Solución Sugerida
1. Documentar API pública de containers
2. Agregar deprecated warnings antes de cambios
3. Mantener compatibilidad hacia atrás por 1-2 versiones
4. Crear migration guide cuando hay cambios breaking

```go
// Ejemplo de deprecation gradual
type MongoDBContainer struct {
    // ...
}

// Deprecated: Use ConnectionString(ctx) instead
// Will be removed in v0.9.0
func (c *MongoDBContainer) URI() string {
    ctx := context.Background()
    uri, _ := c.ConnectionString(ctx)
    return uri
}

// Nuevo método preferido
func (c *MongoDBContainer) ConnectionString(ctx context.Context) (string, error) {
    // ...
}
```

### Prioridad: **MEDIA**

---

## Issue #5: Logger Interface vs Logrus Dependency

### Ubicación
```
logger/logger.go (interface)
bootstrap/bootstrap.go (usa *logrus.Logger directamente)
```

### Problema
- `logger` module define una interface `Logger` genérica
- `bootstrap` usa `*logrus.Logger` directamente
- Inconsistencia que dificulta cambiar de logger

### Código Actual
```go
// logger/logger.go
type Logger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    // ...
}

// bootstrap/interfaces.go - Resources usa logrus directamente
type Resources struct {
    Logger *logrus.Logger  // <-- Debería ser logger.Logger
    // ...
}
```

### Solución
```go
// bootstrap/interfaces.go
import "github.com/EduGoGroup/edugo-shared/logger"

type Resources struct {
    Logger logger.Logger  // Interface, no implementación
    // ...
}

// Esto permite usar cualquier implementación:
// - logrus
// - zap
// - slog (Go 1.21+)
// - mock para tests
```

### Prioridad: **BAJA** (funciona, pero limita flexibilidad)

---

## Plan de Refactoring

### Sprint de Deuda Técnica (Sugerido)

**Semana 1: Preparación**
- [ ] Crear branches de feature para cada refactoring
- [ ] Escribir tests adicionales para código existente
- [ ] Documentar comportamiento actual

**Semana 2: Refactoring Core**
- [ ] Dividir bootstrap.go en archivos más pequeños
- [ ] Implementar extractConfigField genérico
- [ ] Unificar MessagePublisher

**Semana 3: API Cleanup**
- [ ] Documentar API de containers
- [ ] Agregar deprecation warnings
- [ ] Migrar a logger.Logger interface

**Semana 4: Validación**
- [ ] Ejecutar suite completa de tests
- [ ] Verificar coverage >= 80%
- [ ] Code review
- [ ] Merge a main

---

## Métricas de Éxito

| Métrica | Antes | Objetivo |
|---------|-------|----------|
| Líneas en bootstrap.go | 623 | < 150 |
| Duplicación de código | ~320 líneas | < 50 líneas |
| Implementaciones de Publisher | 2 | 1 |
| Tests de integración passing | 0 | 15+ |
