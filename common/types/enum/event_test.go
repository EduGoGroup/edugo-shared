package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventType_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		want      bool
	}{
		{"MaterialUploaded", EventMaterialUploaded, true},
		{"MaterialReprocess", EventMaterialReprocess, true},
		{"MaterialDeleted", EventMaterialDeleted, true},
		{"MaterialPublished", EventMaterialPublished, true},
		{"MaterialArchived", EventMaterialArchived, true},
		{"AssessmentAttemptRecorded", EventAssessmentAttemptRecorded, true},
		{"AssessmentCompleted", EventAssessmentCompleted, true},
		{"StudentEnrolled", EventStudentEnrolled, true},
		{"StudentProgress", EventStudentProgress, true},
		{"UserCreated", EventUserCreated, true},
		{"UserUpdated", EventUserUpdated, true},
		{"UserDeactivated", EventUserDeactivated, true},
		{"Invalid", "invalid_event", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.eventType.IsValid())
		})
	}
}

func TestEventType_String(t *testing.T) {
	assert.Equal(t, "material.uploaded", EventMaterialUploaded.String())
}

func TestEventType_GetRoutingKey(t *testing.T) {
	assert.Equal(t, "material.uploaded", EventMaterialUploaded.GetRoutingKey())
}
