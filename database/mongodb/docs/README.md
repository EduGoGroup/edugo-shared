# Database MongoDB - Documentación Técnica

Documentación completa de componentes, métodos, ciclo de vida y patrones de integración del módulo de conexión a MongoDB.

## Componentes

### Config

Estructura que modela la configuración para conectarse a MongoDB.

```go
type Config struct {
	URI         string        // URI de conexión MongoDB
	Database    string        // Nombre de la database
	Timeout     time.Duration // Timeout para operaciones
	MaxPoolSize uint64        // Máximo de conexiones en pool
	MinPoolSize uint64        // Mínimo de conexiones en pool
}
```

**Campos:**
- `URI`: Connection string MongoDB en formato `mongodb://[user:pass@]host:port[/database][?options]`
  - Soporta replica sets, sharded clusters, MongoDB Atlas (mongodb+srv://...)
- `Database`: Nombre de la base de datos a usar por defecto
- `Timeout`: Timeout aplicado a Connect, Ping y operaciones de desconexión
- `MaxPoolSize`: Máximo de conexiones mantenidas en el pool (0 = ilimitado)
- `MinPoolSize`: Mínimo de conexiones siempre abiertas en el pool

### DefaultConfig()

Retorna configuración con valores seguros para desarrollo local.

```go
func DefaultConfig() Config
```

**Valores por defecto:**
- `URI`: `mongodb://localhost:27017`
- `Database`: `test`
- `Timeout`: 10 segundos
- `MaxPoolSize`: 100
- `MinPoolSize`: 10

### Connect()

Establece conexión a MongoDB, configura pool y valida conectividad.

```go
func Connect(cfg Config) (*mongo.Client, error)
```

**Comportamiento:**
1. Crea `options.Client()` con URI, pool sizes y timeouts
2. Conecta con timeout del contexto
3. Ejecuta Ping contra read preference Primary para validar
4. Desconecta y retorna error si Ping falla
5. Retorna cliente conectado

**Errores:**
- `"failed to connect to mongodb: %w"` - Error de conexión (timeout, URI inválido, credenciales)
- `"failed to ping mongodb: %w"` - Error de Ping (no hay primary accesible, sin permisos)

**Notas:**
- El cliente abre conexión lazy (primera operación real)
- Ping valida que al menos una primary esté accesible
- Pool se inicializa con MinPoolSize conexiones en background

### GetDatabase()

Obtiene una instancia de base de datos del cliente conectado.

```go
func GetDatabase(client *mongo.Client, databaseName string) *mongo.Database
```

**Parámetros:**
- `client`: Cliente MongoDB conectado
- `databaseName`: Nombre de la database a obtener

**Retorno:**
- `*mongo.Database`: Instancia de database para realizar operaciones

**Notas:**
- No valida que la database exista (validación lazy)
- Seguro para concurrencia
- El cliente mantiene caché interno de databases

### HealthCheck()

Verifica que la conexión a MongoDB sea funcional.

```go
func HealthCheck(client *mongo.Client) error
```

**Comportamiento:**
- Ejecuta Ping contra read preference Primary
- Usa timeout configurado (DefaultHealthCheckTimeout = 5 segundos)
- Retorna nil si Ping exitoso, error si falla

**Retorno:**
- `nil`: Conexión healthy
- Error: Conexión no disponible (puede ser red, instancia caída, sin permisos)

**Casos de uso:**
- Validación periódica en background (liveness checks)
- Pre-validación antes de operaciones críticas
- Circuitos de reintentos con backoff

### Close()

Cierra la conexión a MongoDB de forma controlada.

```go
func Close(client *mongo.Client) error
```

**Comportamiento:**
- Valida que client no sea nil
- Ejecuta Disconnect con timeout (DefaultDisconnectTimeout = 10 segundos)
- Drena conexiones del pool ordenadamente
- Aborts operaciones in-flight

**Retorno:**
- `nil`: Desconexión exitosa
- Error: Error durante Disconnect (timeout, error interno)

**Notas:**
- Siempre llamar en defer después de Connect
- Timeout evita bloqueos indefinidos
- Es seguro llamar multiple veces (close es idempotente en driver)

## Constantes

```go
const (
	DefaultTimeout            = 10 * time.Second
	DefaultMaxPoolSize        = 100
	DefaultMinPoolSize        = 10
	DefaultHealthCheckTimeout = 5 * time.Second
	DefaultDisconnectTimeout  = 10 * time.Second
)
```

**Propósito:**
- `DefaultTimeout`: Timeout para Connect, Ping, operaciones
- `DefaultMaxPoolSize`, `DefaultMinPoolSize`: Pool sizing defaults
- `DefaultHealthCheckTimeout`: Timeout específico para HealthCheck (más corto)
- `DefaultDisconnectTimeout`: Timeout para Close (más largo, permite drain)

## Flujos comunes

### Flujo 1: Inicialización en startup

```go
func initDatabase(ctx context.Context) (*mongo.Client, error) {
	// 1. Construir configuración
	cfg := mongodb.Config{
		URI:         os.Getenv("MONGODB_URI"),
		Database:    os.Getenv("MONGODB_DATABASE"),
		Timeout:     15 * time.Second,
		MaxPoolSize: 100,
		MinPoolSize: 20,
	}

	// 2. Conectar con validación
	client, err := mongodb.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("mongodb startup failed: %w", err)
	}

	// 3. Obtener database de trabajo
	db := mongodb.GetDatabase(client, cfg.Database)
	log.Printf("Connected to MongoDB database: %s", db.Name())

	return client, nil
}

func main() {
	client, err := initDatabase(context.Background())
	if err != nil {
		log.Fatalf("startup failed: %v", err)
	}
	defer mongodb.Close(client)

	// Usar client en aplicación...
}
```

### Flujo 2: Health checks en readiness probes

```go
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	if err := mongodb.HealthCheck(mongodbClient); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "mongodb unavailable: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func livelinessHandler(w http.ResponseWriter, r *http.Request) {
	if err := mongodb.HealthCheck(mongodbClient); err != nil {
		// Liveness checks pueden fallar pero no derogan el pod
		// Readiness checks si
		log.Printf("liveness: mongodb health check failed: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func main() {
	client, _ := mongodb.Connect(mongodb.DefaultConfig())
	defer mongodb.Close(client)

	http.HandleFunc("/healthz", livelinessHandler)
	http.HandleFunc("/readyz", readinessHandler)
	http.ListenAndServe(":8080", nil)
}
```

### Flujo 3: Configuración por entorno

```go
func configFromEnv() mongodb.Config {
	cfg := mongodb.DefaultConfig()

	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		cfg.URI = uri
	}
	if db := os.Getenv("MONGODB_DATABASE"); db != "" {
		cfg.Database = db
	}
	if maxPool := os.Getenv("MONGODB_MAX_POOL"); maxPool != "" {
		if v, err := strconv.ParseUint(maxPool, 10, 64); err == nil {
			cfg.MaxPoolSize = v
		}
	}
	if minPool := os.Getenv("MONGODB_MIN_POOL"); minPool != "" {
		if v, err := strconv.ParseUint(minPool, 10, 64); err == nil {
			cfg.MinPoolSize = v
		}
	}

	return cfg
}

// Uso
func main() {
	cfg := configFromEnv()
	client, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer mongodb.Close(client)

	// Usar client...
}
```

### Flujo 4: Reintentos con backoff exponencial

```go
func connectWithRetry(cfg mongodb.Config, maxRetries int) (*mongo.Client, error) {
	var client *mongo.Client
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		client, err = mongodb.Connect(cfg)
		if err == nil {
			return client, nil
		}

		backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		log.Printf("connect attempt %d failed, retrying in %v: %v", attempt+1, backoff, err)
		time.Sleep(backoff)
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}

func main() {
	cfg := mongodb.DefaultConfig()
	client, err := connectWithRetry(cfg, 5)
	if err != nil {
		log.Fatalf("startup failed: %v", err)
	}
	defer mongodb.Close(client)

	// Usar client...
}
```

## Arquitectura

```
Config (URI, Database, Pool, Timeout)
    ↓
Connect()
    ├─ Apply URI, pool, timeout to ClientOptions
    ├─ mongo.Connect(ctx, opts)
    ├─ Ping Primary (validation)
    └─ Return *mongo.Client
    ↓
GetDatabase(client, name)
    └─ Return *mongo.Database
    ↓
Operaciones (collections, inserts, queries, etc.)
    ↓
HealthCheck() [periódicamente]
    └─ Ping Primary
    ↓
Close()
    └─ client.Disconnect(ctx)
```

**Ciclo de vida:**
1. **Init**: DefaultConfig() o Config{}
2. **Connect**: mongo.Connect + Ping validation
3. **GetDatabase**: Obtener database de trabajo
4. **Operations**: Usar collections via database
5. **HealthCheck**: Periódico (opcional)
6. **Close**: defer mongodb.Close(client)

## Dependencias

**Internas:**
- Ninguna en producción (módulo autocontendo)

**Externas:**
- `go.mongodb.org/mongo-driver/v2`: Cliente oficial de MongoDB
  - `mongo`: Cliente, conexión, database
  - `mongo/options`: ClientOptions, read preferences
  - `mongo/readpref`: Read preference settings

**Estándar:**
- `context`: Context para timeouts
- `time`: Duration, timeouts
- `fmt`: Error formatting

## Notas de diseño

- **Bajo nivel**: Expone solo conexión y salud, no abstrae colecciones o queries
- **Defensivo**: Validación en Connect, timeouts en todas las operaciones
- **Pool configurable**: MaxPoolSize/MinPoolSize adaptables a carga
- **Concurrencia**: Driver es thread-safe, Client y Database son seguros para compartir
- **Lazy evaluation**: Database existe solo en operaciones reales, no en GetDatabase()
- **Timeouts distintos**: DefaultHealthCheckTimeout < DefaultDisconnectTimeout para diferenciar probes de cierre

## Testing

**Unitarios:**
- Validación de Config con DefaultConfig()
- Mocking de mongo.Client para probar lógica

**Integración:**
- Requiere MongoDB running (local o Docker/Testcontainers)
- Valida Connect realmente se conecta
- Valida HealthCheck retorna estado real
- Valida Close drena ordena
