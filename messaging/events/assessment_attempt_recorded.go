package events

import (
	"errors"
	"time"
)

// AssessmentAttemptRecordedEvent representa el registro de un intento de evaluacion por un estudiante.
type AssessmentAttemptRecordedEvent struct {
	EventID      string                            `json:"event_id"`
	EventType    string                            `json:"event_type"`
	EventVersion string                            `json:"event_version"`
	Timestamp    time.Time                         `json:"timestamp"`
	Payload      AssessmentAttemptRecordedPayload  `json:"payload"`
}

// AssessmentAttemptRecordedPayload contiene los datos del intento registrado.
type AssessmentAttemptRecordedPayload struct {
	AttemptID    string    `json:"attempt_id"`
	AssessmentID string    `json:"assessment_id"`
	StudentID    string    `json:"student_id"`
	SchoolID     string    `json:"school_id"`
	Score        float64   `json:"score"`
	TotalPoints  float64   `json:"total_points"`
	SubmittedAt  time.Time `json:"submitted_at"`
}

// NewAssessmentAttemptRecordedEvent crea y valida un nuevo evento de intento de evaluacion.
func NewAssessmentAttemptRecordedEvent(eventID, eventType, eventVersion string, payload AssessmentAttemptRecordedPayload) (AssessmentAttemptRecordedEvent, error) {
	if eventID == "" {
		return AssessmentAttemptRecordedEvent{}, errors.New("eventID no puede estar vacío")
	}
	if eventType == "" {
		return AssessmentAttemptRecordedEvent{}, errors.New("eventType no puede estar vacío")
	}
	if eventVersion == "" {
		return AssessmentAttemptRecordedEvent{}, errors.New("eventVersion no puede estar vacío")
	}

	if payload.AttemptID == "" {
		return AssessmentAttemptRecordedEvent{}, errors.New("AttemptID no puede estar vacío")
	}
	if payload.AssessmentID == "" {
		return AssessmentAttemptRecordedEvent{}, errors.New("AssessmentID no puede estar vacío")
	}
	if payload.StudentID == "" {
		return AssessmentAttemptRecordedEvent{}, errors.New("StudentID no puede estar vacío")
	}
	if payload.SchoolID == "" {
		return AssessmentAttemptRecordedEvent{}, errors.New("SchoolID no puede estar vacío")
	}

	return AssessmentAttemptRecordedEvent{
		EventID:      eventID,
		EventType:    eventType,
		EventVersion: eventVersion,
		Timestamp:    time.Now(),
		Payload:      payload,
	}, nil
}
