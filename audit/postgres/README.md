# Audit PostgreSQL

Adaptador de auditoría que normaliza eventos y los persiste en `audit.events` usando GORM.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/audit/postgres
```

## Uso rápido

```go
import "github.com/EduGoGroup/edugo-shared/audit/postgres"

// Crear logger
logger := postgres.NewPostgresAuditLogger(db, "my-service")

// Log directo
logger.Log(ctx, audit.AuditEvent{
    Action:       "create",
    ResourceType: "user",
    ActorID:      "user-123",
    ActorEmail:   "user@example.com",
})

// Desde contexto Gin
logger.LogFromGin(c, "update", "profile", "profile-456")
```

## API Pública

- `NewPostgresAuditLogger(db *gorm.DB, serviceName string) *PostgresAuditLogger`
  - Crea una nueva instancia del logger para el servicio indicado.

- `Log(ctx context.Context, event audit.AuditEvent) error`
  - Persiste un evento de auditoría. Aplica defaults automáticos para `Severity` e `Category`.

- `LogFromGin(c *gin.Context, action, resourceType, resourceID string, opts ...audit.AuditOption) error`
  - Extrae contexto HTTP de Gin y persiste el evento. Permite pasar opciones adicionales.

## Estructura del módulo

```
├── types.go              # API pública
├── doc.go               # Documentación
├── go.mod              # Definición del módulo
└── internal/           # Implementación privada
    ├── models.go       # Modelos GORM
    ├── converter.go    # Conversión de eventos
    └── defaults.go     # Constantes internas
```

## Requisitos

- PostgreSQL con tabla `audit.events` y serializers JSON configurados.
- GORM >= 1.25
- Gin >= 1.9 (solo si usas `LogFromGin`)

## Comandos disponibles

```bash
make build     # Compilar el módulo
make test      # Ejecutar tests
make check     # Lint y validación
```

## Dependencias

- `github.com/EduGoGroup/edugo-shared/audit` - Contrato de auditoría
- `gorm.io/gorm` - ORM para PostgreSQL
- `github.com/gin-gonic/gin` - Framework HTTP (opcional)
