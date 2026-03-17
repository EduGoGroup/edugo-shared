package export

import (
	"bytes"
	"fmt"
)

// MarkdownRenderer renderiza un ExportDocument como Markdown.
type MarkdownRenderer struct{}

func (r *MarkdownRenderer) Extension() string { return "md" }

func (r *MarkdownRenderer) Render(doc ExportDocument) ([]byte, string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("# %s\n\n", doc.Header.Title))
	if doc.Header.Subtitle != "" {
		buf.WriteString(fmt.Sprintf("_%s_\n\n", doc.Header.Subtitle))
	}
	buf.WriteString(fmt.Sprintf("**Generado:** %s\n\n---\n\n", doc.Header.GeneratedAt.Format("2006-01-02 15:04")))
	for _, section := range doc.Sections {
		buf.WriteString(fmt.Sprintf("## %s\n\n", section.Title))
		for _, row := range section.Rows {
			buf.WriteString(fmt.Sprintf("- **%s:** %s\n", row.Label, row.Value))
		}
		if section.Note != "" {
			buf.WriteString(fmt.Sprintf("\n> %s\n", section.Note))
		}
		buf.WriteString("\n")
	}
	if doc.Footer.Text != "" {
		buf.WriteString(fmt.Sprintf("---\n_%s_\n", doc.Footer.Text))
	}
	return buf.Bytes(), "text/markdown", nil
}
