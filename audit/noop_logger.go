package audit

import (
	"context"

	"github.com/gin-gonic/gin"
)

// NoopAuditLogger does nothing — for tests
type NoopAuditLogger struct{}

func NewNoopAuditLogger() *NoopAuditLogger {
	return &NoopAuditLogger{}
}

func (l *NoopAuditLogger) Log(ctx context.Context, event AuditEvent) error {
	return nil
}

func (l *NoopAuditLogger) LogFromGin(c *gin.Context, action, resourceType, resourceID string, opts ...AuditOption) error {
	return nil
}
