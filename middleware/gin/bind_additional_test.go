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
)

// TestBindJSON_CustomJSONTags prueba que BindJSON usa los tags json correctamente
func TestBindJSON_CustomJSONTags(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type customTagRequest struct {
		UserId    string `json:"user_id" binding:"required"`
		FirstName string `json:"first_name" binding:"required"`
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Falta user_id
	c.Request, _ = http.NewRequest("POST", "/test", bytes.NewBufferString(`{"first_name":"John"}`)) //nolint:errcheck
	c.Request.Header.Set("Content-Type", "application/json")

	var req customTagRequest
	err := BindJSON(c, &req)

	assert.Error(t, err)
	var appErr *sharedErrors.AppError
	assert.True(t, errors.As(err, &appErr))
	// Debe usar el tag json "user_id", no snake_case del struct field name
	assert.Contains(t, appErr.Fields, "user_id")
}

// TestToSnakeCase_EdgeCases prueba diferentes conversions de snake_case
func TestToSnakeCase_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"name", "name"},
		{"Name", "name"},
		{"UserID", "user_id"},
		{"HTTPResponse", "http_response"},
		{"getHTTPResponse", "get_http_response"},
		{"ID", "id"},
		{"HTTPSServer", "https_server"},
		{"aB", "a_b"},
		{"ABC", "abc"},
		{"ABc", "a_bc"},
		{"apiKeySecret", "api_key_secret"},
	}

	for _, tt := range tests {
		result := ToSnakeCase(tt.input)
		if result != tt.expected {
			t.Errorf("ToSnakeCase(%q): got %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

// TestBindJSON_NestedStructs prueba BindJSON con structs anidados
func TestBindJSON_NestedStructs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type addressStruct struct {
		Street string `json:"street" binding:"required"`
		City   string `json:"city" binding:"required"`
	}

	type userRequest struct {
		Name    string        `json:"name" binding:"required"`
		Address addressStruct `json:"address" binding:"required"`
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"John","address":{}}`)) //nolint:errcheck
	c.Request.Header.Set("Content-Type", "application/json")

	var req userRequest
	err := BindJSON(c, &req)

	// Debe fallar por los campos requeridos en el address anidado
	assert.Error(t, err)
}

// TestBindJSON_ValidComplexRequest prueba un request complejo válido
func TestBindJSON_ValidComplexRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type complexRequest struct {
		UserId    string `json:"user_id" binding:"required,uuid"`
		Email     string `json:"email" binding:"required,email"`
		Age       int    `json:"age" binding:"required,min=0,max=150"`
		FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/test", //nolint:errcheck
		bytes.NewBufferString(`{"user_id":"550e8400-e29b-41d4-a716-446655440000","email":"test@example.com","age":25,"first_name":"John"}`)) //nolint:errcheck
	c.Request.Header.Set("Content-Type", "application/json")

	var req complexRequest
	err := BindJSON(c, &req)

	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", req.UserId)
	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, 25, req.Age)
	assert.Equal(t, "John", req.FirstName)
}
