package events

import (
	"errors"
	"time"
)

// EventTypeMaterialAssessmentRequested es el event_type del evento que dispara la
// generación de una evaluación a partir de un material (plan 043). Lo consumen
// learning (al solicitar) y el worker (al enrutar) para no depender de literales sueltos.
const EventTypeMaterialAssessmentRequested = "material.assessment_requested"

// MaterialAssessmentRequestedEvent señala que un material debe procesarse para
// generar una evaluación mediante el LLM (pipeline material→evaluación).
// Lo publica learning al solicitar la generación y lo consume el worker (plan 043).
// Es mínimo a propósito: el worker lee el material fresco por M2M; el evento NO
// transporta el contenido del material.
type MaterialAssessmentRequestedEvent struct {
	EventID      string                             `json:"event_id"`
	EventType    string                             `json:"event_type"`
	EventVersion string                             `json:"event_version"`
	Timestamp    time.Time                          `json:"timestamp"`
	Payload      MaterialAssessmentRequestedPayload `json:"payload"`
}

// MaterialAssessmentRequestedPayload identifica el job de generación y el material
// de origen. Nada de contenido: solo referencias que el worker resuelve por M2M.
type MaterialAssessmentRequestedPayload struct {
	JobID      string `json:"job_id"`
	MaterialID string `json:"material_id"`
	SchoolID   string `json:"school_id"`
}

// NewMaterialAssessmentRequestedEvent crea y valida un nuevo evento de solicitud de
// generación de evaluación a partir de un material.
func NewMaterialAssessmentRequestedEvent(eventID, eventType, eventVersion string, payload MaterialAssessmentRequestedPayload) (MaterialAssessmentRequestedEvent, error) {
	if eventID == "" {
		return MaterialAssessmentRequestedEvent{}, errors.New("eventID no puede estar vacío")
	}
	if eventType == "" {
		return MaterialAssessmentRequestedEvent{}, errors.New("eventType no puede estar vacío")
	}
	if eventVersion == "" {
		return MaterialAssessmentRequestedEvent{}, errors.New("eventVersion no puede estar vacío")
	}

	if payload.JobID == "" {
		return MaterialAssessmentRequestedEvent{}, errors.New("JobID no puede estar vacío")
	}
	if payload.MaterialID == "" {
		return MaterialAssessmentRequestedEvent{}, errors.New("MaterialID no puede estar vacío")
	}
	if payload.SchoolID == "" {
		return MaterialAssessmentRequestedEvent{}, errors.New("SchoolID no puede estar vacío")
	}

	return MaterialAssessmentRequestedEvent{
		EventID:      eventID,
		EventType:    eventType,
		EventVersion: eventVersion,
		Timestamp:    time.Now(),
		Payload:      payload,
	}, nil
}
