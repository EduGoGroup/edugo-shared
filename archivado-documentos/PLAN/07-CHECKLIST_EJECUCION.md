# Checklist de Ejecución - edugo-shared Final

## 🎯 Propósito

Este es el **checklist ejecutable** que el desarrollador debe seguir paso a paso para completar el plan definitivo de edugo-shared y llegar a la versión congelada v0.7.0.

**Instrucciones:**
1. Marcar `[x]` cuando completes cada tarea
2. NO saltarse pasos
3. Si un paso falla, DETENER y resolver antes de continuar
4. Actualizar este archivo conforme avanzas (hacer commits del checklist)

---

## 📋 Fase 1: Preparación (Hoy - 2-3 horas)

### Paso 1.1: Lectura Inicial

- [ ] Leer `00-README.md` (este plan)
- [ ] Leer `01-ESTADO_ACTUAL.md` (estado del 15 Nov 2025)
- [ ] Leer `02-NECESIDADES_CONSOLIDADAS.md` (requisitos)
- [ ] Leer `03-MODULOS_FALTANTES.md` (evaluation/)
- [ ] Leer `04-FEATURES_FALTANTES.md` (DLQ, refresh tokens)
- [ ] Leer `05-PLAN_SPRINTS.md` (sprints detallados)
- [ ] Leer `06-VERSION_FINAL_CONGELADA.md` (objetivo final)

**Tiempo:** 1-2 horas

**Validación:** Entiendes el plan completo

---

### Paso 1.2: Verificar Acceso a Repositorio

- [ ] Abrir terminal
- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`
- [ ] `git status` (debe funcionar sin errores)
- [ ] `git remote -v` (verificar que apunta a EduGoGroup/edugo-shared)

**Validación:** Tienes acceso al repo

---

### Paso 1.3: Backup del Estado Actual

- [ ] `git checkout dev`
- [ ] `git pull origin dev`
- [ ] `git checkout -b backup/pre-v0.7.0-$(date +%Y%m%d)`
- [ ] `git push origin backup/pre-v0.7.0-$(date +%Y%m%d)`

**Validación:** Rama de backup creada

---

### Paso 1.4: Preparar Ambiente Local

- [ ] Go version: `go version` (debe ser 1.24+)
- [ ] Docker: `docker --version` (para Testcontainers)
- [ ] Docker running: `docker ps` (debe listar contenedores)

**Validación:** Ambiente listo

---

## 📋 Fase 2: Auditoría (Día 1 - 2-3 horas)

### Sprint 0 - Auditoría Completa

#### Paso 2.1: Sincronizar Ramas

- [ ] `git checkout main && git pull origin main`
- [ ] `git log --oneline -1` (anotar hash: ____________)
- [ ] `git checkout dev && git pull origin dev`
- [ ] `git log --oneline -1` (anotar hash: ____________)
- [ ] `git diff main dev --stat`
- [ ] **Decisión:** ¿Mergear o trabajar desde dev? → **Trabajar desde dev**

**Validación:** Conoces estado de ramas

---

#### Paso 2.2: Fix Dependencias Rotas

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared/auth`
- [ ] `go mod tidy`
- [ ] `go test ./...` (anotar resultado: ____________)
- [ ] `cd ../middleware/gin`
- [ ] `go mod tidy`
- [ ] `go test ./...` (anotar resultado: ____________)
- [ ] `cd ../..`
- [ ] `git add auth/go.mod auth/go.sum middleware/gin/go.mod middleware/gin/go.sum`
- [ ] `git commit -m "fix(deps): execute go mod tidy on auth and middleware/gin"`
- [ ] `git push origin dev`

**Validación:** go mod tidy ejecutado, tests ejecutables

---

#### Paso 2.3: Ejecutar Suite Completa de Tests

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`
- [ ] `make test-all-modules > test-results-baseline.txt 2>&1`
- [ ] `cat test-results-baseline.txt` (revisar resultados)
- [ ] Anotar módulos failing: ________________________________
- [ ] `make coverage-all-modules > coverage-baseline.txt 2>&1`
- [ ] `cat coverage-baseline.txt` (revisar coverage)

**Validación:** Tienes baseline de tests y coverage

---

#### Paso 2.4: Crear Issues en GitHub

- [ ] Ir a https://github.com/EduGoGroup/edugo-shared/issues/new
- [ ] Issue #1: "Create evaluation/ module (P0)"
  ```
  **Objetivo:** Crear módulo evaluation/ con Assessment, Question, Attempt models
  **Prioridad:** P0
  **Bloqueante para:** api-mobile, worker
  **Estimación:** 4-5 horas
  **Milestone:** v0.7.0
  ```
  
- [ ] Issue #2: "Implement DLQ support in messaging/rabbit (P0)"
  ```
  **Objetivo:** Agregar Dead Letter Queue support con retry logic
  **Prioridad:** P0
  **Bloqueante para:** worker
  **Estimación:** 3-5 horas
  **Milestone:** v0.7.0
  ```
  
- [ ] Issue #3: "Increase database/postgres coverage to >80% (P0)"
- [ ] Issue #4: "Add tests to logger, common/* modules (P1)"
- [ ] Issue #5: "Implement/verify refresh tokens in auth (P1)"
- [ ] Issue #6: "Increase coverage in config, bootstrap (P2)"

**Validación:** 6 issues creados

---

## 📋 Fase 3: Sprint 1 - Módulos Críticos (Semana 1)

### Día 1-2: Crear evaluation/ (4-5 horas)

#### Paso 3.1: Setup Módulo

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`
- [ ] `mkdir evaluation`
- [ ] `cd evaluation`
- [ ] `go mod init github.com/EduGoGroup/edugo-shared/evaluation`
- [ ] Editar `go.mod`:
  ```go
  module github.com/EduGoGroup/edugo-shared/evaluation
  
  go 1.24
  
  require (
      github.com/google/uuid v1.6.0
  )
  ```
- [ ] `go mod tidy`

**Validación:** Módulo inicializado

---

#### Paso 3.2: Implementar assessment.go

- [ ] Copiar código de `03-MODULOS_FALTANTES.md` sección "assessment.go"
- [ ] Crear `evaluation/assessment.go`
- [ ] Pegar código
- [ ] Ajustar imports si es necesario
- [ ] `go build .` (debe compilar)

**Validación:** assessment.go compila

---

#### Paso 3.3: Implementar question.go

- [ ] Copiar código de `03-MODULOS_FALTANTES.md` sección "question.go"
- [ ] Crear `evaluation/question.go`
- [ ] Pegar código
- [ ] `go build .`

**Validación:** question.go compila

---

#### Paso 3.4: Implementar attempt.go

- [ ] Copiar código de `03-MODULOS_FALTANTES.md` sección "attempt.go"
- [ ] Crear `evaluation/attempt.go`
- [ ] Pegar código
- [ ] `go build .`

**Validación:** attempt.go compila

---

#### Paso 3.5: Tests de evaluation/

- [ ] Crear `evaluation/assessment_test.go` (copiar de `03-MODULOS_FALTANTES.md`)
- [ ] Crear `evaluation/question_test.go`
- [ ] Crear `evaluation/attempt_test.go`
- [ ] `go test -v -cover ./...`
- [ ] Coverage debe ser >80%
- [ ] Si coverage <80%: agregar más tests

**Validación:** Tests pasan, coverage >80%

---

#### Paso 3.6: Documentación de evaluation/

- [ ] Crear `evaluation/README.md`:
  ```markdown
  # evaluation - EduGo Shared
  
  Módulo que define modelos compartidos para el sistema de evaluaciones.
  
  ## Uso
  
  go
  import "github.com/EduGoGroup/edugo-shared/evaluation"
  
  assessment := evaluation.Assessment{
      Title: "Quiz de Matemáticas",
      PassingScore: 70,
  }
  
  
  ## Modelos
  
  - **Assessment**: Cuestionario
  - **Question**: Pregunta con opciones
  - **Attempt**: Intento de un estudiante
  ```

**Validación:** README.md creado

---

#### Paso 3.7: Commit y Tag de evaluation/

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`
- [ ] `git add evaluation/`
- [ ] `git commit -m "feat(evaluation): create evaluation module with Assessment, Question, Attempt models

Closes #1

- Add Assessment struct with validation
- Add Question struct with QuestionType enum
- Add QuestionOption for multiple choice
- Add Attempt struct with scoring logic
- Add Answer struct for student responses
- Comprehensive unit tests (>90% coverage)"`
- [ ] `git push origin dev`
- [ ] `git tag evaluation/v0.1.0`
- [ ] `git push origin evaluation/v0.1.0`

**Validación:** evaluation/v0.1.0 publicado

---

### Día 3-4: Implementar DLQ en messaging/rabbit/ (3-5 horas)

#### Paso 3.8: Implementar dlq.go

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared/messaging/rabbit`
- [ ] Crear `dlq.go` (copiar de `04-FEATURES_FALTANTES.md`)
- [ ] `go build .`

**Validación:** dlq.go compila

---

#### Paso 3.9: Modificar consumer.go

- [ ] Editar `consumer.go` (seguir spec en `04-FEATURES_FALTANTES.md`)
- [ ] Agregar `DLQ` field a `ConsumerConfig`
- [ ] Implementar `ConsumeWithDLQ()`
- [ ] Implementar `setupDLQ()`
- [ ] Implementar `sendToDLQ()`
- [ ] Implementar `getRetryCount()`
- [ ] `go build .`

**Validación:** consumer.go compila

---

#### Paso 3.10: Tests de DLQ

- [ ] Crear `dlq_test.go` (unit tests de backoff)
- [ ] Crear `consumer_dlq_test.go` (integration con Testcontainers)
- [ ] `go test -v -cover ./...`
- [ ] Verificar que DLQ tests pasan

**Validación:** Tests de DLQ pasan

---

#### Paso 3.11: Commit y Tag de messaging/rabbit/

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`
- [ ] `git add messaging/rabbit/`
- [ ] `git commit -m "feat(messaging/rabbit): add Dead Letter Queue (DLQ) support with retry logic

Closes #2

- Add DLQConfig struct with configurable retries
- Implement ConsumeWithDLQ with automatic retry
- Add exponential backoff support
- Implement sendToDLQ and setupDLQ helpers
- Integration tests with Testcontainers"`
- [ ] `git push origin dev`
- [ ] `git tag messaging/rabbit/v0.6.0`
- [ ] `git push origin messaging/rabbit/v0.6.0`

**Validación:** messaging/rabbit/v0.6.0 publicado

---

### Día 5: Aumentar Coverage en database/postgres/ (4-6 horas)

#### Paso 3.12: Crear Tests de Integración

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared/database/postgres`
- [ ] Crear `postgres_integration_test.go`:
  - Setup con Testcontainers PostgreSQL
  - Test de conexión
  - Test de transacciones (Begin, Commit, Rollback)
  - Test de health check
  - Test de reconnection
- [ ] `go test -v -cover ./...`
- [ ] Coverage debe ser >80% (actualmente 2%)

**Validación:** Coverage >80%

---

#### Paso 3.13: Commit y Tag de database/postgres/

- [ ] `git add database/postgres/`
- [ ] `git commit -m "test(database/postgres): increase coverage from 2% to >80% with integration tests

Closes #3

- Add Testcontainers setup
- Test connection pooling
- Test transactions (Begin, Commit, Rollback)
- Test health checks
- Test reconnection logic"`
- [ ] `git push origin dev`
- [ ] `git tag database/postgres/v0.6.0`
- [ ] `git push origin database/postgres/v0.6.0`

**Validación:** database/postgres/v0.6.0 publicado

---

### Validación de Sprint 1

- [ ] evaluation/v0.1.0 existe: `git tag -l | grep evaluation`
- [ ] messaging/rabbit/v0.6.0 existe: `git tag -l | grep messaging`
- [ ] database/postgres/v0.6.0 existe: `git tag -l | grep postgres`
- [ ] Todos los tests P0 pasan

---

## 📋 Fase 4: Sprint 2 - Features Faltantes (Semana 2)

### Día 1: Agregar Tests a logger/ (3-4 horas)

#### Paso 4.1: Implementar logger_test.go

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared/logger`
- [ ] Crear `logger_test.go`:
  - Test de creación de logger
  - Test de niveles (Debug, Info, Warn, Error)
  - Test de formatos (JSON, Console)
  - Test de context fields
- [ ] `go test -v -cover ./...`
- [ ] Coverage >80%

**Validación:** logger/ tiene >80% coverage

---

#### Paso 4.2: Commit y Tag de logger/

- [ ] `git commit -m "test(logger): add comprehensive unit tests (>80% coverage)

Part of #4"`
- [ ] `git tag logger/v0.6.0`
- [ ] `git push origin dev logger/v0.6.0`

---

### Día 2: Agregar Tests a common/* (6-8 horas)

#### Paso 4.3: Tests de common/errors

- [ ] Crear `common/errors/errors_test.go`
- [ ] `cd common/errors && go test -v -cover ./...`

#### Paso 4.4: Tests de common/types

- [ ] Crear `common/types/uuid_test.go`
- [ ] Crear `common/types/enum/enum_test.go`
- [ ] `cd common/types && go test -v -cover ./...`

#### Paso 4.5: Tests de common/validator

- [ ] Crear `common/validator/validator_test.go`
- [ ] `cd common/validator && go test -v -cover ./...`

#### Paso 4.6: Verificar Coverage Global de common/

- [ ] `cd common && go test -v -cover ./...`
- [ ] Coverage >80% en todos los submódulos

**Validación:** common/* tiene >80% coverage

---

#### Paso 4.7: Commit y Tag de common/

- [ ] `git commit -m "test(common): add tests to errors, types, validator (>80% coverage)

Closes #4"`
- [ ] `git tag common/v0.6.0`
- [ ] `git push origin dev common/v0.6.0`

---

### Día 3: Refresh Tokens en auth/ (2-3 horas)

#### Paso 4.8: Verificar si Refresh Tokens Existe

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared/auth`
- [ ] `grep -r "RefreshToken" .`
- [ ] **Si NO existe:** Implementar según `04-FEATURES_FALTANTES.md`
- [ ] **Si YA existe:** Verificar tests y documentar

#### Paso 4.9: Implementar (si no existe)

- [ ] Editar `jwt.go`
- [ ] Agregar `TokenPair` struct
- [ ] Agregar `RefreshClaims` struct
- [ ] Implementar `GenerateTokenPair()`
- [ ] Implementar `RefreshAccessToken()`
- [ ] Implementar `ValidateRefreshToken()`

#### Paso 4.10: Tests de Refresh Tokens

- [ ] Crear `refresh_token_test.go`
- [ ] `go test -v -cover ./...`

**Validación:** Refresh tokens implementado y testeado

---

#### Paso 4.11: Commit y Tag de auth/

- [ ] `git commit -m "feat(auth): implement refresh token support

Closes #5

- Add TokenPair struct
- Implement GenerateTokenPair()
- Implement RefreshAccessToken()
- Comprehensive tests"`
- [ ] `git tag auth/v0.6.0`
- [ ] `git push origin dev auth/v0.6.0`

---

### Día 4: Validar database/mongodb/ (2-3 horas)

#### Paso 4.12: Tests de MongoDB

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared/database/mongodb`
- [ ] Si no existe `mongodb_integration_test.go`: Crear
- [ ] `go test -v -cover ./...`
- [ ] Coverage >80%

#### Paso 4.13: Commit y Tag

- [ ] `git commit -m "test(database/mongodb): add integration tests with Testcontainers"`
- [ ] `git tag database/mongodb/v0.6.0`
- [ ] `git push origin dev database/mongodb/v0.6.0`

---

### Día 5: Buffer / Refactoring

- [ ] Revisar todos los commits de Sprint 2
- [ ] `make test-all-modules`
- [ ] Arreglar failing tests
- [ ] Actualizar documentación

---

## 📋 Fase 5: Sprint 3 - Consolidación (3 días)

### Día 1: Coverage P2 (4-6 horas)

#### Paso 5.1: Aumentar Coverage en config/

- [ ] `cd config`
- [ ] Agregar tests hasta >80%
- [ ] `go test -v -cover ./...`
- [ ] `git tag config/v0.6.0`

#### Paso 5.2: Aumentar Coverage en bootstrap/

- [ ] `cd bootstrap`
- [ ] Agregar tests hasta >80%
- [ ] `go test -v -cover ./...`
- [ ] `git tag bootstrap/v0.6.0`

- [ ] Commit y push tags

---

### Día 2: Validación Completa (Full Day)

#### Paso 5.3: Suite Completa de Tests

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`
- [ ] `make test-all-modules | tee test-results-final.txt`
- [ ] **Verificar:** 0 tests failing
- [ ] Si hay failing: ARREGLAR antes de continuar

#### Paso 5.4: Coverage Global

- [ ] `make coverage-all-modules | tee coverage-report-final.txt`
- [ ] Calcular coverage promedio
- [ ] **Target:** >85% global
- [ ] Si <85%: agregar más tests

#### Paso 5.5: Validar Compilación de Consumidores

**api-mobile:**
- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile`
- [ ] Editar `go.mod` (agregar evaluation@v0.1.0, etc.)
- [ ] `go get github.com/EduGoGroup/edugo-shared/evaluation@v0.1.0`
- [ ] `go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@v0.6.0`
- [ ] `go build ./cmd/api-mobile`
- [ ] **Debe compilar SIN errores**

**api-admin:**
- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion`
- [ ] `go build ./cmd/api-admin`
- [ ] **Debe compilar SIN errores**

**worker:**
- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-worker`
- [ ] `go get github.com/EduGoGroup/edugo-shared/evaluation@v0.1.0`
- [ ] `go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@v0.6.0`
- [ ] `go build ./cmd/worker`
- [ ] **Debe compilar SIN errores**

**Validación:** Todos compilan exitosamente

---

### Día 3: Release v0.7.0 (2-3 horas)

#### Paso 5.6: Crear Rama de Release

- [ ] `cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`
- [ ] `git checkout dev`
- [ ] `git checkout -b release/v0.7.0`

#### Paso 5.7: Actualizar CHANGELOG.md

- [ ] Editar `CHANGELOG.md`
- [ ] Copiar contenido de `06-VERSION_FINAL_CONGELADA.md` sección "Changelog"
- [ ] Agregar fecha real
- [ ] `git add CHANGELOG.md`
- [ ] `git commit -m "docs: update CHANGELOG for v0.7.0 release"`

#### Paso 5.8: Actualizar README.md

- [ ] Editar `README.md`
- [ ] Actualizar sección de instalación con v0.7.0
- [ ] Agregar nota de versión congelada
- [ ] `git commit -m "docs: update README for v0.7.0 frozen release"`

#### Paso 5.9: Mergear a main

- [ ] `git push origin release/v0.7.0`
- [ ] `git checkout main`
- [ ] `git pull origin main`
- [ ] `git merge release/v0.7.0`
- [ ] `git push origin main`

#### Paso 5.10: Crear Tags Coordinados v0.7.0

```bash
#!/bin/bash
# Script para crear todos los tags v0.7.0

MODULES=(
  "auth"
  "logger"
  "common"
  "config"
  "bootstrap"
  "lifecycle"
  "middleware/gin"
  "messaging/rabbit"
  "database/postgres"
  "database/mongodb"
  "testing"
  "evaluation"
)

for module in "${MODULES[@]}"; do
  echo "Tagging $module/v0.7.0..."
  git tag "$module/v0.7.0"
done

echo "Pushing all tags..."
git push origin --tags

echo "Done! All modules tagged as v0.7.0"
```

- [ ] Crear script `scripts/tag-v0.7.0.sh`
- [ ] `chmod +x scripts/tag-v0.7.0.sh`
- [ ] `./scripts/tag-v0.7.0.sh`
- [ ] Verificar: `git tag -l | grep v0.7.0 | wc -l` (debe ser 12)

**Validación:** 12 módulos en v0.7.0

---

#### Paso 5.11: Mergear main a dev

- [ ] `git checkout dev`
- [ ] `git merge main`
- [ ] `git push origin dev`

**Validación:** main y dev sincronizados

---

#### Paso 5.12: Crear GitHub Release

- [ ] Ir a https://github.com/EduGoGroup/edugo-shared/releases/new
- [ ] Tag: `v0.7.0`
- [ ] Title: `v0.7.0 - Frozen Release for EduGo MVP`
- [ ] Description: Copiar de `CHANGELOG.md` sección v0.7.0
- [ ] Marcar checkbox "Set as the latest release"
- [ ] Click "Publish release"

**Validación:** Release publicado en GitHub

---

## 📋 Fase 6: Validaciones Finales (30 minutos)

### Checklist Final Pre-Congelamiento

- [ ] Todos los módulos tienen tag v0.7.0
- [ ] GitHub Release v0.7.0 publicado
- [ ] CHANGELOG.md actualizado con v0.7.0
- [ ] README.md menciona v0.7.0 y congelamiento
- [ ] Coverage global >85% documentado
- [ ] `make test-all-modules` pasa 100%
- [ ] api-mobile compila con shared v0.7.0
- [ ] api-admin compila con shared v0.7.0
- [ ] worker compila con shared v0.7.0
- [ ] Cada módulo tiene README.md actualizado
- [ ] Issues #1-#6 cerrados en GitHub
- [ ] Este checklist completado al 100%

---

## 🎉 CONGELAMIENTO DECLARADO

### Último Paso

- [ ] Crear issue "v0.7.0 FROZEN - No new features until post-MVP"
- [ ] Pin issue en GitHub
- [ ] Notificar al equipo: "shared v0.7.0 is FROZEN"

---

## 📊 Métricas Finales (Completar al final)

```
Fecha de congelamiento: ___________
Tiempo total invertido: ___________ horas
Módulos en v0.7.0: ___________
Coverage global: ___________%
Tests passing: ___________%
Issues cerrados: ___________
Commits totales: ___________
```

---

**¡FELICITACIONES!**

Has completado exitosamente el Plan Definitivo de edugo-shared.

La versión v0.7.0 está CONGELADA y lista para ser usada por api-mobile, api-admin y worker.

---

**Documento generado:** 15 de Noviembre, 2025  
**Tipo:** Checklist ejecutable  
**Estado:** Listo para usar
