package events

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

// MaterialUploadedEvent representa un evento de carga de material siguiendo el estándar CloudEvents.
// Se utiliza para comunicar entre servicios cuando un profesor sube un nuevo material educativo.
//
// El evento incluye metadatos CloudEvents estándar (EventID, EventType, EventVersion, Timestamp)
// y un payload específico del dominio con información del material subido.
//
// IMPORTANTE: Use NewMaterialUploadedEvent para crear instancias con validación automática.
//
// Ejemplo de uso:
//
//	event, err := NewMaterialUploadedEvent(
//		"evt-123",
//		"material.uploaded",
//		"1.0",
//		MaterialUploadedPayload{
//			MaterialID:    "mat-456",
//			SchoolID:      "school-789",
//			TeacherID:     "teacher-001",
//			FileURL:       "https://s3.amazonaws.com/bucket/file.pdf",
//			FileSizeBytes: 1024,
//			FileType:      "application/pdf",
//		},
//	)
//	if err != nil {
//		// Manejar error de validación
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
//   - MaterialID: Identificador único del material en el sistema (requerido)
//   - SchoolID: Identificador de la escuela a la que pertenece el material (requerido)
//   - TeacherID: Identificador del profesor que subió el material (requerido)
//   - FileURL: URL completa donde está almacenado el archivo (requerido, debe ser URL válida)
//   - FileSizeBytes: Tamaño del archivo en bytes (debe ser >= 0)
//   - FileType: Tipo MIME del archivo (ej: "application/pdf", "image/png") (requerido)
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
	FileSizeBytes uint64                 `json:"file_size_bytes"` // uint64 previene valores negativos
	FileType      string                 `json:"file_type"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// NewMaterialUploadedEvent crea y valida un nuevo evento de material subido.
//
// Esta función constructora valida que todos los campos requeridos estén presentes
// y que FileURL sea una URL válida. Esto previene la creación de eventos inválidos
// que podrían causar problemas en servicios consumidores.
//
// Parámetros:
//   - eventID: Identificador único del evento (requerido)
//   - eventType: Tipo de evento (requerido)
//   - eventVersion: Versión del esquema del evento (requerido)
//   - payload: Datos del material subido (se validan campos requeridos)
//
// Retorna:
//   - MaterialUploadedEvent validado con Timestamp establecido a time.Now()
//   - error si alguna validación falla
//
// Ejemplo:
//
//	event, err := NewMaterialUploadedEvent("evt-123", "material.uploaded", "1.0", payload)
//	if err != nil {
//		log.Fatalf("evento inválido: %v", err)
//	}
func NewMaterialUploadedEvent(eventID, eventType, eventVersion string, payload MaterialUploadedPayload) (MaterialUploadedEvent, error) {
	// Validar campos del evento
	if eventID == "" {
		return MaterialUploadedEvent{}, errors.New("eventID no puede estar vacío")
	}
	if eventType == "" {
		return MaterialUploadedEvent{}, errors.New("eventType no puede estar vacío")
	}
	if eventVersion == "" {
		return MaterialUploadedEvent{}, errors.New("eventVersion no puede estar vacío")
	}

	// Validar campos requeridos del payload
	if payload.MaterialID == "" {
		return MaterialUploadedEvent{}, errors.New("MaterialID no puede estar vacío")
	}
	if payload.SchoolID == "" {
		return MaterialUploadedEvent{}, errors.New("SchoolID no puede estar vacío")
	}
	if payload.TeacherID == "" {
		return MaterialUploadedEvent{}, errors.New("TeacherID no puede estar vacío")
	}
	if payload.FileURL == "" {
		return MaterialUploadedEvent{}, errors.New("FileURL no puede estar vacío")
	}
	if payload.FileType == "" {
		return MaterialUploadedEvent{}, errors.New("FileType no puede estar vacío")
	}

	// Validar que FileURL sea una URL válida
	if _, err := url.Parse(payload.FileURL); err != nil {
		return MaterialUploadedEvent{}, errors.New("FileURL no es una URL válida: " + err.Error())
	}

	// Crear evento con timestamp actual
	return MaterialUploadedEvent{
		EventID:      eventID,
		EventType:    eventType,
		EventVersion: eventVersion,
		Timestamp:    time.Now(),
		Payload:      payload,
	}, nil
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
// Implementación:
//  1. Primero intenta extraer el valor de Metadata["s3_key"]
//  2. Si no está disponible, parsea FileURL para extraer la key de la ruta
//  3. Como último recurso, retorna FileURL completo
//
// Este es un método de compatibilidad para sistemas legacy.
//
// Retorna:
//   - La clave S3 si está en Metadata["s3_key"]
//   - La key parseada de FileURL si es una URL de S3
//   - FileURL completo como fallback
func (e MaterialUploadedEvent) GetS3Key() string {
	// 1. Intentar extraer de metadata primero (preferido)
	if e.Payload.Metadata != nil {
		if key, ok := e.Payload.Metadata["s3_key"].(string); ok && key != "" {
			return key
		}
	}

	// 2. Parsear FileURL para extraer la key
	parsedURL, err := url.Parse(e.Payload.FileURL)
	if err == nil && parsedURL.Path != "" {
		// Eliminar el primer "/" del path para obtener la key
		key := strings.TrimPrefix(parsedURL.Path, "/")
		if key != "" {
			return key
		}
	}

	// 3. Fallback: retornar FileURL completo
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
