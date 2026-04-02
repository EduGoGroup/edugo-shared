# Database PostgreSQL - Documentación Técnica

Documentación completa de componentes, métodos, transacciones y patrones de integración del módulo de conexión a PostgreSQL.

## Componentes

### Config

Estructura que modela la configuración para conectarse a PostgreSQL.

```go
type Config struct {
	Host               string        // Host del servidor PostgreSQL
	User               string        // Usuario para autenticación
	Password           string        // Contraseña para autenticación
	Database           string        // Nombre de la base de datos
	SSLMode            string        // Modo SSL (disable, require, verify-ca, verify-full)
	MaxLifetime        time.Duration // Tiempo máximo de vida de una conexión
	ConnectTimeout     time.Duration // Timeout para establecer conexión
	Port               int           // Puerto PostgreSQL (default 5432)
	MaxConnections     int           // Número máximo de conexiones en pool
	MaxIdleConnections int           // Número máximo de conexiones idle
	SearchPath         string        // search_path PostgreSQL (schemas)
}
```

**Campos:**
- `Host`: Hostname o IP del servidor PostgreSQL
- `User`: Usuario PostgreSQL para autenticación
- `Password`: Contraseña del usuario (segura en configuración)
- `Database`: Nombre de la base de datos a conectar
- `SSLMode`: Modo SSL (disable=dev local, require=producción, verify-ca/verify-full=certificados)
- `Port`: Puerto PostgreSQL (default 5432)
- `MaxConnections`: Máximo de conexiones en pool (default 25)
- `MaxIdleConnections`: Máximo de conexiones idle (default 5)
- `MaxLifetime`: Lifetime máximo de conexión antes de recrarse (default 5 minutos)
- `ConnectTimeout`: Timeout de conexión inicial (default 10 segundos)
- `SearchPath`: Schema path PostgreSQL (ej: "auth,iam,academic,public")

### DefaultConfig()

Retorna configuración con valores seguros para desarrollo local.

```go
func DefaultConfig() Config
```

**Valores por defecto:**
- `Host`: `localhost`
- `Port`: `5432`
- `User`: `postgres`
- `Database`: `postgres`
- `SSLMode`: `disable`
- `MaxConnections`: 25
- `MaxIdleConnections`: 5
- `MaxLifetime`: 5 minutos
- `ConnectTimeout`: 10 segundos
- `SearchPath`: `"auth,iam,academic,content,assessment,ui_config,public"`

### Connect()

Establece conexión nativa (sql.DB) a PostgreSQL con configuración de pool.

```go
func Connect(cfg *Config) (*sql.DB, error)
```

**Comportamiento:**
1. Construye DSN (Data Source Name) con credentials, SSL, timeout, search_path
2. Abre conexión con `sql.Open("postgres", dsn)`
3. Configura pool: MaxOpenConns, MaxIdleConns, ConnMaxLifetime
4. Valida conectividad con PingContext usando ConnectTimeout
5. Retorna error si Ping falla (credenciales, host no accesible, etc.)

**Errores:**
- `"failed to open database: %w"` - Error al abrir conexión
- `"failed to ping database: %w"` - Error al validar conectividad

**Notas:**
- DSN construido defensivamente con búsqueda de quoting de valores
- Pool se inicializa inmediatamente (eager)
- Validación de Ping asegura que PostgreSQL esté accesible

### ConnectGORM()

Establece conexión GORM a PostgreSQL con mismo pool configuration que sql.DB.

```go
func ConnectGORM(cfg *Config, gormLogger ...logger.Interface) (*gorm.DB, error)
```

**Parámetros:**
- `cfg`: Configuración PostgreSQL
- `gormLogger`: Logger GORM opcional (nil = sin logs)

**Comportamiento:**
1. Construye DSN equivalente a Connect()
2. Abre conexión GORM con `gorm.Open(postgres.Open(dsn), cfg)`
3. Obtiene underlying sql.DB con `db.DB()`
4. Aplica pool configuration (MaxOpenConns, MaxIdleConns, ConnMaxLifetime)
5. Retorna *gorm.DB conectado

**Errores:**
- `"failed to open GORM database: %w"` - Error al abrir conexión GORM
- `"failed to get underlying sql.DB: %w"` - Error al obtener sql.DB

**Notas:**
- Logger GORM opcional permite custom logging
- Pool configuration aplicada al underlying sql.DB
- Seguro compartir *gorm.DB entre goroutines

### HealthCheck()

Verifica que la conexión a PostgreSQL sea funcional.

```go
func HealthCheck(db *sql.DB) error
```

**Comportamiento:**
- Ejecuta PingContext con DefaultHealthCheckTimeout (5 segundos)
- Retorna nil si Ping exitoso, error si falla

**Retorno:**
- `nil`: Conexión healthy
- Error: Conexión no disponible (red, PostgreSQL caído, sin permisos)

**Casos de uso:**
- Validación periódica en background (liveness checks)
- Readiness probes en Kubernetes
- Circuitos de reintentos

### GetStats()

Obtiene estadísticas del pool de conexiones.

```go
func GetStats(db *sql.DB) sql.DBStats
```

**Retorno (sql.DBStats):**
```go
type DBStats struct {
	OpenConnections int // Conexiones abiertas actualmente
	InUse           int // Conexiones en uso
	Idle            int // Conexiones idle
	WaitCount       int64 // Total de operaciones esperando por conexión
	WaitDuration    time.Duration // Tiempo total esperando
	MaxIdleClosed   int64 // Total conexiones cerradas por MaxIdleConns
	MaxLifetimeClosed int64 // Total conexiones cerradas por MaxLifetime
}
```

**Casos de uso:**
- Monitoreo de pool health
- Detección de agotamiento de conexiones
- Alertas cuando WaitCount > 0 (conexiones esperando)

### Close()

Cierra la conexión a PostgreSQL.

```go
func Close(db *sql.DB) error
```

**Comportamiento:**
- Valida que db no sea nil
- Cierra todas las conexiones del pool
- Retorna error si cierre falla

**Notas:**
- Siempre llamar en defer después de Connect/ConnectGORM
- Es idempotente (seguro llamar multiple veces)

### TxFunc

Tipo de función para ejecutar dentro de transacción.

```go
type TxFunc func(*sql.Tx) error
```

### WithTransaction()

Ejecuta código dentro de transacción con rollback/commit automático.

```go
func WithTransaction(ctx context.Context, db *sql.DB, fn TxFunc) error
```

**Comportamiento:**
1. Comienza transacción con `db.BeginTx(ctx, nil)`
2. Ejecuta función fn(tx)
3. Si fn retorna error: rollback (con error joining si rollback falla)
4. Si fn exitoso: commit
5. Panic en fn: rollback automático + re-panic
6. Retorna error si begin, commit o rollback fallan

**Errores:**
- `"failed to begin transaction: %w"` - Error iniciar transacción
- `"failed to commit transaction: %w"` - Error commit
- `"tx error: %w" / "rollback error: %w"` - Error función + error rollback (errors.Join)

**Notas:**
- Aislamiento default: ReadCommitted
- Rollback automático en panic (defer protection)
- Context cancellation respetada (ctx debe ser válido)

### WithTransactionIsolation()

Transacción con nivel de aislamiento específico.

```go
func WithTransactionIsolation(ctx context.Context, db *sql.DB, isolation sql.IsolationLevel, fn TxFunc) error
```

**Parámetros:**
- `isolation`: sql.LevelDefault, sql.LevelReadUncommitted, sql.LevelReadCommitted, sql.LevelRepeatableRead, sql.LevelSerializable

**Comportamiento:**
- Idéntico a WithTransaction pero con nivel de aislamiento configurado
- Usa `db.BeginTx(ctx, &sql.TxOptions{Isolation: isolation})`

**Niveles de aislamiento (de menor a mayor restricción):**
- `LevelReadUncommitted`: Dirty reads posibles, máxima concurrencia
- `LevelReadCommitted`: Solo lecturas de datos committeados (default PostgreSQL)
- `LevelRepeatableRead`: Lecturas repetibles de datos
- `LevelSerializable`: Máxima aislación, transacciones serializadas

## Constantes

```go
const (
	DefaultPort                = 5432
	DefaultMaxConnections      = 25
	DefaultMaxIdleConnections  = 5
	DefaultMaxLifetime         = 5 * time.Minute
	DefaultConnectTimeout      = 10 * time.Second
	DefaultHealthCheckTimeout  = 5 * time.Second
)
```

## Flujos comunes

### Flujo 1: Inicialización en startup

```go
func initDatabase(ctx context.Context) (*sql.DB, error) {
	cfg := postgres.Config{
		Host:               os.Getenv("DB_HOST"),
		Port:               5432,
		User:               os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASSWORD"),
		Database:           os.Getenv("DB_NAME"),
		SSLMode:            "require",
		MaxConnections:     50,
		MaxIdleConnections: 10,
		ConnectTimeout:     15 * time.Second,
	}

	db, err := postgres.Connect(&cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres startup failed: %w", err)
	}

	log.Printf("PostgreSQL connected: %s@%s:%d", cfg.User, cfg.Host, cfg.Port)
	return db, nil
}

func main() {
	db, err := initDatabase(context.Background())
	if err != nil {
		log.Fatalf("startup failed: %v", err)
	}
	defer postgres.Close(db)

	// Usar db en aplicación...
}
```

### Flujo 2: Transacción con múltiples operaciones

```go
func createUserWithRole(ctx context.Context, db *sql.DB, username, email, role string) error {
	return postgres.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		// 1. Insertar usuario
		var userID int
		err := tx.QueryRowContext(ctx,
			"INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id",
			username, email,
		).Scan(&userID)
		if err != nil {
			return fmt.Errorf("insert user failed: %w", err)
		}

		// 2. Asignar rol
		_, err = tx.ExecContext(ctx,
			"INSERT INTO user_roles (user_id, role) VALUES ($1, $2)",
			userID, role,
		)
		if err != nil {
			return fmt.Errorf("assign role failed: %w", err)
		}

		return nil // Commit automático
	})
}

// Uso
err := createUserWithRole(context.Background(), db, "john", "john@example.com", "admin")
if err != nil {
	log.Printf("create user failed: %v", err)
}
```

### Flujo 3: Transacción con aislamiento serializable

```go
func transferBalance(ctx context.Context, db *sql.DB, fromID, toID int, amount int64) error {
	return postgres.WithTransactionIsolation(
		ctx, db, sql.LevelSerializable,
		func(tx *sql.Tx) error {
			// 1. Validar balance de origen
			var balance int64
			err := tx.QueryRowContext(ctx,
				"SELECT balance FROM accounts WHERE id = $1 FOR UPDATE",
				fromID,
			).Scan(&balance)
			if err != nil {
				return fmt.Errorf("query from_account failed: %w", err)
			}

			if balance < amount {
				return fmt.Errorf("insufficient balance: %d < %d", balance, amount)
			}

			// 2. Restar de origen
			_, err = tx.ExecContext(ctx,
				"UPDATE accounts SET balance = balance - $1 WHERE id = $2",
				amount, fromID,
			)
			if err != nil {
				return fmt.Errorf("debit from_account failed: %w", err)
			}

			// 3. Sumar a destino
			_, err = tx.ExecContext(ctx,
				"UPDATE accounts SET balance = balance + $1 WHERE id = $2",
				amount, toID,
			)
			if err != nil {
				return fmt.Errorf("credit to_account failed: %w", err)
			}

			return nil // Commit automático
		},
	)
}
```

### Flujo 4: Health checks y monitoreo de pool

```go
func monitorDatabaseHealth(db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// Health check
		if err := postgres.HealthCheck(db); err != nil {
			log.Printf("health check failed: %v", err)
			continue
		}

		// Pool statistics
		stats := postgres.GetStats(db)
		log.Printf(
			"pool: open=%d, in_use=%d, idle=%d, waiting=%d, wait_time=%v",
			stats.OpenConnections,
			stats.InUse,
			stats.Idle,
			stats.WaitCount,
			stats.WaitDuration,
		)

		// Alertar si hay espera por conexiones
		if stats.WaitCount > 0 {
			log.Printf("WARNING: %d queries waiting for connection", stats.WaitCount)
		}
	}
}

func main() {
	db, _ := postgres.Connect(postgres.DefaultConfig())
	defer postgres.Close(db)

	go monitorDatabaseHealth(db, 30*time.Second)

	// Usar db...
}
```

## Arquitectura

```
Config (host, user, password, pool, timeout, SSL)
    ↓
Connect() o ConnectGORM()
    ├─ Build DSN (search_path, SSL, timeout)
    ├─ sql.Open / gorm.Open
    ├─ Configure pool (MaxOpenConns, MaxIdleConns, MaxLifetime)
    ├─ Ping validation
    └─ Return *sql.DB or *gorm.DB
    ↓
Operations
    ├─ Query / QueryRow / Exec vía sql.DB o GORM
    ├─ WithTransaction() para operaciones multi-step
    └─ WithTransactionIsolation() para control de concurrencia
    ↓
HealthCheck() [periódicamente]
    └─ Ping
    ↓
GetStats()
    └─ Monitoreo del pool
    ↓
Close()
    └─ Cierre ordenado
```

**Ciclo de vida:**
1. **Init**: DefaultConfig() o Config{}
2. **Connect**: sql.Open + pool config + Ping
3. **Operations**: Query/Exec, WithTransaction
4. **HealthCheck**: Periódico (opcional)
5. **Close**: defer postgres.Close(db)

## Dependencias

**Internas:**
- Ninguna en producción (módulo autocontendido)

**Externas (sql.DB nativo):**
- `database/sql`: Interfaz estándar SQL
- `github.com/lib/pq`: Driver PostgreSQL para sql.DB
- `context`: Context para operaciones

**Externas (GORM):**
- `gorm.io/gorm`: ORM framework
- `gorm.io/driver/postgres`: Driver GORM para PostgreSQL
- Además todas las dependencias de sql.DB

**Estándar:**
- `context`: Cancellation y timeouts
- `time`: Durations
- `fmt`: Error formatting
- `errors`: Error handling

## Notas de diseño

- **Dual API**: sql.DB para máximo control, GORM para productividad
- **No elige por ti**: Usuario decide qué abstracción usar según necesidades
- **Pool configurable**: MaxConnections/MaxIdleConnections/MaxLifetime adaptables
- **Validación temprana**: Connect valida inmediatamente, no lazy
- **Transacciones seguras**: Rollback automático en error o panic
- **Aislamiento flexible**: WithTransactionIsolation para control de concurrencia
- **Monitoring**: GetStats permite observabilidad del pool
- **Search path**: Soporta múltiples schemas en SearchPath
- **SSL flexible**: SSLMode cubre dev (disable) a producción (verify-full)
- **Timeouts defensivos**: ConnectTimeout evita bloqueos indefinidos

## Testing

**Unitarios:**
- Validación de Config, DefaultConfig, constants
- Mocking de sql.DB para probar lógica

**Integración:**
- Requiere PostgreSQL running (local o Docker/Testcontainers)
- Tests de connection real, pool configuration
- Tests de transacciones con commit/rollback/panic
- Validación de niveles de aislamiento
