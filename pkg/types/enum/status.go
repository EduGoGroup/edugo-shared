package enum

// MaterialStatus representa el estado de un material educativo
type MaterialStatus string

const (
	MaterialStatusDraft     MaterialStatus = "draft"
	MaterialStatusPublished MaterialStatus = "published"
	MaterialStatusArchived  MaterialStatus = "archived"
)

// IsValid verifica si el status es válido
func (s MaterialStatus) IsValid() bool {
	switch s {
	case MaterialStatusDraft, MaterialStatusPublished, MaterialStatusArchived:
		return true
	}
	return false
}

// String retorna la representación en string del status
func (s MaterialStatus) String() string {
	return string(s)
}

// AllMaterialStatuses retorna todos los status válidos
func AllMaterialStatuses() []MaterialStatus {
	return []MaterialStatus{
		MaterialStatusDraft,
		MaterialStatusPublished,
		MaterialStatusArchived,
	}
}

// ProgressStatus representa el estado de progreso de lectura
type ProgressStatus string

const (
	ProgressStatusNotStarted ProgressStatus = "not_started"
	ProgressStatusInProgress ProgressStatus = "in_progress"
	ProgressStatusCompleted  ProgressStatus = "completed"
)

// IsValid verifica si el status es válido
func (p ProgressStatus) IsValid() bool {
	switch p {
	case ProgressStatusNotStarted, ProgressStatusInProgress, ProgressStatusCompleted:
		return true
	}
	return false
}

// String retorna la representación en string del status
func (p ProgressStatus) String() string {
	return string(p)
}

// AllProgressStatuses retorna todos los status de progreso válidos
func AllProgressStatuses() []ProgressStatus {
	return []ProgressStatus{
		ProgressStatusNotStarted,
		ProgressStatusInProgress,
		ProgressStatusCompleted,
	}
}

// ProcessingStatus representa el estado de procesamiento de un material
type ProcessingStatus string

const (
	ProcessingStatusPending    ProcessingStatus = "pending"
	ProcessingStatusProcessing ProcessingStatus = "processing"
	ProcessingStatusCompleted  ProcessingStatus = "completed"
	ProcessingStatusFailed     ProcessingStatus = "failed"
)

// IsValid verifica si el status es válido
func (p ProcessingStatus) IsValid() bool {
	switch p {
	case ProcessingStatusPending, ProcessingStatusProcessing, ProcessingStatusCompleted, ProcessingStatusFailed:
		return true
	}
	return false
}

// String retorna la representación en string del status
func (p ProcessingStatus) String() string {
	return string(p)
}

// AllProcessingStatuses retorna todos los status de procesamiento válidos
func AllProcessingStatuses() []ProcessingStatus {
	return []ProcessingStatus{
		ProcessingStatusPending,
		ProcessingStatusProcessing,
		ProcessingStatusCompleted,
		ProcessingStatusFailed,
	}
}
