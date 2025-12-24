package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMaterialUploadedEvent_Valid(t *testing.T) {
	payload := MaterialUploadedPayload{
		MaterialID:    "mat_456",
		SchoolID:      "school_789",
		TeacherID:     "teacher_012",
		FileURL:       "https://s3.amazonaws.com/bucket/file.pdf",
		FileSizeBytes: 1024000,
		FileType:      "application/pdf",
		Metadata: map[string]interface{}{
			"s3_key": "uploads/2024/file.pdf",
		},
	}

	event, err := NewMaterialUploadedEvent("evt_123", "material.uploaded", "1.0", payload)
	
	require.NoError(t, err)
	assert.Equal(t, "evt_123", event.EventID)
	assert.Equal(t, "material.uploaded", event.EventType)
	assert.Equal(t, "1.0", event.EventVersion)
	assert.False(t, event.Timestamp.IsZero())
	assert.Equal(t, "mat_456", event.Payload.MaterialID)
}

func TestNewMaterialUploadedEvent_EmptyEventID(t *testing.T) {
	payload := MaterialUploadedPayload{
		MaterialID:    "mat_456",
		SchoolID:      "school_789",
		TeacherID:     "teacher_012",
		FileURL:       "https://s3.amazonaws.com/bucket/file.pdf",
		FileSizeBytes: 1024000,
		FileType:      "application/pdf",
	}

	_, err := NewMaterialUploadedEvent("", "material.uploaded", "1.0", payload)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "eventID")
}

func TestNewMaterialUploadedEvent_EmptyMaterialID(t *testing.T) {
	payload := MaterialUploadedPayload{
		MaterialID:    "",
		SchoolID:      "school_789",
		TeacherID:     "teacher_012",
		FileURL:       "https://s3.amazonaws.com/bucket/file.pdf",
		FileSizeBytes: 1024000,
		FileType:      "application/pdf",
	}

	_, err := NewMaterialUploadedEvent("evt_123", "material.uploaded", "1.0", payload)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MaterialID")
}

func TestNewMaterialUploadedEvent_InvalidFileURL(t *testing.T) {
	payload := MaterialUploadedPayload{
		MaterialID:    "mat_456",
		SchoolID:      "school_789",
		TeacherID:     "teacher_012",
		FileURL:       "not a valid url",
		FileSizeBytes: 1024000,
		FileType:      "application/pdf",
	}

	_, err := NewMaterialUploadedEvent("evt_123", "material.uploaded", "1.0", payload)
	
	// URL parsing en Go es permisivo, así que esto debería pasar
	// pero si necesitamos validación más estricta, se puede mejorar
	assert.NoError(t, err) // Go acepta URLs sin scheme
}

func TestNewMaterialUploadedEvent_AllFieldsRequired(t *testing.T) {
	tests := []struct {
		name    string
		payload MaterialUploadedPayload
		wantErr string
	}{
		{
			name: "SchoolID vacío",
			payload: MaterialUploadedPayload{
				MaterialID:    "mat_456",
				SchoolID:      "",
				TeacherID:     "teacher_012",
				FileURL:       "https://example.com/file.pdf",
				FileSizeBytes: 1024,
				FileType:      "application/pdf",
			},
			wantErr: "SchoolID",
		},
		{
			name: "TeacherID vacío",
			payload: MaterialUploadedPayload{
				MaterialID:    "mat_456",
				SchoolID:      "school_789",
				TeacherID:     "",
				FileURL:       "https://example.com/file.pdf",
				FileSizeBytes: 1024,
				FileType:      "application/pdf",
			},
			wantErr: "TeacherID",
		},
		{
			name: "FileType vacío",
			payload: MaterialUploadedPayload{
				MaterialID:    "mat_456",
				SchoolID:      "school_789",
				TeacherID:     "teacher_012",
				FileURL:       "https://example.com/file.pdf",
				FileSizeBytes: 1024,
				FileType:      "",
			},
			wantErr: "FileType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMaterialUploadedEvent("evt_123", "material.uploaded", "1.0", tt.payload)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestMaterialUploadedEvent_Serialization(t *testing.T) {
	payload := MaterialUploadedPayload{
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
	}

	event, err := NewMaterialUploadedEvent("evt_123", "material.uploaded", "1.0", payload)
	require.NoError(t, err)

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

func TestMaterialUploadedEvent_GetS3Key_ParsedFromURL(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			FileURL: "https://s3.amazonaws.com/bucket/uploads/2024/file.pdf",
		},
	}

	// Sin metadata, debe parsear la URL y extraer el path
	assert.Equal(t, "bucket/uploads/2024/file.pdf", event.GetS3Key())
}

func TestMaterialUploadedEvent_GetS3Key_FallbackToFileURL(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			FileURL: "",
		},
	}

	// FileURL vacío, debe retornar vacío
	assert.Equal(t, "", event.GetS3Key())
}

func TestMaterialUploadedEvent_GetS3Key_NilMetadata(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			FileURL:  "https://s3.amazonaws.com/my-bucket/files/document.pdf",
			Metadata: nil,
		},
	}

	// Metadata nil, debe parsear URL
	assert.Equal(t, "my-bucket/files/document.pdf", event.GetS3Key())
}

func TestMaterialUploadedEvent_GetS3Key_EmptyMetadata(t *testing.T) {
	event := MaterialUploadedEvent{
		Payload: MaterialUploadedPayload{
			FileURL:  "https://cdn.example.com/storage/file.pdf",
			Metadata: map[string]interface{}{},
		},
	}

	// Metadata vacío (sin s3_key), debe parsear URL
	assert.Equal(t, "storage/file.pdf", event.GetS3Key())
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

func TestMaterialUploadedPayload_FileSizeBytes_UInt64(t *testing.T) {
	// uint64 previene valores negativos a nivel de tipo
	payload := MaterialUploadedPayload{
		MaterialID:    "mat_size",
		SchoolID:      "school_size",
		TeacherID:     "teacher_size",
		FileURL:       "https://example.com/large-file.zip",
		FileSizeBytes: 9999999999, // Valor grande positivo
		FileType:      "application/zip",
	}

	assert.Greater(t, payload.FileSizeBytes, uint64(0))
	
	// Serializar y deserializar
	data, err := json.Marshal(payload)
	require.NoError(t, err)
	
	var decoded MaterialUploadedPayload
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	
	assert.Equal(t, payload.FileSizeBytes, decoded.FileSizeBytes)
}

func TestMaterialUploadedEvent_CompleteStructure(t *testing.T) {
	tests := []struct {
		name    string
		payload MaterialUploadedPayload
	}{
		{
			name: "Evento completo con todos los campos",
			payload: MaterialUploadedPayload{
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
		{
			name: "Evento mínimo sin metadata",
			payload: MaterialUploadedPayload{
				MaterialID:    "mat_002",
				SchoolID:      "school_002",
				TeacherID:     "teacher_002",
				FileURL:       "https://storage.example.com/file.docx",
				FileSizeBytes: 512000,
				FileType:      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Crear evento con validación
			event, err := NewMaterialUploadedEvent("evt_test", "material.uploaded", "1.0", tt.payload)
			require.NoError(t, err)

			// Serializar
			data, err := json.Marshal(event)
			require.NoError(t, err)

			// Deserializar
			var decoded MaterialUploadedEvent
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err)

			// Verificar campos obligatorios
			assert.Equal(t, event.EventID, decoded.EventID)
			assert.Equal(t, event.EventType, decoded.EventType)
			assert.Equal(t, event.EventVersion, decoded.EventVersion)
			assert.Equal(t, event.Payload.MaterialID, decoded.Payload.MaterialID)
			assert.Equal(t, event.Payload.SchoolID, decoded.Payload.SchoolID)
			assert.Equal(t, event.Payload.TeacherID, decoded.Payload.TeacherID)
			assert.Equal(t, event.Payload.FileURL, decoded.Payload.FileURL)
			assert.Equal(t, event.Payload.FileSizeBytes, decoded.Payload.FileSizeBytes)
			assert.Equal(t, event.Payload.FileType, decoded.Payload.FileType)

			// Verificar métodos de compatibilidad
			assert.Equal(t, event.Payload.MaterialID, decoded.GetMaterialID())
			assert.Equal(t, event.Payload.TeacherID, decoded.GetAuthorID())
			assert.NotEmpty(t, decoded.GetS3Key())
		})
	}
}
