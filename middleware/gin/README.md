# Middleware Gin

Middleware HTTP agnóstico para Gin con validación JWT, autorización por permisos, logging estructurado, auditoría automática y filtros de lista.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/middleware/gin
```

El módulo se descarga como `middleware/gin`, principal consumo vía package `gin`.

## Quick Start

### Ejemplo 1: Configurar cadena de middleware básica con JWT

```go
package main

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"github.com/EduGoGroup/edugo-shared/middleware/gin"
)

func main() {
	router := gin.Default()
	baseLogger := slog.Default()

	// Cadena de middleware recomendada (orden importante)
	router.Use(gin.Recovery())
	router.Use(gin.RequestLogging(baseLogger))
	router.Use(gin.CORS())
	router.Use(gin.JWTAuthMiddleware("your-secret-key"))
	router.Use(gin.PostAuthLogging())

	// Rutas protegidas
	router.GET("/api/v1/users", func(c *gin.Context) {
		userID := gin.MustGetUserID(c)
		email := gin.MustGetEmail(c)
		logger := gin.GetLogger(c)

		logger.Info("user fetched",
			slog.String("user_id", userID),
			slog.String("email", email),
		)
		c.JSON(200, gin.H{"user_id": userID})
	})

	router.Run(":8080")
}
```

### Ejemplo 2: Aplicar autorización por permisos

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/EduGoGroup/edugo-shared/middleware/gin"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

func main() {
	router := gin.Default()

	router.Use(gin.JWTAuthMiddleware("secret-key"))
	router.Use(gin.PostAuthLogging())

	// Requiere permiso específico
	router.DELETE("/api/v1/users/:id",
		gin.RequirePermission(enum.PermissionUserDelete),
		deleteUserHandler,
	)

	// Requiere AL MENOS uno de los permisos
	router.POST("/api/v1/assessments",
		gin.RequireAnyPermission(
			enum.PermissionAssessmentCreate,
			enum.PermissionAssessmentAdmin,
		),
		createAssessmentHandler,
	)

	// Requiere TODOS los permisos
	router.PATCH("/api/v1/settings",
		gin.RequireAllPermissions(
			enum.PermissionSettingsEdit,
			enum.PermissionSettingsReview,
		),
		updateSettingsHandler,
	)

	router.Run(":8080")
}
```

### Ejemplo 3: Configurar auditoría automática y logging

```go
package main

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"github.com/EduGoGroup/edugo-shared/middleware/gin"
	"github.com/EduGoGroup/edugo-shared/audit"
)

func main() {
	router := gin.Default()
	baseLogger := slog.Default()
	auditLogger := audit.NewLogger() // Tu implementación

	// Middleware en orden correcto
	router.Use(gin.Recovery())
	router.Use(gin.RequestLogging(baseLogger))
	router.Use(gin.CORS())
	router.Use(gin.JWTAuthMiddleware("secret-key"))
	router.Use(gin.PostAuthLogging())
	router.Use(gin.AuditMiddleware(auditLogger))

	router.POST("/api/v1/assessments", createAssessmentHandler)
	router.DELETE("/api/v1/assessments/:id", deleteAssessmentHandler)

	router.Run(":8080")
}
```

### Ejemplo 4: Parsear filtros de lista con paginación y búsqueda

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/EduGoGroup/edugo-shared/middleware/gin"
)

func main() {
	router := gin.Default()

	router.Use(gin.JWTAuthMiddleware("secret-key"))

	// Endpoint con filtrado: GET /api/v1/users?page=1&limit=25&search=john&is_active=true
	router.GET("/api/v1/users", listUsersHandler)

	// Endpoint con filtros extra: GET /api/v1/assessments?page=1&school_id=123&status=active
	router.GET("/api/v1/assessments", listAssessmentsHandler)

	router.Run(":8080")
}

func listUsersHandler(c *gin.Context) {
	// Parsear parámetros estándar: page, limit, search, search_fields, is_active
	filters, err := gin.ParseListFilters(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := gin.GetLogger(c)
	logger.Info("listing users",
		slog.Int("page", filters.Page),
		slog.Int("limit", filters.Limit),
	)

	users, total := yourUserRepo.List(c.Request.Context(), filters)
	c.JSON(200, gin.H{"data": users, "total": total})
}

func listAssessmentsHandler(c *gin.Context) {
	// Parsear parámetros estándar + campos extra
	filters, err := gin.ParseListFilters(c, "school_id", "status")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	assessments, total := yourAssessmentRepo.List(c.Request.Context(), filters)
	c.JSON(200, gin.H{"data": assessments, "total": total})
}
```

## Componentes principales

- **JWT Authentication**: Validación segura de tokens JWT con claims poblados en contexto
- **Permission Authorization**: Validación granular de permisos (RequirePermission, RequireAnyPermission, RequireAllPermissions)
- **Request Logging**: Enriquecimiento automático de logs con request_id, correlation_id, user_id, school_id
- **Audit Logging**: Registro automático de operaciones mutantes (POST, PUT, PATCH, DELETE)
- **Context Helpers**: Extractores seguros de user_id, email, role, claims del contexto
- **List Filters**: Parseo y validación de parámetros de paginación, búsqueda y filtrado

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

- **Orden de middleware**: Recovery → RequestLogging → CORS → JWT → PostAuthLogging → Audit → handlers
- **Logging enriquecido**: Cada operación incluye automáticamente request_id, user_id, school_id y duración
- **Auditoría automática**: Solo captura operaciones mutantes; GET/HEAD/OPTIONS son ignoradas por defecto
- **Permisos granulares**: Soporta validación individual, OR lógico y AND lógico de permisos
- **Paginación defensiva**: Límite máximo de 200 items, default 50; is_activo nil significa "todos"
- **Extractores seguros**: Métodos Get* retornan errores; MustGet* entran en pánico (solo post-JWT)
