package export

import "time"

// ExportDocument es el contrato estandar para cualquier exportacion.
type ExportDocument struct {
	Header   ExportHeader    `json:"header"`
	Sections []ExportSection `json:"sections"`
	Footer   ExportFooter    `json:"footer,omitempty"`
}

type ExportHeader struct {
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle,omitempty"`
	GeneratedAt time.Time `json:"generated_at"`
	GeneratedBy string    `json:"generated_by,omitempty"`
}

type ExportSection struct {
	Title string      `json:"title"`
	Rows  []ExportRow `json:"rows,omitempty"`
	Note  string      `json:"note,omitempty"`
}

type ExportRow struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type ExportFooter struct {
	Text string `json:"text,omitempty"`
}
