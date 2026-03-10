# Matriz servicio-modulo

Esta matriz resume que modulos de `edugo-shared` consume cada servicio backend del ecosistema. La evidencia se tomo del `go.mod` de cada servicio y luego se contrasto con sus contenedores de dependencias.

## Matriz compacta

| Servicio | Modulos de `edugo-shared` consumidos directamente |
| --- | --- |
| `edugo-api-iam-platform` | `audit`, `audit/postgres`, `auth`, `common`, `logger`, `middleware/gin`, `repository` |
| `edugo-api-admin-new` | `audit`, `audit/postgres`, `auth`, `common`, `logger`, `middleware/gin`, `repository` |
| `edugo-api-mobile-new` | `audit`, `audit/postgres`, `auth`, `cache/redis`, `common`, `logger`, `messaging/events`, `messaging/rabbit`, `middleware/gin`, `repository` |
| `edugo-worker` | `bootstrap`, `common`, `database/postgres`, `lifecycle`, `logger`, `testing` |

## Lectura por servicio

### IAM Platform

- Usa [`auth`](../../auth/docs/README.md) para `JWTManager` y manejo de tokens.
- Usa [`middleware/gin`](../../middleware/gin/docs/README.md) para validar JWT, auditar mutaciones y verificar permisos.
- Usa [`audit/postgres`](../../audit/postgres/docs/README.md) para persistir auditoria en `audit.events`.
- Usa [`repository`](../../repository/docs/README.md) para usuarios, memberships y schools compartidos.
- Usa [`logger`](../../logger/docs/README.md) y [`common`](../../common/docs/README.md) como capas transversales.

### Admin API

- Reutiliza casi el mismo stack compartido que IAM: auth, auditoria, middleware, repository, logger y common.
- La diferencia funcional es que la autenticacion se resuelve con un `AuthClient` propio que puede validar localmente y hacer fallback remoto a IAM.
- El middleware compartido de permisos sigue siendo el mismo [`middleware/gin`](../../middleware/gin/docs/README.md).

### Mobile API

- Reutiliza auth, auditoria, middleware, logger, common y repository.
- Agrega [`cache/redis`](../../cache/redis/docs/README.md) para caching de screen/service data.
- Agrega [`messaging/events`](../../messaging/events/docs/README.md) y [`messaging/rabbit`](../../messaging/rabbit/docs/README.md) para publicar eventos asincronos.
- Sigue usando auditoria persistida en PostgreSQL mediante [`audit/postgres`](../../audit/postgres/docs/README.md).

### Worker

- No usa middleware HTTP ni auth shared del mismo modo que las APIs.
- Usa [`bootstrap`](../../bootstrap/docs/README.md) para construir recursos de infraestructura.
- Usa [`lifecycle`](../../lifecycle/docs/README.md) para ordenar startup/cleanup.
- Usa [`database/postgres`](../../database/postgres/docs/README.md) y [`common`](../../common/docs/README.md) como base técnica.
- Usa [`testing`](../../testing/docs/README.md) para su integracion con containers.

## Modulos que no aparecieron como dependencia directa en los servicios escaneados

- `config`
- `database/mongodb`
- `screenconfig`

Eso no significa que esten obsoletos. Significa que, en el estado actual del ecosistema escaneado, no aparecen como imports directos en IAM, Admin, Mobile o Worker. Pueden estar disponibles para evolucion futura o ser absorbidos por implementaciones locales de cada servicio.
