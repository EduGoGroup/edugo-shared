# Cache Redis

Módulo de bajo nivel para conectar a Redis con caché genérica respaldada por JSON y soporte de patrones de borrado.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/cache/redis
```

El módulo se descarga como `cache/redis`, principal consumo vía package `redis`.

## Quick Start

### Ejemplo 1: Conexión básica y obtener cliente Redis

```go
package main

import (
	"fmt"
	"log"

	"github.com/EduGoGroup/edugo-shared/cache/redis"
)

func main() {
	// Crear configuración
	cfg := redis.RedisConfig{
		URL: "redis://localhost:6379",
	}

	// Conectar a Redis
	client, err := redis.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	// Usar la conexión
	fmt.Println("Connected to Redis successfully")
}
```

### Ejemplo 2: Crear servicio de caché y guardar datos

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EduGoGroup/edugo-shared/cache/redis"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	cfg := redis.RedisConfig{
		URL: "redis://localhost:6379",
	}

	client, err := redis.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer client.Close()

	// Crear servicio de caché
	cache := redis.NewCacheService(client)

	// Guardar usuario en caché con TTL de 1 hora
	user := User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	ctx := context.Background()
	err = cache.Set(ctx, "user:1", user, 1*time.Hour)
	if err != nil {
		log.Fatalf("failed to set cache: %v", err)
	}

	fmt.Println("User cached successfully")
}
```

### Ejemplo 3: Obtener y deserializar datos del caché

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/EduGoGroup/edugo-shared/cache/redis"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	cfg := redis.RedisConfig{
		URL: "redis://localhost:6379",
	}

	client, err := redis.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer client.Close()

	cache := redis.NewCacheService(client)
	ctx := context.Background()

	// Obtener usuario del caché
	var user User
	err = cache.Get(ctx, "user:1", &user)
	if err != nil {
		log.Fatalf("failed to get cache: %v", err)
	}

	fmt.Printf("Retrieved user: %+v\n", user)
}
```

### Ejemplo 4: Borrar múltiples claves y patrones

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EduGoGroup/edugo-shared/cache/redis"
)

func main() {
	cfg := redis.RedisConfig{
		URL: "redis://localhost:6379",
	}

	client, err := redis.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer client.Close()

	cache := redis.NewCacheService(client)
	ctx := context.Background()

	// Guardar múltiples sesiones
	cache.Set(ctx, "session:user1", map[string]string{"token": "abc123"}, 24*time.Hour)
	cache.Set(ctx, "session:user2", map[string]string{"token": "def456"}, 24*time.Hour)
	cache.Set(ctx, "session:user3", map[string]string{"token": "ghi789"}, 24*time.Hour)

	// Borrar sesiones específicas
	err = cache.Delete(ctx, "session:user1", "session:user2")
	if err != nil {
		log.Printf("failed to delete keys: %v", err)
	}

	// Borrar todas las sesiones por patrón
	err = cache.DeleteByPattern(ctx, "session:*")
	if err != nil {
		log.Printf("failed to delete by pattern: %v", err)
	}

	fmt.Println("Cache cleared successfully")
}
```

## Componentes principales

- **RedisConfig**: Configuración de conexión (URL con soporte rediss:// para TLS)
- **ConnectRedis()**: Establece conexión a Redis con validación de conectividad
- **CacheService**: Interfaz genérica para operaciones de caché con JSON
- **NewCacheService()**: Crea un servicio de caché respaldado por Redis
- **Get()**: Obtiene y deserializa valores del caché
- **Set()**: Serializa a JSON y guarda con TTL
- **Delete()**: Borra claves específicas
- **DeleteByPattern()**: Borra claves por patrón usando SCAN

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

- **Interfaz genérica**: CacheService abstrae Redis con una API agnóstica
- **JSON automático**: Serialización/deserialización transparente en Set/Get
- **TTL flexible**: Configurable por operación, no global
- **Borrado por patrón**: DeleteByPattern usa SCAN para evitar bloqueos con KEYS
- **TLS soportado**: rediss:// URL para conexiones a Upstash y similares
- **Conexión validada**: ConnectRedis valida con Ping, no lazy connection
- **Contexto sensible**: Todas las operaciones respetan contexto y timeouts
