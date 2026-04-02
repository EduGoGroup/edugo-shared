# Metrics

Fachada agnóstica para registro de métricas de aplicación (Prometheus, Datadog, OpenTelemetry, NoOp).

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/metrics
```

El módulo se descarga como `metrics`, principal consumo vía package `metrics`.

## Quick Start

### Ejemplo 1: Inicializar Metrics con NoOp (por defecto)

```go
package main

import (
    "fmt"
    "time"
    "github.com/EduGoGroup/edugo-shared/metrics"
)

func main() {
    // Sin argumentos: usa NoopRecorder (cero overhead)
    m := metrics.New("my-service")

    // Registrar métrica de login (sin hacer nada)
    m.RecordLogin(true, 100*time.Millisecond)

    // Registrar métrica HTTP
    m.RecordHTTPRequest("GET", "/api/users", 200, 50*time.Millisecond)

    // Obtener el recorder subyacente (útil para casos avanzados)
    recorder := m.Recorder()
    fmt.Printf("Recorder type: %T\n", recorder)

    // Obtener nombre del servicio
    fmt.Printf("Service: %s\n", m.Service())
}
```

### Ejemplo 2: Registrar métricas de autenticación

```go
package main

import (
    "time"
    "github.com/EduGoGroup/edugo-shared/metrics"
)

func main() {
    m := metrics.New("auth-service")

    // Registrar intento de login exitoso
    start := time.Now()
    // ... hacer login
    m.RecordLogin(true, time.Since(start))

    // Registrar intento fallido
    start = time.Now()
    // ... login fallido
    m.RecordLogin(false, time.Since(start))

    // Registrar refresh de token
    start = time.Now()
    // ... refresh token
    m.RecordTokenRefresh(true, time.Since(start))

    // Registrar verificación de permisos
    m.RecordPermissionCheck("admin.delete", true)
    m.RecordPermissionCheck("editor.edit", false)

    // Registrar rate limiting
    m.RecordRateLimitHit("api.login")
    m.RecordRateLimitHit("api.create_material")
}
```

### Ejemplo 3: Registrar métricas de operaciones CRUD y HTTP

```go
package main

import (
    "fmt"
    "time"
    "github.com/EduGoGroup/edugo-shared/metrics"
)

func recordDatabaseOperations(m *metrics.Metrics) {
    // Registrar query exitosa
    start := time.Now()
    // ... SELECT * FROM users WHERE id = 123
    m.RecordDBQuery("postgres", "select", "users", time.Since(start), nil)

    // Registrar query con error
    start = time.Now()
    // ... INSERT fallido
    err := fmt.Errorf("constraint violation")
    m.RecordDBQuery("postgres", "insert", "assessments", time.Since(start), err)

    // Registrar conexiones abiertas
    m.SetDBConnectionsOpen("postgres", 25)
    m.SetDBConnectionsOpen("mongodb", 10)
}

func recordHTTPOperations(m *metrics.Metrics) {
    // Registrar request exitoso
    m.RecordHTTPRequest("POST", "/api/v1/users", 201, 120*time.Millisecond)

    // Registrar request con error
    m.RecordHTTPRequest("DELETE", "/api/v1/materials/{id}", 500, 500*time.Millisecond)

    // Registrar requests activos
    m.SetHTTPActiveRequests(42)
}

func main() {
    m := metrics.New("api-service")
    recordDatabaseOperations(m)
    recordHTTPOperations(m)
}
```

### Ejemplo 4: Registrar métricas de lógica de negocio avanzada

```go
package main

import (
    "time"
    "github.com/EduGoGroup/edugo-shared/metrics"
)

func main() {
    m := metrics.New("assessment-service")

    // Registrar operaciones de entidad
    start := time.Now()
    // ... crear assessment
    m.RecordBusinessOperation("assessment", "create", time.Since(start), nil)

    start = time.Now()
    // ... actualizar assessment
    m.RecordBusinessOperation("assessment", "update", time.Since(start), nil)

    // Registrar intentos de evaluación
    start = time.Now()
    // ... student inicia evaluación
    m.RecordAssessmentAttempt("start", time.Since(start), nil)

    start = time.Now()
    // ... student envía evaluación
    m.RecordAssessmentAttempt("submit", time.Since(start), nil)

    // Registrar operaciones de calificación
    start = time.Now()
    // ... calificar pregunta múltiple
    m.RecordGrading("multiple_choice", time.Since(start), nil)

    // Registrar operaciones de messaging
    start = time.Now()
    // ... procesar mensaje de event bus
    m.RecordMessageProcessed("assessment_attempt", time.Since(start), nil)

    // Registrar estado de circuit breaker
    m.SetCircuitBreakerState("payment-service", "closed")

    // Registrar mensajes en cola
    m.SetMessagesInQueue("assessment_events", 15)

    // Registrar operaciones de notificación
    m.RecordNotification("email", nil)
    m.RecordNotification("push", nil)

    // Registrar operaciones de exportación
    m.RecordExport("markdown", 1500, 2500*time.Millisecond, nil)
}
```

## Componentes principales

- **Metrics**: Punto de entrada central, fachada agnóstica con métodos de registro
- **Recorder**: Interfaz extensible para backends (Prometheus, Datadog, OpenTelemetry)
- **NoopRecorder**: Implementación por defecto (cero overhead)
- Métodos especializados: Auth, HTTP, DB, Business, Messaging, Export

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build     # Compilar
make test      # Tests unitarios
make test-race # Race detector
make check     # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- **Agnóstico de backend**: La interfaz `Recorder` permite implementar Prometheus, Datadog, OpenTelemetry sin cambios en código cliente
- **Zero overhead por defecto**: `NoopRecorder` incurre cero overhead, perfecto para desarrollo y testing
- **Seguridad de concurrencia**: La seguridad depende del `Recorder` implementado (NoopRecorder es siempre seguro)
- **Labels estandarizados**: Cada métrica incluye automáticamente el nombre del servicio y estado (success/error)
- **Duración en segundos**: Histogramas de duración se normalizan a segundos para consistencia cross-backend
