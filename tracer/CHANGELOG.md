# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/tracer` se registran aquí.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencias directas: `go.opentelemetry.io/contrib/bridges/otelslog` `v0.18.0` → `v0.21.0`, `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin` `v0.60.0` → `v0.72.0`, `go.opentelemetry.io/otel` `v1.43.0` → `v1.47.0`, `go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc` `v0.19.0` → `v0.23.0`, `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc` `v1.35.0` → `v1.47.0`, `go.opentelemetry.io/otel/exporters/stdout/stdouttrace` `v1.35.0` → `v1.47.0`, `go.opentelemetry.io/otel/sdk` `v1.43.0` → `v1.47.0`, `go.opentelemetry.io/otel/sdk/log` `v0.19.0` → `v1.47.0`, `go.opentelemetry.io/otel/trace` `v1.43.0` → `v1.47.0`.
- 31 dependencias indirectas suben con `go mod tidy` (ver `go.mod`).
- `otel/log` pasa de `v0.19.0` a `v1.47.0`: `global.SetLoggerProvider` (obsoleto) se reemplaza por `otel.SetLoggerProvider`, que es a donde reenvía. Sin cambio de comportamiento.
- Comentarios de documentación en `Config` y `Shutdown` (lint).
- El módulo entra a la banda `v0.900.x` (ADR 0010); el tag anterior `v0.1.x` queda como procedencia.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.
