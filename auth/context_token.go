package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	stdErrors "errors"
	"fmt"
	"sort"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// Plan 053 · T2.1: the EduGo Context Token signed with ES256.
//
// The Context Token carries exactly the same claims as today's HS256 access
// token (Claims, including active_context); only the signature changes. The
// keyring holds an optional active P-256 private key (only the issuer of the
// token, edugo-api-identity, has it) plus every public key accepted for
// verification, indexed by `kid`, so a key can be rotated while tokens signed
// with the previous one live until their `exp`.

// headerKID is the JOSE header that names the verification key.
const headerKID = "kid"

// ContextSigningKey is the active ES256 signing key of the keyring.
type ContextSigningKey struct {
	// ID is the `kid` written in the header of every token signed with Private.
	ID string
	// Private is a P-256 private key.
	Private *ecdsa.PrivateKey
}

// ContextKeyring signs Context Tokens with its active key (if any) and verifies
// them with any of its accepted public keys, chosen by the `kid` header.
//
// A keyring without signing key is verify-only: that is what the domain APIs
// (academic, learning, platform, messaging) need.
type ContextKeyring struct {
	signing *ContextSigningKey
	public  map[string]*ecdsa.PublicKey
}

// NewContextKeyring builds a keyring.
//
// signing may be nil (verify-only keyring). When present, its public half is
// accepted automatically under its own kid. public holds the additional
// accepted public keys by kid (e.g. the previous key during a rotation). Every
// key must be P-256 and every kid must be non-empty and unique; the keyring
// must accept at least one key.
func NewContextKeyring(signing *ContextSigningKey, public map[string]*ecdsa.PublicKey) (*ContextKeyring, error) {
	accepted := make(map[string]*ecdsa.PublicKey, len(public)+1)
	for kid, pub := range public {
		if kid == "" {
			return nil, stdErrors.New("context keyring: every public key needs a kid")
		}
		if !isP256Public(pub) {
			return nil, fmt.Errorf("context keyring: public key %q is not a P-256 key", kid)
		}
		accepted[kid] = pub
	}

	var active *ContextSigningKey
	if signing != nil {
		if signing.ID == "" {
			return nil, stdErrors.New("context keyring: the signing key needs a kid")
		}
		if signing.Private == nil || !isP256Public(&signing.Private.PublicKey) {
			return nil, fmt.Errorf("context keyring: signing key %q is not a P-256 private key", signing.ID)
		}
		if _, dup := accepted[signing.ID]; dup {
			return nil, fmt.Errorf("context keyring: kid %q is both the signing key and a public key", signing.ID)
		}
		accepted[signing.ID] = &signing.Private.PublicKey
		active = &ContextSigningKey{ID: signing.ID, Private: signing.Private}
	}

	if len(accepted) == 0 {
		return nil, stdErrors.New("context keyring: no key configured")
	}

	return &ContextKeyring{signing: active, public: accepted}, nil
}

// CanSign reports whether the keyring has an active signing key.
func (k *ContextKeyring) CanSign() bool { return k != nil && k.signing != nil }

// ActiveKID returns the kid used to sign new tokens, or "" if verify-only.
func (k *ContextKeyring) ActiveKID() string {
	if !k.CanSign() {
		return ""
	}
	return k.signing.ID
}

// KIDs returns the accepted kids, sorted.
func (k *ContextKeyring) KIDs() []string {
	if k == nil {
		return nil
	}
	kids := make([]string, 0, len(k.public))
	for kid := range k.public {
		kids = append(kids, kid)
	}
	sort.Strings(kids)
	return kids
}

// PublicKey returns the accepted public key for kid.
func (k *ContextKeyring) PublicKey(kid string) (*ecdsa.PublicKey, bool) {
	if k == nil {
		return nil, false
	}
	pub, ok := k.public[kid]
	return pub, ok
}

func isP256Public(pub *ecdsa.PublicKey) bool {
	return pub != nil && pub.Curve == elliptic.P256() && pub.X != nil && pub.Y != nil
}

// WithContextKeyring attaches an ES256 keyring to the manager and returns it.
//
// Once attached, ValidateToken dispatches on the `alg` header: HS256 (without
// `kid`) keeps today's path, ES256 with a known `kid` goes to the keyring, and
// anything else is rejected. A nil keyring leaves the manager exactly as
// NewJWTManager built it.
func (m *JWTManager) WithContextKeyring(keyring *ContextKeyring) *JWTManager {
	m.contextKeys = keyring
	return m
}

// NewJWTManagerWithContextKeyring is NewJWTManager plus WithContextKeyring.
func NewJWTManagerWithContextKeyring(secretKey, issuer string, keyring *ContextKeyring) *JWTManager {
	return NewJWTManager(secretKey, issuer).WithContextKeyring(keyring)
}

// ContextKeyring returns the attached keyring, or nil.
func (m *JWTManager) ContextKeyring() *ContextKeyring { return m.contextKeys }

// GenerateContextToken signs an ES256 Context Token with the same parameters,
// validations and claims as GenerateTokenWithContext. It fails if the manager
// has no keyring or the keyring cannot sign.
func (m *JWTManager) GenerateContextToken(
	userID, email string,
	activeContext *UserContext,
	expiresIn time.Duration,
) (string, time.Time, error) {
	if err := validateAccessInputs(userID, email, activeContext); err != nil {
		return "", time.Time{}, err
	}
	if expiresIn < time.Minute {
		return "", time.Time{}, errors.NewValidationError("expiresIn debe ser mayor a 1 minuto")
	}
	return m.GenerateContextTokenUntil(userID, email, activeContext, m.clock().Add(expiresIn))
}

// GenerateContextTokenUntil signs an ES256 Context Token that expires at
// expiresAt. It exists for the exchange (plan 053 T3.1), where the Context
// Token must not outlive the Identity Token (`exp <=` identity's), which may
// be less than a minute away. expiresAt must be in the future.
func (m *JWTManager) GenerateContextTokenUntil(
	userID, email string,
	activeContext *UserContext,
	expiresAt time.Time,
) (string, time.Time, error) {
	if !m.contextKeys.CanSign() {
		return "", time.Time{}, errors.NewInternalError("no hay clave de firma ES256 configurada para el Context Token", nil)
	}
	now := m.clock()
	if !expiresAt.After(now) {
		return "", time.Time{}, errors.NewValidationError("expiresAt debe estar en el futuro")
	}
	claims, err := m.accessClaims(userID, email, activeContext, now, expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header[headerKID] = m.contextKeys.signing.ID
	signedToken, err := token.SignedString(m.contextKeys.signing.Private)
	if err != nil {
		return "", time.Time{}, errors.NewInternalError("no se pudo firmar el Context Token", err)
	}
	return signedToken, expiresAt, nil
}

// accessClaims builds the access-token claims. HS256 and ES256 share it so the
// shape of both tokens is identical by construction (REQ-C4).
func (m *JWTManager) accessClaims(
	userID, email string,
	activeContext *UserContext,
	now, expiresAt time.Time,
) (Claims, error) {
	if err := validateAccessInputs(userID, email, activeContext); err != nil {
		return Claims{}, err
	}
	return Claims{
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
	}, nil
}

func validateAccessInputs(userID, email string, activeContext *UserContext) error {
	if userID == "" {
		return errors.NewValidationError("userID no puede estar vacío")
	}
	if email == "" {
		return errors.NewValidationError("email no puede estar vacío")
	}
	if activeContext == nil {
		return errors.NewValidationError("activeContext no puede ser nil")
	}
	return nil
}

// validateDispatched is ValidateToken's path when a keyring is attached: it
// reads `alg` and `kid` from the unverified header and routes to a parser
// pinned to that single algorithm.
func (m *JWTManager) validateDispatched(tokenString string) (*Claims, error) {
	alg, kid, err := peekHeader(tokenString)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid token")
	}

	switch {
	case alg == jwt.SigningMethodHS256.Alg() && kid == "":
		return m.parseAndValidateToken(tokenString)
	case alg == jwt.SigningMethodES256.Alg() && kid != "":
		return m.parseContextToken(tokenString)
	default:
		// none, RS256, ES384, ES256 without kid, HS256 with kid...
		return nil, errors.NewUnauthorizedError("invalid token")
	}
}

// peekHeader decodes the header without verifying anything. `kid`, when
// present, must be a string.
func peekHeader(tokenString string) (alg, kid string, err error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return "", "", err
	}
	alg, ok := token.Header["alg"].(string)
	if !ok {
		return "", "", stdErrors.New("missing alg")
	}
	if raw, present := token.Header[headerKID]; present {
		kid, ok = raw.(string)
		if !ok || kid == "" {
			return "", "", stdErrors.New("invalid kid")
		}
	}
	return alg, kid, nil
}

// parseContextToken verifies an ES256 Context Token: same issuer, `exp`
// required and leeway as the HS256 path; the key is the one named by `kid`.
func (m *JWTManager) parseContextToken(tokenString string) (*Claims, error) {
	parser := jwt.NewParser(
		jwt.WithIssuer(m.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(30*time.Second),
		jwt.WithTimeFunc(m.clock),
	)

	token, err := parser.ParseWithClaims(tokenString, &Claims{}, m.contextKeyFor)
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

// contextKeyFor picks the accepted public key named by `kid`.
func (m *JWTManager) contextKeyFor(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	kid, ok := token.Header[headerKID].(string)
	if !ok || kid == "" {
		return nil, stdErrors.New("missing kid")
	}
	pub, ok := m.contextKeys.PublicKey(kid)
	if !ok {
		return nil, fmt.Errorf("unknown kid %q", kid)
	}
	return pub, nil
}
