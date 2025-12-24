# Events - Módulo de Eventos CloudEvents

Módulo Go compartido para eventos siguiendo el estándar CloudEvents, utilizado en el ecosistema de microservicios EduGo.

## Propósito

Este módulo proporciona estructuras de datos estandarizadas para eventos de dominio que se comunican entre los diferentes microservicios de EduGo. Siguiendo el estándar CloudEvents, garantiza interoperabilidad y consistencia en la comunicación asíncrona.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/messaging/events
```

## Ejemplo de Uso

### Crear un Evento de Material Subido

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/EduGoGroup/edugo-shared/messaging/events"
)

func main() {
    // Crear evento de material subido
    event := events.MaterialUploadedEvent{
        EventID:      "evt-123456",
        EventType:    "material.uploaded",
        EventVersion: "1.0",
        Timestamp:    time.Now(),
        Payload: events.MaterialUploadedPayload{
            MaterialID:    "mat-789",
            SchoolID:      "school-001",
            TeacherID:     "teacher-456",
            FileURL:       "https://s3.amazonaws.com/edugo-materials/file.pdf",
            FileSizeBytes: 2048576, // 2MB
            FileType:      "application/pdf",
            Metadata: map[string]interface{}{
                "s3_key":   "materials/2024/file.pdf",
                "checksum": "abc123def456",
            },
        },
    }
    
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
- ✅ **Backward Compatibility**: Métodos helper para integración con sistemas legacy
- ✅ **Metadata Flexible**: Campo opcional para información adicional variable
- ✅ **100% Test Coverage**: Totalmente probado con tests unitarios
- ✅ **Documentación GoDoc**: Documentación completa en español

## Eventos Disponibles

### MaterialUploadedEvent

Representa el evento generado cuando un profesor sube un nuevo material educativo.

**Campos CloudEvents estándar:**
- `EventID`: Identificador único del evento
- `EventType`: Tipo de evento (ej: "material.uploaded")
- `EventVersion`: Versión del esquema del evento
- `Timestamp`: Marca de tiempo del evento

**Payload del dominio:**
- `MaterialID`: ID del material
- `SchoolID`: ID de la escuela
- `TeacherID`: ID del profesor
- `FileURL`: URL completa del archivo
- `FileSizeBytes`: Tamaño en bytes
- `FileType`: Tipo MIME
- `Metadata`: Datos adicionales opcionales

**Métodos de compatibilidad:**
- `GetMaterialID()`: Obtiene el ID del material
- `GetS3Key()`: Obtiene la clave S3 (desde metadata o FileURL)
- `GetAuthorID()`: Obtiene el ID del profesor/autor

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
