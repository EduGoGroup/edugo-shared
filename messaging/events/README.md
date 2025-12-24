# Events - Módulo de Eventos CloudEvents

Módulo Go compartido para eventos siguiendo el estándar CloudEvents, utilizado en el ecosistema de microservicios EduGo.

## Propósito

Este módulo proporciona estructuras de datos estandarizadas para eventos de dominio que se comunican entre los diferentes microservicios de EduGo. Siguiendo el estándar CloudEvents, garantiza interoperabilidad y consistencia en la comunicación asíncrona.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/messaging/events
```

## Ejemplo de Uso

### Crear un Evento de Material Subido (Recomendado)

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/EduGoGroup/edugo-shared/messaging/events"
)

func main() {
    // Crear payload con validación automática
    payload := events.MaterialUploadedPayload{
        MaterialID:    "mat-789",
        SchoolID:      "school-001",
        TeacherID:     "teacher-456",
        FileURL:       "https://s3.amazonaws.com/edugo-materials/file.pdf",
        FileSizeBytes: 2048576, // 2MB (uint64 previene valores negativos)
        FileType:      "application/pdf",
        Metadata: map[string]interface{}{
            "s3_key":   "materials/2024/file.pdf",
            "checksum": "abc123def456",
        },
    }

    // Usar constructor con validación (RECOMENDADO)
    event, err := events.NewMaterialUploadedEvent(
        "evt-123456",
        "material.uploaded",
        "1.0",
        payload,
    )
    if err != nil {
        log.Fatalf("error creando evento: %v", err)
    }

    // Timestamp se establece automáticamente a time.Now()

    // Serializar a JSON para enviar por mensaje
    jsonData, err := json.Marshal(event)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(jsonData))

    // Usar métodos de compatibilidad
    fmt.Printf("Material ID: %s\n", event.GetMaterialID())
    fmt.Printf("S3 Key: %s\n", event.GetS3Key())
    fmt.Printf("Author ID: %s\n", event.GetAuthorID())
}
```

## Features Principales

- ✅ **Estándar CloudEvents**: Cumple con la especificación CloudEvents para eventos distribuidos
- ✅ **Type-Safe**: Estructuras fuertemente tipadas en Go
- ✅ **Validación Automática**: Constructor NewMaterialUploadedEvent valida campos requeridos
- ✅ **Seguridad de Tipos**: FileSizeBytes usa uint64 para prevenir valores negativos
- ✅ **S3 Key Parsing**: Extrae automáticamente la S3 key desde FileURL si no está en metadata
- ✅ **Backward Compatibility**: Métodos helper para integración con sistemas legacy
- ✅ **Metadata Flexible**: Campo opcional para información adicional variable
- ✅ **100% Test Coverage**: Totalmente probado con tests unitarios
- ✅ **Documentación GoDoc**: Documentación completa en español

## Eventos Disponibles

### MaterialUploadedEvent

Representa el evento generado cuando un profesor sube un nuevo material educativo.

**Campos CloudEvents estándar:**
- `EventID`: Identificador único del evento (requerido)
- `EventType`: Tipo de evento (ej: "material.uploaded") (requerido)
- `EventVersion`: Versión del esquema del evento (requerido)
- `Timestamp`: Marca de tiempo del evento (auto-generado)

**Payload del dominio:**
- `MaterialID`: ID del material (requerido)
- `SchoolID`: ID de la escuela (requerido)
- `TeacherID`: ID del profesor (requerido)
- `FileURL`: URL completa del archivo (requerido, debe ser URL válida)
- `FileSizeBytes`: Tamaño en bytes (uint64, no puede ser negativo)
- `FileType`: Tipo MIME (requerido)
- `Metadata`: Datos adicionales opcionales

**Métodos de compatibilidad:**
- `GetMaterialID()`: Obtiene el ID del material
- `GetS3Key()`: Obtiene la clave S3 (prioridad: metadata["s3_key"] > parseado de FileURL > FileURL completo)
- `GetAuthorID()`: Obtiene el ID del profesor/autor

**Validaciones:**

El constructor `NewMaterialUploadedEvent` valida:
- Todos los campos requeridos no estén vacíos
- FileURL sea una URL válida
- FileSizeBytes es uint64 (no puede ser negativo por diseño)

## GetS3Key - Estrategia de Resolución

El método `GetS3Key()` utiliza la siguiente estrategia de resolución:

1. **Metadata prioritario**: Si existe `metadata["s3_key"]`, lo retorna
2. **Parsing de URL**: Si no hay metadata, parsea `FileURL` para extraer el path
3. **Fallback**: Retorna `FileURL` completo como último recurso

```go
// Ejemplo 1: Desde metadata (preferido)
event := events.MaterialUploadedEvent{
    Payload: events.MaterialUploadedPayload{
        FileURL: "https://s3.amazonaws.com/bucket/file.pdf",
        Metadata: map[string]interface{}{
            "s3_key": "uploads/2024/file.pdf",
        },
    },
}
fmt.Println(event.GetS3Key()) // Output: "uploads/2024/file.pdf"

// Ejemplo 2: Parseado automático de URL
event := events.MaterialUploadedEvent{
    Payload: events.MaterialUploadedPayload{
        FileURL: "https://s3.amazonaws.com/bucket/uploads/2024/file.pdf",
    },
}
fmt.Println(event.GetS3Key()) // Output: "bucket/uploads/2024/file.pdf"
```

## Documentación

- [GoDoc](https://pkg.go.dev/github.com/EduGoGroup/edugo-shared/messaging/events)
- [CloudEvents Spec](https://cloudevents.io/)

## Desarrollo

### Ejecutar Tests

```bash
cd messaging/events
go test -v
```

### Ejecutar Tests con Cobertura

```bash
cd messaging/events
go test -v -cover
```

### Verificar Formato

```bash
cd messaging/events
go fmt ./...
```

## Licencia

Propiedad de EduGo Group
