package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Every key of these tests is generated here: no private key is versioned.

func newP256(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	return k
}

func newKeyring(t *testing.T, kid string, priv *ecdsa.PrivateKey, extra map[string]*ecdsa.PublicKey) *ContextKeyring {
	t.Helper()
	kr, err := NewContextKeyring(&ContextSigningKey{ID: kid, Private: priv}, extra)
	require.NoError(t, err)
	return kr
}

func esManager(t *testing.T, kr *ContextKeyring) *JWTManager {
	t.Helper()
	m, err := NewContextTokenManager(testIssuer, kr)
	require.NoError(t, err)
	return m
}

func selfContext() *UserContext {
	return &UserContext{
		RoleID:           "role-teacher",
		RoleName:         "teacher",
		SchoolID:         "school-1",
		SchoolName:       "Escuela Uno",
		AcademicUnitID:   "unit-1",
		AcademicUnitName: "1A",
		Landing:          "dashboard-home",
		Grants:           Grants{Allow: []string{"materials:*", "assessments:read"}, Deny: []string{}},
	}
}

func wardContext() *UserContext {
	return &UserContext{
		RoleID:             "role-guardian",
		RoleName:           "guardian",
		SchoolID:           "school-1",
		SchoolName:         "Escuela Uno",
		Landing:            "guardian-home",
		Grants:             Grants{Allow: []string{"grades:read", "attendance:*"}, Deny: []string{"grades:write", "attendance:delete"}},
		SubjectStudentID:   "student-42",
		SubjectStudentName: "Ana",
		ActorMode:          ActorModeWard,
	}
}

func decodeSegment(t *testing.T, token string, i int) map[string]any {
	t.Helper()
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3)
	raw, err := base64.RawURLEncoding.DecodeString(parts[i])
	require.NoError(t, err)
	out := map[string]any{}
	require.NoError(t, json.Unmarshal(raw, &out))
	return out
}

// rewriteHeader keeps payload and signature and replaces the header.
func rewriteHeader(t *testing.T, token string, header map[string]any) string {
	t.Helper()
	parts := strings.Split(token, ".")
	raw, err := json.Marshal(header)
	require.NoError(t, err)
	parts[0] = base64.RawURLEncoding.EncodeToString(raw)
	return strings.Join(parts, ".")
}

func validClaims(now time.Time, issuer string) Claims {
	return Claims{
		UserID:        "user-1",
		Email:         testEmail,
		ActiveContext: selfContext(),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "jti-1",
			Issuer:    issuer,
			Subject:   "user-1",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
}

func signES(t *testing.T, claims Claims, kid string, priv *ecdsa.PrivateKey) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	if kid != "" {
		tok.Header["kid"] = kid
	}
	s, err := tok.SignedString(priv)
	require.NoError(t, err)
	return s
}

// TestContextToken_GoldenShape freezes the wire shape of the Context Token
// (REQ-C4): header alg/kid and the payload (minus the jti, a uuid).
func TestContextToken_GoldenShape(t *testing.T) {
	fixed := time.Date(2026, 11, 14, 3, 0, 0, 0, time.UTC)
	m := esManager(t, newKeyring(t, "kid-a", newP256(t), nil))
	m.now = func() time.Time { return fixed }

	es, esExp, err := m.GenerateContextToken("user-1", testEmail, wardContext(), 15*time.Minute)
	require.NoError(t, err)
	assert.Equal(t, fixed.Add(15*time.Minute), esExp)

	assert.Equal(t, map[string]any{"alg": "ES256", "typ": "JWT", "kid": "kid-a"}, decodeSegment(t, es, 0))
	payload := decodeSegment(t, es, 1)
	assert.NotEmpty(t, payload["jti"])
	delete(payload, "jti")
	assert.Equal(t, frozenWardPayload(fixed), payload, "frozen wire shape (REQ-C4)")

	// What the middleware stores round-trips with the same context.
	m.now = nil
	tok, _, err := m.GenerateContextToken("user-1", testEmail, wardContext(), 15*time.Minute)
	require.NoError(t, err)
	claims, err := m.ValidateToken(tok)
	require.NoError(t, err)
	assert.Equal(t, wardContext(), claims.ActiveContext)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, testEmail, claims.Email)
}

// frozenWardPayload is the literal wire shape of the Context Token for the
// ward case (minus jti). Any drift in the shared claims builder breaks it.
func frozenWardPayload(now time.Time) map[string]any {
	iat := float64(now.Unix())
	return map[string]any{
		"user_id": "user-1",
		"email":   testEmail,
		"active_context": map[string]any{
			"role_id":              "role-guardian",
			"role_name":            "guardian",
			"school_id":            "school-1",
			"school_name":          "Escuela Uno",
			"landing":              "guardian-home",
			"grants":               map[string]any{"allow": []any{"grades:read", "attendance:*"}, "deny": []any{"grades:write", "attendance:delete"}},
			"subject_student_id":   "student-42",
			"subject_student_name": "Ana",
			"actor_mode":           "ward",
		},
		"iss": testIssuer,
		"sub": "user-1",
		"iat": iat,
		"nbf": iat,
		"exp": iat + 15*60,
	}
}

func TestContextToken_ValidatesAndKeepsIssuer(t *testing.T) {
	m := esManager(t, newKeyring(t, "kid-a", newP256(t), nil))
	tok, _, err := m.GenerateContextToken("user-1", testEmail, wardContext(), time.Hour)
	require.NoError(t, err)

	claims, err := m.ValidateToken(tok)
	require.NoError(t, err)
	assert.Equal(t, testIssuer, claims.Issuer)
	assert.Equal(t, "student-42", claims.ActiveContext.SubjectStudentID)
	assert.Equal(t, ActorModeWard, claims.ActiveContext.ActorMode)
	assert.Equal(t, []string{"grades:write", "attendance:delete"}, claims.ActiveContext.Grants.Deny)
}

func TestNewContextTokenManager_Rejects(t *testing.T) {
	kr := newKeyring(t, "kid-a", newP256(t), nil)
	_, err := NewContextTokenManager(testIssuer, nil)
	require.Error(t, err, "keyring required")
	_, err = NewContextTokenManager("  ", kr)
	require.Error(t, err, "issuer required")
	_, err = NewContextTokenManagerFromConfig(testIssuer, ContextKeyConfig{})
	require.Error(t, err, "empty config")

	m, err := NewContextTokenManager(testIssuer, kr)
	require.NoError(t, err)
	assert.Equal(t, testIssuer, m.Issuer())
	assert.Same(t, kr, m.ContextKeyring())
}

func TestNewContextTokenManagerFromConfig_VerifyOnly(t *testing.T) {
	priv := newP256(t)
	pubPEM, err := MarshalECPublicKeyPEM(&priv.PublicKey)
	require.NoError(t, err)
	m, err := NewContextTokenManagerFromConfig(testIssuer, ContextKeyConfig{
		PublicKeysJSON: mustJSON(t, map[string]string{"kid-a": pubPEM}),
	})
	require.NoError(t, err)
	_, err = m.ValidateToken(signES(t, validClaims(time.Now(), testIssuer), "kid-a", priv))
	require.NoError(t, err)
	_, _, err = m.GenerateContextToken("user-1", testEmail, selfContext(), time.Hour)
	require.Error(t, err, "verify-only manager cannot sign")
}

func TestContextToken_AlgConfusion(t *testing.T) {
	priv := newP256(t)
	kr := newKeyring(t, "kid-a", priv, nil)
	m := esManager(t, kr)
	now := time.Now()

	es, _, err := m.GenerateContextToken("user-1", testEmail, selfContext(), time.Hour)
	require.NoError(t, err)

	der, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	require.NoError(t, err)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})

	hmacWith := func(secret []byte, kid string) string {
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims(now, testIssuer))
		if kid != "" {
			tok.Header["kid"] = kid
		}
		s, sErr := tok.SignedString(secret)
		require.NoError(t, sErr)
		return s
	}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	rs256, err := jwt.NewWithClaims(jwt.SigningMethodRS256, validClaims(now, testIssuer)).SignedString(rsaKey)
	require.NoError(t, err)

	p384, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err)
	es384Tok := jwt.NewWithClaims(jwt.SigningMethodES384, validClaims(now, testIssuer))
	es384Tok.Header["kid"] = "kid-a"
	es384, err := es384Tok.SignedString(p384)
	require.NoError(t, err)

	none, err := jwt.NewWithClaims(jwt.SigningMethodNone, validClaims(now, testIssuer)).
		SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	cases := map[string]string{
		"ES256 token relabelled HS256 (kid kept)":   rewriteHeader(t, es, map[string]any{"alg": "HS256", "typ": "JWT", "kid": "kid-a"}),
		"ES256 token relabelled HS256 (kid gone)":   rewriteHeader(t, es, map[string]any{"alg": "HS256", "typ": "JWT"}),
		"HS256 token relabelled ES256 (known kid)":  rewriteHeader(t, hmacWith([]byte(testSecretKey), ""), map[string]any{"alg": "ES256", "typ": "JWT", "kid": "kid-a"}),
		"HS256 user token (retired in T6.2)":        hmacWith([]byte(testSecretKey), ""),
		"HMAC signed with public PEM, no kid":       hmacWith(pubPEM, ""),
		"HMAC signed with public PEM, known kid":    hmacWith(pubPEM, "kid-a"),
		"HMAC signed with public DER, known kid":    hmacWith(der, "kid-a"),
		"valid HS256 secret but with a kid":         hmacWith([]byte(testSecretKey), "kid-a"),
		"alg none":                                  none,
		"RS256":                                     rs256,
		"ES384 with known kid":                      es384,
		"missing alg":                               rewriteHeader(t, es, map[string]any{"typ": "JWT", "kid": "kid-a"}),
		"non-string kid":                            rewriteHeader(t, es, map[string]any{"alg": "ES256", "typ": "JWT", "kid": 7}),
		"garbage":                                   "not.a.jwt",
		"ES256 token into HS256 parser via relabel": rewriteHeader(t, es, map[string]any{"alg": "HS256"}),
	}
	for name, tok := range cases {
		t.Run(name, func(t *testing.T) {
			_, vErr := m.ValidateToken(tok)
			require.Error(t, vErr)
		})
	}

}

func TestContextToken_KIDChecks(t *testing.T) {
	priv := newP256(t)
	other := newP256(t)
	m := esManager(t, newKeyring(t, "kid-a", priv, map[string]*ecdsa.PublicKey{"kid-b": &other.PublicKey}))
	now := time.Now()

	t.Run("missing kid", func(t *testing.T) {
		_, err := m.ValidateToken(signES(t, validClaims(now, testIssuer), "", priv))
		require.Error(t, err)
	})
	t.Run("unknown kid", func(t *testing.T) {
		_, err := m.ValidateToken(signES(t, validClaims(now, testIssuer), "kid-zzz", priv))
		require.Error(t, err)
	})
	t.Run("kid naming a different key than the signer", func(t *testing.T) {
		_, err := m.ValidateToken(signES(t, validClaims(now, testIssuer), "kid-b", priv))
		require.Error(t, err)
	})
	t.Run("signer key the keyring does not hold", func(t *testing.T) {
		_, err := m.ValidateToken(signES(t, validClaims(now, testIssuer), "kid-a", newP256(t)))
		require.Error(t, err)
	})
	t.Run("hand-built valid token passes", func(t *testing.T) {
		_, err := m.ValidateToken(signES(t, validClaims(now, testIssuer), "kid-b", other))
		require.NoError(t, err)
	})
}

func TestContextToken_ClaimChecks(t *testing.T) {
	priv := newP256(t)
	m := esManager(t, newKeyring(t, "kid-a", priv, nil))

	t.Run("expired beyond leeway", func(t *testing.T) {
		c := validClaims(time.Now().Add(-2*time.Hour), testIssuer)
		c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute))
		_, err := m.ValidateToken(signES(t, c, "kid-a", priv))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token expired")
	})
	t.Run("wrong issuer", func(t *testing.T) {
		_, err := m.ValidateToken(signES(t, validClaims(time.Now(), "someone-else"), "kid-a", priv))
		require.Error(t, err)
	})
	t.Run("no exp", func(t *testing.T) {
		c := validClaims(time.Now(), testIssuer)
		c.ExpiresAt = nil
		_, err := m.ValidateToken(signES(t, c, "kid-a", priv))
		require.Error(t, err)
	})
	t.Run("no active_context", func(t *testing.T) {
		c := validClaims(time.Now(), testIssuer)
		c.ActiveContext = nil
		_, err := m.ValidateToken(signES(t, c, "kid-a", priv))
		require.Error(t, err)
	})
	t.Run("GenerateContextTokenUntil allows under a minute but not the past", func(t *testing.T) {
		tok, _, err := m.GenerateContextTokenUntil("user-1", testEmail, selfContext(), time.Now().Add(20*time.Second))
		require.NoError(t, err)
		_, err = m.ValidateToken(tok)
		require.NoError(t, err)
		_, _, err = m.GenerateContextTokenUntil("user-1", testEmail, selfContext(), time.Now().Add(-time.Second))
		require.Error(t, err)
	})
	t.Run("input validation", func(t *testing.T) {
		_, _, err := m.GenerateContextToken("", testEmail, selfContext(), time.Hour)
		require.Error(t, err)
		_, _, err = m.GenerateContextToken("u", "", selfContext(), time.Hour)
		require.Error(t, err)
		_, _, err = m.GenerateContextToken("u", testEmail, nil, time.Hour)
		require.Error(t, err)
		_, _, err = m.GenerateContextToken("u", testEmail, selfContext(), time.Second)
		require.Error(t, err)
	})
}

func TestContextToken_Rotation(t *testing.T) {
	oldKey, newKey := newP256(t), newP256(t)

	before := esManager(t, newKeyring(t, "kid-1", oldKey, nil))
	oldTok, _, err := before.GenerateContextToken("user-1", testEmail, selfContext(), time.Hour)
	require.NoError(t, err)

	// Rotation: kid-2 signs, kid-1 is still accepted.
	during := esManager(t, newKeyring(t, "kid-2", newKey, map[string]*ecdsa.PublicKey{"kid-1": &oldKey.PublicKey}))
	newTok, _, err := during.GenerateContextToken("user-1", testEmail, selfContext(), time.Hour)
	require.NoError(t, err)
	assert.Equal(t, "kid-2", decodeSegment(t, newTok, 0)["kid"])
	_, err = during.ValidateToken(oldTok)
	require.NoError(t, err)
	_, err = during.ValidateToken(newTok)
	require.NoError(t, err)
	assert.Equal(t, []string{"kid-1", "kid-2"}, during.ContextKeyring().KIDs())

	// A verify-only API with both public keys accepts both.
	verifier, err := NewContextKeyring(nil, map[string]*ecdsa.PublicKey{"kid-1": &oldKey.PublicKey, "kid-2": &newKey.PublicKey})
	require.NoError(t, err)
	assert.False(t, verifier.CanSign())
	vm := esManager(t, verifier)
	_, err = vm.ValidateToken(oldTok)
	require.NoError(t, err)
	_, err = vm.ValidateToken(newTok)
	require.NoError(t, err)
	_, _, err = vm.GenerateContextToken("user-1", testEmail, selfContext(), time.Hour)
	require.Error(t, err, "verify-only keyring cannot sign")

	// After: kid-1 removed, its tokens fail.
	after := esManager(t, newKeyring(t, "kid-2", newKey, nil))
	_, err = after.ValidateToken(oldTok)
	require.Error(t, err)
	_, err = after.ValidateToken(newTok)
	require.NoError(t, err)
}

func TestNewContextKeyring_Rejects(t *testing.T) {
	p256 := newP256(t)
	p384, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err)

	var nilKR *ContextKeyring
	assert.False(t, nilKR.CanSign())
	assert.Empty(t, nilKR.KIDs())
	_, ok := nilKR.PublicKey("a")
	assert.False(t, ok)

	_, err = NewContextKeyring(nil, nil)
	require.Error(t, err, "no key")
	_, err = NewContextKeyring(&ContextSigningKey{ID: "", Private: p256}, nil)
	require.Error(t, err, "empty signing kid")
	_, err = NewContextKeyring(&ContextSigningKey{ID: "a", Private: p384}, nil)
	require.Error(t, err, "P-384 signing key")
	_, err = NewContextKeyring(nil, map[string]*ecdsa.PublicKey{"a": &p384.PublicKey})
	require.Error(t, err, "P-384 public key")
	_, err = NewContextKeyring(nil, map[string]*ecdsa.PublicKey{"": &p256.PublicKey})
	require.Error(t, err, "empty public kid")
	_, err = NewContextKeyring(&ContextSigningKey{ID: "a", Private: p256}, map[string]*ecdsa.PublicKey{"a": &p256.PublicKey})
	require.Error(t, err, "duplicate kid")
}

func pemPKCS8(t *testing.T, k *ecdsa.PrivateKey) string {
	t.Helper()
	der, err := x509.MarshalPKCS8PrivateKey(k)
	require.NoError(t, err)
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
}

func TestContextKeyConfig(t *testing.T) {
	signer, previous := newP256(t), newP256(t)
	prevPEM, err := MarshalECPublicKeyPEM(&previous.PublicKey)
	require.NoError(t, err)
	pubJSON, err := json.Marshal(map[string]string{"kid-1": prevPEM})
	require.NoError(t, err)

	t.Run("nothing set is zero", func(t *testing.T) {
		cfg, cErr := ContextKeyConfigFromEnv(func(string) (string, bool) { return "", false })
		require.NoError(t, cErr)
		assert.True(t, cfg.IsZero())
		_, cErr = NewContextKeyringFromConfig(cfg)
		require.Error(t, cErr)
	})

	t.Run("signing key from file, public keys from env", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "signing.pem")
		require.NoError(t, os.WriteFile(path, []byte(pemPKCS8(t, signer)), 0o600))
		env := map[string]string{
			EnvContextSigningKey + "_FILE": path,
			EnvContextSigningKID:           "kid-2",
			EnvContextPublicKeys:           string(pubJSON),
		}
		cfg, cErr := ContextKeyConfigFromEnv(func(k string) (string, bool) { v, ok := env[k]; return v, ok })
		require.NoError(t, cErr)
		kr, cErr := NewContextKeyringFromConfig(cfg)
		require.NoError(t, cErr)
		assert.Equal(t, "kid-2", kr.ActiveKID())
		assert.Equal(t, []string{"kid-1", "kid-2"}, kr.KIDs())
	})

	t.Run("SEC1 PEM accepted", func(t *testing.T) {
		der, mErr := x509.MarshalECPrivateKey(signer)
		require.NoError(t, mErr)
		sec1 := string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der}))
		kr, cErr := NewContextKeyringFromConfig(ContextKeyConfig{SigningKeyPEM: sec1, SigningKID: "k"})
		require.NoError(t, cErr)
		assert.True(t, kr.CanSign())
	})

	t.Run("verify-only from public keys", func(t *testing.T) {
		kr, cErr := NewContextKeyringFromConfig(ContextKeyConfig{PublicKeysJSON: string(pubJSON)})
		require.NoError(t, cErr)
		assert.False(t, kr.CanSign())
	})

	bad := map[string]ContextKeyConfig{
		"key without kid":       {SigningKeyPEM: pemPKCS8(t, signer)},
		"kid without key":       {SigningKID: "k"},
		"not PEM":               {SigningKeyPEM: "nope", SigningKID: "k"},
		"public JSON malformed": {PublicKeysJSON: "{"},
		"public not PEM":        {PublicKeysJSON: `{"k":"nope"}`},
		"public given private":  {PublicKeysJSON: mustJSON(t, map[string]string{"k": pemPKCS8(t, signer)})},
	}
	for name, cfg := range bad {
		t.Run("invalid: "+name, func(t *testing.T) {
			_, cErr := NewContextKeyringFromConfig(cfg)
			require.Error(t, cErr)
		})
	}

	t.Run("value and _FILE together", func(t *testing.T) {
		env := map[string]string{EnvContextSigningKID: "a", EnvContextSigningKID + "_FILE": "/x"}
		_, cErr := ContextKeyConfigFromEnv(func(k string) (string, bool) { v, ok := env[k]; return v, ok })
		require.Error(t, cErr)
	})
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return string(b)
}
