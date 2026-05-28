package internal

import "time"

// AuditEventDB es el modelo GORM para la tabla audit.events.
type AuditEventDB struct {
	ID             string         `gorm:"column:id;primaryKey;default:gen_random_uuid()"`
	ActorID        string         `gorm:"column:actor_id;not null"`
	ActorEmail     string         `gorm:"column:actor_email;not null"`
	ActorRole      string         `gorm:"column:actor_role;not null"`
	ActorIP        *string        `gorm:"column:actor_ip"`
	ActorUserAgent *string        `gorm:"column:actor_user_agent"`
	SchoolID       *string        `gorm:"column:school_id"`
	UnitID         *string        `gorm:"column:unit_id"`
	ServiceName    string         `gorm:"column:service_name;not null"`
	Action         string         `gorm:"column:action;not null"`
	ResourceType   string         `gorm:"column:resource_type;not null"`
	ResourceID     *string        `gorm:"column:resource_id"`
	PermissionUsed *string        `gorm:"column:permission_used"`
	RequestMethod  *string        `gorm:"column:request_method"`
	RequestPath    *string        `gorm:"column:request_path"`
	RequestID      *string        `gorm:"column:request_id"`
	StatusCode     *int           `gorm:"column:status_code"`
	Changes        map[string]any `gorm:"column:changes;serializer:json"`
	Metadata       map[string]any `gorm:"column:metadata;serializer:json"`
	ErrorMessage   *string        `gorm:"column:error_message"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	Severity       string         `gorm:"column:severity;not null;default:info"`
	Category       string         `gorm:"column:category;not null;default:data"`
}

// TableName especifica el nombre de la tabla en la base de datos.
func (AuditEventDB) TableName() string {
	return "audit.events"
}
