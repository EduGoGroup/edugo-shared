package events

import (
	"errors"
	"time"
)

// AssessmentPublishedEvent representa la publicacion de una evaluacion.
type AssessmentPublishedEvent struct {
	EventID      string                     `json:"event_id"`
	EventType    string                     `json:"event_type"`
	EventVersion string                     `json:"event_version"`
	Timestamp    time.Time                  `json:"timestamp"`
	Payload      AssessmentPublishedPayload `json:"payload"`
}

// AssessmentPublishedPayload contiene los datos de la evaluacion publicada.
type AssessmentPublishedPayload struct {
	AssessmentID  string `json:"assessment_id"`
	SchoolID      string `json:"school_id"`
	TeacherID     string `json:"teacher_id"`
	Title         string `json:"title"`
	QuestionCount int    `json:"question_count"`
}

// NewAssessmentPublishedEvent crea y valida un nuevo evento de evaluacion publicada.
func NewAssessmentPublishedEvent(eventID, eventType, eventVersion string, payload AssessmentPublishedPayload) (AssessmentPublishedEvent, error) {
	if eventID == "" {
		return AssessmentPublishedEvent{}, errors.New("eventID no puede estar vacío")
	}
	if eventType == "" {
		return AssessmentPublishedEvent{}, errors.New("eventType no puede estar vacío")
	}
	if eventVersion == "" {
		return AssessmentPublishedEvent{}, errors.New("eventVersion no puede estar vacío")
	}

	if payload.AssessmentID == "" {
		return AssessmentPublishedEvent{}, errors.New("AssessmentID no puede estar vacío")
	}
	if payload.SchoolID == "" {
		return AssessmentPublishedEvent{}, errors.New("SchoolID no puede estar vacío")
	}
	if payload.TeacherID == "" {
		return AssessmentPublishedEvent{}, errors.New("TeacherID no puede estar vacío")
	}
	if payload.Title == "" {
		return AssessmentPublishedEvent{}, errors.New("title no puede estar vacío")
	}

	return AssessmentPublishedEvent{
		EventID:      eventID,
		EventType:    eventType,
		EventVersion: eventVersion,
		Timestamp:    time.Now(),
		Payload:      payload,
	}, nil
}
