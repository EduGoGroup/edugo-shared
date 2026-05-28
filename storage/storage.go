package storage

import (
	"context"
	"io"
	"time"
)

// Client define la interfaz para operaciones CRUD de almacenamiento.
// Permite abstraer S3, MinIO, filesystem local, etc.
type Client interface {
	// Download descarga un archivo desde el storage.
	// El caller debe cerrar el ReadCloser retornado.
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Upload sube un archivo al storage.
	Upload(ctx context.Context, key string, content io.Reader) error

	// Delete elimina un archivo del storage.
	Delete(ctx context.Context, key string) error

	// Exists verifica si un archivo existe.
	Exists(ctx context.Context, key string) (bool, error)

	// GetMetadata obtiene metadatos de un archivo.
	GetMetadata(ctx context.Context, key string) (*FileMetadata, error)
}

// PresignClient define la interfaz para generar URLs pre-firmadas.
type PresignClient interface {
	// GenerateUploadURL crea una URL pre-firmada para subir un objeto.
	// Retorna la URL y el timestamp de expiracion.
	GenerateUploadURL(ctx context.Context, key string) (url string, expiresAt time.Time, err error)

	// GenerateDownloadURL crea una URL pre-firmada para descargar un objeto.
	// Retorna la URL y el timestamp de expiracion.
	GenerateDownloadURL(ctx context.Context, key string) (url string, expiresAt time.Time, err error)
}

// FileMetadata contiene informacion sobre un archivo en storage.
type FileMetadata struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
}
