# Middleware Gin - Documentación Técnica

Documentación completa de componentes, métodos, flujos comunes e integración del middleware HTTP para Gin.

## Componentes

### JWTAuthMiddleware

Valida el token JWT del header `Authorization: Bearer ...` y puebla el contexto con claims.

```go
func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc
```

**Comportamiento:**
- Extrae token del header `Authorization: Bearer ...`
- Valida y parsea el JWT usando `auth.JWTManager`
- Puebla el contexto con claves estándar: `ContextKeyUserID`, `ContextKeyEmail`, `ContextKeyRole`, `ContextKeyClaims`
- Retorna HTTP 401 Unauthorized si el token es inválido o no existe

**Context Keys pobladas:**
- `user_id` (string): ID del usuario
- `email` (string): Email del usuario
- `role` (string): Rol principal del usuario
- `jwt_claims` (*auth.Claims): Estructura completa de claims JWT

### JWTAuthMiddlewareWithBlacklist

Extiende la validación JWT con verificación de tokens en lista negra.

```go
func JWTAuthMiddlewareWithBlacklist(jwtSecret string, checker auth.BlacklistChecker) gin.HandlerFunc
```

**Parámetros:**
- `jwtSecret`: Clave secreta para validación de firma JWT
- `checker`: Implementación de `auth.BlacklistChecker` para verificar tokens revocados

**Casos de uso:**
- Cierre de sesiones explícito
- Invalidación preventiva de tokens comprometidos
- Gestión de permisos modificados en tiempo real

### RequirePermission

Valida que el usuario tenga un permiso específico.

```go
func RequirePermission(permission enum.Permission) gin.HandlerFunc
```

**Validación:**
- Extrae claims del contexto (poblados por JWT middleware)
- Busca el permiso exacto en `ActiveContext.Permissions`
- Retorna HTTP 403 Forbidden si falta el permiso
- Registra advertencia en logs cuando se deniega acceso

**Ejemplo:**
```go
router.DELETE("/api/v1/users/:id",
	gin.RequirePermission(enum.PermissionUserDelete),
	handler,
)
```

### RequireAnyPermission

Valida que el usuario tenga AL MENOS uno de los permisos.

```go
func RequireAnyPermission(permissions ...enum.Permission) gin.HandlerFunc
```

**Validación:**
- Construye mapa de permisos del usuario para búsqueda O(1)
- Itera sobre permisos requeridos hasta encontrar coincidencia
- Retorna HTTP 403 Forbidden si ninguno coincide

**Ejemplo:**
```go
router.POST("/api/v1/assessments",
	gin.RequireAnyPermission(
		enum.PermissionAssessmentCreate,
		enum.PermissionAssessmentAdmin,
	),
	handler,
)
```

### RequireAllPermissions

Valida que el usuario tenga TODOS los permisos.

```go
func RequireAllPermissions(permissions ...enum.Permission) gin.HandlerFunc
```

**Validación:**
- Construye mapa de permisos del usuario
- Itera sobre permisos requeridos verificando cada uno
- Retorna HTTP 403 con lista de permisos faltantes

**Ejemplo:**
```go
router.PATCH("/api/v1/settings",
	gin.RequireAllPermissions(
		enum.PermissionSettingsEdit,
		enum.PermissionSettingsReview,
	),
	handler,
)
```

### RequestLogging

Genera request_id y correlation_id, crea logger enriquecido e inyecta en contexto.

```go
func RequestLogging(baseLogger *slog.Logger) gin.HandlerFunc
```

**Comportamiento:**
- Genera UUID como request_id si no viene en header `X-Request-ID`
- Propaga correlation_id del header `X-Correlation-ID` o usa request_id
- Establece headers de respuesta para trazabilidad
- Crea `slog.Logger` enriquecido con request_id, correlation_id, method, path, IP
- Inyecta logger en `gin.Context` y `context.Context`
- Registra resumen post-petición con status, duración (ms), bytes
- Log level: INFO para 200-399, WARN para 400-499, ERROR para 500+

### PostAuthLogging

Enriquece el logger con información de autenticación después del JWT middleware.

```go
func PostAuthLogging() gin.HandlerFunc
```

**Comportamiento:**
- Extrae user_id, role, school_id del contexto
- Re-enriquece logger con estos campos si están disponibles
- Re-inyecta logger en contexto para logs posteriores

**Cadena de middleware recomendada:**
```go
router.Use(gin.Recovery())
router.Use(gin.RequestLogging(baseLogger))
router.Use(gin.CORS())
router.Use(gin.JWTAuthMiddleware(secret))
router.Use(gin.PostAuthLogging())
router.Use(gin.AuditMiddleware(auditLogger))
```

### AuditMiddleware

Registra automáticamente todas las operaciones mutantes (POST, PUT, PATCH, DELETE).

```go
func AuditMiddleware(logger audit.AuditLogger) gin.HandlerFunc
```

**Comportamiento:**
- Ignora GET, HEAD, OPTIONS (solo lectura)
- Se ejecuta POST-handler para capturar status final
- Extrae user_id, email, role, school_id, unit_id del contexto
- Convierte método HTTP a acción: POST→create, PUT/PATCH→update, DELETE→delete
- Extrae resource_type e resource_id del URL path (ej: /api/v1/roles/123 → role, 123)
- Soporta patrones simples `/api/v1/{resource}[/{id}]`

**Evento de auditoría contiene:**
- `action`: create, update, delete
- `resource_type`: extraído de la ruta (singularizado)
- `resource_id`: ID del recurso si está en la ruta
- `actor_id`, `actor_email`, `actor_role`: Del contexto JWT
- `status_code`, `request_path`, `request_method`: De la petición
- `metadata`: school_id, unit_id del ActiveContext

### ParseListFilters

Parsea parámetros de paginación, búsqueda y filtrado desde query string.

```go
func ParseListFilters(c *gin.Context, extraFields ...string) (sharedrepo.ListFilters, error)
```

**Parámetros parseados:**
- `page`: Número de página (default: 1)
- `limit`: Items por página (default: 50, máx: 200)
- `search`: Término de búsqueda (string libre)
- `search_fields`: Campos donde buscar (comma-separated)
- `is_active`: Filtro booleano (nil = todos, true = activos, false = inactivos)
- Extra fields: Cualquier campo adicional como `school_id`, `status`

**Validación:**
- `limit` debe ser entero positivo, capped a 200
- `page` debe ser entero positivo
- `is_active` debe ser booleano válido (true, false)
- Retorna `*commonerrors.AppError` (HTTP 400) para valores inválidos

**Defensivo:**
- IsActive nil = "mostrar todos" (repositorios manejan via `ApplyIsActive()`)
- Límite default 50, máximo 200 para evitar queries costosas
- Extra fields accesibles via `filters.FieldFilters[fieldName]`

### Context Helpers

Extractores para acceder a claims y datos del usuario poblados por JWT middleware.

**Safe getters (retornan error):**
```go
GetUserID(c *gin.Context) (string, error)
GetEmail(c *gin.Context) (string, error)
GetRole(c *gin.Context) (string, error)
GetClaims(c *gin.Context) (*auth.Claims, error)
GetLogger(c *gin.Context) *slog.Logger
GetRequestID(c *gin.Context) string
```

**Unsafe getters (panic si no existe):**
```go
MustGetUserID(c *gin.Context) string
MustGetEmail(c *gin.Context) string
MustGetRole(c *gin.Context) string
MustGetClaims(c *gin.Context) *auth.Claims
```

**Error types:**
```go
var (
	ErrUserIDNotFound = errors.New("user_id not found in context")
	ErrEmailNotFound  = errors.New("email not found in context")
	ErrRoleNotFound   = errors.New("role not found in context")
	ErrClaimsNotFound = errors.New("claims not found in context")
	ErrInvalidType    = errors.New("invalid type in context")
)
```

## Flujos comunes

### Flujo 1: Autenticación y logging básico

```go
// Middleware setup
router.Use(gin.RequestLogging(slog.Default()))
router.Use(gin.JWTAuthMiddleware(secret))
router.Use(gin.PostAuthLogging())

// Handler
router.GET("/api/v1/profile", func(c *gin.Context) {
	logger := gin.GetLogger(c) // Obtiene logger enriquecido
	userID := gin.MustGetUserID(c)

	logger.Info("fetching profile", slog.String("user_id", userID))
	c.JSON(200, gin.H{"user_id": userID})
})
```

**Flujo:**
1. RequestLogging genera request_id, crea logger
2. JWTAuthMiddleware valida token, puebla contexto
3. PostAuthLogging enriquece logger con user_id
4. Handler obtiene logger pre-enriquecido y lo usa
5. RequestLogging registra resumen post-handler

### Flujo 2: Autorización con permisos múltiples

```go
router.Use(gin.JWTAuthMiddleware(secret))

// Requiere permiso específico
router.DELETE("/api/v1/users/:id",
	gin.RequirePermission(enum.PermissionUserDelete),
	deleteUserHandler,
)

// Requiere cualquiera de varios permisos
router.POST("/api/v1/materials",
	gin.RequireAnyPermission(
		enum.PermissionMaterialCreate,
		enum.PermissionMaterialAdmin,
	),
	createMaterialHandler,
)

// Requiere todos los permisos
router.PATCH("/api/v1/school-settings",
	gin.RequireAllPermissions(
		enum.PermissionSchoolSettingsEdit,
		enum.PermissionSchoolSettingsReview,
	),
	updateSettingsHandler,
)
```

**Flujo:**
1. JWT middleware valida token
2. Permission middleware valida contra `ActiveContext.Permissions`
3. Si falta permiso: HTTP 403, log warning, abort
4. Si OK: c.Next() continúa al handler

### Flujo 3: Auditoría automática de operaciones

```go
router.Use(gin.JWTAuthMiddleware(secret))
router.Use(gin.AuditMiddleware(auditLogger))

// POST /api/v1/assessments
router.POST("/api/v1/assessments", func(c *gin.Context) {
	// Handler procesa...
	c.JSON(201, gin.H{"id": "assess-123"})
	// Después de c.Next(), AuditMiddleware captura:
	// action: "create"
	// resource_type: "assessment"
	// status_code: 201
	// actor_id: user_id del JWT
})

// DELETE /api/v1/assessments/123
router.DELETE("/api/v1/assessments/:id", func(c *gin.Context) {
	c.JSON(204, nil)
	// AuditMiddleware captura:
	// action: "delete"
	// resource_type: "assessment"
	// resource_id: "123"
})
```

**Flujo:**
1. Handler procesa petición
2. c.Next() ejecuta
3. AuditMiddleware extrae resource_type, resource_id, actor_id
4. Registra en auditLogger

### Flujo 4: Paginación y filtrado de listas

```go
router.GET("/api/v1/users", func(c *gin.Context) {
	// Query: ?page=2&limit=25&search=john&is_active=true
	filters, err := gin.ParseListFilters(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// filters.Page = 2
	// filters.Limit = 25
	// filters.Search = "john"
	// filters.IsActive = &true (pointer)

	users, total := userRepo.List(c.Request.Context(), filters)
	c.JSON(200, gin.H{
		"data": users,
		"pagination": gin.H{
			"page": filters.Page,
			"limit": filters.Limit,
			"total": total,
		},
	})
})

// Con filtros extra
router.GET("/api/v1/assessments", func(c *gin.Context) {
	// Query: ?page=1&limit=50&school_id=sch-123&status=active,published
	filters, err := gin.ParseListFilters(c, "school_id", "status")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// filters.FieldFilters["school_id"] = ["sch-123"]
	// filters.FieldFilters["status"] = ["active", "published"]

	assessments, total := assessmentRepo.List(c.Request.Context(), filters)
	c.JSON(200, gin.H{"data": assessments, "total": total})
})
```

## Arquitectura

```
HTTP Request
    ↓
Recovery (panic recovery)
    ↓
RequestLogging (request_id, logger)
    ↓
CORS / Auth headers
    ↓
JWTAuthMiddleware (token validation, context population)
    ↓
PostAuthLogging (enrich logger with user_id, role)
    ↓
Permission Middleware (RequirePermission, RequireAnyPermission, RequireAllPermissions)
    ↓
Handler (business logic)
    ↓
AuditMiddleware (log mutating operations)
    ↓
HTTP Response
```

**Headers importantes:**
- `Authorization: Bearer <token>` - JWT token (extraído por JWT middleware)
- `X-Request-ID` - ID de petición (generado si no existe)
- `X-Correlation-ID` - ID de correlación para trazas distribuidas
- `User-Agent` - Capturado en auditoría
- `X-Forwarded-For` - IP del cliente (fallback en ClientIP())

## Dependencias

**Internas:**
- `github.com/EduGoGroup/edugo-shared/auth` - JWT parsing, Claims
- `github.com/EduGoGroup/edugo-shared/audit` - AuditLogger, AuditEvent
- `github.com/EduGoGroup/edugo-shared/common/types/enum` - Permission enums
- `github.com/EduGoGroup/edugo-shared/common/errors` - AppError, ValidationError
- `github.com/EduGoGroup/edugo-shared/repository` - ListFilters
- `github.com/EduGoGroup/edugo-shared/logger` - Log field constants

**Externas:**
- `github.com/gin-gonic/gin` - HTTP framework
- `log/slog` - Structured logging (standard library)
- `google/uuid` - UUID generation
- Standard library: `strings`, `errors`, `time`, `strconv`, `net/http`

## Operación

```bash
make build     # Compilar módulo
make test      # Tests unitarios
make test-race # Detector de race conditions
make check     # Validar: fmt, vet, lint, test
```

## Notas de diseño

- **Sin estado compartido**: Cada petición obtiene su propio logger, no hay mutación global
- **Context-first**: El contexto es el vehículo para pasar datos entre middlewares y handlers
- **Logging enriquecido**: RequestLogging + PostAuthLogging garantizan que todos los logs incluyan request_id y user_id
- **Auditoría post-handler**: AuditMiddleware captura status final después del handler
- **Validación defensiva**: ParseListFilters tiene defaults sensatos (limit 50, max 200)
- **Permiso nil = todos**: IsActive nil en ListFilters significa "mostrar todos", respeta decisión de repositorio
- **Singularización automática**: extractResourceFromPath() convierte plurales REST a singulares para auditoría
