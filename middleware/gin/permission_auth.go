package gin

import (
	"net/http"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/gin-gonic/gin"
)

// RequirePermission valida que el usuario tenga un permiso específico
func RequirePermission(permission enum.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener claims del contexto (inyectados por JWTAuthMiddleware)
		claims, exists := c.Get(ContextKeyClaims)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
				"code":  "NO_CLAIMS",
			})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*auth.Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid claims type",
				"code":  "INVALID_CLAIMS_TYPE",
			})
			c.Abort()
			return
		}

		// Verificar si el contexto activo tiene el permiso
		if userClaims.ActiveContext == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "no active context",
				"code":  "NO_ACTIVE_CONTEXT",
			})
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
		claims, exists := c.Get(ContextKeyClaims)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*auth.Claims)
		if !ok || userClaims.ActiveContext == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "no active context"})
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
		claims, exists := c.Get(ContextKeyClaims)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*auth.Claims)
		if !ok || userClaims.ActiveContext == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "no active context"})
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
