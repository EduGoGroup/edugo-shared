package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	stdErrors "errors"
	"fmt"
	"os"
	"strings"
)

// Environment variables of the Context Token keyring (plan 053 T2.1). They
// follow the `auth.context.*` -> `AUTH_CONTEXT_*` mapping of the config
// loader. No key is ever versioned: the values come from Secret Manager or a
// local file outside git.
const (
	// EnvContextSigningKey holds the PEM (PKCS#8 or SEC1) of the active P-256
	// private key. Only edugo-api-identity, which issues the token, sets it.
	EnvContextSigningKey = "AUTH_CONTEXT_SIGNING_KEY"
	// EnvContextSigningKID holds the kid of the active signing key.
	EnvContextSigningKID = "AUTH_CONTEXT_SIGNING_KID"
	// EnvContextPublicKeys holds a JSON object {"<kid>": "<PKIX PEM>", ...}
	// with the accepted public keys (the domain APIs, and the previous key
	// during a rotation).
	EnvContextPublicKeys = "AUTH_CONTEXT_PUBLIC_KEYS"
	// envFileSuffix: each variable also accepts <NAME>_FILE with a file path.
	envFileSuffix = "_FILE"
)

// ContextKeyConfig is the raw configuration of the keyring.
type ContextKeyConfig struct {
	// SigningKeyPEM is the PEM of the active P-256 private key (optional).
	SigningKeyPEM string
	// SigningKID is the kid of SigningKeyPEM (required when it is set).
	SigningKID string
	// PublicKeysJSON is a JSON object kid -> PKIX PEM (optional).
	PublicKeysJSON string
}

// IsZero reports whether nothing is configured: the API keeps HS256 only.
func (c ContextKeyConfig) IsZero() bool {
	return strings.TrimSpace(c.SigningKeyPEM) == "" &&
		strings.TrimSpace(c.SigningKID) == "" &&
		strings.TrimSpace(c.PublicKeysJSON) == ""
}

// ContextKeyConfigFromEnv reads the AUTH_CONTEXT_* variables (or their _FILE
// variants) through lookup; pass os.LookupEnv in production. Setting both a
// variable and its _FILE variant is an error.
func ContextKeyConfigFromEnv(lookup func(string) (string, bool)) (ContextKeyConfig, error) {
	if lookup == nil {
		lookup = os.LookupEnv
	}
	var cfg ContextKeyConfig
	var err error
	if cfg.SigningKeyPEM, err = envOrFile(lookup, EnvContextSigningKey); err != nil {
		return ContextKeyConfig{}, err
	}
	if cfg.SigningKID, err = envOrFile(lookup, EnvContextSigningKID); err != nil {
		return ContextKeyConfig{}, err
	}
	if cfg.PublicKeysJSON, err = envOrFile(lookup, EnvContextPublicKeys); err != nil {
		return ContextKeyConfig{}, err
	}
	cfg.SigningKID = strings.TrimSpace(cfg.SigningKID)
	return cfg, nil
}

func envOrFile(lookup func(string) (string, bool), name string) (string, error) {
	value, hasValue := lookup(name)
	path, hasPath := lookup(name + envFileSuffix)
	hasValue = hasValue && value != ""
	hasPath = hasPath && path != ""
	switch {
	case hasValue && hasPath:
		return "", fmt.Errorf("context keys: set %s or %s%s, not both", name, name, envFileSuffix)
	case hasPath:
		raw, err := os.ReadFile(path) // #nosec G304 -- path comes from the operator's config
		if err != nil {
			return "", fmt.Errorf("context keys: read %s%s: %w", name, envFileSuffix, err)
		}
		return string(raw), nil
	default:
		return value, nil
	}
}

// NewContextKeyringFromConfig builds the keyring from cfg. A zero config is
// an error: callers check IsZero first and keep HS256 only. A config that is
// set but invalid is always an error, never a silent fallback to HS256.
func NewContextKeyringFromConfig(cfg ContextKeyConfig) (*ContextKeyring, error) {
	if cfg.IsZero() {
		return nil, stdErrors.New("context keys: nothing configured")
	}

	var signing *ContextSigningKey
	hasKey := strings.TrimSpace(cfg.SigningKeyPEM) != ""
	hasKID := strings.TrimSpace(cfg.SigningKID) != ""
	switch {
	case hasKey && hasKID:
		priv, err := ParseECPrivateKeyPEM([]byte(cfg.SigningKeyPEM))
		if err != nil {
			return nil, err
		}
		signing = &ContextSigningKey{ID: strings.TrimSpace(cfg.SigningKID), Private: priv}
	case hasKey != hasKID:
		return nil, fmt.Errorf("context keys: %s and %s go together", EnvContextSigningKey, EnvContextSigningKID)
	}

	public, err := ParseContextPublicKeysJSON(cfg.PublicKeysJSON)
	if err != nil {
		return nil, err
	}
	return NewContextKeyring(signing, public)
}

// ParseContextPublicKeysJSON parses a JSON object kid -> PKIX PEM. An empty
// string yields an empty map.
func ParseContextPublicKeysJSON(raw string) (map[string]*ecdsa.PublicKey, error) {
	out := map[string]*ecdsa.PublicKey{}
	if strings.TrimSpace(raw) == "" {
		return out, nil
	}
	var pems map[string]string
	if err := json.Unmarshal([]byte(raw), &pems); err != nil {
		return nil, fmt.Errorf("context keys: %s is not a JSON object kid -> PEM: %w", EnvContextPublicKeys, err)
	}
	for kid, p := range pems {
		pub, err := ParseECPublicKeyPEM([]byte(p))
		if err != nil {
			return nil, fmt.Errorf("context keys: public key %q: %w", kid, err)
		}
		out[kid] = pub
	}
	return out, nil
}

// ParseECPrivateKeyPEM parses a P-256 private key in PKCS#8 ("PRIVATE KEY")
// or SEC1 ("EC PRIVATE KEY") PEM.
func ParseECPrivateKeyPEM(data []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, stdErrors.New("context keys: signing key is not PEM")
	}
	var priv *ecdsa.PrivateKey
	switch block.Type {
	case "EC PRIVATE KEY":
		k, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("context keys: parse SEC1 key: %w", err)
		}
		priv = k
	case "PRIVATE KEY":
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("context keys: parse PKCS#8 key: %w", err)
		}
		ec, ok := k.(*ecdsa.PrivateKey)
		if !ok {
			return nil, stdErrors.New("context keys: signing key is not ECDSA")
		}
		priv = ec
	default:
		return nil, fmt.Errorf("context keys: unexpected PEM block %q", block.Type)
	}
	if !isP256Public(&priv.PublicKey) {
		return nil, stdErrors.New("context keys: signing key is not P-256")
	}
	return priv, nil
}

// ParseECPublicKeyPEM parses a P-256 public key in PKIX ("PUBLIC KEY") PEM.
func ParseECPublicKeyPEM(data []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, stdErrors.New("context keys: public key is not a PKIX PEM")
	}
	k, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("context keys: parse PKIX key: %w", err)
	}
	pub, ok := k.(*ecdsa.PublicKey)
	if !ok || !isP256Public(pub) {
		return nil, stdErrors.New("context keys: public key is not P-256")
	}
	return pub, nil
}

// MarshalECPublicKeyPEM encodes a public key as PKIX PEM, the format of
// AUTH_CONTEXT_PUBLIC_KEYS.
func MarshalECPublicKeyPEM(pub *ecdsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("context keys: marshal public key: %w", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})), nil
}
