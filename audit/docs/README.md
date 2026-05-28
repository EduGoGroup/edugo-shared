# Audit — Documentación técnica

Contrato base para construir y despachar eventos de auditoría sin acoplar el almacenamiento.

## Propósito

Proporcionar una interfaz agnóstica que desacople a los consumidores del backend de auditoría, permitiendo múltiples implementaciones (PostgreSQL, Kafka, archivos, etc.).

## Flujo de operación

```
┌─────────────┐
│ AuditEvent  │
└──────┬──────┘
       │
       ├─→ WithSeverity, WithCategory, etc. (opcionales)
       │
       └─→ AuditLogger.Log(ctx, event)
           ↓
       Implementación concreta (postgres, noop, etc.)
```

## Componentes principales

### audit_logger.go — Contrato público

**AuditEvent**
Estructura que centraliza datos auditables:
- Actor: `ActorID`, `ActorEmail`, `ActorRole`
- HTTP: `RequestMethod`, `RequestPath`, `RequestID`, `StatusCode`
- Acción: `Action`, `ResourceType`, `ResourceID`
- Contexto: `ServiceName`, `Severity`, `Category`
- Opcionales: `ActorIP`, `ActorUserAgent`, `SchoolID`, `UnitID`, `PermissionUsed`, `Changes`, `Metadata`, `ErrorMessage`

**AuditLogger**
Interfaz mínima:
```go
type AuditLogger interface {
    Log(ctx context.Context, event AuditEvent) error
}
```

**AuditOption**
Funciones que modifican un `AuditEvent` de forma declarativa:
- `WithChanges(before, after)` — registra cambios de datos
- `WithSeverity(level)` — establece severidad (Info, Warning, Critical)
- `WithCategory(category)` — establece categoría (Auth, Data, Config, Admin)
- `WithMetadata(key, value)` — agrega metadatos adicionales
- `WithPermission(permission)` — registra permiso usado
- `WithError(err)` — registra mensaje de error

### noop_logger.go — Implementación de referencia

**NoopAuditLogger**
Implementación que descarta eventos sin hacer nada. Útil para tests y entornos de desarrollo.

```go
logger := NewNoopAuditLogger()
logger.Log(ctx, event) // Siempre retorna nil
```

## Constantes

Severidades:
- `SeverityInfo` = `"info"`
- `SeverityWarning` = `"warning"`
- `SeverityCritical` = `"critical"`

Categorías:
- `CategoryAuth` = `"auth"`
- `CategoryData` = `"data"`
- `CategoryConfig` = `"config"`
- `CategoryAdmin` = `"admin"`

## Testing

El módulo incluye 12 tests que validan:
- Funcionamiento de `NoopAuditLogger`
- Todas las funciones `WithX`
- Inicialización de mapas en opciones
- Valores de constantes
- Verificación de cumplimiento de interfaz

Ejecutar:
```bash
make test          # Tests locales
make check         # Tests + linting + format
```

## Integración

Para usar este módulo:

1. **Crear instancia de un AuditLogger concreto**
   ```go
   logger := postgres.NewPostgresAuditLogger(db, "my-service")
   ```

2. **Construir y enviar un evento**
   ```go
   event := audit.AuditEvent{
       ActorID:      userID,
       ActorEmail:   email,
       ActorRole:    role,
       ServiceName:  "my-service",
       Action:       "create",
       ResourceType: "user",
       ResourceID:   newUserID,
   }
   logger.Log(ctx, event)
   ```

3. **Enriquecer opcionalmente**
   ```go
   event.Severity = audit.SeverityWarning
   event.Category = audit.CategoryAdmin
   event.Metadata = map[string]any{"ip": "10.0.0.1"}
   ```

## Notas de diseño

- El módulo no contiene persistencia; es solo contrato.
- La estrategia de adaptador permite implementaciones múltiples sin afectar consumidores.
- Los `AuditOption` proporcionan forma declarativa de enriquecimiento.
- No tiene dependencias internas del repositorio, solo stdlib.
- Cada adaptador concreto (postgres, kafka, etc.) vive en su propio módulo con su propio `go.mod`.
