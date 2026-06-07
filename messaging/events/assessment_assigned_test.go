package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAssessmentAssignedEvent_Valid(t *testing.T) {
	due := time.Now().Add(72 * time.Hour)
	payload := AssessmentAssignedPayload{
		AssessmentID:           "assess_001",
		AssignmentID:           "assign_001",
		SchoolID:               "school_001",
		AssignedByMembershipID: "membership_001",
		SubjectOfferingID:      "offering_001",
		DueDate:                &due,
		Title:                  "Examen de Lengua",
	}

	event, err := NewAssessmentAssignedEvent("evt_001", "assessment.assigned", "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "evt_001", event.EventID)
	assert.Equal(t, "assessment.assigned", event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, "assess_001", event.Payload.AssessmentID)
	assert.Equal(t, "offering_001", event.Payload.SubjectOfferingID)
	assert.Equal(t, "Examen de Lengua", event.Payload.Title)
}

func TestNewAssessmentAssignedEvent_EmptyFields(t *testing.T) {
	base := AssessmentAssignedPayload{
		AssessmentID: "a", AssignmentID: "ai", SchoolID: "s",
		AssignedByMembershipID: "ab", SubjectOfferingID: "so",
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
			name: "AssignedByMembershipID vacio", eventID: "evt_1",
			payload: func() AssessmentAssignedPayload { p := base; p.AssignedByMembershipID = ""; return p }(),
			wantErr: "AssignedByMembershipID",
		},
		{
			name: "SubjectOfferingID vacio", eventID: "evt_1",
			payload: func() AssessmentAssignedPayload { p := base; p.SubjectOfferingID = ""; return p }(),
			wantErr: "SubjectOfferingID",
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
		AssessmentID:           "assess_001",
		AssignmentID:           "assign_001",
		SchoolID:               "school_001",
		AssignedByMembershipID: "membership_001",
		SubjectOfferingID:      "offering_3A",
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
	assert.Equal(t, event.Payload.SubjectOfferingID, decoded.Payload.SubjectOfferingID)
	assert.Equal(t, event.Payload.AssignedByMembershipID, decoded.Payload.AssignedByMembershipID)
}

func TestAssessmentAssignedEvent_DueDateSerialization(t *testing.T) {
	due := time.Date(2026, 6, 30, 23, 59, 0, 0, time.UTC)
	payload := AssessmentAssignedPayload{
		AssessmentID:           "assess_001",
		AssignmentID:           "assign_001",
		SchoolID:               "school_001",
		AssignedByMembershipID: "membership_001",
		SubjectOfferingID:      "offering_001",
		DueDate:                &due,
		Title:                  "Evaluacion Parcial",
	}

	event, err := NewAssessmentAssignedEvent("evt_001", "assessment.assigned", "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded AssessmentAssignedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	require.NotNil(t, decoded.Payload.DueDate)
	assert.True(t, decoded.Payload.DueDate.Equal(due))
	assert.Equal(t, "Evaluacion Parcial", decoded.Payload.Title)
}

func TestAssessmentAssignedEvent_OptionalFieldsOmitted(t *testing.T) {
	jsonWithoutOptionals := `{
		"event_id": "evt_001",
		"event_type": "assessment.assigned",
		"event_version": "1.0",
		"timestamp": "2026-03-28T10:00:00Z",
		"payload": {
			"assessment_id": "assess_001",
			"assignment_id": "assign_001",
			"school_id": "school_001",
			"assigned_by_membership_id": "membership_001",
			"subject_offering_id": "offering_001"
		}
	}`

	var decoded AssessmentAssignedEvent
	err := json.Unmarshal([]byte(jsonWithoutOptionals), &decoded)
	require.NoError(t, err)
	assert.Equal(t, "", decoded.Payload.Title)
	assert.Nil(t, decoded.Payload.DueDate)
}
