package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// ExtractUserID extrae el user ID de un token sin validar completamente
// Útil solo para logging o debugging, NO para autenticación
func ExtractUserID(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	return claims.UserID, nil
}
