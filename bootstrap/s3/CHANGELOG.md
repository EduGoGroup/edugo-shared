# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/bootstrap/s3` se registran aqui.

## [Unreleased]

## [0.101.0] - 2026-04-02

### Added

- Factory S3 con soporte para endpoints custom (LocalStack).
- `NewFactory()` constructor para crear instancias de la factory.
- `CreateClient(ctx, S3Config) (*s3.Client, error)` con credenciales estaticas y validacion de bucket.
- `CreatePresignClient(*s3.Client) *s3.PresignClient` para URLs pre-firmadas.
- `ValidateBucket(ctx, *s3.Client, string) error` para verificar acceso al bucket.
- Soporte para `S3Config.Endpoint` y `ForcePathStyle` para LocalStack.
- Targets Makefile: build, test, check, lint, fmt, vet, tidy, deps, release.

### Dependencies

- `github.com/EduGoGroup/edugo-shared/bootstrap` v0.101.0
- `github.com/aws/aws-sdk-go-v2` v1.41.4
- `github.com/aws/aws-sdk-go-v2/service/s3` v1.97.2
