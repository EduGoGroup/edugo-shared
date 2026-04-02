package gin

import (
	"strings"

	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/gin-gonic/gin"
)

// AuditMiddleware registra automáticamente todas las peticiones mutantes
// (POST, PUT, PATCH, DELETE) usando el AuditLogger proporcionado.
// Las peticiones de solo lectura (GET, HEAD, OPTIONS) son ignoradas.
func AuditMiddleware(logger audit.AuditLogger) gin.HandlerFunc { //nolint:gocyclo
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		c.Next()

		resourceType, resourceID := extractResourceFromPath(c.Request.URL.Path)

		event := audit.AuditEvent{
			Action:         methodToAction(c.Request.Method),
			ResourceType:   resourceType,
			ResourceID:     resourceID,
			RequestMethod:  c.Request.Method,
			RequestPath:    c.Request.URL.Path,
			RequestID:      c.GetHeader("X-Request-ID"),
			ActorIP:        c.ClientIP(),
			ActorUserAgent: c.GetHeader("User-Agent"),
			StatusCode:     c.Writer.Status(),
			Severity:       audit.SeverityInfo,
			Category:       audit.CategoryData,
		}

		// Extraer datos del JWT usando las context keys del middleware
		if userID, exists := c.Get("user_id"); exists {
			if v, ok := userID.(string); ok {
				event.ActorID = v
			}
		}
		if email, exists := c.Get("email"); exists {
			if v, ok := email.(string); ok {
				event.ActorEmail = v
			}
		}
		if role, exists := c.Get("role"); exists {
			if v, ok := role.(string); ok {
				event.ActorRole = v
			}
		}

		// Extraer school_id y unit_id de los claims JWT
		if claims, err := GetClaims(c); err == nil && claims.ActiveContext != nil {
			if claims.ActiveContext.SchoolID != "" {
				if event.Metadata == nil {
					event.Metadata = make(map[string]any)
				}
				event.Metadata["school_id"] = claims.ActiveContext.SchoolID
			}
			if claims.ActiveContext.AcademicUnitID != "" {
				if event.Metadata == nil {
					event.Metadata = make(map[string]any)
				}
				event.Metadata["unit_id"] = claims.ActiveContext.AcademicUnitID
			}
		}

		_ = logger.Log(c.Request.Context(), event) //nolint:errcheck
	}
}

func methodToAction(method string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return strings.ToLower(method)
	}
}

// extractResourceFromPath extrae el tipo de recurso e ID del path de la API.
// Asume el patrón /api/v1/{resource}[/{id}].
// Nota: solo soporta rutas simples de un nivel bajo /v1/.
// Ej: /api/v1/roles/123 → ("role", "123")
// Ej: /api/v1/memberships → ("membership", "")
func extractResourceFromPath(path string) (resourceType, resourceID string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	var resourceParts []string
	for i, part := range parts {
		if part == "v1" && i+1 < len(parts) {
			resourceParts = parts[i+1:]
			break
		}
	}

	if len(resourceParts) == 0 {
		return "unknown", ""
	}

	name := strings.ReplaceAll(resourceParts[0], "-", "_")
	resourceType = singularize(name)

	if len(resourceParts) > 1 {
		resourceID = resourceParts[1]
	}

	return resourceType, resourceID
}

// singularize convierte un nombre plural de recurso REST a su forma singular.
// Maneja los casos comunes de recursos en APIs: palabras en "ies", "ses",
// "xes", "ches", "shes" y el plural regular en "s".
func singularize(word string) string {
	switch {
	case strings.HasSuffix(word, "ies"):
		return word[:len(word)-3] + "y"
	case strings.HasSuffix(word, "ses"),
		strings.HasSuffix(word, "xes"),
		strings.HasSuffix(word, "ches"),
		strings.HasSuffix(word, "shes"):
		return word[:len(word)-2]
	case strings.HasSuffix(word, "s") && len(word) > 1:
		return word[:len(word)-1]
	default:
		return word
	}
}
