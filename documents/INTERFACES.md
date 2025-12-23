# API de Interfaces Públicas

## Resumen de Interfaces

Este documento detalla todas las interfaces públicas exportadas por edugo-shared, incluyendo sus métodos, parámetros y valores de retorno.

---

## 1. Logger Interface

**Ubicación:** `logger/logger.go`  
**Import:** `github.com/EduGoGroup/edugo-shared/logger`

```go
// Logger define la interfaz para logging estructurado
type Logger interface {
    // Debug registra un mensaje de nivel debug
    // Parámetros: msg (mensaje), fields (pares key-value opcionales)
    Debug(msg string, fields ...interface{})

    // Info registra un mensaje de nivel info
    Info(msg string, fields ...interface{})

    // Warn registra un mensaje de nivel warning
    Warn(msg string, fields ...interface{})

    // Error registra un mensaje de nivel error
    Error(msg string, fields ...interface{})

    // Fatal registra un mensaje y termina la aplicación
    Fatal(msg string, fields ...interface{})

    // With crea un logger con campos contextuales
    // Retorna: nuevo Logger con los campos agregados
    With(fields ...interface{}) Logger

    // Sync sincroniza el buffer del logger
    // Retorna: error si falla el sync
    Sync() error
}
```

### Uso
```go
log.Info("user created", "user_id", "123", "email", "user@example.com")
log.Error("database error", "error", err, "query", "SELECT * FROM users")

userLog := log.With("user_id", "123")
userLog.Info("action completed") // Incluye user_id automáticamente
```

---

## 2. Factory Interfaces

**Ubicación:** `bootstrap/interfaces.go`  
**Import:** `github.com/EduGoGroup/edugo-shared/bootstrap`

### LoggerFactory

```go
// LoggerFactory crea instancias de logger
type LoggerFactory interface {
    // CreateLogger crea un logger configurado
    // Parámetros:
    //   - ctx: contexto de la operación
    //   - env: ambiente (local, dev, qa, prod)
    //   - version: versión de la aplicación
    // Retorna: logger configurado o error
    CreateLogger(ctx context.Context, env string, version string) (*logrus.Logger, error)
}
```

### PostgreSQLFactory

```go
// PostgreSQLFactory gestiona conexiones PostgreSQL
type PostgreSQLFactory interface {
    // CreateConnection crea una conexión GORM
    // Retorna: *gorm.DB o error
    CreateConnection(ctx context.Context, config PostgreSQLConfig) (*gorm.DB, error)
    
    // CreateRawConnection crea conexión SQL nativa
    // Retorna: *sql.DB o error
    CreateRawConnection(ctx context.Context, config PostgreSQLConfig) (*sql.DB, error)
    
    // Ping verifica conectividad
    // Retorna: error si no hay conectividad
    Ping(ctx context.Context, db *gorm.DB) error
    
    // Close cierra la conexión
    // Retorna: error si falla el cierre
    Close(db *gorm.DB) error
}
```

### MongoDBFactory

```go
// MongoDBFactory gestiona conexiones MongoDB
type MongoDBFactory interface {
    // CreateConnection crea conexión al servidor
    // Retorna: *mongo.Client o error
    CreateConnection(ctx context.Context, config MongoDBConfig) (*mongo.Client, error)
    
    // GetDatabase obtiene una base de datos específica
    // Retorna: *mongo.Database (nunca nil para cliente válido)
    GetDatabase(client *mongo.Client, dbName string) *mongo.Database
    
    // Ping verifica conectividad
    Ping(ctx context.Context, client *mongo.Client) error
    
    // Close cierra la conexión
    Close(ctx context.Context, client *mongo.Client) error
}
```

### RabbitMQFactory

```go
// RabbitMQFactory gestiona conexiones RabbitMQ
type RabbitMQFactory interface {
    // CreateConnection crea conexión AMQP
    // Retorna: *amqp.Connection o error
    CreateConnection(ctx context.Context, config RabbitMQConfig) (*amqp.Connection, error)
    
    // CreateChannel crea un canal de comunicación
    // Retorna: *amqp.Channel o error
    CreateChannel(conn *amqp.Connection) (*amqp.Channel, error)
    
    // DeclareQueue declara una cola
    // Retorna: información de la cola o error
    DeclareQueue(channel *amqp.Channel, queueName string) (amqp.Queue, error)
    
    // Close cierra canal y conexión
    Close(channel *amqp.Channel, conn *amqp.Connection) error
}
```

### S3Factory

```go
// S3Factory gestiona clientes S3
type S3Factory interface {
    // CreateClient crea cliente S3
    // Retorna: *s3.Client o error
    CreateClient(ctx context.Context, config S3Config) (*s3.Client, error)
    
    // CreatePresignClient crea cliente para URLs pre-firmadas
    // Retorna: cliente de presign
    CreatePresignClient(client *s3.Client) interface{}
    
    // ValidateBucket verifica acceso al bucket
    ValidateBucket(ctx context.Context, client *s3.Client, bucket string) error
}
```

---

## 3. Resource Interfaces

### MessagePublisher

**Ubicación:** `bootstrap/interfaces.go`

```go
// MessagePublisher publica mensajes a colas
type MessagePublisher interface {
    // Publish publica un mensaje
    // Parámetros:
    //   - ctx: contexto (para cancelación/timeout)
    //   - queueName: nombre de la cola destino
    //   - body: mensaje serializado en bytes
    // Retorna: error si falla la publicación
    Publish(ctx context.Context, queueName string, body []byte) error
    
    // PublishWithPriority publica con prioridad
    // Parámetros adicionales:
    //   - priority: 0-9 (mayor número = mayor prioridad)
    PublishWithPriority(ctx context.Context, queueName string, body []byte, priority uint8) error
    
    // Close cierra el publicador
    Close() error
}
```

### StorageClient

**Ubicación:** `bootstrap/interfaces.go`

```go
// StorageClient opera sobre almacenamiento de objetos
type StorageClient interface {
    // Upload sube un archivo
    // Parámetros:
    //   - key: path/nombre del archivo en el bucket
    //   - data: contenido del archivo
    //   - contentType: MIME type (e.g., "image/png")
    // Retorna: URL del archivo o error
    Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
    
    // Download descarga un archivo
    // Retorna: contenido del archivo o error
    Download(ctx context.Context, key string) ([]byte, error)
    
    // Delete elimina un archivo
    Delete(ctx context.Context, key string) error
    
    // GetPresignedURL genera URL temporal de acceso
    // Parámetros:
    //   - expirationMinutes: tiempo de validez en minutos
    // Retorna: URL pre-firmada o error
    GetPresignedURL(ctx context.Context, key string, expirationMinutes int) (string, error)
    
    // Exists verifica si un archivo existe
    // Retorna: true si existe, false si no, error si falla la verificación
    Exists(ctx context.Context, key string) (bool, error)
}
```

### DatabaseClient

```go
// DatabaseClient operaciones básicas de base de datos
type DatabaseClient interface {
    // Ping verifica conectividad
    Ping(ctx context.Context) error
    
    // Close cierra la conexión
    Close(ctx context.Context) error
    
    // GetStats obtiene estadísticas de conexión
    // Retorna: mapa con métricas (open_connections, idle, etc.)
    GetStats(ctx context.Context) (map[string]interface{}, error)
}
```

### HealthChecker

```go
// HealthChecker verifica salud de recursos
type HealthChecker interface {
    // Check verifica el estado
    // Retorna: error si el recurso no está saludable
    Check(ctx context.Context) error
    
    // GetResourceName retorna identificador del recurso
    GetResourceName() string
    
    // IsRequired indica si es obligatorio
    IsRequired() bool
}
```

---

## 4. Messaging Interfaces

**Ubicación:** `messaging/rabbit/`

### Publisher

```go
// Publisher interface para publicar mensajes
type Publisher interface {
    // Publish publica mensaje a un exchange
    // Parámetros:
    //   - exchange: nombre del exchange ("" para default)
    //   - routingKey: clave de enrutamiento
    //   - body: cualquier struct serializable a JSON
    Publish(ctx context.Context, exchange, routingKey string, body interface{}) error
    
    // PublishWithPriority publica con prioridad
    PublishWithPriority(ctx context.Context, exchange, routingKey string, body interface{}, priority uint8) error
    
    // Close libera recursos
    Close() error
}
```

### Consumer

```go
// MessageHandler procesa mensajes recibidos
// Retorna: error si el procesamiento falla (mensaje será NACK'd)
type MessageHandler func(ctx context.Context, body []byte) error

// Consumer interface para consumir mensajes
type Consumer interface {
    // Consume inicia consumo de una cola
    // El handler se ejecuta para cada mensaje
    // Bloquea hasta que ctx se cancele o haya error
    Consume(ctx context.Context, queueName string, handler MessageHandler) error
    
    // Close libera recursos
    Close() error
}
```

---

## 5. Auth Interfaces

**Ubicación:** `auth/jwt.go`

### JWTManager (Struct, no interface)

```go
// JWTManager gestiona tokens JWT
type JWTManager struct { /* campos privados */ }

// NewJWTManager crea una nueva instancia
// Parámetros:
//   - secretKey: clave secreta para firmar tokens (mínimo 32 chars recomendado)
//   - issuer: identificador del emisor del token
func NewJWTManager(secretKey, issuer string) *JWTManager

// GenerateToken genera un nuevo JWT
// Parámetros:
//   - userID: identificador único del usuario
//   - email: email del usuario
//   - role: rol del sistema (enum.SystemRole)
//   - expiresIn: duración de validez del token
// Retorna: token string o error
func (m *JWTManager) GenerateToken(userID, email string, role enum.SystemRole, expiresIn time.Duration) (string, error)

// ValidateToken valida y decodifica un token
// Retorna: claims extraídos o error (token inválido/expirado)
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error)

// RefreshToken genera nuevo token desde uno existente válido
func (m *JWTManager) RefreshToken(tokenString string, expiresIn time.Duration) (string, error)
```

### Claims Structure

```go
// Claims representa los datos del token JWT
type Claims struct {
    UserID string          `json:"user_id"` // ID único del usuario
    Email  string          `json:"email"`   // Email del usuario
    Role   enum.SystemRole `json:"role"`    // Rol en el sistema
    jwt.RegisteredClaims                    // Claims estándar (exp, iat, etc.)
}
```

---

## 6. Lifecycle Interface

**Ubicación:** `lifecycle/manager.go`

### Manager (Struct)

```go
// Manager gestiona ciclo de vida de recursos
type Manager struct { /* campos privados */ }

// NewManager crea un nuevo manager
func NewManager(log logger.Logger) *Manager

// Register registra un recurso con startup y cleanup
// Parámetros:
//   - name: identificador único del recurso
//   - startup: función ejecutada en Startup() (puede ser nil)
//   - cleanup: función ejecutada en Cleanup() (puede ser nil)
func (m *Manager) Register(name string, startup func(ctx context.Context) error, cleanup func() error)

// RegisterSimple registra solo función de cleanup
func (m *Manager) RegisterSimple(name string, cleanup func() error)

// Startup ejecuta startup de todos los recursos (orden FIFO)
// Retorna: error del primer recurso que falle
func (m *Manager) Startup(ctx context.Context) error

// Cleanup ejecuta cleanup de todos los recursos (orden LIFO)
// Continúa aunque algunos fallen, acumula errores
func (m *Manager) Cleanup() error

// Count retorna número de recursos registrados
func (m *Manager) Count() int

// Clear elimina todos los recursos sin ejecutar cleanup (para testing)
func (m *Manager) Clear()
```

---

## 7. Error Types

**Ubicación:** `common/errors/errors.go`

### AppError

```go
// AppError es el error personalizado de la aplicación
type AppError struct {
    Code       ErrorCode              // Código único (VALIDATION_ERROR, etc.)
    Message    string                 // Mensaje legible
    Details    string                 // Detalles adicionales
    StatusCode int                    // Código HTTP sugerido
    Fields     map[string]interface{} // Campos de contexto
    Internal   error                  // Error interno (no expuesto)
}

// Error implementa interface error
func (e *AppError) Error() string

// Unwrap permite usar errors.Is/errors.As
func (e *AppError) Unwrap() error

// Métodos fluent para construcción
func (e *AppError) WithDetails(details string) *AppError
func (e *AppError) WithField(key string, value interface{}) *AppError
func (e *AppError) WithInternal(err error) *AppError
```

### Constructores de Errores

```go
// Crear error genérico
func New(code ErrorCode, message string) *AppError

// Envolver error existente
func Wrap(err error, code ErrorCode, message string) *AppError

// Constructores específicos
func NewValidationError(message string) *AppError
func NewNotFoundError(resource string) *AppError
func NewAlreadyExistsError(resource string) *AppError
func NewUnauthorizedError(message string) *AppError
func NewForbiddenError(message string) *AppError
func NewInternalError(message string, err error) *AppError
func NewDatabaseError(operation string, err error) *AppError
func NewBusinessRuleError(message string) *AppError
func NewConflictError(message string) *AppError
func NewRateLimitError() *AppError

// Utilidades
func IsAppError(err error) bool
func GetAppError(err error) (*AppError, bool)
```

---

## 8. Config Structures

**Ubicación:** `config/base.go`

```go
// BaseConfig configuración base para servicios
type BaseConfig struct {
    Environment string          `mapstructure:"environment"`  // local|dev|qa|prod
    ServiceName string          `mapstructure:"service_name"`
    Server      ServerConfig    `mapstructure:"server"`
    Database    DatabaseConfig  `mapstructure:"database"`
    MongoDB     MongoDBConfig   `mapstructure:"mongodb"`
    Logger      LoggerConfig    `mapstructure:"logger"`
    Bootstrap   BootstrapConfig `mapstructure:"bootstrap"`
}

// ServerConfig configuración HTTP
type ServerConfig struct {
    Port         int           `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
    IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
    Host         string        `mapstructure:"host"`
}

// DatabaseConfig configuración PostgreSQL
type DatabaseConfig struct {
    Host            string        `mapstructure:"host"`
    Port            int           `mapstructure:"port"`
    User            string        `mapstructure:"user"`
    Password        string        `mapstructure:"password"`
    Database        string        `mapstructure:"database"`
    SSLMode         string        `mapstructure:"ssl_mode"`
    MaxOpenConns    int           `mapstructure:"max_open_conns"`
    MaxIdleConns    int           `mapstructure:"max_idle_conns"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
    ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

// Método helper
func (d *DatabaseConfig) ConnectionString() string
func (d *DatabaseConfig) ConnectionStringWithDB(database string) string
```

---

## 9. Enum Types

**Ubicación:** `common/types/enum/`

### SystemRole

```go
type SystemRole string

const (
    SystemRoleAdmin    SystemRole = "admin"
    SystemRoleTeacher  SystemRole = "teacher"
    SystemRoleStudent  SystemRole = "student"
    SystemRoleGuardian SystemRole = "guardian"
)

func (r SystemRole) IsValid() bool
func (r SystemRole) String() string
func AllSystemRoles() []SystemRole
func AllSystemRolesStrings() []string
```

### MaterialStatus

```go
type MaterialStatus string

const (
    MaterialStatusDraft     MaterialStatus = "draft"
    MaterialStatusPublished MaterialStatus = "published"
    MaterialStatusArchived  MaterialStatus = "archived"
)

func (s MaterialStatus) IsValid() bool
func (s MaterialStatus) String() string
func AllMaterialStatuses() []MaterialStatus
```

### ProgressStatus

```go
type ProgressStatus string

const (
    ProgressStatusNotStarted ProgressStatus = "not_started"
    ProgressStatusInProgress ProgressStatus = "in_progress"
    ProgressStatusCompleted  ProgressStatus = "completed"
)

func (p ProgressStatus) IsValid() bool
func (p ProgressStatus) String() string
func AllProgressStatuses() []ProgressStatus
```

### ProcessingStatus

```go
type ProcessingStatus string

const (
    ProcessingStatusPending    ProcessingStatus = "pending"
    ProcessingStatusProcessing ProcessingStatus = "processing"
    ProcessingStatusCompleted  ProcessingStatus = "completed"
    ProcessingStatusFailed     ProcessingStatus = "failed"
)

func (p ProcessingStatus) IsValid() bool
func (p ProcessingStatus) String() string
func AllProcessingStatuses() []ProcessingStatus
```
