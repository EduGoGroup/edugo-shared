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
	testSecretKey = "test-secret-key-at-least-32-chars-long-for-security"
	testIssuer    = "edugo-test"
	testEmail     = "test@example.com"
)

func TestNewJWTManager(t *testing.T) {
	t.Run("crea JWTManager correctamente", func(t *testing.T) {
		manager := NewJWTManager(testSecretKey, testIssuer)

		assert.NotNil(t, manager)
		assert.Equal(t, []byte(testSecretKey), manager.secretKey)
		assert.Equal(t, testIssuer, manager.issuer)
	})
}

func TestGenerateTokenWithContext(t *testing.T) {
	manager := NewJWTManager(testSecretKey, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	expiresIn := 24 * time.Hour

	activeContext := &UserContext{
		RoleID:      "role-123",
		RoleName:    "Teacher",
		SchoolID:    "school-456",
		SchoolName:  "Test School",
		Permissions: []string{"users:read", "materials:create", "assessments:grade"},
	}

	t.Run("genera token con contexto válido exitosamente", func(t *testing.T) {
		token, expiresAt, err := manager.GenerateTokenWithContext(userID, email, activeContext, expiresIn)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.False(t, expiresAt.IsZero())
		assert.True(t, expiresAt.After(time.Now()))
	})

	t.Run("token contiene claims con ActiveContext correctos", func(t *testing.T) {
		token, _, err := manager.GenerateTokenWithContext(userID, email, activeContext, expiresIn)
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
		assert.Equal(t, activeContext.Permissions, claims.ActiveContext.Permissions)
	})

	t.Run("rechaza userID vacío", func(t *testing.T) {
		_, _, err := manager.GenerateTokenWithContext("", email, activeContext, expiresIn)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "userID no puede estar vacío")
	})

	t.Run("rechaza email vacío", func(t *testing.T) {
		_, _, err := manager.GenerateTokenWithContext(userID, "", activeContext, expiresIn)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "email no puede estar vacío")
	})

	t.Run("rechaza activeContext nil", func(t *testing.T) {
		_, _, err := manager.GenerateTokenWithContext(userID, email, nil, expiresIn)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "activeContext no puede ser nil")
	})

	t.Run("rechaza expiresIn menor a 1 minuto", func(t *testing.T) {
		_, _, err := manager.GenerateTokenWithContext(userID, email, activeContext, 30*time.Second)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "expiresIn debe ser mayor a 1 minuto")
	})

	t.Run("acepta expiresIn de exactamente 1 minuto", func(t *testing.T) {
		token, _, err := manager.GenerateTokenWithContext(userID, email, activeContext, time.Minute)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("genera tokens únicos con mismo contexto", func(t *testing.T) {
		token1, _, err1 := manager.GenerateTokenWithContext(userID, email, activeContext, expiresIn)
		time.Sleep(10 * time.Millisecond)
		token2, _, err2 := manager.GenerateTokenWithContext(userID, email, activeContext, expiresIn)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2, "Los tokens deben ser únicos debido al timestamp y JTI")
	})

	t.Run("token sin SchoolID (context mínimo)", func(t *testing.T) {
		minimalContext := &UserContext{
			RoleID:      "role-789",
			RoleName:    "Admin",
			Permissions: []string{"users:create", "schools:manage"},
		}

		token, _, err := manager.GenerateTokenWithContext(userID, email, minimalContext, expiresIn)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, minimalContext.RoleID, claims.ActiveContext.RoleID)
		assert.Equal(t, minimalContext.RoleName, claims.ActiveContext.RoleName)
		assert.Empty(t, claims.ActiveContext.SchoolID)
		assert.Empty(t, claims.ActiveContext.SchoolName)
		assert.Equal(t, minimalContext.Permissions, claims.ActiveContext.Permissions)
	})

	t.Run("token con unidad académica", func(t *testing.T) {
		contextWithUnit := &UserContext{
			RoleID:           "role-abc",
			RoleName:         "Unit Coordinator",
			SchoolID:         "school-def",
			SchoolName:       "Test School",
			AcademicUnitID:   "unit-ghi",
			AcademicUnitName: "Computer Science",
			Permissions:      []string{"units:manage", "materials:create"},
		}

		token, _, err := manager.GenerateTokenWithContext(userID, email, contextWithUnit, expiresIn)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, contextWithUnit.AcademicUnitID, claims.ActiveContext.AcademicUnitID)
		assert.Equal(t, contextWithUnit.AcademicUnitName, claims.ActiveContext.AcademicUnitName)
	})

	t.Run("tiempo de expiración es correcto", func(t *testing.T) {
		expiresIn := 2 * time.Hour
		_, expiresAt, err := manager.GenerateTokenWithContext(userID, email, activeContext, expiresIn)

		require.NoError(t, err)

		expectedExpiration := time.Now().Add(expiresIn)
		// Tolerancia de 5 segundos
		assert.WithinDuration(t, expectedExpiration, expiresAt, 5*time.Second)
	})

	t.Run("token con lista vacía de permisos", func(t *testing.T) {
		contextNoPerms := &UserContext{
			RoleID:      "role-no-perms",
			RoleName:    "Guest",
			Permissions: []string{},
		}

		token, _, err := manager.GenerateTokenWithContext(userID, email, contextNoPerms, expiresIn)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		assert.NotNil(t, claims.ActiveContext.Permissions)
		assert.Empty(t, claims.ActiveContext.Permissions)
	})
}

func TestValidateToken(t *testing.T) {
	manager := NewJWTManager(testSecretKey, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	activeContext := &UserContext{
		RoleID:      "role-123",
		RoleName:    "Teacher",
		Permissions: []string{"users:read"},
	}

	t.Run("valida token válido exitosamente", func(t *testing.T) {
		token, _, err := manager.GenerateTokenWithContext(userID, email, activeContext, 24*time.Hour)
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
		wrongManager := NewJWTManager("wrong-secret-key", testIssuer)
		token, _, err := wrongManager.GenerateTokenWithContext(userID, email, activeContext, 24*time.Hour)
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
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		expiredTokenString, err := expiredToken.SignedString([]byte(testSecretKey))
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

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testSecretKey))
		require.NoError(t, err)

		_, err = manager.ValidateToken(tokenString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing active context")
	})
}

func TestExtractUserID(t *testing.T) {
	manager := NewJWTManager(testSecretKey, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	activeContext := &UserContext{
		RoleID:      "role-123",
		RoleName:    "Teacher",
		Permissions: []string{"users:read"},
	}

	t.Run("extrae userID de token válido", func(t *testing.T) {
		token, _, err := manager.GenerateTokenWithContext(userID, email, activeContext, 24*time.Hour)
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
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		expiredTokenString, err := expiredToken.SignedString([]byte(testSecretKey))
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
			Permissions:      []string{"users:read", "materials:create"},
		}

		assert.NotNil(t, ctx)
		assert.Equal(t, "role-123", ctx.RoleID)
		assert.Equal(t, "Teacher", ctx.RoleName)
		assert.Equal(t, "school-456", ctx.SchoolID)
		assert.Equal(t, "Test School", ctx.SchoolName)
		assert.Equal(t, "unit-789", ctx.AcademicUnitID)
		assert.Equal(t, "Mathematics", ctx.AcademicUnitName)
		assert.Len(t, ctx.Permissions, 2)
	})

	t.Run("UserContext solo con campos requeridos", func(t *testing.T) {
		ctx := &UserContext{
			RoleID:      "role-abc",
			RoleName:    "Admin",
			Permissions: []string{"users:create"},
		}

		assert.Equal(t, "role-abc", ctx.RoleID)
		assert.Equal(t, "Admin", ctx.RoleName)
		assert.Empty(t, ctx.SchoolID)
		assert.Empty(t, ctx.SchoolName)
		assert.Empty(t, ctx.AcademicUnitID)
		assert.Empty(t, ctx.AcademicUnitName)
		assert.Len(t, ctx.Permissions, 1)
	})
}

func TestConcurrentTokenGeneration(t *testing.T) {
	manager := NewJWTManager(testSecretKey, testIssuer)
	activeContext := &UserContext{
		RoleID:      "role-123",
		RoleName:    "Teacher",
		Permissions: []string{"users:read"},
	}

	t.Run("genera tokens concurrentemente sin errores", func(t *testing.T) {
		const numTokens = 100
		results := make(chan bool, numTokens)

		for i := 0; i < numTokens; i++ {
			go func() {
				userID := uuid.New().String()
				_, _, err := manager.GenerateTokenWithContext(
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
		for i := 0; i < numTokens; i++ {
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
	manager1 := NewJWTManager("secret-key-1", "issuer-1")
	manager2 := NewJWTManager("secret-key-2", "issuer-2")

	userID := uuid.New().String()
	activeContext := &UserContext{
		RoleID:      "role-123",
		RoleName:    "Teacher",
		Permissions: []string{"users:read"},
	}

	t.Run("tokens con diferentes secrets no son intercambiables", func(t *testing.T) {
		token, _, err := manager1.GenerateTokenWithContext(userID, testEmail, activeContext, time.Hour)
		require.NoError(t, err)

		_, err = manager2.ValidateToken(token)
		assert.Error(t, err)
	})
}
