package export

import "time"

// Document es el contrato estándar para cualquier exportación.
// Define la estructura completa de un documento exportable con encabezado,
// secciones de contenido y pie de página.
type Document struct {
	Header   Header    `json:"header"`
	Sections []Section `json:"sections"`
	Footer   Footer    `json:"footer,omitempty"`
}

// Header contiene los metadatos del encabezado del documento exportado,
// incluyendo título, subtítulo opcional y datos de generación.
type Header struct {
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle,omitempty"`
	GeneratedAt time.Time `json:"generated_at"`
	GeneratedBy string    `json:"generated_by,omitempty"`
}

// Section representa una sección del documento con un título,
// filas de datos clave-valor y una nota opcional.
type Section struct {
	Title string `json:"title"`
	Rows  []Row  `json:"rows,omitempty"`
	Note  string `json:"note,omitempty"`
}

// Row representa una fila de datos con un par etiqueta-valor.
type Row struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Footer contiene el texto opcional del pie de página del documento.
type Footer struct {
	Text string `json:"text,omitempty"`
}
