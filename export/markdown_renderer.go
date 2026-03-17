package export

import (
	"bytes"
	"fmt"
)

// MarkdownRenderer renderiza un Document como Markdown.
type MarkdownRenderer struct{}

// Extension devuelve la extensión de archivo para el formato Markdown ("md").
func (r *MarkdownRenderer) Extension() string { return "md" }

// Render convierte un Document a bytes en formato Markdown.
// Devuelve el contenido renderizado, el MIME type "text/markdown" y un error (siempre nil).
// Los campos opcionales (Subtitle, Note, Footer) se omiten si están vacíos.
func (r *MarkdownRenderer) Render(doc Document) ([]byte, string, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# %s\n\n", doc.Header.Title)
	if doc.Header.Subtitle != "" {
		fmt.Fprintf(&buf, "_%s_\n\n", doc.Header.Subtitle)
	}
	fmt.Fprintf(&buf, "**Generado:** %s\n\n---\n\n", doc.Header.GeneratedAt.Format("2006-01-02 15:04"))
	for _, section := range doc.Sections {
		fmt.Fprintf(&buf, "## %s\n\n", section.Title)
		for _, row := range section.Rows {
			fmt.Fprintf(&buf, "- **%s:** %s\n", row.Label, row.Value)
		}
		if section.Note != "" {
			fmt.Fprintf(&buf, "\n> %s\n", section.Note)
		}
		buf.WriteString("\n")
	}
	if doc.Footer.Text != "" {
		fmt.Fprintf(&buf, "---\n_%s_\n", doc.Footer.Text)
	}
	return buf.Bytes(), "text/markdown", nil
}
