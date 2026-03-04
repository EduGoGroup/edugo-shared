package audit

import (
	"context"
)

// Niveles de severidad para eventos de auditoría.
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
)

// Categorías de eventos de auditoría.
const (
	CategoryAuth   = "auth"
	CategoryData   = "data"
	CategoryConfig = "config"
	CategoryAdmin  = "admin"
)

// AuditEvent representa una acción auditable en el sistema.
type AuditEvent struct {
	ActorID        string
	ActorEmail     string
	ActorRole      string
	ActorIP        string
	ActorUserAgent string
	SchoolID       string
	UnitID         string
	ServiceName    string
	Action         string
	ResourceType   string
	ResourceID     string
	PermissionUsed string
	RequestMethod  string
	RequestPath    string
	RequestID      string
	StatusCode     int
	Changes        map[string]interface{}
	Metadata       map[string]interface{}
	ErrorMessage   string
	Severity       string
	Category       string
}

// AuditLogger es la interfaz para registrar eventos de auditoría.
// El método LogFromGin es un método de conveniencia disponible en la
// implementación concreta PostgresAuditLogger, pero no forma parte
// de este contrato para evitar acoplamiento con el framework Gin.
type AuditLogger interface {
	Log(ctx context.Context, event AuditEvent) error
}

// AuditOption permite configurar campos opcionales de un AuditEvent.
type AuditOption func(*AuditEvent)

// WithChanges registra los valores antes y después de una modificación.
func WithChanges(before, after interface{}) AuditOption {
	return func(e *AuditEvent) {
		e.Changes = map[string]interface{}{"before": before, "after": after}
	}
}

// WithSeverity establece el nivel de severidad del evento.
func WithSeverity(severity string) AuditOption {
	return func(e *AuditEvent) {
		e.Severity = severity
	}
}

// WithCategory establece la categoría del evento.
func WithCategory(category string) AuditOption {
	return func(e *AuditEvent) {
		e.Category = category
	}
}

// WithMetadata agrega metadatos adicionales al evento.
func WithMetadata(key string, value interface{}) AuditOption {
	return func(e *AuditEvent) {
		if e.Metadata == nil {
			e.Metadata = make(map[string]interface{})
		}
		e.Metadata[key] = value
	}
}

// WithPermission registra el permiso utilizado en la acción.
func WithPermission(permission string) AuditOption {
	return func(e *AuditEvent) {
		e.PermissionUsed = permission
	}
}

// WithError registra el mensaje de error asociado al evento.
func WithError(err error) AuditOption {
	return func(e *AuditEvent) {
		if err != nil {
			e.ErrorMessage = err.Error()
		}
	}
}
