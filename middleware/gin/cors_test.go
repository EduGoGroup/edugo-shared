package gin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCORSMiddleware_WildcardInDevelopment(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "*",
		AllowedMethods: "GET,POST",
		AllowedHeaders: "Origin,Content-Type",
	}

	handler := CORSMiddleware(cfg, "development")

	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected wildcard origin, got %q", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORSMiddleware_WildcardInLocal(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "*",
		AllowedMethods: "GET,POST",
		AllowedHeaders: "Origin,Content-Type",
	}

	handler := CORSMiddleware(cfg, "local")

	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected wildcard origin in local env, got %q", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORSMiddleware_ExplicitOriginsEmptyEnv(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "https://app.edugo.com",
		AllowedMethods: "GET,POST",
		AllowedHeaders: "Origin,Content-Type",
	}

	handler := CORSMiddleware(cfg, "")

	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "https://app.edugo.com")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://app.edugo.com" {
		t.Errorf("expected allowed origin with empty env, got %q", got)
	}
}

func TestCORSMiddleware_ExplicitOriginsInProduction(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "https://app.edugo.com,https://admin.edugo.com",
		AllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowedHeaders: "Origin,Content-Type,Authorization",
	}

	handler := CORSMiddleware(cfg, "production")

	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	// Allowed origin
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "https://app.edugo.com")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://app.edugo.com" {
		t.Errorf("expected allowed origin, got %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("expected credentials header, got %q", got)
	}

	// Disallowed origin
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	req2.Header.Set("Origin", "https://evil.com")
	r.ServeHTTP(w2, req2)

	if got := w2.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no origin header for disallowed origin, got %q", got)
	}
}

func TestCORSMiddleware_PreflightResponse(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "https://app.edugo.com",
		AllowedMethods: "GET,POST,PUT",
		AllowedHeaders: "Origin,Content-Type,Authorization",
	}

	handler := CORSMiddleware(cfg, "production")
	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "https://app.edugo.com")
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("expected 204 for preflight, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got != "GET,POST,PUT" {
		t.Errorf("expected methods header, got %q", got)
	}
	if got := w.Header().Get("Access-Control-Max-Age"); got != "86400" {
		t.Errorf("expected max-age header, got %q", got)
	}
}

func TestCORSMiddleware_ExposeHeadersOnNormalResponse(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "https://app.edugo.com",
		AllowedMethods: "GET,POST",
		AllowedHeaders: "Origin,Content-Type",
	}

	handler := CORSMiddleware(cfg, "production")
	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "https://app.edugo.com")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Expose-Headers"); got != "Content-Length,ETag,X-Request-ID,X-Correlation-ID" {
		t.Errorf("expected Expose-Headers on normal response, got %q", got)
	}
}

func TestCORSMiddleware_PreflightDisallowedOrigin(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "https://app.edugo.com",
		AllowedMethods: "GET,POST,PUT",
		AllowedHeaders: "Origin,Content-Type,Authorization",
	}

	handler := CORSMiddleware(cfg, "production")
	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "https://evil.com")
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("expected 204 for preflight, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got != "" {
		t.Errorf("expected no Allow-Methods for disallowed origin, got %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no Allow-Origin for disallowed origin, got %q", got)
	}
}

func TestCORSMiddleware_VaryOriginHeader(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "https://app.edugo.com",
		AllowedMethods: "GET,POST",
		AllowedHeaders: "Origin,Content-Type",
	}

	handler := CORSMiddleware(cfg, "production")
	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "https://app.edugo.com")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Vary"); got != "Origin" {
		t.Errorf("expected Vary: Origin for explicit origins, got %q", got)
	}
}

func TestCORSMiddleware_NoOriginHeader(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "https://app.edugo.com",
		AllowedMethods: "GET,POST",
		AllowedHeaders: "Origin,Content-Type",
	}

	handler := CORSMiddleware(cfg, "production")
	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	// No Origin header
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no CORS headers without Origin, got %q", got)
	}
}

func TestParseCSV(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"*", []string{"*"}},
		{"a,b,c", []string{"a", "b", "c"}},
		{" a , b , c ", []string{"a", "b", "c"}},
		{"a,,b", []string{"a", "b"}},
	}

	for _, tt := range tests {
		result := parseCSV(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("parseCSV(%q): got %v, expected %v", tt.input, result, tt.expected)
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("parseCSV(%q)[%d]: got %q, expected %q", tt.input, i, result[i], tt.expected[i])
			}
		}
	}
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		origin   string
		allowed  []string
		expected bool
	}{
		{"", []string{"*"}, false},
		{"http://a.com", []string{"*"}, true},
		{"http://a.com", []string{"http://a.com"}, true},
		{"http://a.com", []string{"http://b.com"}, false},
		{"http://a.com", []string{"http://b.com", "http://a.com"}, true},
	}

	for _, tt := range tests {
		result := isOriginAllowed(tt.origin, tt.allowed, true)
		if result != tt.expected {
			t.Errorf("isOriginAllowed(%q, %v): got %v, expected %v", tt.origin, tt.allowed, result, tt.expected)
		}
	}
}

func TestAppendVaryHeader(t *testing.T) {
	tests := []struct {
		name         string
		existingVary string
		appendValue  string
		expected     string
	}{
		{"empty", "", "Origin", "Origin"},
		{"single_different", "Accept-Encoding", "Origin", "Accept-Encoding,Origin"},
		{"already_has_origin", "Origin", "Origin", "Origin"},
		{"multiple_values", "Accept-Encoding,Content-Type", "Origin", "Accept-Encoding,Content-Type,Origin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.existingVary != "" {
				c.Writer.Header().Set("Vary", tt.existingVary)
			}
			appendVaryHeader(c, tt.appendValue)
			result := w.Header().Get("Vary")
			if result != tt.expected {
				t.Errorf("appendVaryHeader: got %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestIsOriginAllowed_WithAllowsWildcard(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowed        []string
		allowsWildcard bool
		expected       bool
	}{
		{"wildcard_allowed", "http://any.com", []string{"*"}, true, true},
		{"wildcard_not_allowed_in_prod", "http://any.com", []string{"*"}, false, false},
		{"explicit_origin_allowed", "http://a.com", []string{"http://a.com"}, false, true},
		{"explicit_origin_not_allowed", "http://a.com", []string{"http://b.com"}, false, false},
		{"mix_wildcard_and_explicit_dev", "http://any.com", []string{"*", "http://a.com"}, true, true},
		{"mix_wildcard_and_explicit_prod_wildcard", "http://x.com", []string{"*", "http://a.com"}, false, false},
		{"mix_wildcard_and_explicit_prod_explicit", "http://a.com", []string{"*", "http://a.com"}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOriginAllowed(tt.origin, tt.allowed, tt.allowsWildcard)
			if result != tt.expected {
				t.Errorf("isOriginAllowed(%q, %v, %v): got %v, expected %v",
					tt.origin, tt.allowed, tt.allowsWildcard, result, tt.expected)
			}
		})
	}
}

func TestCORSMiddleware_WildcardNotAllowedInProduction(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "*",
		AllowedMethods: "GET,POST",
		AllowedHeaders: "Origin,Content-Type",
	}

	// Ambiente vacío se trata como non-development, par bloquear wildcard
	handler := CORSMiddleware(cfg, "")

	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	// Con wildcard NO permitido en non-dev, no debe reflejar wildcard
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("expected no wildcard in non-development, got %q", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORSMiddleware_ReflectOriginInProduction(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "https://app.edugo.com,https://admin.edugo.com",
		AllowedMethods: "GET,POST",
		AllowedHeaders: "Origin,Content-Type",
	}

	handler := CORSMiddleware(cfg, "production")

	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "https://app.edugo.com")
	r.ServeHTTP(w, req)

	// En prod, refleja el origin específico
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://app.edugo.com" {
		t.Errorf("expected reflected origin, got %q", got)
	}

	// Debe tener Vary: Origin
	if got := w.Header().Get("Vary"); !strings.Contains(got, "Origin") {
		t.Errorf("expected Vary header with Origin, got %q", got)
	}
}

func TestCORSMiddleware_OptionsRequestAllowsWildcardInDev(t *testing.T) {
	cfg := CORSConfig{
		AllowedOrigins: "*",
		AllowedMethods: "GET,POST,PUT",
		AllowedHeaders: "Origin,Content-Type,Authorization",
	}

	handler := CORSMiddleware(cfg, "development")

	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil) //nolint:errcheck
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	// En dev con wildcard, debe permitir el preflight
	if w.Code != 204 {
		t.Errorf("expected 204 for wildcard preflight in dev, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected wildcard origin in dev preflight, got %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got != "GET,POST,PUT" {
		t.Errorf("expected methods in response, got %q", got)
	}
}
