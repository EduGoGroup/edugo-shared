# Database PostgreSQL

Capa de bajo nivel para conectar, validar y manejar transacciones en PostgreSQL con soporte nativo (sql.DB) y GORM.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/database/postgres
```

El módulo se descarga como `database/postgres`, principal consumo vía package `postgres`.

## Quick Start

### Ejemplo 1: Conexión básica con sql.DB

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/EduGoGroup/edugo-shared/database/postgres"
)

func main() {
	// Crear configuración
	cfg := postgres.DefaultConfig()
	cfg.Host = "db.example.com"
	cfg.User = "app_user"
	cfg.Password = "secure_password"
	cfg.Database = "edugo_prod"
	cfg.SSLMode = "require"

	// Conectar a PostgreSQL
	db, err := postgres.Connect(&cfg)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer postgres.Close(db)

	// Usar la conexión
	fmt.Println("Connected to PostgreSQL successfully")
}
```

### Ejemplo 2: Conexión GORM con logger personalizado

```go
package main

import (
	"fmt"
	"log"

	"github.com/EduGoGroup/edugo-shared/database/postgres"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := postgres.DefaultConfig()
	cfg.Host = "postgres.local"
	cfg.Database = "schools"

	// Conectar con logger GORM (opcional)
	db, err := postgres.ConnectGORM(&cfg) // Sin logger
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer db.Close()

	// Usar GORM db
	fmt.Println("GORM connection established")
}
```

### Ejemplo 3: Transacciones con manejo de errores

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/EduGoGroup/edugo-shared/database/postgres"
)

func main() {
	cfg := postgres.DefaultConfig()
	db, err := postgres.Connect(&cfg)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer postgres.Close(db)

	// Ejecutar transacción
	err = postgres.WithTransaction(context.Background(), db, func(tx *sql.Tx) error {
		// INSERT
		result, err := tx.ExecContext(context.Background(),
			"INSERT INTO users (name, email) VALUES ($1, $2)",
			"John Doe", "john@example.com",
		)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}

		id, _ := result.LastInsertId()
		fmt.Printf("Created user ID: %d\n", id)
		return nil
	})

	if err != nil {
		log.Printf("transaction failed: %v", err)
	}
}
```

### Ejemplo 4: Health checks periódicos y monitoreo de pool

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/EduGoGroup/edugo-shared/database/postgres"
)

func main() {
	cfg := postgres.DefaultConfig()
	db, err := postgres.Connect(&cfg)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer postgres.Close(db)

	// Health check periódico
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := postgres.HealthCheck(db); err != nil {
			log.Printf("health check failed: %v", err)
			continue
		}

		// Obtener estadísticas del pool
		stats := postgres.GetStats(db)
		fmt.Printf("Pool stats: OpenConnections=%d, InUse=%d, Idle=%d\n",
			stats.OpenConnections, stats.InUse, stats.Idle)
	}
}
```

## Componentes principales

- **Config**: Configuración de conexión (host, user, password, pool, SSL, timeouts)
- **DefaultConfig()**: Valores por defecto para desarrollo local
- **Connect()**: Establece conexión sql.DB nativa con validación
- **ConnectGORM()**: Establece conexión GORM con configuración de pool
- **HealthCheck()**: Verifica estado de la conexión
- **GetStats()**: Obtiene estadísticas del pool de conexiones
- **Close()**: Cierra la conexión
- **WithTransaction()**: Ejecuta código dentro de transacción con rollback/commit automático
- **WithTransactionIsolation()**: Transacción con nivel de aislamiento específico

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build     # Compilar
make test      # Tests unitarios
make test-race # Race detector
make check     # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- **Dual API**: Soporte para sql.DB nativo y GORM, no decide por el usuario
- **Pool configurable**: MaxConnections, MaxIdleConnections, MaxLifetime adaptables a carga
- **Validación de conectividad**: Connect valida inmediatamente con Ping, no lazy
- **Transacciones seguras**: WithTransaction rollback automático en error o panic
- **Schema search_path**: SearchPath configurable para múltiples esquemas
- **SSL flexible**: SSLMode soporta disable/require/verify-ca/verify-full
- **Timeouts defensivos**: ConnectTimeout y health check timeouts evitan bloqueos indefinidos
- **Sin abstracción**: Solo conecta y maneja transacciones, no define repositorios
