# Documentación técnica - Audit PostgreSQL

## Descripción general

Adaptador de auditoría que implementa `audit.AuditLogger` persistiendo eventos en la tabla `audit.events` mediante GORM. Normaliza datos, aplica defaults y convierte el contrato público `AuditEvent` al modelo GORM interno `AuditEventDB`.

## Flujo de operación

```
┌─────────────┐
│ AuditEvent  │
└──────┬──────┘
       │
       ├─→ Defaults (Severity, Category, ServiceName)
       │
       ├─→ Actor defaults (si faltan)
       │
       ├─→ ToDBModel (conversión a AuditEventDB)
       │
       └─→ GORM Create en audit.events
```

## Componentes principales

### types.go - API Pública

**PostgresAuditLogger**
- Implementación de `audit.AuditLogger`
- Campos: `db *gorm.DB`, `serviceName string`

**NewPostgresAuditLogger(db *gorm.DB, serviceName string)**
- Constructor que vincula la instancia a una conexión GORM y nombre de servicio.

**Log(ctx context.Context, event audit.AuditEvent) error**
- Aplica defaults de severity, category, serviceName y actor.
- Convierte a modelo GORM y persiste con contexto.

**LogFromGin(c *gin.Context, action, resourceType, resourceID string, opts ...audit.AuditOption) error**
- Extrae automáticamente: método HTTP, path, IP, user-agent, headers.
- Busca `user_id`, `email` y `role` del contexto Gin.
- Delega en `Log` con el evento construido.

### internal/models.go - Modelo GORM

**AuditEventDB**
- Estructura con tags GORM mapeada a tabla `audit.events`.
- Campos opcionales como punteros (`ActorIP`, `ResourceID`, etc.).
- Serialización JSON automática para `Changes` y `Metadata`.
- Timestamps automáticos con `autoCreateTime`.

**TableName()**
- Retorna `"audit.events"` para mapeo de tabla.

### internal/converter.go - Conversión

**ToDBModel(event audit.AuditEvent) AuditEventDB**
- Convierte `AuditEvent` (público) a `AuditEventDB` (GORM).
- Maneja conversión de campos obligatorios vs opcionales.
- Preserva valores no vacíos, ignora ceros.

### internal/defaults.go - Constantes

- `DefaultActorID`: UUID nil para actores anónimos.
- `DefaultActorEmail`: `"system"` cuando no se proporciona.
- `DefaultActorRole`: `"unknown"` cuando no se proporciona.

## Esquema esperado

La tabla `audit.events` debe existir con estructura similar:

```sql
CREATE TABLE audit.events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id VARCHAR NOT NULL,
    actor_email VARCHAR NOT NULL,
    actor_role VARCHAR NOT NULL,
    actor_ip INET,
    actor_user_agent TEXT,
    school_id UUID,
    unit_id UUID,
    service_name VARCHAR NOT NULL,
    action VARCHAR NOT NULL,
    resource_type VARCHAR NOT NULL,
    resource_id UUID,
    permission_used VARCHAR,
    request_method VARCHAR,
    request_path TEXT,
    request_id UUID,
    status_code INT,
    changes JSONB,
    metadata JSONB,
    error_message TEXT,
    severity VARCHAR NOT NULL DEFAULT 'info',
    category VARCHAR NOT NULL DEFAULT 'data',
    created_at TIMESTAMP DEFAULT NOW()
);
```

## Integración

1. **Crear instancia**
   ```go
   logger := postgres.NewPostgresAuditLogger(gormDB, "my-service")
   ```

2. **Usar directamente**
   ```go
   logger.Log(ctx, event)
   ```

3. **Desde Gin**
   ```go
   logger.LogFromGin(c, "action", "resourceType", "resourceID")
   ```

## Notas de diseño

- Implementa estrategia de adaptador: múltiples backends pueden implementar `audit.AuditLogger`.
- `internal/` previene importaciones accidentales de detalles de implementación.
- Campos opcionales como punteros permiten distinguir entre "no proporcionado" y "vacío".
- Serialización JSON automática en GORM para `Changes` y `Metadata`.
- No incluye migración de esquema: la tabla debe existir previamente.
