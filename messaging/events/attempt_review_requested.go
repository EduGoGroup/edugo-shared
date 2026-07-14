package events

import (
	"errors"
	"time"
)

// AttemptReviewRequestedEvent señala que un intento entregado quedó pendiente de
// revisión humana (o asistida por LLM): tiene ≥1 respuesta open_ended sin calificar
// o short_answer que falló el matching fuzzy. Lo publica learning al hacer submit y
// lo consume el worker para arrancar la revisión asistida (planes 039/040).
type AttemptReviewRequestedEvent struct {
	EventID      string                        `json:"event_id"`
	EventType    string                        `json:"event_type"`
	EventVersion string                        `json:"event_version"`
	Timestamp    time.Time                     `json:"timestamp"`
	Payload      AttemptReviewRequestedPayload `json:"payload"`
}

// AttemptReviewRequestedPayload contiene el intento y la lista de respuestas que
// requieren revisión.
type AttemptReviewRequestedPayload struct {
	AttemptID    string                   `json:"attempt_id"`
	AssessmentID string                   `json:"assessment_id"`
	SchoolID     string                   `json:"school_id"`
	Answers      []AttemptReviewAnswerRef `json:"answers"`
}

// AttemptReviewAnswerRef identifica una respuesta a revisar y su tipo de pregunta,
// para que el consumidor decida el carril de revisión (open_ended vs short_answer).
type AttemptReviewAnswerRef struct {
	AnswerID     string `json:"answer_id"`
	QuestionType string `json:"question_type"`
}

// NewAttemptReviewRequestedEvent crea y valida un nuevo evento de solicitud de
// revisión de intento.
func NewAttemptReviewRequestedEvent(eventID, eventType, eventVersion string, payload AttemptReviewRequestedPayload) (AttemptReviewRequestedEvent, error) {
	if eventID == "" {
		return AttemptReviewRequestedEvent{}, errors.New("eventID no puede estar vacío")
	}
	if eventType == "" {
		return AttemptReviewRequestedEvent{}, errors.New("eventType no puede estar vacío")
	}
	if eventVersion == "" {
		return AttemptReviewRequestedEvent{}, errors.New("eventVersion no puede estar vacío")
	}

	if payload.AttemptID == "" {
		return AttemptReviewRequestedEvent{}, errors.New("AttemptID no puede estar vacío")
	}
	if payload.AssessmentID == "" {
		return AttemptReviewRequestedEvent{}, errors.New("AssessmentID no puede estar vacío")
	}
	if payload.SchoolID == "" {
		return AttemptReviewRequestedEvent{}, errors.New("SchoolID no puede estar vacío")
	}
	if len(payload.Answers) == 0 {
		return AttemptReviewRequestedEvent{}, errors.New("Answers no puede estar vacío")
	}
	for _, a := range payload.Answers {
		if a.AnswerID == "" {
			return AttemptReviewRequestedEvent{}, errors.New("AnswerID no puede estar vacío")
		}
		if a.QuestionType == "" {
			return AttemptReviewRequestedEvent{}, errors.New("QuestionType no puede estar vacío")
		}
	}

	return AttemptReviewRequestedEvent{
		EventID:      eventID,
		EventType:    eventType,
		EventVersion: eventVersion,
		Timestamp:    time.Now(),
		Payload:      payload,
	}, nil
}
