package events

import (
	"errors"
	"time"
)

// AssessmentReviewedEvent representa la revision de un intento de evaluacion por un profesor.
type AssessmentReviewedEvent struct {
	EventID      string                    `json:"event_id"`
	EventType    string                    `json:"event_type"`
	EventVersion string                    `json:"event_version"`
	Timestamp    time.Time                 `json:"timestamp"`
	Payload      AssessmentReviewedPayload `json:"payload"`
}

// AssessmentReviewedPayload contiene los datos de la revision.
type AssessmentReviewedPayload struct {
	AttemptID    string `json:"attempt_id"`
	AssessmentID string `json:"assessment_id"`
	ReviewerID   string `json:"reviewer_id"`
	SchoolID     string `json:"school_id"`
	// AcademicUnitID es la unidad académica activa del emisor (F4.6). Viaja para que
	// el push lleve unit_id y el deep-link multi-tenant cambie al contexto exacto
	// (escuela+unidad). omitempty: un evento viejo sin el campo deserializa a vacío.
	AcademicUnitID string  `json:"academic_unit_id,omitempty"`
	FinalScore     float64 `json:"final_score"`
	TotalPoints    float64 `json:"total_points"`
	Status         string  `json:"status"`
	StudentID      string  `json:"student_id,omitempty"`
	Title          string  `json:"title,omitempty"`
}

// NewAssessmentReviewedEvent crea y valida un nuevo evento de revision de evaluacion.
func NewAssessmentReviewedEvent(eventID, eventType, eventVersion string, payload AssessmentReviewedPayload) (AssessmentReviewedEvent, error) {
	if eventID == "" {
		return AssessmentReviewedEvent{}, errors.New("eventID no puede estar vacío")
	}
	if eventType == "" {
		return AssessmentReviewedEvent{}, errors.New("eventType no puede estar vacío")
	}
	if eventVersion == "" {
		return AssessmentReviewedEvent{}, errors.New("eventVersion no puede estar vacío")
	}

	if payload.AttemptID == "" {
		return AssessmentReviewedEvent{}, errors.New("AttemptID no puede estar vacío")
	}
	if payload.AssessmentID == "" {
		return AssessmentReviewedEvent{}, errors.New("AssessmentID no puede estar vacío")
	}
	if payload.ReviewerID == "" {
		return AssessmentReviewedEvent{}, errors.New("ReviewerID no puede estar vacío")
	}
	if payload.SchoolID == "" {
		return AssessmentReviewedEvent{}, errors.New("SchoolID no puede estar vacío")
	}
	if payload.Status == "" {
		return AssessmentReviewedEvent{}, errors.New("status no puede estar vacío")
	}

	return AssessmentReviewedEvent{
		EventID:      eventID,
		EventType:    eventType,
		EventVersion: eventVersion,
		Timestamp:    time.Now(),
		Payload:      payload,
	}, nil
}
