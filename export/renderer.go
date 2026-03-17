package export

// Renderer convierte un ExportDocument a bytes.
type Renderer interface {
	Render(doc ExportDocument) ([]byte, string, error) // bytes, mime-type, error
	Extension() string                                 // "md", "pdf"
}
