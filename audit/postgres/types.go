package postgres

import (
	"context"

	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/EduGoGroup/edugo-shared/audit/postgres/internal"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostgresAuditLogger implementa audit.AuditLogger usando PostgreSQL mediante GORM.
// Persiste los eventos en la tabla audit.events.
type PostgresAuditLogger struct { //nolint:revive
	db          *gorm.DB
	serviceName string
}

// NewPostgresAuditLogger crea un nuevo PostgresAuditLogger para el servicio indicado.
// El parámetro serviceName identifica el servicio que registra el evento
// (por ejemplo: "iam-platform", "admin-api", "mobile-api").
func NewPostgresAuditLogger(db *gorm.DB, serviceName string) *PostgresAuditLogger {
	return &PostgresAuditLogger{db: db, serviceName: serviceName}
}

// Log persiste un AuditEvent en la base de datos.
// Aplica valores por defecto de Severity y Category si no están definidos.
func (l *PostgresAuditLogger) Log(ctx context.Context, event audit.AuditEvent) error {
	if event.Severity == "" {
		event.Severity = audit.SeverityInfo
	}
	if event.Category == "" {
		event.Category = audit.CategoryData
	}
	event.ServiceName = l.serviceName

	if event.ActorID == "" {
		event.ActorID = internal.DefaultActorID
	}
	if event.ActorEmail == "" {
		event.ActorEmail = internal.DefaultActorEmail
	}
	if event.ActorRole == "" {
		event.ActorRole = internal.DefaultActorRole
	}

	record := internal.ToDBModel(event)
	return l.db.WithContext(ctx).Create(&record).Error
}

// LogFromGin es un método de conveniencia que extrae los datos del contexto Gin
// y delega en Log. No forma parte de la interfaz audit.AuditLogger.
func (l *PostgresAuditLogger) LogFromGin(c *gin.Context, action, resourceType, resourceID string, opts ...audit.AuditOption) error {
	event := audit.AuditEvent{
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
		Severity:       audit.SeverityInfo,
		Category:       audit.CategoryData,
	}

	for _, opt := range opts {
		opt(&event)
	}

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
