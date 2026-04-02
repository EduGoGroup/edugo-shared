# Cache Redis - Documentación Técnica

Documentación completa de componentes, métodos y patrones de integración del módulo de caché Redis.

## Componentes

### RedisConfig

Estructura que modela la configuración para conectarse a Redis.

```go
type RedisConfig struct {
	URL string // URI de conexión Redis
}
```

**Campos:**
- `URL`: Connection string Redis en formato `redis://[user:password@]host:port[/db]`
  - Soporta `redis://` para conexiones estándar
  - Soporta `rediss://` para TLS (ej: Upstash)
  - Puede incluir database index (ej: `/0`, `/1`)

**Ejemplos:**
- `redis://localhost:6379` - Local development
- `redis://password@redis.example.com:6380` - Con credenciales
- `redis://user:password@redis.example.com:6379/1` - Con database específica
- `rediss://default:password@abc123.upstash.io:6379` - Upstash con TLS

### ConnectRedis()

Establece conexión a Redis, parsea URL y valida conectividad.

```go
func ConnectRedis(cfg RedisConfig) (*goredis.Client, error)
```

**Comportamiento:**
1. Parsea URL con `redis.ParseURL(cfg.URL)`
2. Crea cliente con `redis.NewClient(opts)`
3. Valida conectividad con `Ping(ctx)` (timeout 5 segundos)
4. Retorna cliente si Ping exitoso, error si falla

**Errores:**
- `"parsing redis URL: %w"` - URL inválida o formato incorrecto
- `"pinging redis: %w"` - Conexión fallida (host no accesible, credenciales incorrectas)

**Notas:**
- El cliente es thread-safe, seguro compartir entre goroutines
- Ping se ejecuta con timeout de 5 segundos
- La URL parseada puede incluir timeouts adicionales vía parámetros (ej: `?dial_timeout=10s`)

### CacheService

Interfaz genérica para operaciones de caché con JSON.

```go
type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}
```

**Métodos:**
- `Get(ctx, key, dest)`: Obtiene valor deserializado en dest
- `Set(ctx, key, value, ttl)`: Guarda valor serializado con TTL
- `Delete(ctx, keys...)`: Borra una o más claves específicas
- `DeleteByPattern(ctx, pattern)`: Borra claves que coinciden patrón

### NewCacheService()

Crea un servicio de caché respaldado por Redis.

```go
func NewCacheService(client *goredis.Client) CacheService
```

**Parámetros:**
- `client`: Cliente Redis conectado

**Retorno:**
- `CacheService`: Interfaz para operaciones de caché

**Notas:**
- No valida que el cliente esté conectado (asume client válido)
- El cliente debe ser el retorno de `ConnectRedis()`

### Get()

Obtiene valor del caché y deserializa de JSON.

```go
func (s *redisCacheService) Get(ctx context.Context, key string, dest interface{}) error
```

**Comportamiento:**
1. Obtiene valor en bytes con `client.Get(ctx, key).Result()`
2. Deserializa JSON en `dest` con `json.Unmarshal`
3. Retorna error si clave no existe (redis.Nil) o JSON inválido

**Errores:**
- `redis.Nil`: Clave no existe
- `json.SyntaxError`: Valor no es JSON válido
- Error de conexión si Redis no está disponible

**Notas:**
- Respeta contexto (timeout, cancellation)
- `dest` debe ser un puntero (ej: `&user`, `&stringValue`)
- No tiene expiración lazy (usa TTL de Redis)

### Set()

Serializa a JSON y guarda en caché con TTL.

```go
func (s *redisCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
```

**Comportamiento:**
1. Serializa `value` a JSON con `json.Marshal`
2. Guarda en Redis con TTL vía `client.Set(ctx, key, data, ttl)`
3. Retorna error si serialización o escritura falla

**Errores:**
- `json.UnsupportedTypeError`: Tipo no puede serializarse (ej: channel)
- `"marshaling cache value: %w"` - Error durante serialización
- Error de conexión si Redis no está disponible

**Notas:**
- TTL es por entrada, no global
- TTL `0` es permitido (clave indefinida en Redis)
- TTL negativo se trata como `0` en Redis (no expiración)

### Delete()

Borra una o más claves específicas.

```go
func (s *redisCacheService) Delete(ctx context.Context, keys ...string) error
```

**Comportamiento:**
1. Valida que `keys` no esté vacío (retorna nil si está)
2. Ejecuta `client.Del(ctx, keys...)` para borrar todas las claves
3. Retorna error si operación falla

**Errores:**
- Error de conexión si Redis no está disponible

**Notas:**
- Seguro pasar cero claves (no-op)
- Retorna nil si algunas claves no existen (no es error)
- Operación atómica en Redis

### DeleteByPattern()

Borra claves que coinciden con patrón usando SCAN.

```go
func (s *redisCacheService) DeleteByPattern(ctx context.Context, pattern string) error
```

**Comportamiento:**
1. Usa `client.Scan(ctx, 0, pattern, 100)` para iterar claves coincidentes
2. Acumula claves en slice
3. Ejecuta `client.Del(ctx, keys...)` para borrar todas
4. Retorna error si SCAN o DELETE falla

**Patrones:**
- `*` - Todas las claves
- `session:*` - Claves que empiezan con "session:"
- `*:user:*` - Claves con "user" en el medio
- `cache:item:*:metadata` - Wildcard en posiciones específicas

**Errores:**
- `"scanning keys: %w"` - Error durante iteración con SCAN
- Error de conexión si Redis no está disponible

**Notas:**
- Usa SCAN (no bloquea Redis) en lugar de KEYS
- Batch size de 100 claves para eficiencia
- Seguro para patrones amplios (ej: `*`)
- No es atómico: SCAN + DEL pueden interleaved con otras operaciones

## Constantes

```go
// No hay constantes públicas definidas en este módulo.
// Timeout de Ping está hardcoded: 5 segundos
// Batch size de SCAN está hardcoded: 100
```

## Flujos comunes

### Flujo 1: Inicialización en startup

```go
func initCache(url string) (redis.CacheService, error) {
	cfg := redis.RedisConfig{
		URL: url,
	}

	client, err := redis.ConnectRedis(cfg)
	if err != nil {
		return nil, fmt.Errorf("redis startup failed: %w", err)
	}

	cache := redis.NewCacheService(client)
	log.Printf("Connected to Redis: %s", url)
	return cache, nil
}

func main() {
	cache, err := initCache("redis://localhost:6379")
	if err != nil {
		log.Fatalf("startup failed: %v", err)
	}

	// Usar cache en aplicación...
}
```

### Flujo 2: Caché de objetos con serialización JSON

```go
type CacheEntry struct {
	UserID    int       `json:"user_id"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func cacheUserData(ctx context.Context, cache redis.CacheService, userID int, data string) error {
	entry := CacheEntry{
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now(),
	}

	// Guardar con TTL de 1 hora
	key := fmt.Sprintf("user:%d:data", userID)
	return cache.Set(ctx, key, entry, 1*time.Hour)
}

func getUserFromCache(ctx context.Context, cache redis.CacheService, userID int) (*CacheEntry, error) {
	key := fmt.Sprintf("user:%d:data", userID)
	var entry CacheEntry
	if err := cache.Get(ctx, key, &entry); err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("cache get failed: %w", err)
	}
	return &entry, nil
}
```

### Flujo 3: Invalidación de caché por patrón

```go
func invalidateUserCache(ctx context.Context, cache redis.CacheService, userID int) error {
	// Opción 1: Borrar clave específica
	pattern := fmt.Sprintf("user:%d:*", userID)
	return cache.DeleteByPattern(ctx, pattern)
}

func invalidateSessionCache(ctx context.Context, cache redis.CacheService) error {
	// Borrar todas las sesiones
	return cache.DeleteByPattern(ctx, "session:*")
}

func clearCache(ctx context.Context, cache redis.CacheService) error {
	// Vaciar completamente Redis
	return cache.DeleteByPattern(ctx, "*")
}
```

### Flujo 4: Manejo de errores y timeouts

```go
func getWithFallback(ctx context.Context, cache redis.CacheService, key string, fallback interface{}) error {
	// Crear contexto con timeout para cache.Get
	cacheCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var value interface{}
	err := cache.Get(cacheCtx, key, &value)
	if err == redis.Nil {
		// Cache miss: usar fallback
		log.Printf("cache miss for key: %s", key)
		return nil // fallback es válido
	}
	if err != nil {
		// Error de conexión o desserialización
		log.Printf("cache error (using fallback): %v", err)
		return nil // fallback es válido
	}

	// Cache hit
	*fallback.(*interface{}) = value
	return nil
}

func setWithRetry(ctx context.Context, cache redis.CacheService, key string, value interface{}, ttl time.Duration) error {
	// Retry una sola vez si falla
	err := cache.Set(ctx, key, value, ttl)
	if err != nil {
		log.Printf("cache set failed, retrying: %v", err)
		return cache.Set(ctx, key, value, ttl)
	}
	return nil
}
```

## Arquitectura

```
RedisConfig (URL)
    ↓
ConnectRedis()
    ├─ ParseURL
    ├─ redis.NewClient(opts)
    ├─ Ping validation (5s timeout)
    └─ Return *redis.Client
    ↓
NewCacheService(client)
    └─ Return CacheService interface
    ↓
Operations
    ├─ Get: Redis GET + JSON Unmarshal
    ├─ Set: JSON Marshal + Redis SET + TTL
    ├─ Delete: Redis DEL multiple keys
    └─ DeleteByPattern: Redis SCAN + DEL
    ↓
client.Close()
    └─ Cierre ordenado
```

**Ciclo de vida:**
1. **Init**: RedisConfig{URL: "..."}
2. **Connect**: ConnectRedis + Ping validation
3. **Service**: NewCacheService(client)
4. **Operations**: Get, Set, Delete, DeleteByPattern
5. **Close**: client.Close() en defer

## Dependencias

**Internas:**
- Ninguna en producción (módulo autocontendido)

**Externas:**
- `github.com/redis/go-redis/v9`: Cliente oficial de Redis
  - `redis`: Cliente y operaciones
  - `redis/internal`: Parser de URL

**Estándar:**
- `context`: Context para operaciones y timeouts
- `encoding/json`: Serialización/deserialización
- `time`: Durations
- `fmt`: Error formatting

## Notas de diseño

- **Interfaz genérica**: CacheService abstrae Redis completamente
- **JSON automático**: Transparente en Get/Set, no requiere conversión manual
- **TTL flexible**: Por operación, permite patrones dinámicos
- **SCAN seguro**: DeleteByPattern usa SCAN para evitar bloqueos
- **TLS soportado**: rediss:// para conexiones seguras
- **Sin pooling adicional**: Redis client maneja pool internamente
- **Contexto sensible**: Respeta cancellation y timeouts
- **Errores claros**: Wrapped con contexto en operaciones de conexión
- **Concurrencia**: Client y Database son thread-safe

## Testing

**Unitarios:**
- Validación de RedisConfig, constantes
- Mocking de redis.Client para probar lógica

**Integración:**
- Requiere Redis running (local o Docker)
- Valida Connect realmente se conecta
- Valida Set/Get con serialización JSON
- Valida Delete y DeleteByPattern funcionan
