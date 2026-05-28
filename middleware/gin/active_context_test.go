package gin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/gin-gonic/gin"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// withClaims monta un middleware que inyecta los claims directamente en el
// gin.Context, simulando lo que hace el JWTAuth middleware tras validar el
// token.
func withClaims(claims *auth.Claims) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ContextKeyClaims, claims)
		c.Next()
	}
}

func TestRequireActiveSchool_AllowsValidContext(t *testing.T) {
	r := newTestRouter()
	r.GET("/x",
		withClaims(&auth.Claims{
			ActiveContext: &auth.UserContext{SchoolID: "sch-001", AcademicUnitID: "unit-001"},
		}),
		RequireActiveSchool(),
		func(c *gin.Context) {
			s := MustActiveSchoolID(c)
			u, _ := GetActiveUnitID(c)
			c.JSON(http.StatusOK, gin.H{"school": s, "unit": u})
		},
	)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestRequireActiveSchool_RejectsMissingSchool(t *testing.T) {
	r := newTestRouter()
	r.GET("/x",
		withClaims(&auth.Claims{
			ActiveContext: &auth.UserContext{}, // sin SchoolID
		}),
		RequireActiveSchool(),
		func(c *gin.Context) { c.Status(http.StatusOK) },
	)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusPreconditionRequired {
		t.Fatalf("want 428, got %d", rec.Code)
	}
}

func TestRequireActiveSchool_RejectsNoClaims(t *testing.T) {
	r := newTestRouter()
	r.GET("/x",
		RequireActiveSchool(),
		func(c *gin.Context) { c.Status(http.StatusOK) },
	)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusPreconditionRequired {
		t.Fatalf("want 428, got %d", rec.Code)
	}
}

func TestRequireActiveContext_RejectsMissingUnit(t *testing.T) {
	r := newTestRouter()
	r.GET("/x",
		withClaims(&auth.Claims{
			ActiveContext: &auth.UserContext{SchoolID: "sch-001"}, // sin AcademicUnitID
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

func TestRequireActiveContext_AllowsFullContext(t *testing.T) {
	r := newTestRouter()
	r.GET("/x",
		withClaims(&auth.Claims{
			ActiveContext: &auth.UserContext{SchoolID: "sch-001", AcademicUnitID: "unit-001"},
		}),
		RequireActiveContext(),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"school": MustActiveSchoolID(c),
				"unit":   MustActiveUnitID(c),
			})
		},
	)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestGetActiveSchoolID_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, err := GetActiveSchoolID(c)
	if !errors.Is(err, ErrNoActiveSchool) {
		t.Fatalf("want ErrNoActiveSchool, got %v", err)
	}
}

func TestMustActiveSchoolID_PanicsWhenMissing(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	_ = MustActiveSchoolID(c)
}
