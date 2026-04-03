package gin

import (
	"net/http"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/gin-gonic/gin"
)

// ContextKeyActiveContext es la key para el UserContext en el contexto Gin.
// Complementa las keys existentes (ContextKeyUserID, ContextKeyEmail, etc.)
const ContextKeyActiveContext = "active_context"

// RemoteAuthMiddleware valida el token JWT del header Authorization
// usando el AuthClient (local + fallback remoto opcional), e inyecta
// la informacion del usuario en el contexto Gin.
func RemoteAuthMiddleware(authClient *AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
				"code":  "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format, expected 'Bearer {token}'",
				"code":  "INVALID_AUTH_FORMAT",
			})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]

		info, err := authClient.ValidateToken(c.Request.Context(), tokenString)
		if err != nil || !info.Valid {
			errMsg := "invalid or expired token"
			if info != nil && info.Error != "" {
				errMsg = info.Error
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMsg,
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, info.UserID)
		c.Set(ContextKeyEmail, info.Email)
		if info.ActiveContext != nil {
			c.Set(ContextKeyRole, info.ActiveContext.RoleName)
			c.Set(ContextKeyActiveContext, info.ActiveContext)
		}

		// Crear claims para que RequirePermission pueda leer los permisos
		c.Set(ContextKeyClaims, &auth.Claims{
			UserID:        info.UserID,
			Email:         info.Email,
			ActiveContext: info.ActiveContext,
		})

		c.Next()
	}
}
