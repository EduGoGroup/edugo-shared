package audit

// AuditOption permite configurar campos opcionales de un AuditEvent.
type AuditOption func(*AuditEvent) //nolint:revive

// WithChanges registra los valores antes y después de una modificación.
func WithChanges(before, after any) AuditOption {
	return func(e *AuditEvent) {
		e.Changes = map[string]any{"before": before, "after": after}
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
func WithMetadata(key string, value any) AuditOption {
	return func(e *AuditEvent) {
		if e.Metadata == nil {
			e.Metadata = make(map[string]any)
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
