package export

import "time"

// ExportDocument es el contrato estándar para cualquier exportación.
// Define la estructura completa de un documento exportable con encabezado,
// secciones de contenido y pie de página.
type ExportDocument struct {
	Header   ExportHeader    `json:"header"`
	Sections []ExportSection `json:"sections"`
	Footer   ExportFooter    `json:"footer,omitempty"`
}

// ExportHeader contiene los metadatos del encabezado del documento exportado,
// incluyendo título, subtítulo opcional y datos de generación.
type ExportHeader struct {
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle,omitempty"`
	GeneratedAt time.Time `json:"generated_at"`
	GeneratedBy string    `json:"generated_by,omitempty"`
}

// ExportSection representa una sección del documento con un título,
// filas de datos clave-valor y una nota opcional.
type ExportSection struct {
	Title string      `json:"title"`
	Rows  []ExportRow `json:"rows,omitempty"`
	Note  string      `json:"note,omitempty"`
}

// ExportRow representa una fila de datos con un par etiqueta-valor.
type ExportRow struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ExportFooter contiene el texto opcional del pie de página del documento.
type ExportFooter struct {
	Text string `json:"text,omitempty"`
}
