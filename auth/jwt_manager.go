package auth

import (
	stdErrors "errors"
	"strings"
	"time"

	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// JWTManager firma y verifica el token de sesión de usuario de EduGo: el
// Context Token ES256 (plan 053). No hay otro: el HS256 de usuario se retiró
// en la T6.2 y un token HS256, o ES256 con un `kid` que el keyring no conoce,
// se rechaza siempre.
//
// Con un keyring solo de verificación (las APIs de dominio) el manager solo
// valida; con clave de firma (edugo-api-identity) además emite.
type JWTManager struct {
	issuer string
	// contextKeys is the ES256 keyring of the Context Token. Never nil.
	contextKeys *ContextKeyring
	// now is the clock; nil means time.Now. Only tests override it.
	now func() time.Time
}

// clock returns the current time from the manager's clock.
func (m *JWTManager) clock() time.Time {
	if m.now != nil {
		return m.now()
	}
	return time.Now()
}

// NewContextTokenManager crea el manager del Context Token ES256.
//
// keyring es obligatorio: sin claves no hay forma de verificar un token de
// usuario y el arranque debe abortar. Si el keyring no tiene clave de firma,
// el manager solo verifica (GenerateContextToken devuelve error).
func NewContextTokenManager(issuer string, keyring *ContextKeyring) (*JWTManager, error) {
	if strings.TrimSpace(issuer) == "" {
		return nil, stdErrors.New("context token: issuer requerido")
	}
	if keyring == nil {
		return nil, stdErrors.New("context token: keyring requerido")
	}
	return &JWTManager{issuer: issuer, contextKeys: keyring}, nil
}

// NewContextTokenManagerFromConfig es NewContextKeyringFromConfig más
// NewContextTokenManager. Una configuración vacía es un error.
func NewContextTokenManagerFromConfig(issuer string, cfg ContextKeyConfig) (*JWTManager, error) {
	keyring, err := NewContextKeyringFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	return NewContextTokenManager(issuer, keyring)
}

// Issuer returns the `iss` the manager signs and requires.
func (m *JWTManager) Issuer() string { return m.issuer }

// ValidateToken verifica un Context Token ES256 y retorna sus claims.
//
// Solo acepta `alg` ES256 con un `kid` del keyring; cualquier otra
// combinación (HS256, none, RS256, ES256 sin kid o con kid desconocido...) se
// rechaza. Exige `exp`, el issuer del manager y el claim active_context.
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	claims, err := m.validateDispatched(tokenString)
	if err != nil {
		return nil, err
	}

	// Validar que ActiveContext esté presente
	if claims.ActiveContext == nil {
		return nil, errors.NewUnauthorizedError("token missing active context")
	}

	return claims, nil
}
