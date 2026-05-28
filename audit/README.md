# Audit

Contrato base para construir y despachar eventos de auditoría sin acoplar el almacenamiento.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/audit
```

## Quick Start

```go
// Crear un evento de auditoría
event := audit.AuditEvent{
    ActorID:      "user-123",
    ActorEmail:   "user@example.com",
    ActorRole:    "admin",
    ServiceName:  "my-service",
    Action:       "create",
    ResourceType: "document",
    ResourceID:   "doc-456",
}

// Enriquecer con opciones
event.Severity = audit.SeverityInfo
event.Category = audit.CategoryData

// O usar helpers para mayor claridad
logger.Log(ctx, event)
```

## API Pública

### AuditEvent
Estructura que centraliza los datos auditables: actor (ID, email, rol), acción, recurso, request, metadatos y cambios.

### AuditLogger
Interfaz mínima que define el contrato `Log(ctx context.Context, event AuditEvent) error`. Cualquier backend (PostgreSQL, Kafka, etc.) puede implementarlo.

### AuditOption
Funciones declarativas para enriquecer eventos:
- `WithChanges(before, after)` — registra cambios de datos
- `WithSeverity(level)` — establece nivel de severidad
- `WithCategory(category)` — establece categoría del evento
- `WithMetadata(key, value)` — agrega metadatos adicionales
- `WithPermission(permission)` — registra permiso usado
- `WithError(err)` — registra error asociado

### NoopAuditLogger
Implementación inerte que descarta eventos, útil para tests y entornos de desarrollo.

### Constantes
Severidades: `SeverityInfo`, `SeverityWarning`, `SeverityCritical`
Categorías: `CategoryAuth`, `CategoryData`, `CategoryConfig`, `CategoryAdmin`

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build    # Compilar módulo
make test     # Ejecutar tests
make check    # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- Este módulo define solo el contrato y una implementación noop.
- La persistencia en PostgreSQL y la extracción desde Gin viven en adaptadores separados.
- El módulo es agnóstico al backend, permitiendo múltiples implementaciones de `AuditLogger`.
