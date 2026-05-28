# AGENTS.md — edugo-shared

> Detalle local. Reglas globales del ecosistema en `../../AGENTS.md` (no las repitas).
> Norte actual del proyecto en `docs/ACTIVE.md`. Mapa de módulos y arquitectura en `README.md`
> y `docs/` (phase-1/2/3, roadmap).

## Propósito

**Librería compartida multi-módulo** en Go que consumen las 4 APIs y el worker. No es un servicio:
no tiene `cmd/` ni se "levanta". Es un repo con ~17 módulos Go independientes (cada uno con su `go.mod`,
`README.md`, `CHANGELOG.md` y `docs/`), versionados y publicados por separado.

## Módulos (foco de cada uno)

| Módulo | Foco |
| --- | --- |
| `common` | Base: env, errores, validator, UUID, enums (incl. `enum.Permission*`). |
| `logger` | Logging estructurado (interfaz + backends Zap/Logrus/slog). |
| `config` | Carga de configuración desde archivo + entorno. |
| `auth` | JWT con contexto activo, passwords, refresh tokens, blacklist. |
| `middleware/gin` | Auth, `RequirePermission`, contexto, CORS, request logging, auditoría para Gin. |
| `screenconfig` | Transformación/validación de configuración SDUI (la usa `edugo-api-platform`). |
| `audit`, `audit/postgres` | Contrato de eventos auditables + persistencia GORM. |
| `bootstrap`, `lifecycle` | Inicialización ordenada de recursos + cleanup LIFO. |
| `database/postgres`, `database/mongodb` | Conexión, pool, transacciones. |
| `cache/redis` | Conexión Redis + caché JSON genérico. |
| `messaging/events`, `messaging/rabbit` | Schemas de eventos de dominio + RabbitMQ (publish/consume/DLQ). |
| `repository` | Repos GORM compartidos (users, schools, memberships). |
| `metrics`, `tracer`, `health`, `resilience`, `export` | Observabilidad y utilidades transversales. |
| `testing` | Testcontainers reutilizables + helpers de integración. |

## Cómo construir y testear

El `Makefile` **raíz** orquesta los módulos por niveles de dependencia (foundation → runtime) y separa
el set de integración. Cada módulo con `go.mod` usa el mismo contrato vía `scripts/module-common.mk`.
El inventario único vive en `scripts/module-manifest.tsv` (alimenta Makefile, coverage y CI).
- Por módulo: `make -C <modulo> test` / `lint` / `fmt` / `release-check`.
- Raíz: orquesta build/test/coverage de todos.

## Para agregar / cambiar algo

- **Cambio interno a un módulo**: edítalo, corre su `make test`/`release-check`, actualiza su
  `CHANGELOG.md`.
- **Módulo nuevo**: crea su carpeta con `go.mod`, `README.md`, `CHANGELOG.md`, `docs/`, y regístralo en
  `scripts/module-manifest.tsv` para que entre al orquestador, coverage y release.
- El `go.work` del ecosistema (en `EduBack/`) incluye todos los módulos para integración local sin release.

## Convenciones y gotchas locales

- **Multi-módulo, no monolito**: cada `go.mod` es independiente; un consumidor puede fijar versiones
  distintas por módulo. No asumas que todo comparte una sola versión.
- **`enum.Permission*`** (en `common`) es la fuente compartida de permisos que usan las APIs con
  `RequirePermission`. Cambiarlo impacta a todos los servicios.
- **Releases**: tags raíz (`vX.Y.Z`) o modulares (`modulo/vX.Y.Z`); `release.yml` usa el `CHANGELOG.md`
  correcto. No mezcles cambios de varios módulos sin pensar el tag.
- `middleware/gin` evita importar `metrics` directamente: la conexión se hace por duck-typing desde el
  consumidor (ver `SetPermissionMetricsRecorder` en las APIs).
- Reglas globales: código en inglés, logs/docs en español, fechas UTC.
