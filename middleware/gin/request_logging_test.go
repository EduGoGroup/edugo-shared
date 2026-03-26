package gin

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLoggingTestRouter(buf *bytes.Buffer) (*gin.Engine, *slog.Logger) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	testLogger := slog.New(handler)

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

	var entry map[string]interface{}
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
}

func TestRequestLogging_LogsWarnFor4xx(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusNotFound) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var entry map[string]interface{}
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

	var entry map[string]interface{}
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

func TestRequestLogging_IncludesUserIDPostAuth(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)

	// Simular middleware JWT que establece user_id
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, "user-abc")
		c.Next()
	})

	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "user-abc", entry[logger.FieldUserID])
}

func TestRequestLogging_IncludesRolePostAuth(t *testing.T) {
	buf := &bytes.Buffer{}
	r, _ := setupLoggingTestRouter(buf)

	// Simular middleware JWT que establece user_id y role
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, "user-abc")
		c.Set(ContextKeyRole, "teacher")
		c.Next()
	})

	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "teacher", entry["role"])
}

func TestGetLogger_DefaultWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	l := GetLogger(c)
	assert.Equal(t, slog.Default(), l)
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
		_ = c.Error(assert.AnError)
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	output := buf.String()
	assert.True(t, strings.Contains(output, logger.FieldError))
}
