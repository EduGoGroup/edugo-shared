package audit

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// auditEventDB es el modelo GORM para audit.events
type auditEventDB struct {
	ID             string                 `gorm:"column:id;primaryKey;default:gen_random_uuid()"`
	ActorID        string                 `gorm:"column:actor_id;not null"`
	ActorEmail     string                 `gorm:"column:actor_email;not null"`
	ActorRole      string                 `gorm:"column:actor_role;not null"`
	ActorIP        *string                `gorm:"column:actor_ip"`
	ActorUserAgent *string                `gorm:"column:actor_user_agent"`
	SchoolID       *string                `gorm:"column:school_id"`
	UnitID         *string                `gorm:"column:unit_id"`
	ServiceName    string                 `gorm:"column:service_name;not null"`
	Action         string                 `gorm:"column:action;not null"`
	ResourceType   string                 `gorm:"column:resource_type;not null"`
	ResourceID     *string                `gorm:"column:resource_id"`
	PermissionUsed *string                `gorm:"column:permission_used"`
	RequestMethod  *string                `gorm:"column:request_method"`
	RequestPath    *string                `gorm:"column:request_path"`
	RequestID      *string                `gorm:"column:request_id"`
	StatusCode     *int                   `gorm:"column:status_code"`
	Changes        map[string]interface{} `gorm:"column:changes;serializer:json"`
	Metadata       map[string]interface{} `gorm:"column:metadata;serializer:json"`
	ErrorMessage   *string                `gorm:"column:error_message"`
	CreatedAt      time.Time              `gorm:"column:created_at;autoCreateTime"`
	Severity       string                 `gorm:"column:severity;not null;default:info"`
	Category       string                 `gorm:"column:category;not null;default:data"`
}

func (auditEventDB) TableName() string {
	return "audit.events"
}

// PostgresAuditLogger implements AuditLogger using PostgreSQL
type PostgresAuditLogger struct {
	db          *gorm.DB
	serviceName string
}

// NewPostgresAuditLogger creates a new PostgresAuditLogger
func NewPostgresAuditLogger(db *gorm.DB, serviceName string) *PostgresAuditLogger {
	return &PostgresAuditLogger{db: db, serviceName: serviceName}
}

func (l *PostgresAuditLogger) Log(ctx context.Context, event AuditEvent) error {
	if event.Severity == "" {
		event.Severity = SeverityInfo
	}
	if event.Category == "" {
		event.Category = CategoryData
	}
	event.ServiceName = l.serviceName

	record := toDBModel(event)
	return l.db.WithContext(ctx).Create(&record).Error
}

func (l *PostgresAuditLogger) LogFromGin(c *gin.Context, action, resourceType, resourceID string, opts ...AuditOption) error {
	event := AuditEvent{
		Action:         action,
		ResourceType:   resourceType,
		ResourceID:     resourceID,
		ServiceName:    l.serviceName,
		RequestMethod:  c.Request.Method,
		RequestPath:    c.Request.URL.Path,
		RequestID:      c.GetHeader("X-Request-ID"),
		ActorIP:        c.ClientIP(),
		ActorUserAgent: c.GetHeader("User-Agent"),
		StatusCode:     c.Writer.Status(),
		Severity:       SeverityInfo,
		Category:       CategoryData,
	}

	for _, opt := range opts {
		opt(&event)
	}

	// Extraer datos del JWT usando las context keys del middleware
	if userID, exists := c.Get("user_id"); exists {
		if v, ok := userID.(string); ok {
			event.ActorID = v
		}
	}
	if email, exists := c.Get("email"); exists {
		if v, ok := email.(string); ok {
			event.ActorEmail = v
		}
	}
	if role, exists := c.Get("role"); exists {
		if v, ok := role.(string); ok {
			event.ActorRole = v
		}
	}

	return l.Log(c.Request.Context(), event)
}

func toDBModel(event AuditEvent) auditEventDB {
	r := auditEventDB{
		ActorID:      event.ActorID,
		ActorEmail:   event.ActorEmail,
		ActorRole:    event.ActorRole,
		ServiceName:  event.ServiceName,
		Action:       event.Action,
		ResourceType: event.ResourceType,
		Changes:      event.Changes,
		Metadata:     event.Metadata,
		Severity:     event.Severity,
		Category:     event.Category,
	}
	if event.ActorIP != "" {
		r.ActorIP = &event.ActorIP
	}
	if event.ActorUserAgent != "" {
		r.ActorUserAgent = &event.ActorUserAgent
	}
	if event.SchoolID != "" {
		r.SchoolID = &event.SchoolID
	}
	if event.UnitID != "" {
		r.UnitID = &event.UnitID
	}
	if event.ResourceID != "" {
		r.ResourceID = &event.ResourceID
	}
	if event.PermissionUsed != "" {
		r.PermissionUsed = &event.PermissionUsed
	}
	if event.RequestMethod != "" {
		r.RequestMethod = &event.RequestMethod
	}
	if event.RequestPath != "" {
		r.RequestPath = &event.RequestPath
	}
	if event.RequestID != "" {
		r.RequestID = &event.RequestID
	}
	if event.StatusCode != 0 {
		r.StatusCode = &event.StatusCode
	}
	if event.ErrorMessage != "" {
		r.ErrorMessage = &event.ErrorMessage
	}
	return r
}
