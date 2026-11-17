// Package authtest is the single test helper of the EduGo ecosystem for user
// tokens (plan 053 T6.2). It generates an ES256 (P-256) key pair in memory, so
// no key is ever versioned, and gives what a test needs on both sides:
//
//   - to mint a Context Token with user and context claims (Token, Sign);
//   - to configure the verifier exactly like production (PublicKeysJSON, the
//     value of AUTH_CONTEXT_PUBLIC_KEYS, or Verifier).
//
// Default returns one key pair per process, so a test that mints tokens in one
// file and builds the server in another sees the same keys. New returns a
// fresh pair (for negative tests: a token signed by a key the server does not
// know).
//
// It does not import the testing package and never panics except in MustNew
// and Default, so it can be used from TestMain and package-level helpers.
package authtest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/EduGoGroup/edugo-shared/auth"
)

const (
	// DefaultIssuer is the `iss` of the ecosystem's user tokens.
	DefaultIssuer = "edugo-central"
	// DefaultKID is the kid of every key pair built by this package.
	DefaultKID = "authtest"
)

// KeyPair is an in-memory ES256 key pair plus the managers built from it.
type KeyPair struct {
	issuer         string
	kid            string
	private        *ecdsa.PrivateKey
	privatePEM     string
	publicKeysJSON string
	manager        *auth.JWTManager
	verifier       *auth.JWTManager
}

// New generates a fresh P-256 key pair for issuer (DefaultIssuer if empty).
func New(issuer string) (*KeyPair, error) {
	if issuer == "" {
		issuer = DefaultIssuer
	}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("authtest: generate key: %w", err)
	}

	pubPEM, err := auth.MarshalECPublicKeyPEM(&priv.PublicKey)
	if err != nil {
		return nil, err
	}
	pubJSON, err := json.Marshal(map[string]string{DefaultKID: pubPEM})
	if err != nil {
		return nil, fmt.Errorf("authtest: marshal public keys: %w", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("authtest: marshal private key: %w", err)
	}

	keyring, err := auth.NewContextKeyring(&auth.ContextSigningKey{ID: DefaultKID, Private: priv}, nil)
	if err != nil {
		return nil, err
	}
	manager, err := auth.NewContextTokenManager(issuer, keyring)
	if err != nil {
		return nil, err
	}
	// The verifier goes through the same path as a domain API in production.
	verifier, err := auth.NewContextTokenManagerFromConfig(issuer, auth.ContextKeyConfig{PublicKeysJSON: string(pubJSON)})
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		issuer:         issuer,
		kid:            DefaultKID,
		private:        priv,
		privatePEM:     string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})),
		publicKeysJSON: string(pubJSON),
		manager:        manager,
		verifier:       verifier,
	}, nil
}

// MustNew is New that panics on error.
func MustNew(issuer string) *KeyPair {
	k, err := New(issuer)
	if err != nil {
		panic(err)
	}
	return k
}

var (
	defaultOnce sync.Once
	defaultKeys *KeyPair
)

// Default returns the process-wide key pair for DefaultIssuer, generated on
// first use.
func Default() *KeyPair {
	defaultOnce.Do(func() { defaultKeys = MustNew(DefaultIssuer) })
	return defaultKeys
}

// Issuer returns the `iss` of the tokens this pair signs.
func (k *KeyPair) Issuer() string { return k.issuer }

// KID returns the kid written in the header of every token.
func (k *KeyPair) KID() string { return k.kid }

// PublicKeysJSON returns {"<kid>": "<PKIX PEM>"}, the value of
// AUTH_CONTEXT_PUBLIC_KEYS for a verifier of this pair.
func (k *KeyPair) PublicKeysJSON() string { return k.publicKeysJSON }

// PrivateKeyPEM returns the PKCS#8 PEM of the private key, the value of
// AUTH_CONTEXT_SIGNING_KEY for a process that must sign (with
// AUTH_CONTEXT_SIGNING_KID = KID()).
func (k *KeyPair) PrivateKeyPEM() string { return k.privatePEM }

// Env returns the verify-only environment of a child process:
// AUTH_CONTEXT_PUBLIC_KEYS.
func (k *KeyPair) Env() map[string]string {
	return map[string]string{auth.EnvContextPublicKeys: k.publicKeysJSON}
}

// Manager returns a manager that signs and verifies with this pair.
func (k *KeyPair) Manager() *auth.JWTManager { return k.manager }

// Verifier returns a verify-only manager built from PublicKeysJSON, as a
// domain API builds it.
func (k *KeyPair) Verifier() *auth.JWTManager { return k.verifier }

// Token mints a valid Context Token for the user and context that expires in
// ttl (at least one minute).
func (k *KeyPair) Token(userID, email string, activeContext *auth.UserContext, ttl time.Duration) (string, error) {
	token, _, err := k.manager.GenerateContextToken(userID, email, activeContext, ttl)
	return token, err
}

// Sign signs arbitrary claims with this pair (ES256, kid in the header),
// without any validation: for negative tests (expired, wrong issuer, no
// active_context...).
func (k *KeyPair) Sign(claims auth.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = k.kid
	signed, err := token.SignedString(k.private)
	if err != nil {
		return "", fmt.Errorf("authtest: sign: %w", err)
	}
	return signed, nil
}
