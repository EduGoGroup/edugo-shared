package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear JWTManager para tests
	jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")

	// Generar token válido
	token, err := jwtManager.GenerateToken("user-123", "test@edugo.com", "student", time.Hour)
	if err != nil {
		t.Fatalf("Error al generar token: %v", err)
	}

	// Setup router con middleware
	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtManager))
	router.GET("/test", func(c *gin.Context) {
		userID, errUserID := GetUserID(c)
		email, errEmail := GetEmail(c)
		role, errRole := GetRole(c)

		// En test real estos errores no deberían ocurrir con token válido
		if errUserID != nil || errEmail != nil || errRole != nil {
			c.JSON(500, gin.H{"error": "failed to get claims"})
			return
		}

		c.JSON(200, gin.H{
			"user_id": userID,
			"email":   email,
			"role":    role,
		})
	})

	// Request con token válido
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verificar respuesta
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verificar que los claims están en la respuesta
	body := w.Body.String()
	if !containsString(body, "user-123") {
		t.Error("Response debe contener user_id")
	}
	if !containsString(body, "test@edugo.com") {
		t.Error("Response debe contener email")
	}
	if !containsString(body, "student") {
		t.Error("Response debe contener role")
	}
}

func TestJWTAuthMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", "test-issuer")

	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtManager))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	// No agregar header Authorization

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Debe retornar 401
	if w.Code != 401 {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	// Verificar mensaje de error
	if !containsString(w.Body.String(), "authorization header required") {
		t.Errorf("Error message incorrecto: %s", w.Body.String())
	}
}

func TestJWTAuthMiddleware_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", "test-issuer")

	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtManager))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	testCases := []struct {
		name   string
		header string
	}{
		{"Sin Bearer", "token123"},
		{"Bearer sin espacio", "Bearertoken123"},
		{"Mayúsculas incorrectas", "bearer token123"},
		{"Prefijo incorrecto", "Token token123"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tc.header)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != 401 {
				t.Errorf("Expected 401, got %d for header: %s", w.Code, tc.header)
			}

			if !containsString(w.Body.String(), "invalid authorization header format") {
				t.Errorf("Expected format error, got: %s", w.Body.String())
			}
		})
	}
}

func TestJWTAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", "test-issuer")

	// Generar token expirado (TTL negativo)
	token, err := jwtManager.GenerateToken("user123", "test@test.com", "student", -time.Hour)
	if err != nil {
		t.Fatalf("Error al generar token: %v", err)
	}

	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtManager))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Expected 401 for expired token, got %d", w.Code)
	}

	if !containsString(w.Body.String(), "invalid or expired token") {
		t.Errorf("Expected expired error, got: %s", w.Body.String())
	}
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", "test-issuer")

	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtManager))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	testCases := []struct {
		name  string
		token string
	}{
		{"Token malformado", "invalid.token.here"},
		{"Token vacío", ""},
		{"Solo puntos", "..."},
		{"Random string", "abcdefghijklmnop"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tc.token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != 401 {
				t.Errorf("Expected 401, got %d for token: %s", w.Code, tc.token)
			}
		})
	}
}

func TestJWTAuthMiddleware_WrongSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear token con un secret
	jwtManager1 := auth.NewJWTManager("secret-1", "issuer")
	token, err := jwtManager1.GenerateToken("user123", "test@test.com", "student", time.Hour)
	require.NoError(t, err)

	// Intentar validar con otro secret
	jwtManager2 := auth.NewJWTManager("secret-2", "issuer")

	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtManager2))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Debe rechazar (secret diferente)
	if w.Code != 401 {
		t.Errorf("Token con secret diferente debe ser rechazado, got %d", w.Code)
	}
}

func TestJWTAuthMiddleware_AbortChain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", "test-issuer")

	handlerCalled := false

	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtManager))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(200, gin.H{"ok": true})
	})

	// Request sin token
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Handler NO debe ser llamado
	if handlerCalled {
		t.Error("Handler no debe ser llamado cuando falla autenticación")
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && anySubstring(s, substr))
}

func anySubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
