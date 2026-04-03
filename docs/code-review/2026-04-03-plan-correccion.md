# Plan de Correccion Detallado - Code Review (2026-04-03)

## Objetivo

Eliminar brechas de validacion y restaurar consistencia operativa multi-modulo sin introducir cambios funcionales en APIs publicas.

---

## Fase 1 - Corregir bloqueantes de build/test (prioridad inmediata)

## 1.1 Normalizar `go.sum` en modulos con fallo

### Modulos
- `config`
- `testing`

### Acciones
- Ejecutar `go mod tidy` en cada modulo afectado.
- Verificar que se agreguen las entradas faltantes `.../go.mod`.
- Re-ejecutar en cada modulo:
  - `go build ./...`
  - `go test -short ./...`
  - `go vet ./...`

### Criterio de salida
- Ambos modulos quedan en verde para build/test/vet local.

---

## Fase 2 - Cerrar brecha manifest vs modulos reales (prioridad critica)

## 2.1 Actualizar `scripts/module-manifest.tsv`

### Acciones
- Incluir los 11 modulos faltantes:
  - `bootstrap/mongodb`
  - `bootstrap/postgres`
  - `bootstrap/rabbitmq`
  - `bootstrap/s3`
  - `health`
  - `lifecycle/shutdown`
  - `resilience/circuitbreaker`
  - `resilience/ratelimiter`
  - `resilience/retry`
  - `storage`
  - `storage/s3`
- Asignar nivel (`level`) segun dependencias reales.
- Definir flags `integration` y `coverage_validation` por modulo en base a su naturaleza de tests.

## 2.2 Validar impacto en orquestadores

### Acciones
- Ejecutar:
  - `./scripts/list-modules.sh --set all`
  - `make build-all`
  - `make test-all`
- Revisar que `.github/workflows/ci.yml` y `.github/workflows/test.yml` absorban el nuevo universo sin cambios adicionales estructurales.

### Criterio de salida
- `all` refleja el 100% de modulos reales (`go.mod`).
- Ningun modulo queda fuera de matriz por error de inventario.

---

## Fase 3 - Fortalecer experiencia local de lint

## 3.1 Robustecer `lint-all` o flujo de setup

### Opciones de correccion
- Opcion A: pre-check explicito en `Makefile` con mensaje claro si falta `golangci-lint`.
- Opcion B: documentar y estandarizar `make install-tools` como prerequisito obligatorio en README raiz (con validacion clara).

### Criterio de salida
- Contributor nuevo puede ejecutar lint sin ambiguedad de pasos.

---

## Fase 4 - Consistencia de version de Go y documentacion

## 4.1 Estandarizar directiva `go` en modulos

### Acciones
- Definir version objetivo (por ejemplo `1.25.0` o `1.25.5`).
- Homogeneizar la directiva `go` en todos los `go.mod`, evitando mezcla `1.25` / `1.25.0` / `1.25.5`.

## 4.2 Actualizar documentacion transversal

### Archivos minimos
- `README.md` (estado operativo)
- `docs/phase-3/README.md` (conteo/alcance de modulos)
- `docs/phase-3/validation-contract.md` (si cambian sets o criterios)

### Criterio de salida
- Documentacion alineada con la realidad del repositorio.

---

## Priorizacion recomendada

1. **Fase 1** (bloqueante tecnico inmediato)
2. **Fase 2** (riesgo critico de cobertura CI)
3. **Fase 3** (productividad y onboarding)
4. **Fase 4** (consistencia y gobernanza)

---

## Riesgos y mitigaciones

- **Riesgo**: incorporar modulos faltantes puede exponer fallos ocultos en CI.
  - **Mitigacion**: habilitar por lotes pequeños y validar cada lote con build/test/vet.
- **Riesgo**: cambios de `go.mod/go.sum` masivos.
  - **Mitigacion**: limitar cambios a modulos afectados y validar por modulo.
- **Riesgo**: thresholds de cobertura no listos para nuevos modulos.
  - **Mitigacion**: iniciar con umbrales conservadores y plan incremental.

---

## Checklist de ejecucion sugerido

- [ ] Resolver `go.sum` en `config` y `testing`
- [ ] Confirmar build/test/vet verdes en ambos
- [ ] Agregar 11 modulos faltantes al manifest
- [ ] Ajustar nivel/integration/coverage_validation por modulo
- [ ] Validar `list-modules` + `make build-all` + `make test-all`
- [ ] Definir y aplicar estrategia para `golangci-lint` local
- [ ] Homogeneizar directiva `go` en todos los modulos
- [ ] Actualizar README/docs de fase 3
