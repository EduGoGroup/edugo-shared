# Improvement Assessment Index

Este directorio contiene el informe de mejora técnica de `edugo-shared`, con foco en dos preguntas:

1. Qué se puede mejorar en el código por módulo.
2. Qué implicaciones tendría cada cambio para los consumidores actuales.

## Documentos

- [01 - Executive report](01-executive-report.md)
- [02 - Module findings](02-module-findings.md)
- [03 - Consumer impact matrix](03-consumer-impact-matrix.md)
- [04 - Migration strategy](04-migration-strategy.md)
- [05 - Priority module deep dive](05-priority-module-bootstrap.md)

## Alcance y evidencia

Se usó evidencia del propio repositorio:

- `scripts/module-manifest.tsv` (inventario de módulos)
- `PLAN_REFACTORING.md` (hallazgos previos y plan de fases)
- `docs/phase-1/*` y `docs/phase-2/*` (catálogo, consumidores y flujos)
- estructura de carpetas y `go.mod` de los módulos

## Consumidores considerados

- `edugo-api-iam-platform`
- `edugo-api-admin-new`
- `edugo-api-mobile-new`
- `edugo-worker`
- consumo indirecto de `kmp_new` y `apple_new` vía APIs

## Convenciones de impacto

- **Sin impacto**: no requiere cambio en consumidor.
- **Compatible**: requiere ajuste menor sin romper API pública.
- **Breaking**: requiere cambios de código del consumidor y migración coordinada.
