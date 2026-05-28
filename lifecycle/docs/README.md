# Lifecycle — Documentación técnica

Orquestador de ciclo de vida para recursos de infraestructura con startup secuencial y cleanup LIFO.

## Propósito

Proporcionar mecanismo robusto para registrar, iniciar y limpiar recursos de infraestructura en orden controlado, asegurando que la limpieza ocurra en orden inverso al registro (LIFO) incluso si algunos recursos fallan.

## Componentes principales

### Manager — Orquestador

Gestor central que mantiene lista de recursos y coordina startup/cleanup.

**Estructura:**
```go
type Manager struct {
    resources []Resource        // Lista de recursos registrados
    mu        sync.Mutex        // Protección de concurrencia
    logger    logger.Logger     // Logger opcional
    startTime time.Time         // Tiempo de inicio (para métricas)
}
```

**Métodos principales:**
- `NewManager(log logger.Logger) *Manager` — Constructor
- `Register(name string, startup func(ctx context.Context) error, cleanup func() error)` — Registrar recurso con startup y cleanup
- `RegisterSimple(name string, cleanup func() error)` — Registrar recurso solo con cleanup
- `Startup(ctx context.Context) error` — Ejecutar startup de todos los recursos en orden
- `Cleanup() error` — Ejecutar cleanup de todos los recursos en orden inverso (LIFO)
- `Count() int` — Cantidad de recursos registrados
- `Clear()` — Limpiar lista sin ejecutar cleanup (para testing)

### Resource — Definición de recurso

Estructura que define un recurso registrable con funciones de startup y cleanup.

**Estructura:**
```go
type Resource struct {
    Name    string                          // Identificador del recurso
    Startup func(ctx context.Context) error // Inicialización (opcional)
    Cleanup func() error                    // Limpieza (obligatorio si se registra)
}
```

## Flujos comunes

### 1. Registrar múltiples recursos al inicializar

```go
func setupInfrastructure(ctx context.Context, log logger.Logger) (*lifecycle.Manager, error) {
    mgr := lifecycle.NewManager(log)

    // Recurso con startup y cleanup
    mgr.Register("database",
        func(ctx context.Context) error {
            return db.Connect(ctx)
        },
        func() error {
            return db.Close()
        },
    )

    // Otro recurso
    mgr.Register("cache",
        func(ctx context.Context) error {
            return cache.Init(ctx)
        },
        func() error {
            return cache.Flush()
        },
    )

    // Recurso sin startup (ya está inicializado externamente)
    mgr.RegisterSimple("logger",
        func() error {
            return log.Close()
        },
    )

    // Ejecutar startup en orden de registro
    if err := mgr.Startup(ctx); err != nil {
        return nil, fmt.Errorf("infrastructure startup failed: %w", err)
    }

    return mgr, nil
}
```

### 2. Cleanup con manejo de errores acumulados

```go
func shutdown(mgr *lifecycle.Manager) {
    // Cleanup ejecuta en orden inverso, acumulando errores
    if err := mgr.Cleanup(); err != nil {
        log.Printf("cleanup warnings: %v", err)
        // Continuamos aunque haya errores; se reportan para observabilidad
    }
}
```

### 3. Startup fallido aborta el proceso

```go
// Si un recurso falla en startup, el proceso se detiene inmediatamente
// El error incluye el nombre del recurso que falló
if err := mgr.Startup(ctx); err != nil {
    // Error: "failed to startup resource database: connection refused"
    return nil, err
}
```

### 4. Consultar estado del manager

```go
// Cantidad de recursos
count := mgr.Count() // 3

// Limpiar para reutilizar en testing
mgr.Clear()
newCount := mgr.Count() // 0
```

## Arquitectura

Flujo de ciclo de vida:

```
1. NewManager(logger)
   ↓
2. Register(name, startup, cleanup)
   ├─ Agrega recurso a lista
   ├─ Logger trazabilidad
   └─ Retorna inmediatamente
   ↓
3. Startup(ctx)
   ├─ Itera en orden de registro
   ├─ Ejecuta startup de cada recurso
   ├─ Si falla: retorna error inmediatamente
   └─ Si todo ok: retorna nil
   ↓
4. [Aplicación ejecutándose]
   ↓
5. Cleanup()
   ├─ Itera en orden inverso (LIFO)
   ├─ Ejecuta cleanup de cada recurso
   ├─ Continúa incluso si falla
   ├─ Acumula todos los errores
   └─ Retorna error compuesto si hubo fallos
```

## Características

- **Startup secuencial**: Recursos se inicializan en orden de registro
- **Cleanup LIFO**: Orden inverso garantiza dependencias se limpian correctamente
- **Tolerancia a cleanup**: Continúa limpiando incluso si falla, para no dejar recursos abiertos
- **Thread-safe**: Protección con mutex para operaciones concurrentes
- **Logger opcional**: Trazabilidad sin obligatoriedad
- **Métricas**: Calcula duración de startup y cleanup

## Dependencias

- **Internas**: `logger` (opcional, solo para logging)
- **Externas**: Ninguna

## Testing

Suite de tests completa:
- Registro de recursos y validación de orden
- Startup exitoso y fallido
- Cleanup en orden inverso
- Acumulación de errores en cleanup
- Operaciones con lista vacía
- Thread-safety con race detector

Ejecutar:
```bash
make test          # Tests básicos
make test-race     # Tests con race detector
make check         # Tests + linting + format
```

## Notas de diseño

- **Módulo pequeño y estable**: Responsabilidad única — orquestar ciclo de vida
- **LIFO como contrato**: Inversión del orden de registro es lo que hace al módulo valioso
- **Tolerancia en cleanup**: Permite que recursos parcialmente inicializados se limpien completamente
- **Logger opcional**: Adición de trazabilidad sin acoplamiento fuerte
- **Sin framework específico**: Funciona con cualquier tipo de recurso que implemente las interfaces de startup/cleanup
