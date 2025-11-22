# Resumen Ejecutivo: Plan de Mejora de Cobertura
## edugo-shared - Go 1.25

---

## 📊 Estado Actual vs Objetivo

```
Módulos que NECESITAN Atención (<80%):

messaging/rabbit:     [████████░░░░░░░░░░] 30% → 80% (GAP: +50%)  🔴 CRÍTICA
bootstrap:            [████████████░░░░░░] 40% → 80% (GAP: +40%)  🔴 ALTA
database/postgres:    [█████████████████░] 59% → 80% (GAP: +21%)  🟡 MEDIA
testing:              [█████████████████░] 59% → 80% (GAP: +21%)  🟡 MEDIA

Módulos CERCA del Umbral (80-85%):

database/mongodb:     [████████████████▓░] 82% → 85% ✅ MANTENER
config:               [████████████████▓░] 83% → 85% ✅ MANTENER
auth:                 [█████████████████░] 85% → 90% ✅ BUENO

Módulos EXCELENTES (>90%):

evaluation:           [████████████████████] 100%    ✅ PERFECTO
middleware/gin:       [███████████████████▓] 98.5%   ✅ EXCELENTE
logger:               [███████████████████░] 95.8%   ✅ EXCELENTE
lifecycle:            [██████████████████▓░] 91.8%   ✅ EXCELENTE
```

---

## 🎯 Plan de Acción Rápido

### Fase 1: CRÍTICA (Sprints 5-6, 3 semanas)
**Objetivo:** Llevar módulos críticos de <50% a 80%

| Semana | Módulo | Tests | Horas | Resultado |
|--------|--------|-------|-------|-----------|
| 1-2 | messaging/rabbit | 34 | 10-14h | 30% → 75-80% |
| 2-3 | bootstrap | 41 | 11-15h | 40% → 85-90% |

### Fase 2: MEDIA (Sprint 6, 1 semana)
**Objetivo:** Completar módulos medios a 80%

| Semana | Módulo | Tests | Horas | Resultado |
|--------|--------|-------|-------|-----------|
| 3 | database/postgres | 12 | 3-4h | 59% → 80% |
| 3 | testing | 18 | 4h | 59% → 80% |

---

## ⏱️ Esfuerzo Total

```
TOTAL: 105 tests nuevos
ESFUERZO: 28-37 horas (3-4 semanas)
RESULTADO: 4 módulos de <80% → 80%+
```

---

## 📋 Checklist por Sprint

### Sprint 5 - Semana 1: messaging/rabbit

**Archivos a crear:**
```bash
messaging/rabbit/
├── consumer_dlq_test.go           # NUEVO - 10 tests
├── dlq_integration_test.go        # NUEVO - 8 tests
├── connection_advanced_test.go    # NUEVO - 6 tests
├── publisher_advanced_test.go     # NUEVO - 5 tests
└── consumer_advanced_test.go      # NUEVO - 5 tests
```

**Comandos:**
```bash
cd messaging/rabbit
go test -v -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | tail -1
# Esperado: 75-80%
```

- [ ] Día 1-2: DLQ Consumer Integration (10 tests, 3-4h)
- [ ] Día 3: DLQ Helpers (8 tests, 2-3h)
- [ ] Día 4: Connection Advanced (6 tests, 2-3h)
- [ ] Día 5: Publisher/Consumer Advanced (10 tests, 3-4h)
- [ ] Verificar cobertura ≥ 75%
- [ ] Commit y push

---

### Sprint 5 - Semana 2: bootstrap

**Archivos a crear:**
```bash
bootstrap/
├── factory_postgresql_integration_test.go  # NUEVO - 8 tests
├── factory_mongodb_integration_test.go     # NUEVO - 8 tests
├── factory_rabbitmq_integration_test.go    # NUEVO - 8 tests
├── factory_s3_test.go                      # NUEVO - 6 tests
├── bootstrap_integration_test.go           # NUEVO - 6 tests
├── options_test.go                         # AMPLIAR - 3 tests
└── resource_implementations_test.go        # NUEVO - 2 tests
```

**Comandos:**
```bash
cd bootstrap
go test -v -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | tail -1
# Esperado: 85-90%
```

- [ ] Día 1: PostgreSQL Factory (8 tests, 2-3h)
- [ ] Día 2: MongoDB Factory (8 tests, 2-3h)
- [ ] Día 3: RabbitMQ Factory (8 tests, 2-3h)
- [ ] Día 4: S3 Factory (6 tests, 2h)
- [ ] Día 5: Bootstrap Integration + Options (11 tests, 3-4h)
- [ ] Verificar cobertura ≥ 85%
- [ ] Commit y push

---

### Sprint 6 - Semana 3: database/postgres + testing

**Archivos a modificar/crear:**
```bash
database/postgres/
├── transaction_test.go           # AMPLIAR - 8 tests nuevos
└── connection_test.go            # AMPLIAR - 4 tests nuevos

testing/containers/
├── helpers_test.go               # NUEVO - 8 tests
├── options_test.go               # NUEVO - 6 tests
└── manager_test.go               # AMPLIAR - 4 tests nuevos
```

**Comandos:**
```bash
# database/postgres
cd database/postgres
go test -v -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | tail -1
# Esperado: 80%

# testing
cd ../../testing
go test -v -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | tail -1
# Esperado: 80%
```

- [ ] Día 1-2: database/postgres (12 tests, 3-4h)
- [ ] Día 3-4: testing (18 tests, 4h)
- [ ] Día 5: Verificación final, ajustes
- [ ] Ejecutar `make coverage-all-modules`
- [ ] Actualizar `.coverage-thresholds.yml`
- [ ] Commit y push

---

## 🔍 Tests Críticos por Módulo

### messaging/rabbit (34 tests)
```
✓ DLQ Consumer Integration        10 tests  +15% cobertura
✓ DLQ Helpers & Utilities           8 tests  +10% cobertura
✓ Connection Advanced               6 tests  +8% cobertura
✓ Publisher/Consumer Advanced      10 tests  +12% cobertura
                                   --------  -------------
                              TOTAL: 34 tests  +45% cobertura
```

### bootstrap (41 tests)
```
✓ PostgreSQL Factory Integration    8 tests  +10% cobertura
✓ MongoDB Factory Integration        8 tests  +10% cobertura
✓ RabbitMQ Factory Integration       8 tests  +10% cobertura
✓ S3 Factory (mocked)                6 tests  +8% cobertura
✓ Bootstrap Integration              6 tests  +7% cobertura
✓ Options & Resources                5 tests  +5% cobertura
                                    --------  -------------
                               TOTAL: 41 tests  +50% cobertura
```

### database/postgres (12 tests)
```
✓ Transaction Advanced               8 tests  +15% cobertura
✓ Connection Edge Cases              4 tests  +6% cobertura
                                    --------  -------------
                               TOTAL: 12 tests  +21% cobertura
```

### testing (18 tests)
```
✓ Helpers (ExecSQL, WaitFor, Retry)  8 tests  +10% cobertura
✓ Options Builder Pattern            6 tests  +7% cobertura
✓ Manager Edge Cases                 4 tests  +4% cobertura
                                    --------  -------------
                               TOTAL: 18 tests  +21% cobertura
```

---

## 🛠️ Comandos Útiles

### Verificar Cobertura Individual
```bash
# En cada módulo
cd <módulo>
go test -v -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | tail -1
go tool cover -html=coverage.out -o coverage.html
```

### Verificar Cobertura Global
```bash
# En raíz del proyecto
make coverage-all-modules
make validate-coverage
```

### Tests con Race Detection
```bash
go test -v -race ./...
```

### CI Local
```bash
./scripts/test-ci-local.sh
```

---

## ✅ Criterios de Aceptación

**Por Módulo:**
- [ ] Cobertura ≥ 80%
- [ ] Todos los tests pasan
- [ ] No race conditions
- [ ] Tests ejecutan en < 2 minutos
- [ ] Coverage HTML generado

**Global:**
- [ ] 4 módulos alcanzaron 80%
- [ ] CI pipeline pasa
- [ ] Coverage validation exitosa
- [ ] `.coverage-thresholds.yml` actualizado
- [ ] Documentación actualizada

---

## 📈 Métricas de Progreso

### Seguimiento Semanal

**Semana 1:**
```
messaging/rabbit:  30% → [____] %  (Objetivo: 75-80%)
```

**Semana 2:**
```
bootstrap:         40% → [____] %  (Objetivo: 85-90%)
```

**Semana 3:**
```
database/postgres: 59% → [____] %  (Objetivo: 80%)
testing:           59% → [____] %  (Objetivo: 80%)
```

### Tracking de Tests
```
Semana 1:  [ ] 34/34 tests (messaging/rabbit)
Semana 2:  [ ] 41/41 tests (bootstrap)
Semana 3:  [ ] 30/30 tests (postgres + testing)
           ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
TOTAL:     [ ] 105/105 tests completados
```

---

## 🚨 Riesgos y Mitigación

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| Tests lentos (>5min) | Media | Alto | Usar testcontainers singleton, parallel tests |
| Race conditions | Media | Alto | Ejecutar con `-race` regularmente |
| Containers no inician | Baja | Alto | Verificar Docker antes de empezar |
| Scope creep | Media | Medio | Seguir plan estrictamente, no agregar features |

---

## 📞 Soporte

**Documentación:**
- Plan Detallado: `docs/coverage-improvement-plan.md`
- Thresholds Config: `.coverage-thresholds.yml`
- CI Workflows: `.github/workflows/coverage-validation.yml`

**Herramientas:**
- `make help` - Ver todos los comandos
- `make coverage-status` - Ver estado actual
- `./scripts/analyze-coverage.sh` - Análisis detallado

---

*Última actualización: 2025-11-22*
*Siguiente revisión: Fin de Sprint 5*
