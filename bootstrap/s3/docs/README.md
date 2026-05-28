# Bootstrap S3 — Documentacion tecnica

## Descripcion general

Sub-modulo que implementa la creacion de clientes AWS S3 usando AWS SDK v2 con credenciales estaticas, validacion de bucket y soporte para endpoints custom (LocalStack).

## Componentes principales

### Factory

```go
type Factory struct{}
```

Struct sin estado.

### CreateClient

```go
func (f *Factory) CreateClient(ctx context.Context, cfg bootstrap.S3Config) (*s3.Client, error)
```

Pasos:
1. Crea `StaticCredentialsProvider` con AccessKeyID y SecretAccessKey
2. Carga config AWS con region y credenciales
3. Si `cfg.Endpoint` no esta vacio, configura `BaseEndpoint` y `UsePathStyle`
4. Si `cfg.Bucket` no esta vacio, ejecuta `ValidateBucket` via `HeadBucket`
5. Retorna el cliente configurado

### CreatePresignClient

```go
func (f *Factory) CreatePresignClient(client *s3.Client) *s3.PresignClient
```

Wrapper sobre `s3.NewPresignClient(client)`.

### ValidateBucket

```go
func (f *Factory) ValidateBucket(ctx context.Context, client *s3.Client, bucket string) error
```

Ejecuta `HeadBucket` para verificar que el bucket existe y es accesible con las credenciales proporcionadas.

## S3Config (del modulo raiz bootstrap)

```go
type S3Config struct {
    Bucket          string
    Region          string
    AccessKeyID     string
    SecretAccessKey string
    Endpoint        string // Para LocalStack
    ForcePathStyle  bool   // Para LocalStack
}
```

## Flujos comunes

### 1. AWS S3 produccion

```go
client, err := factory.CreateClient(ctx, bootstrap.S3Config{
    Bucket:          "prod-bucket",
    Region:          "us-east-1",
    AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
    SecretAccessKey:  os.Getenv("AWS_SECRET_ACCESS_KEY"),
})
```

### 2. LocalStack desarrollo

```go
client, err := factory.CreateClient(ctx, bootstrap.S3Config{
    Bucket:          "dev-bucket",
    Region:          "us-east-1",
    AccessKeyID:     "test",
    SecretAccessKey:  "test",
    Endpoint:        "http://localhost:4566",
    ForcePathStyle:  true,
})
```

## Dependencias

### Internas
- `github.com/EduGoGroup/edugo-shared/bootstrap` — S3Config

### Externas
- `github.com/aws/aws-sdk-go-v2` — AWS SDK base
- `github.com/aws/aws-sdk-go-v2/config` — AWS config loader
- `github.com/aws/aws-sdk-go-v2/credentials` — Static credentials
- `github.com/aws/aws-sdk-go-v2/service/s3` — S3 service client

## Notas de diseño

- **Endpoint custom**: Permite usar LocalStack para desarrollo local sin cambiar el codigo del consumidor.
- **ForcePathStyle**: Requerido para LocalStack ya que no soporta virtual-hosted-style URLs.
- **ValidateBucket opcional**: Solo ejecuta si `cfg.Bucket` no esta vacio, permitiendo crear clientes sin bucket predefinido.
- **Ningun consumidor actual usa esta factory**: Worker usa `infrastructure.NewFactory()` y Mobile usa `storage.NewS3Client()`. Existe para futura adopcion o uso directo de S3.
