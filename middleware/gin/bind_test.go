package gin

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	sharedErrors "github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testBindRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"min=1,max=150"`
}

func TestBindJSON_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"John","email":"john@example.com","age":25}`)) //nolint:errcheck
	c.Request.Header.Set("Content-Type", "application/json")

	var req testBindRequest
	err := BindJSON(c, &req)

	assert.NoError(t, err)
	assert.Equal(t, "John", req.Name)
	assert.Equal(t, "john@example.com", req.Email)
	assert.Equal(t, 25, req.Age)
}

func TestBindJSON_MissingRequiredField(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/test", bytes.NewBufferString(`{"email":"john@example.com","age":25}`)) //nolint:errcheck
	c.Request.Header.Set("Content-Type", "application/json")

	var req testBindRequest
	err := BindJSON(c, &req)

	require.Error(t, err)
	var appErr *sharedErrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Contains(t, appErr.Fields, "name")
	assert.Equal(t, "field is required", appErr.Fields["name"])
}

func TestBindJSON_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"John","email":"not-an-email","age":25}`)) //nolint:errcheck
	c.Request.Header.Set("Content-Type", "application/json")

	var req testBindRequest
	err := BindJSON(c, &req)

	require.Error(t, err)
	var appErr *sharedErrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Contains(t, appErr.Fields, "email")
	assert.Equal(t, "invalid email format", appErr.Fields["email"])
}

func TestBindJSON_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/test", bytes.NewBufferString(`{invalid json`)) //nolint:errcheck
	c.Request.Header.Set("Content-Type", "application/json")

	var req testBindRequest
	err := BindJSON(c, &req)

	require.Error(t, err)
	var appErr *sharedErrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, "invalid request body", appErr.Message)
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"Name", "name"},
		{"FirstName", "first_name"},
		{"HTTPSServer", "https_server"},
		{"userID", "user_id"},
		{"APIKey", "api_key"},
		{"simpleTest", "simple_test"},
		{"JSONData", "json_data"},
		{"getHTTPResponse", "get_http_response"},
	}

	for _, tt := range tests {
		result := ToSnakeCase(tt.input)
		assert.Equal(t, tt.expected, result, "ToSnakeCase(%q)", tt.input)
	}
}

func TestValidationMessage_MinMaxDistinguishesKind(t *testing.T) {
	// Este test verifica que min/max distingue entre length y value.
	// No podemos crear FieldError directamente, pero podemos verificar
	// la logica via BindJSON con structs que tengan min en string vs int.

	type stringMinReq struct {
		Name string `json:"name" binding:"min=3"`
	}
	type intMinReq struct {
		Age int `json:"age" binding:"min=1"`
	}

	gin.SetMode(gin.TestMode)

	// String min -> "minimum length is 3"
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"ab"}`)) //nolint:errcheck
	c1.Request.Header.Set("Content-Type", "application/json")

	var sr stringMinReq
	err := BindJSON(c1, &sr)
	require.Error(t, err)
	var appErr1 *sharedErrors.AppError
	require.True(t, errors.As(err, &appErr1))
	assert.Equal(t, "minimum length is 3", appErr1.Fields["name"])

	// Int min -> "minimum value is 1"
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"age":0}`)) //nolint:errcheck
	c2.Request.Header.Set("Content-Type", "application/json")

	var ir intMinReq
	err = BindJSON(c2, &ir)
	require.Error(t, err)
	var appErr2 *sharedErrors.AppError
	require.True(t, errors.As(err, &appErr2))
	assert.Equal(t, "minimum value is 1", appErr2.Fields["age"])
}
