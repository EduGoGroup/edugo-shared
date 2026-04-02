# Lifecycle

Manager de ciclo de vida para recursos de infraestructura con startup secuencial y cleanup LIFO.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/lifecycle
```

El módulo se versionan y consume de forma independiente gracias a su `go.mod` propio.

## Quick Start

### Crear manager y registrar recursos

```go
// Crear manager con logger
mgr := lifecycle.NewManager(log)

// Registrar recurso con startup y cleanup
mgr.Register("database",
    func(ctx context.Context) error {
        return db.Connect(ctx)
    },
    func() error {
        return db.Close()
    },
)

// Registrar recurso sin startup (ya está inicializado)
mgr.RegisterSimple("cache",
    func() error {
        return cache.Close()
    },
)
```

### Ejecutar startup y cleanup en orden

```go
// Startup ejecuta en orden de registro
if err := mgr.Startup(ctx); err != nil {
    return fmt.Errorf("startup failed: %w", err)
}

// Cleanup ejecuta en orden inverso (LIFO), acumulando errores
if err := mgr.Cleanup(); err != nil {
    log.Printf("cleanup errors: %v", err)
}
```

### Consultar estado

```go
// Cantidad de recursos registrados
count := mgr.Count()

// Limpiar lista (sin ejecutar cleanup, útil para testing)
mgr.Clear()
```

## Componentes principales

- **Manager**: Orquestador de ciclo de vida con mutex (startup ordenado, cleanup LIFO)
- **Resource**: Estructura con nombre, startup opcional y cleanup
- **Métodos**: Register, RegisterSimple, Startup, Cleanup, Count, Clear

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

- **Módulo pequeño y estable**: Orientado a coordinación in-process
- **Valor principal**: Orden de cleanup (LIFO) y agregación de errores
- **Startup**: Aborta si un recurso falla
- **Cleanup**: Continúa incluso si fallan recursos, acumulando todos los errores
