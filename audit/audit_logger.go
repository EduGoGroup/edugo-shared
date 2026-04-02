package audit

import "context"

// AuditLogger es la interfaz para registrar eventos de auditoría.
// El método LogFromGin es un método de conveniencia disponible en la
// implementación concreta PostgresAuditLogger, pero no forma parte
// de este contrato para evitar acoplamiento con el framework Gin.
type AuditLogger interface { //nolint:revive
	Log(ctx context.Context, event AuditEvent) error
}
