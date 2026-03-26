package gin

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// requestPath retorna la ruta registrada en el router (e.g. /api/v1/users/:id)
// para evitar alta cardinalidad en dashboards. Si FullPath no está disponible,
// cae al URL.Path literal como fallback.
func requestPath(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if fullPath := c.FullPath(); fullPath != "" {
		return fullPath
	}
	if c.Request != nil && c.Request.URL != nil {
		return c.Request.URL.Path
	}
	return ""
}

// requestMethod retorna el método HTTP de la petición o cadena vacía si no hay request.
func requestMethod(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return ""
	}
	return c.Request.Method
}

// getValidatedClaims es una función helper que extrae y valida los claims del contexto.
// Retorna los claims validados o nil si hay algún error, y envía la respuesta HTTP correspondiente.
func getValidatedClaims(c *gin.Context) *auth.Claims {
	claims, exists := c.Get(ContextKeyClaims)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
			"code":  "NO_CLAIMS",
		})
		return nil
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid claims type",
			"code":  "INVALID_CLAIMS_TYPE",
		})
		return nil
	}

	if userClaims.ActiveContext == nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "no active context",
			"code":  "NO_ACTIVE_CONTEXT",
		})
		return nil
	}

	return userClaims
}

// RequirePermission valida que el usuario tenga un permiso específico
func RequirePermission(permission enum.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := getValidatedClaims(c)
		if userClaims == nil {
			c.Abort()
			return
		}

		hasPermission := false
		for _, perm := range userClaims.ActiveContext.Permissions {
			if perm == permission.String() {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			reqLogger := GetLogger(c)
			reqLogger.Warn("permission denied",
				slog.String("required_permission", permission.String()),
				slog.String(logger.FieldPath, requestPath(c)),
				slog.String(logger.FieldMethod, requestMethod(c)),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "forbidden",
				"code":     "INSUFFICIENT_PERMISSIONS",
				"required": permission.String(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission valida que el usuario tenga AL MENOS uno de los permisos
func RequireAnyPermission(permissions ...enum.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := getValidatedClaims(c)
		if userClaims == nil {
			c.Abort()
			return
		}

		userPerms := make(map[string]bool)
		for _, perm := range userClaims.ActiveContext.Permissions {
			userPerms[perm] = true
		}

		for _, requiredPerm := range permissions {
			if userPerms[requiredPerm.String()] {
				c.Next()
				return
			}
		}

		permNames := make([]string, len(permissions))
		for i, p := range permissions {
			permNames[i] = p.String()
		}
		reqLogger := GetLogger(c)
		reqLogger.Warn("permission denied",
			slog.String("required_permissions", strings.Join(permNames, ",")),
			slog.String(logger.FieldPath, requestPath(c)),
			slog.String(logger.FieldMethod, requestMethod(c)),
		)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
			"code":  "INSUFFICIENT_PERMISSIONS",
		})
		c.Abort()
	}
}

// RequireAllPermissions valida que el usuario tenga TODOS los permisos
func RequireAllPermissions(permissions ...enum.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := getValidatedClaims(c)
		if userClaims == nil {
			c.Abort()
			return
		}

		userPerms := make(map[string]bool)
		for _, perm := range userClaims.ActiveContext.Permissions {
			userPerms[perm] = true
		}

		missing := []string{}
		for _, requiredPerm := range permissions {
			if !userPerms[requiredPerm.String()] {
				missing = append(missing, requiredPerm.String())
			}
		}

		if len(missing) > 0 {
			reqLogger := GetLogger(c)
			reqLogger.Warn("permission denied",
				slog.String("missing_permissions", strings.Join(missing, ",")),
				slog.String(logger.FieldPath, requestPath(c)),
				slog.String(logger.FieldMethod, requestMethod(c)),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "insufficient permissions",
				"code":    "INSUFFICIENT_PERMISSIONS",
				"missing": missing,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
