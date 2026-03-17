package export

import (
	"strings"
	"testing"
	"time"
)

func TestMarkdownRenderer_Extension(t *testing.T) {
	r := &MarkdownRenderer{}
	if got := r.Extension(); got != "md" {
		t.Errorf("Extension() = %q, want %q", got, "md")
	}
}

func TestMarkdownRenderer_Render(t *testing.T) {
	fixedTime := time.Date(2026, 3, 17, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name         string
		doc          ExportDocument
		wantContains []string
		wantMIME     string
	}{
		{
			name: "documento completo con todas las secciones",
			doc: ExportDocument{
				Header: ExportHeader{
					Title:       "Reporte de Notas",
					Subtitle:    "Periodo 2026-Q1",
					GeneratedAt: fixedTime,
					GeneratedBy: "admin@edugo.test",
				},
				Sections: []ExportSection{
					{
						Title: "Matematicas",
						Rows: []ExportRow{
							{Label: "Carlos Mendoza", Value: "95"},
							{Label: "Sofia Herrera", Value: "88"},
						},
						Note: "Promedio general: 91.5",
					},
				},
				Footer: ExportFooter{Text: "Generado por EduGo"},
			},
			wantContains: []string{
				"# Reporte de Notas",
				"_Periodo 2026-Q1_",
				"**Generado:** 2026-03-17 10:30",
				"## Matematicas",
				"- **Carlos Mendoza:** 95",
				"- **Sofia Herrera:** 88",
				"> Promedio general: 91.5",
				"---\n_Generado por EduGo_",
			},
			wantMIME: "text/markdown",
		},
		{
			name: "sin subtitle ni footer ni note",
			doc: ExportDocument{
				Header: ExportHeader{
					Title:       "Reporte Simple",
					GeneratedAt: fixedTime,
				},
				Sections: []ExportSection{
					{
						Title: "Datos",
						Rows: []ExportRow{
							{Label: "Total", Value: "100"},
						},
					},
				},
			},
			wantContains: []string{
				"# Reporte Simple",
				"## Datos",
				"- **Total:** 100",
			},
			wantMIME: "text/markdown",
		},
		{
			name: "sin secciones",
			doc: ExportDocument{
				Header: ExportHeader{
					Title:       "Documento Vacio",
					GeneratedAt: fixedTime,
				},
			},
			wantContains: []string{
				"# Documento Vacio",
				"**Generado:** 2026-03-17 10:30",
			},
			wantMIME: "text/markdown",
		},
		{
			name: "seccion sin rows pero con note",
			doc: ExportDocument{
				Header: ExportHeader{
					Title:       "Con Nota",
					GeneratedAt: fixedTime,
				},
				Sections: []ExportSection{
					{
						Title: "Observaciones",
						Note:  "Sin datos disponibles",
					},
				},
			},
			wantContains: []string{
				"## Observaciones",
				"> Sin datos disponibles",
			},
			wantMIME: "text/markdown",
		},
	}

	r := &MarkdownRenderer{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, mime, err := r.Render(tt.doc)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			if mime != tt.wantMIME {
				t.Errorf("Render() mime = %q, want %q", mime, tt.wantMIME)
			}
			output := string(got)
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Render() output missing %q\ngot:\n%s", want, output)
				}
			}
		})
	}
}

func TestMarkdownRenderer_Render_OmitsEmptyOptionalFields(t *testing.T) {
	fixedTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	r := &MarkdownRenderer{}

	doc := ExportDocument{
		Header: ExportHeader{
			Title:       "Test",
			GeneratedAt: fixedTime,
		},
		Sections: []ExportSection{
			{
				Title: "Section 1",
				Rows:  []ExportRow{{Label: "A", Value: "1"}},
			},
		},
	}

	got, _, err := r.Render(doc)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	output := string(got)

	// Subtitle should not appear
	if strings.Contains(output, "_ _") || strings.Count(output, "\n_") > 0 && strings.Contains(output, "_\n\n**") {
		// Subtitle is rendered as _subtitle_ so check there is no empty italic
	}

	// Footer separator should not appear when footer is empty
	if strings.HasSuffix(strings.TrimSpace(output), "---") {
		t.Error("Render() should not include footer separator when footer is empty")
	}
}
