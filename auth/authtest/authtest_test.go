package authtest_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/auth/authtest"
)

func userContext() *auth.UserContext {
	return &auth.UserContext{RoleID: "r1", RoleName: "teacher", SchoolID: "s1", Grants: auth.Grants{Allow: []string{"materials:*"}}}
}

func TestTokenVerifiesWithPublicKeysJSON(t *testing.T) {
	keys := authtest.MustNew("")
	assert.Equal(t, authtest.DefaultIssuer, keys.Issuer())

	token, err := keys.Token("u1", "u1@edugo.test", userContext(), time.Hour)
	require.NoError(t, err)

	// A verifier configured the way an API is: only the public keys JSON.
	cfg, err := auth.ContextKeyConfigFromEnv(func(name string) (string, bool) {
		v, ok := keys.Env()[name]
		return v, ok
	})
	require.NoError(t, err)
	verifier, err := auth.NewContextTokenManagerFromConfig(keys.Issuer(), cfg)
	require.NoError(t, err)

	claims, err := verifier.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "u1", claims.UserID)
	assert.Equal(t, "s1", claims.ActiveContext.SchoolID)

	_, err = keys.Verifier().ValidateToken(token)
	require.NoError(t, err)
}

func TestPrivateKeyPEMSigns(t *testing.T) {
	keys := authtest.MustNew("iss")
	signer, err := auth.NewContextTokenManagerFromConfig("iss", auth.ContextKeyConfig{
		SigningKeyPEM: keys.PrivateKeyPEM(),
		SigningKID:    keys.KID(),
	})
	require.NoError(t, err)
	token, _, err := signer.GenerateContextToken("u1", "u1@edugo.test", userContext(), time.Hour)
	require.NoError(t, err)
	_, err = keys.Verifier().ValidateToken(token)
	require.NoError(t, err)
}

func TestDefaultIsStableAndNewIsFresh(t *testing.T) {
	assert.Same(t, authtest.Default(), authtest.Default())

	other := authtest.MustNew("")
	token, err := other.Token("u1", "u1@edugo.test", userContext(), time.Hour)
	require.NoError(t, err)
	_, err = authtest.Default().Verifier().ValidateToken(token)
	require.Error(t, err, "a token from another pair must be rejected")
}

func TestSignArbitraryClaims(t *testing.T) {
	keys := authtest.MustNew("")
	now := time.Now()
	expired, err := keys.Sign(auth.Claims{
		UserID:        "u1",
		Email:         "u1@edugo.test",
		ActiveContext: userContext(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    keys.Issuer(),
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
		},
	})
	require.NoError(t, err)
	_, err = keys.Verifier().ValidateToken(expired)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "token expired")
}
