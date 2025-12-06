# DEPENDENCIES - SHARED

## Matriz de Dependencias

```
┌──────────────────────────────────────┐
│     SHARED Library (Foundation)      │
├──────────────────────────────────────┤
│ Dependencias Externas Go:            │
│ ├─ sirupsen/logrus (Logger)          │
│ ├─ golang-jwt/jwt/v5 (Auth)          │
│ ├─ streadway/amqp (RabbitMQ)         │
│ ├─ gorm.io/gorm (ORM)                │
│ ├─ go.mongodb.org/mongo-driver       │
│ ├─ spf13/viper (Config)              │
│ └─ Otras utilidades                  │
│                                      │
│ Dependencias del Sistema:            │
│ ├─ PostgreSQL 15+ (para tests)       │
│ ├─ MongoDB 7.0+ (para tests)         │
│ ├─ RabbitMQ 3.12+ (para tests)       │
│ └─ Go 1.21+ (compilación)            │
└──────────────────────────────────────┘
```

---

## Dependencias Go Externas

### Logger
```bash
go get github.com/sirupsen/logrus@latest
```

**Uso:**
```go
import "github.com/sirupsen/logrus"

log := logrus.New()
log.WithFields(logrus.Fields{
    "user_id": 42,
    "action": "login",
}).Info("User logged in")
```

### Authentication (JWT)
```bash
go get github.com/golang-jwt/jwt/v5@latest
```

**Uso:**
```go
import "github.com/golang-jwt/jwt/v5"

claims := jwt.MapClaims{
    "user_id": 42,
    "exp": time.Now().Add(24*time.Hour).Unix(),
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
```

### Messaging (RabbitMQ)
```bash
go get github.com/streadway/amqp@latest
```

**Uso:**
```go
import "github.com/streadway/amqp"

conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
ch, _ := conn.Channel()
```

### Database (GORM + Drivers)
```bash
go get gorm.io/gorm@latest
go get gorm.io/driver/postgres@latest
go get go.mongodb.org/mongo-driver@latest
```

**PostgreSQL:**
```go
import "gorm.io/gorm"
import "gorm.io/driver/postgres"

db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

**MongoDB:**
```go
import "go.mongodb.org/mongo-driver/mongo"

client, _ := mongo.Connect(context.Background(), opts)
```

### Configuration
```bash
go get github.com/spf13/viper@latest
```

**Uso:**
```go
import "github.com/spf13/viper"

viper.SetConfigName("config")
viper.ReadInConfig()
port := viper.GetInt("api.port")
```

---

## Dependencias del Sistema para Testing

### PostgreSQL 15+
```bash
# Para tests de database module
docker run -d \
  --name postgres \
  -e POSTGRES_USER=test_user \
  -e POSTGRES_PASSWORD=test_pass \
  -e POSTGRES_DB=edugo_test \
  postgres:15-alpine
```

### MongoDB 7.0+
```bash
# Para tests de database module
docker run -d \
  --name mongodb \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=admin \
  mongo:7.0
```

### RabbitMQ 3.12+
```bash
# Para tests de messaging module
docker run -d \
  --name rabbitmq \
  rabbitmq:3.12-management-alpine
```

---

## Versionamiento de Dependencias

```go
module github.com/EduGoGroup/edugo-shared

go 1.21

require (
    github.com/sirupsen/logrus v1.9.3
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/streadway/amqp v1.0.0
    gorm.io/gorm v1.25.0
    gorm.io/driver/postgres v1.5.2
    go.mongodb.org/mongo-driver v1.12.0
    github.com/spf13/viper v1.16.0
)
```

---

## Update Policy

```
Mensuales:
- go mod tidy
- go get -u ./...
- Testear cambios

Antes de Release:
- Verify todas las dependencias están pinned
- Validate compatibilidad con proyectos que usan SHARED
- Testing completo (>80% coverage)
```

---

## Checklist de Instalación para Desarrollo

```markdown
## Dependencias del Sistema
- [ ] Go 1.21+ instalado
- [ ] PostgreSQL 15+ (local o Docker)
- [ ] MongoDB 7.0+ (local o Docker)
- [ ] RabbitMQ 3.12+ (local o Docker)

## Dependencias Go
- [ ] go mod download ejecutado
- [ ] go mod tidy ejecutado
- [ ] go build ./... compila
- [ ] go test ./... pasa

## Servicios para Testing
- [ ] PostgreSQL accesible en localhost:5432
- [ ] MongoDB accesible en localhost:27017
- [ ] RabbitMQ accesible en localhost:5672

## Verificación
- [ ] make test pasa
- [ ] make coverage >= 80%
- [ ] make lint pasa
```

---

## Resolución de Problemas

### Dependencia desactualizada
```bash
# Revisar qué está desactualizado
go get -u ./...
go mod tidy

# Verificar compatible
go test ./...
```

### Incompatibilidad de versiones
```bash
# Pinear versión compatible
go get github.com/package@v1.2.3
go mod tidy
go test ./...
```

### Fallo de test por dependencia externa
```bash
# Verificar que servicio está corriendo
docker ps | grep postgres
docker ps | grep mongo
docker ps | grep rabbitmq

# Si no está, levantar
docker-compose up -d postgres mongo rabbitmq
```

---

## Dependencias Heredadas por Otros Proyectos

Cuando un proyecto importa SHARED:

```go
// En api-mobile/main.go
import "github.com/EduGoGroup/edugo-shared/logger"

// Automáticamente obtiene:
// - sirupsen/logrus (dependencia de logger)
// - Todas las otras deps de SHARED
```

**Implicación:**
- API Mobile nunca debe tener conflicto de versiones
- Si hay conflicto, SHARED debe ser el authority
- Todos proyectos deben usar misma versión de SHARED

---

## Actualización de SHARED

Proceso cuando nueva versión de dependencia externa:

```
1. SHARED dev actualiza dependencia
   go get -u github.com/package@latest

2. Testing completo en SHARED
   go test ./...
   go test -cover ./...

3. Tag nueva versión de SHARED
   git tag v1.3.1

4. Otros proyectos actualizar:
   go get -u github.com/EduGoGroup/edugo-shared@v1.3.1
   go test ./...

5. Commit y deploy
```

---

## Dependencias Mínimas Recomendadas

```
Para mínimo funcionamiento:

SHARED + Logger:
- Go 1.21
- logrus

SHARED + Database:
- Go 1.21
- PostgreSQL 15
- MongoDB 7.0
- gorm

SHARED + Auth:
- Go 1.21
- golang-jwt/jwt

SHARED + Messaging:
- Go 1.21
- RabbitMQ 3.12
- streadway/amqp
```
