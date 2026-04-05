# Migration Strategy: ejecución segura por fases

## Objetivo

Evolucionar `edugo-shared` sin interrumpir consumidores activos, separando cambios sin ruptura de cambios potencialmente breaking.

## Fase 1: Quick wins (sin ruptura)

Acciones:

- Limpieza de tests triviales/duplicados.
- Consolidación de suites en módulos con muchos archivos de test.
- Split explícito de pruebas unitarias vs integración.

Criterio de salida:

- No hay cambios de API pública.
- Todos los módulos siguen compilando y testeando.

Impacto esperado en consumidores:

- Ninguno.

Playbook de ejecución:

1. Ejecutar baseline: `make test-all && make test-integration-all`.
2. Eliminar/combinar tests triviales por módulo (empezar por `common`, `logger`, `bootstrap`, `testing`).
3. Revalidar por módulo tocado: `go test -short ./...`.
4. Revalidar global: `make test-parallel`.
5. Documentar en `CHANGELOG.md` de cada módulo afectado.

## Fase 2: Reestructuración interna compatible

Acciones:

- `bootstrap`, `logger`, `screenconfig`: mover implementación a `internal/`.
- Mantener constructores, interfaces y tipos públicos en rutas actuales.
- Actualizar documentación de estructura por módulo.

Criterio de salida:

- Imports de consumidores siguen funcionando sin cambios.
- Sólo se modifica organización interna.

Impacto esperado en consumidores:

- Bajo o nulo, con validación puntual en `edugo-worker` para `bootstrap`.

Playbook de ejecución:

1. Crear estructura `internal/` en módulo objetivo (`bootstrap` primero).
2. Mover implementación interna sin cambiar firmas públicas.
3. Mantener factories públicas como facade/wrapper.
4. Ejecutar tests del módulo + consumidores más sensibles (primero `edugo-worker`).
5. Verificar que imports antiguos siguen compilando.

## Fase 3: Reducción de duplicación compartida

Acciones:

- Crear utilidades comunes (`common/health`, `common/retry`) y migración progresiva.
- Revisar wrappers de config duplicados en `bootstrap`.

Criterio de salida:

- Adopción incremental; no forzar migración masiva en una sola release.

Impacto esperado en consumidores:

- Compatible en general; potencialmente breaking si se retiran wrappers sin capa de compatibilidad.

Playbook de ejecución:

1. Crear paquetes comunes (`common/health`, `common/retry`) con API mínima.
2. Integrar primero de forma interna en módulos con más duplicación.
3. Mantener compatibilidad con wrappers locales hasta estabilizar.
4. Medir regresiones (latencia de reconexión, health checks) en integration tests.

## Fase 4: Cambios potencialmente breaking (solo con evidencia)

Acciones:

- Eliminar API pública marcada como dead code (por ejemplo funciones `All*()`), **solo** tras confirmar no uso externo adicional.

Requisitos previos:

- Escaneo de uso en repos consumidores objetivo.
- Nota de migración por módulo.
- Versionado semántico apropiado por módulo afectado.

Impacto esperado en consumidores:

- Bajo a medio, dependiendo de uso real no detectado.

Playbook de ejecución:

1. Escaneo cross-repo de uso (IAM/Admin/Mobile/Worker y otros repos internos).
2. Publicar propuesta de deprecación (release N) antes de remoción (release N+1).
3. Incluir snippet de migración por símbolo removido.
4. Ejecutar validación coordinada con equipos consumidores antes del tag final.

## Gobernanza de releases recomendada

1. Publicar changelog por módulo con sección de compatibilidad.
2. Marcar explícitamente “breaking changes” cuando aplique.
3. Mantener guía de migración corta por módulo crítico:
   - `auth`
   - `middleware/gin`
   - `repository`
   - `logger`
   - `bootstrap`

## Checklist mínimo antes de publicar

- [ ] Validación de compilación por módulo.
- [ ] Ejecución de suites unitarias e integración relevantes.
- [ ] Verificación de contratos públicos clave.
- [ ] Actualización de `README.md`, `docs/README.md` y `CHANGELOG.md` de módulos afectados.
- [ ] Comunicación a equipos consumidores sobre cambios y ventana de adopción.

## Módulo de mayor inversión recomendada

`bootstrap` es el módulo con mejor retorno de inversión por:

- mayor volumen de archivos y responsabilidades mezcladas,
- concentración de wiring de infraestructura crítica,
- dependencia directa del `edugo-worker`.

Implicación:

- Si `bootstrap` mejora su arquitectura interna sin romper su API, se reduce riesgo operativo del Worker y se simplifica evolución de infraestructura compartida.
