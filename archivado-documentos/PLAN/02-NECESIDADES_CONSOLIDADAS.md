# Necesidades Consolidadas de Proyectos Consumidores

## 🎯 Objetivo de Este Documento

Consolidar TODAS las necesidades que los proyectos consumidores (api-mobile, api-admin, worker) tienen de `edugo-shared`, identificando:
- ✅ Módulos que existen y cumplen requisitos
- ⚠️ Módulos que existen pero les faltan features
- ❌ Módulos que NO existen y deben crearse

**Fuentes de información:**
- `/Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/api-mobile/01-Context/DEPENDENCIES.md`
- `/Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/api-admin/01-Context/DEPENDENCIES.md`
- `/Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/worker/01-Context/DEPENDENCIES.md`

---

## 📱 API Mobile

### Dependencias Declaradas

**Versión de shared requerida:** v1.3.0+ (según documentación)  
**Versión actual de shared:** v0.5.0 (mayoría de módulos)  
**⚠️ Gap de versión:** Documentación desactualizada o versión real es menor

---

### Módulos Requeridos

#### 1. logger/ → ✅ EXISTE (v0.5.0)

**Para qué:** Logging estructurado en toda la aplicación

**Features necesarias:**
- ✅ Structured logging (Info, Warn, Error, Debug)
- ✅ JSON y Console formats
- ✅ Context-aware logging

**Estado en shared:** ✅ Implementado en v0.5.0

**Gap detectado:** 🔴 **SIN TESTS** (0% coverage)

**Acción:** Agregar tests unitarios

---

#### 2. database/postgres/ → ✅ EXISTE (v0.5.0)

**Para qué:** Conexión a PostgreSQL, GORM, transacciones

**Features necesarias:**
- ✅ Connection pooling
- ✅ Health checks
- ✅ Transaction support
- ✅ GORM integration

**Estado en shared:** ✅ Implementado en v0.5.0

**Gap detectado:** 🔴 **Coverage 2%** (casi sin tests)

**Acción:** Aumentar coverage a >80% con tests de integración

---

#### 3. database/mongodb/ → ✅ EXISTE (v0.5.0)

**Para qué:** Persistencia de resultados de evaluaciones

**Features necesarias:**
- ✅ MongoDB client configuration
- ✅ Connection pooling
- ✅ Collections access

**Estado en shared:** ✅ Implementado en v0.5.0

**Gap detectado:** ⚠️ Tests no verificados

**Acción:** Validar tests con Testcontainers

---

#### 4. auth/ → ✅ EXISTE (v0.5.0)

**Para qué:** JWT validation, claims extraction, roles

**Features necesarias:**
- ✅ JWT generation
- ✅ JWT validation
- ✅ Claims extraction
- ✅ Roles: admin, teacher, student, guardian
- ❌ **Refresh tokens** (mencionados en docs pero NO verificados)

**Estado en shared:** ⚠️ Implementado parcialmente

**Gap detectado:**
- 🔴 `go mod tidy` requerido (tests no ejecutables)
- ⚠️ Refresh tokens: **NO CONFIRMADO** si existe

**Acción:**
1. Fix dependencias (go mod tidy)
2. Verificar si refresh tokens está implementado
3. Si NO existe: Implementar refresh token support

---

#### 5. messaging/rabbit/ → ✅ EXISTE (v0.5.0)

**Para qué:** Publicar eventos a RabbitMQ (opcional para MVP)

**Features necesarias:**
- ✅ Publisher interface
- ✅ Connection management
- ❌ **Consumer con prefetch** (existe pero sin DLQ)
- ❌ **Dead Letter Queue (DLQ)** NO implementado

**Estado en shared:** ⚠️ Implementado parcialmente

**Gap detectado:** 🔴 **Sin soporte DLQ** (crítico para worker)

**Acción:** Implementar DLQ support

---

#### 6. common/errors → ✅ EXISTE (v0.5.0)

**Para qué:** Error handling estructurado con HTTP status codes

**Features necesarias:**
- ✅ NotFoundError (404)
- ✅ ValidationError (400)
- ✅ UnauthorizedError (401)
- ✅ InternalError (500)

**Estado en shared:** ✅ Implementado en v0.5.0

**Gap detectado:** 🔴 **SIN TESTS** (0% coverage)

**Acción:** Agregar tests unitarios

---

#### 7. common/types → ✅ EXISTE (v0.5.0)

**Para qué:** UUID, Enums (SystemRole, Status, etc.)

**Features necesarias:**
- ✅ UUID wrapper con JSON marshaling
- ✅ SystemRole enum (admin, teacher, student, guardian)
- ✅ Status enum

**Estado en shared:** ✅ Implementado en v0.5.0

**Gap detectado:** 🔴 **SIN TESTS** (0% coverage)

**Acción:** Agregar tests de marshaling/unmarshaling

---

#### 8. common/validator → ✅ EXISTE (v0.5.0)

**Para qué:** Validación de emails, UUIDs, campos requeridos

**Features necesarias:**
- ✅ Email validation
- ✅ UUID validation
- ✅ Required fields

**Estado en shared:** ✅ Implementado en v0.5.0

**Gap detectado:** 🔴 **SIN TESTS** (0% coverage)

**Acción:** Agregar tests de validación

---

#### 9. config/ → ✅ EXISTE (v0.5.0)

**Para qué:** Configuración multi-ambiente (local, dev, qa, prod)

**Features necesarias:**
- ✅ Viper integration
- ✅ Environment variables loading
- ✅ Multi-environment support

**Estado en shared:** ✅ Implementado en v0.5.0

**Gap detectado:** ⚠️ Coverage 32.9% (bajo)

**Acción:** Aumentar coverage a >80%

---

#### 10. evaluation/ → ❌ **NO EXISTE** (CRÍTICO)

**Para qué:** Modelos compartidos de evaluaciones (Assessment, Question, Attempt)

**Features necesarias:**
- ❌ `Assessment` struct (ID, MaterialID, MongoDocID, Title, etc.)
- ❌ `Question` struct (ID, Text, Type, Options, Points)
- ❌ `QuestionOption` struct (Text, IsCorrect, Position)
- ❌ `Attempt` struct (ID, AssessmentID, StudentID, Score, Submitted)
- ❌ `Answer` struct (QuestionID, AnswerText, IsCorrect, Points)

**Estado en shared:** 🔴 **NO EXISTE**

**Justificación:**
- api-mobile necesita estos modelos para endpoints de evaluaciones
- worker necesita estos modelos para generar quizzes
- Compartir modelos evita duplicación y garantiza consistencia

**Impacto:** **BLOQUEANTE** para implementar sistema de evaluaciones

**Acción:** **CREAR módulo `evaluation/` en v0.7.0**

**Tiempo estimado:** 4-6 horas

---

### Resumen de Gaps en API Mobile

| Módulo | Existe | Features OK | Tests OK | Acción Requerida |
|--------|--------|-------------|----------|------------------|
| logger | ✅ | ✅ | 🔴 0% | Agregar tests |
| database/postgres | ✅ | ✅ | 🔴 2% | Aumentar tests |
| database/mongodb | ✅ | ✅ | ⚠️ | Validar tests |
| auth | ✅ | ⚠️ | 🔴 | Fix deps, verificar refresh tokens |
| messaging/rabbit | ✅ | ⚠️ | ⚠️ | Agregar DLQ |
| common/errors | ✅ | ✅ | 🔴 0% | Agregar tests |
| common/types | ✅ | ✅ | 🔴 0% | Agregar tests |
| common/validator | ✅ | ✅ | 🔴 0% | Agregar tests |
| config | ✅ | ✅ | ⚠️ 32.9% | Aumentar tests |
| **evaluation** | 🔴 NO | - | - | **CREAR módulo** |

**Total gaps:** 10 (9 mejoras + 1 creación)

---

## 🏫 API Admin

### Dependencias Declaradas

**Versión de shared requerida:** v1.3.0+ (según documentación)

---

### Módulos Requeridos

#### 1. logger/ → ✅ EXISTE (v0.5.0)
**Para qué:** Logging de operaciones administrativas

**Estado:** Mismo que api-mobile (0% tests)

---

#### 2. database/postgres/ → ✅ EXISTE (v0.5.0)
**Para qué:** CRUD de escuelas, unidades académicas, jerarquías

**Features especiales necesarias:**
- ✅ GORM integration (existe)
- ✅ Support para CTEs recursivas (PostgreSQL feature, no de shared)

**Estado:** Mismo que api-mobile (2% coverage)

**Nota:** Las CTEs recursivas son feature de PostgreSQL, no requieren nada especial en shared

---

#### 3. auth/ → ✅ EXISTE (v0.5.0)
**Para qué:** Validar que solo admins puedan crear escuelas

**Estado:** Mismo que api-mobile (requiere fix deps)

---

#### 4. common/errors, types, validator → ✅ EXISTE (v0.5.0)
**Para qué:** Error handling, validación

**Estado:** Mismo que api-mobile (0% tests)

---

### Resumen de Gaps en API Admin

| Módulo | Gap vs API Mobile |
|--------|-------------------|
| Todos | **Mismo estado** que api-mobile |

**Conclusión:** API Admin NO introduce nuevos requisitos, solo usa los mismos módulos que api-mobile.

---

## ⚙️ Worker

### Dependencias Declaradas

**Versión de shared requerida:** v1.4.0+ (según documentación)

---

### Módulos Requeridos

#### 1. logger/ → ✅ EXISTE (v0.5.0)
**Para qué:** Logging de procesamiento de materiales

**Estado:** Mismo (0% tests)

---

#### 2. database/postgres/ → ⚠️ OPCIONAL (v0.5.0)
**Para qué:** Auditoría de procesamiento (NO crítico)

**Estado:** Mismo (2% coverage)

---

#### 3. database/mongodb/ → ✅ EXISTE (v0.5.0)
**Para qué:** Guardar resultados de assessments generados por IA

**Features necesarias:**
- ✅ MongoDB client
- ✅ Collection access
- ✅ InsertOne, UpdateOne

**Estado:** Mismo que api-mobile

---

#### 4. messaging/rabbit/ → ✅ EXISTE (v0.5.0)
**Para qué:** **Consumer** (NO publisher) de eventos de procesamiento

**Features necesarias:**
- ✅ Consumer interface
- ✅ Prefetch configuration
- ❌ **Dead Letter Queue (DLQ)** - **CRÍTICO** para Worker
  - Cuando falla procesamiento, mensaje debe ir a DLQ
  - Worker debe poder reintentar con exponential backoff

**Estado:** 🔴 **SIN DLQ** (BLOQUEANTE para Worker)

**Impacto:** Worker no puede manejar errores de procesamiento correctamente

**Acción:** **IMPLEMENTAR DLQ** en messaging/rabbit/ v0.6.0

**Tiempo estimado:** 3-4 horas

---

#### 5. evaluation/ → ❌ **NO EXISTE** (NECESARIO)

**Para qué:** Worker genera `Assessment` y lo guarda en MongoDB

**Features necesarias:**
- Mismo módulo que api-mobile necesita
- Worker **escribe** assessments
- api-mobile **lee** assessments

**Estado:** 🔴 **NO EXISTE**

**Acción:** Crear módulo (mismo que para api-mobile)

---

### Resumen de Gaps en Worker

| Módulo | Existe | Gap Específico del Worker |
|--------|--------|---------------------------|
| messaging/rabbit | ✅ | 🔴 **DLQ crítico** para manejo de errores |
| evaluation | 🔴 NO | Necesario para generar assessments |
| Otros | ✅ | Mismos gaps que api-mobile |

**Conclusión:** Worker tiene 1 gap CRÍTICO adicional: **DLQ en messaging/rabbit**

---

## 📊 Matriz de Dependencias Consolidada

### Módulos Existentes con Gaps

| Módulo | api-mobile | api-admin | worker | Gap Principal | Prioridad |
|--------|------------|-----------|--------|---------------|-----------|
| logger/ | ✅ | ✅ | ✅ | 🔴 0% tests | P1 |
| database/postgres/ | ✅ | ✅ | ⚠️ Opcional | 🔴 2% tests | P0 |
| database/mongodb/ | ✅ | ❌ | ✅ | ⚠️ Tests no validados | P1 |
| auth/ | ✅ | ✅ | ❌ | 🔴 Deps rotas, refresh tokens? | P0 |
| messaging/rabbit/ | ✅ Opcional | ❌ | ✅ **Crítico** | 🔴 **Sin DLQ** | **P0** |
| common/errors | ✅ | ✅ | ✅ | 🔴 0% tests | P1 |
| common/types | ✅ | ✅ | ✅ | 🔴 0% tests | P1 |
| common/validator | ✅ | ✅ | ❌ | 🔴 0% tests | P1 |
| config/ | ✅ | ✅ | ✅ | ⚠️ 32.9% tests | P2 |
| bootstrap/ | ⚠️ | ⚠️ | ⚠️ | ⚠️ 29.9% tests | P2 |
| lifecycle/ | ⚠️ | ⚠️ | ⚠️ | ✅ 91.8% tests | ✅ OK |

**Leyenda:**
- ✅ Requerido y funcionando
- ⚠️ Opcional o con warnings
- ❌ No requerido
- 🔴 Gap crítico
- P0: Crítico (bloquea desarrollo)
- P1: Importante (afecta calidad)
- P2: Mejora (nice to have)

---

### Módulos Nuevos Requeridos

| Módulo | api-mobile | api-admin | worker | Justificación | Prioridad |
|--------|------------|-----------|--------|---------------|-----------|
| **evaluation/** | ✅ **CRÍTICO** | ❌ | ✅ **CRÍTICO** | Modelos compartidos de evaluaciones | **P0** |

---

## 🔍 Análisis de Features Faltantes

### 1. messaging/rabbit - DLQ Support (P0)

**Requerido por:** Worker (crítico), api-mobile (opcional)

**Feature actual:**
```go
// consumer.go (actual)
type Consumer struct {
    connection *amqp.Connection
    channel    *amqp.Channel
}

func (c *Consumer) Consume(queue string, handler func([]byte) error) {
    msgs, _ := c.channel.Consume(queue, ...)
    
    for msg := range msgs {
        if err := handler(msg.Body); err != nil {
            msg.Nack(false, true)  // Requeue indefinidamente
        }
        msg.Ack(false)
    }
}
```

**Problema:** Si un mensaje falla 10 veces, se reencola infinitamente

**Feature necesaria:**
```go
// consumer.go (necesario)
type ConsumerConfig struct {
    Queue            string
    MaxRetries       int    // 3
    DLQExchange      string // "assessment.dlx"
    DLQRoutingKey    string // "assessment.dlq"
    RetryBackoff     time.Duration // Exponential
}

func (c *Consumer) ConsumeWithDLQ(config ConsumerConfig, handler func([]byte) error) {
    // Lógica de retry
    // Si falla > MaxRetries, enviar a DLQ
}
```

**Impacto si no se implementa:**
- Worker crashea en mensajes con errores
- Mensajes se reencolan infinitamente
- No hay visibilidad de mensajes fallidos

**Tiempo estimado:** 3-4 horas

---

### 2. auth - Refresh Tokens (P1)

**Requerido por:** api-mobile (importante), api-admin (importante)

**Feature actual:**
```go
// jwt.go (actual)
func (j *JWTManager) GenerateToken(userID, email string, role enum.SystemRole, expiration time.Duration) (string, error) {
    // Solo access token
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
    // Solo valida access token
}
```

**Feature necesaria:**
```go
// jwt.go (necesario)
type TokenPair struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int64
}

func (j *JWTManager) GenerateTokenPair(userID, email string, role enum.SystemRole) (*TokenPair, error) {
    // Access token: 15 minutos
    // Refresh token: 7 días
}

func (j *JWTManager) RefreshAccessToken(refreshToken string) (string, error) {
    // Validar refresh token
    // Generar nuevo access token
}
```

**Impacto si no se implementa:**
- Usuarios deben re-loguearse cada vez que expira access token
- UX degradada en app móvil

**Tiempo estimado:** 2-3 horas

**⚠️ Nota:** Verificar si ya existe en código (no confirmado en análisis)

---

## 📋 Resumen de Gaps por Prioridad

### P0 - Críticos (Bloquean desarrollo)

1. **Crear módulo `evaluation/`**
   - Requerido por: api-mobile, worker
   - Impacto: Sin esto NO se puede implementar sistema de evaluaciones
   - Tiempo: 4-6 horas

2. **Implementar DLQ en `messaging/rabbit/`**
   - Requerido por: worker
   - Impacto: Worker no puede manejar errores correctamente
   - Tiempo: 3-4 horas

3. **Fix dependencias en `auth/` y `middleware/gin/`**
   - Requerido por: Todos
   - Impacto: Tests no ejecutables
   - Tiempo: 10 minutos (go mod tidy)

4. **Aumentar coverage en `database/postgres/` de 2% a >80%**
   - Requerido por: api-mobile, api-admin
   - Impacto: Alto riesgo de bugs en producción
   - Tiempo: 4-6 horas

**Total tiempo P0:** 12-16 horas (~2 días)

---

### P1 - Importantes (Afectan calidad)

5. **Agregar tests a `logger/` (0% → >80%)**
   - Tiempo: 3-4 horas

6. **Agregar tests a `common/*` (0% → >80%)**
   - Submódulos: errors, types, validator
   - Tiempo: 6-8 horas

7. **Verificar e implementar refresh tokens en `auth/`**
   - Si no existe: Implementar
   - Tiempo: 2-3 horas

8. **Validar tests en `database/mongodb/`**
   - Tiempo: 2 horas

**Total tiempo P1:** 13-17 horas (~2 días)

---

### P2 - Mejoras (Nice to have)

9. **Aumentar coverage en `config/` (32.9% → >80%)**
   - Tiempo: 2-3 horas

10. **Aumentar coverage en `bootstrap/` (29.9% → >80%)**
    - Tiempo: 2-3 horas

**Total tiempo P2:** 4-6 horas (~1 día)

---

## 🎯 Plan de Acción Consolidado

### Sprint 1: Gaps P0 (Críticos)
**Duración:** 1 semana

- [ ] Crear módulo `evaluation/` v0.1.0
- [ ] Implementar DLQ en `messaging/rabbit/` v0.6.0
- [ ] Fix dependencias (go mod tidy) en auth, middleware/gin
- [ ] Aumentar coverage en `database/postgres/` a >80%

**Entregables:**
- evaluation/v0.1.0 publicado
- messaging/rabbit/v0.6.0 publicado
- database/postgres/v0.6.0 publicado con >80% coverage
- Todos los tests pasando

---

### Sprint 2: Gaps P1 (Importantes)
**Duración:** 1 semana

- [ ] Agregar tests a `logger/` (>80% coverage)
- [ ] Agregar tests a `common/*` (>80% coverage)
- [ ] Implementar/verificar refresh tokens en `auth/`
- [ ] Validar tests en `database/mongodb/`

**Entregables:**
- logger/v0.6.0 con tests
- common/v0.6.0 con tests
- auth/v0.6.0 con refresh tokens (si no existe)
- database/mongodb/v0.6.0 validado

---

### Sprint 3: Consolidación y Congelamiento
**Duración:** 3 días

- [ ] Aumentar coverage en config/ y bootstrap/ (P2)
- [ ] Ejecutar suite completa de tests
- [ ] Validar coverage global >85%
- [ ] Release coordinado: **todos los módulos a v0.7.0**
- [ ] Congelar versión

**Entregables:**
- Todos los módulos en v0.7.0
- Coverage global >85%
- Documentación actualizada
- go.mod.example para cada consumidor

---

## ✅ Criterios de Éxito

### Para considerar shared "LISTO"

- ✅ Módulo `evaluation/` existe y publicado
- ✅ DLQ implementado en `messaging/rabbit/`
- ✅ Coverage global >85%
- ✅ Todos los tests pasando (0 failing)
- ✅ 0 dependencias rotas
- ✅ Refresh tokens implementados (o confirmado que NO son necesarios)
- ✅ api-mobile puede compilar con shared v0.7.0
- ✅ api-admin puede compilar con shared v0.7.0
- ✅ worker puede compilar con shared v0.7.0

---

**Documento generado:** 15 de Noviembre, 2025  
**Basado en:** Documentación de api-mobile, api-admin, worker  
**Próximo documento:** `03-MODULOS_FALTANTES.md`
