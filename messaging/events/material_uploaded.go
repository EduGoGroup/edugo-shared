package events

import (
	"time"
)

// MaterialUploadedEvent representa el evento de material subido siguiendo CloudEvents estándar
type MaterialUploadedEvent struct {
	EventID      string                  `json:"event_id"`
	EventType    string                  `json:"event_type"`
	EventVersion string                  `json:"event_version"`
	Timestamp    time.Time               `json:"timestamp"`
	Payload      MaterialUploadedPayload `json:"payload"`
}

// MaterialUploadedPayload contiene los datos específicos del evento de material subido
type MaterialUploadedPayload struct {
	MaterialID    string                 `json:"material_id"`
	SchoolID      string                 `json:"school_id"`
	TeacherID     string                 `json:"teacher_id"`
	FileURL       string                 `json:"file_url"`
	FileSizeBytes int64                  `json:"file_size_bytes"`
	FileType      string                 `json:"file_type"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// GetMaterialID retorna el ID del material (método de compatibilidad)
func (e MaterialUploadedEvent) GetMaterialID() string {
	return e.Payload.MaterialID
}

// GetS3Key retorna la S3 key del archivo (método de compatibilidad)
// Extrae de FileURL o Metadata según el formato
func (e MaterialUploadedEvent) GetS3Key() string {
	// Intentar extraer de metadata primero
	if e.Payload.Metadata != nil {
		if key, ok := e.Payload.Metadata["s3_key"].(string); ok {
			return key
		}
	}
	
	// Fallback: retornar FileURL completo
	// En producción, aquí se podría parsear la URL para extraer solo la key
	return e.Payload.FileURL
}

// GetAuthorID retorna el ID del autor/profesor (método de compatibilidad)
func (e MaterialUploadedEvent) GetAuthorID() string {
	return e.Payload.TeacherID
}
