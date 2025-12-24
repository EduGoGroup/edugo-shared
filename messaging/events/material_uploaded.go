package events

import (
	"time"
)

// MaterialUploadedEvent representa un evento de carga de material siguiendo el estándar CloudEvents.
// Se utiliza para comunicar entre servicios cuando un profesor sube un nuevo material educativo.
//
// El evento incluye metadatos CloudEvents estándar (EventID, EventType, EventVersion, Timestamp)
// y un payload específico del dominio con información del material subido.
//
// Ejemplo de uso:
//
//	event := MaterialUploadedEvent{
//		EventID:      "evt-123",
//		EventType:    "material.uploaded",
//		EventVersion: "1.0",
//		Timestamp:    time.Now(),
//		Payload: MaterialUploadedPayload{
//			MaterialID:    "mat-456",
//			SchoolID:      "school-789",
//			TeacherID:     "teacher-001",
//			FileURL:       "https://s3.amazonaws.com/bucket/file.pdf",
//			FileSizeBytes: 1024,
//			FileType:      "application/pdf",
//		},
//	}
type MaterialUploadedEvent struct {
	EventID      string                  `json:"event_id"`
	EventType    string                  `json:"event_type"`
	EventVersion string                  `json:"event_version"`
	Timestamp    time.Time               `json:"timestamp"`
	Payload      MaterialUploadedPayload `json:"payload"`
}

// MaterialUploadedPayload contiene los datos específicos del evento de material subido.
//
// Campos principales:
//   - MaterialID: Identificador único del material en el sistema
//   - SchoolID: Identificador de la escuela a la que pertenece el material
//   - TeacherID: Identificador del profesor que subió el material
//   - FileURL: URL completa donde está almacenado el archivo
//   - FileSizeBytes: Tamaño del archivo en bytes
//   - FileType: Tipo MIME del archivo (ej: "application/pdf", "image/png")
//   - Metadata: Información adicional opcional (ej: s3_key, checksums, tags)
//
// El campo Metadata se usa para información variable que no está en el esquema base.
// Por ejemplo, para almacenar la clave S3 separada de la URL completa, o metadatos
// específicos del tipo de archivo como dimensiones de imagen, duración de video, etc.
type MaterialUploadedPayload struct {
	MaterialID    string                 `json:"material_id"`
	SchoolID      string                 `json:"school_id"`
	TeacherID     string                 `json:"teacher_id"`
	FileURL       string                 `json:"file_url"`
	FileSizeBytes int64                  `json:"file_size_bytes"`
	FileType      string                 `json:"file_type"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// GetMaterialID retorna el ID del material.
// Este es un método de compatibilidad para sistemas legacy que esperan esta interfaz.
//
// Retorna:
//   - El identificador único del material (MaterialID del Payload)
func (e MaterialUploadedEvent) GetMaterialID() string {
	return e.Payload.MaterialID
}

// GetS3Key retorna la clave S3 del archivo almacenado.
//
// Primero intenta extraer el valor de Metadata["s3_key"].
// Si no está disponible, retorna FileURL completo como fallback.
// Este es un método de compatibilidad para sistemas legacy.
//
// Retorna:
//   - La clave S3 si está en Metadata["s3_key"]
//   - FileURL completo si Metadata es nil o no contiene "s3_key"
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

// GetAuthorID retorna el ID del autor/profesor que subió el material.
// Este es un método de compatibilidad para sistemas legacy que esperan esta interfaz.
//
// Retorna:
//   - El identificador del profesor (TeacherID del Payload)
func (e MaterialUploadedEvent) GetAuthorID() string {
	return e.Payload.TeacherID
}
