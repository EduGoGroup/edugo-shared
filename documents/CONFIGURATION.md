# Guía de Configuración

## Visión General

La configuración en edugo-shared se maneja a través de archivos YAML y variables de entorno, con validación automática usando tags de struct.

---

## Estructura de Configuración

### Archivo config.yaml Completo

```yaml
# ============================================================================
# CONFIGURACIÓN BASE - EDUGO SHARED
# ============================================================================

# Ambiente de ejecución
# Valores válidos: local, dev, qa, prod
environment: local

# Nombre del servicio (para logging y métricas)
service_name: my-service

# ============================================================================
# SERVIDOR HTTP
# ============================================================================
server:
  # Puerto de escucha
  port: 8080
  
  # Host de binding (vacío = todas las interfaces)
  host: ""
  
  # Timeout de lectura de request
  read_timeout: 30s
  
  # Timeout de escritura de response
  write_timeout: 30s
  
  # Timeout de conexiones idle
  idle_timeout: 120s

# ============================================================================
# BASE DE DATOS POSTGRESQL
# ============================================================================
database:
  # Host del servidor PostgreSQL
  host: localhost
  
  # Puerto de conexión
  port: 5432
  
  # Usuario de base de datos
  user: edugo_user
  
  # Contraseña (preferiblemente desde variable de entorno)
  password: edugo_pass
  
  # Nombre de la base de datos
  database: edugo_db
  
  # Modo SSL: disable, require, verify-ca, verify-full
  ssl_mode: disable
  
  # Conexiones máximas abiertas simultáneamente
  max_open_conns: 25
  
  # Conexiones idle máximas en el pool
  max_idle_conns: 5
  
  # Tiempo máximo de vida de una conexión
  conn_max_lifetime: 5m
  
  # Tiempo máximo idle antes de cerrar
  conn_max_idle_time: 5m

# ============================================================================
# BASE DE DATOS MONGODB
# ============================================================================
mongodb:
  # URI de conexión completa
  uri: mongodb://localhost:27017
  
  # Nombre de la base de datos
  database: edugo_db
  
  # Tamaño máximo del pool de conexiones
  max_pool_size: 100
  
  # Tamaño mínimo del pool
  min_pool_size: 10
  
  # Timeout de conexión inicial
  connect_timeout: 30s

# ============================================================================
# LOGGER
# ============================================================================
logger:
  # Nivel de logging: debug, info, warn, error
  level: info
  
  # Formato de salida: json, console
  format: json

# ============================================================================
# BOOTSTRAP (Recursos Opcionales)
# ============================================================================
bootstrap:
  optional_resources:
    # Habilitar/deshabilitar RabbitMQ
    rabbitmq: true
    
    # Habilitar/deshabilitar S3
    s3: false

# ============================================================================
# RABBITMQ (si está habilitado)
# ============================================================================
rabbitmq:
  # URL AMQP completa
  url: amqp://edugo_user:edugo_pass@localhost:5672/

# ============================================================================
# AWS S3 (si está habilitado)
# ============================================================================
s3:
  # Nombre del bucket
  bucket: edugo-files
  
  # Región AWS
  region: us-east-1
  
  # Credenciales (preferiblemente desde variables de entorno)
  access_key_id: ""
  secret_access_key: ""
  
  # Endpoint personalizado (para MinIO/LocalStack)
  # endpoint: http://localhost:9000
```

---

## Variables de Entorno

Las variables de entorno pueden sobrescribir valores del archivo YAML.

### Convención de Nombres

```
EDUGO_{SECCION}_{CAMPO}
```

### Variables Disponibles

```bash
# ============================================================================
# AMBIENTE
# ============================================================================
EDUGO_ENVIRONMENT=local          # local, dev, qa, prod
EDUGO_SERVICE_NAME=my-service

# ============================================================================
# SERVIDOR
# ============================================================================
EDUGO_SERVER_PORT=8080
EDUGO_SERVER_HOST=
EDUGO_SERVER_READ_TIMEOUT=30s
EDUGO_SERVER_WRITE_TIMEOUT=30s
EDUGO_SERVER_IDLE_TIMEOUT=120s

# ============================================================================
# POSTGRESQL
# ============================================================================
EDUGO_DATABASE_HOST=localhost
EDUGO_DATABASE_PORT=5432
EDUGO_DATABASE_USER=edugo_user
EDUGO_DATABASE_PASSWORD=edugo_pass
EDUGO_DATABASE_NAME=edugo_db
EDUGO_DATABASE_SSL_MODE=disable
EDUGO_DATABASE_MAX_OPEN_CONNS=25
EDUGO_DATABASE_MAX_IDLE_CONNS=5
EDUGO_DATABASE_CONN_MAX_LIFETIME=5m

# ============================================================================
# MONGODB
# ============================================================================
EDUGO_MONGODB_URI=mongodb://localhost:27017
EDUGO_MONGODB_DATABASE=edugo_db
EDUGO_MONGODB_MAX_POOL_SIZE=100
EDUGO_MONGODB_MIN_POOL_SIZE=10
EDUGO_MONGODB_CONNECT_TIMEOUT=30s

# ============================================================================
# RABBITMQ
# ============================================================================
EDUGO_RABBITMQ_URL=amqp://edugo_user:edugo_pass@localhost:5672/

# ============================================================================
# S3
# ============================================================================
EDUGO_S3_BUCKET=edugo-files
EDUGO_S3_REGION=us-east-1
EDUGO_S3_ACCESS_KEY_ID=your-access-key
EDUGO_S3_SECRET_ACCESS_KEY=your-secret-key
EDUGO_S3_ENDPOINT=                    # Solo para MinIO/LocalStack

# ============================================================================
# LOGGER
# ============================================================================
EDUGO_LOGGER_LEVEL=info              # debug, info, warn, error
EDUGO_LOGGER_FORMAT=json             # json, console

# ============================================================================
# JWT (Para servicios que usan auth)
# ============================================================================
EDUGO_JWT_SECRET_KEY=your-super-secret-key-min-32-chars
EDUGO_JWT_ISSUER=edugo-auth
EDUGO_JWT_EXPIRATION=24h
```

---

## Validación de Configuración

La configuración se valida automáticamente usando tags `validate`:

```go
type BaseConfig struct {
    Environment string `validate:"required,oneof=local dev qa prod"`
    ServiceName string `validate:"required"`
    // ...
}

type ServerConfig struct {
    Port         int           `validate:"required,min=1,max=65535"`
    ReadTimeout  time.Duration `validate:"required"`
    WriteTimeout time.Duration `validate:"required"`
    IdleTimeout  time.Duration `validate:"required"`
}

type DatabaseConfig struct {
    Host     string `validate:"required"`
    Port     int    `validate:"required,min=1,max=65535"`
    User     string `validate:"required"`
    Password string `validate:"required"`
    Database string `validate:"required"`
    SSLMode  string `validate:"required,oneof=disable require verify-ca verify-full"`
}
```

### Errores de Validación Comunes

| Error | Causa | Solución |
|-------|-------|----------|
| `environment: oneof` | Valor no válido | Usar: local, dev, qa, prod |
| `port: min` | Puerto inválido | Puerto entre 1-65535 |
| `ssl_mode: oneof` | Modo SSL no válido | Usar: disable, require, etc. |
| `required` | Campo vacío | Proporcionar valor |

---

## Configuración por Ambiente

### Local (Desarrollo)

```yaml
environment: local

server:
  port: 8080

database:
  host: localhost
  port: 5432
  ssl_mode: disable

mongodb:
  uri: mongodb://localhost:27017

logger:
  level: debug
  format: console
```

### Development

```yaml
environment: dev

server:
  port: 8080

database:
  host: dev-postgres.internal
  port: 5432
  ssl_mode: require

mongodb:
  uri: mongodb://dev-mongo.internal:27017

logger:
  level: debug
  format: json
```

### QA

```yaml
environment: qa

server:
  port: 8080

database:
  host: qa-postgres.internal
  port: 5432
  ssl_mode: verify-ca

mongodb:
  uri: mongodb://qa-mongo.internal:27017

logger:
  level: info
  format: json
```

### Production

```yaml
environment: prod

server:
  port: 8080

database:
  host: prod-postgres.internal
  port: 5432
  ssl_mode: verify-full
  max_open_conns: 50
  max_idle_conns: 10

mongodb:
  uri: mongodb+srv://cluster.mongodb.net
  max_pool_size: 200
  min_pool_size: 50

logger:
  level: info
  format: json
```

---

## Carga de Configuración

### Código de Ejemplo

```go
package main

import (
    "log"
    
    "github.com/EduGoGroup/edugo-shared/config"
)

func main() {
    // Cargar desde archivo YAML
    cfg, err := config.Load("config.yaml")
    if err != nil {
        log.Fatalf("Error cargando configuración: %v", err)
    }
    
    // Validar configuración
    if err := config.Validate(cfg); err != nil {
        log.Fatalf("Configuración inválida: %v", err)
    }
    
    // Usar configuración
    log.Printf("Ambiente: %s", cfg.Environment)
    log.Printf("Puerto: %d", cfg.Server.Port)
}
```

### Prioridad de Configuración

```
1. Variables de entorno    (mayor prioridad)
2. Archivo YAML
3. Valores por defecto     (menor prioridad)
```

---

## Archivos .env

Para desarrollo local, puedes usar archivos `.env`:

```bash
# .env.local
EDUGO_ENVIRONMENT=local
EDUGO_DATABASE_HOST=localhost
EDUGO_DATABASE_PASSWORD=dev_password

# .env.dev
EDUGO_ENVIRONMENT=dev
EDUGO_DATABASE_HOST=dev-db.internal
EDUGO_DATABASE_PASSWORD=${DB_PASSWORD}  # Desde secretos
```

### Cargar con direnv

```bash
# .envrc
export_env .env.local
```

---

## Secretos

### Recomendaciones

1. **Nunca** commitear secretos en archivos de configuración
2. Usar variables de entorno para contraseñas
3. En producción, usar gestores de secretos (AWS Secrets Manager, Vault, etc.)

### Ejemplo con AWS Secrets Manager

```go
// Cargar secretos antes de inicializar
secretValue := awsSecretsManager.GetSecret("edugo/prod/database")
os.Setenv("EDUGO_DATABASE_PASSWORD", secretValue)

// Luego cargar configuración normalmente
cfg, err := config.Load("config.yaml")
```

---

## Configuración de Testing

```yaml
# config.test.yaml
environment: local

database:
  host: localhost
  port: 5433  # Puerto diferente para tests
  database: edugo_test

mongodb:
  uri: mongodb://localhost:27018
  database: edugo_test

logger:
  level: error  # Menos ruido en tests
  format: console
```

### Con Testcontainers

```go
func TestMain(m *testing.M) {
    // Los containers generan su propia configuración
    config := containers.NewConfig().
        WithPostgreSQL(nil).  // Usa defaults
        WithMongoDB(nil).
        Build()
    
    // El manager proporciona strings de conexión
    manager, _ := containers.GetManager(nil, config)
    
    // Acceder a la configuración generada
    connStr, _ := manager.PostgreSQL().ConnectionString(ctx)
}
```
