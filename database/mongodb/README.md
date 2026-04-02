# Database MongoDB

Módulo de bajo nivel para conectar, validar y cerrar conexiones a MongoDB con configuración defensiva y health checks.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/database/mongodb
```

El módulo se descarga como `database/mongodb`, principal consumo vía package `mongodb`.

## Quick Start

### Ejemplo 1: Configuración básica y conexión

```go
package main

import (
	"fmt"
	"log"

	"github.com/EduGoGroup/edugo-shared/database/mongodb"
)

func main() {
	// Crear configuración con defaults
	cfg := mongodb.DefaultConfig()
	cfg.URI = "mongodb://localhost:27017"
	cfg.Database = "edugo"
	cfg.Timeout = 10 * time.Second

	// Conectar a MongoDB
	client, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer mongodb.Close(client)

	// Usar la conexión
	fmt.Println("Connected to MongoDB successfully")
}
```

### Ejemplo 2: Obtener database y validar conexión

```go
package main

import (
	"fmt"
	"log"

	"github.com/EduGoGroup/edugo-shared/database/mongodb"
)

func main() {
	cfg := mongodb.DefaultConfig()
	cfg.URI = "mongodb://user:password@mongodb.example.com:27017"
	cfg.Database = "schools_db"

	client, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer mongodb.Close(client)

	// Obtener database
	db := mongodb.GetDatabase(client, cfg.Database)

	// Validar health
	if err := mongodb.HealthCheck(client); err != nil {
		log.Fatalf("health check failed: %v", err)
	}

	fmt.Printf("Connected to database: %s\n", db.Name())
}
```

### Ejemplo 3: Configuración con pool personalizado

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/EduGoGroup/edugo-shared/database/mongodb"
)

func main() {
	// Configurar pool para producción
	cfg := mongodb.Config{
		URI:         "mongodb+srv://admin:password@cluster.mongodb.net/",
		Database:    "production",
		Timeout:     15 * time.Second,
		MaxPoolSize: 200,  // Conexiones máximas
		MinPoolSize: 50,   // Conexiones mínimas siempre activas
	}

	client, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer mongodb.Close(client)

	fmt.Println("Production MongoDB pool initialized")
}
```

### Ejemplo 4: Health checks periódicos en aplicación

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/EduGoGroup/edugo-shared/database/mongodb"
)

func main() {
	cfg := mongodb.DefaultConfig()
	client, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer mongodb.Close(client)

	// Health check periódico
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := mongodb.HealthCheck(client); err != nil {
			log.Printf("health check failed: %v", err)
			// Re-conectar o alertar
			continue
		}
		fmt.Println("health check: OK")
	}
}
```

## Componentes principales

- **Config**: Configuración de conexión (URI, database, timeouts, pool sizes)
- **DefaultConfig()**: Valores por defecto para desarrollo local
- **Connect()**: Establece conexión a MongoDB con validación
- **GetDatabase()**: Obtiene instancia de base de datos del cliente
- **HealthCheck()**: Verifica estado de la conexión
- **Close()**: Cierra conexión de forma controlada

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

- **Sin abstracción de colecciones**: Módulo de bajo nivel, solo conecta y vigila salud
- **Pool configurable**: MaxPoolSize y MinPoolSize adaptables a carga esperada
- **Validación de conectividad**: Connect valida contra read preference Primary
- **Timeouts defensivos**: DefaultTimeout y health check timeouts evitan bloqueos indefinidos
- **Cierre ordenado**: Close usa timeout configurado para desconexión controlada
