package gin

import (
	"net/http"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// Constantes para keys de contexto
// Usar constantes previene errores de typo en strings mágicos
const (
	ContextKeyUserID = "user_id"
	ContextKeyEmail  = "email"
	ContextKeyRole   = "role"
	ContextKeyClaims = "jwt_claims"
)

// JWTAuthMiddleware crea el middleware de autenticación de usuario para Gin.
// Valida el header Authorization, verifica el Context Token ES256 con el
// manager (auth.NewContextTokenManager) y guarda los claims en el contexto de
// Gin para uso en handlers. No hay blacklist: se retiró en el plan 053 (T6.2).
func JWTAuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Verificar header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			GetLogger(c).Warn("missing authorization header",
				logger.FieldPath, requestPath(c),
				logger.FieldMethod, requestMethod(c),
				logger.FieldIP, c.ClientIP(),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
				"code":  "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		// 2. Extraer token del header "Bearer {token}"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			GetLogger(c).Warn("invalid authorization header format",
				logger.FieldPath, requestPath(c),
				logger.FieldMethod, requestMethod(c),
				logger.FieldIP, c.ClientIP(),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format, expected 'Bearer {token}'",
				"code":  "INVALID_AUTH_FORMAT",
			})
			c.Abort()
			return
		}

		// 3. Verificar el Context Token ES256
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			GetLogger(c).Warn("jwt validation failed",
				logger.FieldPath, requestPath(c),
				logger.FieldMethod, requestMethod(c),
				logger.FieldIP, c.ClientIP(),
				logger.FieldError, err.Error(),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// 4. Guardar claims en contexto para uso en handlers
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyRole, claims.ActiveContext.RoleName)
		c.Set(ContextKeyClaims, claims)

		c.Next()
	}
}
