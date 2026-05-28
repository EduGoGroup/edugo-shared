package gin

import (
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// PermissionMetricsRecorder es la interfaz mínima que permite a este paquete
// emitir métricas de evaluación de permisos sin introducir una dependencia de
// build hacia edugo-shared/metrics.
type PermissionMetricsRecorder interface {
	RecordPermissionCheck(permission string, granted bool)
}

var permissionMetricsRecorder atomic.Value // PermissionMetricsRecorder

// SetPermissionMetricsRecorder registra el recorder que recibirá las
// observaciones de RequirePermission.
func SetPermissionMetricsRecorder(r PermissionMetricsRecorder) {
	if r == nil {
		permissionMetricsRecorder.Store((PermissionMetricsRecorder)(nil))
		return
	}
	permissionMetricsRecorder.Store(r)
}

func loadPermissionMetricsRecorder() PermissionMetricsRecorder {
	v := permissionMetricsRecorder.Load()
	if v == nil {
		return nil
	}
	r, _ := v.(PermissionMetricsRecorder)
	return r
}

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

func requestMethod(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return ""
	}
	return c.Request.Method
}

// getValidatedClaims extrae y valida los claims del contexto.
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

func recordPermissionCheck(permission string, granted bool) {
	rec := loadPermissionMetricsRecorder()
	if rec == nil {
		return
	}
	rec.RecordPermissionCheck(permission, granted)
}

// RequirePermission valida que los grants del usuario cubran el
// permission dado. Aplica deny precedence + glob matching path-based
// (auth.EvaluateGrants).
func RequirePermission(permission enum.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := getValidatedClaims(c)
		if userClaims == nil {
			c.Abort()
			return
		}

		granted := auth.EvaluateGrants(userClaims.ActiveContext.Grants, permission.String())
		recordPermissionCheck(permission.String(), granted)

		if !granted {
			reqLogger := GetLogger(c)
			reqLogger.Warn("permission denied",
				"required_permission", permission.String(),
				logger.FieldPath, requestPath(c),
				logger.FieldMethod, requestMethod(c),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
				"code":  "INSUFFICIENT_PERMISSIONS",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission valida que el usuario tenga AL MENOS uno de los permisos.
func RequireAnyPermission(permissions ...enum.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := getValidatedClaims(c)
		if userClaims == nil {
			c.Abort()
			return
		}

		grants := userClaims.ActiveContext.Grants
		for _, requiredPerm := range permissions {
			if auth.EvaluateGrants(grants, requiredPerm.String()) {
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
			"required_permissions", strings.Join(permNames, ","),
			logger.FieldPath, requestPath(c),
			logger.FieldMethod, requestMethod(c),
		)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
			"code":  "INSUFFICIENT_PERMISSIONS",
		})
		c.Abort()
	}
}

// RequireAllPermissions valida que el usuario tenga TODOS los permisos.
func RequireAllPermissions(permissions ...enum.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := getValidatedClaims(c)
		if userClaims == nil {
			c.Abort()
			return
		}

		grants := userClaims.ActiveContext.Grants
		missing := []string{}
		for _, requiredPerm := range permissions {
			if !auth.EvaluateGrants(grants, requiredPerm.String()) {
				missing = append(missing, requiredPerm.String())
			}
		}

		if len(missing) > 0 {
			reqLogger := GetLogger(c)
			reqLogger.Warn("permission denied",
				"missing_permissions", strings.Join(missing, ","),
				logger.FieldPath, requestPath(c),
				logger.FieldMethod, requestMethod(c),
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
