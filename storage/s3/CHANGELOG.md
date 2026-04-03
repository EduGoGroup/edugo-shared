# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/storage/s3` se registran aqui.

## [Unreleased]

## [0.101.0] - 2026-04-02

### Added

- `Client` struct implementando `storage.Client` para AWS S3 con retry y backoff exponencial.
- `NewClient(s3Client, bucket, opts...)` constructor con opciones funcionales.
- `WithMaxRetries(n)` y `WithBaseBackoff(d)` para configurar retry.
- Operaciones CRUD: Download, Upload, Delete, Exists, GetMetadata.
- `PresignClient` struct implementando `storage.PresignClient` para URLs pre-firmadas.
- `NewPresignClient(s3Client, bucket, expiry)` con expiry default de 15 minutos.
- GenerateUploadURL y GenerateDownloadURL para presigned URLs.
- Targets Makefile: build, test, check, lint, fmt, vet, tidy, deps, release.

### Dependencies

- `github.com/EduGoGroup/edugo-shared/storage` v0.101.0
- `github.com/aws/aws-sdk-go-v2` v1.41.4
- `github.com/aws/aws-sdk-go-v2/service/s3` v1.97.2
