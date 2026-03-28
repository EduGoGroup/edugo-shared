package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAssessmentReviewedEvent_Valid(t *testing.T) {
	payload := AssessmentReviewedPayload{
		AttemptID:    "attempt_001",
		AssessmentID: "assess_001",
		ReviewerID:   "teacher_001",
		SchoolID:     "school_001",
		FinalScore:   88.5,
		TotalPoints:  100.0,
		Status:       "approved",
	}

	event, err := NewAssessmentReviewedEvent("evt_001", "assessment.reviewed", "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "evt_001", event.EventID)
	assert.Equal(t, "assessment.reviewed", event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, 88.5, event.Payload.FinalScore)
	assert.Equal(t, "approved", event.Payload.Status)
}

func TestNewAssessmentReviewedEvent_EmptyFields(t *testing.T) {
	base := AssessmentReviewedPayload{
		AttemptID: "at", AssessmentID: "a", ReviewerID: "r", SchoolID: "s",
		FinalScore: 80, TotalPoints: 100, Status: "approved",
	}

	tests := []struct {
		name    string
		eventID string
		payload AssessmentReviewedPayload
		wantErr string
	}{
		{
			name: "eventID vacio", eventID: "",
			payload: base, wantErr: "eventID",
		},
		{
			name: "AttemptID vacio", eventID: "evt_1",
			payload: func() AssessmentReviewedPayload { p := base; p.AttemptID = ""; return p }(),
			wantErr: "AttemptID",
		},
		{
			name: "AssessmentID vacio", eventID: "evt_1",
			payload: func() AssessmentReviewedPayload { p := base; p.AssessmentID = ""; return p }(),
			wantErr: "AssessmentID",
		},
		{
			name: "ReviewerID vacio", eventID: "evt_1",
			payload: func() AssessmentReviewedPayload { p := base; p.ReviewerID = ""; return p }(),
			wantErr: "ReviewerID",
		},
		{
			name: "SchoolID vacio", eventID: "evt_1",
			payload: func() AssessmentReviewedPayload { p := base; p.SchoolID = ""; return p }(),
			wantErr: "SchoolID",
		},
		{
			name: "Status vacio", eventID: "evt_1",
			payload: func() AssessmentReviewedPayload { p := base; p.Status = ""; return p }(),
			wantErr: "status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAssessmentReviewedEvent(tt.eventID, "assessment.reviewed", "1.0", tt.payload)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestAssessmentReviewedEvent_Serialization(t *testing.T) {
	payload := AssessmentReviewedPayload{
		AttemptID:    "attempt_001",
		AssessmentID: "assess_001",
		ReviewerID:   "teacher_001",
		SchoolID:     "school_001",
		FinalScore:   95.0,
		TotalPoints:  100.0,
		Status:       "approved",
	}

	event, err := NewAssessmentReviewedEvent("evt_001", "assessment.reviewed", "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded AssessmentReviewedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, event.EventID, decoded.EventID)
	assert.Equal(t, event.Payload.AttemptID, decoded.Payload.AttemptID)
	assert.Equal(t, event.Payload.FinalScore, decoded.Payload.FinalScore)
	assert.Equal(t, event.Payload.Status, decoded.Payload.Status)
}
