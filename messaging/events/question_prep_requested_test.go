package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validQuestionPrepRequestedPayload() QuestionPrepRequestedPayload {
	return QuestionPrepRequestedPayload{
		QuestionID:   "question_001",
		AssessmentID: "assess_001",
		Reason:       QuestionPrepReasonCreated,
	}
}

func TestNewQuestionPrepRequestedEvent_Valid(t *testing.T) {
	payload := validQuestionPrepRequestedPayload()

	event, err := NewQuestionPrepRequestedEvent("evt_001", EventTypeQuestionPrepRequested, "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "evt_001", event.EventID)
	assert.Equal(t, "question.prep_requested", event.EventType)
	assert.Equal(t, EventTypeQuestionPrepRequested, event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, "question_001", event.Payload.QuestionID)
	assert.Equal(t, "assess_001", event.Payload.AssessmentID)
	assert.Equal(t, "created", event.Payload.Reason)
}

func TestNewQuestionPrepRequestedEvent_ValidReasons(t *testing.T) {
	reasons := []string{
		QuestionPrepReasonCreated,
		QuestionPrepReasonUpdated,
		QuestionPrepReasonFeedback,
		QuestionPrepReasonBackfill,
	}
	for _, reason := range reasons {
		t.Run(reason, func(t *testing.T) {
			payload := validQuestionPrepRequestedPayload()
			payload.Reason = reason

			event, err := NewQuestionPrepRequestedEvent("evt_1", EventTypeQuestionPrepRequested, "1.0", payload)

			require.NoError(t, err)
			assert.Equal(t, reason, event.Payload.Reason)
		})
	}
}

func TestNewQuestionPrepRequestedEvent_EmptyFields(t *testing.T) {
	tests := []struct {
		name         string
		eventID      string
		eventType    string
		eventVersion string
		payload      QuestionPrepRequestedPayload
		wantErr      string
	}{
		{
			name: "eventID vacio", eventID: "", eventType: EventTypeQuestionPrepRequested, eventVersion: "1.0",
			payload: validQuestionPrepRequestedPayload(), wantErr: "eventID",
		},
		{
			name: "eventType vacio", eventID: "evt_1", eventType: "", eventVersion: "1.0",
			payload: validQuestionPrepRequestedPayload(), wantErr: "eventType",
		},
		{
			name: "eventVersion vacio", eventID: "evt_1", eventType: EventTypeQuestionPrepRequested, eventVersion: "",
			payload: validQuestionPrepRequestedPayload(), wantErr: "eventVersion",
		},
		{
			name: "QuestionID vacio", eventID: "evt_1", eventType: EventTypeQuestionPrepRequested, eventVersion: "1.0",
			payload: func() QuestionPrepRequestedPayload {
				p := validQuestionPrepRequestedPayload()
				p.QuestionID = ""
				return p
			}(),
			wantErr: "QuestionID",
		},
		{
			name: "AssessmentID vacio", eventID: "evt_1", eventType: EventTypeQuestionPrepRequested, eventVersion: "1.0",
			payload: func() QuestionPrepRequestedPayload {
				p := validQuestionPrepRequestedPayload()
				p.AssessmentID = ""
				return p
			}(),
			wantErr: "AssessmentID",
		},
		{
			name: "Reason vacio", eventID: "evt_1", eventType: EventTypeQuestionPrepRequested, eventVersion: "1.0",
			payload: func() QuestionPrepRequestedPayload {
				p := validQuestionPrepRequestedPayload()
				p.Reason = ""
				return p
			}(),
			wantErr: "Reason",
		},
		{
			name: "Reason invalido", eventID: "evt_1", eventType: EventTypeQuestionPrepRequested, eventVersion: "1.0",
			payload: func() QuestionPrepRequestedPayload {
				p := validQuestionPrepRequestedPayload()
				p.Reason = "otro"
				return p
			}(),
			wantErr: "Reason",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewQuestionPrepRequestedEvent(tt.eventID, tt.eventType, tt.eventVersion, tt.payload)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestQuestionPrepRequestedEvent_Serialization(t *testing.T) {
	payload := validQuestionPrepRequestedPayload()

	event, err := NewQuestionPrepRequestedEvent("evt_001", EventTypeQuestionPrepRequested, "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded QuestionPrepRequestedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, event.EventID, decoded.EventID)
	assert.Equal(t, event.EventType, decoded.EventType)
	assert.Equal(t, event.Payload.QuestionID, decoded.Payload.QuestionID)
	assert.Equal(t, event.Payload.AssessmentID, decoded.Payload.AssessmentID)
	assert.Equal(t, event.Payload.Reason, decoded.Payload.Reason)
}
