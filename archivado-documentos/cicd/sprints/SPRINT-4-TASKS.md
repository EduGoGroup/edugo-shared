# Sprint 4: Completar Optimización de Coverage - edugo-shared

**Duración:** 1-2 días  
**Objetivo:** Completar mejoras de coverage pendientes del Sprint 3  
**Estado:** En Ejecución

---

## 📋 Resumen del Sprint

| Métrica | Objetivo |
|---------|----------|
| **Tareas Totales** | 2 |
| **Tiempo Estimado** | 5-7 horas |
| **Prioridad** | Media-Alta |
| **Módulos a mejorar** | 2 |

---

## 🎯 Objetivos del Sprint

1. **messaging/rabbit:** 14.4% → 30% (+15.6 puntos)
2. **bootstrap:** 35.7% → 40% (+4.3 puntos)

**Meta global:** Alcanzar objetivos diferidos del Sprint 3

---

## 📝 TAREAS DETALLADAS

---

### ✅ Tarea 1.1: Mejorar messaging/rabbit 14.4% → 30%

**Prioridad:** 🔴 Alta  
**Estimación:** ⏱️ 4-5 horas  
**Coverage actual:** 14.4%  
**Coverage objetivo:** 30%  
**Gap:** +15.6 puntos

#### Funciones Sin Cobertura (0%)

**Connection (0%):**
- Connect, GetChannel, GetConnection
- Close, IsClosed
- DeclareExchange, DeclareQueue, BindQueue
- SetPrefetchCount, HealthCheck

**Consumer (0%):**
- NewConsumer, Consume, Close
- ConsumeWithDLQ, setupDLQ, sendToDLQ

**Publisher (0%):**
- NewPublisher, Publish, PublishWithPriority, Close

#### Tests a Agregar

**Total estimado:** ~20-25 tests

1. **Tests de Connection** (~5 tests)
   - Conexión básica
   - GetChannel, GetConnection
   - DeclareQueue básico
   - IsClosed

2. **Tests de Consumer básico** (~5 tests)
   - NewConsumer con diferentes configs
   - UnmarshalMessage (ya tiene algunos)
   - Close

3. **Tests de Publisher básico** (~5 tests)
   - NewPublisher
   - Publish básico
   - PublishWithPriority
   - Close

4. **Tests de DLQ** (~3 tests)
   - setupDLQ configuración
   - sendToDLQ lógica
   - getRetryCount (ya al 100%)

5. **Tests sin containers** (~5 tests)
   - Mocks y configuración
   - Error handling
   - Edge cases

**Nota:** Usar containers de testing/containers para tests de integración

---

### ✅ Tarea 1.2: Mejorar bootstrap 35.7% → 40%

**Prioridad:** 🟡 Media  
**Estimación:** ⏱️ 1-2 horas  
**Coverage actual:** 35.7%  
**Coverage objetivo:** 40%  
**Gap:** +4.3 puntos

#### Tests a Agregar

**Total estimado:** ~5 tests

1. **Tests de Logger Factory** (~2 tests)
   - CreateLogger con diferentes configs
   - Error handling

2. **Tests adicionales de init functions** (~3 tests)
   - initLogger con configuración
   - Error paths en extract functions

---

## ✅ Criterios de Éxito

1. ✅ messaging/rabbit alcanza 30% coverage
2. ✅ bootstrap alcanza 40% coverage
3. ✅ Todos los tests pasan
4. ✅ Build exitoso
5. ✅ CI/CD pasa
6. ✅ 11/11 módulos cumplen umbrales actualizados

---

**Generado por:** Claude Code  
**Fecha:** 21 Nov 2025
