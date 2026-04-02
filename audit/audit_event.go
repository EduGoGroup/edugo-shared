package audit

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
type AuditEvent struct { //nolint:revive
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
	Changes        map[string]any
	Metadata       map[string]any
	ErrorMessage   string
	Severity       string
	Category       string
}
