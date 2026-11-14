# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/storage/s3` se registran aqui.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Migración a la banda `0.900.x` (ADR 0010).
- `go.mod`: **eliminado** el `replace github.com/EduGoGroup/edugo-shared/storage => ../`; el módulo requiere el tag publicado `edugo-shared/storage v0.900.0` (HU-E1.2). Sin cambios de API.
- Dependencias directas: `github.com/aws/aws-sdk-go-v2` `v1.41.5` → `v1.48.0`, `service/s3` `v1.98.0` → `v1.114.10`.
- Lint 2.14.0: comentarios de doc en los métodos exportados. Sin cambio de comportamiento.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

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
