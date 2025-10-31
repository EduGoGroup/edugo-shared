// Package auth provides JWT token generation, validation, and management
// functionality for authentication and authorization in the EduGo shared library.
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// Claims representa los claims personalizados del JWT
type Claims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   enum.SystemRole `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager maneja la generación y validación de tokens JWT
type JWTManager struct {
	issuer    string
	secretKey []byte
}

// NewJWTManager crea un nuevo JWTManager
func NewJWTManager(secretKey, issuer string) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateToken genera un nuevo JWT token
func (m *JWTManager) GenerateToken(userID, email string, role enum.SystemRole, expiresIn time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(expiresIn)

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida un JWT token y retorna los claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, errors.NewUnauthorizedError("token expired")
		}
		return nil, errors.NewUnauthorizedError("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.NewUnauthorizedError("invalid token claims")
	}

	return claims, nil
}

// RefreshToken genera un nuevo token basado en uno existente (no expirado)
func (m *JWTManager) RefreshToken(tokenString string, expiresIn time.Duration) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return m.GenerateToken(claims.UserID, claims.Email, claims.Role, expiresIn)
}

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

// ExtractRole extrae el rol de un token sin validar completamente
// Útil solo para logging o debugging, NO para autenticación
func ExtractRole(tokenString string) (enum.SystemRole, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	return claims.Role, nil
}
