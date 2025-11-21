# Estrategia de Testing y Coverage

**Sprint:** SPRINT-2  
**Fecha:** 2025-11-20  
**Versión:** 1.0

---

## 🎯 Objetivos

1. Alcanzar coverage global del 65%
2. Resolver módulo común (error en tests)
3. Mejorar messaging/rabbit de 2.9% → 15% (crítico)
4. Ningún módulo crítico por debajo del umbral
5. Plan de mejora gradual para módulos bajos

---

## 📊 Estado Actual vs Objetivo

| Módulo | Actual | Umbral | Gap | Prioridad | Sprint |
|--------|--------|--------|-----|-----------|--------|
| **EXCELENTES (>80%)** |
| evaluation | 100% | 95% | ✅ +5% | Mantener | - |
| middleware/gin | 98.5% | 95% | ✅ +3.5% | Mantener | - |
| logger | 95.8% | 90% | ✅ +5.8% | Mantener | - |
| lifecycle | 91.8% | 85% | ✅ +6.8% | Mantener | - |
| auth | 85.0% | 80% | ✅ +5% | Mantener | - |
| config | 82.9% | 75% | ✅ +7.9% | Mantener | - |
| **BUENOS (60-80%)** |
| testing | 59.0% | 55% | ✅ +4% | Mantener | - |
| **ACEPTABLES (40-60%)** |
| database/postgres | 58.8% | 60% | 🟡 -1.2% | Media | Sprint 3 |
| database/mongodb | 54.5% | 55% | 🟡 -0.5% | Media | Sprint 3 |
| **BAJOS (20-40%)** |
| bootstrap | 29.5% | 40% | 🟠 -10.5% | Alta | Sprint 3 |
| **CRÍTICOS (<20%)** |
| messaging/rabbit | 2.9% | 15% | 🔴 -12.1% | **CRÍTICA** | **Sprint 2** |
| **CON ERRORES** |
| common | ERROR | 70% | ❌ N/A | Alta | **Sprint 2** |

---

## 🚦 Clasificación de Módulos

### ✅ Excelentes (>80%) - 6 módulos
**Acción:** Mantener y proteger con umbrales altos

### 🟢 Buenos (60-80%) - 1 módulo
**Acción:** Mantener

### 🟡 Aceptables (40-60%) - 2 módulos
**Acción:** Mejorar en Sprint 3 (fácil, solo 1-2% cada uno)

### 🟠 Bajos (20-40%) - 1 módulo
**Acción:** Plan de mejora en Sprint 3

### 🔴 Críticos (<20%) - 1 módulo
**Acción:** **PRIORIDAD MÁXIMA** - Sprint 2

### ❌ Con Errores - 1 módulo
**Acción:** Resolver en Sprint 2

---

## 📋 Plan de Acción Inmediato (SPRINT-2)

### 🔴 Prioridad 1: Resolver common (1-2h)
**Problema:** Tests fallan con error  
**Acciones:**
1. Investigar error en tests de common
2. Resolver problema de covdata tool
3. Ejecutar tests exitosamente
4. Medir coverage real
5. Establecer umbral en 70%

**Estimación:** 1-2 horas

---

### 🔴 Prioridad 2: messaging/rabbit 2.9% → 15% (3-4h)
**Problema:** Solo 2.9% de cobertura

**Funciones sin cobertura (0%):**
- DLQ configuration
- Consumer/Publisher
- Error handling

**Acciones:**
1. **Tests de configuración DLQ:**
   - `TestDefaultDLQConfig` ✅ (ya existe)
   - `TestCalculateBackoff` ✅ (ya existe)
   - Agregar tests para casos edge

2. **Tests de consumer básico:**
   - Configuración de consumer
   - Manejo de mensajes
   - Error handling básico

3. **Tests de publisher:**
   - Publicación básica
   - Headers y propiedades

**Estimación:** 3-4 horas para llegar a 15-20%

---

## 📋 Plan Futuro

### Sprint 3: Optimizar Medios y Bajos (4-5h)

1. **database/postgres:** 58.8% → 60% (1h)
   - 2-3 tests de transacciones
   - Edge cases

2. **database/mongodb:** 54.5% → 55% (1h)
   - Tests de operaciones CRUD
   - Manejo de conexiones

3. **bootstrap:** 29.5% → 40% (2-3h)
   - Tests de factories (CreateConnection, etc)
   - Tests de cleanup lifecycle
   - Tests de health checks

### Sprint 4: Incrementar Umbrales (3-4h)

1. **messaging/rabbit:** 15% → 30%
2. **bootstrap:** 40% → 60%
3. Review general

---

## 🎓 Guías Rápidas de Testing

### Para messaging/rabbit

```go
func TestConsumerConfig(t *testing.T) {
    config := &ConsumerConfig{
        AutoAck: false,
        DLQ: DLQConfig{
            Enabled: true,
            MaxRetries: 3,
        },
    }
    
    assert.False(t, config.AutoAck)
    assert.Equal(t, 3, config.DLQ.MaxRetries)
}
```

### Para bootstrap factories

```go
func TestPostgreSQLFactory_CreateConnection(t *testing.T) {
    container := testing.NewPostgresContainer(t)
    defer container.Close()
    
    factory := NewDefaultPostgreSQLFactory()
    db, err := factory.CreateConnection(container.Config())
    
    assert.NoError(t, err)
    assert.NotNil(t, db)
    assert.NoError(t, db.Ping())
}
```

---

## 📈 Métricas de Seguimiento

### Estado Actual
- Módulos cumpliendo umbral: 8/12 (66.7%)
- Módulos excelentes (>80%): 6/12 (50%)
- Módulos críticos (<20%): 1/12 (8.3%)
- Coverage global estimado: ~60%

### Objetivo Sprint 2
- Módulos cumpliendo umbral: 10/12 (83.3%)
- Módulos con errores: 0/12 (0%)
- Módulos críticos: 0/12 (0%)

### Objetivo a 3 Meses
- Módulos cumpliendo umbral: 12/12 (100%)
- Módulos excelentes (>80%): 8/12 (66.7%)
- Coverage global: >65%

---

**Mantenido por:** EduGo Team  
**Revisión:** Cada sprint  
**Próxima revisión:** Post-Sprint 2
