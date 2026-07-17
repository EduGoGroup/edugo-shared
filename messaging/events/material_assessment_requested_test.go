package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validMaterialAssessmentRequestedPayload() MaterialAssessmentRequestedPayload {
	return MaterialAssessmentRequestedPayload{
		JobID:      "job_001",
		MaterialID: "material_001",
		SchoolID:   "school_001",
	}
}

func TestNewMaterialAssessmentRequestedEvent_Valid(t *testing.T) {
	payload := validMaterialAssessmentRequestedPayload()

	event, err := NewMaterialAssessmentRequestedEvent("evt_001", EventTypeMaterialAssessmentRequested, "1.0", payload)

	require.NoError(t, err)
	assert.Equal(t, "evt_001", event.EventID)
	assert.Equal(t, "material.assessment_requested", event.EventType)
	assert.Equal(t, EventTypeMaterialAssessmentRequested, event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, "job_001", event.Payload.JobID)
	assert.Equal(t, "material_001", event.Payload.MaterialID)
	assert.Equal(t, "school_001", event.Payload.SchoolID)
}

func TestNewMaterialAssessmentRequestedEvent_EmptyFields(t *testing.T) {
	tests := []struct {
		name         string
		eventID      string
		eventType    string
		eventVersion string
		payload      MaterialAssessmentRequestedPayload
		wantErr      string
	}{
		{
			name: "eventID vacio", eventID: "", eventType: EventTypeMaterialAssessmentRequested, eventVersion: "1.0",
			payload: validMaterialAssessmentRequestedPayload(), wantErr: "eventID",
		},
		{
			name: "eventType vacio", eventID: "evt_1", eventType: "", eventVersion: "1.0",
			payload: validMaterialAssessmentRequestedPayload(), wantErr: "eventType",
		},
		{
			name: "eventVersion vacio", eventID: "evt_1", eventType: EventTypeMaterialAssessmentRequested, eventVersion: "",
			payload: validMaterialAssessmentRequestedPayload(), wantErr: "eventVersion",
		},
		{
			name: "JobID vacio", eventID: "evt_1", eventType: EventTypeMaterialAssessmentRequested, eventVersion: "1.0",
			payload: func() MaterialAssessmentRequestedPayload {
				p := validMaterialAssessmentRequestedPayload()
				p.JobID = ""
				return p
			}(),
			wantErr: "JobID",
		},
		{
			name: "MaterialID vacio", eventID: "evt_1", eventType: EventTypeMaterialAssessmentRequested, eventVersion: "1.0",
			payload: func() MaterialAssessmentRequestedPayload {
				p := validMaterialAssessmentRequestedPayload()
				p.MaterialID = ""
				return p
			}(),
			wantErr: "MaterialID",
		},
		{
			name: "SchoolID vacio", eventID: "evt_1", eventType: EventTypeMaterialAssessmentRequested, eventVersion: "1.0",
			payload: func() MaterialAssessmentRequestedPayload {
				p := validMaterialAssessmentRequestedPayload()
				p.SchoolID = ""
				return p
			}(),
			wantErr: "SchoolID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMaterialAssessmentRequestedEvent(tt.eventID, tt.eventType, tt.eventVersion, tt.payload)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestMaterialAssessmentRequestedEvent_Serialization(t *testing.T) {
	payload := validMaterialAssessmentRequestedPayload()

	event, err := NewMaterialAssessmentRequestedEvent("evt_001", EventTypeMaterialAssessmentRequested, "1.0", payload)
	require.NoError(t, err)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded MaterialAssessmentRequestedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, event.EventID, decoded.EventID)
	assert.Equal(t, event.EventType, decoded.EventType)
	assert.Equal(t, event.Payload.JobID, decoded.Payload.JobID)
	assert.Equal(t, event.Payload.MaterialID, decoded.Payload.MaterialID)
	assert.Equal(t, event.Payload.SchoolID, decoded.Payload.SchoolID)
}
