# GitHub Copilot - Instrucciones Personalizadas: EduGo Shared

## 🌍 IDIOMA / LANGUAGE

**IMPORTANTE**: Todos los comentarios, sugerencias, code reviews y respuestas en chat deben estar **SIEMPRE EN ESPAÑOL**.

- ✅ Comentarios en Pull Requests: **español**
- ✅ Sugerencias de código: **español**
- ✅ Explicaciones en chat: **español**
- ✅ Mensajes de error: **español**

---

## 📚 Naturaleza del Proyecto

Este proyecto es una **librería compartida** (NO una aplicación), diseñada para ser consumida por otros proyectos del ecosistema EduGo:

- **edugo-api-mobile**: API REST móvil
- **edugo-api-administracion**: API REST administrativa
- **edugo-worker**: Worker background jobs
- **Futuros proyectos**: Cualquier servicio del ecosistema

### ⚠️ DIFERENCIA CRÍTICA: Librería vs Aplicación

```go
// ❌ INCORRECTO: Código específico de aplicación
func StartServer() {
    router := gin.New()
    router.Run(":8080")  // ❌ No debe tener punto de entrada
}

// ✅ CORRECTO: Código reutilizable
func NewJWTMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Lógica reutilizable
    }
}
```

---

## 🏗️ Arquitectura Multi-Módulo

Este proyecto implementa una **arquitectura modular independiente** con **7 módulos Go separados**:

```
edugo-shared/
├── common/                    # Módulo: Utilidades base
│   ├── config/               # Configuración compartida
│   ├── errors/               # Tipos de error
│   ├── types/                # Tipos comunes (enum, etc)
│   └── validator/            # Validación de datos
│
├── logger/                    # Módulo: Logging estructurado
│   ├── logger.go             # Interface Logger
│   └── zap_logger.go         # Implementación con Zap
│
├── auth/                      # Módulo: Autenticación
│   ├── jwt.go                # JWT Manager
│   ├── password.go           # Hash de passwords (bcrypt)
│   └── refresh_token.go      # Refresh tokens
│
├── middleware/
│   └── gin/                   # Módulo: Middleware para Gin
│       ├── context.go        # Helpers de contexto
│       └── jwt_auth.go       # Middleware JWT
│
├── messaging/
│   └── rabbit/                # Módulo: RabbitMQ helpers
│       └── ...               # Conexión, producers, consumers
│
└── database/
    ├── postgres/              # Módulo: PostgreSQL utilities
    └── mongodb/               # Módulo: MongoDB utilities
```

### Principios Modulares

1. **Independencia**: Cada módulo tiene su propio `go.mod`
2. **Versionado Independiente**: Se pueden versionar módulos por separado
3. **Dependencias Mínimas**: Solo importar lo necesario
4. **Zero External Config**: No depender de archivos de configuración externos

---

## 📦 Sistema de Versionado

### Versionado Semántico Estricto

Este proyecto usa **Semantic Versioning 2.0.0**:

```
vMAJOR.MINOR.PATCH

- MAJOR: Breaking changes (incompatibilidad con versiones anteriores)
- MINOR: Nuevas features (retrocompatibles)
- PATCH: Bug fixes (retrocompatibles)
```

### Versionado Híbrido: Global + Por Módulo

#### Tags Globales (para releases coordinadas)
```bash
v2.0.5     # Todos los módulos en sincronía
v2.1.0     # Nueva feature en múltiples módulos
```

#### Tags por Módulo (para cambios aislados)
```bash
middleware/gin/v0.0.1           # Nuevo módulo middleware/gin
auth/v2.1.0                     # Nueva feature solo en auth
logger/v2.0.6                   # Bugfix solo en logger
```

### ⚠️ REGLA CRÍTICA: Retrocompatibilidad

```go
// ❌ BREAKING CHANGE SIN MAJOR VERSION BUMP
// Antes: v2.0.5
func NewJWTManager(secret string) *JWTManager
// Después: v2.0.6 ← ERROR, debería ser v3.0.0
func NewJWTManager(secret string, expiration time.Duration) *JWTManager

// ✅ CORRECTO: Agregar función nueva (MINOR bump)
// v2.0.5 → v2.1.0
func NewJWTManager(secret string) *JWTManager           // Mantener
func NewJWTManagerWithExpiration(secret string, exp time.Duration) *JWTManager  // Nueva

// ✅ CORRECTO: Deprecar antes de eliminar
// v2.1.0: Agregar nueva función y marcar vieja como deprecated
// v3.0.0: Eliminar función deprecated (MAJOR bump)
```

---

## 🔄 Flujo de Desarrollo y Release

### Estrategia de Branches

```
feature/xyz → dev → PR a main (con tag manual) → Release automático
```

| Branch | Propósito | Workflows Activos |
|--------|-----------|-------------------|
| **feature/*** | Desarrollo de nuevas features | Ninguno (desarrollo local) |
| **dev** | Integración y testing | ci.yml, test.yml |
| **main** | Releases estables | release.yml, sync-dev-to-main.yml |

### Proceso de Release (MANUAL)

1. **Desarrollar** en feature branch
2. **PR a dev** → Ejecuta CI y tests
3. **Mergear** a dev cuando pase CI
4. **PR de dev a main** (cuando esté listo para release)
5. **Tú agregas tags manualmente** en el PR (ejemplo en descripción)
6. **Mergear PR** → GitHub Actions crea release automático
7. **Sync automático** de main → dev

#### Ejemplo de PR dev → main

```markdown
## Release v2.1.0

### Cambios
- Nueva feature: Middleware de rate limiting
- Bugfix: JWT expiration validation

### Tags a crear después del merge
- `v2.1.0` (tag global)
- `middleware/gin/v0.1.0` (tag del módulo actualizado)

### Breaking Changes
- Ninguno (retrocompatible)
```

---

## 🎯 Convenciones de Código

### Naming Conventions

```go
// Packages
package auth       // ✅ Lowercase, singular
package logger     // ✅ Lowercase, singular

// Tipos exportados (públicos)
type JWTManager struct { ... }     // ✅ PascalCase
type Logger interface { ... }      // ✅ PascalCase
type ErrorType int                 // ✅ PascalCase

// Tipos no exportados (privados)
type jwtClaims struct { ... }      // ✅ camelCase
type zapLogger struct { ... }      // ✅ camelCase

// Funciones exportadas
func NewJWTManager() *JWTManager   // ✅ PascalCase, prefijo "New" para constructores
func GenerateToken() string        // ✅ PascalCase

// Funciones privadas
func validateClaims() error        // ✅ camelCase
```

### Documentación (godoc)

```go
// ✅ CORRECTO: Documentar TODAS las funciones/tipos exportados
// JWTManager gestiona la generación y validación de tokens JWT.
// Utiliza el algoritmo HS256 y soporta claims personalizados.
type JWTManager struct {
    secretKey  string
    expiration time.Duration
}

// NewJWTManager crea una nueva instancia de JWTManager.
//
// Parámetros:
//   - secretKey: Clave secreta para firmar tokens (mínimo 32 caracteres)
//   - expiration: Duración de validez del token (recomendado: 15min)
//
// Retorna un error si la clave secreta es demasiado corta.
func NewJWTManager(secretKey string, expiration time.Duration) (*JWTManager, error) {
    if len(secretKey) < 32 {
        return nil, ErrInvalidSecretKey
    }
    return &JWTManager{
        secretKey:  secretKey,
        expiration: expiration,
    }, nil
}

// ❌ INCORRECTO: Sin documentación
func NewJWTManager(secretKey string, expiration time.Duration) (*JWTManager, error) {
    // ...
}
```

### Context Siempre Primero

```go
// ✅ CORRECTO: context.Context como primer parámetro
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field)
func (r *Repository) Save(ctx context.Context, entity *Entity) error

// ❌ INCORRECTO: Sin context o en posición incorrecta
func (l *Logger) Info(msg string, ctx context.Context)
func (r *Repository) Save(entity *Entity) error
```

### Manejo de Errores

```go
// ✅ CORRECTO: Usar tipos de error del módulo common/errors
import "github.com/EduGoGroup/edugo-shared/common/errors"

func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
    if tokenString == "" {
        return nil, errors.NewValidationError("token cannot be empty")
    }

    token, err := jwt.Parse(tokenString, m.keyFunc)
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, errors.NewUnauthorizedError("token expired")
        }
        return nil, errors.NewInternalError("failed to parse token", err)
    }

    return extractClaims(token), nil
}

// ❌ INCORRECTO: Usar fmt.Errorf o errors.New
return nil, fmt.Errorf("token expired")
return nil, errors.New("invalid token")
```

---

## 🧪 Testing Exhaustivo

### Cobertura Alta (>80%)

Este proyecto es una librería crítica, por lo tanto **REQUIERE alta cobertura de tests**:

```bash
# Meta de cobertura por módulo
common:           >85%
logger:           >80%
auth:             >90% (crítico para seguridad)
middleware/gin:   >85%
messaging/rabbit: >75%
database/*:       >70%
```

### Estructura de Tests

```go
// ✅ CORRECTO: Test exhaustivo con tabla de casos
func TestJWTManager_GenerateToken(t *testing.T) {
    tests := []struct {
        name      string
        userID    string
        email     string
        roles     []string
        wantErr   bool
        errType   error
    }{
        {
            name:    "valid token generation",
            userID:  "user-123",
            email:   "test@edugo.com",
            roles:   []string{"student"},
            wantErr: false,
        },
        {
            name:    "empty user id",
            userID:  "",
            email:   "test@edugo.com",
            roles:   []string{"student"},
            wantErr: true,
            errType: errors.ErrValidation,
        },
        {
            name:    "invalid email format",
            userID:  "user-123",
            email:   "invalid-email",
            roles:   []string{"student"},
            wantErr: true,
            errType: errors.ErrValidation,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            manager := NewJWTManager("test-secret-key-32-characters", 15*time.Minute)
            token, err := manager.GenerateToken(tt.userID, tt.email, tt.roles)

            if tt.wantErr {
                assert.Error(t, err)
                assert.True(t, errors.Is(err, tt.errType))
                assert.Empty(t, token)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, token)
            }
        })
    }
}
```

### Tests de Integración (cuando aplique)

```go
// Para módulos database/*, messaging/*
func TestPostgresRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Setup: Levantar testcontainer
    ctx := context.Background()
    container := setupPostgresContainer(t)
    defer container.Terminate(ctx)

    // Test real contra base de datos
    repo := NewPostgresRepository(container.ConnectionString())
    err := repo.Save(ctx, testEntity)
    assert.NoError(t, err)
}
```

### Mocks NO Recomendados

```go
// ⚠️ EVITAR mocks en librería compartida
// Razón: Las librerías deben ser testeadas con implementaciones reales

// ❌ NO hacer esto en edugo-shared
type MockJWTManager struct { ... }

// ✅ CORRECTO: Tests con implementaciones reales
manager := auth.NewJWTManager("secret", 15*time.Minute)
token, err := manager.GenerateToken("user-123", "test@test.com", []string{"admin"})
```

---

## 🔒 Seguridad

### Secrets y Configuración

```go
// ❌ INCORRECTO: Hardcodear secrets
const jwtSecret = "my-secret-key"

// ✅ CORRECTO: Secrets vienen del consumidor
func NewJWTManager(secretKey string, expiration time.Duration) *JWTManager {
    // La aplicación que consume la librería pasa el secret
}
```

### Validación de Entrada

```go
// ✅ CORRECTO: Validar TODAS las entradas públicas
func NewJWTManager(secretKey string, expiration time.Duration) (*JWTManager, error) {
    if len(secretKey) < 32 {
        return nil, errors.NewValidationError("secret key must be at least 32 characters")
    }
    if expiration < 1*time.Minute {
        return nil, errors.NewValidationError("expiration must be at least 1 minute")
    }
    // ...
}
```

### Vulnerabilidades Comunes a Evitar

```go
// ❌ SQL Injection (en database modules)
query := fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", userID)  // ❌ Vulnerable

// ✅ Usar prepared statements
query := "SELECT * FROM users WHERE id = $1"
row := db.QueryRowContext(ctx, query, userID)

// ❌ Command Injection
cmd := exec.Command("sh", "-c", userInput)  // ❌ Vulnerable

// ✅ Validar y sanitizar
if !isValidInput(userInput) {
    return errors.NewValidationError("invalid input")
}
```

---

## 📊 Dependencias Externas

### Política de Dependencias

```go
// ✅ PERMITIDO: Dependencias estables y mantenidas
go.uber.org/zap                    // Logging
github.com/golang-jwt/jwt/v5       // JWT
golang.org/x/crypto/bcrypt         // Hashing passwords
github.com/gin-gonic/gin           // Web framework (solo middleware/gin)

// ⚠️ EVALUAR: Dependencias no críticas
github.com/google/uuid             // UUIDs (evaluar alternatives)

// ❌ EVITAR: Dependencias poco mantenidas o experimentales
github.com/abandoned/library       // ❌ No mantenida
github.com/experimental/beta       // ❌ No estable
```

### Actualizaciones de Dependencias

```bash
# Revisar dependencias desactualizadas
go list -u -m all

# Actualizar con precaución (validar breaking changes)
go get -u github.com/golang-jwt/jwt/v5

# Siempre ejecutar tests después de actualizar
make test-all-modules
```

---

## 🚀 CI/CD y Workflows

### Workflows Activos

| Workflow | Trigger | Propósito |
|----------|---------|-----------|
| **ci.yml** | PR a dev/main | Tests en 7 módulos (matrix) |
| **test.yml** | PR + manual | Cobertura de código |
| **release.yml** | Push de tag `v*` | Crear GitHub Release |
| **sync-dev-to-main.yml** | Push a main | Sincronizar main → dev |

### Matrix Strategy

Los workflows usan matrix para paralelizar tests:

```yaml
strategy:
  fail-fast: false
  matrix:
    module:
      - common
      - logger
      - auth
      - middleware/gin
      - messaging/rabbit
      - database/postgres
      - database/mongodb
```

### Comandos Make para CI Local

```bash
# Ejecutar CI completo localmente
make ci-all-modules

# Tests con race detection
make test-race-all-modules

# Cobertura de todos los módulos
make coverage-all-modules

# Lint de todos los módulos
make lint-all-modules

# Validación completa pre-PR
make check-all-modules
```

---

## 📖 Documentación

### README.md de Módulos

Cada módulo debe tener un README.md explicando:

```markdown
# auth

Módulo de autenticación para proyectos EduGo.

## Instalación

\`\`\`bash
go get github.com/EduGoGroup/edugo-shared/auth
\`\`\`

## Uso

\`\`\`go
import "github.com/EduGoGroup/edugo-shared/auth"

jwtManager := auth.NewJWTManager("secret-key", 15*time.Minute)
token, err := jwtManager.GenerateToken("user-123", "user@test.com", []string{"admin"})
\`\`\`

## Features

- ✅ Generación de JWT tokens (HS256)
- ✅ Validación de tokens
- ✅ Refresh tokens
- ✅ Hash de passwords con bcrypt

## API Reference

Ver [godoc](https://pkg.go.dev/github.com/EduGoGroup/edugo-shared/auth)
```

### CHANGELOG.md

Mantener actualizado con cada release:

```markdown
## [2.1.0] - 2025-11-01

### Added
- Nuevo módulo `middleware/gin` con middleware JWT para Gin
- Función `NewJWTManagerWithExpiration` en módulo auth

### Changed
- Actualizar dependencia jwt a v5.3.0

### Deprecated
- Función `OldJWTManager` (usar `NewJWTManager` en su lugar)

### Fixed
- Bug en validación de refresh tokens expirados

### Security
- Actualizar golang.org/x/crypto para patch de seguridad
```

---

## 🎯 Casos de Uso de Consumidores

### Cómo Consumir edugo-shared

```go
// Proyecto: edugo-api-mobile
package main

import (
    "github.com/EduGoGroup/edugo-shared/auth"
    "github.com/EduGoGroup/edugo-shared/logger"
    "github.com/EduGoGroup/edugo-shared/middleware/gin"
)

func main() {
    // 1. Inicializar logger
    log := logger.NewZapLogger()
    defer log.Sync()

    // 2. Configurar JWT
    jwtManager := auth.NewJWTManager(
        os.Getenv("JWT_SECRET"),
        15 * time.Minute,
    )

    // 3. Configurar router con middleware
    router := gin.Default()

    // Rutas protegidas
    protected := router.Group("/api/v1")
    protected.Use(middleware.JWTAuthMiddleware(jwtManager))
    {
        protected.GET("/profile", getProfile)
    }

    router.Run(":8080")
}
```

### Actualizar a Nueva Versión

```bash
# Proyecto consumidor (api-mobile, api-administracion, worker)
cd edugo-api-mobile

# Actualizar a tag específico
go get github.com/EduGoGroup/edugo-shared@v2.1.0

# O actualizar módulo específico
go get github.com/EduGoGroup/edugo-shared/middleware/gin@v0.1.0

# Limpiar dependencias
go mod tidy

# Verificar actualización
go list -m github.com/EduGoGroup/edugo-shared
```

---

## ⚠️ Reglas de Oro

1. **Retrocompatibilidad es CRÍTICA**
   - NUNCA romper la API pública sin MAJOR version bump
   - Deprecar antes de eliminar

2. **Tests Exhaustivos**
   - Cobertura >80% (>90% en módulos críticos como auth)
   - Tests con tabla de casos

3. **Documentación Completa**
   - Godoc para todas las exportaciones públicas
   - README por módulo
   - CHANGELOG actualizado

4. **Zero Dependencies Innecesarias**
   - Solo agregar dependencias críticas y bien mantenidas
   - Evaluar alternativas antes de agregar

5. **Seguridad Primero**
   - Validar todas las entradas públicas
   - No hardcodear secrets
   - Seguir OWASP best practices

6. **Versionado Manual Coordinado**
   - Tags globales para releases mayores
   - Tags por módulo para cambios aislados
   - Documentar breaking changes

7. **Código como Documentación**
   - Nombres descriptivos
   - Funciones pequeñas y enfocadas
   - Comentarios solo cuando sea necesario (el código debe ser auto-explicativo)

---

## 📞 Soporte y Contribución

### Reportar Bugs

```markdown
**Descripción**: Breve descripción del bug

**Módulo afectado**: auth / logger / middleware/gin / etc

**Versión**: v2.0.5

**Pasos para reproducir**:
1. Importar módulo
2. Llamar función X con parámetro Y
3. Ver error Z

**Comportamiento esperado**: ...

**Comportamiento actual**: ...

**Logs/Stacktrace**: ...
```

### Sugerir Features

```markdown
**Feature**: Nombre de la feature

**Módulo**: auth / logger / nuevo módulo

**Descripción**: Descripción detallada

**Caso de uso**: Cómo mejoraría a los consumidores

**Breaking change**: Sí/No

**Propuesta de API**:
\`\`\`go
func NewFeature() *Feature { ... }
\`\`\`
```

---

## 📚 Referencias

- [Effective Go](https://go.dev/doc/effective_go)
- [Semantic Versioning](https://semver.org/)
- [Go Modules Reference](https://go.dev/ref/mod)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)

---

**Última actualización**: 2025-11-01
**Versión de Go**: 1.25.3
**Versión actual de la librería**: v2.0.5
