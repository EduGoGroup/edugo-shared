# Plan de Trabajo: Mejora de Cobertura de Tests
## Proyecto: edugo-shared (Go 1.25)

**Fecha:** 2025-11-22
**Objetivo:** Alcanzar 80% de cobertura mínima en todos los módulos
**Estado Actual:** 4 módulos por debajo del 80%

---

## 📊 Resumen Ejecutivo

### Módulos que Requieren Atención (< 80%)

| Módulo | Cobertura Actual | Objetivo | Gap | Prioridad | Esfuerzo Estimado |
|--------|------------------|----------|-----|-----------|-------------------|
| **messaging/rabbit** | 30.0% | 80% | **+50%** | 🔴 CRÍTICA | 8-10 horas |
| **bootstrap** | 40.0% | 80% | **+40%** | 🔴 ALTA | 6-8 horas |
| **database/postgres** | 58.8% | 80% | **+21.2%** | 🟡 MEDIA | 3-4 horas |
| **testing** | 59.0% | 80% | **+21%** | 🟡 MEDIA | 3-4 horas |

**Total Esfuerzo Estimado:** 20-26 horas

### Módulos Cerca del Umbral (80-85%)

| Módulo | Cobertura Actual | Estado | Acción |
|--------|------------------|--------|--------|
| **database/mongodb** | 81.8% | ✅ Cumple | Mantener y agregar tests de edge cases |
| **config** | 82.9% | ✅ Cumple | Mantener cobertura actual |
| **auth** | 85.0% | ✅ Bueno | Objetivo: 90% |

---

## 🎯 Estrategia General

### Fase 1: Módulos Críticos (Prioridad Alta)
**Duración:** Sprint 5 (2 semanas)
- messaging/rabbit: 30% → 80%
- bootstrap: 40% → 80%

### Fase 2: Módulos Medios (Prioridad Media)
**Duración:** Sprint 6 (1 semana)
- database/postgres: 58.8% → 80%
- testing: 59% → 80%

### Fase 3: Consolidación (Prioridad Baja)
**Duración:** Sprint 7 (1 semana)
- Revisar módulos cerca del umbral
- Incrementar umbrales de módulos que cumplen
- Documentación y mejoras

---

## 📋 Plan Detallado por Módulo

---

## 1️⃣ messaging/rabbit (30% → 80%)

### 🔍 Análisis de Código Sin Cubrir

**Archivos Existentes:**
- ✅ `config.go` + `config_test.go` - Ya tiene tests
- ✅ `connection.go` + `connection_test.go` - Tests con containers (Sprint 4)
- ✅ `consumer.go` + `consumer_test.go` - Tests con containers (~20 tests)
- ✅ `publisher.go` + `publisher_test.go` - Tests con containers (~21 tests)
- ❌ `consumer_dlq.go` - **SIN COBERTURA COMPLETA**
- ❌ `dlq.go` - **SIN COBERTURA COMPLETA**
- ✅ `consumer_dlq_helpers_test.go` - Helpers parciales

### 🎯 Tests Necesarios

#### A. Tests de DLQ Consumer Integration (~10 tests, +15% cobertura)

```go
// consumer_dlq_test.go (NUEVO)

func TestConsumerDLQ_Integration(t *testing.T) {
    // Setup con RabbitMQ container
    // Tests:
    // 1. Crear consumer DLQ básico
    // 2. Consumir mensajes de DLQ
    // 3. Procesamiento exitoso de mensaje DLQ
    // 4. Manejo de errores en procesamiento DLQ
    // 5. Reintento de mensajes DLQ
    // 6. DLQ con múltiples consumidores
    // 7. Cancelación de consumer DLQ
    // 8. Cierre limpio de consumer DLQ
}

func TestConsumerDLQ_ErrorHandling(t *testing.T) {
    // Tests de edge cases:
    // 1. Cola DLQ no existe
    // 2. Conexión perdida durante consumo
}
```

**Esfuerzo:** 3-4 horas
**Cobertura Esperada:** +15%

#### B. Tests de DLQ Helpers y Utilities (~8 tests, +10% cobertura)

```go
// dlq_integration_test.go (NUEVO)

func TestDLQ_MoveToRetry(t *testing.T) {
    // Tests:
    // 1. Mover mensaje a cola de retry exitosamente
    // 2. Mover mensaje con delay
    // 3. Mover mensaje preservando headers
    // 4. Error al mover mensaje (cola destino no existe)
}

func TestDLQ_GetDeadLetterCount(t *testing.T) {
    // Tests:
    // 1. Contar mensajes en DLQ correctamente
    // 2. DLQ vacía retorna 0
    // 3. Error de conexión durante conteo
}

func TestDLQ_PurgeDeadLetterQueue(t *testing.T) {
    // Tests:
    // 1. Purgar DLQ con mensajes
    // 2. Purgar DLQ vacía (no error)
}
```

**Esfuerzo:** 2-3 horas
**Cobertura Esperada:** +10%

#### C. Tests de Connection Edge Cases (~6 tests, +8% cobertura)

```go
// connection_advanced_test.go (NUEVO)

func TestConnection_Reconnection(t *testing.T) {
    // Tests:
    // 1. Reconexión automática después de pérdida
    // 2. Reconexión con backoff exponencial
    // 3. Fallo de reconexión después de max intentos
}

func TestConnection_ChannelPool(t *testing.T) {
    // Tests:
    // 1. Crear múltiples canales concurrentemente
    // 2. Cerrar canales en orden
    // 3. Manejo de canal cerrado inesperadamente
}
```

**Esfuerzo:** 2-3 horas
**Cobertura Esperada:** +8%

#### D. Tests de Publisher y Consumer Advanced (~10 tests, +12% cobertura)

```go
// publisher_advanced_test.go (NUEVO)

func TestPublisher_RetryMechanism(t *testing.T) {
    // Tests:
    // 1. Retry después de fallo temporal
    // 2. Fallo permanente después de max retries
    // 3. Backoff entre retries
}

func TestPublisher_BatchPublishing(t *testing.T) {
    // Tests:
    // 1. Publicar batch de mensajes
    // 2. Partial failure en batch
    // 3. Rollback de batch en error
}

// consumer_advanced_test.go (NUEVO)

func TestConsumer_ConcurrentProcessing(t *testing.T) {
    // Tests:
    // 1. Procesar múltiples mensajes concurrentemente
    // 2. Respetar límite de concurrencia
    // 3. Graceful shutdown con mensajes en proceso
}

func TestConsumer_Acknowledgment(t *testing.T) {
    // Tests:
    // 1. ACK exitoso
    // 2. NACK y requeue
    // 3. NACK sin requeue (a DLQ)
}
```

**Esfuerzo:** 3-4 horas
**Cobertura Esperada:** +12%

### 📈 Resumen messaging/rabbit

| Categoría | Tests | Esfuerzo | Cobertura |
|-----------|-------|----------|-----------|
| DLQ Consumer Integration | 10 | 3-4h | +15% |
| DLQ Helpers | 8 | 2-3h | +10% |
| Connection Advanced | 6 | 2-3h | +8% |
| Publisher/Consumer Advanced | 10 | 3-4h | +12% |
| **TOTAL** | **34** | **10-14h** | **+45%** |

**Cobertura Final Esperada:** 75-80%

---

## 2️⃣ bootstrap (40% → 80%)

### 🔍 Análisis de Código Sin Cubrir

**Archivos con Cobertura Parcial:**
- ✅ `factory_logger.go` + tests - Ya cubierto (~100%)
- ❌ `factory_postgresql.go` - **Falta CreateConnection, CreateRawConnection, Ping, Close**
- ❌ `factory_mongodb.go` - **Falta CreateConnection, Ping, Close**
- ❌ `factory_rabbitmq.go` - **Falta CreateConnection, CreateChannel, DeclareQueue, Close**
- ❌ `factory_s3.go` - **COMPLETAMENTE SIN TESTS DE INTEGRACIÓN**
- ❌ `bootstrap.go` - **Faltan tests de integración completa**
- ❌ `resource_implementations.go` - **Sin tests**
- ❌ `options.go` - **Cobertura parcial**

### 🎯 Tests Necesarios

#### A. Tests de PostgreSQL Factory con Containers (~8 tests, +10% cobertura)

```go
// factory_postgresql_integration_test.go (NUEVO)

func TestPostgreSQLFactory_CreateConnection_Integration(t *testing.T) {
    // Setup con PostgreSQL container
    // Tests:
    // 1. CreateConnection exitosa con config válida
    // 2. CreateConnection falla con config inválida
    // 3. CreateConnection con SSL mode
    // 4. Ping exitoso
    // 5. Ping falla con DB desconectada
    // 6. Close exitoso
    // 7. Close falla si ya está cerrada
}

func TestPostgreSQLFactory_CreateRawConnection_Integration(t *testing.T) {
    // Tests:
    // 1. CreateRawConnection exitosa
    // 2. Configuración de connection pool
    // 3. Verificación de ping
}

func TestPostgreSQLFactory_ConnectionPool(t *testing.T) {
    // Tests:
    // 1. MaxOpenConns configurado correctamente
    // 2. MaxIdleConns configurado correctamente
    // 3. ConnMaxLifetime configurado
}
```

**Esfuerzo:** 2-3 horas
**Cobertura Esperada:** +10%

#### B. Tests de MongoDB Factory con Containers (~8 tests, +10% cobertura)

```go
// factory_mongodb_integration_test.go (NUEVO)

func TestMongoDBFactory_CreateConnection_Integration(t *testing.T) {
    // Setup con MongoDB container
    // Tests:
    // 1. CreateConnection exitosa
    // 2. CreateConnection falla con URI inválida
    // 3. Connection timeout funciona
    // 4. Ping exitoso
    // 5. Ping falla con timeout
    // 6. GetDatabase retorna database correcta
    // 7. Close exitoso
    // 8. Close con timeout
}

func TestMongoDBFactory_ConnectionPoolSettings(t *testing.T) {
    // Tests:
    // 1. MaxPoolSize configurado
    // 2. MinPoolSize configurado
    // 3. Connection idle time configurado
}
```

**Esfuerzo:** 2-3 horas
**Cobertura Esperada:** +10%

#### C. Tests de RabbitMQ Factory con Containers (~8 tests, +10% cobertura)

```go
// factory_rabbitmq_integration_test.go (NUEVO)

func TestRabbitMQFactory_CreateConnection_Integration(t *testing.T) {
    // Setup con RabbitMQ container
    // Tests:
    // 1. CreateConnection exitosa
    // 2. CreateConnection con timeout
    // 3. CreateConnection falla con URL inválida
    // 4. CreateConnection cancelación por contexto
}

func TestRabbitMQFactory_CreateChannel_Integration(t *testing.T) {
    // Tests:
    // 1. CreateChannel exitoso
    // 2. QoS configurado correctamente
    // 3. Error al configurar QoS
}

func TestRabbitMQFactory_DeclareQueue_Integration(t *testing.T) {
    // Tests:
    // 1. DeclareQueue exitoso
    // 2. Cola con configuración correcta (TTL, priority, etc)
    // 3. Error al declarar cola duplicada con config diferente
}

func TestRabbitMQFactory_Close_Integration(t *testing.T) {
    // Tests:
    // 1. Close exitoso de canal y conexión
    // 2. Close con canal nil (no error)
    // 3. Close con conexión ya cerrada
}
```

**Esfuerzo:** 2-3 horas
**Cobertura Esperada:** +10%

#### D. Tests de S3 Factory (Mocked) (~6 tests, +8% cobertura)

```go
// factory_s3_test.go (NUEVO con mocks)

func TestS3Factory_CreateClient_Mock(t *testing.T) {
    // Usar localstack o mocks
    // Tests:
    // 1. CreateClient exitoso con credenciales válidas
    // 2. CreateClient falla con credenciales inválidas
    // 3. ValidateBucket exitoso
    // 4. ValidateBucket falla si bucket no existe
    // 5. CreatePresignClient retorna cliente correcto
    // 6. Region configurada correctamente
}
```

**Esfuerzo:** 2 horas
**Cobertura Esperada:** +8%

#### E. Tests de Bootstrap Integration (~6 tests, +7% cobertura)

```go
// bootstrap_integration_test.go (NUEVO)

func TestBootstrap_InitializeAllResources(t *testing.T) {
    // Setup con containers
    // Tests:
    // 1. Inicializar logger solamente
    // 2. Inicializar logger + PostgreSQL
    // 3. Inicializar logger + PostgreSQL + MongoDB
    // 4. Inicializar todos los recursos
    // 5. Cleanup exitoso de todos los recursos
    // 6. Error en inicialización (cleanup parcial)
}
```

**Esfuerzo:** 2 horas
**Cobertura Esperada:** +7%

#### F. Tests de Options y Resource Implementations (~5 tests, +5% cobertura)

```go
// options_test.go (AMPLIAR)

func TestBootstrapOptions_AllOptions(t *testing.T) {
    // Tests:
    // 1. WithLogger option
    // 2. WithPostgreSQL option
    // 3. WithMongoDB option
    // 4. WithRabbitMQ option
    // 5. WithS3 option
    // 6. Múltiples options combinadas
}

// resource_implementations_test.go (NUEVO)

func TestResourceImplementations(t *testing.T) {
    // Tests de cada resource implementation
    // Verificar que implementan correctamente las interfaces
}
```

**Esfuerzo:** 1-2 horas
**Cobertura Esperada:** +5%

### 📈 Resumen bootstrap

| Categoría | Tests | Esfuerzo | Cobertura |
|-----------|-------|----------|-----------|
| PostgreSQL Factory | 8 | 2-3h | +10% |
| MongoDB Factory | 8 | 2-3h | +10% |
| RabbitMQ Factory | 8 | 2-3h | +10% |
| S3 Factory | 6 | 2h | +8% |
| Bootstrap Integration | 6 | 2h | +7% |
| Options & Resources | 5 | 1-2h | +5% |
| **TOTAL** | **41** | **11-15h** | **+50%** |

**Cobertura Final Esperada:** 90%

---

## 3️⃣ database/postgres (58.8% → 80%)

### 🔍 Análisis de Código Sin Cubrir

**Archivos:**
- ✅ `config.go` + `config_test.go` - Ya cubierto
- ✅ `connection.go` + `connection_test.go` - Ya cubierto
- ⚠️ `transaction.go` + `transaction_test.go` - **Falta cobertura de edge cases**

### 🎯 Tests Necesarios

#### A. Tests Adicionales de Transacciones (~8 tests, +15% cobertura)

```go
// transaction_test.go (AMPLIAR)

func TestWithTransaction_AdvancedScenarios(t *testing.T) {
    // Tests adicionales:
    // 1. Múltiples transacciones concurrentes
    // 2. Transacción con múltiples tablas
    // 3. Transacción con foreign key constraints
    // 4. Rollback con error de commit
    // 5. Nested transaction simulation (savepoints)
}

func TestWithTransactionIsolation_AllLevels(t *testing.T) {
    // Tests de todos los niveles de aislamiento:
    // 1. LevelDefault
    // 2. LevelReadUncommitted
    // 3. LevelRepeatableRead
    // 4. Verificar comportamiento de cada nivel
}

func TestTransaction_EdgeCases(t *testing.T) {
    // Tests de edge cases:
    // 1. Transacción con DB nil
    // 2. Transacción con contexto cancelado
    // 3. Timeout durante transacción
    // 4. Error al comenzar transacción (DB cerrada)
    // 5. Error al hacer commit (conflicto)
}
```

**Esfuerzo:** 2-3 horas
**Cobertura Esperada:** +15%

#### B. Tests de Connection Edge Cases (~4 tests, +6% cobertura)

```go
// connection_test.go (AMPLIAR)

func TestConnection_AdvancedScenarios(t *testing.T) {
    // Tests adicionales:
    // 1. Reconexión después de pérdida
    // 2. Connection pool exhausted
    // 3. Idle connections cleanup
    // 4. Connection con statement timeout
}
```

**Esfuerzo:** 1 hora
**Cobertura Esperada:** +6%

### 📈 Resumen database/postgres

| Categoría | Tests | Esfuerzo | Cobertura |
|-----------|-------|----------|-----------|
| Transaction Advanced | 8 | 2-3h | +15% |
| Connection Edge Cases | 4 | 1h | +6% |
| **TOTAL** | **12** | **3-4h** | **+21%** |

**Cobertura Final Esperada:** 80%

---

## 4️⃣ testing (59% → 80%)

### 🔍 Análisis de Código Sin Cubrir

**Archivos:**
- ✅ `containers/manager.go` + `manager_test.go` - Ya cubierto
- ✅ `containers/postgres.go` + `postgres_test.go` - Ya cubierto
- ✅ `containers/mongodb.go` + `mongodb_test.go` - Ya cubierto
- ✅ `containers/rabbitmq.go` + `rabbitmq_test.go` - Ya cubierto
- ⚠️ `containers/helpers.go` - **Falta cobertura de edge cases**
- ⚠️ `containers/options.go` - **Sin tests completos**

### 🎯 Tests Necesarios

#### A. Tests de Helpers (~8 tests, +10% cobertura)

```go
// helpers_test.go (NUEVO)

func TestExecSQLFile(t *testing.T) {
    // Tests:
    // 1. Ejecutar SQL file exitosamente
    // 2. Error al leer archivo (no existe)
    // 3. Error al ejecutar SQL (syntax error)
    // 4. SQL file con múltiples statements
}

func TestWaitForHealthy(t *testing.T) {
    // Tests:
    // 1. Health check exitoso inmediatamente
    // 2. Health check exitoso después de retries
    // 3. Timeout esperando health check
    // 4. Contexto cancelado durante espera
}

func TestRetryOperation(t *testing.T) {
    // Tests:
    // 1. Operación exitosa al primer intento
    // 2. Operación exitosa después de retries
    // 3. Operación falla después de max retries
    // 4. Delay entre retries funciona
}
```

**Esfuerzo:** 2 horas
**Cobertura Esperada:** +10%

#### B. Tests de Options y Builder Pattern (~6 tests, +7% cobertura)

```go
// options_test.go (NUEVO)

func TestConfig_Builder(t *testing.T) {
    // Tests:
    // 1. NewConfig con defaults
    // 2. WithPostgreSQL configuración
    // 3. WithMongoDB configuración
    // 4. WithRabbitMQ configuración
    // 5. Múltiples with methods encadenados
    // 6. Build retorna configuración correcta
}
```

**Esfuerzo:** 1 hora
**Cobertura Esperada:** +7%

#### C. Tests de Manager Edge Cases (~4 tests, +4% cobertura)

```go
// manager_test.go (AMPLIAR)

func TestManager_EdgeCases(t *testing.T) {
    // Tests adicionales:
    // 1. Cleanup parcial en error de setup
    // 2. Cleanup con containers nil
    // 3. Clean methods con containers no habilitados
    // 4. GetManager llamado múltiples veces (singleton)
}
```

**Esfuerzo:** 1 hora
**Cobertura Esperada:** +4%

### 📈 Resumen testing

| Categoría | Tests | Esfuerzo | Cobertura |
|-----------|-------|----------|-----------|
| Helpers | 8 | 2h | +10% |
| Options Builder | 6 | 1h | +7% |
| Manager Edge Cases | 4 | 1h | +4% |
| **TOTAL** | **18** | **4h** | **+21%** |

**Cobertura Final Esperada:** 80%

---

## 📅 Cronograma de Ejecución

### Sprint 5 (Semanas 1-2): Módulos Críticos

#### Semana 1: messaging/rabbit
- **Días 1-2:** DLQ Consumer Integration (3-4h)
- **Días 3:** DLQ Helpers (2-3h)
- **Día 4:** Connection Advanced (2-3h)
- **Día 5:** Publisher/Consumer Advanced (3-4h)
- **Total:** 10-14 horas

#### Semana 2: bootstrap
- **Días 1-2:** PostgreSQL + MongoDB Factories (4-6h)
- **Día 3:** RabbitMQ Factory (2-3h)
- **Día 4:** S3 Factory (2h)
- **Día 5:** Bootstrap Integration + Options (3-4h)
- **Total:** 11-15 horas

### Sprint 6 (Semana 3): Módulos Medios

#### Semana 3: database/postgres + testing
- **Días 1-2:** database/postgres (3-4h)
- **Días 3-4:** testing module (4h)
- **Día 5:** Review y ajustes
- **Total:** 7-8 horas

---

## 🎬 Checklist de Implementación

### Preparación
- [ ] Revisar que todos los containers estén funcionando
- [ ] Verificar versión de Go (1.25)
- [ ] Actualizar dependencias si es necesario
- [ ] Crear branch para cada módulo

### Por Cada Módulo
- [ ] Crear archivos de test nuevos
- [ ] Implementar tests según plan
- [ ] Ejecutar tests localmente
- [ ] Verificar cobertura con `go test -coverprofile`
- [ ] Revisar que no haya race conditions (`go test -race`)
- [ ] Commit y push de cambios
- [ ] Actualizar `.coverage-thresholds.yml` con nuevos números

### Validación Final
- [ ] Ejecutar `make coverage-all-modules`
- [ ] Validar que todos los módulos cumplan el 80%
- [ ] Ejecutar CI completo
- [ ] Revisar reportes de cobertura
- [ ] Actualizar documentación

---

## 📝 Notas Importantes

### Consideraciones Técnicas

1. **Testcontainers:** Todos los tests de integración deben usar testcontainers para PostgreSQL, MongoDB y RabbitMQ
2. **Cleanup:** Asegurar cleanup adecuado en `t.Cleanup()` o `defer`
3. **Timeouts:** Configurar timeouts apropiados para evitar tests colgados
4. **Aislamiento:** Cada test debe ser independiente y poder ejecutarse en paralelo
5. **Naming:** Seguir convención `TestFunctionName_Scenario` para nombres de tests

### Herramientas

```bash
# Ejecutar tests de un módulo con cobertura
cd <módulo>
go test -v -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out

# Ver cobertura en HTML
go tool cover -html=coverage.out -o coverage.html

# Ejecutar tests con race detection
go test -v -race ./...

# Ver cobertura de todos los módulos
make coverage-all-modules
```

### Mejores Prácticas

1. **Arrange-Act-Assert:** Estructura clara de tests
2. **Table-Driven Tests:** Para múltiples casos similares
3. **Subtests:** Usar `t.Run()` para organizar tests relacionados
4. **Mocks vs Real:** Preferir containers reales para integración, mocks para unit tests
5. **Error Messages:** Mensajes de error descriptivos con `require.NoError()` y `assert.Equal()`

---

## 🎯 Métricas de Éxito

### Objetivos Cuantitativos
- ✅ messaging/rabbit: ≥ 80% (actualmente 30%)
- ✅ bootstrap: ≥ 80% (actualmente 40%)
- ✅ database/postgres: ≥ 80% (actualmente 58.8%)
- ✅ testing: ≥ 80% (actualmente 59%)

### Objetivos Cualitativos
- ✅ Todos los tests pasan en CI
- ✅ No race conditions detectadas
- ✅ Tests ejecutan en < 5 minutos
- ✅ Código cubierto incluye edge cases críticos
- ✅ Documentación actualizada

---

## 🚀 Próximos Pasos

Después de alcanzar 80% en todos los módulos:

### Sprint 7: Optimización
1. Revisar módulos con >90% para alcanzar 95-100%
2. Incrementar umbrales en `.coverage-thresholds.yml`
3. Agregar tests de performance/benchmark
4. Revisar y mejorar tests existentes
5. Documentar casos de uso complejos

### Mejora Continua
1. Configurar pre-commit hooks para validar cobertura
2. Agregar coverage gates en CI/CD
3. Monitorear cobertura en cada PR
4. Revisar tests deprecated o redundantes

---

## 📚 Referencias

- **Documentación Go Testing:** https://golang.org/pkg/testing/
- **Testcontainers Go:** https://golang.testcontainers.org/
- **Best Practices:** https://github.com/golang/go/wiki/TestComments
- **Coverage Tools:** https://go.dev/blog/cover

---

## ✅ Aprobación

- [ ] Plan revisado por: __________________
- [ ] Fecha de inicio: __________________
- [ ] Fecha estimada de finalización: __________________
- [ ] Responsable: __________________

---

*Generado automáticamente el 2025-11-22*
