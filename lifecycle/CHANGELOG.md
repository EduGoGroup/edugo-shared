# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/lifecycle` se registran aquí.

## [0.100.0] - 2026-04-01

### Added

- **Manager**: Orquestador de ciclo de vida con startup secuencial y cleanup LIFO.
- **Resource**: Estructura que define nombre, startup (opcional) y cleanup.
- **Register method**: Registrar recurso con ambas fases (startup y cleanup).
- **RegisterSimple method**: Registrar recurso solo con cleanup (sin startup).
- **Startup method**: Ejecutar startup de todos los recursos en orden de registro (aborta si falla).
- **Cleanup method**: Ejecutar cleanup de todos los recursos en orden inverso (continúa incluso si falla).
- **Count method**: Retornar cantidad de recursos registrados.
- **Clear method**: Limpiar lista sin ejecutar cleanup (para testing).
- Thread-safety: Protección con mutex para operaciones concurrentes.
- Logger opcional: Trazabilidad sin obligatoriedad (logger.Logger).
- Métricas de duración: Calcula y reporta tiempos de startup y cleanup.
- Suite completa de tests unitarios con race detector.
- Documentación técnica detallada en docs/README.md con flujos comunes y arquitectura.
- Makefile con targets: build, test, test-race, check, lint, fmt, vet, tidy, deps, release.

### Design Notes

- Módulo pequeño y estable orientado a coordinación in-process.
- LIFO (Last In, First Out) es el contrato principal: inversión de orden de registro garantiza dependencias limpias correctamente.
- Tolerancia en cleanup: continúa limpiando incluso si fallan recursos para evitar dejar recursos abiertos.
- Logger opcional: adición de trazabilidad sin acoplamiento fuerte.
- Sin framework específico: funciona con cualquier tipo de recurso que implemente las interfaces de startup/cleanup.
