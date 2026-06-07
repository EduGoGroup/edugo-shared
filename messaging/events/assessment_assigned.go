package events

import (
	"errors"
	"time"
)

// AssessmentAssignedEvent representa la asignacion de una evaluacion a una
// sesion de materia (subject_offering). El targeting siempre es por oferta:
// los destinatarios son los alumnos inscritos en la sesion.
type AssessmentAssignedEvent struct {
	EventID      string                    `json:"event_id"`
	EventType    string                    `json:"event_type"`
	EventVersion string                    `json:"event_version"`
	Timestamp    time.Time                 `json:"timestamp"`
	Payload      AssessmentAssignedPayload `json:"payload"`
}

// AssessmentAssignedPayload contiene los datos de la asignacion sobre una sesion.
type AssessmentAssignedPayload struct {
	AssessmentID           string     `json:"assessment_id"`
	AssignmentID           string     `json:"assignment_id"`
	SchoolID               string     `json:"school_id"`
	AssignedByMembershipID string     `json:"assigned_by_membership_id"`
	SubjectOfferingID      string     `json:"subject_offering_id"`
	DueDate                *time.Time `json:"due_date,omitempty"`
	Title                  string     `json:"title,omitempty"`
}

// NewAssessmentAssignedEvent crea y valida un nuevo evento de asignacion de evaluacion.
// Campos requeridos: AssessmentID, AssignmentID, SchoolID, AssignedByMembershipID,
// SubjectOfferingID. DueDate y Title son opcionales.
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
	if payload.AssignedByMembershipID == "" {
		return AssessmentAssignedEvent{}, errors.New("AssignedByMembershipID no puede estar vacío")
	}
	if payload.SubjectOfferingID == "" {
		return AssessmentAssignedEvent{}, errors.New("SubjectOfferingID no puede estar vacío")
	}

	return AssessmentAssignedEvent{
		EventID:      eventID,
		EventType:    eventType,
		EventVersion: eventVersion,
		Timestamp:    time.Now(),
		Payload:      payload,
	}, nil
}
