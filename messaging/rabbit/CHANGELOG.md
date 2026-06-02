# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/messaging/rabbit` deben registrarse aqui.

## [Unreleased]

## [0.900.0] - 2026-06-02

### Changed

- Reversionado a la serie `0.900.x` para escapar del fantasma de checksum del tag `v0.1.0`
  (su contenido cambió tras la limpieza de historial del repo, dejando go.sum de consumidores
  desincronizado; ver bug 0022 / ADR 0015). No hay cambios de código ni de API respecto al
  contenido actual de `main`: es una versión fresca con hash limpio para los consumidores
  (`edugo-api-learning`, `edugo-worker`).

## [0.100.0] - 2026-04-02

### Removed

- Deleted `publisher_unit_test.go` (100% covered by integration tests)

### Changed

- Added `//go:build integration` tag to `publisher_test.go`, `connection_test.go`, `consumer_test.go`

## [Unreleased]

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

### Added

- Baseline de documentacion de fase 1 con `README.md`, `docs/README.md` y organizacion local por modulo.
