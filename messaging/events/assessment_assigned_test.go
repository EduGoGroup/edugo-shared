package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAssessmentAssignedEvent_Valid(t *testing.T) {
	payload := AssessmentAssignedPayload{
		AssessmentID: "assess_001",
		AssignmentID: "assign_001",
		SchoolID:     "school_001",
		AssignedByID: "teacher_001",
		TargetType:   "student",
		TargetID:     "student_001",
		Title:        "Examen de Lengua",
	}

	event, err := NewAssessmentAssignedEvent("evt_001", "assessment.assigned", "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "evt_001", event.EventID)
	assert.Equal(t, "assessment.assigned", event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, "assess_001", event.Payload.AssessmentID)
	assert.Equal(t, "student", event.Payload.TargetType)
	assert.Equal(t, "Examen de Lengua", event.Payload.Title)
}

func TestNewAssessmentAssignedEvent_EmptyFields(t *testing.T) {
	base := AssessmentAssignedPayload{
		AssessmentID: "a", AssignmentID: "ai", SchoolID: "s",
		AssignedByID: "ab", TargetType: "student", TargetID: "t",
	}

	tests := []struct {
		name    string
		eventID string
		payload AssessmentAssignedPayload
		wantErr string
	}{
		{
			name: "eventID vacio", eventID: "",
			payload: base, wantErr: "eventID",
		},
		{
			name: "AssessmentID vacio", eventID: "evt_1",
			payload: func() AssessmentAssignedPayload { p := base; p.AssessmentID = ""; return p }(),
			wantErr: "AssessmentID",
		},
		{
			name: "AssignmentID vacio", eventID: "evt_1",
			payload: func() AssessmentAssignedPayload { p := base; p.AssignmentID = ""; return p }(),
			wantErr: "AssignmentID",
		},
		{
			name: "SchoolID vacio", eventID: "evt_1",
			payload: func() AssessmentAssignedPayload { p := base; p.SchoolID = ""; return p }(),
			wantErr: "SchoolID",
		},
		{
			name: "AssignedByID vacio", eventID: "evt_1",
			payload: func() AssessmentAssignedPayload { p := base; p.AssignedByID = ""; return p }(),
			wantErr: "AssignedByID",
		},
		{
			name: "TargetType vacio", eventID: "evt_1",
			payload: func() AssessmentAssignedPayload { p := base; p.TargetType = ""; return p }(),
			wantErr: "TargetType",
		},
		{
			name: "TargetID vacio", eventID: "evt_1",
			payload: func() AssessmentAssignedPayload { p := base; p.TargetID = ""; return p }(),
			wantErr: "TargetID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAssessmentAssignedEvent(tt.eventID, "assessment.assigned", "1.0", tt.payload)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestAssessmentAssignedEvent_Serialization(t *testing.T) {
	payload := AssessmentAssignedPayload{
		AssessmentID: "assess_001",
		AssignmentID: "assign_001",
		SchoolID:     "school_001",
		AssignedByID: "teacher_001",
		TargetType:   "unit",
		TargetID:     "unit_3A",
	}

	event, err := NewAssessmentAssignedEvent("evt_001", "assessment.assigned", "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded AssessmentAssignedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, event.EventID, decoded.EventID)
	assert.Equal(t, event.Payload.AssessmentID, decoded.Payload.AssessmentID)
	assert.Equal(t, event.Payload.TargetType, decoded.Payload.TargetType)
	assert.Equal(t, event.Payload.TargetID, decoded.Payload.TargetID)
}

func TestAssessmentAssignedEvent_TitleSerialization(t *testing.T) {
	payload := AssessmentAssignedPayload{
		AssessmentID: "assess_001",
		AssignmentID: "assign_001",
		SchoolID:     "school_001",
		AssignedByID: "teacher_001",
		TargetType:   "student",
		TargetID:     "student_001",
		Title:        "Evaluacion Parcial",
	}

	event, err := NewAssessmentAssignedEvent("evt_001", "assessment.assigned", "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded AssessmentAssignedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "Evaluacion Parcial", decoded.Payload.Title)
}

func TestAssessmentAssignedEvent_BackwardCompatibility(t *testing.T) {
	jsonWithoutTitle := `{
		"event_id": "evt_001",
		"event_type": "assessment.assigned",
		"event_version": "1.0",
		"timestamp": "2026-03-28T10:00:00Z",
		"payload": {
			"assessment_id": "assess_001",
			"assignment_id": "assign_001",
			"school_id": "school_001",
			"assigned_by_id": "teacher_001",
			"target_type": "student",
			"target_id": "student_001"
		}
	}`

	var decoded AssessmentAssignedEvent
	err := json.Unmarshal([]byte(jsonWithoutTitle), &decoded)
	require.NoError(t, err)
	assert.Equal(t, "", decoded.Payload.Title)
}
