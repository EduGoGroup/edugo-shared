package internal

import "github.com/EduGoGroup/edugo-shared/audit"

// ToDBModel convierte un audit.AuditEvent al modelo GORM AuditEventDB.
func ToDBModel(event audit.AuditEvent) AuditEventDB {
	r := AuditEventDB{
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
