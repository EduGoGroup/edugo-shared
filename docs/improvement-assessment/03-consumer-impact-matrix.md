# Consumer Impact Matrix

Esta matriz conecta mejoras propuestas con impacto esperado en consumidores verificados.

## Consumidores base

- `edugo-api-iam-platform`
- `edugo-api-admin-new`
- `edugo-api-mobile-new`
- `edugo-worker`
- `kmp_new` y `apple_new` (impacto indirecto por comportamiento de APIs)

## Impacto por tipo de mejora

| Mejora | IAM | Admin | Mobile | Worker | Frontends | Tipo de impacto |
| --- | --- | --- | --- | --- | --- | --- |
| Limpieza de tests triviales en `edugo-shared` | Nulo | Nulo | Nulo | Nulo | Nulo | Sin impacto |
| Split unit/integration en módulos de infraestructura | Nulo | Nulo | Nulo | Nulo | Nulo | Sin impacto |
| Reestructurar `bootstrap` hacia `internal/` manteniendo API pública | Nulo | Nulo | Nulo | Bajo (validar wiring) | Nulo | Compatible |
| Reestructurar `logger` con wrappers públicos estables | Bajo | Bajo | Bajo | Bajo | Nulo | Compatible |
| Reestructurar `screenconfig` sin cambiar contratos públicos | Nulo (hoy sin consumo directo) | Nulo | Nulo | Nulo | Nulo | Sin impacto |
| Extraer `common/health` y adopción gradual | Nulo | Nulo | Nulo | Bajo | Nulo | Compatible |
| Extraer `common/retry` y uso en `messaging/rabbit` | Nulo | Nulo | Bajo | Nulo | Nulo | Compatible |
| Eliminar funciones públicas `All*()` marcadas como dead code | Bajo* | Bajo* | Bajo* | Bajo* | Nulo | Potencial breaking |
| Eliminar wrappers de config en `bootstrap` (si cambia contrato) | Nulo | Nulo | Nulo | Alto | Nulo | Breaking (Worker) |

\* Bajo si realmente no hay uso externo; alto si existe uso no inventariado.

## Impacto por módulo consumido

| Módulo | Consumidores directos | Sensibilidad al cambio |
| --- | --- | --- |
| `auth` | IAM, Admin, Mobile | Alta |
| `middleware/gin` | IAM, Admin, Mobile | Alta |
| `repository` | IAM, Admin, Mobile | Alta |
| `audit` / `audit/postgres` | IAM, Admin, Mobile | Media |
| `logger` | IAM, Admin, Mobile, Worker | Alta |
| `common` | IAM, Admin, Mobile, Worker | Alta |
| `bootstrap` | Worker | Alta (concentrada) |
| `database/postgres` | Worker | Media |
| `lifecycle` | Worker | Media |
| `testing` | Worker (suite de integración) | Media |
| `cache/redis` | Mobile | Media |
| `messaging/events` | Mobile | Alta (compatibilidad de schema) |
| `messaging/rabbit` | Mobile | Media |
| `config`, `database/mongodb`, `screenconfig` | Sin consumo directo verificado | Baja |

## Reglas de compatibilidad recomendadas

1. Cambios de estructura interna solo con API pública intacta.
2. Cambios a contratos públicos en `auth`, `middleware/gin`, `repository`, `logger`, `common` requieren migración documentada.
3. Cambios de schema en `messaging/events` deben mantener backward compatibility por versión.
4. Cualquier posible breaking en módulos de Worker debe validarse primero en `edugo-worker` antes de release.
