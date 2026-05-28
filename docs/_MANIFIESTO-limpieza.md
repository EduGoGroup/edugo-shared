# Manifiesto de limpieza — edugo-shared/docs (2026-05-22)

Triage de la documentación interna de `edugo-shared`. Sesgo conservador.
No se tocó código (.go), go.mod, scripts ni .git.

## Veredictos

| Archivo | Veredicto | Razón | Destino |
| --- | --- | --- | --- |
| `README.md` (raíz) | KEEP | Índice modular vigente, refleja estructura actual | Se mantiene |
| `CHANGELOG.md` (raíz) | KEEP | Changelog activo (entradas Unreleased OTel recientes) | Se mantiene |
| `PLAN_REFACTORING.md` (raíz) | RESCUE | Plan ~75-80% implementado (Fase 0/1/3 hechas, Fase 2 pendiente) | Resumen → `EduGo/docs/planes-rescatados/shared-refactoring.md`; original borrado |
| `docs/README.md` | KEEP | Índice de docs; corregido para quitar links a code-review borrado | Editado |
| `docs/phase-1/*` | KEEP | Describe arquitectura/módulos reales del repo | Se mantiene |
| `docs/phase-2/*` | KEEP | Matriz servicio-módulo y flujos de integración (consumidores) | Se mantiene |
| `docs/phase-3/*` | KEEP | Contrato de validación y workflow de release modular vigente | Se mantiene |
| `docs/roadmap/*` | KEEP | Estado de fases documentales y roadmap | Se mantiene |
| `docs/shared-modules-roadmap.md` | KEEP | Inventario de extracción de módulos reutilizables; la mayoría "Done" verificado en código, items "Evaluate" con valor futuro | Se mantiene |
| `docs/improvement-assessment/` (6 archivos) | DELETE | Duplica PLAN_REFACTORING con datos stale (describe bootstrap con 26 archivos/init_*.go que ya no existen). Valor vivo (Fase 2) absorbido en plan rescatado | Borrado |
| `docs/code-review/2026-04-03-hallazgos.md` | RESCUE | Hallazgos de manifest CI / versión Go / doc drift aún válidos al 2026-05-22 | Resumen → plan rescatado; original borrado |
| `docs/code-review/2026-04-03-plan-correccion.md` | RESCUE | Plan de corrección de esos hallazgos, mayormente pendiente | Resumen → plan rescatado; original borrado |
| `docs/.DS_Store` | DELETE | Basura del sistema | Borrado |

## Notas de verificación

- Fase 0/1/3 de PLAN_REFACTORING confirmadas en código: dead code `AllX()` eliminado,
  bootstrap consolidado, build tags integration, `common/health` y `common/retry` existen.
- Fase 2 (internal/) NO implementada salvo en `audit/postgres` → único pendiente vivo.
- Code-review: manifest aún lista ~19 módulos vs 31 `go.mod` reales; mezcla `go 1.25`/`1.25.0`;
  README/phase-3 aún dicen "17 modulos". Hallazgos vigentes → rescatados.

## docs/ resultante (vigente)

README.md, phase-1/, phase-2/, phase-3/, roadmap/, shared-modules-roadmap.md, cicd/ (placeholder).
