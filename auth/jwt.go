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
)

// UserContext representa el contexto activo del usuario en el JWT.
// Encapsula el rol actual, la escuela y unidad académica asociadas, y los permisos
// específicos del usuario en ese contexto.
//
// Campos opcionales (omitempty):
//   - SchoolID, SchoolName: Solo para usuarios con contexto de escuela
//   - AcademicUnitID, AcademicUnitName: Solo para usuarios con contexto de unidad académica
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
	UserID        string       `json:"user_id"`
	Email         string       `json:"email"`
	ActiveContext *UserContext `json:"active_context"`
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

// GenerateTokenWithContext genera un JWT con contexto RBAC.
//
// Parámetros:
//   - userID: ID del usuario (requerido, no puede estar vacío)
//   - email: Email del usuario (requerido, no puede estar vacío)
//   - activeContext: Contexto activo del usuario con rol, escuela y permisos (requerido)
//   - expiresIn: Duración hasta la expiración del token (mínimo 1 minuto)
//
// Retorna:
//   - string: Token JWT firmado
//   - time.Time: Fecha y hora de expiración del token
//   - error: Error de validación si los parámetros son inválidos, o error interno si falla la firma
func (m *JWTManager) GenerateTokenWithContext(
	userID, email string,
	activeContext *UserContext,
	expiresIn time.Duration,
) (string, time.Time, error) {
	// Validar entradas
	if userID == "" {
		return "", time.Time{}, errors.NewValidationError("userID no puede estar vacío")
	}
	if email == "" {
		return "", time.Time{}, errors.NewValidationError("email no puede estar vacío")
	}
	if activeContext == nil {
		return "", time.Time{}, errors.NewValidationError("activeContext no puede ser nil")
	}
	if expiresIn < time.Minute {
		return "", time.Time{}, errors.NewValidationError("expiresIn debe ser mayor a 1 minuto")
	}

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
	if err != nil {
		return "", time.Time{}, errors.NewInternalError("no se pudo firmar el token JWT", err)
	}

	return signedToken, expiresAt, nil
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

	// Validar que ActiveContext esté presente
	if claims.ActiveContext == nil {
		return nil, errors.NewUnauthorizedError("token missing active context")
	}

	return claims, nil
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
