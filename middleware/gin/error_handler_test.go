package gin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testLogger struct {
	messages []string
}

func (l *testLogger) Debug(msg string, args ...any) { l.messages = append(l.messages, msg) }
func (l *testLogger) Info(msg string, args ...any)  { l.messages = append(l.messages, msg) }
func (l *testLogger) Warn(msg string, args ...any)  { l.messages = append(l.messages, msg) }
func (l *testLogger) Error(msg string, args ...any) { l.messages = append(l.messages, msg) }
func (l *testLogger) Fatal(msg string, args ...any) { l.messages = append(l.messages, msg) }
func (l *testLogger) With(_ ...any) logger.Logger   { return l }
func (l *testLogger) Sync() error                   { return nil }

func TestErrorHandler_PanicRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := &testLogger{}

	r := gin.New()
	r.Use(ErrorHandler(log))
	r.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "internal server error", resp.Error)
	assert.Equal(t, "INTERNAL_ERROR", resp.Code)
	assert.Contains(t, log.messages, "panic recovered")
}

func TestErrorHandler_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := &testLogger{}

	r := gin.New()
	r.Use(ErrorHandler(log))
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.NewNotFoundError("user"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "user not found", resp.Error)
	assert.Equal(t, "NOT_FOUND", resp.Code)
}

func TestErrorHandler_AppErrorWithFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := &testLogger{}

	r := gin.New()
	r.Use(ErrorHandler(log))
	r.GET("/test", func(c *gin.Context) {
		appErr := errors.NewValidationErrorWithFields("validation failed", map[string]string{
			"name":  "field is required",
			"email": "invalid email format",
		})
		_ = c.Error(appErr)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "validation failed", resp.Error)
	assert.Equal(t, "VALIDATION_ERROR", resp.Code)
	assert.Equal(t, "field is required", resp.Details["name"])
	assert.Equal(t, "invalid email format", resp.Details["email"])
}

func TestErrorHandler_UnexpectedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := &testLogger{}

	r := gin.New()
	r.Use(ErrorHandler(log))
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(fmt.Errorf("something went wrong"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "internal server error", resp.Error)
	assert.Equal(t, "INTERNAL_ERROR", resp.Code)
}

func TestErrorHandler_NoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := &testLogger{}

	r := gin.New()
	r.Use(ErrorHandler(log))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, log.messages)
}

func TestHandleError_DirectCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	HandleError(c, errors.NewValidationError("bad input"))

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "bad input", resp.Error)
}
