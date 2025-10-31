package enum

// EventType representa los tipos de eventos del sistema
type EventType string

const (
	// EventMaterialUploaded represents a material upload event
	EventMaterialUploaded EventType = "material.uploaded"
	// EventMaterialReprocess represents a material reprocessing event
	EventMaterialReprocess EventType = "material.reprocess"
	// EventMaterialDeleted represents a material deletion event
	EventMaterialDeleted EventType = "material.deleted"
	// EventMaterialPublished represents a material publishing event
	EventMaterialPublished EventType = "material.published"
	// EventMaterialArchived represents a material archival event
	EventMaterialArchived EventType = "material.archived"

	// EventAssessmentAttemptRecorded represents an assessment attempt recording event
	EventAssessmentAttemptRecorded EventType = "assessment.attempt_recorded"
	// EventAssessmentCompleted represents an assessment completion event
	EventAssessmentCompleted EventType = "assessment.completed"

	// EventStudentEnrolled represents a student enrollment event
	EventStudentEnrolled EventType = "student.enrolled"
	// EventStudentProgress represents a student progress event
	EventStudentProgress EventType = "student.progress"

	// EventUserCreated represents a user creation event
	EventUserCreated EventType = "user.created"
	// EventUserUpdated represents a user update event
	EventUserUpdated EventType = "user.updated"
	// EventUserDeactivated represents a user deactivation event
	EventUserDeactivated EventType = "user.deactivated"
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
