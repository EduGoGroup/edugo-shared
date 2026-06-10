package gin

import (
	"net/http"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// Constantes para keys de contexto del plano M2M (service JWT).
// Deliberadamente separadas de las del usuario (ContextKeyUserID, etc.): un
// service token NO tiene user_id ni active_context.
const (
	ContextKeyClientID      = "client_id"
	ContextKeyScopes        = "scopes"
	ContextKeyServiceClaims = "service_jwt_claims"
)

// ServiceTokenValidator es la interfaz mínima que ServiceJWTAuthMiddleware
// necesita para validar un service JWT. La implementa `auth.ServiceJWTManager`.
// Se expresa como interfaz para facilitar el testing con dobles.
type ServiceTokenValidator interface {
	ValidateServiceToken(tokenString string) (*auth.ServiceClaims, error)
}

// ServiceJWTAuthMiddleware crea un middleware de autenticación M2M para rutas
// internas (`/api/v1/internal/*`). Es INDEPENDIENTE del middleware de usuario
// (JWTAuthMiddleware): no se mezclan en el mismo handler (D14 del plan 020).
//
// Rechaza con 401 (UNAUTHORIZED):
//   - header Authorization ausente o mal formado;
//   - token inválido/expirado, firma incorrecta, iss/aud distintos;
//   - token que no es de servicio (token_use != "service"), p. ej. un JWT de usuario.
//
// Rechaza con 403 (FORBIDDEN):
//   - token de servicio válido pero sin el `requiredScope`.
//
// En éxito inyecta en el contexto Gin: client_id, scopes y los claims completos.
func ServiceJWTAuthMiddleware(validator ServiceTokenValidator, requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Verificar header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			GetLogger(c).Warn("missing authorization header (service)",
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
			GetLogger(c).Warn("invalid authorization header format (service)",
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

		// 3. Validar service token (firma + iss + aud + exp + token_use=service)
		claims, err := validator.ValidateServiceToken(tokenString)
		if err != nil {
			GetLogger(c).Warn("service jwt validation failed",
				logger.FieldPath, requestPath(c),
				logger.FieldMethod, requestMethod(c),
				logger.FieldIP, c.ClientIP(),
				logger.FieldError, err.Error(),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired service token",
				"code":  "INVALID_SERVICE_TOKEN",
			})
			c.Abort()
			return
		}

		// 4. Verificar scope requerido (403 si falta)
		if requiredScope != "" && !claims.HasScope(requiredScope) {
			GetLogger(c).Warn("service token missing required scope",
				"client_id", claims.ClientID,
				"required_scope", requiredScope,
				logger.FieldPath, requestPath(c),
				logger.FieldMethod, requestMethod(c),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "insufficient scope",
				"code":  "INSUFFICIENT_SCOPE",
			})
			c.Abort()
			return
		}

		// 5. Inyectar identidad de servicio en el contexto
		c.Set(ContextKeyClientID, claims.ClientID)
		c.Set(ContextKeyScopes, claims.Scopes)
		c.Set(ContextKeyServiceClaims, claims)

		c.Next()
	}
}

// GetClientID extrae el client_id del service token desde el contexto Gin.
// Retorna error si la key no existe o el tipo no es string.
func GetClientID(c *gin.Context) (string, error) {
	clientID, exists := c.Get(ContextKeyClientID)
	if !exists {
		return "", ErrClientIDNotFound
	}

	clientIDStr, ok := clientID.(string)
	if !ok {
		return "", ErrInvalidType
	}

	return clientIDStr, nil
}

// GetServiceClaims extrae los claims del service token desde el contexto Gin.
// Retorna error si la key no existe o el tipo no es *auth.ServiceClaims.
func GetServiceClaims(c *gin.Context) (*auth.ServiceClaims, error) {
	claims, exists := c.Get(ContextKeyServiceClaims)
	if !exists {
		return nil, ErrServiceClaimsNotFound
	}

	claimsTyped, ok := claims.(*auth.ServiceClaims)
	if !ok {
		return nil, ErrInvalidType
	}

	return claimsTyped, nil
}
