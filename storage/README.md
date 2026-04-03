# Storage

Interfaces para operaciones de almacenamiento de archivos. Implementacion S3 en sub-modulo.

## Arquitectura

```
storage/              # Interfaces (0 deps externas)
└── storage/s3/       # Implementacion AWS S3
```

## Instalacion

```bash
# Solo interfaces
go get github.com/EduGoGroup/edugo-shared/storage

# Implementacion S3
go get github.com/EduGoGroup/edugo-shared/storage/s3
```

## Interfaces

### Client (CRUD)

```go
type Client interface {
    Download(ctx, key) (io.ReadCloser, error)
    Upload(ctx, key, io.Reader) error
    Delete(ctx, key) error
    Exists(ctx, key) (bool, error)
    GetMetadata(ctx, key) (*FileMetadata, error)
}
```

### PresignClient (URLs pre-firmadas)

```go
type PresignClient interface {
    GenerateUploadURL(ctx, key) (url string, expiresAt time.Time, err error)
    GenerateDownloadURL(ctx, key) (url string, expiresAt time.Time, err error)
}
```

## Sub-modulos

| Sub-modulo | Descripcion | Consumidores |
|-----------|-------------|-------------|
| [s3](s3/) | AWS S3 con retry y presigning | Mobile, Worker |

## Navegacion

- [Changelog](CHANGELOG.md)

## Comandos disponibles

```bash
make build     # Compilar
make test      # Tests
make check     # Lint y validacion
```
