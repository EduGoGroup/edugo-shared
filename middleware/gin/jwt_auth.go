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

// JWTAuthMiddleware crea un middleware de autenticación JWT para Gin
// Valida el header Authorization, extrae y valida el token JWT
// y guarda los claims en el contexto de Gin para uso en handlers
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

		// 3. Validar token con JWTManager
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

// JWTAuthMiddlewareWithBlacklist creates a JWT auth middleware that also checks
// if the token has been revoked via the blacklist.
func JWTAuthMiddlewareWithBlacklist(jwtManager *auth.JWTManager, blacklist auth.TokenBlacklist) gin.HandlerFunc {
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

		// 3. Validar token con JWTManager
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

		// 4. Verificar si el token fue revocado
		if blacklist != nil && blacklist.IsRevoked(claims.ID) {
			GetLogger(c).Warn("revoked token used",
				logger.FieldPath, requestPath(c),
				logger.FieldMethod, requestMethod(c),
				logger.FieldIP, c.ClientIP(),
				"jti", claims.ID,
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token has been revoked",
				"code":  "TOKEN_REVOKED",
			})
			c.Abort()
			return
		}

		// 5. Guardar claims en contexto para uso en handlers
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyRole, claims.ActiveContext.RoleName)
		c.Set(ContextKeyClaims, claims)

		c.Next()
	}
}
