# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/storage` se registran aqui.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- El módulo entra a la banda `v0.900.x` (ADR 0010); el tag anterior `v0.1.x` queda como procedencia.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.101.0] - 2026-04-02

### Added

- Interface `Client` para operaciones CRUD de almacenamiento (Download, Upload, Delete, Exists, GetMetadata).
- Interface `PresignClient` para generacion de URLs pre-firmadas (GenerateUploadURL, GenerateDownloadURL).
- Struct `FileMetadata` con Key, Size, ContentType, LastModified, ETag.
- Sub-modulo `storage/s3` con implementacion para AWS S3.
- Documentacion completa en README.md y docs/README.md.
- Targets Makefile: build, test, check, lint, fmt, vet, tidy, deps, release.
