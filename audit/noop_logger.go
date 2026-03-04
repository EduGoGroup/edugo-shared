package audit

import (
	"context"
)

// NoopAuditLogger es una implementación vacía de AuditLogger para uso en tests.
// Descarta todos los eventos sin registrarlos.
type NoopAuditLogger struct{}

// NewNoopAuditLogger crea un nuevo NoopAuditLogger para entornos de prueba.
func NewNoopAuditLogger() *NoopAuditLogger {
	return &NoopAuditLogger{}
}

// Log descarta el evento sin hacer nada.
func (l *NoopAuditLogger) Log(ctx context.Context, event AuditEvent) error {
	return nil
}
