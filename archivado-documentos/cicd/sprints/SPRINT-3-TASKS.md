# Sprint 3: Mejora Gradual de Coverage - edugo-shared

**Duración:** 2 días  
**Objetivo:** Incrementar coverage de módulos que están cerca de sus objetivos  
**Estado:** En Ejecución

---

## 📋 Resumen del Sprint

| Métrica | Objetivo |
|---------|----------|
| **Tareas Totales** | 5 |
| **Tiempo Estimado** | 7-9 horas |
| **Prioridad** | Media |
| **Módulos a mejorar** | 4 |
| **Commits Esperados** | 4-5 |

---

## 🎯 Objetivos del Sprint

1. **bootstrap:** 29.5% → 40% (objetivo: +10.5 puntos)
2. **database/postgres:** 58.8% → 60% (objetivo: +1.2 puntos)
3. **database/mongodb:** 54.5% → 55% (objetivo: +0.5 puntos)
4. **messaging/rabbit:** 14.4% → 30% (objetivo: +15.6 puntos)

**Meta global:** Incrementar umbrales de 4 módulos

---

## 📝 TAREAS DETALLADAS

---

### ✅ Tarea 1.1: Mejorar bootstrap 29.5% → 40%

**Prioridad:** 🟠 Alta  
**Estimación:** ⏱️ 2-3 horas  
**Coverage actual:** 29.5%  
**Coverage objetivo:** 40%  
**Gap:** +10.5 puntos

#### Análisis de Funciones Sin Cobertura

Funciones al 0% que necesitan tests:
- Factories: `CreateConnection`, `CreateRawConnection`, `Ping`, `Close`
- Health checks: `performHealthChecks` (0%)
- Cleanup: `registerPostgreSQLCleanup`, `registerMongoDBCleanup`, `registerRabbitMQCleanup` (0%)

#### Tests a Agregar

**Total estimado:** ~12 tests

1. **Tests de PostgreSQL Factory** (~4 tests)
2. **Tests de MongoDB Factory** (~3 tests)
3. **Tests de RabbitMQ Factory** (~3 tests)
4. **Tests de Health Checks** (~2 tests)

---

### ✅ Tarea 1.2: Mejorar database/postgres 58.8% → 60%

**Prioridad:** 🟡 Media  
**Estimación:** ⏱️ 1 hora  
**Coverage actual:** 58.8%  
**Coverage objetivo:** 60%  
**Gap:** +1.2 puntos

#### Tests a Agregar

**Total estimado:** ~3 tests

1. **Test de transacción anidada**
2. **Test de rollback en panic**
3. **Test de error handling en transacciones**

---

### ✅ Tarea 1.3: Mejorar database/mongodb 54.5% → 55%

**Prioridad:** 🟡 Media  
**Estimación:** ⏱️ 1 hora  
**Coverage actual:** 54.5%  
**Coverage objetivo:** 55%  
**Gap:** +0.5 puntos

#### Tests a Agregar

**Total estimado:** ~2 tests

1. **Test de operaciones con contexto**
2. **Test de error handling en operaciones**

---

### ✅ Tarea 2.1: Mejorar messaging/rabbit 14.4% → 30%

**Prioridad:** 🟠 Alta  
**Estimación:** ⏱️ 3-4 horas  
**Coverage actual:** 14.4%  
**Coverage objetivo:** 30%  
**Gap:** +15.6 puntos

#### Tests a Agregar

**Total estimado:** ~15 tests

1. **Tests de Connection** (~4 tests)
2. **Tests de Consumer básico** (~5 tests)
3. **Tests de Publisher básico** (~4 tests)
4. **Tests de error handling** (~2 tests)

---

## 📊 Métricas Objetivo del Sprint

| Módulo | Actual | Objetivo | Tests a agregar | Prioridad |
|--------|--------|----------|-----------------|-----------|
| bootstrap | 29.5% | 40% | ~12 | Alta |
| database/postgres | 58.8% | 60% | ~3 | Media |
| database/mongodb | 54.5% | 55% | ~2 | Media |
| messaging/rabbit | 14.4% | 30% | ~15 | Alta |

**Total tests a agregar:** ~32 tests

---

## ✅ Criterios de Éxito

1. ✅ bootstrap alcanza 40% coverage
2. ✅ database/postgres alcanza 60% coverage
3. ✅ database/mongodb alcanza 55% coverage
4. ✅ messaging/rabbit alcanza 30% coverage
5. ✅ Todos los tests pasan
6. ✅ Build exitoso
7. ✅ CI/CD pasa
8. ✅ Umbrales actualizados en `.coverage-thresholds.yml`

---

**Generado por:** Claude Code  
**Fecha:** 20 Nov 2025
