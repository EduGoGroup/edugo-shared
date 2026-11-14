# Changelog

Todos los cambios relevantes del repositorio `edugo-shared` deben registrarse aqui.

## [Unreleased]

### Added

- `govulncheck` v1.9.0 en la CI (`ci.yml`, matriz `validate`) y en `make vuln` / `make ci-local`; una vulnerabilidad alcanzable pone la CI en rojo (plan 054, H.2).
- Nueva base documental de fase 1 para el repositorio y sus modulos independientes.
- `README.md` raiz, `docs/` general y `CHANGELOG.md` raiz alineados con la estructura modular actual.
- `logger.ResolveOtelLevel(envVarValue, deploymentEnv) slog.Level` — resuelve el nivel del exporter OTel/Loki desde `OTEL_LOG_LEVEL` con fallback estricto `info` (DA-MPH-5).
- `logger.SlogConfig.OtelLevel` — campo nuevo para propagar el nivel OTel al provider.

### Changed

- **BREAKING** `tracer.NewSlogHandler` cambia firma: `NewSlogHandler(name string)` → `NewSlogHandler(name string, minLevel slog.Level)`. El segundo argumento aplica un filtro de nivel propio al branch OTLP, desacoplándolo de `LOGGING_LEVEL` (stdout). Migración: en el call-site, calcular `otelLevel := logger.ResolveOtelLevel(os.Getenv("OTEL_LOG_LEVEL"), cfg.Environment)` y pasarlo. Motivación en DA-MPH-5 (plan multi-platform-hardening).
