package postgres

import (
	"context"
	"testing"

	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/EduGoGroup/edugo-shared/audit/postgres/internal"
)

func TestNewPostgresAuditLogger(t *testing.T) {
	logger := NewPostgresAuditLogger(nil, "test-service")

	if logger == nil {
		t.Fatal("Expected logger to be created, got nil")
	}
}

func TestToDBModelConvertsRequiredFields(t *testing.T) {
	event := audit.AuditEvent{
		ActorID:      "actor-123",
		ActorEmail:   "actor@example.com",
		ActorRole:    "admin",
		ServiceName:  "test-service",
		Action:       "create",
		ResourceType: "user",
		Severity:     audit.SeverityInfo,
		Category:     audit.CategoryData,
	}

	model := internal.ToDBModel(event)

	if model.ActorID != event.ActorID {
		t.Error("ActorID not converted")
	}
	if model.ActorEmail != event.ActorEmail {
		t.Error("ActorEmail not converted")
	}
	if model.ActorRole != event.ActorRole {
		t.Error("ActorRole not converted")
	}
	if model.ServiceName != event.ServiceName {
		t.Error("ServiceName not converted")
	}
	if model.Action != event.Action {
		t.Error("Action not converted")
	}
	if model.ResourceType != event.ResourceType {
		t.Error("ResourceType not converted")
	}
	if model.Severity != event.Severity {
		t.Error("Severity not converted")
	}
	if model.Category != event.Category {
		t.Error("Category not converted")
	}
}

func TestToDBModelConvertsOptionalFields(t *testing.T) {
	event := audit.AuditEvent{
		ActorID:        "actor-123",
		ActorEmail:     "actor@example.com",
		ActorRole:      "admin",
		ActorIP:        "192.168.1.1",
		ActorUserAgent: "Mozilla/5.0",
		SchoolID:       "school-456",
		UnitID:         "unit-789",
		ServiceName:    "test-service",
		Action:         "create",
		ResourceType:   "user",
		ResourceID:     "resource-111",
		RequestMethod:  "POST",
		StatusCode:     201,
	}

	model := internal.ToDBModel(event)

	if model.ActorIP == nil || *model.ActorIP != event.ActorIP {
		t.Error("ActorIP not converted to pointer")
	}
	if model.ActorUserAgent == nil || *model.ActorUserAgent != event.ActorUserAgent {
		t.Error("ActorUserAgent not converted to pointer")
	}
	if model.SchoolID == nil || *model.SchoolID != event.SchoolID {
		t.Error("SchoolID not converted to pointer")
	}
	if model.UnitID == nil || *model.UnitID != event.UnitID {
		t.Error("UnitID not converted to pointer")
	}
	if model.ResourceID == nil || *model.ResourceID != event.ResourceID {
		t.Error("ResourceID not converted to pointer")
	}
	if model.RequestMethod == nil || *model.RequestMethod != event.RequestMethod {
		t.Error("RequestMethod not converted to pointer")
	}
	if model.StatusCode == nil || *model.StatusCode != event.StatusCode {
		t.Error("StatusCode not converted to pointer")
	}
}

func TestToDBModelHandlesEmptyOptionalFields(t *testing.T) {
	event := audit.AuditEvent{
		ActorID:      "actor-123",
		ActorEmail:   "actor@example.com",
		ActorRole:    "user",
		ServiceName:  "test-service",
		Action:       "read",
		ResourceType: "file",
		// All optional fields are empty
		ActorIP:        "",
		ActorUserAgent: "",
		SchoolID:       "",
		UnitID:         "",
		ResourceID:     "",
		PermissionUsed: "",
		RequestMethod:  "",
		RequestPath:    "",
		RequestID:      "",
		StatusCode:     0,
		ErrorMessage:   "",
	}

	model := internal.ToDBModel(event)

	// Empty strings should become nil pointers
	if model.ActorIP != nil {
		t.Error("ActorIP should be nil for empty string")
	}
	if model.ActorUserAgent != nil {
		t.Error("ActorUserAgent should be nil for empty string")
	}
	if model.SchoolID != nil {
		t.Error("SchoolID should be nil for empty string")
	}
	if model.UnitID != nil {
		t.Error("UnitID should be nil for empty string")
	}
	if model.ResourceID != nil {
		t.Error("ResourceID should be nil for empty string")
	}
	if model.PermissionUsed != nil {
		t.Error("PermissionUsed should be nil for empty string")
	}
	if model.RequestMethod != nil {
		t.Error("RequestMethod should be nil for empty string")
	}
	if model.RequestPath != nil {
		t.Error("RequestPath should be nil for empty string")
	}
	if model.RequestID != nil {
		t.Error("RequestID should be nil for empty string")
	}
	if model.ErrorMessage != nil {
		t.Error("ErrorMessage should be nil for empty string")
	}

	// Zero status code should be nil
	if model.StatusCode != nil {
		t.Error("StatusCode should be nil for zero value")
	}
}

func TestToDBModelPreservesNonZeroValues(t *testing.T) {
	event := audit.AuditEvent{
		ActorID:       "actor-123",
		ActorEmail:    "actor@example.com",
		ActorRole:     "user",
		ActorIP:       "10.0.0.1",
		ServiceName:   "test-service",
		Action:        "delete",
		ResourceType:  "record",
		ResourceID:    "rec-999",
		RequestMethod: "DELETE",
		StatusCode:    204,
		ErrorMessage:  "Record deleted",
		Severity:      audit.SeverityWarning,
		Category:      audit.CategoryAuth,
	}

	model := internal.ToDBModel(event)

	if model.ActorIP == nil || *model.ActorIP != "10.0.0.1" {
		t.Error("ActorIP should be preserved")
	}
	if model.StatusCode == nil || *model.StatusCode != 204 {
		t.Error("StatusCode should be preserved")
	}
	if model.ErrorMessage == nil || *model.ErrorMessage != "Record deleted" {
		t.Error("ErrorMessage should be preserved")
	}
}

func TestAuditEventDBTableName(t *testing.T) {
	model := internal.AuditEventDB{}
	expected := "audit.events"

	if model.TableName() != expected {
		t.Errorf("Expected TableName to return %q, got %q", expected, model.TableName())
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Verify context is actually canceled
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Fatal("Context should be canceled")
	}
}

func TestMapSerializationFields(t *testing.T) {
	changes := map[string]any{
		"field1": "value1",
		"field2": 123,
		"nested": map[string]any{
			"key": "val",
		},
	}

	metadata := map[string]any{
		"correlation_id": "corr-123",
		"trace_id":       "trace-456",
	}

	event := audit.AuditEvent{
		ActorID:      "actor-123",
		ActorEmail:   "actor@example.com",
		ActorRole:    "user",
		ServiceName:  "test-service",
		Action:       "modify",
		ResourceType: "config",
		Changes:      changes,
		Metadata:     metadata,
	}

	model := internal.ToDBModel(event)

	if model.Changes == nil {
		t.Error("Changes should be preserved")
	}
	if model.Metadata == nil {
		t.Error("Metadata should be preserved")
	}
}

func TestSeverityAndCategoryPreservation(t *testing.T) {
	severities := []string{
		audit.SeverityInfo,
		audit.SeverityWarning,
		audit.SeverityCritical,
	}

	categories := []string{
		audit.CategoryAuth,
		audit.CategoryData,
		audit.CategoryConfig,
		audit.CategoryAdmin,
	}

	for _, severity := range severities {
		event := audit.AuditEvent{
			ActorID:      "actor-123",
			ActorEmail:   "actor@example.com",
			ActorRole:    "user",
			ServiceName:  "test-service",
			Action:       "test",
			ResourceType: "test",
			Severity:     severity,
		}

		model := internal.ToDBModel(event)
		if model.Severity != severity {
			t.Errorf("Severity %q not preserved", severity)
		}
	}

	for _, category := range categories {
		event := audit.AuditEvent{
			ActorID:      "actor-123",
			ActorEmail:   "actor@example.com",
			ActorRole:    "user",
			ServiceName:  "test-service",
			Action:       "test",
			ResourceType: "test",
			Category:     category,
		}

		model := internal.ToDBModel(event)
		if model.Category != category {
			t.Errorf("Category %q not preserved", category)
		}
	}
}
