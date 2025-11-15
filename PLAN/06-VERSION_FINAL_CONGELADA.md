# Versión Final Congelada: v0.7.0

## 🎯 Definición de la Versión Congelada

**Versión:** v0.7.0  
**Fecha objetivo de congelamiento:** +3 semanas desde inicio de sprints  
**Estado:** **FROZEN** - No modificable hasta post-MVP

---

## 📦 Módulos Incluidos (Todos en v0.7.0)

| Módulo | Versión | Features Clave | Estado Actual |
|--------|---------|----------------|---------------|
| auth/ | v0.7.0 | JWT, roles, refresh tokens | v0.5.0 → v0.7.0 |
| logger/ | v0.7.0 | Structured logging, Zap | v0.5.0 → v0.7.0 |
| common/ | v0.7.0 | Errors, types, validator | v0.5.0 → v0.7.0 |
| config/ | v0.7.0 | Multi-environment config | v0.5.0 → v0.7.0 |
| bootstrap/ | v0.7.0 | App initialization | v0.5.0 → v0.7.0 |
| lifecycle/ | v0.7.0 | Graceful shutdown | v0.5.0 → v0.7.0 |
| middleware/gin/ | v0.7.0 | JWT, logging, CORS middlewares | v0.5.0 → v0.7.0 |
| messaging/rabbit/ | v0.7.0 | Publisher, consumer, **DLQ** | v0.5.0 → v0.7.0 |
| database/postgres/ | v0.7.0 | GORM, transactions, tests | v0.5.0 → v0.7.0 |
| database/mongodb/ | v0.7.0 | MongoDB client, pooling | v0.5.0 → v0.7.0 |
| testing/ | v0.7.0 | Testcontainers, helpers | v0.6.2 → v0.7.0 |
| **evaluation/** | **v0.7.0** | **Assessment, Question, Attempt (NUEVO)** | **v0.1.0 → v0.7.0** |

**Total módulos:** 12 (11 existentes + 1 nuevo)

---

## 🔒 Contrato de Congelamiento

### Qué significa "CONGELADO"

#### ✅ Permitido (v0.7.x)
- 🐛 **Bug fixes críticos** (v0.7.1, v0.7.2, v0.7.3...)
  - Errores que rompen funcionalidad existente
  - Vulnerabilidades de seguridad
  - Crashes o deadlocks
  
- 📝 **Documentación**
  - Mejorar godoc comments
  - Agregar ejemplos en README
  - Aclarar confusiones

- 🧪 **Tests**
  - Agregar tests para aumentar coverage
  - Arreglar tests flaky

#### ❌ NO Permitido hasta post-MVP
- ⛔ **Nuevas features**
  - NO agregar nuevos métodos públicos
  - NO agregar nuevos módulos
  - NO agregar nuevas structs exportadas

- ⛔ **Breaking changes**
  - NO cambiar signatures de funciones públicas
  - NO renombrar structs/campos exportados
  - NO cambiar comportamiento de APIs existentes
  - NO modificar go.mod de forma incompatible

- ⛔ **Refactoring mayor**
  - NO reestructurar módulos
  - NO cambiar arquitectura interna si afecta API pública

---

### Proceso de Bug Fixes (v0.7.x)

#### 1. Identificar Bug Crítico
**Criterios de criticidad:**
- ¿Rompe funcionalidad existente? → CRÍTICO
- ¿Causa crash/panic? → CRÍTICO
- ¿Vulnerabilidad de seguridad? → CRÍTICO
- ¿Solo afecta performance? → NO crítico (post-MVP)

#### 2. Fix y Release
```bash
# Crear rama de hotfix
git checkout main
git checkout -b hotfix/v0.7.1-bug-description

# Implementar fix
# ...

# Tests
go test ./...

# Commit
git commit -m "fix(module): description of bug fix (#issue)"

# Merge a main
git checkout main
git merge hotfix/v0.7.1-bug-description

# Tag SOLO el módulo afectado
git tag module/v0.7.1
git push origin main module/v0.7.1

# Merge a dev
git checkout dev
git merge main
git push origin dev
```

**⚠️ Importante:** Solo se taggea el módulo afectado, NO todos los módulos.

---

## 📖 Cómo Consumir (Proyectos Dependientes)

### go.mod Recomendado para api-mobile

```go
module github.com/EduGoGroup/edugo-api-mobile

go 1.24

require (
    // Módulos de shared en v0.7.0 coordinado
    github.com/EduGoGroup/edugo-shared/auth          v0.7.0
    github.com/EduGoGroup/edugo-shared/logger        v0.7.0
    github.com/EduGoGroup/edugo-shared/common        v0.7.0
    github.com/EduGoGroup/edugo-shared/config        v0.7.0
    github.com/EduGoGroup/edugo-shared/bootstrap     v0.7.0
    github.com/EduGoGroup/edugo-shared/lifecycle     v0.7.0
    github.com/EduGoGroup/edugo-shared/middleware/gin v0.7.0
    github.com/EduGoGroup/edugo-shared/messaging/rabbit v0.7.0
    github.com/EduGoGroup/edugo-shared/database/postgres v0.7.0
    github.com/EduGoGroup/edugo-shared/database/mongodb v0.7.0
    github.com/EduGoGroup/edugo-shared/testing       v0.7.0  // Solo para tests
    github.com/EduGoGroup/edugo-shared/evaluation    v0.7.0  // NUEVO
    
    // Otros paquetes
    github.com/gin-gonic/gin v1.10.0
    gorm.io/gorm v1.25.0
    // ...
)
```

**Instalación:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Opción 1: Individual (recomendado)
go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0
go get github.com/EduGoGroup/edugo-shared/evaluation@v0.7.0
# ... otros

# Opción 2: Script (crear script/update-shared.sh)
#!/bin/bash
SHARED_VERSION="v0.7.0"
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
    echo "Updating $module to $SHARED_VERSION..."
    go get "github.com/EduGoGroup/edugo-shared/$module@$SHARED_VERSION"
done

go mod tidy
```

---

### go.mod Recomendado para api-admin

```go
require (
    github.com/EduGoGroup/edugo-shared/auth          v0.7.0
    github.com/EduGoGroup/edugo-shared/logger        v0.7.0
    github.com/EduGoGroup/edugo-shared/common        v0.7.0
    github.com/EduGoGroup/edugo-shared/config        v0.7.0
    github.com/EduGoGroup/edugo-shared/database/postgres v0.7.0
    github.com/EduGoGroup/edugo-shared/lifecycle     v0.7.0
    // MongoDB NO requerido por api-admin
    // evaluation NO requerido por api-admin
)
```

---

### go.mod Recomendado para worker

```go
require (
    github.com/EduGoGroup/edugo-shared/logger        v0.7.0
    github.com/EduGoGroup/edugo-shared/common        v0.7.0
    github.com/EduGoGroup/edugo-shared/config        v0.7.0
    github.com/EduGoGroup/edugo-shared/messaging/rabbit v0.7.0  // CON DLQ
    github.com/EduGoGroup/edugo-shared/database/mongodb v0.7.0
    github.com/EduGoGroup/edugo-shared/evaluation    v0.7.0  // Para generar assessments
    github.com/EduGoGroup/edugo-shared/lifecycle     v0.7.0
    // auth NO crítico para worker
    // database/postgres OPCIONAL (solo auditoría)
)
```

---

## 📋 Changelog de v0.7.0

```markdown
# Changelog - edugo-shared

## [0.7.0] - 2025-11-XX - 🔒 FROZEN RELEASE

### 🎉 Versión Congelada para MVP de EduGo

Esta versión marca la **base estable y congelada** para el ecosistema EduGo.
- ✅ Todos los módulos en v0.7.0
- ✅ Coverage global >85%
- ✅ Tests completos y pasando
- ✅ Validado con api-mobile, api-admin, worker
- ⚠️ **NO se agregarán features nuevas hasta post-MVP**

---

### 🆕 Added

#### NEW MODULE: evaluation/ (v0.7.0)
- Assessment struct con validación
- Question struct con QuestionType enum (multiple_choice, true_false, short_answer)
- QuestionOption struct para opciones de respuesta
- Attempt struct con scoring automático
- Answer struct para respuestas de estudiantes
- Helper methods: Validate, IsPublished, GetCorrectOptions, CalculatePercentage, CheckPassed
- Comprehensive unit tests (>90% coverage)

**Uso:**
go
import "github.com/EduGoGroup/edugo-shared/evaluation"

assessment := evaluation.Assessment{
    ID: uuid.New(),
    Title: "Quiz de Matemáticas",
    PassingScore: 70,
}


#### messaging/rabbit/ (v0.5.0 → v0.7.0)
- **Dead Letter Queue (DLQ) support**
- DLQConfig struct con retry configurable
- ConsumeWithDLQ() method con retry automático
- Exponential backoff support
- sendToDLQ() y setupDLQ() helpers

**Uso:**
go
config := rabbit.ConsumerConfig{
    DLQ: rabbit.DLQConfig{
        Enabled: true,
        MaxRetries: 3,
        DLXExchange: "dlx",
    },
}
consumer.ConsumeWithDLQ(handler)


#### auth/ (v0.5.0 → v0.7.0)
- **Refresh token support**
- GenerateTokenPair() retorna access + refresh tokens
- RefreshAccessToken() para renovar access token
- TokenPair struct (access, refresh, expires_in)

**Uso:**
go
pair, _ := jwtManager.GenerateTokenPair(userID, email, role)
// pair.AccessToken (15 min)
// pair.RefreshToken (7 días)

newAccessToken, _ := jwtManager.RefreshAccessToken(pair.RefreshToken)


---

### 🧪 Changed

#### Tests Coverage Improvements
- **database/postgres/**: 2% → >80% coverage
- **logger/**: 0% → >80% coverage
- **common/errors**: 0% → >80% coverage
- **common/types**: 0% → >80% coverage
- **common/validator**: 0% → >80% coverage
- **config/**: 32.9% → >80% coverage
- **bootstrap/**: 29.9% → >80% coverage
- **database/mongodb/**: Validado con integration tests

#### All modules bumped to v0.7.0
- Coordinated release para todos los 12 módulos
- Versionado consistente en todo el ecosistema

---

### 🐛 Fixed
- **auth/, middleware/gin/**: Fixed broken dependencies (go mod tidy ejecutado)
- **database/postgres/**: Mejorado manejo de conexiones en pool
- **messaging/rabbit/**: Mensajes con errores ya no se reencolan infinitamente

---

### 📦 Dependencies
- Go 1.24.10 (todos los módulos)
- github.com/google/uuid v1.6.0
- github.com/golang-jwt/jwt/v5 (auth)
- go.uber.org/zap (logger)
- github.com/rabbitmq/amqp091-go (messaging/rabbit)
- gorm.io/gorm (database/postgres)
- go.mongodb.org/mongo-driver (database/mongodb)
- github.com/testcontainers/testcontainers-go (testing)

---

### 🎯 Metrics
- **Total modules:** 12
- **Global coverage:** >85%
- **Total test files:** 30+
- **Lines of code:** ~5000
- **Tests passing:** 100%

---

### ✅ Validation
- ✅ api-mobile compila sin errores con shared v0.7.0
- ✅ api-admin compila sin errores con shared v0.7.0
- ✅ worker compila sin errores con shared v0.7.0
- ✅ All tests passing en CI/CD
- ✅ golangci-lint: 0 warnings

---

### 🚀 Migration Guide

**From v0.5.0 to v0.7.0:**

bash
# Actualizar todos los módulos
go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0
go get github.com/EduGoGroup/edugo-shared/evaluation@v0.7.0
go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@v0.7.0
go mod tidy


**Breaking changes:** NINGUNO (backward compatible con v0.5.0)

**New features to adopt:**
1. Usar `evaluation/` para modelos de assessments
2. Habilitar DLQ en RabbitMQ consumers:
   go
   config.DLQ.Enabled = true
   
3. Usar refresh tokens en autenticación:
   go
   pair, _ := jwtManager.GenerateTokenPair(...)
   

---

## 🗺️ Roadmap Post-Congelamiento

### v0.8.0 (Post-MVP) - FUTURO
**Cuando:** Después de lanzar MVP a producción

**Features candidatas:**
- ⏱️ Performance optimizations
  - Caché de tokens JWT
  - Connection pooling mejorado
  - Batch operations en DB

- 📊 Observability
  - Prometheus metrics
  - OpenTelemetry tracing
  - Health check endpoints

- 🔒 Security enhancements
  - Token rotation automático
  - Rate limiting helpers
  - Audit log middleware

- 🧪 Testing utilities
  - Más helpers de Testcontainers
  - Mock generators
  - Fixtures library

**Nota:** Estas features NO se implementarán en v0.7.x

---

### v1.0.0 (Producción Estable) - FUTURO
**Cuando:** Cuando todos los servicios estén en producción y estables

**Garantías de v1.0.0:**
- 🔒 API pública 100% estable
- 📚 Documentación completa
- 🛡️ Soporte LTS (Long Term Support)
- ⚠️ Breaking changes solo en v2.0.0

---

## 📞 Soporte y Mantenimiento

### Durante el período congelado (v0.7.x)

**Issues:**
- Bug reports: https://github.com/EduGoGroup/edugo-shared/issues
- Template: Usar `bug-report.md`

**Pull Requests:**
- Solo para bug fixes críticos
- Requiere aprobación de 2+ maintainers
- Debe pasar todos los tests en CI/CD

**Releases:**
- Bug fixes: v0.7.1, v0.7.2, v0.7.3...
- Frecuencia: Según necesidad (no scheduled)

---

## ✅ Checklist de Congelamiento

Antes de declarar v0.7.0 como FROZEN:

- [ ] Todos los módulos tienen tag v0.7.0
- [ ] GitHub Release publicado
- [ ] CHANGELOG.md actualizado
- [ ] README.md actualizado con v0.7.0
- [ ] Coverage global >85% verificado
- [ ] Tests 100% passing en CI/CD
- [ ] api-mobile compila exitosamente
- [ ] api-admin compila exitosamente
- [ ] worker compila exitosamente
- [ ] Documentación de cada módulo actualizada
- [ ] go.mod.example creado para cada consumidor
- [ ] Este documento (06-VERSION_FINAL_CONGELADA.md) revisado

---

**Documento generado:** 15 de Noviembre, 2025  
**Última actualización:** Pre-congelamiento  
**Próximo documento:** `07-CHECKLIST_EJECUCION.md`
