# Storage S3

Implementacion AWS S3 de las interfaces storage.Client y storage.PresignClient.

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/storage/s3
```

## Uso rapido

```go
import (
    bootstraps3 "github.com/EduGoGroup/edugo-shared/bootstrap/s3"
    storages3 "github.com/EduGoGroup/edugo-shared/storage/s3"
)

// Crear cliente S3 via bootstrap
s3Client, _ := bootstraps3.NewFactory().CreateClient(ctx, s3Config)

// CRUD operations
client := storages3.NewClient(s3Client, "my-bucket")
err := client.Upload(ctx, "docs/file.pdf", reader)
body, err := client.Download(ctx, "docs/file.pdf")

// Presigned URLs
presigner := storages3.NewPresignClient(s3Client, "my-bucket", 15*time.Minute)
url, expiresAt, err := presigner.GenerateUploadURL(ctx, "uploads/photo.jpg")
```

## API Publica

### Client (CRUD)
- `NewClient(s3Client, bucket, opts...) *Client` — CRUD client con retry.
- `WithMaxRetries(n)` — Maximo de reintentos (default: 3).
- `WithBaseBackoff(d)` — Backoff base (default: 100ms, exponencial).

### PresignClient (URLs)
- `NewPresignClient(s3Client, bucket, expiry) *PresignClient` — Presign client.
- Expiry default: 15 minutos.

## Navegacion

- [Changelog](CHANGELOG.md)

## Comandos disponibles

```bash
make build     # Compilar
make test      # Tests
make check     # Lint y validacion
```

## Dependencias

- `github.com/EduGoGroup/edugo-shared/storage` — Interfaces
- `github.com/aws/aws-sdk-go-v2/service/s3` — AWS S3 SDK
