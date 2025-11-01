package gin

import (
	"errors"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/gin-gonic/gin"
)

// Errores de contexto
var (
	ErrUserIDNotFound = errors.New("user_id not found in context")
	ErrEmailNotFound  = errors.New("email not found in context")
	ErrRoleNotFound   = errors.New("role not found in context")
	ErrClaimsNotFound = errors.New("claims not found in context")
	ErrInvalidType    = errors.New("invalid type in context")
)

// GetUserID extrae el user_id del contexto Gin
// Retorna error si la key no existe o el tipo no es string
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return "", ErrUserIDNotFound
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", ErrInvalidType
	}

	return userIDStr, nil
}

// MustGetUserID extrae el user_id o entra en pánico
// Útil cuando sabes que el middleware JWT ya validó el token
func MustGetUserID(c *gin.Context) string {
	userID, err := GetUserID(c)
	if err != nil {
		panic(err)
	}
	return userID
}

// GetEmail extrae el email del contexto Gin
// Retorna error si la key no existe o el tipo no es string
func GetEmail(c *gin.Context) (string, error) {
	email, exists := c.Get(ContextKeyEmail)
	if !exists {
		return "", ErrEmailNotFound
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", ErrInvalidType
	}

	return emailStr, nil
}

// MustGetEmail extrae el email o entra en pánico
// Útil cuando sabes que el middleware JWT ya validó el token
func MustGetEmail(c *gin.Context) string {
	email, err := GetEmail(c)
	if err != nil {
		panic(err)
	}
	return email
}

// GetRole extrae el role del contexto Gin
// Retorna error si la key no existe o el tipo no es string
func GetRole(c *gin.Context) (string, error) {
	role, exists := c.Get(ContextKeyRole)
	if !exists {
		return "", ErrRoleNotFound
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", ErrInvalidType
	}

	return roleStr, nil
}

// MustGetRole extrae el role o entra en pánico
// Útil cuando sabes que el middleware JWT ya validó el token
func MustGetRole(c *gin.Context) string {
	role, err := GetRole(c)
	if err != nil {
		panic(err)
	}
	return role
}

// GetClaims extrae todos los claims JWT del contexto
// Retorna error si la key no existe o el tipo no es *auth.Claims
func GetClaims(c *gin.Context) (*auth.Claims, error) {
	claims, exists := c.Get(ContextKeyClaims)
	if !exists {
		return nil, ErrClaimsNotFound
	}

	claimsTyped, ok := claims.(*auth.Claims)
	if !ok {
		return nil, ErrInvalidType
	}

	return claimsTyped, nil
}

// MustGetClaims extrae los claims o entra en pánico
// Útil cuando sabes que el middleware JWT ya validó el token
func MustGetClaims(c *gin.Context) *auth.Claims {
	claims, err := GetClaims(c)
	if err != nil {
		panic(err)
	}
	return claims
}
