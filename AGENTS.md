# AGENTS.md — edugo-shared

> Detalle local. Reglas globales del ecosistema en `../../AGENTS.md` (no las repitas).
> Norte actual del proyecto en `docs/ACTIVE.md`. Mapa de módulos y arquitectura en `README.md`
> y `docs/` (phase-3, roadmap; las fases 1 y 2 son históricas y viven en `docs/archive/`).

## Propósito

**Librería compartida multi-módulo** en Go que consumen las 4 APIs y el worker. No es un servicio:
no tiene `cmd/` ni se "levanta". Es un repo con 33 módulos Go independientes (cada uno con su `go.mod`,
`README.md`, `CHANGELOG.md` y `docs/`), versionados y publicados por separado.

## Módulos (foco de cada uno)

| Módulo | Foco |
| --- | --- |
| `common` | Base: env, errores, validator, UUID, enums (incl. `enum.Permission*`). |
| `logger` | Logging estructurado (interfaz + backends Zap/Logrus/slog). |
| `config` | Carga de configuración desde archivo + entorno. |
| `auth` | Context Token ES256 con contexto activo (+ `auth/authtest` para tests), JWT M2M, passwords. |
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
- Gates antes de un PR, desde la raíz y con `GOWORK=off`: `make ci-local` (fmt-check + vet + lint + tests +
  `make vuln`). `make vuln` corre `govulncheck` (versión en `GOVULNCHECK_VERSION` del `Makefile`, espejo del
  `env` de `ci.yml`) en los 33 módulos de `scripts/list-modules.sh --set all`; la CI lo corre en la matriz
  `validate` y **falla** si hay una vulnerabilidad alcanzable (plan 054, H.2; Dependabot apagado).
- **Integración: gate local obligatorio, no CI** (decisión del dueño, 2026-11-14; plan 054, S.6). Antes de
  cada PR que toque un módulo con `integration=true` en `scripts/module-manifest.tsv` (`testing`,
  `bootstrap`, `database/postgres`, `database/mongodb`, `messaging/rabbit`), con Docker arrancado:
  `GOWORK=off make test-integration`. Recorre esos módulos **en serie** con
  `go test -tags integration -race -count=1 ./...` (sin `-short`), así que corre los ficheros
  `//go:build integration` y también los tests con Docker sin tag. El código de salida se lee sin pipe.
  - **La CI no corre la integración.** El job `docker-untagged` de `ci.yml` solo corre los tests con
    Docker que no llevan tag (los que se saltan con `-short`) de `testing` y `database/postgres`; lo que
    lleva `//go:build integration` solo lo ve este gate local.
  - No hay variable de entorno que active nada: `INTEGRATION_TESTS` no la leía ningún test y se retiró.
    Los dos interruptores reales son `-tags integration` y la ausencia de `-short`.
  - Con la máquina cargada, los containers de RabbitMQ vencen su timeout de arranque (60 s) y el gate
    falla sin que el código tenga culpa: se corre con la máquina libre y nunca en paralelo.
- **Lint y vet ven los ficheros de integración**: `.golangci.yml` fija `build-tags: [integration]` y no
  recorta hallazgos (`max-same-issues: 0`, `max-issues-per-linter: 0`). `make vet-all` y la CI corren
  `go vet` dos veces: sin tag y con `-tags integration`.
  - **No crees ficheros `//go:build !integration`**: con el tag fijo en la configuración dejarían de
    lintarse (`--build-tags` en la línea de comandos no lo quita). Hoy no hay ninguno. Un test sin Docker
    va sin build tag y corre en los dos modos.

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

## Versionado de módulos: migrar hacia adelante, nunca re-taguear (regla anti-conflicto cloud)

> **Por qué existe esta regla:** un "clean reset" del repo + re-taguear el **mismo** número de versión
> con contenido distinto rompió `audit`/`audit/postgres@v0.1.0`. En local **no se nota** (el `go.work`
> los resuelve del filesystem, sin checksum); en **CI/cloud sí**: Go descarga el tag de GitHub, lo hashea
> y lo compara contra el `go.sum` del consumidor → `SECURITY ERROR: checksum mismatch` → el deploy falla.
> El "historial" que muerde vive en el `go.sum` de cada consumidor (repos privados: **no hay sumdb global**;
> la fuente de verdad es el `go.sum` commiteado).

**Estándar:** los módulos migran al esquema **`0.900.X`** camino a producción (allí saltan a `1.0.0`). Un
módulo aún en `v0.1.0` está **"bajo estándar"**.

**Reglas (obligatorias):**
1. **Nunca re-pushees un tag existente.** Si el contenido de un módulo cambia, **sube el número de versión**
   (tag nuevo e inmutable). Re-taguear el mismo número es exactamente lo que rompió `audit`.
2. **Si necesitas modificar un módulo que aún está en `v0.1.0`, migrálo a `0.900.X` en el mismo cambio:**
   bump del tag **y** bump del `require` en **todos** sus consumidores + `go mod tidy` (regenera `go.sum`
   con el hash fresco). No lo dejes en `v0.1.0`. Recuerda que migrar un módulo que es dependencia de otro
   módulo `0.900.X` (p. ej. `audit` ← `middleware/gin`) obliga a **bumpear también a ese intermediario**
   y a sus consumidores (efecto cascada).
3. Tras migrar, **retira cualquier `replace` transitorio del `go.work`**: el objetivo es que CI/cloud
   resuelva contra el tag real, no contra el filesystem. El `replace` en `go.work` es muleta local, no fix.

**Un `v0.1.0` "sano" puede quedarse:** si nunca se re-taguea, es inmutable y seguro. La migración a
`0.900.X` es **incremental, al tocar** el módulo — no un barrido masivo del catálogo.
