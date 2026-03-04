package gin

import (
	"strings"

	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/gin-gonic/gin"
)

// AuditMiddleware logs all mutating requests automatically
func AuditMiddleware(logger audit.AuditLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		c.Next()

		resourceType, resourceID := extractResourceFromPath(c.Request.URL.Path)

		opts := []audit.AuditOption{
			audit.WithCategory(audit.CategoryData),
		}

		// Extraer school_id y unit_id de los claims JWT
		if claims, err := GetClaims(c); err == nil && claims.ActiveContext != nil {
			if claims.ActiveContext.SchoolID != "" {
				opts = append(opts, audit.WithMetadata("school_id", claims.ActiveContext.SchoolID))
			}
			if claims.ActiveContext.AcademicUnitID != "" {
				opts = append(opts, audit.WithMetadata("unit_id", claims.ActiveContext.AcademicUnitID))
			}
		}

		_ = logger.LogFromGin(c,
			methodToAction(c.Request.Method),
			resourceType,
			resourceID,
			opts...,
		)
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

// extractResourceFromPath extrae el tipo de recurso e ID del path
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
	if strings.HasSuffix(name, "s") {
		name = name[:len(name)-1]
	}
	resourceType = name

	if len(resourceParts) > 1 {
		resourceID = resourceParts[1]
	}

	return resourceType, resourceID
}
