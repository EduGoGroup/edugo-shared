package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testServiceSecret   = "service-secret-key"
	testServiceIssuer   = "edugo-identity"
	testServiceAudience = "edugo-api-platform"
)

func newTestServiceManager() *ServiceJWTManager {
	return NewServiceJWTManager(testServiceSecret, testServiceIssuer, testServiceAudience)
}

func TestServiceJWTManager_GenerateAndValidate(t *testing.T) {
	m := newTestServiceManager()

	token, expiresAt, err := m.GenerateServiceToken("edugo-worker", []string{"notifications.dispatch"}, time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	assert.WithinDuration(t, time.Now().Add(time.Hour), expiresAt, time.Minute)

	claims, err := m.ValidateServiceToken(token)
	require.NoError(t, err)
	assert.Equal(t, TokenUseService, claims.TokenUse)
	assert.Equal(t, "edugo-worker", claims.ClientID)
	assert.Equal(t, []string{"notifications.dispatch"}, claims.Scopes)
	assert.Equal(t, "service:edugo-worker", claims.Subject)
	assert.Contains(t, claims.Audience, testServiceAudience)
	assert.True(t, claims.HasScope("notifications.dispatch"))
	assert.False(t, claims.HasScope("other.scope"))
}

func TestServiceJWTManager_GenerateValidations(t *testing.T) {
	m := newTestServiceManager()

	_, _, err := m.GenerateServiceToken("", []string{"notifications.dispatch"}, time.Hour)
	assert.Error(t, err)

	_, _, err = m.GenerateServiceToken("edugo-worker", nil, time.Second)
	assert.Error(t, err)
}

func TestServiceJWTManager_RejectsUserToken(t *testing.T) {
	// Un JWT de usuario (mismo secret, sin token_use="service") NO debe pasar
	// la validación de servicio.
	userManager := NewJWTManager(testServiceSecret, testServiceIssuer)
	userToken, _, err := userManager.GenerateTokenWithContext(
		"user-123", "user@edugo.com",
		&UserContext{RoleID: "r", RoleName: "Student"},
		time.Hour,
	)
	require.NoError(t, err)

	m := newTestServiceManager()
	_, err = m.ValidateServiceToken(userToken)
	assert.Error(t, err, "un JWT de usuario no debe validar como service token")
}

func TestServiceJWTManager_RejectsWrongAudience(t *testing.T) {
	signer := NewServiceJWTManager(testServiceSecret, testServiceIssuer, "otra-api")
	token, _, err := signer.GenerateServiceToken("edugo-worker", []string{"notifications.dispatch"}, time.Hour)
	require.NoError(t, err)

	validator := newTestServiceManager() // espera aud "edugo-api-platform"
	_, err = validator.ValidateServiceToken(token)
	assert.Error(t, err, "aud distinto debe rechazarse")
}

func TestServiceJWTManager_RejectsWrongIssuer(t *testing.T) {
	signer := NewServiceJWTManager(testServiceSecret, "otro-issuer", testServiceAudience)
	token, _, err := signer.GenerateServiceToken("edugo-worker", []string{"notifications.dispatch"}, time.Hour)
	require.NoError(t, err)

	validator := newTestServiceManager()
	_, err = validator.ValidateServiceToken(token)
	assert.Error(t, err, "iss distinto debe rechazarse")
}

func TestServiceJWTManager_RejectsWrongSecret(t *testing.T) {
	signer := NewServiceJWTManager("otro-secret", testServiceIssuer, testServiceAudience)
	token, _, err := signer.GenerateServiceToken("edugo-worker", []string{"notifications.dispatch"}, time.Hour)
	require.NoError(t, err)

	validator := newTestServiceManager()
	_, err = validator.ValidateServiceToken(token)
	assert.Error(t, err, "firma con otro secret debe rechazarse")
}

func TestServiceJWTManager_RejectsExpired(t *testing.T) {
	m := newTestServiceManager()

	// Firmar manualmente un token ya expirado (GenerateServiceToken impone mínimo 1m).
	now := time.Now()
	claims := ServiceClaims{
		TokenUse: TokenUseService,
		ClientID: "edugo-worker",
		Scopes:   []string{"notifications.dispatch"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    testServiceIssuer,
			Subject:   "service:edugo-worker",
			Audience:  jwt.ClaimStrings{testServiceAudience},
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testServiceSecret))
	require.NoError(t, err)

	_, err = m.ValidateServiceToken(signed)
	assert.Error(t, err, "token expirado debe rechazarse")
}

// signServiceClaims firma manualmente unos ServiceClaims con un método dado y
// el secret de prueba. Útil para construir tokens que GenerateServiceToken no
// permite (sin exp, alg distinto, etc.).
func signServiceClaims(t *testing.T, method jwt.SigningMethod, claims ServiceClaims) string {
	t.Helper()
	signed, err := jwt.NewWithClaims(method, claims).SignedString([]byte(testServiceSecret))
	require.NoError(t, err)
	return signed
}

// SH-1: un token bien firmado pero SIN claim `exp` debe rechazarse.
func TestServiceJWTManager_RejectsTokenWithoutExp(t *testing.T) {
	m := newTestServiceManager()

	now := time.Now()
	claims := ServiceClaims{
		TokenUse: TokenUseService,
		ClientID: "edugo-worker",
		Scopes:   []string{"notifications.dispatch"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       uuid.New().String(),
			Issuer:   testServiceIssuer,
			Subject:  "service:edugo-worker",
			Audience: jwt.ClaimStrings{testServiceAudience},
			IssuedAt: jwt.NewNumericDate(now),
			// SIN ExpiresAt a propósito.
		},
	}
	signed := signServiceClaims(t, jwt.SigningMethodHS256, claims)

	_, err := m.ValidateServiceToken(signed)
	assert.Error(t, err, "un service token sin claim exp debe rechazarse")
}

// SH-1/seguridad: alg=none (firma vacía) debe rechazarse.
func TestServiceJWTManager_RejectsAlgNone(t *testing.T) {
	m := newTestServiceManager()

	now := time.Now()
	claims := ServiceClaims{
		TokenUse: TokenUseService,
		ClientID: "edugo-worker",
		Scopes:   []string{"notifications.dispatch"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    testServiceIssuer,
			Subject:   "service:edugo-worker",
			Audience:  jwt.ClaimStrings{testServiceAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
	// alg=none requiere la clave centinela UnsafeAllowNoneSignatureType al firmar.
	signed, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).
		SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = m.ValidateServiceToken(signed)
	assert.Error(t, err, "alg=none debe rechazarse")
}

// Seguridad: alg-confusion (HS512 con el MISMO secret) debe rechazarse porque
// WithValidMethods solo admite HS256.
func TestServiceJWTManager_RejectsAlgConfusion(t *testing.T) {
	m := newTestServiceManager()

	now := time.Now()
	claims := ServiceClaims{
		TokenUse: TokenUseService,
		ClientID: "edugo-worker",
		Scopes:   []string{"notifications.dispatch"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    testServiceIssuer,
			Subject:   "service:edugo-worker",
			Audience:  jwt.ClaimStrings{testServiceAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
	signed := signServiceClaims(t, jwt.SigningMethodHS512, claims)

	_, err := m.ValidateServiceToken(signed)
	assert.Error(t, err, "alg-confusion HS512 debe rechazarse")
}

// SH-2: un token expirado hace pocos segundos (dentro del leeway de 30s) debe
// ACEPTARSE para tolerar clock-skew.
func TestServiceJWTManager_AcceptsWithinLeeway(t *testing.T) {
	m := newTestServiceManager()

	now := time.Now()
	claims := ServiceClaims{
		TokenUse: TokenUseService,
		ClientID: "edugo-worker",
		Scopes:   []string{"notifications.dispatch"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       uuid.New().String(),
			Issuer:   testServiceIssuer,
			Subject:  "service:edugo-worker",
			Audience: jwt.ClaimStrings{testServiceAudience},
			IssuedAt: jwt.NewNumericDate(now.Add(-time.Minute)),
			// Expirado hace 5s: dentro del leeway de 30s.
			ExpiresAt: jwt.NewNumericDate(now.Add(-5 * time.Second)),
		},
	}
	signed := signServiceClaims(t, jwt.SigningMethodHS256, claims)

	claimsOut, err := m.ValidateServiceToken(signed)
	require.NoError(t, err, "token expirado dentro del leeway debe aceptarse")
	assert.Equal(t, "edugo-worker", claimsOut.ClientID)
}

func TestServiceJWTManager_RejectsNonServiceTokenUse(t *testing.T) {
	m := newTestServiceManager()

	// Token correctamente firmado/aud/iss pero con token_use distinto.
	now := time.Now()
	claims := ServiceClaims{
		TokenUse: "refresh",
		ClientID: "edugo-worker",
		Scopes:   []string{"notifications.dispatch"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    testServiceIssuer,
			Subject:   "service:edugo-worker",
			Audience:  jwt.ClaimStrings{testServiceAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testServiceSecret))
	require.NoError(t, err)

	// ParseServiceToken no exige token_use; ValidateServiceToken sí.
	_, parseErr := m.ParseServiceToken(signed)
	assert.NoError(t, parseErr)

	_, valErr := m.ValidateServiceToken(signed)
	assert.Error(t, valErr, "token_use != service debe rechazarse en Validate")
}
