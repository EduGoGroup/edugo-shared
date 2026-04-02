# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/storage` se registran aqui.

## [Unreleased]

## [0.101.0] - 2026-04-02

### Added

- Interface `Client` para operaciones CRUD de almacenamiento (Download, Upload, Delete, Exists, GetMetadata).
- Interface `PresignClient` para generacion de URLs pre-firmadas (GenerateUploadURL, GenerateDownloadURL).
- Struct `FileMetadata` con Key, Size, ContentType, LastModified, ETag.
- Sub-modulo `storage/s3` con implementacion para AWS S3.
- Documentacion completa en README.md y docs/README.md.
- Targets Makefile: build, test, check, lint, fmt, vet, tidy, deps, release.
