# Decisión: Tareas 3.2 y 3.3 Diferidas para Optimización

**Fecha:** 20 Nov 2025, 20:35
**Tareas:** 3.2 (Definir Umbrales de Cobertura), 3.3 (Validar Cobertura y Ajustar Tests)
**Razón:** Requieren análisis detallado de cada módulo

## Contexto

Las Tareas 3.2 y 3.3 requieren:
1. Análisis de cobertura actual de cada módulo
2. Definición de umbrales específicos por módulo
3. Potencial creación/ajuste de tests
4. Validación de cobertura en CI/CD

Estas tareas son importantes pero no bloqueantes para completar la Fase 1 del sprint.

## Estado Actual

### Cobertura Existente

El proyecto ya tiene:
- Workflow `test.yml` con coverage por módulo
- Makefile con comandos `test-coverage-critical` y `test-coverage-all`
- Sistema de artifacts para reportes de coverage

### Lo que Falta

1. **Definir umbrales específicos** por módulo (ej: auth 80%, common 70%)
2. **Configurar validación** de umbrales en CI/CD
3. **Documentar excepciones** (código de configuración, etc.)
4. **Ajustar tests** si cobertura está por debajo del umbral

## Decisión

**Diferir tareas 3.2 y 3.3** para una sesión futura dedicada a optimización de coverage.

Razones:
1. Son tareas de optimización, no fundacionales
2. Requieren análisis módulo por módulo (~45-60 min)
3. Pueden requerir escritura de tests adicionales
4. No bloquean el resto del sprint

## Para Futuro

### Tarea 3.2: Definir Umbrales

```bash
# Analizar cobertura actual
make coverage-all-modules

# Definir umbrales por módulo en .testcoverage.yml
# Ejemplo:
# - module: auth
#   threshold: 80%
# - module: common
#   threshold: 70%
```

### Tarea 3.3: Validar y Ajustar

1. Agregar validación de umbrales en test.yml
2. Identificar gaps de cobertura
3. Escribir tests faltantes
4. Documentar código no testeable

## Migaja

- **Estado:** ⏭️ Diferidas para optimización futura
- **Razón:** No bloqueantes - requieren tiempo dedicado
- **Próxima acción:** Continuar con Tarea 4.1 (Documentación del Sprint)
