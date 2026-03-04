package audit

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Severity levels
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
)

// Categories
const (
	CategoryAuth   = "auth"
	CategoryData   = "data"
	CategoryConfig = "config"
	CategoryAdmin  = "admin"
)

// AuditEvent represents an auditable action
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

// AuditLogger interface for writing audit events
type AuditLogger interface {
	Log(ctx context.Context, event AuditEvent) error
	LogFromGin(c *gin.Context, action, resourceType, resourceID string, opts ...AuditOption) error
}

// AuditOption for optional fields
type AuditOption func(*AuditEvent)

func WithChanges(before, after interface{}) AuditOption {
	return func(e *AuditEvent) {
		e.Changes = map[string]interface{}{"before": before, "after": after}
	}
}

func WithSeverity(severity string) AuditOption {
	return func(e *AuditEvent) {
		e.Severity = severity
	}
}

func WithCategory(category string) AuditOption {
	return func(e *AuditEvent) {
		e.Category = category
	}
}

func WithMetadata(key string, value interface{}) AuditOption {
	return func(e *AuditEvent) {
		if e.Metadata == nil {
			e.Metadata = make(map[string]interface{})
		}
		e.Metadata[key] = value
	}
}

func WithPermission(permission string) AuditOption {
	return func(e *AuditEvent) {
		e.PermissionUsed = permission
	}
}

func WithError(err error) AuditOption {
	return func(e *AuditEvent) {
		if err != nil {
			e.ErrorMessage = err.Error()
		}
	}
}
