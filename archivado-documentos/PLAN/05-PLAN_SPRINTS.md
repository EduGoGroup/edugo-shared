# Plan de Sprints para edugo-shared

## 🎯 Objetivo

Ejecutar el plan de implementación en **sprints iterativos** que culminarán en la versión congelada **v0.7.0** de todos los módulos de edugo-shared.

**Duración total estimada:** 2-3 semanas  
**Fecha de inicio:** A definir  
**Fecha objetivo de congelamiento:** +3 semanas desde inicio

---

## 📅 Sprint 0: Auditoría y Alineación (2-3 horas)

### Objetivo
Preparar el terreno, arreglar dependencias rotas, y tener baseline limpio para comenzar desarrollo.

### Tareas

#### 1. Verificar sincronización de ramas
- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`
- [ ] `git checkout main && git pull origin main`
- [ ] `git checkout dev && git pull origin dev`
- [ ] `git diff main dev --stat`
- [ ] Si hay diferencias funcionales: mergear dev a main o viceversa
- [ ] **Decisión:** Trabajar desde `dev` (es la rama más actualizada)

**Tiempo:** 15 minutos

---

#### 2. Fix dependencias rotas
- [ ] `cd auth && go mod tidy && go test ./...`
- [ ] `cd ../middleware/gin && go mod tidy && go test ./...`
- [ ] Verificar que tests pasan (o al menos se ejecutan)
- [ ] Commit: `fix(deps): execute go mod tidy on auth and middleware/gin`

**Tiempo:** 30 minutos

---

#### 3. Ejecutar suite completa de tests
- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`
- [ ] `make test-all-modules > test-results-baseline.txt 2>&1`
- [ ] Documentar en `01-ESTADO_ACTUAL.md`:
  - Módulos con tests passing
  - Módulos con tests failing
  - Coverage actual por módulo

**Tiempo:** 1 hora

---

#### 4. Crear issues en GitHub
- [ ] Issue #1: "Create evaluation/ module (P0)"
- [ ] Issue #2: "Implement DLQ support in messaging/rabbit (P0)"
- [ ] Issue #3: "Increase database/postgres coverage to >80% (P0)"
- [ ] Issue #4: "Add tests to logger, common/* modules (P1)"
- [ ] Issue #5: "Implement/verify refresh tokens in auth (P1)"
- [ ] Issue #6: "Increase coverage in config, bootstrap (P2)"

**Tiempo:** 30 minutos

---

### Entregables Sprint 0
- [x] Ramas sincronizadas
- [x] Dependencias arregladas (go mod tidy)
- [x] Baseline de tests documentado
- [x] Issues creados en GitHub

**Total tiempo:** 2-3 horas

---

## 🚀 Sprint 1: Módulos Críticos Nuevos (1 semana)

### Objetivo
Crear módulo `evaluation/` y agregar features críticas (DLQ).

**Duración:** 5 días laborables  
**Prioridad:** P0 (bloquean desarrollo de consumidores)

---

### Día 1-2: Crear módulo evaluation/

#### Tareas
- [ ] Crear carpeta `evaluation/`
- [ ] Crear `go.mod`:
```bash
cd evaluation
go mod init github.com/EduGoGroup/edugo-shared/evaluation
```

- [ ] Implementar `assessment.go`:
  - Struct `Assessment`
  - Métodos: `Validate()`, `IsPublished()`
  - Tests: `assessment_test.go`

- [ ] Implementar `question.go`:
  - Struct `Question`, `QuestionOption`
  - Enum `QuestionType`
  - Métodos: `Validate()`, `GetCorrectOptions()`
  - Tests: `question_test.go`

- [ ] Implementar `attempt.go`:
  - Struct `Attempt`, `Answer`
  - Métodos: `CalculatePercentage()`, `CheckPassed()`
  - Tests: `attempt_test.go`

- [ ] Crear `README.md` con ejemplos de uso

- [ ] Ejecutar tests:
```bash
go test -v -cover ./...
# Target: >80% coverage
```

- [ ] Commit y push:
```bash
git add evaluation/
git commit -m "feat(evaluation): create evaluation module with Assessment, Question, Attempt models"
git push origin dev
```

- [ ] Crear tag:
```bash
git tag evaluation/v0.1.0
git push origin evaluation/v0.1.0
```

**Tiempo:** 2 días (4-5 horas implementación + tests)

---

### Día 3-4: Implementar DLQ en messaging/rabbit/

#### Tareas
- [ ] Crear `dlq.go`:
  - Struct `DLQConfig`
  - Método `DefaultDLQConfig()`
  - Método `calculateBackoff()`

- [ ] Modificar `consumer.go`:
  - Agregar `DLQ` field a `ConsumerConfig`
  - Crear método `ConsumeWithDLQ()`
  - Implementar `setupDLQ()`
  - Implementar `sendToDLQ()`
  - Helper `getRetryCount()`

- [ ] Crear tests:
  - `dlq_test.go` (tests unitarios de backoff)
  - `consumer_dlq_test.go` (integración con Testcontainers)

- [ ] Ejecutar tests:
```bash
cd messaging/rabbit
go test -v -cover ./...
```

- [ ] Commit y tag:
```bash
git commit -m "feat(messaging/rabbit): add Dead Letter Queue (DLQ) support with retry logic"
git tag messaging/rabbit/v0.6.0
git push origin dev messaging/rabbit/v0.6.0
```

**Tiempo:** 2 días (3-5 horas)

---

### Día 5: Aumentar coverage en database/postgres/

#### Tareas
- [ ] Crear `postgres_integration_test.go`:
  - Setup con Testcontainers
  - Test de conexión
  - Test de transacciones
  - Test de health check
  - Test de reconnection

- [ ] Ejecutar coverage:
```bash
cd database/postgres
go test -v -cover ./...
# Target: >80% (actualmente 2%)
```

- [ ] Si coverage <80%: agregar más tests

- [ ] Commit y tag:
```bash
git commit -m "test(database/postgres): increase coverage from 2% to >80% with integration tests"
git tag database/postgres/v0.6.0
git push origin dev database/postgres/v0.6.0
```

**Tiempo:** 1 día (4-6 horas)

---

### Entregables Sprint 1
- [x] evaluation/v0.1.0 publicado
- [x] messaging/rabbit/v0.6.0 con DLQ publicado
- [x] database/postgres/v0.6.0 con >80% coverage
- [x] Todos los tests P0 pasando

**Validación:**
```bash
# Verificar que api-mobile puede importar evaluation
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
go get github.com/EduGoGroup/edugo-shared/evaluation@v0.1.0
# Debe compilar sin errores
```

---

## 🧪 Sprint 2: Features Faltantes (1 semana)

### Objetivo
Agregar tests a módulos sin tests, implementar refresh tokens, validar MongoDB.

**Duración:** 5 días laborables  
**Prioridad:** P1 (afectan calidad)

---

### Día 1-2: Agregar tests a logger/ y common/*

#### Tareas - logger/
- [ ] Crear `logger_test.go`:
  - Test de creación de logger
  - Test de niveles (Debug, Info, Warn, Error)
  - Test de formatos (JSON, Console)
  - Test de context fields

- [ ] Coverage target: >80%

- [ ] Commit y tag:
```bash
git commit -m "test(logger): add comprehensive unit tests (>80% coverage)"
git tag logger/v0.6.0
git push origin dev logger/v0.6.0
```

**Tiempo:** 1 día (3-4 horas)

---

#### Tareas - common/*
- [ ] `common/errors/errors_test.go`:
  - Test de cada tipo de error
  - Test de HTTP status codes
  - Test de error messages

- [ ] `common/types/uuid_test.go`:
  - Test de JSON marshaling/unmarshaling
  - Test de validación

- [ ] `common/types/enum/enum_test.go`:
  - Test de todos los enums
  - Test de string conversion

- [ ] `common/validator/validator_test.go`:
  - Test de validación de email
  - Test de validación de UUID
  - Test de campos requeridos

- [ ] Coverage target: >80% por submódulo

- [ ] Commit y tag:
```bash
git commit -m "test(common): add tests to errors, types, validator (>80% coverage)"
git tag common/v0.6.0
git push origin dev common/v0.6.0
```

**Tiempo:** 1 día (6-8 horas)

---

### Día 3: Implementar/Verificar Refresh Tokens en auth/

#### Opción A: Si NO existe

- [ ] Implementar `GenerateTokenPair()`
- [ ] Implementar `RefreshAccessToken()`
- [ ] Implementar `ValidateRefreshToken()`
- [ ] Struct `TokenPair`, `RefreshClaims`
- [ ] Tests: `refresh_token_test.go`

#### Opción B: Si YA existe

- [ ] Ejecutar tests: `go test -v ./...`
- [ ] Verificar coverage >80%
- [ ] Documentar en README.md

---

- [ ] Commit y tag:
```bash
git commit -m "feat(auth): implement refresh token support"
# o "test(auth): verify and document refresh token feature"
git tag auth/v0.6.0
git push origin dev auth/v0.6.0
```

**Tiempo:** 1 día (2-3 horas si no existe, 1 hora si existe)

---

### Día 4: Validar tests en database/mongodb/

#### Tareas
- [ ] Crear `mongodb_integration_test.go` (si no existe):
  - Setup con Testcontainers
  - Test de conexión
  - Test de InsertOne, FindOne
  - Test de UpdateOne
  - Test de health check

- [ ] Ejecutar tests:
```bash
cd database/mongodb
go test -v -cover ./...
```

- [ ] Coverage target: >80%

- [ ] Commit y tag:
```bash
git commit -m "test(database/mongodb): add integration tests with Testcontainers"
git tag database/mongodb/v0.6.0
git push origin dev database/mongodb/v0.6.0
```

**Tiempo:** 1 día (2-3 horas)

---

### Día 5: Buffer / Refactoring

#### Tareas
- [ ] Revisar todos los tests agregados en Sprint 2
- [ ] Ejecutar `make test-all-modules`
- [ ] Arreglar tests failing
- [ ] Refactorizar código si es necesario
- [ ] Actualizar documentación (README.md de cada módulo)

**Tiempo:** 1 día

---

### Entregables Sprint 2
- [x] logger/v0.6.0 con tests
- [x] common/v0.6.0 con tests
- [x] auth/v0.6.0 con refresh tokens
- [x] database/mongodb/v0.6.0 validado
- [x] Coverage >80% en todos los módulos P1

---

## 🎯 Sprint 3: Consolidación y Congelamiento (3 días)

### Objetivo
Aumentar coverage en módulos P2, validar integración completa, release coordinado a v0.7.0.

**Duración:** 3 días  
**Prioridad:** P2 + congelamiento

---

### Día 1: Coverage P2 (config, bootstrap)

#### Tareas - config/
- [ ] Aumentar coverage de 32.9% a >80%
- [ ] Agregar tests de multi-environment
- [ ] Tests de Viper integration

#### Tareas - bootstrap/
- [ ] Aumentar coverage de 29.9% a >80%
- [ ] Tests de inicialización
- [ ] Tests de dependency injection

- [ ] Commit:
```bash
git commit -m "test(config,bootstrap): increase coverage to >80%"
git tag config/v0.6.0 bootstrap/v0.6.0
git push origin dev config/v0.6.0 bootstrap/v0.6.0
```

**Tiempo:** 1 día (4-6 horas)

---

### Día 2: Validación Completa

#### Tareas
- [ ] Ejecutar suite completa de tests:
```bash
make test-all-modules | tee test-results-final.txt
```

- [ ] Verificar que TODOS los tests pasan (0 failing)

- [ ] Calcular coverage global:
```bash
make coverage-all-modules | tee coverage-report.txt
```

- [ ] **Validación:** Coverage global >85%

- [ ] Validar que proyectos consumidores compilan:
```bash
# api-mobile
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
# Actualizar go.mod con últimas versiones
go get github.com/EduGoGroup/edugo-shared/evaluation@v0.1.0
go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@v0.6.0
go build ./cmd/api-mobile
# Debe compilar sin errores

# api-admin
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion
go build ./cmd/api-admin

# worker
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-worker
go get github.com/EduGoGroup/edugo-shared/evaluation@v0.1.0
go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@v0.6.0
go build ./cmd/worker
```

- [ ] Si hay errores de compilación: ARREGLAR antes de congelar

**Tiempo:** 1 día

---

### Día 3: Release Coordinado v0.7.0

#### Tareas

**1. Crear rama de release:**
```bash
git checkout dev
git checkout -b release/v0.7.0
```

**2. Actualizar CHANGELOG.md:**
```markdown
# Changelog

## [0.7.0] - 2025-11-XX - FROZEN RELEASE

### 🎉 Version Congelada
Esta versión es la BASE CONGELADA para el ecosistema EduGo MVP.
NO se agregarán features nuevas hasta post-MVP.

### Added
- **NEW MODULE**: evaluation/ v0.1.0 (Assessment, Question, Attempt models)
- messaging/rabbit: Dead Letter Queue (DLQ) support
- auth: Refresh token support
- Comprehensive tests across all modules (>85% global coverage)

### Changed
- ALL modules bumped to v0.7.0 (coordinated release)
- database/postgres: Coverage increased from 2% to >80%
- logger: Coverage increased from 0% to >80%
- common/*: Coverage increased from 0% to >80%

### Fixed
- auth, middleware/gin: Fixed broken dependencies (go mod tidy)
```

**3. Mergear a main:**
```bash
git add -A
git commit -m "chore: release v0.7.0 - frozen version for EduGo MVP"
git checkout main
git merge release/v0.7.0
git push origin main
```

**4. Crear tags coordinados:**
```bash
# Todos los módulos a v0.7.0
git tag auth/v0.7.0
git tag logger/v0.7.0
git tag common/v0.7.0
git tag config/v0.7.0
git tag bootstrap/v0.7.0
git tag lifecycle/v0.7.0
git tag middleware/gin/v0.7.0
git tag messaging/rabbit/v0.7.0
git tag database/postgres/v0.7.0
git tag database/mongodb/v0.7.0
git tag testing/v0.7.0  # Bump desde v0.6.2
git tag evaluation/v0.7.0  # Bump desde v0.1.0

# Push todos los tags
git push origin --tags
```

**5. Mergear main a dev:**
```bash
git checkout dev
git merge main
git push origin dev
```

**6. Crear GitHub Release:**
- [ ] Ir a https://github.com/EduGoGroup/edugo-shared/releases/new
- [ ] Tag: `v0.7.0`
- [ ] Title: "v0.7.0 - Frozen Release for EduGo MVP"
- [ ] Description: Copiar CHANGELOG.md
- [ ] Marcar como "Latest release"
- [ ] Publish

**7. Actualizar README.md:**
```markdown
## Installation

### Recommended: Use v0.7.0 (Frozen Release)

bash
go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0
go get github.com/EduGoGroup/edugo-shared/evaluation@v0.7.0
# ... otros módulos


### ⚠️ Important: v0.7.0 is FROZEN
- No new features will be added until post-MVP
- Only critical bug fixes (v0.7.1, v0.7.2...)
- Breaking changes NOT allowed
```

**Tiempo:** 1 día (2-3 horas)

---

### Entregables Sprint 3
- [x] Todos los módulos en v0.7.0
- [x] Coverage global >85%
- [x] Tests 100% passing
- [x] api-mobile, api-admin, worker compilan exitosamente
- [x] GitHub Release publicado
- [x] **SHARED CONGELADO**

---

## 📊 Resumen de Sprints

| Sprint | Duración | Entregables | Prioridad |
|--------|----------|-------------|-----------|
| Sprint 0 | 2-3 horas | Baseline limpio, issues creados | Preparación |
| Sprint 1 | 1 semana | evaluation/, DLQ, postgres tests | P0 |
| Sprint 2 | 1 semana | Tests en logger/common/auth/mongodb | P1 |
| Sprint 3 | 3 días | Coverage P2, release v0.7.0 | P2 + congelamiento |

**Total tiempo:** 2-3 semanas (10-15 días laborables)

---

## ✅ Criterios de Éxito del Plan Completo

### Para considerar el plan EXITOSO:

- ✅ Todos los sprints completados
- ✅ Módulo `evaluation/` existe y funciona
- ✅ DLQ implementado en `messaging/rabbit/`
- ✅ Refresh tokens implementados en `auth/`
- ✅ Coverage global >85%
- ✅ 0 tests failing
- ✅ 0 dependencias rotas
- ✅ api-mobile compila con shared v0.7.0
- ✅ api-admin compila con shared v0.7.0
- ✅ worker compila con shared v0.7.0
- ✅ Todos los módulos en v0.7.0
- ✅ GitHub Release publicado
- ✅ **shared CONGELADO**

---

## 🚨 Riesgos y Mitigaciones

### Riesgo 1: Tests de integración toman mucho tiempo
**Mitigación:** Usar Testcontainers en paralelo, optimizar setup

### Riesgo 2: Refresh tokens ya existe pero no documentado
**Mitigación:** Verificar código antes de implementar

### Riesgo 3: Proyectos consumidores no compilan con v0.7.0
**Mitigación:** Validar temprano (día 2 de Sprint 3), arreglar antes de congelar

### Riesgo 4: Coverage <85% después de Sprint 2
**Mitigación:** Usar Sprint 3 día 1 como buffer para agregar tests

---

## 📞 Próximos Pasos

1. **Aprobar este plan** con el equipo
2. **Ejecutar Sprint 0** (2-3 horas)
3. **Iniciar Sprint 1** (fecha a definir)
4. **Seguir checklist** en `07-CHECKLIST_EJECUCION.md`

---

**Documento generado:** 15 de Noviembre, 2025  
**Próximo documento:** `06-VERSION_FINAL_CONGELADA.md`
