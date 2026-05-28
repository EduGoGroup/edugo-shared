# Logger — Documentación técnica

Abstracción de logging estructurado con múltiples implementaciones (Zap, Logrus, Slog) intercambiables.

## Propósito

Proporcionar interfaz común de logging que abstrae detalles de implementación, permitiendo múltiples backends (Zap, Logrus, Slog) y facilitar migración gradual entre ellos.

## Componentes principales

### Logger — Interfaz común

Define el contrato que todas las implementaciones deben cumplir.

**Métodos:**
- `Debug(msg string, fields ...interface{})` — Log de nivel debug
- `Info(msg string, fields ...interface{})` — Log de nivel info
- `Warn(msg string, fields ...interface{})` — Log de nivel warning
- `Error(msg string, fields ...interface{})` — Log de nivel error
- `Fatal(msg string, fields ...interface{})` — Log de nivel fatal y exit
- `With(fields ...interface{}) Logger` — Retornar nuevo logger con campos contextuales
- `Sync() error` — Sincronizar buffers (importante antes de exit)

**Características:**
- Campos variadicos key/value (alternancia: "key1", value1, "key2", value2, ...)
- `With()` retorna nuevo logger sin mutar el original
- Recomendado sincronizar en defer antes de exit

### ZapLogger — Implementación Zap

Backend usando go.uber.org/zap con formato JSON o console.

**Constructor:**
- `NewZapLogger(level, format string) Logger`
  - `level`: "debug", "info", "warn", "error", "fatal" (default: "info")
  - `format`: "json" o "console" (default: "console")

**Características:**
- Encoder personalizado con timestamp ISO8601, caller, y nivel lowercase
- Console format con colores cuando es formato console
- Sincronización automática en métodos individuales (SugaredLogger)

### LogrusLogger — Implementación Logrus

Adaptador que implementa Logger delegando a github.com/sirupsen/logrus.

**Constructor:**
- `NewLogrusLogger(logger *logrus.Logger) Logger`

**Características:**
- Convierte pares key/value a logrus.Fields
- `With()` retorna nueva instancia con logrus.Entry con campos acumulados
- Compatible con configuración existente de Logrus

### SlogProvider — Factory moderna

Factory que crea `*slog.Logger` con configuración desde variables de entorno.

**Constructores:**
- `NewSlogProvider(serviceName, env, version string, format LogFormat) *SlogProvider`
- `NewSlogProviderFromEnv() *SlogProvider` — Lee de variables de entorno

**Métodos:**
- `Logger() *slog.Logger` — Retornar logger configurado

**Campos base automáticos:**
- `service`: nombre del servicio
- `env`: ambiente (dev, staging, prod)
- `version`: versión de la aplicación

### SlogAdapter — Adaptador Slog

Implementa interfaz Logger delegando a `*slog.Logger` para migración gradual.

**Constructor:**
- `NewSlogAdapter(slogLogger *slog.Logger) Logger`

**Características:**
- Convierte campos variadicos a slog.Attr
- Preserva tipos en atributos
- Permite usar slog en código nuevo mientras se migra

### Context helpers — Propagación de logger

Funciones para propagar logger via context.Context.

**Funciones:**
- `NewContext(ctx context.Context, logger *slog.Logger) context.Context` — Agregar logger al contexto
- `FromContext(ctx context.Context) *slog.Logger` — Extraer logger del contexto
- `L(ctx context.Context) *slog.Logger` — Alias corto para FromContext

### Fields helpers — Constructores tipados

Helpers que retornan slog.Attr para campos comunes.

**Funciones:**
- `WithRequestID(id string) slog.Attr`
- `WithUserID(id string) slog.Attr`
- `WithCorrelationID(id string) slog.Attr`
- `WithError(err error) slog.Attr`
- `WithDuration(d time.Duration) slog.Attr`
- `WithComponent(name string) slog.Attr`
- `WithSchoolID(id string) slog.Attr`
- `WithRole(role string) slog.Attr`
- `WithResource(resource string) slog.Attr`
- `WithResourceID(id string) slog.Attr`
- `WithAction(action string) slog.Attr`
- `WithIP(ip string) slog.Attr`

## Flujos comunes

### 1. Crear y usar logger Zap

```go
func initLogger() logger.Logger {
    log := logger.NewZapLogger("info", "json")
    // En defer/cleanup:
    defer log.Sync()
    return log
}

func handleRequest(log logger.Logger, userID string) {
    log.Info("procesando solicitud", "user_id", userID)

    // Con contexto
    reqLog := log.With("request_id", uuid.New())
    reqLog.Debug("detalles de solicitud", "status", "processing")
}
```

### 2. Migrar desde Logrus existente

```go
// Código anterior (Logrus directo)
logrusLog := logrus.New()

// Migración: adaptar a interfaz común
log := logger.NewLogrusLogger(logrusLog)
log.Info("ya usamos interfaz común")
```

### 3. Usar Slog (recomendado)

```go
func initSlog(ctx context.Context) context.Context {
    // Factory con configuración desde ENV
    provider := logger.NewSlogProviderFromEnv()
    slogLogger := provider.Logger()

    // Propagar via context
    return logger.NewContext(ctx, slogLogger)
}

func handler(ctx context.Context, userID string) {
    // Extraer logger del contexto
    log := logger.FromContext(ctx)

    // Usar helpers tipados
    log.InfoContext(ctx, "usuario activo",
        logger.WithUserID(userID),
        logger.WithAction("login"),
        logger.WithDuration(50*time.Millisecond),
    )
}
```

### 4. Context propagation en goroutines

```go
func processAsync(ctx context.Context, items []Item) {
    log := logger.FromContext(ctx)

    for _, item := range items {
        go func(item Item) {
            // Logger propagado automáticamente via contexto
            itemLog := logger.FromContext(ctx)
            itemLog.InfoContext(ctx, "procesando ítem",
                logger.WithResourceID(item.ID),
            )
        }(item)
    }
}
```

## Arquitectura

Opciones de implementación por escenario:

```
Escenario 1: Nuevo proyecto
├─ Usar SlogProvider
├─ Propagar via context.Context
└─ Usar Fields helpers

Escenario 2: Migración desde Logrus
├─ Wrappear con NewLogrusLogger
├─ Mantener compatibilidad
└─ Migrar a Slog gradualmente

Escenario 3: Código existente Zap
├─ Usar NewZapLogger directo
├─ Mantener sincronización en defer
└─ Considerar Slog para nuevo código

Escenario 4: Múltiples backends
└─ Interface Logger abstrae diferencias
```

## Dependencias

- **Internas**: Ninguna
- **Externas**:
  - `go.uber.org/zap` (para ZapLogger)
  - `github.com/sirupsen/logrus` (para LogrusLogger)
  - `log/slog` (Go stdlib, para Slog)

## Testing

Suite de tests comprensiva:
- Creación de loggers (Zap, Logrus, Slog)
- Niveles de logging
- Campos contextuales y With()
- Sincronización
- Context propagation
- Fields helpers
- Benchmarks para cada implementación

Ejecutar:
```bash
make test          # Tests básicos
make test-race     # Tests con race detector
make check         # Tests + linting + format
```

## Notas de diseño

- **Sin políticas definidas**: Módulo no impone cómo logguear, solo la mecánica
- **Abstracción intercambiable**: Zap, Logrus y Slog son opciones equivalentes
- **Campos variadicos**: Minimiza acoplamiento en interfaces de log
- **With() inmutable**: Retorna nuevo logger, no muta el original
- **Sync importante**: Especialmente en Zap, Sync() es crítico antes de exit
- **Migración gradual**: SlogAdapter permite usar slog mientras se migra
- **Context-aware**: Propagación de logger via context es práctica moderna
