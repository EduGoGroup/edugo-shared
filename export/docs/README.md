# Export — Documentación técnica

Módulo agnóstico para exportación de documentos estructurados en múltiples formatos.

## Propósito

El módulo `export` proporciona una abstracción agnóstica que permite exportar datos estructurados en diversos formatos (Markdown, PDF, HTML, Excel) sin acoplamiento a implementaciones específicas. Define un modelo estándar de `Document` (encabezado, secciones, pie de página) que los renderers pueden transformar a cualquier formato de salida.

## Componentes principales

### Document — Estructura estándar de documento

Contrato central que define la estructura de cualquier documento exportable. Incluye encabezado obligatorio, secciones de contenido y pie de página opcional.

**Estructura:**
```go
type Document struct {
    Header   Header    `json:"header"`
    Sections []Section `json:"sections"`
    Footer   Footer    `json:"footer,omitempty"`
}
```

**Características:**
- Encabezado obligatorio con metadatos del documento
- Múltiples secciones para organizar contenido
- Pie de página opcional para mensajes de cierre
- Serializable a JSON para almacenamiento o transmisión

**Ejemplo:**
```go
doc := export.Document{
    Header: export.Header{
        Title:       "Reporte de Evaluación",
        Subtitle:    "Q1 2026",
        GeneratedAt: time.Now(),
        GeneratedBy: "Sistema EduGo",
    },
    Sections: []export.Section{
        {
            Title: "Estadísticas",
            Rows: []export.Row{
                {Label: "Participantes", Value: "45"},
                {Label: "Promedio", Value: "4.3/5.0"},
            },
        },
    },
}
```

### Header — Metadatos del documento

Encabezado que contiene información de contexto del documento.

**Métodos:**
- `Header.Title` — Título principal del documento (obligatorio)
- `Header.Subtitle` — Subtítulo opcional
- `Header.GeneratedAt` — Marca temporal de generación
- `Header.GeneratedBy` — Identificación de qué sistema generó el documento

### Section — Organización de contenido

Representa una sección temática del documento con título, datos y nota opcional.

**Métodos:**
- `Section.Title` — Nombre de la sección
- `Section.Rows` — Datos clave-valor
- `Section.Note` — Nota adicional (opcional)

**Características:**
- Filas de datos clave-valor para presentar información estructurada
- Notas opcionales para contexto adicional
- Múltiples secciones por documento

**Ejemplo:**
```go
section := export.Section{
    Title: "Calificaciones por Competencia",
    Rows: []export.Row{
        {Label: "Lectura Crítica", Value: "4.5"},
        {Label: "Escritura", Value: "4.2"},
    },
    Note: "Basado en evaluaciones acumulativas",
}
```

### Renderer — Interfaz de transformación

Interfaz que define cómo convertir un `Document` a bytes en un formato específico.

**Función/Métodos:**
```go
type Renderer interface {
    Render(doc Document) ([]byte, string, error) // bytes, mime-type, error
    Extension() string                           // "md", "pdf", "html"
}
```

**Métodos principales:**
- `Render(doc Document)` — Transforma el documento a bytes, retorna MIME type
- `Extension()` — Devuelve la extensión de archivo del formato

**Características:**
- Interfaz agnóstica: cualquier formato puede implementarse
- MIME types estandarizados
- Extensiones de archivo explícitas
- Manejo de errores integrado

**Ejemplo:**
```go
renderer := &export.MarkdownRenderer{}
content, mimeType, err := renderer.Render(doc)
// mimeType = "text/markdown"
// content = bytes del Markdown
```

### MarkdownRenderer — Implementación para Markdown

Implementación concreta que renderiza documentos a formato Markdown.

**Métodos principales:**
- `Extension()` — Retorna `"md"`
- `Render(doc)` — Convierte a Markdown con estructura de headings, listas, notas

**Características:**
- Encabezado renderizado como `# Título` y `## Subtítulo`
- Metadata del encabezado como texto informativo
- Secciones como `## Título de Sección`
- Filas como listas con formato `- **Label:** Value`
- Notas como bloques de cita Markdown
- Pie de página como línea separadora

**Ejemplo:**
```go
renderer := &export.MarkdownRenderer{}
bytes, "text/markdown", _ := renderer.Render(doc)
// Devuelve Markdown formateado listo para guardar en .md

// Salida:
// # Reporte de Evaluación
//
// _Q1 2026_
//
// **Generado:** 2026-04-02 14:30
//
// ---
//
// ## Estadísticas
//
// - **Participantes:** 45
// - **Promedio:** 4.3/5.0
```

## Flujos comunes

### 1. Crear documento simple y exportar a Markdown

```go
package main

import (
    "io/ioutil"
    "time"
    "github.com/EduGoGroup/edugo-shared/export"
)

func main() {
    // Crear documento
    doc := export.Document{
        Header: export.Header{
            Title:       "Boleta de Calificaciones",
            GeneratedAt: time.Now(),
        },
        Sections: []export.Section{
            {
                Title: "Notas Finales",
                Rows: []export.Row{
                    {Label: "Matemática", Value: "4.5"},
                    {Label: "Lengua", Value: "4.8"},
                },
            },
        },
    }

    // Renderizar a Markdown
    renderer := &export.MarkdownRenderer{}
    content, _, err := renderer.Render(doc)
    if err != nil {
        panic(err)
    }

    // Guardar a archivo
    ioutil.WriteFile("boleta.md", content, 0644)
}
```

### 2. Generar documento dinámicamente desde datos de base de datos

```go
package main

import (
    "database/sql"
    "time"
    "github.com/EduGoGroup/edugo-shared/export"
)

func generateStudentReport(db *sql.DB, studentID int) (export.Document, error) {
    // Leer datos de estudiante
    var name, grade string
    var average float64

    row := db.QueryRow(`
        SELECT name, current_grade, average_score
        FROM students WHERE id = ?
    `, studentID)

    err := row.Scan(&name, &grade, &average)
    if err != nil {
        return export.Document{}, err
    }

    // Construir documento
    doc := export.Document{
        Header: export.Header{
            Title:       "Reporte de Estudiante",
            Subtitle:    name,
            GeneratedAt: time.Now(),
            GeneratedBy: "Sistema de Calificaciones",
        },
        Sections: []export.Section{
            {
                Title: "Información Académica",
                Rows: []export.Row{
                    {Label: "Nombre", Value: name},
                    {Label: "Grado Actual", Value: grade},
                    {Label: "Promedio", Value: fmt.Sprintf("%.2f", average)},
                },
            },
        },
    }

    return doc, nil
}

func main() {
    // db := openDatabase()
    // doc, _ := generateStudentReport(db, 123)
    // renderer := &export.MarkdownRenderer{}
    // content, _, _ := renderer.Render(doc)
}
```

### 3. Implementar formato personalizado (PDF, HTML, Excel)

```go
package main

import (
    "bytes"
    "fmt"
    "github.com/EduGoGroup/edugo-shared/export"
)

// CSVRenderer renderiza documentos a CSV
type CSVRenderer struct{}

func (r *CSVRenderer) Extension() string {
    return "csv"
}

func (r *CSVRenderer) Render(doc export.Document) ([]byte, string, error) {
    var buf bytes.Buffer

    // Escribir encabezado como primera fila
    fmt.Fprintf(&buf, "Reporte,%s\n", doc.Header.Title)
    fmt.Fprintf(&buf, "Generado,%s\n\n", doc.Header.GeneratedAt.Format("2006-01-02"))

    // Escribir cada sección como tabla
    for _, section := range doc.Sections {
        fmt.Fprintf(&buf, "%s\n", section.Title)
        for _, row := range section.Rows {
            fmt.Fprintf(&buf, "%s,%s\n", row.Label, row.Value)
        }
        fmt.Fprintf(&buf, "\n")
    }

    return buf.Bytes(), "text/csv", nil
}

func main() {
    doc := export.Document{
        Header: export.Header{Title: "Datos"},
        Sections: []export.Section{
            {
                Title: "Tabla1",
                Rows: []export.Row{
                    {Label: "A", Value: "1"},
                    {Label: "B", Value: "2"},
                },
            },
        },
    }

    renderer := &CSVRenderer{}
    content, mimeType, _ := renderer.Render(doc)
    // content ahora contiene CSV bien formateado
}
```

### 4. Flujo de exportación en servicio HTTP

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/EduGoGroup/edugo-shared/export"
)

// ExportHandler maneja solicitudes de exportación
func ExportHandler(w http.ResponseWriter, r *http.Request) {
    // Leer parámetro de formato
    format := r.URL.Query().Get("format")
    if format == "" {
        format = "markdown"
    }

    // Construir documento (simulado)
    doc := export.Document{
        Header: export.Header{
            Title:       "Reporte API",
            GeneratedAt: time.Now(),
        },
        Sections: []export.Section{
            {
                Title: "Datos",
                Rows: []export.Row{
                    {Label: "Status", Value: "OK"},
                },
            },
        },
    }

    // Renderizar según formato
    var renderer export.Renderer
    switch format {
    case "markdown":
        renderer = &export.MarkdownRenderer{}
    default:
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "Formato no soportado: %s", format)
        return
    }

    content, mimeType, err := renderer.Render(doc)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Enviar respuesta
    w.Header().Set("Content-Type", mimeType)
    w.Header().Set("Content-Disposition", fmt.Sprintf(
        "attachment; filename=report.%s", renderer.Extension(),
    ))
    w.Write(content)
}

func main() {
    http.HandleFunc("/export", ExportHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Arquitectura

Flujo de transformación de datos a documento exportado:

```
Datos (Base de datos, API, archivos)
    ↓
Construir Document (Header, Sections, Footer)
    ↓
Seleccionar Renderer (Markdown, PDF, HTML, CSV, etc.)
    ↓
Render(Document) → (bytes, MIME type, error)
    ↓
Enviar respuesta HTTP o guardar a archivo
```

Modelo de componentes:

```
Document
├─ Header (obligatorio)
│  ├─ Title (requerido)
│  ├─ Subtitle (opcional)
│  ├─ GeneratedAt
│  └─ GeneratedBy
├─ Sections (múltiples, requerido)
│  ├─ Title
│  ├─ Rows (múltiples clave-valor)
│  └─ Note (opcional)
└─ Footer (opcional)
   └─ Text

Renderer Interface
├─ MarkdownRenderer (incluido)
├─ [PDFRenderer] (futuro)
├─ [HTMLRenderer] (futuro)
└─ [CustomRenderer] (usuario)
```

## Dependencias

- **Internas**: Ninguna
- **Externas**:
  - `time` (estándar) — Para timestamps en Header

## Testing

Suite de tests unitarios:

- Creación y validación de estructura `Document`
- Renderización de todos los campos opcionales (Subtitle, Note, Footer)
- Generación correcta de MIME types
- Extensiones de archivo correctas
- Manejo de secciones múltiples y filas vacías
- Salida Markdown bien formateada

Ejecutar:
```bash
make test      # Tests unitarios
make test-race # Race detector
make check     # Validar + tests
```

## Notas de diseño

- **Agnóstico de formato**: La interfaz `Renderer` permite cualquier implementación sin modificar `Document`
- **Composición sobre herencia**: `Document` es una composición simple de tipos, no una jerarquía
- **Serialización JSON**: El modelo es serializable para almacenamiento o APIs
- **Extensibilidad clara**: Nuevos formatos se agregan implementando la interfaz `Renderer`
- **Zero overhead predeterminado**: Sin renderer, `Document` solo es una estructura de datos
