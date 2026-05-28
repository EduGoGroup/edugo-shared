# Bootstrap S3

Factory para crear clientes AWS S3 con credenciales estaticas, validacion de bucket y soporte para endpoints custom (LocalStack).

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/bootstrap/s3
```

## Uso rapido

```go
import (
    "github.com/EduGoGroup/edugo-shared/bootstrap"
    s3bootstrap "github.com/EduGoGroup/edugo-shared/bootstrap/s3"
)

factory := s3bootstrap.NewFactory()

// AWS S3
client, err := factory.CreateClient(ctx, bootstrap.S3Config{
    Bucket:          "my-bucket",
    Region:          "us-east-1",
    AccessKeyID:     "AKIA...",
    SecretAccessKey:  "secret",
})

// LocalStack
client, err := factory.CreateClient(ctx, bootstrap.S3Config{
    Bucket:          "local-bucket",
    Region:          "us-east-1",
    AccessKeyID:     "test",
    SecretAccessKey:  "test",
    Endpoint:        "http://localhost:4566",
    ForcePathStyle:  true,
})

// URLs pre-firmadas
presignClient := factory.CreatePresignClient(client)
```

## API Publica

- `NewFactory() *Factory` — Crea una nueva factory de S3.
- `CreateClient(ctx, S3Config) (*s3.Client, error)` — Cliente S3 con credenciales y validacion de bucket.
- `CreatePresignClient(*s3.Client) *s3.PresignClient` — Cliente para URLs pre-firmadas.
- `ValidateBucket(ctx, *s3.Client, string) error` — Verifica que el bucket existe y es accesible.

## Navegacion

- [Documentacion tecnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Comandos disponibles

```bash
make build     # Compilar el modulo
make test      # Ejecutar tests
make check     # Lint y validacion
```

## Dependencias

- `github.com/EduGoGroup/edugo-shared/bootstrap` — Config structs
- `github.com/aws/aws-sdk-go-v2` — AWS SDK base
- `github.com/aws/aws-sdk-go-v2/service/s3` — Servicio S3
