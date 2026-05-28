package events

import (
	"errors"
	"time"
)

// AssessmentAssignedEvent representa la asignacion de una evaluacion a estudiantes o unidades.
type AssessmentAssignedEvent struct {
	EventID      string                    `json:"event_id"`
	EventType    string                    `json:"event_type"`
	EventVersion string                    `json:"event_version"`
	Timestamp    time.Time                 `json:"timestamp"`
	Payload      AssessmentAssignedPayload `json:"payload"`
}

// AssessmentAssignedPayload contiene los datos de la asignacion.
type AssessmentAssignedPayload struct {
	AssessmentID string `json:"assessment_id"`
	AssignmentID string `json:"assignment_id"`
	SchoolID     string `json:"school_id"`
	AssignedByID string `json:"assigned_by_id"`
	TargetType   string `json:"target_type"` // "student" o "unit"
	TargetID     string `json:"target_id"`
	Title        string `json:"title,omitempty"`
}

// NewAssessmentAssignedEvent crea y valida un nuevo evento de asignacion de evaluacion.
func NewAssessmentAssignedEvent(eventID, eventType, eventVersion string, payload AssessmentAssignedPayload) (AssessmentAssignedEvent, error) {
	if eventID == "" {
		return AssessmentAssignedEvent{}, errors.New("eventID no puede estar vacío")
	}
	if eventType == "" {
		return AssessmentAssignedEvent{}, errors.New("eventType no puede estar vacío")
	}
	if eventVersion == "" {
		return AssessmentAssignedEvent{}, errors.New("eventVersion no puede estar vacío")
	}

	if payload.AssessmentID == "" {
		return AssessmentAssignedEvent{}, errors.New("AssessmentID no puede estar vacío")
	}
	if payload.AssignmentID == "" {
		return AssessmentAssignedEvent{}, errors.New("AssignmentID no puede estar vacío")
	}
	if payload.SchoolID == "" {
		return AssessmentAssignedEvent{}, errors.New("SchoolID no puede estar vacío")
	}
	if payload.AssignedByID == "" {
		return AssessmentAssignedEvent{}, errors.New("AssignedByID no puede estar vacío")
	}
	if payload.TargetType == "" {
		return AssessmentAssignedEvent{}, errors.New("TargetType no puede estar vacío")
	}
	if payload.TargetID == "" {
		return AssessmentAssignedEvent{}, errors.New("TargetID no puede estar vacío")
	}

	return AssessmentAssignedEvent{
		EventID:      eventID,
		EventType:    eventType,
		EventVersion: eventVersion,
		Timestamp:    time.Now(),
		Payload:      payload,
	}, nil
}
