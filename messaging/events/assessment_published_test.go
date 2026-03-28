package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAssessmentPublishedEvent_Valid(t *testing.T) {
	payload := AssessmentPublishedPayload{
		AssessmentID:  "assess_001",
		SchoolID:      "school_001",
		TeacherID:     "teacher_001",
		Title:         "Examen Final Matematicas",
		QuestionCount: 20,
	}

	event, err := NewAssessmentPublishedEvent("evt_001", "assessment.published", "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "evt_001", event.EventID)
	assert.Equal(t, "assessment.published", event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, "assess_001", event.Payload.AssessmentID)
	assert.Equal(t, 20, event.Payload.QuestionCount)
}

func TestNewAssessmentPublishedEvent_EmptyFields(t *testing.T) {
	tests := []struct {
		name    string
		eventID string
		payload AssessmentPublishedPayload
		wantErr string
	}{
		{
			name:    "eventID vacio",
			eventID: "",
			payload: AssessmentPublishedPayload{
				AssessmentID: "a", SchoolID: "s", TeacherID: "t", Title: "T", QuestionCount: 1,
			},
			wantErr: "eventID",
		},
		{
			name:    "AssessmentID vacio",
			eventID: "evt_1",
			payload: AssessmentPublishedPayload{
				AssessmentID: "", SchoolID: "s", TeacherID: "t", Title: "T", QuestionCount: 1,
			},
			wantErr: "AssessmentID",
		},
		{
			name:    "SchoolID vacio",
			eventID: "evt_1",
			payload: AssessmentPublishedPayload{
				AssessmentID: "a", SchoolID: "", TeacherID: "t", Title: "T", QuestionCount: 1,
			},
			wantErr: "SchoolID",
		},
		{
			name:    "TeacherID vacio",
			eventID: "evt_1",
			payload: AssessmentPublishedPayload{
				AssessmentID: "a", SchoolID: "s", TeacherID: "", Title: "T", QuestionCount: 1,
			},
			wantErr: "TeacherID",
		},
		{
			name:    "Title vacio",
			eventID: "evt_1",
			payload: AssessmentPublishedPayload{
				AssessmentID: "a", SchoolID: "s", TeacherID: "t", Title: "", QuestionCount: 1,
			},
			wantErr: "Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAssessmentPublishedEvent(tt.eventID, "assessment.published", "1.0", tt.payload)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestAssessmentPublishedEvent_Serialization(t *testing.T) {
	payload := AssessmentPublishedPayload{
		AssessmentID:  "assess_001",
		SchoolID:      "school_001",
		TeacherID:     "teacher_001",
		Title:         "Examen Final",
		QuestionCount: 15,
	}

	event, err := NewAssessmentPublishedEvent("evt_001", "assessment.published", "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded AssessmentPublishedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, event.EventID, decoded.EventID)
	assert.Equal(t, event.Payload.AssessmentID, decoded.Payload.AssessmentID)
	assert.Equal(t, event.Payload.Title, decoded.Payload.Title)
	assert.Equal(t, event.Payload.QuestionCount, decoded.Payload.QuestionCount)
}
