package events

import (
	"errors"
	"time"
)

// EventTypeQuestionPrepRequested es el event_type del evento que dispara la
// preparación de una pregunta para el LLM (plan 042). Lo consumen learning
// (al publicar) y el worker (al enrutar) para no depender de literales sueltos.
const EventTypeQuestionPrepRequested = "question.prep_requested"

// Razones válidas por las que learning solicita (re)preparar una pregunta.
// El conjunto es cerrado: el worker las usa para decidir el carril de preparación
// (p. ej. relleno masivo vs edición puntual) y para observabilidad.
const (
	// QuestionPrepReasonCreated: la pregunta se acaba de crear.
	QuestionPrepReasonCreated = "created"
	// QuestionPrepReasonUpdated: cambió el enunciado/rúbrica de una pregunta existente.
	QuestionPrepReasonUpdated = "updated"
	// QuestionPrepReasonFeedback: el profesor dio feedback sobre la prep previa y se re-prepara.
	QuestionPrepReasonFeedback = "feedback"
	// QuestionPrepReasonBackfill: relleno de preguntas sin prep (migración/arranque del carril).
	QuestionPrepReasonBackfill = "backfill"
)

// QuestionPrepRequestedEvent señala que una pregunta short_answer/open_ended fue
// creada o editada y debe prepararse para el LLM (descomposición/enriquecimiento).
// Lo publica learning en cada create/update de pregunta y lo consume el worker
// (plan 042). Es mínimo a propósito: el worker lee la pregunta fresca por M2M; el
// evento NO transporta el contenido de la pregunta.
type QuestionPrepRequestedEvent struct {
	EventID      string                       `json:"event_id"`
	EventType    string                       `json:"event_type"`
	EventVersion string                       `json:"event_version"`
	Timestamp    time.Time                    `json:"timestamp"`
	Payload      QuestionPrepRequestedPayload `json:"payload"`
}

// QuestionPrepRequestedPayload identifica la pregunta a preparar y el motivo de la
// solicitud. Nada de contenido: solo referencias que el worker resuelve por M2M.
type QuestionPrepRequestedPayload struct {
	QuestionID   string `json:"question_id"`
	AssessmentID string `json:"assessment_id"`
	Reason       string `json:"reason"`
}

// NewQuestionPrepRequestedEvent crea y valida un nuevo evento de solicitud de
// preparación de pregunta.
func NewQuestionPrepRequestedEvent(eventID, eventType, eventVersion string, payload QuestionPrepRequestedPayload) (QuestionPrepRequestedEvent, error) {
	if eventID == "" {
		return QuestionPrepRequestedEvent{}, errors.New("eventID no puede estar vacío")
	}
	if eventType == "" {
		return QuestionPrepRequestedEvent{}, errors.New("eventType no puede estar vacío")
	}
	if eventVersion == "" {
		return QuestionPrepRequestedEvent{}, errors.New("eventVersion no puede estar vacío")
	}

	if payload.QuestionID == "" {
		return QuestionPrepRequestedEvent{}, errors.New("QuestionID no puede estar vacío")
	}
	if payload.AssessmentID == "" {
		return QuestionPrepRequestedEvent{}, errors.New("AssessmentID no puede estar vacío")
	}
	if !isValidQuestionPrepReason(payload.Reason) {
		return QuestionPrepRequestedEvent{}, errors.New("valor de Reason inválido: debe ser created|updated|feedback|backfill")
	}

	return QuestionPrepRequestedEvent{
		EventID:      eventID,
		EventType:    eventType,
		EventVersion: eventVersion,
		Timestamp:    time.Now(),
		Payload:      payload,
	}, nil
}

// isValidQuestionPrepReason valida que la razón sea uno de los valores del conjunto
// cerrado. También rechaza la cadena vacía.
func isValidQuestionPrepReason(reason string) bool {
	switch reason {
	case QuestionPrepReasonCreated, QuestionPrepReasonUpdated, QuestionPrepReasonFeedback, QuestionPrepReasonBackfill:
		return true
	default:
		return false
	}
}
