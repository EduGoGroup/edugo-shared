# Export

Módulo agnóstico para exportación de documentos estructurados en múltiples formatos (Markdown, PDF, HTML).

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/export
```

El módulo se descarga como `export`, principal consumo vía package `export`.

## Quick Start

### Ejemplo 1: Crear y renderizar documento Markdown

```go
package main

import (
    "fmt"
    "time"
    "github.com/EduGoGroup/edugo-shared/export"
)

func main() {
    // Crear un documento con encabezado, secciones y pie de página
    doc := export.Document{
        Header: export.Header{
            Title:       "Reporte de Evaluación",
            Subtitle:    "Semestre 2026-I",
            GeneratedAt: time.Now(),
            GeneratedBy: "Sistema EduGo",
        },
        Sections: []export.Section{
            {
                Title: "Calificaciones",
                Rows: []export.Row{
                    {Label: "Matemática", Value: "4.5/5.0"},
                    {Label: "Lengua", Value: "4.8/5.0"},
                },
            },
        },
        Footer: export.Footer{
            Text: "Generado automáticamente por EduGo",
        },
    }

    // Renderizar a Markdown
    renderer := &export.MarkdownRenderer{}
    content, mimeType, err := renderer.Render(doc)
    if err != nil {
        panic(err)
    }

    fmt.Printf("MIME Type: %s\n", mimeType)
    fmt.Printf("Contenido:\n%s\n", string(content))
}
```

### Ejemplo 2: Exportar reporte de asistencia

```go
package main

import (
    "fmt"
    "time"
    "github.com/EduGoGroup/edugo-shared/export"
)

func main() {
    // Construir documento programáticamente desde datos
    sections := make([]export.Section, 0)

    // Sección de estudiantes presentes
    sections = append(sections, export.Section{
        Title: "Asistencia Presente",
        Rows: []export.Row{
            {Label: "Juan Pérez", Value: "✓ Presente"},
            {Label: "María García", Value: "✓ Presente"},
            {Label: "Carlos López", Value: "✓ Presente"},
        },
        Note: "Total: 3 estudiantes presentes",
    })

    // Sección de estudiantes ausentes
    sections = append(sections, export.Section{
        Title: "Asistencia Ausentes",
        Rows: []export.Row{
            {Label: "Ana Martínez", Value: "✗ Ausente"},
        },
        Note: "Total: 1 estudiante ausente",
    })

    doc := export.Document{
        Header: export.Header{
            Title:       "Registro de Asistencia",
            Subtitle:    "Clase de Matemática - 2026-04-02",
            GeneratedAt: time.Now(),
        },
        Sections: sections,
    }

    renderer := &export.MarkdownRenderer{}
    content, _, _ := renderer.Render(doc)
    fmt.Println(string(content))
}
```

### Ejemplo 3: Implementar renderer personalizado (extensibilidad)

```go
package main

import (
    "bytes"
    "fmt"
    "github.com/EduGoGroup/edugo-shared/export"
)

// HTMLRenderer renderiza documentos a HTML
type HTMLRenderer struct{}

func (r *HTMLRenderer) Extension() string {
    return "html"
}

func (r *HTMLRenderer) Render(doc export.Document) ([]byte, string, error) {
    var buf bytes.Buffer

    fmt.Fprintf(&buf, "<html><body>\n")
    fmt.Fprintf(&buf, "<h1>%s</h1>\n", doc.Header.Title)

    for _, section := range doc.Sections {
        fmt.Fprintf(&buf, "<h2>%s</h2>\n<ul>\n", section.Title)
        for _, row := range section.Rows {
            fmt.Fprintf(&buf, "<li><strong>%s:</strong> %s</li>\n", row.Label, row.Value)
        }
        fmt.Fprintf(&buf, "</ul>\n")
    }

    fmt.Fprintf(&buf, "</body></html>\n")
    return buf.Bytes(), "text/html", nil
}

func main() {
    // Usar el renderer personalizado
    doc := export.Document{
        Header: export.Header{
            Title: "Reporte HTML",
        },
        Sections: []export.Section{
            {
                Title: "Datos",
                Rows: []export.Row{
                    {Label: "Campo1", Value: "Valor1"},
                },
            },
        },
    }

    renderer := &HTMLRenderer{}
    content, mimeType, _ := renderer.Render(doc)
    fmt.Printf("Renderizado como: %s\n%s\n", mimeType, string(content))
}
```

### Ejemplo 4: Construir documento dinámicamente desde múltiples fuentes

```go
package main

import (
    "fmt"
    "time"
    "github.com/EduGoGroup/edugo-shared/export"
)

// ReportBuilder construye reportes dinámicamente
type ReportBuilder struct {
    title    string
    subtitle string
    sections []export.Section
}

func NewReportBuilder(title string) *ReportBuilder {
    return &ReportBuilder{title: title, sections: []export.Section{}}
}

func (rb *ReportBuilder) AddSection(title string, rows []export.Row, note string) *ReportBuilder {
    rb.sections = append(rb.sections, export.Section{
        Title: title,
        Rows:  rows,
        Note:  note,
    })
    return rb
}

func (rb *ReportBuilder) Build() export.Document {
    return export.Document{
        Header: export.Header{
            Title:       rb.title,
            Subtitle:    rb.subtitle,
            GeneratedAt: time.Now(),
        },
        Sections: rb.sections,
    }
}

func main() {
    report := NewReportBuilder("Reporte de Calificaciones").
        AddSection("Primer Parcial", []export.Row{
            {Label: "Promedio", Value: "4.2"},
        }, "5 estudiantes evaluados").
        AddSection("Segundo Parcial", []export.Row{
            {Label: "Promedio", Value: "4.5"},
        }, "5 estudiantes evaluados").
        Build()

    renderer := &export.MarkdownRenderer{}
    content, _, _ := renderer.Render(report)
    fmt.Println(string(content))
}
```

## Componentes principales

- **Document**: Estructura estándar de documento con encabezado, secciones y pie de página
- **Renderer**: Interfaz extensible para múltiples formatos de salida
- **MarkdownRenderer**: Implementación para renderizar documentos a Markdown
- **Header/Section/Footer**: Componentes estructurales del documento

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build     # Compilar
make test      # Tests unitarios
make test-race # Race detector
make check     # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- **Agnóstico de formato**: La interfaz `Renderer` permite implementar cualquier formato (PDF, HTML, Excel, etc.)
- **Estructura flexible**: El modelo `Document` es lo suficientemente flexible para representar reportes complejos
- **Facilidad de extensión**: Nuevos formatos se implementan sin modificar código existente
- **Zero overhead con NoOp**: Usar como base sin renderer personalizado tiene overhead mínimo
