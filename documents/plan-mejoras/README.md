# Plan Maestro de Mejoras - EduGo Shared

> Plan estructurado y orquestado para implementar todas las correcciones y mejoras detectadas en el proyecto.

---

## Resumen Ejecutivo

Este plan organiza la implementación de mejoras en **5 fases progresivas**, cada una con pasos atómicos que pueden ejecutarse de forma independiente y con commits granulares.

### Principios del Plan

1. **Atomicidad**: Cada paso es un cambio pequeño y reversible
2. **No Regresión**: Cada paso debe dejar el proyecto en estado funcional
3. **Trazabilidad**: Cada paso tiene criterios de éxito medibles
4. **Commits Frecuentes**: Después de cada paso exitoso, se puede hacer commit
5. **Orden de Dependencias**: Las fases respetan dependencias entre cambios
6. **Documentación Limpia**: Sin historial de "antes/después", solo el estado actual

---

## Flujo de Trabajo Estandarizado por Fase

Cada fase DEBE seguir este flujo de trabajo obligatorio:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FLUJO DE TRABAJO POR FASE                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. CREAR RAMA DESDE DEV                                                    │
│     └── git checkout dev && git pull origin dev                             │
│     └── git checkout -b fase-X-descripcion                                  │
│                                                                             │
│  2. IMPLEMENTAR CAMBIOS                                                     │
│     └── Ejecutar pasos de la fase                                           │
│     └── Commits atómicos por cada paso completado                           │
│                                                                             │
│  3. ACTUALIZAR DOCUMENTACIÓN                                                │
│     └── Actualizar /documents con enfoque LIMPIO                            │
│     └── Solo describir CÓMO se hace, NO cómo se hacía antes                 │
│     └── Eliminar secciones de "migración" o "cambios"                       │
│                                                                             │
│  4. CREAR PULL REQUEST A DEV                                                │
│     └── git push origin fase-X-descripcion                                  │
│     └── Crear PR hacia rama dev                                             │
│                                                                             │
│  5. ESPERAR REVISIÓN DE GITHUB COPILOT                                      │
│     └── Esperar que termine la revisión automática                          │
│     └── DESCARTAR comentarios sobre traducción inglés/español               │
│     └── CORREGIR comentarios importantes                                    │
│     └── DOCUMENTAR lo que no es necesario corregir o es deuda futura        │
│                                                                             │
│  6. ESPERAR PIPELINES (Máximo 10 minutos)                                   │
│     └── Revisar estado cada 1 minuto                                        │
│     └── Si hay errores: CORREGIR (regla de 3 intentos)                      │
│     └── Errores propios o heredados: TODOS se corrigen                      │
│                                                                             │
│  7. MERGE A DEV                                                             │
│     └── Solo cuando pipelines pasen                                         │
│     └── Solo cuando revisión esté aprobada                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Regla de 3 Intentos para Errores

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         REGLA DE 3 INTENTOS                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Para cualquier error (propio o heredado):                                  │
│                                                                             │
│  INTENTO 1:                                                                 │
│    └── Analizar causa raíz del error                                        │
│    └── Implementar corrección                                               │
│    └── Verificar solución                                                   │
│                                                                             │
│  INTENTO 2 (si falla):                                                      │
│    └── Revisar análisis inicial                                             │
│    └── Considerar efectos secundarios                                       │
│    └── Implementar corrección alternativa                                   │
│                                                                             │
│  INTENTO 3 (si falla):                                                      │
│    └── Revisar impacto en otras partes del código                           │
│    └── Implementar solución más conservadora                                │
│                                                                             │
│  SI FALLA DESPUÉS DE 3 INTENTOS:                                            │
│    └── DETENER el proceso                                                   │
│    └── Documentar el problema                                               │
│    └── Informar al usuario con:                                             │
│        • Análisis del error                                                 │
│        • Los 3 intentos realizados                                          │
│        • Posibles soluciones                                                │
│        • Impacto de no resolver                                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Documentación Limpia (Sin Historial)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PRINCIPIO DE DOCUMENTACIÓN LIMPIA                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ❌ NO HACER:                                                               │
│     "Antes se usaba X, ahora se usa Y"                                      │
│     "La API cambió de A a B"                                                │
│     "En la versión anterior..."                                             │
│     "Migration guide: de v0.7 a v0.8"                                       │
│                                                                             │
│  ✅ HACER:                                                                  │
│     "Para conectar a MongoDB, usar ConnectionString(ctx)"                   │
│     "La API expone los siguientes métodos..."                               │
│     "Ejemplo de uso: ..."                                                   │
│                                                                             │
│  RAZÓN: La documentación debe reflejar el estado ACTUAL,                    │
│         no la historia de cambios. El historial está en Git.                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Visión General de Fases

```
┌────────────────────────────────────────────────────────────────────────────┐
│                        PLAN DE MEJORAS EDUGO-SHARED                         │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  FASE 1: CORRECCIONES CRÍTICAS (Prioridad ALTA)                            │
│  ══════════════════════════════════════════════                            │
│  Rama: fase-1-correcciones-criticas                                        │
│  Duración estimada: 2-3 días                                               │
│  - 1.1 Implementar GetPresignedURL                                         │
│  - 1.2 Corregir error handling en Exists()                                 │
│  - 1.3 Manejar errores de Ack/Nack                                         │
│  - 1.4 Implementar extractEnvAndVersion                                    │
│  │                                                                          │
│  ▼                                                                          │
│  FASE 2: RESTAURACIÓN DE TESTS (Prioridad ALTA)                            │
│  ═══════════════════════════════════════════════                           │
│  Rama: fase-2-restauracion-tests                                           │
│  Duración estimada: 2-3 días                                               │
│  - 2.1 Actualizar y restaurar tests MongoDB                                │
│  - 2.2 Actualizar y restaurar tests PostgreSQL                             │
│  - 2.3 Actualizar y restaurar tests RabbitMQ                               │
│  - 2.4 Verificar coverage >= 80%                                           │
│  │                                                                          │
│  ▼                                                                          │
│  FASE 3: REFACTORING ESTRUCTURAL (Prioridad MEDIA)                         │
│  ══════════════════════════════════════════════════                        │
│  Rama: fase-3-refactoring-estructural                                      │
│  Duración estimada: 3-4 días                                               │
│  - 3.1 Dividir bootstrap.go en archivos más pequeños                       │
│  - 3.2 Crear extractConfigField genérico                                   │
│  - 3.3 Unificar MessagePublisher                                           │
│  - 3.4 Corregir tipo de presignClient (interface{} → *s3.PresignClient)    │
│  - 3.5 Agregar control de goroutine en Consumer                            │
│  │                                                                          │
│  ▼                                                                          │
│  FASE 4: MEJORAS DE CALIDAD (Prioridad MEDIA)                              │
│  ═══════════════════════════════════════════════                           │
│  Rama: fase-4-mejoras-calidad                                              │
│  Duración estimada: 2 días                                                 │
│  - 4.1 Limpiar imports comentados                                          │
│  - 4.2 Documentar API de containers                                        │
│  - 4.3 Migrar a logger.Logger interface                                    │
│  - 4.4 Crear constantes para campos de log                                 │
│  │                                                                          │
│  ▼                                                                          │
│  FASE 5: DEUDA TÉCNICA (Prioridad BAJA)                                    │
│  ══════════════════════════════════════                                    │
│  Rama: fase-5-deuda-tecnica-[item]                                         │
│  Duración estimada: Ongoing (Boy Scout Rule)                               │
│  - 5.1 Agregar documentación GoDoc                                         │
│  - 5.2 Crear funciones Example                                             │
│  - 5.3 Agregar benchmarks                                                  │
│  - 5.4 Configurar .golangci.yml                                            │
│  - 5.5 Convertir tests a table-driven                                      │
│                                                                             │
└────────────────────────────────────────────────────────────────────────────┘
```

---

## Índice de Fases

| Fase | Documento | Prioridad | Rama | Estado |
|------|-----------|-----------|------|--------|
| 1 | [FASE-1_CORRECCIONES_CRITICAS.md](./FASE-1_CORRECCIONES_CRITICAS.md) | **ALTA** | `fase-1-correcciones-criticas` | ⏳ Pendiente |
| 2 | [FASE-2_RESTAURACION_TESTS.md](./FASE-2_RESTAURACION_TESTS.md) | **ALTA** | `fase-2-restauracion-tests` | ⏳ Pendiente |
| 3 | [FASE-3_REFACTORING.md](./FASE-3_REFACTORING.md) | **MEDIA** | `fase-3-refactoring-estructural` | ⏳ Pendiente |
| 4 | [FASE-4_MEJORAS_CALIDAD.md](./FASE-4_MEJORAS_CALIDAD.md) | **MEDIA** | `fase-4-mejoras-calidad` | ⏳ Pendiente |
| 5 | [FASE-5_DEUDA_TECNICA.md](./FASE-5_DEUDA_TECNICA.md) | **BAJA** | `fase-5-deuda-tecnica-*` | ⏳ Pendiente |

---

## Métricas de Éxito Global

### Antes de Iniciar

| Métrica | Estado Actual |
|---------|---------------|
| TODOs pendientes | 3 |
| Tests deshabilitados | 3 archivos (.skip) |
| Errores silenciados | 4+ instancias |
| Líneas en bootstrap.go | 623 |
| Duplicación de código | ~320 líneas |
| Coverage estimado bootstrap | ~65% |

### Objetivo Final

| Métrica | Objetivo |
|---------|----------|
| TODOs pendientes | 0 |
| Tests deshabilitados | 0 |
| Errores silenciados | 0 |
| Líneas en bootstrap.go | < 150 |
| Duplicación de código | < 50 líneas |
| Coverage bootstrap | >= 80% |

---

## Dependencias entre Fases

```
Fase 1 ──┐
         ├──► Fase 2 ──► Fase 3 ──► Fase 4 ──► Fase 5
         │
No tiene dependencias previas

Fase 2 depende de:
  - Fase 1.4 (extractEnvAndVersion) para tests de logger

Fase 3 depende de:
  - Fase 1 completa (APIs estables antes de refactorizar)
  - Fase 2 completa (tests funcionando para validar refactoring)

Fase 4 depende de:
  - Fase 3 completa (estructura estable)

Fase 5:
  - Puede ejecutarse en paralelo después de Fase 1
  - Ideal como "Boy Scout Rule" continuo
```

---

## Comandos del Flujo de Trabajo

### Inicio de Fase

```bash
# 1. Asegurarse de estar en dev actualizado
git checkout dev
git pull origin dev

# 2. Crear rama de la fase
git checkout -b fase-X-descripcion

# 3. Verificar estado inicial
make build
make test-all-modules
```

### Durante la Fase

```bash
# Commit atómico por paso completado
git add <archivos>
git commit -m "tipo(módulo): descripción del paso"

# Verificar que no hay regresiones
make test-all-modules
make lint
```

### Fin de Fase

```bash
# 1. Push de la rama
git push origin fase-X-descripcion

# 2. Crear PR en GitHub hacia dev
# 3. Esperar revisión de GitHub Copilot
# 4. Esperar pipelines (revisar cada minuto, máx 10 min)

# 5. Si hay errores en pipeline:
#    - Corregir (máx 3 intentos)
#    - Push de corrección
#    - Esperar pipelines nuevamente

# 6. Merge cuando todo esté verde
```

### Convención de Commits

```bash
# Correcciones de bugs
fix(bootstrap): implementar GetPresignedURL en StorageClient

# Nuevas funcionalidades
feat(auth): agregar refresh token rotation

# Refactoring
refactor(bootstrap): dividir bootstrap.go en múltiples archivos

# Tests
test(bootstrap): restaurar tests de integración MongoDB

# Documentación
docs(api): actualizar documentación de containers API

# Limpieza
chore(consumer): eliminar imports comentados
```

---

## Revisión de GitHub Copilot

### Comentarios a DESCARTAR

- Sugerencias de traducir código/comentarios de español a inglés
- Sugerencias de traducir de inglés a español
- Preferencias de estilo que no afectan funcionalidad

### Comentarios a CORREGIR

- Errores de lógica detectados
- Problemas de seguridad
- Malas prácticas de Go
- Errores de tipado
- Memory leaks potenciales
- Race conditions

### Comentarios a DOCUMENTAR (Deuda Futura)

- Mejoras que requieren cambios breaking
- Optimizaciones que pueden esperar
- Refactorizaciones grandes fuera del scope

---

## Manejo de Errores en Pipeline

### Análisis del Error

Antes de corregir, siempre analizar:

1. **¿Es por código nuevo de esta fase?**
   - Corregir inmediatamente

2. **¿Es por cambio de configuración?**
   - Revisar si el cambio es correcto
   - Corregir configuración o adaptar código

3. **¿Es por código existente (heredado)?**
   - Corregir igualmente (Boy Scout Rule)
   - Documentar que era error pre-existente

### Impacto de la Corrección

Antes de aplicar fix, considerar:

- ¿Esta corrección puede romper otra cosa?
- ¿Afecta comportamiento de producción?
- ¿Requiere actualización de tests?

---

## Priorización de Ejecución

### Modo Rápido (Solo Críticos)
Ejecutar únicamente:
- Fase 1 completa
- Fase 2 completa

**Tiempo estimado**: 4-6 días

### Modo Completo
Ejecutar todas las fases en orden.

**Tiempo estimado**: 10-15 días

### Modo Continuo (Boy Scout)
- Fase 1 y 2 como proyecto
- Fase 3, 4 y 5 como mejora continua

---

## Checklist General

Ver archivo [CHECKLIST.md](./CHECKLIST.md) para seguimiento detallado de todos los pasos.

---

## Historial de Cambios

| Fecha | Versión | Cambio |
|-------|---------|--------|
| 2024-12-22 | 1.0.0 | Creación inicial del plan |
| 2024-12-22 | 1.1.0 | Agregar flujo de trabajo con ramas, PR, Copilot review y pipelines |

---

## Autor

Generado automáticamente a partir del análisis de:
- `documents/README.md`
- `documents/mejoras/CODIGO_INCOMPLETO.md`
- `documents/mejoras/MALAS_PRACTICAS.md`
- `documents/mejoras/REFACTORING.md`
- `documents/mejoras/TESTS_SKIPPED.md`
- `documents/mejoras/DEUDA_TECNICA.md`
