package gin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// capturingLogger captura los eventos registrados para verificarlos en los tests.
type capturingLogger struct {
	events []audit.AuditEvent
}

func (l *capturingLogger) Log(ctx context.Context, event audit.AuditEvent) error {
	l.events = append(l.events, event)
	return nil
}

func TestAuditMiddleware_SoloRegistraMetodosMutantes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	metodosMutantes := []string{"POST", "PUT", "PATCH", "DELETE"}
	metodosLectura := []string{"GET", "HEAD", "OPTIONS"}

	for _, method := range metodosMutantes {
		t.Run(method, func(t *testing.T) {
			logger := &capturingLogger{}
			router := gin.New()
			router.Use(AuditMiddleware(logger))
			router.Handle(method, "/api/v1/users", func(c *gin.Context) {
				c.JSON(200, gin.H{"ok": true})
			})

			req, err := http.NewRequest(method, "/api/v1/users", nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Len(t, logger.events, 1, "método %s debe registrar un evento", method)
		})
	}

	for _, method := range metodosLectura {
		t.Run(method, func(t *testing.T) {
			logger := &capturingLogger{}
			router := gin.New()
			router.Use(AuditMiddleware(logger))
			router.Handle(method, "/api/v1/users", func(c *gin.Context) {
				c.JSON(200, gin.H{"ok": true})
			})

			req, err := http.NewRequest(method, "/api/v1/users", nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Empty(t, logger.events, "método %s NO debe registrar evento", method)
		})
	}
}

func TestAuditMiddleware_AccionCorrecta(t *testing.T) {
	gin.SetMode(gin.TestMode)

	casos := []struct {
		metodo string
		accion string
	}{
		{"POST", "create"},
		{"PUT", "update"},
		{"PATCH", "update"},
		{"DELETE", "delete"},
	}

	for _, tc := range casos {
		t.Run(tc.metodo, func(t *testing.T) {
			logger := &capturingLogger{}
			router := gin.New()
			router.Use(AuditMiddleware(logger))
			router.Handle(tc.metodo, "/api/v1/roles", func(c *gin.Context) {
				c.JSON(200, gin.H{})
			})

			req, err := http.NewRequest(tc.metodo, "/api/v1/roles", nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Len(t, logger.events, 1)
			assert.Equal(t, tc.accion, logger.events[0].Action)
		})
	}
}

func TestAuditMiddleware_ExtraeRecursoDelPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	casos := []struct {
		path            string
		recursoEsperado string
		idEsperado      string
	}{
		{"/api/v1/roles/123", "role", "123"},
		{"/api/v1/memberships", "membership", ""},
		{"/api/v1/users/abc-def", "user", "abc-def"},
		{"/api/v1/schools", "school", ""},
		{"/api/v1/materials/xyz", "material", "xyz"},
	}

	for _, tc := range casos {
		t.Run(tc.path, func(t *testing.T) {
			logger := &capturingLogger{}
			router := gin.New()
			router.Use(AuditMiddleware(logger))
			router.POST(tc.path, func(c *gin.Context) {
				c.JSON(200, gin.H{})
			})

			req, err := http.NewRequest("POST", tc.path, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Len(t, logger.events, 1)
			assert.Equal(t, tc.recursoEsperado, logger.events[0].ResourceType)
			assert.Equal(t, tc.idEsperado, logger.events[0].ResourceID)
		})
	}
}

func TestExtractResourceFromPath(t *testing.T) {
	casos := []struct {
		path            string
		recursoEsperado string
		idEsperado      string
	}{
		{"/api/v1/roles/123", "role", "123"},
		{"/api/v1/memberships", "membership", ""},
		{"/api/v1/users", "user", ""},
		{"/api/v1/schools/abc", "school", "abc"},
		{"/api/v1/unknown-path", "unknown_path", ""},
		{"/sin/version", "unknown", ""},
		{"/api/v1/categories", "category", ""},
		{"/api/v1/addresses", "address", ""},
		{"/api/v1/processes", "process", ""},
		{"/api/v1/classes", "class", ""},
	}

	for _, tc := range casos {
		t.Run(tc.path, func(t *testing.T) {
			recurso, id := extractResourceFromPath(tc.path)
			assert.Equal(t, tc.recursoEsperado, recurso)
			assert.Equal(t, tc.idEsperado, id)
		})
	}
}

func TestSingularize(t *testing.T) {
	casos := []struct {
		entrada  string
		esperado string
	}{
		{"roles", "role"},
		{"users", "user"},
		{"memberships", "membership"},
		{"categories", "category"},
		{"addresses", "address"},
		{"processes", "process"},
		{"classes", "class"},
		{"statuses", "status"},
		{"materials", "material"},
		{"schools", "school"},
		{"resource", "resource"}, // sin "s" final → sin cambio
	}

	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			resultado := singularize(tc.entrada)
			assert.Equal(t, tc.esperado, resultado)
		})
	}
}
