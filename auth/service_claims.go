package auth

import (
	stdErrors "errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// TokenUseService es el valor del claim `token_use` que identifica un
// service JWT (M2M / B2B), distinto del JWT de usuario. Sirve para que un
// token de servicio NO pueda colarse por el middleware de usuario ni
// viceversa (D14 del plan 020).
const TokenUseService = "service"

// ServiceClaims representa los claims de un service JWT (autenticación
// machine-to-machine entre APIs del backend). Es deliberadamente distinto de
// Claims (JWT de usuario): NO lleva `active_context` ni `user_id` porque el
// caller ya resolvió a quién afecta la operación; el filtro escuela/unidad no
// aplica a un servicio (D14).
//
// El `sub` (RegisteredClaims.Subject) se usa como `service:<client_id>` por
// convención; el `aud` (RegisteredClaims.Audience) identifica al servicio
// destino (ej. "edugo-api-platform").
type ServiceClaims struct {
	TokenUse string   `json:"token_use"`
	ClientID string   `json:"client_id"`
	Scopes   []string `json:"scopes"`
	jwt.RegisteredClaims
}

// HasScope indica si el service token incluye el scope dado.
func (c *ServiceClaims) HasScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// ServiceJWTManager genera y valida service JWTs (HS256). Es independiente de
// JWTManager (JWT de usuario): usa su PROPIO secret (`SERVICE_JWT_SECRET`,
// distinto del secret de usuarios) y valida `aud` además de `iss`, evitando
// que un secret comprometido de un plano afecte al otro (D14).
type ServiceJWTManager struct {
	issuer    string
	audience  string
	secretKey []byte
}

// NewServiceJWTManager crea un ServiceJWTManager.
//
// Parámetros:
//   - secretKey: secret HS256 para firmar/validar (SERVICE_JWT_SECRET, ≠ secret de usuarios).
//   - issuer: emisor esperado (ej. "edugo-identity").
//   - audience: servicio destino esperado (ej. "edugo-api-platform").
func NewServiceJWTManager(secretKey, issuer, audience string) *ServiceJWTManager {
	return &ServiceJWTManager{
		secretKey: []byte(secretKey),
		issuer:    issuer,
		audience:  audience,
	}
}

// GenerateServiceToken firma un service JWT para un cliente M2M.
//
// Lo usan los callers (ej. edugo-worker, edugo-api-learning) para obtener un
// token con el que invocar endpoints internos del gateway. El gateway solo
// valida; no firma.
//
// Parámetros:
//   - clientID: identificador del cliente M2M (debe existir en auth.service_clients).
//   - scopes: scopes concedidos al token (ej. ["notifications.dispatch"]).
//   - expiresIn: duración hasta la expiración (mínimo 1 minuto; recomendado ~15 min).
func (m *ServiceJWTManager) GenerateServiceToken(
	clientID string,
	scopes []string,
	expiresIn time.Duration,
) (string, time.Time, error) {
	if clientID == "" {
		return "", time.Time{}, errors.NewValidationError("clientID no puede estar vacío")
	}
	if expiresIn < time.Minute {
		return "", time.Time{}, errors.NewValidationError("expiresIn debe ser mayor a 1 minuto")
	}

	now := time.Now()
	expiresAt := now.Add(expiresIn)

	claims := ServiceClaims{
		TokenUse: TokenUseService,
		ClientID: clientID,
		Scopes:   scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    m.issuer,
			Subject:   "service:" + clientID,
			Audience:  jwt.ClaimStrings{m.audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", time.Time{}, errors.NewInternalError("no se pudo firmar el service JWT", err)
	}

	return signedToken, expiresAt, nil
}

// ParseServiceToken parsea y valida la firma, el método HS256, el issuer, la
// audience y la expiración del token, retornando los claims SIN verificar el
// `token_use`. Útil cuando el caller quiere inspeccionar antes de decidir.
func (m *ServiceJWTManager) ParseServiceToken(tokenString string) (*ServiceClaims, error) {
	parser := jwt.NewParser(
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(m.audience),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		// SH-1: exige claim `exp`. Un service token sin expiración (aunque
		// esté bien firmado) se rechaza: no aceptamos credenciales M2M eternas
		// aunque GenerateServiceToken siempre fije `exp`, la validación no debe
		// confiar en el emisor.
		jwt.WithExpirationRequired(),
		// SH-2: tolerancia de reloj entre el caller que firma (Cloud Run) y el
		// que valida (platform); evita rechazos por skew de pocos segundos en
		// `nbf`/`iat`/`exp`.
		jwt.WithLeeway(30*time.Second),
	)

	token, err := parser.ParseWithClaims(tokenString, &ServiceClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		if stdErrors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.NewUnauthorizedError("service token expired")
		}
		return nil, errors.NewUnauthorizedError("invalid service token")
	}

	claims, ok := token.Claims.(*ServiceClaims)
	if !ok || !token.Valid {
		return nil, errors.NewUnauthorizedError("invalid service token claims")
	}

	if claims.Issuer != m.issuer {
		return nil, errors.NewUnauthorizedError("invalid service token issuer")
	}

	return claims, nil
}

// ValidateServiceToken valida un service JWT completo: firma + iss + aud + exp
// (vía ParseServiceToken) y además exige `token_use == "service"`. Esto evita
// que un JWT de usuario (token_use vacío o "refresh") se acepte en rutas de
// servicio.
func (m *ServiceJWTManager) ValidateServiceToken(tokenString string) (*ServiceClaims, error) {
	claims, err := m.ParseServiceToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenUse != TokenUseService {
		return nil, errors.NewUnauthorizedError("invalid token type: expected service token")
	}

	return claims, nil
}
