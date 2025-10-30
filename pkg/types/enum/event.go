package enum

// EventType representa los tipos de eventos del sistema
type EventType string

const (
	// Eventos de materiales
	EventMaterialUploaded   EventType = "material.uploaded"
	EventMaterialReprocess  EventType = "material.reprocess"
	EventMaterialDeleted    EventType = "material.deleted"
	EventMaterialPublished  EventType = "material.published"
	EventMaterialArchived   EventType = "material.archived"

	// Eventos de assessments
	EventAssessmentAttemptRecorded EventType = "assessment.attempt_recorded"
	EventAssessmentCompleted       EventType = "assessment.completed"

	// Eventos de estudiantes
	EventStudentEnrolled EventType = "student.enrolled"
	EventStudentProgress EventType = "student.progress"

	// Eventos de usuarios
	EventUserCreated      EventType = "user.created"
	EventUserUpdated      EventType = "user.updated"
	EventUserDeactivated  EventType = "user.deactivated"
)

// IsValid verifica si el tipo de evento es válido
func (e EventType) IsValid() bool {
	switch e {
	case EventMaterialUploaded, EventMaterialReprocess, EventMaterialDeleted,
		EventMaterialPublished, EventMaterialArchived,
		EventAssessmentAttemptRecorded, EventAssessmentCompleted,
		EventStudentEnrolled, EventStudentProgress,
		EventUserCreated, EventUserUpdated, EventUserDeactivated:
		return true
	}
	return false
}

// String retorna la representación en string del evento
func (e EventType) String() string {
	return string(e)
}

// GetRoutingKey retorna la routing key para RabbitMQ
func (e EventType) GetRoutingKey() string {
	return string(e)
}

// AllEventTypes retorna todos los tipos de eventos válidos
func AllEventTypes() []EventType {
	return []EventType{
		EventMaterialUploaded,
		EventMaterialReprocess,
		EventMaterialDeleted,
		EventMaterialPublished,
		EventMaterialArchived,
		EventAssessmentAttemptRecorded,
		EventAssessmentCompleted,
		EventStudentEnrolled,
		EventStudentProgress,
		EventUserCreated,
		EventUserUpdated,
		EventUserDeactivated,
	}
}
