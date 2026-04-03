# Roadmap de Modulos Compartidos — edugo-shared

Inventario de codigo duplicado o reutilizable detectado en los servicios.
Cada entrada indica donde vive hoy, que servicios lo usan, y el estado de extraccion.

> Actualizar este documento conforme se vayan creando modulos.

---

## Estado

| Icono | Significado |
|-------|-------------|
| Done | Ya existe en shared |
| Ready | Listo para extraer, sin bloqueos |
| Planned | Decidido, pendiente de implementar |
| Evaluate | Requiere analisis adicional |

---

## 1. Middleware HTTP (Gin)

Modulo: `middleware/gin`

| Componente | Estado | Origen | Servicios que lo usan | Notas |
|------------|--------|--------|-----------------------|-------|
| CORS Middleware | Done | Mobile (mejor), IAM, Admin | IAM, Admin, Mobile | PR #131. Corrige bug Admin (credentials fuera de origin check). Agrega `Vary: Origin` |
| BindJSON | Done | Mobile (mejor), IAM, Admin | IAM, Admin, Mobile | PR #131. Distingue length vs value en min/max |
| JWT Auth | Done | — | IAM | `jwt_auth.go` ya en shared |
| Permission Auth | Done | — | IAM, Admin | `permission_auth.go` ya en shared |
| Request Logging | Done | — | IAM, Admin, Mobile | `request_logging.go` ya en shared |
| Audit Middleware | Done | — | IAM, Admin | `audit.go` ya en shared |
| Error Handler | Done | IAM, Admin, Mobile | IAM, Admin, Mobile | PR #131. Combina panic recovery + c.Errors + HandleError(). Usa GetLogger(c) para logs correlacionados |
| Remote Auth | Done | Admin, Mobile | Admin, Mobile | PR #131. Extraido de Admin/Mobile (identicos). IAM usa JWT Auth directo |
| Metrics (Prometheus) | Evaluate | Mobile | Mobile | Counters y histograms HTTP. Util para estandarizar metricas en todos los servicios |

---

## 2. Clientes HTTP

Modulo propuesto: `client/auth` y `client/iam`

| Componente | Estado | Origen | Servicios que lo usan | Notas |
|------------|--------|--------|-----------------------|-------|
| Auth Client | Done | Admin, Mobile, Worker | Admin, Mobile (Worker usa version propia con circuit breaker) | PR #131. Vive en `middleware/gin/auth_client.go`. JWT local + fallback remoto + cache SHA256+TTL |
| IAM Client | Evaluate | Admin, Mobile | Admin, Mobile | HTTP client para roles y permisos (grant, revoke, get user roles). Admin y Mobile tienen versiones similares |

---

## 3. Resiliencia

Modulo propuesto: `resilience/`

| Componente | Estado | Origen | Servicios que lo usan | Notas |
|------------|--------|--------|-----------------------|-------|
| Circuit Breaker | Ready | Worker | Worker (NLP client) | Patron Closed/Open/HalfOpen con conteo de fallos, configurable. Reutilizable para cualquier llamada service-to-service |
| Rate Limiter | Ready | Worker | Worker | Token bucket + multi-limiter por entidad. Util para APIs (throttling de requests) |
| Retry Policy | Ready | Worker | Worker (processors) | Backoff exponencial con clasificacion de errores (transient vs permanent). Aplica a cualquier operacion con reintentos |

---

## 4. Lifecycle / Infraestructura

Modulo propuesto: `lifecycle/` o directamente en modulos existentes

| Componente | Estado | Origen | Servicios que lo usan | Notas |
|------------|--------|--------|-----------------------|-------|
| Graceful Shutdown | Ready | Worker | Worker | Orquestacion ordenada de cleanup con timeout. Escucha SIGTERM/SIGINT. Todos los servicios podrian beneficiarse |
| Health Check Framework | Ready | Worker | Worker | Interface extensible con checks para Postgres, MongoDB, RabbitMQ. Agrega checks por composicion. APIs actualmente tienen health handlers manuales |
| Metrics Server | Evaluate | Worker | Worker | HTTP server para /metrics (Prometheus) y /health. Podria combinarse con Health Check Framework |

---

## 5. DTOs / Respuestas HTTP

Modulo propuesto: `http/dto` o dentro de `middleware/gin`

| Componente | Estado | Origen | Servicios que lo usan | Notas |
|------------|--------|--------|-----------------------|-------|
| ErrorResponse | Evaluate | IAM, Admin, Mobile | IAM, Admin, Mobile | Struct casi identico en los 3: Code, Message, Details. Ya existe `AppError` en `common/errors` — evaluar si agregar DTOs de response HTTP o si el error handler compartido es suficiente |
| SuccessResponse | Evaluate | IAM, Admin, Mobile | IAM, Admin, Mobile | Wrapper generico de respuesta exitosa. Similar en los 3 |
| PaginatedResponse | Evaluate | IAM, Admin, Mobile | IAM, Admin, Mobile | Paginacion con metadata (Page, PerPage, Total, TotalPages). `ListFilters` ya esta en shared — complementar con el response DTO |

---

## 6. Procesamiento de documentos

Modulo propuesto: `processing/`

| Componente | Estado | Origen | Servicios que lo usan | Notas |
|------------|--------|--------|-----------------------|-------|
| PDF Extractor | Evaluate | Worker | Worker | Extraccion de texto + metadata (page count, word count, deteccion de scanned). Potencialmente util si Admin necesita procesar PDFs |
| NLP Client Interface | Evaluate | Worker | Worker | Interface para generar resumenes, quizzes, extraer secciones. Si se integra NLP en otros servicios seria reutilizable |

---

## Prioridad sugerida

### Ola 1 (alto impacto, bajo riesgo)
1. **Error Handler** — 3 variaciones, unificar elimina inconsistencias en manejo de errores
2. **Remote Auth** — identico en Admin y Mobile, facil de extraer
3. **Auth Client** — 3 servicios lo duplican, patron de caching reutilizable

### Ola 2 (infraestructura de worker reutilizable)
4. **Circuit Breaker** — patron de resiliencia estandar
5. **Rate Limiter** — util cuando las APIs necesiten throttling
6. **Graceful Shutdown** — todos los servicios se benefician
7. **Health Check Framework** — estandarizar health endpoints

### Ola 3 (evaluar necesidad real)
8. **Response DTOs** — evaluar si el error handler compartido hace esto innecesario
9. **Metrics Middleware** — estandarizar metricas HTTP
10. **IAM Client** — evaluar si Admin y Mobile divergen demasiado
11. **PDF / NLP** — solo si hay demanda fuera del worker

---

## Decisiones de diseno

- **No crear modulos especulativos**: solo extraer cuando hay 2+ consumidores reales o demanda clara
- **Thin wrappers en consumidores**: para `bindJSON` se uso un wrapper de 1 linea para evitar tocar 40+ handlers. Aplicar el mismo patron cuando el refactor de imports sea masivo
- **Config structs sin env tags**: shared define structs planos, cada servicio mantiene sus propias configs con tags de `env` y convierte al invocar
- **Patron de sub-modulos**: cada modulo con su propio `go.mod` para evitar arrastrar dependencias innecesarias (leccion del refactoring de bootstrap)
