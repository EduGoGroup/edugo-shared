package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validAttemptReviewRequestedPayload() AttemptReviewRequestedPayload {
	return AttemptReviewRequestedPayload{
		AttemptID:    "attempt_001",
		AssessmentID: "assess_001",
		SchoolID:     "school_001",
		Answers: []AttemptReviewAnswerRef{
			{AnswerID: "ans_001", QuestionType: "open_ended"},
			{AnswerID: "ans_002", QuestionType: "short_answer"},
		},
	}
}

func TestNewAttemptReviewRequestedEvent_Valid(t *testing.T) {
	payload := validAttemptReviewRequestedPayload()

	event, err := NewAttemptReviewRequestedEvent("evt_001", "attempt.review_requested", "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "evt_001", event.EventID)
	assert.Equal(t, "attempt.review_requested", event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, "attempt_001", event.Payload.AttemptID)
	assert.Len(t, event.Payload.Answers, 2)
	assert.Equal(t, "open_ended", event.Payload.Answers[0].QuestionType)
}

func TestNewAttemptReviewRequestedEvent_EmptyFields(t *testing.T) {
	tests := []struct {
		name         string
		eventID      string
		eventType    string
		eventVersion string
		payload      AttemptReviewRequestedPayload
		wantErr      string
	}{
		{
			name: "eventID vacio", eventID: "", eventType: "attempt.review_requested", eventVersion: "1.0",
			payload: validAttemptReviewRequestedPayload(), wantErr: "eventID",
		},
		{
			name: "eventType vacio", eventID: "evt_1", eventType: "", eventVersion: "1.0",
			payload: validAttemptReviewRequestedPayload(), wantErr: "eventType",
		},
		{
			name: "eventVersion vacio", eventID: "evt_1", eventType: "attempt.review_requested", eventVersion: "",
			payload: validAttemptReviewRequestedPayload(), wantErr: "eventVersion",
		},
		{
			name: "AttemptID vacio", eventID: "evt_1", eventType: "attempt.review_requested", eventVersion: "1.0",
			payload: func() AttemptReviewRequestedPayload {
				p := validAttemptReviewRequestedPayload()
				p.AttemptID = ""
				return p
			}(),
			wantErr: "AttemptID",
		},
		{
			name: "AssessmentID vacio", eventID: "evt_1", eventType: "attempt.review_requested", eventVersion: "1.0",
			payload: func() AttemptReviewRequestedPayload {
				p := validAttemptReviewRequestedPayload()
				p.AssessmentID = ""
				return p
			}(),
			wantErr: "AssessmentID",
		},
		{
			name: "SchoolID vacio", eventID: "evt_1", eventType: "attempt.review_requested", eventVersion: "1.0",
			payload: func() AttemptReviewRequestedPayload {
				p := validAttemptReviewRequestedPayload()
				p.SchoolID = ""
				return p
			}(),
			wantErr: "SchoolID",
		},
		{
			name: "Answers vacio", eventID: "evt_1", eventType: "attempt.review_requested", eventVersion: "1.0",
			payload: func() AttemptReviewRequestedPayload {
				p := validAttemptReviewRequestedPayload()
				p.Answers = nil
				return p
			}(),
			wantErr: "answers",
		},
		{
			name: "AnswerID vacio", eventID: "evt_1", eventType: "attempt.review_requested", eventVersion: "1.0",
			payload: func() AttemptReviewRequestedPayload {
				p := validAttemptReviewRequestedPayload()
				p.Answers = []AttemptReviewAnswerRef{{AnswerID: "", QuestionType: "open_ended"}}
				return p
			}(),
			wantErr: "AnswerID",
		},
		{
			name: "QuestionType vacio", eventID: "evt_1", eventType: "attempt.review_requested", eventVersion: "1.0",
			payload: func() AttemptReviewRequestedPayload {
				p := validAttemptReviewRequestedPayload()
				p.Answers = []AttemptReviewAnswerRef{{AnswerID: "ans_001", QuestionType: ""}}
				return p
			}(),
			wantErr: "QuestionType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAttemptReviewRequestedEvent(tt.eventID, tt.eventType, tt.eventVersion, tt.payload)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestAttemptReviewRequestedEvent_Serialization(t *testing.T) {
	payload := validAttemptReviewRequestedPayload()

	event, err := NewAttemptReviewRequestedEvent("evt_001", "attempt.review_requested", "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded AttemptReviewRequestedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, event.EventID, decoded.EventID)
	assert.Equal(t, event.Payload.AttemptID, decoded.Payload.AttemptID)
	assert.Equal(t, event.Payload.SchoolID, decoded.Payload.SchoolID)
	require.Len(t, decoded.Payload.Answers, 2)
	assert.Equal(t, "ans_002", decoded.Payload.Answers[1].AnswerID)
	assert.Equal(t, "short_answer", decoded.Payload.Answers[1].QuestionType)
}
