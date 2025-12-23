# Mejoras Pendientes - EduGo Shared

> Este directorio contiene análisis de código que necesita mejoras, refactoring o eliminación.

## Índice de Mejoras

| Prioridad | Documento | Descripción |
|-----------|-----------|-------------|
| **ALTA** | [CODIGO_INCOMPLETO.md](./CODIGO_INCOMPLETO.md) | Funciones con TODOs y código sin implementar |
| **ALTA** | [TESTS_SKIPPED.md](./TESTS_SKIPPED.md) | Tests deshabilitados que necesitan arreglarse |
| **MEDIA** | [MALAS_PRACTICAS.md](./MALAS_PRACTICAS.md) | Código con malas prácticas a corregir |
| **MEDIA** | [REFACTORING.md](./REFACTORING.md) | Código que necesita refactorización |
| **BAJA** | [DEUDA_TECNICA.md](./DEUDA_TECNICA.md) | Deuda técnica general del proyecto |

---

## Resumen Ejecutivo

### Estado Actual del Código

| Categoría | Cantidad | Impacto |
|-----------|----------|---------|
| TODOs pendientes | 3 | Alto - Funcionalidad incompleta |
| Tests deshabilitados | 3 archivos | Medio - Coverage afectado |
| Código sin usar | 2 imports | Bajo - Limpieza |
| Malas prácticas | 4 | Medio - Mantenibilidad |

### Acciones Inmediatas Recomendadas

1. **Implementar `GetPresignedURL`** en `resource_implementations.go`
2. **Implementar `extractEnvAndVersion`** en `bootstrap.go`
3. **Arreglar tests de integración** en archivos `.skip`
4. **Manejar errores de Ack/Nack** en `consumer.go`

---

## Cómo Usar Esta Documentación

### Para Desarrolladores

1. Antes de empezar a trabajar, revisa este directorio
2. Busca tareas que puedas resolver como parte de tu trabajo
3. Al resolver un issue, actualiza el documento correspondiente
4. Crea un PR con la mejora y referencia el documento

### Para Tech Leads

1. Usa estos documentos para planificar sprints de deuda técnica
2. Prioriza según el impacto en el negocio
3. Asigna tareas pequeñas a desarrolladores junior
4. Revisa periódicamente el estado de las mejoras

---

## Proceso de Resolución

```
1. Identificar → 2. Documentar → 3. Priorizar → 4. Resolver → 5. Verificar → 6. Cerrar
```

### Template para Resolución

```markdown
## [RESUELTO] Título del Issue

**Fecha resolución:** YYYY-MM-DD
**PR:** #XXX
**Autor:** @username

### Solución Implementada
Descripción de la solución...

### Tests Agregados
- test_xxx.go: TestNuevoTest
```

---

## Contribuir

Si encuentras código que debería estar aquí:

1. Crea una sección en el documento correspondiente
2. Incluye:
   - Ubicación exacta del archivo y línea
   - Descripción del problema
   - Impacto
   - Solución sugerida
3. Asigna una prioridad
