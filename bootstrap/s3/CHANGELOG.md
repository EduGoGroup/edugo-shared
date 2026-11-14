# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/bootstrap/s3` se registran aqui.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Migración a la banda `0.900.x` (ADR 0010).
- `go.mod`: **eliminado** el `replace github.com/EduGoGroup/edugo-shared/bootstrap => ../`; el módulo requiere el tag publicado `edugo-shared/bootstrap v0.900.0` (HU-E1.2). Sin cambios de API.
- Dependencias directas: `github.com/aws/aws-sdk-go-v2` `v1.41.5` → `v1.48.0`, `config` `v1.32.14` → `v1.33.18`, `credentials` `v1.19.14` → `v1.21.0`, `service/s3` `v1.98.0` → `v1.114.10`.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

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
