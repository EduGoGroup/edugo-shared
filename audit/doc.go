// Package audit proporciona el contrato base para registrar eventos de auditoría.
//
// # Conceptos principales
//
// AuditEvent es la estructura que centraliza datos auditables: actor, acción, recurso, request y contexto.
//
// AuditLogger es la interfaz mínima que define el contrato Log(ctx, event) error.
// Permite múltiples implementaciones (PostgreSQL, Kafka, archivos, etc.) sin acoplar consumidores.
//
// AuditOption son funciones declarativas para enriquecer eventos con severidad, categoría, cambios, metadatos, permisos y errores.
//
// # Ejemplo de uso
//
//	event := audit.AuditEvent{
//	    ActorID:      "user-123",
//	    ActorEmail:   "user@example.com",
//	    ActorRole:    "admin",
//	    ServiceName:  "my-service",
//	    Action:       "create",
//	    ResourceType: "document",
//	}
//	event.Severity = audit.SeverityInfo
//	event.Category = audit.CategoryData
//	logger.Log(ctx, event)
//
// # Implementaciones disponibles
//
// - NoopAuditLogger: implementación noop para tests y desarrollo.
// - audit/postgres: adaptador que persiste en PostgreSQL.
package audit
