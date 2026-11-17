package auth

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// testSecretKey only builds HS256 tokens that the manager must reject.
	testSecretKey = "test-secret-key-at-least-32-chars-long-for-security"
	testIssuer    = "edugo-test"
	testEmail     = "test@example.com"
	testKID       = "kid-test"
)

// testManager builds an ES256 manager with a fresh in-memory key and returns
// the key so tests can hand-sign tokens.
func testManager(t *testing.T, issuer string) (*JWTManager, *ecdsa.PrivateKey) {
	t.Helper()
	priv := newP256(t)
	kr, err := NewContextKeyring(&ContextSigningKey{ID: testKID, Private: priv}, nil)
	require.NoError(t, err)
	m, err := NewContextTokenManager(issuer, kr)
	require.NoError(t, err)
	return m, priv
}

// signTestES signs claims with priv under testKID.
func signTestES(priv *ecdsa.PrivateKey, claims Claims) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tok.Header["kid"] = testKID
	return tok.SignedString(priv)
}

func TestGenerateContextToken(t *testing.T) {
	manager, _ := testManager(t, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	expiresIn := 24 * time.Hour

	activeContext := &UserContext{
		RoleID:     "role-123",
		RoleName:   "Teacher",
		SchoolID:   "school-456",
		SchoolName: "Test School",
		Grants:     Grants{Allow: []string{"admin.users.read", "content.materials.create", "content.assessments.grade"}},
	}

	t.Run("genera token con contexto válido exitosamente", func(t *testing.T) {
		token, expiresAt, err := manager.GenerateContextToken(userID, email, activeContext, expiresIn)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.False(t, expiresAt.IsZero())
		assert.True(t, expiresAt.After(time.Now()))
	})

	t.Run("token contiene claims con ActiveContext correctos", func(t *testing.T) {
		token, _, err := manager.GenerateContextToken(userID, email, activeContext, expiresIn)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.NotNil(t, claims.ActiveContext)
		assert.Equal(t, activeContext.RoleID, claims.ActiveContext.RoleID)
		assert.Equal(t, activeContext.RoleName, claims.ActiveContext.RoleName)
		assert.Equal(t, activeContext.SchoolID, claims.ActiveContext.SchoolID)
		assert.Equal(t, activeContext.SchoolName, claims.ActiveContext.SchoolName)
		assert.Equal(t, activeContext.Grants, claims.ActiveContext.Grants)
	})

	t.Run("rechaza userID vacío", func(t *testing.T) {
		_, _, err := manager.GenerateContextToken("", email, activeContext, expiresIn)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "userID no puede estar vacío")
	})

	t.Run("rechaza email vacío", func(t *testing.T) {
		_, _, err := manager.GenerateContextToken(userID, "", activeContext, expiresIn)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "email no puede estar vacío")
	})

	t.Run("rechaza activeContext nil", func(t *testing.T) {
		_, _, err := manager.GenerateContextToken(userID, email, nil, expiresIn)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "activeContext no puede ser nil")
	})

	t.Run("rechaza expiresIn menor a 1 minuto", func(t *testing.T) {
		_, _, err := manager.GenerateContextToken(userID, email, activeContext, 30*time.Second)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "expiresIn debe ser mayor a 1 minuto")
	})

	t.Run("acepta expiresIn de exactamente 1 minuto", func(t *testing.T) {
		token, _, err := manager.GenerateContextToken(userID, email, activeContext, time.Minute)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("genera tokens únicos con mismo contexto", func(t *testing.T) {
		token1, _, err1 := manager.GenerateContextToken(userID, email, activeContext, expiresIn)
		time.Sleep(10 * time.Millisecond)
		token2, _, err2 := manager.GenerateContextToken(userID, email, activeContext, expiresIn)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2, "Los tokens deben ser únicos debido al timestamp y JTI")
	})

	t.Run("token sin SchoolID (context mínimo)", func(t *testing.T) {
		minimalContext := &UserContext{
			RoleID:   "role-789",
			RoleName: "Admin",
			Grants:   Grants{Allow: []string{"admin.users.create", "admin.schools.manage"}},
		}

		token, _, err := manager.GenerateContextToken(userID, email, minimalContext, expiresIn)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, minimalContext.RoleID, claims.ActiveContext.RoleID)
		assert.Equal(t, minimalContext.RoleName, claims.ActiveContext.RoleName)
		assert.Empty(t, claims.ActiveContext.SchoolID)
		assert.Empty(t, claims.ActiveContext.SchoolName)
		assert.Equal(t, minimalContext.Grants, claims.ActiveContext.Grants)
	})

	t.Run("token con unidad académica", func(t *testing.T) {
		contextWithUnit := &UserContext{
			RoleID:           "role-abc",
			RoleName:         "Unit Coordinator",
			SchoolID:         "school-def",
			SchoolName:       "Test School",
			AcademicUnitID:   "unit-ghi",
			AcademicUnitName: "Computer Science",
			Grants:           Grants{Allow: []string{"academic.units.create", "content.materials.create"}},
		}

		token, _, err := manager.GenerateContextToken(userID, email, contextWithUnit, expiresIn)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, contextWithUnit.AcademicUnitID, claims.ActiveContext.AcademicUnitID)
		assert.Equal(t, contextWithUnit.AcademicUnitName, claims.ActiveContext.AcademicUnitName)
	})

	t.Run("tiempo de expiración es correcto", func(t *testing.T) {
		expiresIn := 2 * time.Hour
		_, expiresAt, err := manager.GenerateContextToken(userID, email, activeContext, expiresIn)

		require.NoError(t, err)

		expectedExpiration := time.Now().Add(expiresIn)
		// Tolerancia de 5 segundos
		assert.WithinDuration(t, expectedExpiration, expiresAt, 5*time.Second)
	})

	t.Run("token con lista vacía de permisos", func(t *testing.T) {
		contextNoPerms := &UserContext{
			RoleID:   "role-no-perms",
			RoleName: "Guest",
			Grants:   Grants{Allow: []string{}, Deny: []string{}},
		}

		token, _, err := manager.GenerateContextToken(userID, email, contextNoPerms, expiresIn)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		assert.NotNil(t, claims.ActiveContext.Grants.Allow)
		assert.Empty(t, claims.ActiveContext.Grants.Allow)
		assert.Empty(t, claims.ActiveContext.Grants.Deny)
	})
}

func TestValidateToken(t *testing.T) {
	manager, managerKey := testManager(t, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	activeContext := &UserContext{
		RoleID:   "role-123",
		RoleName: "Teacher",
		Grants:   Grants{Allow: []string{"admin.users.read"}},
	}

	t.Run("valida token válido exitosamente", func(t *testing.T) {
		token, _, err := manager.GenerateContextToken(userID, email, activeContext, 24*time.Hour)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.NotNil(t, claims.ActiveContext)
	})

	t.Run("rechaza token vacío", func(t *testing.T) {
		_, err := manager.ValidateToken("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("rechaza token malformado", func(t *testing.T) {
		_, err := manager.ValidateToken("invalid-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("rechaza token con firma incorrecta", func(t *testing.T) {
		wrongManager, _ := testManager(t, testIssuer)
		token, _, err := wrongManager.GenerateContextToken(userID, email, activeContext, 24*time.Hour)
		require.NoError(t, err)

		_, err = manager.ValidateToken(token)
		assert.Error(t, err)
	})

	t.Run("rechaza token expirado", func(t *testing.T) {
		// Crear un token manualmente que ya esté expirado
		now := time.Now()
		expiredClaims := Claims{
			UserID:        userID,
			Email:         email,
			ActiveContext: activeContext,
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        uuid.New().String(),
				Issuer:    testIssuer,
				Subject:   userID,
				IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
				ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)),
				NotBefore: jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			},
		}
		expiredTokenString, err := signTestES(managerKey, expiredClaims)
		require.NoError(t, err)

		_, err = manager.ValidateToken(expiredTokenString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token expired")
	})

	t.Run("rechaza token sin ActiveContext", func(t *testing.T) {
		// Crear un token manualmente sin ActiveContext
		now := time.Now()
		claims := Claims{
			UserID: userID,
			Email:  email,
			// ActiveContext es nil
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        uuid.New().String(),
				Issuer:    testIssuer,
				Subject:   userID,
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				NotBefore: jwt.NewNumericDate(now),
			},
		}

		tokenString, err := signTestES(managerKey, claims)
		require.NoError(t, err)

		_, err = manager.ValidateToken(tokenString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing active context")
	})
}

func TestExtractUserID(t *testing.T) {
	manager, managerKey := testManager(t, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	activeContext := &UserContext{
		RoleID:   "role-123",
		RoleName: "Teacher",
		Grants:   Grants{Allow: []string{"admin.users.read"}},
	}

	t.Run("extrae userID de token válido", func(t *testing.T) {
		token, _, err := manager.GenerateContextToken(userID, email, activeContext, 24*time.Hour)
		require.NoError(t, err)

		extractedID, err := ExtractUserID(token)
		require.NoError(t, err)
		assert.Equal(t, userID, extractedID)
	})

	t.Run("extrae userID de token expirado (sin validar)", func(t *testing.T) {
		// Crear un token manualmente que ya esté expirado
		now := time.Now()
		expiredClaims := Claims{
			UserID:        userID,
			Email:         email,
			ActiveContext: activeContext,
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        uuid.New().String(),
				Issuer:    testIssuer,
				Subject:   userID,
				IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
				ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)),
				NotBefore: jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			},
		}
		expiredTokenString, err := signTestES(managerKey, expiredClaims)
		require.NoError(t, err)

		extractedID, err := ExtractUserID(expiredTokenString)
		require.NoError(t, err)
		assert.Equal(t, userID, extractedID,
			"ExtractUserID debe funcionar incluso con tokens expirados")
	})

	t.Run("falla con token malformado", func(t *testing.T) {
		_, err := ExtractUserID("invalid-token")
		assert.Error(t, err)
	})

	t.Run("falla con token vacío", func(t *testing.T) {
		_, err := ExtractUserID("")
		assert.Error(t, err)
	})
}

func TestUserContext(t *testing.T) {
	t.Run("UserContext con todos los campos", func(t *testing.T) {
		ctx := &UserContext{
			RoleID:           "role-123",
			RoleName:         "Teacher",
			SchoolID:         "school-456",
			SchoolName:       "Test School",
			AcademicUnitID:   "unit-789",
			AcademicUnitName: "Mathematics",
			Grants:           Grants{Allow: []string{"admin.users.read", "content.materials.create"}},
		}

		assert.NotNil(t, ctx)
		assert.Equal(t, "role-123", ctx.RoleID)
		assert.Equal(t, "Teacher", ctx.RoleName)
		assert.Equal(t, "school-456", ctx.SchoolID)
		assert.Equal(t, "Test School", ctx.SchoolName)
		assert.Equal(t, "unit-789", ctx.AcademicUnitID)
		assert.Equal(t, "Mathematics", ctx.AcademicUnitName)
		assert.Len(t, ctx.Grants.Allow, 2)
	})

	t.Run("UserContext solo con campos requeridos", func(t *testing.T) {
		ctx := &UserContext{
			RoleID:   "role-abc",
			RoleName: "Admin",
			Grants:   Grants{Allow: []string{"admin.users.create"}},
		}

		assert.Equal(t, "role-abc", ctx.RoleID)
		assert.Equal(t, "Admin", ctx.RoleName)
		assert.Empty(t, ctx.SchoolID)
		assert.Empty(t, ctx.SchoolName)
		assert.Empty(t, ctx.AcademicUnitID)
		assert.Empty(t, ctx.AcademicUnitName)
		assert.Len(t, ctx.Grants.Allow, 1)
	})
}

func TestConcurrentTokenGeneration(t *testing.T) {
	manager, _ := testManager(t, testIssuer)
	activeContext := &UserContext{
		RoleID:   "role-123",
		RoleName: "Teacher",
		Grants:   Grants{Allow: []string{"admin.users.read"}},
	}

	t.Run("genera tokens concurrentemente sin errores", func(t *testing.T) {
		const numTokens = 100
		results := make(chan bool, numTokens)

		for range numTokens {
			go func() {
				userID := uuid.New().String()
				_, _, err := manager.GenerateContextToken(
					userID,
					testEmail,
					activeContext,
					time.Hour,
				)
				results <- (err == nil)
			}()
		}

		// Esperar a que todos terminen
		successCount := 0
		for range numTokens {
			select {
			case success := <-results:
				if success {
					successCount++
				}
			case <-time.After(5 * time.Second):
				t.Fatal("Timeout esperando generación de tokens")
			}
		}

		assert.Equal(t, numTokens, successCount)
	})
}

func TestTokenSecurity(t *testing.T) {
	manager1, _ := testManager(t, "issuer-1")
	manager2, _ := testManager(t, "issuer-2")

	userID := uuid.New().String()
	activeContext := &UserContext{
		RoleID:   "role-123",
		RoleName: "Teacher",
		Grants:   Grants{Allow: []string{"admin.users.read"}},
	}

	t.Run("tokens con claves distintas no son intercambiables", func(t *testing.T) {
		token, _, err := manager1.GenerateContextToken(userID, testEmail, activeContext, time.Hour)
		require.NoError(t, err)

		_, err = manager2.ValidateToken(token)
		assert.Error(t, err)
	})
}
