package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestAPIVersionHeader_SetsBothHeaders verifica que el middleware adjunta
// ambos headers (versión y build) con los valores recibidos y que continúa
// la cadena de handlers (c.Next()) llegando al handler final.
func TestAPIVersionHeader_SetsBothHeaders(t *testing.T) {
	r := newTestRouter()

	nextCalled := false
	r.GET("/x",
		APIVersionHeader("1.2.3", "abc1234"),
		func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		},
	)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))

	if !nextCalled {
		t.Fatal("APIVersionHeader no continuó la cadena: el handler final no se ejecutó")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if got := rec.Header().Get(HeaderAPIVersion); got != "1.2.3" {
		t.Errorf("%s = %q, want %q", HeaderAPIVersion, got, "1.2.3")
	}
	if got := rec.Header().Get(HeaderAPIBuild); got != "abc1234" {
		t.Errorf("%s = %q, want %q", HeaderAPIBuild, got, "abc1234")
	}
}

// TestAPIVersionHeader_DevDefaults verifica el caso de desarrollo local donde
// ambos valores valen "dev" porque las vars no se inyectan vía ldflags.
func TestAPIVersionHeader_DevDefaults(t *testing.T) {
	r := newTestRouter()
	r.GET("/x",
		APIVersionHeader("dev", "dev"),
		func(c *gin.Context) { c.Status(http.StatusNoContent) },
	)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", rec.Code)
	}
	if got := rec.Header().Get(HeaderAPIVersion); got != "dev" {
		t.Errorf("%s = %q, want %q", HeaderAPIVersion, got, "dev")
	}
	if got := rec.Header().Get(HeaderAPIBuild); got != "dev" {
		t.Errorf("%s = %q, want %q", HeaderAPIBuild, got, "dev")
	}
}
