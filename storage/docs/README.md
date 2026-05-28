# Storage — Documentacion tecnica

## Descripcion general

Modulo raiz que define interfaces para operaciones de almacenamiento de archivos. No tiene dependencias externas — solo stdlib de Go.

## Componentes principales

### Client

```go
type Client interface {
    Download(ctx context.Context, key string) (io.ReadCloser, error)
    Upload(ctx context.Context, key string, content io.Reader) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    GetMetadata(ctx context.Context, key string) (*FileMetadata, error)
}
```

Interface para operaciones CRUD sobre archivos. El caller de Download debe cerrar el ReadCloser retornado.

### PresignClient

```go
type PresignClient interface {
    GenerateUploadURL(ctx context.Context, key string) (url string, expiresAt time.Time, err error)
    GenerateDownloadURL(ctx context.Context, key string) (url string, expiresAt time.Time, err error)
}
```

Interface para generacion de URLs pre-firmadas. Util para delegar uploads/downloads al cliente (browser, mobile app).

### FileMetadata

```go
type FileMetadata struct {
    Key          string
    Size         int64
    ContentType  string
    LastModified string
    ETag         string
}
```

## Dependencias

Ninguna. Solo stdlib de Go.

## Notas de diseño

- **Interfaces sin implementacion**: El modulo raiz solo define contratos. Las implementaciones viven en sub-modulos (storage/s3).
- **Client vs PresignClient**: Son interfaces separadas porque los casos de uso son diferentes. Mobile solo necesita presigning; Worker necesita CRUD completo.
- **Sin validacion de dominio**: Las interfaces no imponen restricciones de tipo de archivo o tamano. Eso es responsabilidad del consumidor.
