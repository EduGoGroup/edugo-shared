# Storage S3 — Documentacion tecnica

## Descripcion general

Sub-modulo que implementa las interfaces `storage.Client` y `storage.PresignClient` usando AWS S3 SDK v2.

## Componentes principales

### Client

```go
type Client struct {
    s3Client    *s3.Client
    bucket      string
    maxRetries  int           // default: 3
    baseBackoff time.Duration // default: 100ms
}
```

Implementa `storage.Client`. Incluye retry con backoff exponencial para Download y Upload.

**Retry logic**: backoff = baseBackoff * 2^attempt. Default: 100ms, 200ms, 400ms.

### PresignClient

```go
type PresignClient struct {
    presigner *s3.PresignClient
    bucket    string
    expiry    time.Duration // default: 15m
}
```

Implementa `storage.PresignClient`. Genera URLs pre-firmadas para upload y download.

## Flujos comunes

### 1. Worker con CRUD + validacion de dominio

```go
// Crear S3 client via bootstrap
s3Client, _ := bootstraps3.NewFactory().CreateClient(ctx, cfg)

// Crear storage client
client := storages3.NewClient(s3Client, cfg.Bucket, storages3.WithMaxRetries(5))

// El worker agrega su propia validacion de dominio encima
body, err := client.Download(ctx, key)
// ... validar extension, tamano, content-type a nivel de dominio ...
```

### 2. Mobile con presigned URLs

```go
s3Client, _ := bootstraps3.NewFactory().CreateClient(ctx, cfg)
presigner := storages3.NewPresignClient(s3Client, cfg.Bucket, 15*time.Minute)

url, expiresAt, _ := presigner.GenerateUploadURL(ctx, "materials/"+materialID+"/original")
// Retornar URL al cliente mobile para upload directo
```

## Dependencias

### Internas
- `github.com/EduGoGroup/edugo-shared/storage` — Client, PresignClient, FileMetadata

### Externas
- `github.com/aws/aws-sdk-go-v2/service/s3` — AWS S3

## Notas de diseño

- **Sin validacion de dominio**: El client no valida extensiones, tamanos ni content-types. Eso es responsabilidad del consumidor (worker valida PDF, mobile no valida).
- **Retry solo en Download y Upload**: Delete, Exists y GetMetadata son idempotentes y rapidos, no necesitan retry.
- **PresignClient separado de Client**: Diferentes casos de uso, diferentes dependencias internas (PresignClient usa s3.PresignClient, Client usa s3.Client directo).
- **Recibe *s3.Client inyectado**: No crea su propio cliente AWS. Usa bootstrap/s3 o cualquier otro mecanismo para crear el cliente.
