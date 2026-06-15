package auth

import (
	stdErrors "errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-shared/common/errors"
)

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

// parseAndValidateToken es un helper interno que parsea y valida un JWT,
// retornando los claims sin verificaciones específicas de tipo de token.
func (m *JWTManager) parseAndValidateToken(tokenString string) (*Claims, error) {
	parser := jwt.NewParser(
		jwt.WithIssuer(m.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		// SH-1: exige claim `exp`. Un token sin expiración (aunque esté bien
		// firmado) se rechaza. Los tres generadores (GenerateTokenWithContext,
		// GenerateMinimalToken) siempre fijan `exp`, así que no rompe tokens
		// legítimos; solo cierra el hueco de aceptar un token eterno.
		jwt.WithExpirationRequired(),
		// SH-2: tolerancia de reloj para evitar rechazos por skew de pocos
		// segundos en `nbf`/`iat`/`exp` entre instancias.
		jwt.WithLeeway(30*time.Second),
	)

	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
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

	return claims, nil
}

// ValidateToken valida un JWT token y retorna los claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	claims, err := m.parseAndValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Validar que ActiveContext esté presente
	if claims.ActiveContext == nil {
		return nil, errors.NewUnauthorizedError("token missing active context")
	}

	return claims, nil
}

// GenerateMinimalToken genera un JWT minimal para refresh tokens.
//
// Guarda el snapshot del contexto activo (schoolID + academicUnitID + roleID) en
// claims raíz — no en ActiveContext — para que el refresh use case lo lea y emita
// el siguiente access_token con los mismos claims, sin reconstruir nada.
//
// La terna sujeto + modo de actor (subjectStudentID + actorMode, ADR 0026)
// también se persiste en el snapshot para que la rotación del access token
// preserve al acudido que el representante está viendo. Ambos pueden venir
// vacíos en contexto propio (actorMode = ActorModeSelf, que se omite por
// omitempty).
//
// Los demás snapshot pueden venir vacíos (ej. super_admin recién logueado que
// todavía no eligió escuela). En ese caso el siguiente access_token también irá
// sin esos claims y el cliente deberá completar la cascada con switch-context.
func (m *JWTManager) GenerateMinimalToken(
	userID, email, schoolID, academicUnitID, roleID, subjectStudentID, actorMode string,
	expiresIn time.Duration,
) (string, time.Time, error) {
	if userID == "" {
		return "", time.Time{}, errors.NewValidationError("userID no puede estar vacío")
	}
	if email == "" {
		return "", time.Time{}, errors.NewValidationError("email no puede estar vacío")
	}
	if expiresIn < time.Minute {
		return "", time.Time{}, errors.NewValidationError("expiresIn debe ser mayor a 1 minuto")
	}

	now := time.Now()
	expiresAt := now.Add(expiresIn)

	claims := Claims{
		UserID:           userID,
		Email:            email,
		ActiveContext:    nil,
		TokenUse:         "refresh",
		SchoolID:         schoolID,
		AcademicUnitID:   academicUnitID,
		RoleID:           roleID,
		SubjectStudentID: subjectStudentID,
		ActorMode:        actorMode,
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

// ValidateMinimalToken valida un JWT sin requerir ActiveContext.
// Diseñado para validar refresh tokens. Verifica que token_use sea "refresh".
func (m *JWTManager) ValidateMinimalToken(tokenString string) (*Claims, error) {
	claims, err := m.parseAndValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Verificar que el token sea de tipo refresh
	if claims.TokenUse != "refresh" {
		return nil, errors.NewUnauthorizedError("invalid token type: expected refresh token")
	}

	return claims, nil
}
