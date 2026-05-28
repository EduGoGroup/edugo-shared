package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupGinContext(queryString string) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, err := http.NewRequest("GET", "/?"+queryString, nil)
	if err != nil {
		panic(err)
	}
	c.Request = req
	return c
}

func TestParseListFilters_EmptyRequest(t *testing.T) {
	c := setupGinContext("")
	filters, err := ParseListFilters(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters.IsActive != nil {
		t.Error("empty request should leave IsActive nil (show all)")
	}
	if filters.Limit != 50 {
		t.Errorf("empty request should default Limit to 50, got %d", filters.Limit)
	}
	if filters.Page != 0 || filters.Search != "" {
		t.Error("empty request should have zero-value Page and Search")
	}
}

func TestParseListFilters_LimitCappedAt200(t *testing.T) {
	c := setupGinContext("limit=500")
	filters, err := ParseListFilters(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters.Limit != 200 {
		t.Errorf("limit should be capped at 200, got %d", filters.Limit)
	}
}

func TestParseListFilters_ExplicitIsActiveFalse(t *testing.T) {
	c := setupGinContext("is_active=false")
	filters, err := ParseListFilters(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters.IsActive == nil || *filters.IsActive {
		t.Error("explicit is_active=false should be preserved")
	}
}

func TestParseListFilters_AllParams(t *testing.T) {
	c := setupGinContext("is_active=true&page=2&limit=25&search=john&search_fields=name,email")
	filters, err := ParseListFilters(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters.IsActive == nil || !*filters.IsActive {
		t.Error("is_active should be true")
	}
	if filters.Page != 2 {
		t.Errorf("page = %d, want 2", filters.Page)
	}
	if filters.Limit != 25 {
		t.Errorf("limit = %d, want 25", filters.Limit)
	}
	if filters.Search != "john" {
		t.Errorf("search = %q, want john", filters.Search)
	}
	if len(filters.SearchFields) != 2 || filters.SearchFields[0] != "name" || filters.SearchFields[1] != "email" {
		t.Errorf("search_fields = %v, want [name email]", filters.SearchFields)
	}
}

func TestParseListFilters_InvalidIsActive(t *testing.T) {
	c := setupGinContext("is_active=notbool")
	_, err := ParseListFilters(c)
	if err == nil {
		t.Error("expected error for invalid is_active")
	}
}

func TestParseListFilters_InvalidPage(t *testing.T) {
	c := setupGinContext("page=abc")
	_, err := ParseListFilters(c)
	if err == nil {
		t.Error("expected error for invalid page")
	}
}

func TestParseListFilters_NegativePage(t *testing.T) {
	c := setupGinContext("page=-1")
	_, err := ParseListFilters(c)
	if err == nil {
		t.Error("expected error for negative page")
	}
}

func TestParseListFilters_InvalidLimit(t *testing.T) {
	c := setupGinContext("limit=abc")
	_, err := ParseListFilters(c)
	if err == nil {
		t.Error("expected error for invalid limit")
	}
}

func TestParseListFilters_NegativeLimit(t *testing.T) {
	c := setupGinContext("limit=-5")
	_, err := ParseListFilters(c)
	if err == nil {
		t.Error("expected error for negative limit")
	}
}

func TestParseListFilters_SearchWithoutFields(t *testing.T) {
	c := setupGinContext("search=test")
	filters, err := ParseListFilters(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters.Search != "test" {
		t.Errorf("search = %q, want test", filters.Search)
	}
	if len(filters.SearchFields) != 0 {
		t.Errorf("search_fields should be empty, got %v", filters.SearchFields)
	}
}

func TestParseListFilters_FieldsTrimmed(t *testing.T) {
	c := setupGinContext("search=test&search_fields=+name+,+email+")
	filters, err := ParseListFilters(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters.SearchFields) != 2 || filters.SearchFields[0] != "name" || filters.SearchFields[1] != "email" {
		t.Errorf("search_fields = %v, want [name email] (trimmed)", filters.SearchFields)
	}
}

func TestParseListFilters_ExtraFieldsSingle(t *testing.T) {
	c := setupGinContext("status=active")
	filters, err := ParseListFilters(c, "status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters.FieldFilters == nil || len(filters.FieldFilters["status"]) != 1 || filters.FieldFilters["status"][0] != "active" {
		t.Errorf("FieldFilters[status] = %v, want [active]", filters.FieldFilters["status"])
	}
}

func TestParseListFilters_ExtraFieldsMulti(t *testing.T) {
	c := setupGinContext("status=active,pending")
	filters, err := ParseListFilters(c, "status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters.FieldFilters == nil || len(filters.FieldFilters["status"]) != 2 {
		t.Errorf("FieldFilters[status] = %v, want [active pending]", filters.FieldFilters["status"])
	}
}

func TestParseListFilters_ExtraFieldsAbsent(t *testing.T) {
	c := setupGinContext("")
	filters, err := ParseListFilters(c, "status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filters.FieldFilters != nil {
		t.Errorf("FieldFilters should be nil when no extra fields present, got %v", filters.FieldFilters)
	}
}
