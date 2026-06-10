package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	svcSecret   = "service-secret-key"
	svcIssuer   = "edugo-identity"
	svcAudience = "edugo-api-platform"
	dispatch    = "notifications.dispatch"
)

func newServiceManager() *auth.ServiceJWTManager {
	return auth.NewServiceJWTManager(svcSecret, svcIssuer, svcAudience)
}

// setupServiceRouter monta un router con el middleware M2M y un handler que
// devuelve el client_id inyectado en el contexto.
func setupServiceRouter(validator ServiceTokenValidator, requiredScope string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ServiceJWTAuthMiddleware(validator, requiredScope))
	router.POST("/api/v1/internal/notifications/dispatch", func(c *gin.Context) {
		clientID, err := GetClientID(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no client_id"})
			return
		}
		claims, err := GetServiceClaims(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no claims"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"client_id": clientID, "scopes": claims.Scopes})
	})
	return router
}

func doServiceRequest(router *gin.Engine, authHeader string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/internal/notifications/dispatch", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestServiceJWTAuthMiddleware_ValidToken(t *testing.T) {
	m := newServiceManager()
	token, _, err := m.GenerateServiceToken("edugo-worker", []string{dispatch}, time.Hour)
	require.NoError(t, err)

	router := setupServiceRouter(m, dispatch)
	w := doServiceRequest(router, "Bearer "+token)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Contains(t, w.Body.String(), "edugo-worker")
}

func TestServiceJWTAuthMiddleware_MissingHeader(t *testing.T) {
	router := setupServiceRouter(newServiceManager(), dispatch)
	w := doServiceRequest(router, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "MISSING_AUTH_HEADER")
}

func TestServiceJWTAuthMiddleware_MalformedHeader(t *testing.T) {
	router := setupServiceRouter(newServiceManager(), dispatch)
	w := doServiceRequest(router, "Token abc.def.ghi")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_AUTH_FORMAT")
}

func TestServiceJWTAuthMiddleware_InvalidToken(t *testing.T) {
	router := setupServiceRouter(newServiceManager(), dispatch)
	w := doServiceRequest(router, "Bearer not-a-jwt")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_SERVICE_TOKEN")
}

// TestServiceJWTAuthMiddleware_RejectsUserToken: un JWT de usuario (mismo
// secret, sin token_use="service") NO debe pasar el middleware de servicio.
func TestServiceJWTAuthMiddleware_RejectsUserToken(t *testing.T) {
	userManager := auth.NewJWTManager(svcSecret, svcIssuer)
	userToken, _, err := userManager.GenerateTokenWithContext(
		"user-123", "user@edugo.com",
		&auth.UserContext{RoleID: "r", RoleName: "Student"},
		time.Hour,
	)
	require.NoError(t, err)

	router := setupServiceRouter(newServiceManager(), dispatch)
	w := doServiceRequest(router, "Bearer "+userToken)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "token de usuario debe rechazarse en rutas de servicio")
	assert.Contains(t, w.Body.String(), "INVALID_SERVICE_TOKEN")
}

// TestServiceJWTAuthMiddleware_MissingScope: token de servicio válido pero sin
// el scope requerido → 403.
func TestServiceJWTAuthMiddleware_MissingScope(t *testing.T) {
	m := newServiceManager()
	token, _, err := m.GenerateServiceToken("edugo-worker", []string{"other.scope"}, time.Hour)
	require.NoError(t, err)

	router := setupServiceRouter(m, dispatch)
	w := doServiceRequest(router, "Bearer "+token)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "INSUFFICIENT_SCOPE")
}

// TestServiceJWTAuthMiddleware_NoScopeRequired: si requiredScope es vacío, no
// se exige scope (cualquier service token válido pasa).
func TestServiceJWTAuthMiddleware_NoScopeRequired(t *testing.T) {
	m := newServiceManager()
	token, _, err := m.GenerateServiceToken("edugo-api-learning", nil, time.Hour)
	require.NoError(t, err)

	router := setupServiceRouter(m, "")
	w := doServiceRequest(router, "Bearer "+token)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
}
