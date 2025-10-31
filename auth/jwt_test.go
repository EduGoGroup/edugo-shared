package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

const (
	testSecretKey            = "test-secret-key-at-least-32-chars-long-for-security"
	testIssuer               = "edugo-test"
	testEmail                = "test@example.com"
	invalidTokenHere         = "invalid-token-here"
	notValidToken            = "not-valid-token"
	unauthorizedTokenExpired = "UNAUTHORIZED: token expired"
	unauthorizedInvalidToken = "UNAUTHORIZED: invalid token"
)

func TestNewJWTManager(t *testing.T) {
	t.Run("crea JWTManager correctamente", func(t *testing.T) {
		manager := NewJWTManager(testSecretKey, testIssuer)

		assert.NotNil(t, manager)
		assert.Equal(t, []byte(testSecretKey), manager.secretKey)
		assert.Equal(t, testIssuer, manager.issuer)
	})
}

func TestGenerateToken(t *testing.T) {
	manager := NewJWTManager(testSecretKey, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	role := enum.SystemRoleTeacher
	expiresIn := 24 * time.Hour

	t.Run("genera token válido exitosamente", func(t *testing.T) {
		token, err := manager.GenerateToken(userID, email, role, expiresIn)

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verificar que el token tiene 3 partes (header.payload.signature)
		assert.Equal(t, 3, len(splitToken(token)))
	})

	t.Run("genera tokens únicos", func(t *testing.T) {
		token1, err1 := manager.GenerateToken(userID, email, role, expiresIn)
		time.Sleep(10 * time.Millisecond) // Pequeña pausa
		token2, err2 := manager.GenerateToken(userID, email, role, expiresIn)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2, "Los tokens deben ser únicos debido al timestamp y JTI")
	})

	t.Run("genera token con todos los roles", func(t *testing.T) {
		roles := []enum.SystemRole{
			enum.SystemRoleAdmin,
			enum.SystemRoleTeacher,
			enum.SystemRoleStudent,
			enum.SystemRoleGuardian,
		}

		for _, testRole := range roles {
			token, err := manager.GenerateToken(userID, email, testRole, expiresIn)

			require.NoError(t, err, "Debe generar token para role: %s", testRole)
			assert.NotEmpty(t, token)

			// Validar que el token contiene el role correcto
			claims, err := manager.ValidateToken(token)
			require.NoError(t, err)
			assert.Equal(t, testRole, claims.Role)
		}
	})

	t.Run("genera token con claims correctos", func(t *testing.T) {
		token, err := manager.GenerateToken(userID, email, role, expiresIn)

		require.NoError(t, err)

		// Parsear token para verificar claims
		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, testIssuer, claims.Issuer)
		assert.Equal(t, userID, claims.Subject)
		assert.NotEmpty(t, claims.ID) // JTI debe estar presente

		// Verificar tiempos
		now := time.Now()
		assert.True(t, claims.IssuedAt.Time.Before(now.Add(1*time.Second)))
		assert.True(t, claims.ExpiresAt.Time.After(now))
		assert.True(t, claims.NotBefore.Time.Before(now.Add(1*time.Second)))
	})

	t.Run("genera token con expiración personalizada", func(t *testing.T) {
		shortExpiry := 1 * time.Hour
		token, err := manager.GenerateToken(userID, email, role, shortExpiry)

		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)

		expectedExpiry := time.Now().Add(shortExpiry)
		actualExpiry := claims.ExpiresAt.Time

		// Verificar que la expiración es aproximadamente correcta (±2 segundos)
		diff := actualExpiry.Sub(expectedExpiry)
		assert.True(t, diff < 2*time.Second && diff > -2*time.Second,
			"Expiración debe ser aproximadamente %v, got %v", expectedExpiry, actualExpiry)
	})
}

func TestValidateToken(t *testing.T) {
	manager := NewJWTManager(testSecretKey, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	role := enum.SystemRoleTeacher

	t.Run("valida token válido exitosamente", func(t *testing.T) {
		token, err := manager.GenerateToken(userID, email, role, 24*time.Hour)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)

		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("rechaza token vacío", func(t *testing.T) {
		claims, err := manager.ValidateToken("")

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("rechaza token malformado", func(t *testing.T) {
		invalidToken := invalidTokenHere

		claims, err := manager.ValidateToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("rechaza token con firma incorrecta", func(t *testing.T) {
		// Crear token con secret diferente
		wrongManager := NewJWTManager("wrong-secret-key", testIssuer)
		token, err := wrongManager.GenerateToken(userID, email, role, 24*time.Hour)
		require.NoError(t, err)

		// Intentar validar con manager original
		claims, err := manager.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("rechaza token expirado", func(t *testing.T) {
		// Generar token que ya expiró
		token, err := manager.GenerateToken(userID, email, role, -1*time.Hour)
		require.NoError(t, err)

		claims, err := manager.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		// El error puede ser "token expired" o "invalid token" dependiendo de la implementación
		assert.True(t, err.Error() == unauthorizedTokenExpired || err.Error() == unauthorizedInvalidToken,
			"Error debe indicar token expirado o inválido, got: %s", err.Error())
	})

	t.Run("rechaza token con método de firma incorrecto", func(t *testing.T) {
		// Crear token con método de firma diferente (RS256 en lugar de HS256)
		claims := Claims{
			UserID: userID,
			Email:  email,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    testIssuer,
			},
		}

		// Crear token con algoritmo None (inseguro)
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

		validatedClaims, err := manager.ValidateToken(tokenString)

		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
	})
}

func TestRefreshToken(t *testing.T) {
	manager := NewJWTManager(testSecretKey, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	role := enum.SystemRoleStudent

	t.Run("refresca token válido exitosamente", func(t *testing.T) {
		originalToken, err := manager.GenerateToken(userID, email, role, 24*time.Hour)
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond) // Pequeña pausa

		newToken, err := manager.RefreshToken(originalToken, 48*time.Hour)

		require.NoError(t, err)
		assert.NotEmpty(t, newToken)
		assert.NotEqual(t, originalToken, newToken, "Nuevo token debe ser diferente")

		// Validar nuevo token
		newClaims, err := manager.ValidateToken(newToken)
		require.NoError(t, err)
		assert.Equal(t, userID, newClaims.UserID)
		assert.Equal(t, email, newClaims.Email)
		assert.Equal(t, role, newClaims.Role)

		// Verificar que la nueva expiración es mayor
		originalClaims, _ := manager.ValidateToken(originalToken)
		assert.True(t, newClaims.ExpiresAt.After(originalClaims.ExpiresAt.Time),
			"Nuevo token debe tener expiración mayor")
	})

	t.Run("falla al refrescar token expirado", func(t *testing.T) {
		expiredToken, err := manager.GenerateToken(userID, email, role, -1*time.Hour)
		require.NoError(t, err)

		newToken, err := manager.RefreshToken(expiredToken, 24*time.Hour)

		assert.Error(t, err)
		assert.Empty(t, newToken)
		// El error puede ser "token expired" o "invalid token" dependiendo de la implementación
		assert.True(t, err.Error() == unauthorizedTokenExpired || err.Error() == unauthorizedInvalidToken,
			"Error debe indicar token expirado o inválido, got: %s", err.Error())
	})

	t.Run("falla al refrescar token inválido", func(t *testing.T) {
		invalidToken := invalidTokenHere

		newToken, err := manager.RefreshToken(invalidToken, 24*time.Hour)

		assert.Error(t, err)
		assert.Empty(t, newToken)
	})

	t.Run("maintains original claims when refreshing", func(t *testing.T) {
		originalToken, err := manager.GenerateToken(userID, email, role, 1*time.Hour)
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond) // Pausa para asegurar timestamp diferente

		refreshedToken, err := manager.RefreshToken(originalToken, 24*time.Hour)
		require.NoError(t, err)

		originalClaims, _ := manager.ValidateToken(originalToken)
		refreshedClaims, _ := manager.ValidateToken(refreshedToken)

		// Verificar que los datos del usuario se mantienen
		assert.Equal(t, originalClaims.UserID, refreshedClaims.UserID)
		assert.Equal(t, originalClaims.Email, refreshedClaims.Email)
		assert.Equal(t, originalClaims.Role, refreshedClaims.Role)

		// Pero los metadatos cambian
		assert.NotEqual(t, originalClaims.ID, refreshedClaims.ID, "JTI debe ser diferente")
		// IssuedAt debe ser igual o más reciente (permitir pequeñas diferencias de timing)
		assert.True(t, refreshedClaims.IssuedAt.Time.After(originalClaims.IssuedAt.Time) ||
			refreshedClaims.IssuedAt.Time.Equal(originalClaims.IssuedAt.Time),
			"IssuedAt debe ser igual o más reciente")
	})
}

func TestExtractUserID(t *testing.T) {
	manager := NewJWTManager(testSecretKey, testIssuer)
	userID := uuid.New().String()
	email := testEmail
	role := enum.SystemRoleAdmin

	t.Run("extrae userID de token válido", func(t *testing.T) {
		token, err := manager.GenerateToken(userID, email, role, 24*time.Hour)
		require.NoError(t, err)

		extractedUserID, err := ExtractUserID(token)

		require.NoError(t, err)
		assert.Equal(t, userID, extractedUserID)
	})

	t.Run("extrae userID de token expirado (sin validar)", func(t *testing.T) {
		// Esta función NO valida expiración, solo extrae el claim
		expiredToken, err := manager.GenerateToken(userID, email, role, -1*time.Hour)
		require.NoError(t, err)

		extractedUserID, err := ExtractUserID(expiredToken)

		require.NoError(t, err)
		assert.Equal(t, userID, extractedUserID,
			"ExtractUserID debe funcionar incluso con tokens expirados")
	})

	t.Run("extrae userID de token con firma incorrecta (sin validar)", func(t *testing.T) {
		// Esta función NO valida firma, solo extrae el claim
		wrongManager := NewJWTManager("wrong-secret", testIssuer)
		token, err := wrongManager.GenerateToken(userID, email, role, 24*time.Hour)
		require.NoError(t, err)

		extractedUserID, err := ExtractUserID(token)

		require.NoError(t, err)
		assert.Equal(t, userID, extractedUserID,
			"ExtractUserID debe funcionar incluso con firma incorrecta")
	})

	t.Run("falla con token malformado", func(t *testing.T) {
		invalidToken := notValidToken

		extractedUserID, err := ExtractUserID(invalidToken)

		assert.Error(t, err)
		assert.Empty(t, extractedUserID)
	})

	t.Run("falla con token vacío", func(t *testing.T) {
		extractedUserID, err := ExtractUserID("")

		assert.Error(t, err)
		assert.Empty(t, extractedUserID)
	})
}

func TestExtractRole(t *testing.T) {
	manager := NewJWTManager(testSecretKey, testIssuer)
	userID := uuid.New().String()
	email := testEmail

	t.Run("extrae role de token válido", func(t *testing.T) {
		roles := []enum.SystemRole{
			enum.SystemRoleAdmin,
			enum.SystemRoleTeacher,
			enum.SystemRoleStudent,
			enum.SystemRoleGuardian,
		}

		for _, expectedRole := range roles {
			token, err := manager.GenerateToken(userID, email, expectedRole, 24*time.Hour)
			require.NoError(t, err)

			extractedRole, err := ExtractRole(token)

			require.NoError(t, err)
			assert.Equal(t, expectedRole, extractedRole,
				"Debe extraer role correcto: %s", expectedRole)
		}
	})

	t.Run("extrae role de token expirado (sin validar)", func(t *testing.T) {
		role := enum.SystemRoleTeacher
		expiredToken, err := manager.GenerateToken(userID, email, role, -1*time.Hour)
		require.NoError(t, err)

		extractedRole, err := ExtractRole(expiredToken)

		require.NoError(t, err)
		assert.Equal(t, role, extractedRole,
			"ExtractRole debe funcionar incluso con tokens expirados")
	})

	t.Run("falla con token malformado", func(t *testing.T) {
		invalidToken := notValidToken

		extractedRole, err := ExtractRole(invalidToken)

		assert.Error(t, err)
		assert.Empty(t, extractedRole)
	})

	t.Run("falla con token vacío", func(t *testing.T) {
		extractedRole, err := ExtractRole("")

		assert.Error(t, err)
		assert.Empty(t, extractedRole)
	})
}

func TestConcurrentTokenGeneration(t *testing.T) {
	t.Run("genera tokens concurrentemente sin errores", func(t *testing.T) {
		manager := NewJWTManager(testSecretKey, testIssuer)
		concurrentUsers := 100
		results := make(chan string, concurrentUsers)
		errors := make(chan error, concurrentUsers)

		// Generar tokens concurrentemente
		for i := 0; i < concurrentUsers; i++ {
			go func(index int) {
				userID := uuid.New().String()
				email := "user" + string(rune(index)) + "@edugo.com"
				token, err := manager.GenerateToken(userID, email, enum.SystemRoleStudent, 24*time.Hour)
				if err != nil {
					errors <- err
				} else {
					results <- token
				}
			}(i)
		}

		// Recolectar resultados
		tokens := make(map[string]bool)
		for i := 0; i < concurrentUsers; i++ {
			select {
			case token := <-results:
				tokens[token] = true
			case err := <-errors:
				t.Fatalf("Error generando token: %v", err)
			}
		}

		// Verificar que todos los tokens son únicos
		assert.Equal(t, concurrentUsers, len(tokens),
			"Todos los tokens deben ser únicos incluso con generación concurrente")
	})
}

func TestTokenSecurity(t *testing.T) {
	t.Run("tokens con diferentes secrets no son intercambiables", func(t *testing.T) {
		manager1 := NewJWTManager("secret-key-1", "issuer1")
		manager2 := NewJWTManager("secret-key-2", "issuer2")

		userID := uuid.New().String()
		email := testEmail
		role := enum.SystemRoleTeacher

		token1, err := manager1.GenerateToken(userID, email, role, 24*time.Hour)
		require.NoError(t, err)

		// Intentar validar con manager diferente
		claims, err := manager2.ValidateToken(token1)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("tokens no pueden ser modificados sin invalidar firma", func(t *testing.T) {
		manager := NewJWTManager(testSecretKey, testIssuer)
		userID := uuid.New().String()
		email := testEmail
		role := enum.SystemRoleStudent

		token, err := manager.GenerateToken(userID, email, role, 24*time.Hour)
		require.NoError(t, err)

		// Intentar modificar el token (cambiar un carácter)
		modifiedToken := token[:len(token)-5] + "XXXXX"

		claims, err := manager.ValidateToken(modifiedToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

// Helper function para dividir token en partes
func splitToken(token string) []string {
	parts := []string{}
	current := ""
	for _, char := range token {
		if char == '.' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
