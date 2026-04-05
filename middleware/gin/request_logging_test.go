package gin

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLoggingTestRouter(buf *bytes.Buffer) (*gin.Engine, logger.Logger) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	testLogger := logger.NewSlogAdapter(slog.New(handler))

	r.Use(RequestLogging(testLogger))
	return r, testLogger
}

func TestRequestLogging_GeneratesRequestID(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// Debe tener header X-Request-ID en la respuesta
	requestID := w.Header().Get(HeaderRequestID)
	assert.NotEmpty(t, requestID)

	// Debe tener header X-Correlation-ID (por defecto igual al request ID)
	correlationID := w.Header().Get(HeaderCorrelationID)
	assert.Equal(t, requestID, correlationID)
}

func TestRequestLogging_PropagatesExistingRequestID(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(HeaderRequestID, "custom-id-123")
	r.ServeHTTP(w, req)

	assert.Equal(t, "custom-id-123", w.Header().Get(HeaderRequestID))
}

func TestRequestLogging_PropagatesCorrelationID(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(HeaderRequestID, "req-1")
	req.Header.Set(HeaderCorrelationID, "corr-1")
	r.ServeHTTP(w, req)

	assert.Equal(t, "corr-1", w.Header().Get(HeaderCorrelationID))
}

func TestRequestLogging_LogsFields(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	r.GET("/test/:id", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test/42", nil)
	r.ServeHTTP(w, req)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))

	assert.Equal(t, "request completed", entry["msg"])
	assert.Equal(t, "INFO", entry["level"])
	assert.NotEmpty(t, entry[logger.FieldRequestID])
	assert.NotEmpty(t, entry[logger.FieldCorrelationID])
	assert.Equal(t, "GET", entry[logger.FieldMethod])
	assert.Equal(t, "/test/:id", entry[logger.FieldPath])
	assert.Equal(t, float64(200), entry[logger.FieldStatusCode])

	// duration_ms debe ser numérico (int64 ms), no string
	durationVal, exists := entry[logger.FieldDuration]
	require.True(t, exists, "duration_ms debe estar presente")
	_, isFloat := durationVal.(float64)
	assert.True(t, isFloat, "duration_ms debe ser numérico, got %T", durationVal)

	// bytes debe estar presente con constante
	_, hasBytesField := entry[logger.FieldBytes]
	assert.True(t, hasBytesField, "bytes field debe estar presente")
}

func TestRequestLogging_FallbackPathFor404(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	// No registrar la ruta — cualquier request dará 404
	r.GET("/exists", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/unknown/path", nil)
	r.ServeHTTP(w, req)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))

	// Debe usar la URL real como fallback en vez de string vacío
	assert.Equal(t, "/unknown/path", entry[logger.FieldPath])
}

func TestRequestLogging_LogsWarnFor4xx(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusNotFound) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "WARN", entry["level"])
}

func TestRequestLogging_LogsErrorFor5xx(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusInternalServerError) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "ERROR", entry["level"])
}

func TestRequestLogging_InjectsLoggerInContext(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	r.GET("/test", func(c *gin.Context) {
		l := GetLogger(c)
		assert.NotNil(t, l)

		// El logger desde context.Context también debe funcionar
		ctxLogger := logger.FromContext(c.Request.Context())
		assert.NotNil(t, ctxLogger)

		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostAuthLogging_EnrichesLoggerWithUserInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	testLogger := logger.NewSlogAdapter(slog.New(handler))

	r.Use(RequestLogging(testLogger))

	// Simular JWT middleware que establece user_id, role y claims
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, "user-abc")
		c.Set(ContextKeyRole, "teacher")
		c.Set(ContextKeyClaims, &auth.Claims{
			UserID: "user-abc",
			ActiveContext: &auth.UserContext{
				RoleName: "teacher",
				SchoolID: "school-xyz",
			},
		})
		c.Next()
	})

	// PostAuthLogging enriquece el logger DESPUÉS del JWT
	r.Use(PostAuthLogging())

	r.GET("/test", func(c *gin.Context) {
		// Verificar que el logger en context tiene user_id
		ctxLogger := logger.FromContext(c.Request.Context())
		assert.NotNil(t, ctxLogger)

		// Nested overwrite: el logger post-auth debe ser distinto al base
		assert.NotEqual(t, testLogger, ctxLogger, "PostAuthLogging must overwrite the context logger")

		// Log dentro del handler — debe incluir user_id, role, school_id
		ctxLogger.Info("handler log")
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// Parsear todas las líneas de log
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.GreaterOrEqual(t, len(lines), 2, "debe haber al menos 2 log entries")

	// El log del handler debe tener user_id, role y school_id
	var handlerEntry map[string]any
	for _, line := range lines {
		var entry map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &entry))
		if entry["msg"] == "handler log" {
			handlerEntry = entry
			break
		}
	}
	require.NotNil(t, handlerEntry, "debe existir el log del handler")
	assert.Equal(t, "user-abc", handlerEntry[logger.FieldUserID])
	assert.Equal(t, "teacher", handlerEntry[logger.FieldRole])
	assert.Equal(t, "school-xyz", handlerEntry[logger.FieldSchoolID])
}

func TestPostAuthLogging_SummaryLogIncludesAuthFields(t *testing.T) {
	buf := &bytes.Buffer{}
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	testLogger := logger.NewSlogAdapter(slog.New(handler))

	r.Use(RequestLogging(testLogger))
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, "user-abc")
		c.Set(ContextKeyRole, "admin")
		c.Next()
	})
	r.Use(PostAuthLogging())

	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// El log "request completed" debe tener user_id y role
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	var summaryEntry map[string]any
	for _, line := range lines {
		var entry map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &entry))
		if entry["msg"] == "request completed" {
			summaryEntry = entry
			break
		}
	}
	require.NotNil(t, summaryEntry, "debe existir el log de request completed")
	assert.Equal(t, "user-abc", summaryEntry[logger.FieldUserID])
	assert.Equal(t, "admin", summaryEntry[logger.FieldRole])
}

func TestGetLogger_DefaultWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	l := GetLogger(c)
	assert.NotNil(t, l)
}

func TestGetRequestID_EmptyWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	assert.Empty(t, GetRequestID(c))
}

func TestGetRequestID_ReturnsStoredValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ContextKeyRequestID, "test-req-id")
	assert.Equal(t, "test-req-id", GetRequestID(c))
}

func TestRequestLogging_LogsGinErrors(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	r.GET("/test", func(c *gin.Context) {
		c.Error(assert.AnError) //nolint:errcheck // error deliberado para test
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	output := buf.String()
	assert.True(t, strings.Contains(output, logger.FieldError))
}
