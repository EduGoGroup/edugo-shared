# Plan de Refactoring - edugo-shared

## Resumen Ejecutivo

Tras analizar los 17 modulos del repositorio, se identificaron 3 ejes de mejora:
1. **Estructura de modulos**: 3 modulos necesitan reorganizacion interna
2. **Calidad de tests**: ~35 archivos de test son triviales/padding (35-40% del total)
3. **Duplicacion de codigo**: Patrones repetidos en 4+ modulos que pueden centralizarse

Se realizo un analisis de cobertura para cada test candidato a eliminacion.
La conclusion es que la mayoria de los tests triviales ya tienen cobertura
desde tests funcionales/integration existentes. Donde hay gap, se propone
mitigacion especifica.

Se verifico que las funciones `All*()` marcadas como dead code NO son usadas
por ningun consumidor externo (iam-platform, admin-api, mobile-api, worker, apple_new).

**Todo el refactoring viaja en un solo PR.**
**Todos los modulos afectados se versionan como v0.100.0.**

---

## Versionamiento y Documentacion

Todos los modulos afectados por este refactoring deben:

1. **Actualizar CHANGELOG.md** — Agregar entrada para v0.100.0 con:
   - Dead code eliminado (funciones `All*()`)
   - Tests triviales eliminados/consolidados
   - Reestructuracion interna (si aplica en Fase 2)
   - Reduccion de duplicacion (si aplica en Fase 3)

2. **Actualizar README.md** — Reflejar:
   - Nueva estructura de directorios (Fase 2)
   - Funciones publicas removidas (dead code)
   - Cambios en patrones de uso (si los hay)

3. **Actualizar docs/README.md** — Documentacion tecnica alineada con los cambios

4. **Version**: Todos los modulos afectados se publican como **v0.100.0**

### Modulos que requieren actualizacion de docs

| Modulo | CHANGELOG | README | docs/README | Motivo |
|--------|:---------:|:------:|:-----------:|--------|
| common | SI | SI | SI | Dead code eliminado en enums |
| auth | SI | NO | NO | Test trivial eliminado |
| config | SI | NO | NO | Test trivial eliminado |
| logger | SI | SI | SI | Tests consolidados, estructura interna (Fase 2) |
| bootstrap | SI | SI | SI | Tests consolidados, estructura interna (Fase 2) |
| screenconfig | SI | SI | SI | Estructura interna (Fase 2) |
| database/postgres | SI | NO | NO | Test trivial eliminado |
| database/mongodb | SI | NO | NO | Ajuste menor en tests |
| messaging/rabbit | SI | NO | NO | Test trivial eliminado |
| testing | SI | NO | NO | Test invalido eliminado |
| cache/redis | SI | NO | NO | CI split (build tags) |
| messaging/events | SI | NO | NO | CI split (build tags) |

---

## FASE 0: Limpieza de Tests (con analisis de cobertura)
**Riesgo: BAJO | Esfuerzo: BAJO | Impacto en CI: ALTO**

### 0.1 Tests de enums — Accion por archivo

#### `common/types/enum/role_test.go`
- **ELIMINAR**: `TestAllSystemRoles`, `TestAllSystemRolesStrings` — funciones sin uso externo (dead code)
- **MANTENER**: `TestSystemRole_IsValid`, `TestSystemRole_String` — cobertura unica de validacion
- **Cobertura**: `.String()` e `.IsValid()` se cubren parcialmente desde middleware/gin/permission_auth_test.go
- **Gap**: ~0% — funciones "All" son dead code

#### `common/types/enum/status_test.go`
- **ELIMINAR**: `TestAllMaterialStatuses`, `TestAllProgressStatuses`, `TestAllProcessingStatuses` — funciones sin uso externo
- **MANTENER**: Tests de `.IsValid()` y `.String()` — sin cobertura externa
- **Gap**: ~0% — funciones "All" son dead code
- **ACCION EXTRA**: Evaluar si `AllMaterialStatuses()`, `AllProgressStatuses()`, `AllProcessingStatuses()` deben eliminarse del codigo fuente (dead code)

#### `common/types/enum/event_test.go`
- **ELIMINAR**: `TestAllEventTypes` — funcion sin uso externo
- **MANTENER**: `TestEventType_GetRoutingKey` — **CRITICO**, es la unica cobertura de routing keys para RabbitMQ
- **MANTENER**: `TestEventType_IsValid`, `TestEventType_String` — sin cobertura externa
- **Gap**: 0% si se mantiene GetRoutingKey

#### `common/types/enum/permission_test.go`
- **ELIMINAR**: Subtests redundantes de `TestAllPermissionsSlice` (lineas 86-147)
- **MANTENER**: Tests de integridad del mapa de permisos (lineas 149-307) — validan consistencia
- **MANTENER**: `.String()` e `.IsValid()` — cubiertos parcialmente por middleware pero el test de integridad es unico
- **Gap**: ~0% — middleware tests cubren los paths funcionales

#### `common/types/enum/assessment_test.go`
- **ELIMINAR**: `TestAllAssessmentTypes` — funcion sin uso externo
- **MANTENER**: Tests de `.IsValid()` y `.String()` — sin cobertura externa
- **Gap**: 0%
- **ACCION EXTRA**: Evaluar si `AllAssessmentTypes()` es dead code

### 0.2 Tests triviales individuales — Accion por test

| Test | Accion | Motivo | Gap |
|------|--------|--------|-----|
| `auth/blacklist_test.go::TestNoOpBlacklist` | ELIMINAR | Cubierto por middleware/gin/jwt_auth_test.go | 0% |
| `database/postgres/connection_unit_test.go::TestDefaultConstants_Unit` | ELIMINAR | Cubierto por config_test.go::TestDefaultConfig | 0% |
| `database/mongodb/connection_unit_test.go::TestDefaultConstants_Unit` | MANTENER | Unica cobertura de DefaultHealthCheckTimeout y DefaultDisconnectTimeout | ~2% |
| `config/base_test.go` — tests de struct fields | ELIMINAR | Trivial, cubierto por loader_test.go | 0% |
| `config/base_test.go` — TestConnectionString | MANTENER | Logica real de construccion DSN, sin cobertura externa | ~2% |
| `messaging/rabbit/publisher_unit_test.go` (4 tests) | ELIMINAR | 100% cubierto por publisher_test.go (integration) | 0% |
| `testing/containers/helpers_unit_test.go` | ELIMINAR | Tests INVALIDOS (pasan db=nil), cubiertos por helpers_test.go | 0% |
| `testing/containers/options_test.go` | MANTENER | Builder pattern bien testeado, unica cobertura de defaults | ~5% |
| `audit/audit_logger_test.go` (opciones) | MANTENER | Unica cobertura de WithSeverity, WithMetadata, WithError (nil check) | ~8% |

### 0.3 Logger tests — Accion por archivo

| Archivo | Accion | Motivo | Gap |
|---------|--------|--------|-----|
| `logger/fields_test.go` — 15 tests de setters | ELIMINAR | Cubiertos por middleware/gin/request_logging_test.go | 0% |
| `logger/fields_test.go` — benchmarks triviales | ELIMINAR | BenchmarkTypedHelpers, BenchmarkWithRequestID (microsegundos) | 0% |
| `logger/fields_test.go` — BenchmarkNewSlogProvider | MANTENER | Benchmark util de inicializacion | 0% |
| `logger/context_test.go` | REDUCIR | Eliminar tests basicos, mantener TestContext_NestedOverwrite | ~5% |
| `logger/context_test.go` mitigacion | AGREGAR | 1-2 lineas en request_logging_test.go para cubrir nested context | 0% |
| `logger/logger_test.go` | MANTENER | Tests funcionales de NewZapLogger con niveles/formatos — NO son triviales | 0% |
| `logger/logrus_logger_test.go` | REDUCIR | Merge TestNewLogrusLogger trivial en TestLogrusLogger_Levels existente | 0% |
| `logger/slog_adapter_test.go` | REDUCIR | Eliminar TestSlogAdapter_SlogLogger (getter trivial), mantener resto | ~3% |
| `logger/slog_provider_test.go` | MANTENER | Tests funcionales de configuracion — NO son triviales | 0% |
| `logger/zap_logger_test.go` | REDUCIR | Eliminar TestNewZapLogger duplicado (ya existe en logger_test.go) | 0% |

### 0.4 Resumen de impacto Fase 0

```
Tests eliminados:    ~25-30 funciones de test
Lineas removidas:    ~350-400
Cobertura perdida:   <3% (mitigada con 5-10 lineas en tests existentes)
Tiempo CI ahorrado:  ~30-40s
Dead code detectado: AllSystemRoles(), AllMaterialStatuses(), AllProgressStatuses(),
                     AllProcessingStatuses(), AllAssessmentTypes(), AllEventTypes()
                     → Candidatos a eliminar del codigo fuente
```

---

## FASE 1: Consolidacion de Tests y CI Split
**Riesgo: BAJO | Esfuerzo: MEDIO | Impacto en CI: ALTO**

### 1.1 Consolidar bootstrap tests (5 archivos -> 2)
- [ ] Merge `cleanup_test.go` + `options_test.go` + `resources_test.go` en `bootstrap_test.go`
- [ ] Mantener `bootstrap_integration_test.go` separado (requiere Docker)
- [ ] Eliminar escenarios duplicados que ya se prueban en bootstrap_test.go

### 1.2 Eliminar `testing/containers/helpers_unit_test.go`
- Tests invalidos que pasan db=nil
- Cobertura real esta en `helpers_test.go` (integration)

### 1.3 Consolidar screenconfig tests
- [ ] Revisar 4 archivos de test y eliminar tests que validan JSON marshaling trivial
- [ ] Mantener tests de logica de negocio (platform overrides, slot resolution, GetRoutingKey)

### 1.4 Separar CI: unit vs integration
- [ ] Crear `make test-unit` (sin Docker, rapido)
- [ ] Crear `make test-integration` (con Docker, lento)
- [ ] Aplicar build tags `//go:build integration` a tests que requieren Docker:
  - `bootstrap/factory_postgresql_integration_test.go`
  - `bootstrap/factory_mongodb_integration_test.go`
  - `bootstrap/factory_rabbitmq_integration_test.go`
  - `database/postgres/connection_test.go`
  - `database/mongodb/mongodb_integration_test.go`
  - `cache/redis/cache_test.go`
  - `messaging/rabbit/publisher_test.go`
  - `messaging/rabbit/connection_test.go`
  - `testing/containers/*_test.go` (postgres, mongodb, rabbitmq)

**Resultado esperado**: CI unit tests en <10s, integration tests aislados ~50-100s

---

## FASE 2: Reestructuracion de Modulos (Interfaz publica + implementacion interna)
**Riesgo: MEDIO | Esfuerzo: ALTO | Impacto: Calidad del codigo**

Principio: Las interfaces y tipos publicos se mantienen en la raiz del modulo.
La implementacion se mueve a `internal/`. Los consumidores no se ven afectados
porque importan las interfaces/factories, no las implementaciones directas.

Modulos modelo (ya bien estructurados): audit/, cache/redis/, common/, database/, messaging/

### 2.1 Bootstrap (18 archivos -> 5 raiz + internal/)
**Prioridad: ALTA** - Es el modulo mas grande y desordenado

Estructura objetivo:
```
bootstrap/
  bootstrap.go           # Bootstrap() - orquestador principal
  resources.go           # Resources struct + helpers
  options.go             # BootstrapOption pattern
  interfaces.go          # Factory interfaces (publicas)
  internal/
    initializers/
      logger.go          # (era init_logger.go)
      postgresql.go      # (era init_postgresql.go)
      mongodb.go         # (era init_mongodb.go)
      rabbitmq.go        # (era init_rabbitmq.go)
      s3.go              # (era init_s3.go)
    factories/
      logger.go          # (era factory_logger.go)
      postgresql.go      # (era factory_postgresql.go)
      mongodb.go         # (era factory_mongodb.go)
      rabbitmq.go        # (era factory_rabbitmq.go)
      s3.go              # (era factory_s3.go)
    health/
      checker.go         # (era health_check.go)
    cleanup/
      registrars.go      # (era cleanup_registrars.go)
    config/
      extractors.go      # (era config_extractors.go)
```

Pasos:
- [ ] Crear estructura internal/
- [ ] Mover archivos init_*.go a internal/initializers/
- [ ] Mover archivos factory_*.go a internal/factories/
- [ ] Mover helpers, health_check, cleanup a internal/
- [ ] Actualizar imports en bootstrap.go
- [ ] Verificar que la API publica no cambie
- [ ] Correr todos los tests

### 2.2 Logger (9 archivos -> 3 raiz + internal/)
**Prioridad: MEDIA**

Alternativa sin breaking change (recomendada):
```
logger/
  logger.go              # Interface Logger (se mantiene)
  fields.go              # Constantes de campos (se mantiene)
  context.go             # Context helpers (se mantiene)
  slog_provider.go       # Factory publica (delegando a internal/slog)
  logrus_logger.go       # Factory publica (delegando a internal/logrus)
  zap_logger.go          # Factory publica (delegando a internal/zap)
  internal/
    slog/
      provider.go
      adapter.go
    logrus/
      logger.go
    zap/
      logger.go
```

> Los factories publicas se mantienen como wrappers finos.
> Los consumidores siguen usando `logger.NewZapLogger()` sin cambios.

- [ ] Crear internal/ con implementaciones
- [ ] Reducir archivos raiz a wrappers delegando a internal/
- [ ] Verificar que la API publica no cambie
- [ ] Correr todos los tests

### 2.3 Screenconfig (10 archivos -> 3 raiz + internal/)
**Prioridad: MEDIA**

Estructura objetivo:
```
screenconfig/
  types.go               # Enums: Pattern, ScreenType, Platform (se mantiene)
  dto.go                 # DTOs publicos (se mantiene)
  screenconfig.go        # Funciones publicas facade
  internal/
    menu/
      tree.go            # (era menu_tree.go)
    permissions/
      checker.go         # (era permissions.go)
    slots/
      resolver.go        # (era slots.go)
    overrides/
      platform.go        # (era platform_overrides.go)
    validation/
      validator.go       # (era validation.go)
```

- [ ] Crear interfaces publicas para cada concern
- [ ] Mover implementaciones a internal/
- [ ] Exponer funciones facade en screenconfig.go
- [ ] Verificar que la API publica no cambie

---

## FASE 3: Reduccion de Duplicacion entre Modulos
**Riesgo: MEDIO-BAJO | Esfuerzo: BAJO | Impacto: Mantenibilidad**

### 3.1 Extraer health check comun a common/
**4 modulos implementan el mismo patron de health check con timeout de 5s**

- [ ] Crear `common/health/checker.go`:
  ```go
  type Checker interface {
      Check(ctx context.Context) error
  }
  const DefaultTimeout = 5 * time.Second
  ```
- [ ] database/mongodb, database/postgres, cache/redis, messaging/rabbit implementan la interfaz
- [ ] No es obligatorio migrar los existentes — disponible para nuevos modulos

### 3.2 Extraer retry/backoff a common/
**messaging/rabbit tiene logica de backoff que podria ser reutilizable**

- [ ] Crear `common/retry/backoff.go`:
  ```go
  func NextBackoff(current, max time.Duration) time.Duration
  ```
- [ ] messaging/rabbit consume de common/retry
- [ ] Disponible para futuros modulos

### 3.3 Eliminar config wrappers duplicados en bootstrap
**bootstrap/interfaces.go tiene PostgreSQLConfig, MongoDBConfig, RabbitMQConfig
que son copias de los tipos originales en database/ y messaging/**

- [ ] Evaluar si bootstrap puede usar directamente los tipos de database/postgres, etc.
- [ ] Si la dependencia de version es aceptable, eliminar los wrappers
- [ ] Si no, documentar por que existen como wrappers

### 3.4 Limpiar dead code detectado en Fase 0
Funciones exportadas sin ningun uso en el repositorio:
- [ ] `AllSystemRoles()` / `AllSystemRolesStrings()` en enum/role.go
- [ ] `AllMaterialStatuses()` / `AllProgressStatuses()` / `AllProcessingStatuses()` en enum/status.go
- [ ] `AllAssessmentTypes()` en enum/assessment.go
- [ ] `AllEventTypes()` en enum/event.go

> VERIFICADO: Ninguno de los consumidores externos (iam-platform, admin-api,
> mobile-api, edugo-worker, apple_new) usa estas funciones. Son dead code confirmado.
> Seguro eliminar del codigo fuente.

---

## Resumen de Impacto por Fase

| Fase | Archivos | Riesgo | Esfuerzo | Cobertura | Beneficio |
|------|----------|--------|----------|-----------|-----------|
| 0 | ~15 test files | Bajo | 1-2 dias | -3% max (mitigable) | -400 lineas, -40s CI |
| 1 | ~10 test files | Bajo | 2-3 dias | 0% | CI split, -duplicacion |
| 2 | ~30 archivos | Medio | 5-7 dias | 0% | Estructura profesional |
| 3 | ~10 archivos | Bajo | 2-3 dias | 0% | Menos duplicacion |

## Orden de Ejecucion Recomendado

```
Fase 0 (limpieza tests) -> Fase 1 (consolidar + CI split) -> Fase 2 (reestructurar) -> Fase 3 (dedup)
```

Todo viaja en un solo PR con version v0.100.0 para todos los modulos afectados.
Se ejecutan las fases en orden dentro del mismo PR.
Fase 0 y 1 son de bajo riesgo y dan resultados inmediatos.
Fase 2 requiere mas cuidado pero es donde esta el mayor valor arquitectonico.
Fase 3 es opcional/incremental.

Al finalizar todas las fases, actualizar CHANGELOG.md, README.md y docs/README.md
de cada modulo afectado antes de crear el release v0.100.0.
