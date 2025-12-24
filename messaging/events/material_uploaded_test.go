package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaterialUploadedEvent_Serialization(t *testing.T) {
	now := time.Now().UTC()
	
	event := MaterialUploadedEvent{
		EventID:      "evt_123",
		EventType:    "material.uploaded",
		EventVersion: "1.0",
		Timestamp:    now,
		Payload: MaterialUploadedPayload{
			MaterialID:    "mat_456",
			SchoolID:      "school_789",
			TeacherID:     "teacher_012",
			FileURL:       "https://s3.amazonaws.com/bucket/file.pdf",
			FileSizeBytes: 1024000,
			FileType:      "application/pdf",
			Metadata: map[string]interface{}{
				"s3_key":       "uploads/2024/file.pdf",
				"content_type": "application/pdf",
			},
		},
	}

	// Serializar a JSON
	data, err := json.Marshal(event)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Deserializar desde JSON
	var decoded MaterialUploadedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	// Verificar todos los campos
	assert.Equal(t, event.EventID, decoded.EventID)
	assert.Equal(t, event.EventType, decoded.EventType)
	assert.Equal(t, event.EventVersion, decoded.EventVersion)
	assert.Equal(t, event.Payload.MaterialID, decoded.Payload.MaterialID)
	assert.Equal(t, event.Payload.SchoolID, decoded.Payload.SchoolID)
	assert.Equal(t, event.Payload.TeacherID, decoded.Payload.TeacherID)
	assert.Equal(t, event.Payload.FileURL, decoded.Payload.FileURL)
	assert.Equal(t, event.Payload.FileSizeBytes, decoded.Payload.FileSizeBytes)
	assert.Equal(t, event.Payload.FileType, decoded.Payload.FileType)
	assert.Equal(t, "uploads/2024/file.pdf", decoded.Payload.Metadata["s3_key"])
}

func TestMaterialUploadedEvent_GetMaterialID(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			MaterialID: "mat_123",
		},
	}

	assert.Equal(t, "mat_123", event.GetMaterialID())
}

func TestMaterialUploadedEvent_GetAuthorID(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			TeacherID: "teacher_456",
		},
	}

	assert.Equal(t, "teacher_456", event.GetAuthorID())
}

func TestMaterialUploadedEvent_GetS3Key_FromMetadata(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			FileURL: "https://s3.amazonaws.com/bucket/file.pdf",
			Metadata: map[string]interface{}{
				"s3_key": "uploads/2024/12/file.pdf",
			},
		},
	}

	assert.Equal(t, "uploads/2024/12/file.pdf", event.GetS3Key())
}

func TestMaterialUploadedEvent_GetS3Key_FallbackToFileURL(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			FileURL: "https://s3.amazonaws.com/bucket/uploads/file.pdf",
		},
	}

	// Sin metadata, debe retornar FileURL completo
	assert.Equal(t, "https://s3.amazonaws.com/bucket/uploads/file.pdf", event.GetS3Key())
}

func TestMaterialUploadedEvent_GetS3Key_NilMetadata(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			FileURL:  "https://s3.amazonaws.com/bucket/file.pdf",
			Metadata: nil,
		},
	}

	// Metadata nil, debe retornar FileURL
	assert.Equal(t, "https://s3.amazonaws.com/bucket/file.pdf", event.GetS3Key())
}

func TestMaterialUploadedEvent_GetS3Key_EmptyMetadata(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			FileURL:  "https://s3.amazonaws.com/bucket/file.pdf",
			Metadata: map[string]interface{}{},
		},
	}

	// Metadata vacío (sin s3_key), debe retornar FileURL
	assert.Equal(t, "https://s3.amazonaws.com/bucket/file.pdf", event.GetS3Key())
}

func TestMaterialUploadedPayload_JSONWithoutMetadata(t *testing.T) {
	payload := MaterialUploadedPayload{
		MaterialID:    "mat_789",
		SchoolID:      "school_012",
		TeacherID:     "teacher_345",
		FileURL:       "https://example.com/file.pdf",
		FileSizeBytes: 2048,
		FileType:      "application/pdf",
		// Metadata omitido (nil)
	}

	data, err := json.Marshal(payload)
	require.NoError(t, err)

	// Metadata debe ser omitido del JSON si es nil (gracias a omitempty)
	assert.NotContains(t, string(data), "metadata")

	// Deserializar debe funcionar correctamente
	var decoded MaterialUploadedPayload
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, payload.MaterialID, decoded.MaterialID)
	assert.Nil(t, decoded.Metadata)
}

func TestMaterialUploadedEvent_CompleteStructure(t *testing.T) {
	tests := []struct {
		name  string
		event MaterialUploadedEvent
	}{
		{
			name: "Evento completo con todos los campos",
			event: MaterialUploadedEvent{
				EventID:      "evt_complete",
				EventType:    "material.uploaded",
				EventVersion: "1.0",
				Timestamp:    time.Now().UTC(),
				Payload: MaterialUploadedPayload{
					MaterialID:    "mat_001",
					SchoolID:      "school_001",
					TeacherID:     "teacher_001",
					FileURL:       "https://cdn.example.com/materials/document.pdf",
					FileSizeBytes: 1024768,
					FileType:      "application/pdf",
					Metadata: map[string]interface{}{
						"s3_key":       "materials/2024/12/document.pdf",
						"content_type": "application/pdf",
						"uploaded_by":  "teacher_001",
						"is_public":    false,
					},
				},
			},
		},
		{
			name: "Evento mínimo sin metadata",
			event: MaterialUploadedEvent{
				EventID:      "evt_minimal",
				EventType:    "material.uploaded",
				EventVersion: "1.0",
				Timestamp:    time.Now().UTC(),
				Payload: MaterialUploadedPayload{
					MaterialID:    "mat_002",
					SchoolID:      "school_002",
					TeacherID:     "teacher_002",
					FileURL:       "https://storage.example.com/file.docx",
					FileSizeBytes: 512000,
					FileType:      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serializar
			data, err := json.Marshal(tt.event)
			require.NoError(t, err)

			// Deserializar
			var decoded MaterialUploadedEvent
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err)

			// Verificar campos obligatorios
			assert.Equal(t, tt.event.EventID, decoded.EventID)
			assert.Equal(t, tt.event.EventType, decoded.EventType)
			assert.Equal(t, tt.event.EventVersion, decoded.EventVersion)
			assert.Equal(t, tt.event.Payload.MaterialID, decoded.Payload.MaterialID)
			assert.Equal(t, tt.event.Payload.SchoolID, decoded.Payload.SchoolID)
			assert.Equal(t, tt.event.Payload.TeacherID, decoded.Payload.TeacherID)
			assert.Equal(t, tt.event.Payload.FileURL, decoded.Payload.FileURL)
			assert.Equal(t, tt.event.Payload.FileSizeBytes, decoded.Payload.FileSizeBytes)
			assert.Equal(t, tt.event.Payload.FileType, decoded.Payload.FileType)

			// Verificar métodos de compatibilidad
			assert.Equal(t, tt.event.Payload.MaterialID, decoded.GetMaterialID())
			assert.Equal(t, tt.event.Payload.TeacherID, decoded.GetAuthorID())
			assert.NotEmpty(t, decoded.GetS3Key())
		})
	}
}

func TestMaterialUploadedPayload_RequiredFields(t *testing.T) {
	payload := MaterialUploadedPayload{
		MaterialID:    "mat_required",
		SchoolID:      "school_required",
		TeacherID:     "teacher_required",
		FileURL:       "https://example.com/required.pdf",
		FileSizeBytes: 1024,
		FileType:      "application/pdf",
	}

	// Todos los campos requeridos deben estar presentes
	assert.NotEmpty(t, payload.MaterialID)
	assert.NotEmpty(t, payload.SchoolID)
	assert.NotEmpty(t, payload.TeacherID)
	assert.NotEmpty(t, payload.FileURL)
	assert.Greater(t, payload.FileSizeBytes, int64(0))
	assert.NotEmpty(t, payload.FileType)
}
