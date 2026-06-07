package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAssessmentAttemptRecordedEvent_Valid(t *testing.T) {
	payload := AssessmentAttemptRecordedPayload{
		AttemptID:           "attempt_001",
		AssessmentID:        "assess_001",
		StudentMembershipID: "membership_001",
		SubjectID:           "subject_001",
		SchoolID:            "school_001",
		Score:               85.5,
		MaxScore:            100.0,
		Status:              "completed",
		SubmittedAt:         time.Now(),
		TeacherID:           "teacher_001",
		Title:               "Examen de Matematicas",
	}

	event, err := NewAssessmentAttemptRecordedEvent("evt_001", "assessment.attempt_recorded", "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "evt_001", event.EventID)
	assert.Equal(t, "assessment.attempt_recorded", event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, "membership_001", event.Payload.StudentMembershipID)
	assert.Equal(t, "subject_001", event.Payload.SubjectID)
	assert.Equal(t, 85.5, event.Payload.Score)
	assert.Equal(t, 100.0, event.Payload.MaxScore)
	assert.Equal(t, "completed", event.Payload.Status)
	assert.Equal(t, "teacher_001", event.Payload.TeacherID)
	assert.Equal(t, "Examen de Matematicas", event.Payload.Title)
}

func TestNewAssessmentAttemptRecordedEvent_PendingReviewStatus(t *testing.T) {
	payload := AssessmentAttemptRecordedPayload{
		AttemptID:           "attempt_002",
		AssessmentID:        "assess_002",
		StudentMembershipID: "membership_002",
		SubjectID:           "subject_002",
		SchoolID:            "school_001",
		Score:               0,
		MaxScore:            100.0,
		Status:              "pending_review",
		SubmittedAt:         time.Now(),
	}

	event, err := NewAssessmentAttemptRecordedEvent("evt_002", "assessment.attempt_recorded", "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "pending_review", event.Payload.Status)
}

func TestNewAssessmentAttemptRecordedEvent_EmptyFields(t *testing.T) {
	base := AssessmentAttemptRecordedPayload{
		AttemptID: "at", AssessmentID: "a", StudentMembershipID: "m", SubjectID: "su", SchoolID: "s",
		Score: 80, MaxScore: 100, Status: "completed", SubmittedAt: time.Now(),
	}

	tests := []struct {
		name    string
		eventID string
		payload AssessmentAttemptRecordedPayload
		wantErr string
	}{
		{
			name: "eventID vacio", eventID: "",
			payload: base, wantErr: "eventID",
		},
		{
			name: "AttemptID vacio", eventID: "evt_1",
			payload: func() AssessmentAttemptRecordedPayload { p := base; p.AttemptID = ""; return p }(),
			wantErr: "AttemptID",
		},
		{
			name: "AssessmentID vacio", eventID: "evt_1",
			payload: func() AssessmentAttemptRecordedPayload { p := base; p.AssessmentID = ""; return p }(),
			wantErr: "AssessmentID",
		},
		{
			name: "StudentMembershipID vacio", eventID: "evt_1",
			payload: func() AssessmentAttemptRecordedPayload { p := base; p.StudentMembershipID = ""; return p }(),
			wantErr: "StudentMembershipID",
		},
		{
			name: "SubjectID vacio", eventID: "evt_1",
			payload: func() AssessmentAttemptRecordedPayload { p := base; p.SubjectID = ""; return p }(),
			wantErr: "SubjectID",
		},
		{
			name: "SchoolID vacio", eventID: "evt_1",
			payload: func() AssessmentAttemptRecordedPayload { p := base; p.SchoolID = ""; return p }(),
			wantErr: "SchoolID",
		},
		{
			name: "MaxScore negativo", eventID: "evt_1",
			payload: func() AssessmentAttemptRecordedPayload { p := base; p.MaxScore = -1; return p }(),
			wantErr: "MaxScore",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAssessmentAttemptRecordedEvent(tt.eventID, "assessment.attempt_recorded", "1.0", tt.payload)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestAssessmentAttemptRecordedEvent_Serialization(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	payload := AssessmentAttemptRecordedPayload{
		AttemptID:           "attempt_001",
		AssessmentID:        "assess_001",
		StudentMembershipID: "membership_001",
		SubjectID:           "subject_001",
		SchoolID:            "school_001",
		Score:               92.0,
		MaxScore:            100.0,
		Status:              "completed",
		SubmittedAt:         now,
	}

	event, err := NewAssessmentAttemptRecordedEvent("evt_001", "assessment.attempt_recorded", "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	// Verifica los json tags nuevos en el documento serializado.
	var raw map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &raw))
	var rawPayload map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw["payload"], &rawPayload))
	assert.Contains(t, rawPayload, "student_membership_id")
	assert.Contains(t, rawPayload, "subject_id")
	assert.Contains(t, rawPayload, "max_score")
	assert.Contains(t, rawPayload, "status")
	assert.NotContains(t, rawPayload, "student_id")
	assert.NotContains(t, rawPayload, "total_points")

	var decoded AssessmentAttemptRecordedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, event.EventID, decoded.EventID)
	assert.Equal(t, event.Payload.AttemptID, decoded.Payload.AttemptID)
	assert.Equal(t, "membership_001", decoded.Payload.StudentMembershipID)
	assert.Equal(t, "subject_001", decoded.Payload.SubjectID)
	assert.Equal(t, 92.0, decoded.Payload.Score)
	assert.Equal(t, 100.0, decoded.Payload.MaxScore)
	assert.Equal(t, "completed", decoded.Payload.Status)
	assert.True(t, decoded.Payload.SubmittedAt.Equal(now))
}

func TestAssessmentAttemptRecordedEvent_OptionalFieldsSerialization(t *testing.T) {
	payload := AssessmentAttemptRecordedPayload{
		AttemptID:           "attempt_001",
		AssessmentID:        "assess_001",
		StudentMembershipID: "membership_001",
		SubjectID:           "subject_001",
		SchoolID:            "school_001",
		Score:               90.0,
		MaxScore:            100.0,
		Status:              "completed",
		SubmittedAt:         time.Now().Truncate(time.Second),
		TeacherID:           "teacher_001",
		Title:               "Examen Final",
	}

	event, err := NewAssessmentAttemptRecordedEvent("evt_001", "assessment.attempt_recorded", "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded AssessmentAttemptRecordedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "teacher_001", decoded.Payload.TeacherID)
	assert.Equal(t, "Examen Final", decoded.Payload.Title)
}

func TestAssessmentAttemptRecordedEvent_Deserialization(t *testing.T) {
	jsonDoc := `{
		"event_id": "evt_001",
		"event_type": "assessment.attempt_recorded",
		"event_version": "1.0",
		"timestamp": "2026-03-28T10:00:00Z",
		"payload": {
			"attempt_id": "attempt_001",
			"assessment_id": "assess_001",
			"student_membership_id": "membership_001",
			"subject_id": "subject_001",
			"school_id": "school_001",
			"score": 85.0,
			"max_score": 100.0,
			"status": "completed",
			"submitted_at": "2026-03-28T10:00:00Z"
		}
	}`

	var decoded AssessmentAttemptRecordedEvent
	err := json.Unmarshal([]byte(jsonDoc), &decoded)
	require.NoError(t, err)
	assert.Equal(t, "membership_001", decoded.Payload.StudentMembershipID)
	assert.Equal(t, "subject_001", decoded.Payload.SubjectID)
	assert.Equal(t, 100.0, decoded.Payload.MaxScore)
	assert.Equal(t, "completed", decoded.Payload.Status)
	assert.Equal(t, "", decoded.Payload.TeacherID)
	assert.Equal(t, "", decoded.Payload.Title)
}
