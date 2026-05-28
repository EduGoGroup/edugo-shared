package audit

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoopAuditLogger_Log(t *testing.T) {
	logger := NewNoopAuditLogger()
	err := logger.Log(context.Background(), AuditEvent{
		Action:       "user.login",
		ResourceType: "auth",
		Severity:     SeverityInfo,
		Category:     CategoryAuth,
	})
	assert.NoError(t, err, "NoopAuditLogger nunca debe retornar error")
}

func TestWithChanges(t *testing.T) {
	event := AuditEvent{}
	opt := WithChanges("valorAntes", "valorDespues")
	opt(&event)

	require.NotNil(t, event.Changes)
	assert.Equal(t, "valorAntes", event.Changes["before"])
	assert.Equal(t, "valorDespues", event.Changes["after"])
}

func TestWithSeverity(t *testing.T) {
	event := AuditEvent{}
	WithSeverity(SeverityWarning)(&event)
	assert.Equal(t, SeverityWarning, event.Severity)
}

func TestWithCategory(t *testing.T) {
	event := AuditEvent{}
	WithCategory(CategoryAdmin)(&event)
	assert.Equal(t, CategoryAdmin, event.Category)
}

func TestWithMetadata(t *testing.T) {
	event := AuditEvent{}
	WithMetadata("school_id", "sch-123")(&event)
	WithMetadata("unit_id", "unit-456")(&event)

	require.NotNil(t, event.Metadata)
	assert.Equal(t, "sch-123", event.Metadata["school_id"])
	assert.Equal(t, "unit-456", event.Metadata["unit_id"])
}

func TestWithPermission(t *testing.T) {
	event := AuditEvent{}
	WithPermission("audit:read")(&event)
	assert.Equal(t, "audit:read", event.PermissionUsed)
}

func TestWithError(t *testing.T) {
	event := AuditEvent{}
	err := errors.New("acceso denegado")
	WithError(err)(&event)
	assert.Equal(t, "acceso denegado", event.ErrorMessage)
}

func TestWithError_NilError(t *testing.T) {
	event := AuditEvent{}
	WithError(nil)(&event)
	assert.Empty(t, event.ErrorMessage, "nil error no debe asignar mensaje")
}

func TestWithMetadata_InicializaMap(t *testing.T) {
	event := AuditEvent{Metadata: nil}
	WithMetadata("clave", "valor")(&event)
	require.NotNil(t, event.Metadata)
	assert.Equal(t, "valor", event.Metadata["clave"])
}

func TestConstantes_Severidad(t *testing.T) {
	assert.Equal(t, "info", SeverityInfo)
	assert.Equal(t, "warning", SeverityWarning)
	assert.Equal(t, "critical", SeverityCritical)
}

func TestConstantes_Categoria(t *testing.T) {
	assert.Equal(t, "auth", CategoryAuth)
	assert.Equal(t, "data", CategoryData)
	assert.Equal(t, "config", CategoryConfig)
	assert.Equal(t, "admin", CategoryAdmin)
}

func TestAuditLogger_InterfaceSatisfiedByNoop(t *testing.T) {
	var _ AuditLogger = NewNoopAuditLogger()
}
