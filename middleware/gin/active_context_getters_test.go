package gin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/gin-gonic/gin"
)

// Cubre las ramas de error de los getters de contexto activo (valor ausente o
// de tipo/valor inválido), el panic de MustActiveUnitID y ramas de rechazo de
// RequireActiveContext, complementando los casos felices ya probados en
// active_context_test.go.

func TestGetActiveUnitID_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	if _, err := GetActiveUnitID(c); !errors.Is(err, ErrNoActiveUnit) {
		t.Fatalf("want ErrNoActiveUnit, got %v", err)
	}
}

func TestGetActiveUnitID_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ContextKeyActiveUnitID, 123) // no es string
	if _, err := GetActiveUnitID(c); !errors.Is(err, ErrNoActiveUnit) {
		t.Fatalf("want ErrNoActiveUnit por tipo inválido, got %v", err)
	}
}

func TestGetActiveUnitID_EmptyString(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ContextKeyActiveUnitID, "")
	if _, err := GetActiveUnitID(c); !errors.Is(err, ErrNoActiveUnit) {
		t.Fatalf("want ErrNoActiveUnit por string vacío, got %v", err)
	}
}

func TestGetActiveUnitID_Set(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ContextKeyActiveUnitID, "unit-001")
	got, err := GetActiveUnitID(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "unit-001" {
		t.Fatalf("got %q, want unit-001", got)
	}
}

func TestGetActiveSchoolID_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ContextKeyActiveSchoolID, 42) // no es string
	if _, err := GetActiveSchoolID(c); !errors.Is(err, ErrNoActiveSchool) {
		t.Fatalf("want ErrNoActiveSchool por tipo inválido, got %v", err)
	}
}

func TestRequireActiveContext_RejectsNoClaims(t *testing.T) {
	r := newTestRouter()
	r.GET("/x",
		RequireActiveContext(),
		func(c *gin.Context) { c.Status(http.StatusOK) },
	)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusPreconditionRequired {
		t.Fatalf("want 428, got %d", rec.Code)
	}
}

func TestRequireActiveContext_RejectsMissingSchool(t *testing.T) {
	r := newTestRouter()
	r.GET("/x",
		withClaims(&auth.Claims{
			ActiveContext: &auth.UserContext{AcademicUnitID: "unit-001"}, // sin SchoolID
		}),
		RequireActiveContext(),
		func(c *gin.Context) { c.Status(http.StatusOK) },
	)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusPreconditionRequired {
		t.Fatalf("want 428, got %d", rec.Code)
	}
}

func TestMustActiveUnitID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ContextKeyActiveUnitID, "unit-001")
	if got := MustActiveUnitID(c); got != "unit-001" {
		t.Fatalf("got %q, want unit-001", got)
	}
}

func TestMustActiveUnitID_PanicsWhenMissing(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	_ = MustActiveUnitID(c)
}
