# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/metrics` se registran aquí.

## [0.100.0] - 2026-04-02

### Added

- **Metrics**: Fachada central agnóstica para registrar métricas (inicialización simple: `New("service")`)
- **Recorder**: Interfaz extensible para backends (Prometheus, Datadog, OpenTelemetry, custom)
- **NoopRecorder**: Implementación por defecto con cero overhead (segura para concurrencia)
- **Auth metrics**: RecordLogin, RecordTokenRefresh (con histogramas de duración), RecordRateLimitHit, RecordPermissionCheck
- **HTTP metrics**: RecordHTTPRequest (método, path, status, duración), SetHTTPActiveRequests
- **Database metrics**: RecordDBQuery (dbType, operation, table, duración, error), SetDBConnectionsOpen
- **Business metrics**: RecordBusinessOperation (entity/operation CRUD), RecordAssessmentAttempt, RecordGrading, RecordReview, RecordNotification, RecordExport
- **Messaging metrics**: RecordMessageProcessed, RecordMessageRetry, SetMessagesInQueue, SetCircuitBreakerState
- Suite completa de tests unitarios sin dependencias externas
- Documentación técnica detallada en docs/README.md con componentes, métodos especializados, flujos comunes y middlewares
- Makefile con targets: fmt, vet, lint, test, build, check

### Design Notes

- **Agnóstico de backend**: La interfaz `Recorder` permite futuras implementaciones (Prometheus, Datadog, OTel) sin cambios en código cliente
- **Zero overhead por defecto**: `NoopRecorder` incurre cero costo, perfecto para desarrollo local y testing
- **Resiliencia garantizada**: Los métodos de Metrics no retornan error; el Recorder debe ser resiliente
- **Labels estandarizados**: Service name se agrega automáticamente a todas las métricas para multidimensionalidad
- **Duración en segundos**: Histogramas de duración se normalizan a segundos (time.Duration → float64 con .Seconds())
- **Status automático**: success/error se deduce automáticamente del parámetro error en métodos de duración

## [0.2.0] - 2026-03-26

### Added

- Modulo completo de metricas con fachada `Metrics` y `NoopRecorder` por defecto.
- Metricas de autenticacion: login, token refresh (con histograma de duracion), rate limit, permission checks.
- Metricas de operaciones de negocio: CRUD con duracion y status.
- Metricas HTTP: request duration y response status.
- Metricas de base de datos: query duration y status.
- Metricas de messaging: message processing, DLQ routing, circuit breaker.
- Interfaz `Recorder` extensible para futuros backends (Prometheus, Datadog, OTel).
- Makefile con targets estandar (fmt, vet, lint, test, build).

