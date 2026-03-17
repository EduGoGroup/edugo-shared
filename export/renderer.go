package export

// Renderer convierte un Document a bytes.
type Renderer interface {
	Render(doc Document) ([]byte, string, error) // bytes, mime-type, error
	Extension() string                           // "md", "pdf"
}
