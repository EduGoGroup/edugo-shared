// Package auth provides JWT token generation, validation, and management
// functionality for authentication and authorization in the EduGo shared library.
package auth

import (
	stdErrors "errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// UserContext representa el contexto activo del usuario en el JWT
type UserContext struct {
	RoleID           string   `json:"role_id"`
	RoleName         string   `json:"role_name"`
	SchoolID         string   `json:"school_id,omitempty"`
	SchoolName       string   `json:"school_name,omitempty"`
	AcademicUnitID   string   `json:"academic_unit_id,omitempty"`
	AcademicUnitName string   `json:"academic_unit_name,omitempty"`
	Permissions      []string `json:"permissions"`
}

// Claims representa los claims personalizados del JWT
type Claims struct {
	UserID        string          `json:"user_id"`
	Email         string          `json:"email"`
	ActiveContext *UserContext     `json:"active_context,omitempty"`
	Role          enum.SystemRole `json:"role,omitempty"`     // Deprecated: usar ActiveContext
	SchoolID      string          `json:"school_id,omitempty"` // Deprecated: usar ActiveContext
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

// Deprecated: Usar GenerateTokenWithContext en su lugar
func (m *JWTManager) GenerateToken(userID, email string, role enum.SystemRole, expiresIn time.Duration) (string, error) {
	return m.GenerateTokenWithSchool(userID, email, role, "", expiresIn)
}

// Deprecated: Usar GenerateTokenWithContext en su lugar
func (m *JWTManager) GenerateTokenWithSchool(userID, email string, role enum.SystemRole, schoolID string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(expiresIn)

	claims := Claims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		SchoolID: schoolID,
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

// GenerateTokenWithContext genera un JWT con contexto RBAC
func (m *JWTManager) GenerateTokenWithContext(
	userID, email string,
	activeContext *UserContext,
	expiresIn time.Duration,
) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(expiresIn)

	claims := Claims{
		UserID:        userID,
		Email:         email,
		ActiveContext: activeContext,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secretKey)
	return signedToken, expiresAt, err
}

// ValidateToken valida un JWT token y retorna los claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	parser := jwt.NewParser(
		jwt.WithIssuer(m.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		if stdErrors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.NewUnauthorizedError("token expired")
		}
		return nil, errors.NewUnauthorizedError("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.NewUnauthorizedError("invalid token claims")
	}

	if claims.Issuer != m.issuer {
		return nil, errors.NewUnauthorizedError("invalid token issuer")
	}

	// Validar role solo si está presente (tokens legacy)
	// Tokens RBAC nuevos usan ActiveContext en su lugar
	if claims.Role != "" && !claims.Role.IsValid() {
		return nil, errors.NewUnauthorizedError("invalid role in token")
	}

	return claims, nil
}

// RefreshToken genera un nuevo token basado en uno existente (no expirado)
func (m *JWTManager) RefreshToken(tokenString string, expiresIn time.Duration) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return m.GenerateTokenWithSchool(claims.UserID, claims.Email, claims.Role, claims.SchoolID, expiresIn)
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

// ExtractSchoolID extrae el school ID de un token sin validar completamente
// Útil solo para logging o debugging, NO para autenticación
func ExtractSchoolID(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	return claims.SchoolID, nil
}
