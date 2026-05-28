# Logger

Abstracción de logging estructurado con implementaciones en Zap, Logrus y Slog.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/logger
```

El módulo se versiona y se consume de forma independiente gracias a su `go.mod` propio.

## Quick Start

### Crear logger Zap con nivel y formato

```go
// JSON format, info level
logger := logger.NewZapLogger("info", "json")
defer logger.Sync()

// Emitir logs estructurados
logger.Info("usuario autenticado", "user_id", 123, "email", "user@example.com")
logger.Error("fallo la conexión", "host", "db.example.com", "error", err)
```

### Usar contexto con fields

```go
// Agregar campos contextuales que se incluyen en todos los logs subsecuentes
requestLogger := logger.With("request_id", "abc123", "user_id", 456)
requestLogger.Info("procesando solicitud")
requestLogger.Info("solicitud completada")
// Ambos logs incluyen request_id y user_id
```

### Adaptar Logrus existente

```go
// Si ya tienes una instancia Logrus, adaptarla a la interfaz común
logrusInstance := logrus.New()
logger := logger.NewLogrusLogger(logrusInstance)
logger.Info("migrando a interfaz común", "status", "ok")
```

### Usar Slog (recomendado para nuevos proyectos)

```go
// Factory moderna con configuración desde variables de entorno
slogLogger := logger.NewSlogProviderFromEnv().Logger()

// Helpers tipados para campos comunes
slogLogger.InfoContext(ctx, "operación iniciada",
    logger.WithRequestID("req-123"),
    logger.WithUserID("user-456"),
    logger.WithDuration(150*time.Millisecond),
)
```

## Componentes principales

- **Logger**: Interfaz común (Debug, Info, Warn, Error, Fatal, With, Sync)
- **ZapLogger**: Implementación con go.uber.org/zap (JSON o console)
- **LogrusLogger**: Adaptador para github.com/sirupsen/logrus
- **SlogProvider**: Factory moderna para *slog.Logger con campos base
- **SlogAdapter**: Adaptador de slog a interfaz Logger (migración gradual)
- **Context helpers**: Propagación de logger via context.Context
- **Fields helpers**: Constructores tipados para atributos (WithUserID, WithRequestID, etc.)

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build    # Compilar módulo
make test     # Ejecutar tests
make test-race # Tests con race detector
make check    # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- **Sin políticas de logging**: Módulo solo proporciona abstracción e implementaciones
- **Múltiples backends**: Zap, Logrus y Slog intercambiables según necesidad
- **Campos variadicos**: Key/value pairs para minimizar acoplamiento
- **Migración gradual**: SlogAdapter permite usar slog en nuevo código mientras se migra
