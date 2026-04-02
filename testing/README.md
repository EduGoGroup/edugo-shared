# Testing

Infraestructura de testing basada en Testcontainers para PostgreSQL, MongoDB y RabbitMQ, expuesta principalmente via el package `containers`.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/testing
```

El módulo se descarga como `testing`, pero la API consumible está concentrada en el package `testing/containers`. Requiere Docker para ejecutar containers reales durante tests.

## Quick Start

### Crear configuración de containers con ConfigBuilder

```go
import (
    "context"
    "github.com/EduGoGroup/edugo-shared/testing/containers"
)

// Crear configuración fluida
config := containers.NewConfig().
    WithPostgres(true).
    WithMongoDB(false).
    WithRabbitMQ(false).
    WithNetwork("test-net")

// Obtener manager singleton
mgr, err := config.GetManager(context.Background())
if err != nil {
    log.Fatal(err)
}
defer mgr.Cleanup(context.Background())
```

### Obtener acceso al container de PostgreSQL

```go
// Manager expone getter para acceso a PostgreSQL
postgresContainer := mgr.PostgreSQL()
if postgresContainer == nil {
    log.Fatal("PostgreSQL not initialized")
}

// Obtener conexión GORM
db := postgresContainer.DB()

// Ejecutar operaciones
var user User
db.Where("id = ?", "user-123").First(&user)
```

### Reutilizar container entre múltiples tests

```go
var mgr *containers.Manager

// Setup (una sola vez para toda la suite)
func TestMain(m *testing.M) {
    ctx := context.Background()
    config := containers.NewConfig().
        WithPostgres(true).
        WithMongoDB(true)

    var err error
    mgr, err = config.GetManager(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer mgr.Cleanup(ctx)

    code := m.Run()
    os.Exit(code)
}

// Limpiar estado entre tests
func cleanupState(t *testing.T) {
    ctx := context.Background()

    // Truncar tablas PostgreSQL
    if err := mgr.TruncatePostgres(ctx, []string{"users", "schools"}); err != nil {
        t.Fatalf("truncate failed: %v", err)
    }

    // Dropear colecciones MongoDB
    if err := mgr.DropMongoDB(ctx, []string{"profiles", "settings"}); err != nil {
        t.Fatalf("drop failed: %v", err)
    }
}

func TestUserCreation(t *testing.T) {
    cleanupState(t)

    db := mgr.PostgreSQL().DB()
    // Test con estado limpio
}
```

### Esperar a que container esté listo con retries

```go
ctx := context.Background()

// WaitForHealthy reintenta hasta que el container está listo
if err := mgr.PostgreSQL().WaitForHealthy(ctx); err != nil {
    log.Fatalf("PostgreSQL failed to become healthy: %v", err)
}

// RetryOperation ejecuta operación con reintentos
result := containers.RetryOperation(func() (interface{}, error) {
    return mgr.PostgreSQL().DB().Raw("SELECT 1").Row().Scan(), nil
}, 5, 100*time.Millisecond)

if result.Err != nil {
    log.Fatal(result.Err)
}
```

## Componentes principales

- **ConfigBuilder**: Constructor fluido para configurar containers habilitados
- **Manager**: Singleton que gestiona ciclo de vida de todos los containers
- **PostgresContainer**: Wrapper para PostgreSQL con acceso GORM y utilities
- **MongoDBContainer**: Wrapper para MongoDB con acceso client
- **RabbitMQContainer**: Wrapper para RabbitMQ con acceso connection
- **Cleanup helpers**: TruncatePostgres, DropMongoDB, PurgeRabbitMQ
- **Health utilities**: WaitForHealthy, RetryOperation, ExecSQLFile

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build     # Compilar módulo (sin Docker)
make test      # Tests unitarios (sin Docker)
make test-race # Tests con race detector (sin Docker)
make check     # Validar (fmt, vet, lint, test)
make test-all  # Tests con integración Docker (requiere Docker)
```

## Notas de diseño

- **Singleton Manager**: Centraliza vida útil de containers para abaratar suites largas
- **ConfigBuilder fluido**: API clara para habilitar solo containers necesarios
- **Cleanup entre tests**: Helpers tipados (TruncatePostgres, DropMongoDB, PurgeRabbitMQ)
- **Sin framework específico**: Funciona con cualquier testing framework (testing, testify, ginkgo, etc.)
- **Health checks**: Retries automáticos para esperar a que containers estén listos
- **Docker agnóstico**: Testcontainers abstrae detalles de Docker
