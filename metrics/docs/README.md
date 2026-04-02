# Metrics — Documentación técnica

Fachada agnóstica para registro de métricas de aplicación (contadores, histogramas, gauges).

## Propósito

El módulo `metrics` proporciona una abstracción agnóstica que permite registrar métricas de aplicación (autenticación, HTTP, base de datos, lógica de negocio) sin acoplamiento a backends específicos (Prometheus, Datadog, OpenTelemetry). Define una interfaz `Recorder` que implementaciones concretas pueden extender para soportar cualquier backend de métricas.

## Componentes principales

### Metrics — Fachada central

Punto de entrada para registrar todas las métricas de la aplicación. Una instancia por servicio, pasada a componentes que necesiten instrumentación.

**Función/Métodos:**
```go
func New(service string, recorder ...Recorder) *Metrics

func (m *Metrics) Service() string              // Nombre del servicio
func (m *Metrics) Recorder() Recorder           // Acceso al recorder subyacente
```

**Métodos de registro (delegados al Recorder):**
- Autenticación: `RecordLogin`, `RecordTokenRefresh`, `RecordRateLimitHit`, `RecordPermissionCheck`
- HTTP: `RecordHTTPRequest`, `SetHTTPActiveRequests`
- Base de datos: `RecordDBQuery`, `SetDBConnectionsOpen`
- Negocio: `RecordBusinessOperation`, `RecordAssessmentAttempt`, `RecordGrading`, `RecordReview`, `RecordNotification`, `RecordExport`
- Messaging: `RecordMessageProcessed`, `RecordMessageRetry`, `SetMessagesInQueue`, `SetCircuitBreakerState`

**Características:**
- Inicialización simple: `New("service-name")` sin argumentos usa `NoopRecorder` (cero overhead)
- Inicialización con backend: `New("service-name", prometheusRecorder)` para métricas reales
- Labels automáticos: Nombre del servicio se agrega a todas las métricas
- Seguridad de concurrencia: Delegada al `Recorder` (NoopRecorder es siempre seguro)

**Ejemplo:**
```go
m := metrics.New("api-service")
m.RecordLogin(true, 150*time.Millisecond)           // Sin hacer nada si usa NoopRecorder
m.RecordHTTPRequest("GET", "/users", 200, 50*ms)   // Delegado al recorder
recorder := m.Recorder()                             // Acceso para casos avanzados
```

### Recorder — Interfaz de backend

Interfaz que implementaciones concretas deben satisfacer para registrar métricas.

**Métodos:**
```go
type Recorder interface {
    CounterAdd(name string, value float64, labels map[string]string)
    HistogramObserve(name string, value float64, labels map[string]string)
    GaugeSet(name string, value float64, labels map[string]string)
}
```

**Métodos principales:**
- `CounterAdd(name, value, labels)` — Incrementa un contador (eventos, fallos, etc.)
- `HistogramObserve(name, value, labels)` — Registra una muestra en un histograma (duraciones, tamaños)
- `GaugeSet(name, value, labels)` — Establece un gauge (conexiones activas, cola de mensajes)

**Características:**
- Labels como mapa clave-valor para multidimensionalidad
- Nombres estandarizados (constantes públicas definidas por tipo)
- Sin retorno de error: implementaciones deben ser resilientes
- Seguridad de concurrencia: responsabilidad del implementador

**Ejemplo:**
```go
// Implementación personalizada
type PrometheusRecorder struct {
    counter   prometheus.CounterVec
    histogram prometheus.HistogramVec
    gauge     prometheus.GaugeVec
}

func (r *PrometheusRecorder) CounterAdd(name string, value float64, labels map[string]string) {
    r.counter.WithLabelValues(extractValues(labels)...).Add(value)
}
```

### NoopRecorder — Implementación por defecto

Implementación que no hace nada. Todos los métodos son no-ops de cero costo, perfecto para desarrollo y testing.

**Métodos:**
```go
type NoopRecorder struct{}

func (n *NoopRecorder) CounterAdd(string, float64, map[string]string) {}
func (n *NoopRecorder) HistogramObserve(string, float64, map[string]string) {}
func (n *NoopRecorder) GaugeSet(string, float64, map[string]string) {}
```

**Características:**
- Cero overhead: métodos vacíos, no asignan memoria
- Seguro para concurrencia: sin estado compartido
- Perfecto para testing y desarrollo local
- Permite activar/desactivar métricas con cambio mínimo de código

**Ejemplo:**
```go
// Uso implícito
m := metrics.New("service")  // Usa NoopRecorder por defecto

// Uso explícito
m := metrics.New("service", &metrics.NoopRecorder{})
```

### Métodos de autenticación

Registro de eventos de autenticación (logins, tokens, permisos).

**Métodos principales:**
- `RecordLogin(success bool, duration time.Duration)` — Intento de login
- `RecordTokenRefresh(success bool, duration time.Duration)` — Refresh de token
- `RecordRateLimitHit(resource string)` — Cuando se activa un rate limit
- `RecordPermissionCheck(permission string, granted bool)` — Verificación de permisos

**Características:**
- Duración en `time.Duration`, convertida a segundos internamente
- Status automático: "success"/"failure" para login/token, "granted"/"denied" para permisos
- Labels: service, status, resource (para rate limit), permission (para permisos)

**Constantes de nombres:**
```
MetricAuthLoginsTotal = "auth_logins_total"
MetricAuthLoginDuration = "auth_login_duration_seconds"
MetricAuthTokenRefreshTotal = "auth_token_refresh_total"
MetricAuthTokenRefreshDuration = "auth_token_refresh_duration"
MetricAuthRateLimitHits = "auth_rate_limit_hits_total"
MetricAuthPermissionChecks = "auth_permission_checks_total"
```

**Ejemplo:**
```go
start := time.Now()
user, err := auth.Login(username, password)
m.RecordLogin(err == nil, time.Since(start))

m.RecordTokenRefresh(tokenValid, 45*time.Millisecond)
m.RecordRateLimitHit("api.login")
m.RecordPermissionCheck("admin.delete_user", hasPermission)
```

### Métodos HTTP

Registro de solicitudes HTTP (método, ruta, status, duración).

**Métodos principales:**
- `RecordHTTPRequest(method, path string, status int, duration time.Duration)` — Request completada
- `SetHTTPActiveRequests(count int)` — Número actual de requests en vuelo

**Características:**
- Status como número entero (200, 404, 500, etc.)
- Path como template (no incluye parámetros específicos para evitar cardinalidad alta)
- Labels: service, method, path, status

**Constantes de nombres:**
```
MetricHTTPRequestsTotal = "http_requests_total"
MetricHTTPRequestDuration = "http_request_duration_seconds"
MetricHTTPActiveRequests = "http_active_requests"
```

**Ejemplo:**
```go
start := time.Now()
// ... servir HTTP request
m.RecordHTTPRequest("POST", "/api/v1/users", 201, time.Since(start))

// En middleware
m.SetHTTPActiveRequests(activeRequestCount)
```

### Métodos de base de datos

Registro de operaciones de base de datos (queries, conexiones).

**Métodos principales:**
- `RecordDBQuery(dbType, operation, table string, duration time.Duration, err error)` — Query completada
- `SetDBConnectionsOpen(dbType string, count int)` — Pool de conexiones actual

**Características:**
- dbType: "postgres", "mongodb", etc.
- operation: "select", "insert", "update", "delete"
- Status automático: "success"/"error" basado en err
- Labels: service, db_type, operation, table, status

**Constantes de nombres:**
```
MetricDBQueriesTotal = "db_queries_total"
MetricDBQueryDuration = "db_query_duration_seconds"
MetricDBConnectionsOpen = "db_connections_open"
```

**Ejemplo:**
```go
start := time.Now()
users, err := db.Query("SELECT * FROM users WHERE id = ?", id)
m.RecordDBQuery("postgres", "select", "users", time.Since(start), err)

m.SetDBConnectionsOpen("postgres", 25)
m.SetDBConnectionsOpen("mongodb", 10)
```

### Métodos de negocio

Registro de operaciones de lógica de negocio (CRUD, evaluaciones, calificación).

**Métodos principales:**
- `RecordBusinessOperation(entity, operation string, duration, err)` — CRUD genérico
- `RecordAssessmentAttempt(action string, duration, err)` — Intento de evaluación
- `RecordGrading(questionType string, duration, err)` — Calificación
- `RecordReview(action string, duration, err)` — Revisión de evaluación
- `RecordNotification(channel string, err)` — Envío de notificación
- `RecordExport(format string, rows int, duration, err)` — Exportación de datos

**Características:**
- Duración en `time.Duration`, convertida a segundos
- Status automático: "success"/"error"
- Entity-specific labels: entity, operation, question_type, action, channel, format

**Constantes de nombres:**
```
MetricBusinessOpsTotal = "business_operations_total"
MetricBusinessOpsDuration = "business_operations_duration_seconds"
MetricAssessmentAttempts = "assessment_attempts_total"
MetricGradingTotal = "grading_operations_total"
MetricGradingDuration = "grading_duration_seconds"
MetricReviewTotal = "review_operations_total"
MetricReviewDuration = "review_duration_seconds"
MetricNotificationTotal = "notification_operations_total"
MetricExportTotal = "export_operations_total"
MetricExportDuration = "export_duration_seconds"
MetricExportRows = "export_rows_total"
```

**Ejemplo:**
```go
start := time.Now()
err := createAssessment(assessmentData)
m.RecordBusinessOperation("assessment", "create", time.Since(start), err)

start = time.Now()
m.RecordAssessmentAttempt("submit", time.Since(start), nil)

m.RecordGrading("multiple_choice", 50*time.Millisecond, nil)

m.RecordNotification("email", sendError)

m.RecordExport("markdown", 1500, 2*time.Second, nil)
```

### Métodos de messaging

Registro de eventos de message queue y circuit breaker.

**Métodos principales:**
- `RecordMessageProcessed(eventType string, duration, err)` — Mensaje procesado
- `RecordMessageRetry(eventType string, attempt int)` — Reintento de mensaje
- `SetMessagesInQueue(queueName string, count int)` — Mensajes pendientes
- `SetCircuitBreakerState(targetService, state string)` — Estado del circuit breaker

**Características:**
- eventType: "material_uploaded", "assessment_attempt", etc.
- state: "closed", "half_open", "open"
- Labels: service, event_type, queue, target_service, state

**Constantes de nombres:**
```
MetricMessagesProcessedTotal = "messages_processed_total"
MetricMessageDuration = "message_processing_duration_seconds"
MetricMessageRetriesTotal = "message_retries_total"
MetricMessagesInQueue = "messages_in_queue"
MetricCircuitBreakerState = "circuit_breaker_state"
```

**Ejemplo:**
```go
start := time.Now()
err := processEvent(eventData)
m.RecordMessageProcessed("assessment_attempt", time.Since(start), err)

if err != nil {
    m.RecordMessageRetry("assessment_attempt", retryCount)
}

m.SetMessagesInQueue("events", pendingCount)
m.SetCircuitBreakerState("payment-service", "closed")
```

## Flujos comunes

### 1. Inicializar service con métricas

```go
package main

import (
    "github.com/EduGoGroup/edugo-shared/metrics"
)

type APIService struct {
    metrics *metrics.Metrics
}

func NewAPIService() *APIService {
    // Sin argumentos: NoopRecorder (perfecto para dev)
    m := metrics.New("api-service")

    return &APIService{metrics: m}
}

// En producción, inyectar recorder real
func NewAPIServiceWithPrometheus(promRecorder metrics.Recorder) *APIService {
    m := metrics.New("api-service", promRecorder)
    return &APIService{metrics: m}
}

func main() {
    svc := NewAPIService()
    // Métricas deshabilitadas en dev (NoOp), activadas en prod
}
```

### 2. Middleware HTTP con instrumentación

```go
package main

import (
    "net/http"
    "sync/atomic"
    "time"
    "github.com/EduGoGroup/edugo-shared/metrics"
)

func MetricsMiddleware(m *metrics.Metrics) func(http.Handler) http.Handler {
    var activeRequests int64

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Incrementar y registrar requests activos
            atomic.AddInt64(&activeRequests, 1)
            m.SetHTTPActiveRequests(int(atomic.LoadInt64(&activeRequests)))

            start := time.Now()

            // Crear wrapper para capturar status
            sw := &statusWriter{ResponseWriter: w}

            // Llamar al handler
            next.ServeHTTP(sw, r)

            // Registrar métrica
            m.RecordHTTPRequest(
                r.Method,
                r.URL.Path,
                sw.status,
                time.Since(start),
            )

            // Decrementar requests activos
            atomic.AddInt64(&activeRequests, -1)
            m.SetHTTPActiveRequests(int(atomic.LoadInt64(&activeRequests)))
        })
    }
}

type statusWriter struct {
    http.ResponseWriter
    status int
}

func (sw *statusWriter) WriteHeader(status int) {
    sw.status = status
    sw.ResponseWriter.WriteHeader(status)
}

func main() {
    m := metrics.New("http-api")
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    handler := MetricsMiddleware(m)(mux)
    http.ListenAndServe(":8080", handler)
}
```

### 3. Instrumentar CRUD genérico

```go
package main

import (
    "github.com/EduGoGroup/edugo-shared/metrics"
)

type Repository struct {
    db      Database
    metrics *metrics.Metrics
}

func (r *Repository) Create(entity string, data map[string]interface{}) error {
    start := time.Now()
    err := r.db.Insert(entity, data)
    r.metrics.RecordBusinessOperation(entity, "create", time.Since(start), err)
    return err
}

func (r *Repository) Read(entity string, id interface{}) (interface{}, error) {
    start := time.Now()
    data, err := r.db.SelectOne(entity, id)
    r.metrics.RecordBusinessOperation(entity, "read", time.Since(start), err)
    return data, err
}

func (r *Repository) Update(entity string, id interface{}, data map[string]interface{}) error {
    start := time.Now()
    err := r.db.Update(entity, id, data)
    r.metrics.RecordBusinessOperation(entity, "update", time.Since(start), err)
    return err
}

func (r *Repository) Delete(entity string, id interface{}) error {
    start := time.Now()
    err := r.db.Delete(entity, id)
    r.metrics.RecordBusinessOperation(entity, "delete", time.Since(start), err)
    return err
}
```

### 4. Instrumentar autenticación completa

```go
package main

import (
    "time"
    "github.com/EduGoGroup/edugo-shared/metrics"
)

type AuthService struct {
    metrics *metrics.Metrics
}

func (a *AuthService) Login(username, password string) (string, error) {
    start := time.Now()

    // Verificar credenciales
    user, err := validateCredentials(username, password)
    success := err == nil

    a.metrics.RecordLogin(success, time.Since(start))

    if !success {
        a.metrics.RecordRateLimitHit(username)
        return "", err
    }

    // Generar token
    token := generateToken(user.ID)
    return token, nil
}

func (a *AuthService) RefreshToken(oldToken string) (string, error) {
    start := time.Now()

    // Validar y renovar token
    newToken, err := validateAndRefresh(oldToken)
    success := err == nil

    a.metrics.RecordTokenRefresh(success, time.Since(start))

    if !success {
        return "", err
    }

    return newToken, nil
}

func (a *AuthService) Authorize(userID int, permission string) bool {
    allowed := checkPermission(userID, permission)
    a.metrics.RecordPermissionCheck(permission, allowed)
    return allowed
}
```

## Arquitectura

Flujo de instrumentación de aplicación:

```
Aplicación (HTTP handlers, Database layer, Business logic)
    ↓
Crear Metrics (con o sin recorder)
    ↓
Registrar eventos (RecordLogin, RecordHTTPRequest, etc.)
    ↓
Metrics.recordXXX() → Recorder.CounterAdd/HistogramObserve/GaugeSet
    ↓
Backend concreto (NoOp, Prometheus, Datadog, OTel)
    ↓
Sistema de monitoreo (Grafana, Datadog dashboard, alertas)
```

Modelo de componentes:

```
Metrics (punto de entrada)
├─ Auth methods
│  ├─ RecordLogin → [counter, histogram]
│  ├─ RecordTokenRefresh → [counter, histogram]
│  ├─ RecordRateLimitHit → [counter]
│  └─ RecordPermissionCheck → [counter]
├─ HTTP methods
│  ├─ RecordHTTPRequest → [counter, histogram]
│  └─ SetHTTPActiveRequests → [gauge]
├─ DB methods
│  ├─ RecordDBQuery → [counter, histogram]
│  └─ SetDBConnectionsOpen → [gauge]
├─ Business methods
│  ├─ RecordBusinessOperation → [counter, histogram]
│  ├─ RecordAssessmentAttempt → [counter, histogram]
│  ├─ RecordGrading → [counter, histogram]
│  ├─ RecordReview → [counter, histogram]
│  ├─ RecordNotification → [counter]
│  └─ RecordExport → [counter×2, histogram]
├─ Messaging methods
│  ├─ RecordMessageProcessed → [counter, histogram]
│  ├─ RecordMessageRetry → [counter]
│  ├─ SetMessagesInQueue → [gauge]
│  └─ SetCircuitBreakerState → [gauge]
└─ Recorder (interfaz)
   ├─ CounterAdd(name, value, labels)
   ├─ HistogramObserve(name, value, labels)
   └─ GaugeSet(name, value, labels)
```

## Dependencias

- **Internas**: Ninguna
- **Externas**:
  - `time` (estándar) — Para `time.Duration` y conversiones

## Testing

Suite de tests unitarios:

- Creación de Metrics con y sin recorder
- RecordLogin/TokenRefresh con success/failure
- RecordDBQuery con error handling
- RecordHTTPRequest con status codes variados
- RecordBusinessOperation con múltiples entities y operations
- RecordAssessmentAttempt con acciones diversas
- SetDBConnectionsOpen/SetHTTPActiveRequests/SetMessagesInQueue
- Verificación de labels en las llamadas al recorder
- Conversión correcta de time.Duration a segundos
- Test race detector para concurrencia

Ejecutar:
```bash
make test      # Tests unitarios
make test-race # Race detector
make check     # Validar + tests
```

## Notas de diseño

- **Agnóstico de backend**: La interfaz `Recorder` permite cualquier implementación (Prometheus, Datadog, OTel) sin cambios en código cliente
- **Zero overhead por defecto**: `NoopRecorder` incurre cero costo, perfecto para desarrollo local
- **Resiliencia**: Los métodos de Metrics no retornan error; el Recorder debe ser resiliente
- **Labels estandarizados**: Service name se agrega automáticamente a todas las métricas para multidimensionalidad
- **Duración en segundos**: Histogramas de duración se normalizan a segundos para consistencia cross-backend
- **Status automático**: success/error se deduce automáticamente del parámetro `error`
