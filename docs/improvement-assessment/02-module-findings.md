# Module Findings: oportunidades de mejora por módulo

## Metodología aplicada

Se consolidó evidencia de:

- `scripts/module-manifest.tsv` (19 módulos listados).
- estructura real de carpetas y archivos `.go`/`_test.go`.
- `PLAN_REFACTORING.md` como base de hallazgos previos.
- `docs/phase-1/module-catalog.md` y `docs/phase-2/*`.

## Señales cuantitativas rápidas

Observación del estado estructural:

- Todos los módulos del manifiesto tienen `README.md`, `CHANGELOG.md` y carpeta `docs/`.
- Solo `audit` y `audit/postgres` ya usan carpeta `internal/` de forma explícita.
- Módulos con mayor volumen de archivos Go:
  - `bootstrap` (26)
  - `common` (21)
  - `logger` (15)
  - `testing` (15)
  - `messaging/rabbit` (14)

## Hallazgos priorizados (globales)

## 1) Testing: reducir redundancia y separar unit vs integration (prioridad alta)

Mejora:

- Eliminar tests triviales/duplicados.
- Consolidar archivos de test fragmentados en módulos grandes.
- Aislar pruebas de integración con build tags (`integration`) y objetivos de make separados.

Beneficio:

- CI más rápido y menos ruido.
- Menor costo de mantenimiento de pruebas.

Riesgo:

- Bajo, siempre que se preserve cobertura funcional relevante.

Módulos más impactados:

- `common`, `logger`, `bootstrap`, `testing`, `messaging/rabbit`, `database/postgres`, `database/mongodb`, `cache/redis`.

## 2) Reestructuración interna en módulos grandes (prioridad media-alta)

Mejora:

- Mantener interfaces/factories públicas en raíz.
- Mover implementación concreta a `internal/`.
- Usar facades/wrappers públicos estables.

Módulos objetivo:

- `bootstrap` (más crítico por tamaño y mezcla de responsabilidades).
- `logger`.
- `screenconfig`.

Beneficio:

- API pública más clara.
- Menor acoplamiento accidental.
- Menor riesgo de uso indebido de detalles internos.

Riesgo:

- Medio (si cambian rutas o firmas públicas por error).

## 3) Duplicación transversal (prioridad media)

Mejora:

- Extraer componentes comunes de health-check y retry/backoff en `common`.
- Revisar wrappers de config duplicados en `bootstrap`.

Beneficio:

- Menos código repetido.
- Comportamiento operativo más consistente entre módulos.

Riesgo:

- Medio-bajo.

## Ejemplos concretos de redundancia y cómo atacarlos

## A) Patrón repetido de inicialización en `bootstrap`

Evidencia:

- `bootstrap/init_postgresql.go`
- `bootstrap/init_mongodb.go`
- `bootstrap/init_rabbitmq.go`
- `bootstrap/init_s3.go`

Patrón repetido en los cuatro archivos:

1. validar factory,
2. extraer config,
3. log de inicio,
4. crear recurso,
5. registrar cleanup (si aplica),
6. log de éxito,
7. manejar error con mensaje similar.

Cómo atacarlo:

- Crear un orquestador interno por etapas en `bootstrap/internal/initializers`.
- Mantener funciones públicas intactas; internamente delegar a una plantilla común (pipeline) para resource init.
- Estandarizar contrato de inicialización y reporte de errores por recurso.

Impacto consumidor:

- **Sin impacto** si no se cambia `Bootstrap(...)` ni los tipos `Factories`/`Resources`.

## B) Registro de cleanup por tipo con misma estructura

Evidencia:

- `bootstrap/cleanup_registrars.go`: `registerPostgreSQLCleanup`, `registerMongoDBCleanup`, `registerRabbitMQCleanup`.

Redundancia:

- Las 3 funciones hacen el mismo flujo: cast de lifecycle manager, guard clauses, cast de recurso concreto y `RegisterSimple(...)`.

Cómo atacarlo:

- Introducir helper interno para registro seguro (`registerTypedCleanup`) con callbacks de cast/close.
- Dejar funciones específicas como wrappers para mantener legibilidad.

Impacto consumidor:

- **Sin impacto** (interno a `bootstrap`).

## C) Config wrappers duplicados respecto a módulos fuente

Evidencia:

- `bootstrap/interfaces.go`: `PostgreSQLConfig`, `MongoDBConfig`, `RabbitMQConfig`.
- Módulos fuente:
  - `database/postgres/config.go`
  - `database/mongodb/config.go`
  - `messaging/rabbit/config.go`

Redundancia:

- Bootstrap mantiene versiones simplificadas de config que duplican semántica de campos base.

Cómo atacarlo (gradual):

- Fase 1: mantener wrappers (sin ruptura), documentar explícitamente mapping.
- Fase 2: agregar adaptadores tipados `bootstrap -> module config` en `internal/adapters`.
- Fase 3 (opcional): evaluar convergencia a tipos compartidos si el acoplamiento es aceptable.

Impacto consumidor:

- **Compatible** mientras `bootstrap` mantenga su contrato actual.
- **Breaking potencial** solo si se eliminan wrappers públicos.

## D) Lógica de timeout de health check repetida

Evidencia:

- `database/postgres/connection.go`: `DefaultHealthCheckTimeout` + `HealthCheck`.
- `database/mongodb/connection.go`: `DefaultHealthCheckTimeout` + `HealthCheck`.
- `bootstrap/health_check.go`: secuencia de checks con patrón homogéneo.

Cómo atacarlo:

- Definir `common/health` con interfaz mínima de checker y timeout default.
- Migrar por adopción progresiva, empezando con adaptadores internos sin tocar API pública.

Impacto consumidor:

- **Compatible**.

## E) Backoff de reconexión aislado en `messaging/rabbit`

Evidencia:

- `messaging/rabbit/connection.go`: `nextBackoff(...)` privado para reconexión.

Oportunidad:

- Reutilizable para otros módulos con reconexión (futura infraestructura).

Cómo atacarlo:

- Extraer a `common/retry/backoff.go` y mantener wrapper local durante transición.

Impacto consumidor:

- **Compatible**.

## 4) Dead code exportado en enums y helpers (prioridad media)

Mejora:

- Evaluar eliminación de funciones `All*()` ya marcadas como no consumidas en análisis previo.

Beneficio:

- API pública más pequeña y más fácil de mantener.

Riesgo:

- Medio-alto si hay consumidores externos no inventariados.

## Priorización sugerida por módulo

| Módulo | Oportunidad principal | Prioridad | Riesgo |
| --- | --- | --- | --- |
| `bootstrap` | Reestructurar API pública vs implementación + consolidar tests | Alta | Medio |
| `logger` | Reducir tests triviales y separar implementaciones internas | Alta | Medio |
| `common` | Limpiar dead code exportado y mantener tests funcionales | Alta | Medio |
| `messaging/rabbit` | Separar unit/integration y consolidar cobertura útil | Alta | Bajo |
| `testing` | Eliminar tests inválidos/triviales y reforzar tests de integración válidos | Alta | Bajo |
| `screenconfig` | Reordenar responsabilidades por concern | Media | Medio |
| `database/postgres` | Mejorar split unit/integration y evitar duplicación | Media | Bajo |
| `database/mongodb` | Igual que postgres, con foco en timeouts/health | Media | Bajo |
| `repository` | Expandir cobertura en comportamiento CRUD/integración | Media | Bajo |
| `audit/postgres` | Aumentar cobertura dedicada del módulo | Media | Bajo |
| `config` | Mantener estable; mejoras menores de pruebas y ergonomía | Baja | Bajo |
| `auth` | Mantener estable; hardening incremental | Baja | Medio (por criticidad) |
| `middleware/gin` | Mantener API estable y ampliar casos límite | Baja | Medio (por uso masivo) |
| `cache/redis` | Ajustes de testing y observabilidad | Baja | Bajo |
| `lifecycle` | Bajo costo de mantenimiento, mantener simple | Baja | Bajo |
| `metrics` | Revisar adopción real y estandarizar uso | Baja | Bajo |
| `export` | Revisar propósito y adopción (podría estar subutilizado) | Baja | Bajo |
| `messaging/events` | Mantener compatibilidad de schemas | Media | Medio |
| `audit` | Mantener contrato estable y cobertura de edge cases | Media | Bajo |

## Decisiones guardrail recomendadas

- No cambiar firmas públicas de módulos críticos en fases iniciales.
- Evitar cambios breaking concurrentes en `auth`, `middleware/gin` y `repository`.
- Asegurar changelog por módulo en cada release de refactor.
