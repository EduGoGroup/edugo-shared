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
		AttemptID:    "attempt_001",
		AssessmentID: "assess_001",
		StudentID:    "student_001",
		SchoolID:     "school_001",
		Score:        85.5,
		TotalPoints:  100.0,
		SubmittedAt:  time.Now(),
		TeacherID:    "teacher_001",
		Title:        "Examen de Matematicas",
	}

	event, err := NewAssessmentAttemptRecordedEvent("evt_001", "assessment.attempt_recorded", "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "evt_001", event.EventID)
	assert.Equal(t, "assessment.attempt_recorded", event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, 85.5, event.Payload.Score)
	assert.Equal(t, 100.0, event.Payload.TotalPoints)
	assert.Equal(t, "teacher_001", event.Payload.TeacherID)
	assert.Equal(t, "Examen de Matematicas", event.Payload.Title)
}

func TestNewAssessmentAttemptRecordedEvent_EmptyFields(t *testing.T) {
	base := AssessmentAttemptRecordedPayload{
		AttemptID: "at", AssessmentID: "a", StudentID: "st", SchoolID: "s",
		Score: 80, TotalPoints: 100, SubmittedAt: time.Now(),
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
			name: "StudentID vacio", eventID: "evt_1",
			payload: func() AssessmentAttemptRecordedPayload { p := base; p.StudentID = ""; return p }(),
			wantErr: "StudentID",
		},
		{
			name: "SchoolID vacio", eventID: "evt_1",
			payload: func() AssessmentAttemptRecordedPayload { p := base; p.SchoolID = ""; return p }(),
			wantErr: "SchoolID",
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
		AttemptID:    "attempt_001",
		AssessmentID: "assess_001",
		StudentID:    "student_001",
		SchoolID:     "school_001",
		Score:        92.0,
		TotalPoints:  100.0,
		SubmittedAt:  now,
	}

	event, err := NewAssessmentAttemptRecordedEvent("evt_001", "assessment.attempt_recorded", "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded AssessmentAttemptRecordedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, event.EventID, decoded.EventID)
	assert.Equal(t, event.Payload.AttemptID, decoded.Payload.AttemptID)
	assert.Equal(t, event.Payload.Score, decoded.Payload.Score)
	assert.True(t, decoded.Payload.SubmittedAt.Equal(now))
}

func TestAssessmentAttemptRecordedEvent_OptionalFieldsSerialization(t *testing.T) {
	payload := AssessmentAttemptRecordedPayload{
		AttemptID:    "attempt_001",
		AssessmentID: "assess_001",
		StudentID:    "student_001",
		SchoolID:     "school_001",
		Score:        90.0,
		TotalPoints:  100.0,
		SubmittedAt:  time.Now().Truncate(time.Second),
		TeacherID:    "teacher_001",
		Title:        "Examen Final",
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

func TestAssessmentAttemptRecordedEvent_BackwardCompatibility(t *testing.T) {
	jsonWithoutNewFields := `{
		"event_id": "evt_001",
		"event_type": "assessment.attempt_recorded",
		"event_version": "1.0",
		"timestamp": "2026-03-28T10:00:00Z",
		"payload": {
			"attempt_id": "attempt_001",
			"assessment_id": "assess_001",
			"student_id": "student_001",
			"school_id": "school_001",
			"score": 85.0,
			"total_points": 100.0,
			"submitted_at": "2026-03-28T10:00:00Z"
		}
	}`

	var decoded AssessmentAttemptRecordedEvent
	err := json.Unmarshal([]byte(jsonWithoutNewFields), &decoded)
	require.NoError(t, err)
	assert.Equal(t, "", decoded.Payload.TeacherID)
	assert.Equal(t, "", decoded.Payload.Title)
}
