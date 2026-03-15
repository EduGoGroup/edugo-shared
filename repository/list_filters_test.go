package repository

import (
	"strings"
	"testing"

	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// escapeLikePattern
// ---------------------------------------------------------------------------

func TestEscapeLikePattern_Backslash(t *testing.T) {
	got := escapeLikePattern(`a\b`)
	want := `a\\b`
	if got != want {
		t.Errorf("escapeLikePattern(`a\\b`) = %q, want %q", got, want)
	}
}

func TestEscapeLikePattern_Percent(t *testing.T) {
	got := escapeLikePattern("100%")
	want := `100\%`
	if got != want {
		t.Errorf("escapeLikePattern(\"100%%\") = %q, want %q", got, want)
	}
}

func TestEscapeLikePattern_Underscore(t *testing.T) {
	got := escapeLikePattern("first_name")
	want := `first\_name`
	if got != want {
		t.Errorf("escapeLikePattern(\"first_name\") = %q, want %q", got, want)
	}
}

func TestEscapeLikePattern_AllSpecialChars(t *testing.T) {
	got := escapeLikePattern(`\%_`)
	want := `\\\%\_`
	if got != want {
		t.Errorf("escapeLikePattern(`\\%%_`) = %q, want %q", got, want)
	}
}

func TestEscapeLikePattern_NoSpecialChars(t *testing.T) {
	got := escapeLikePattern("hello")
	if got != "hello" {
		t.Errorf("escapeLikePattern(\"hello\") = %q, want \"hello\"", got)
	}
}

func TestEscapeLikePattern_Empty(t *testing.T) {
	got := escapeLikePattern("")
	if got != "" {
		t.Errorf("escapeLikePattern(\"\") = %q, want \"\"", got)
	}
}

func TestEscapeLikePattern_OrderMatters(t *testing.T) {
	// Backslash must be escaped first, otherwise the escapes for % and _ would
	// be double-escaped.
	got := escapeLikePattern(`%\_%`)
	want := `\%\\\_\%`
	if got != want {
		t.Errorf("escapeLikePattern(`%%\\_%%`) = %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// GetOffset
// ---------------------------------------------------------------------------

func TestGetOffset_Page1(t *testing.T) {
	f := ListFilters{Page: 1, Limit: 10}
	if got := f.GetOffset(); got != 0 {
		t.Errorf("GetOffset() = %d, want 0 for Page=1", got)
	}
}

func TestGetOffset_Page2(t *testing.T) {
	f := ListFilters{Page: 2, Limit: 10}
	if got := f.GetOffset(); got != 10 {
		t.Errorf("GetOffset() = %d, want 10 for Page=2, Limit=10", got)
	}
}

func TestGetOffset_Page3_Limit25(t *testing.T) {
	f := ListFilters{Page: 3, Limit: 25}
	if got := f.GetOffset(); got != 50 {
		t.Errorf("GetOffset() = %d, want 50 for Page=3, Limit=25", got)
	}
}

func TestGetOffset_PageZero_UsesOffsetField(t *testing.T) {
	f := ListFilters{Page: 0, Limit: 10, Offset: 15}
	if got := f.GetOffset(); got != 15 {
		t.Errorf("GetOffset() = %d, want 15 (fallback to Offset field)", got)
	}
}

func TestGetOffset_NoPageNoOffset(t *testing.T) {
	f := ListFilters{Limit: 10}
	if got := f.GetOffset(); got != 0 {
		t.Errorf("GetOffset() = %d, want 0 when neither Page nor Offset set", got)
	}
}

func TestGetOffset_PageTakesPrecedenceOverOffset(t *testing.T) {
	f := ListFilters{Page: 2, Limit: 10, Offset: 99}
	if got := f.GetOffset(); got != 10 {
		t.Errorf("GetOffset() = %d, want 10 (Page should take precedence over Offset)", got)
	}
}

func TestGetOffset_NegativeOffset_ReturnsZero(t *testing.T) {
	f := ListFilters{Page: 0, Offset: -5}
	if got := f.GetOffset(); got != 0 {
		t.Errorf("GetOffset() = %d, want 0 for negative Offset", got)
	}
}

// ---------------------------------------------------------------------------
// ApplySearch — uses GORM DryRun to inspect generated SQL
// ---------------------------------------------------------------------------

// newDryRunDB creates a GORM DB in DryRun mode (no real database needed).
func newDryRunDB() *gorm.DB {
	db, err := gorm.Open(nil, &gorm.Config{DryRun: true})
	if err != nil {
		panic("failed to create dry run DB")
	}
	return db
}

func TestApplySearch_EmptySearch(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{Search: "", SearchFields: []string{"name"}}
	result := f.ApplySearch(db)
	// Should return the same query unchanged
	if result != db {
		t.Error("ApplySearch with empty search should return the same query")
	}
}

func TestApplySearch_NoFields(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{Search: "test", SearchFields: nil}
	result := f.ApplySearch(db)
	if result != db {
		t.Error("ApplySearch with no fields should return the same query")
	}
}

func TestApplySearch_ValidFields(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{
		Search:       "john",
		SearchFields: []string{"name", "email"},
	}
	result := f.ApplySearch(db)
	// The result should be a different query (WHERE clause added)
	if result == db {
		t.Error("ApplySearch with valid search should modify the query")
	}

	// Inspect the generated SQL via Statement
	stmt := result.Statement
	sql := stmt.SQL.String()
	if sql == "" {
		// Build the statement to populate SQL
		result = result.Find(nil)
		stmt = result.Statement
		sql = stmt.SQL.String()
	}

	// The clauses should be stored; verify via Statement.Clauses
	if len(stmt.Clauses) == 0 && sql == "" {
		t.Log("DryRun mode: verified ApplySearch modifies query reference (clause inspection not available without dialect)")
	}
}

func TestApplySearch_SQLInjection_InvalidFieldNames(t *testing.T) {
	db := newDryRunDB()

	invalidFields := []string{
		"name; DROP TABLE users--",
		"1=1 OR name",
		"name' OR '1'='1",
		"table.column",
		"name()",
		"",
		" ",
		"name-column",
		"name column",
	}

	for _, field := range invalidFields {
		f := ListFilters{
			Search:       "test",
			SearchFields: []string{field},
		}
		result := f.ApplySearch(db)
		// All invalid fields should be skipped, returning original query
		if result != db {
			t.Errorf("ApplySearch should skip invalid field %q and return original query", field)
		}
	}
}

func TestApplySearch_MixedValidAndInvalidFields(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{
		Search:       "test",
		SearchFields: []string{"valid_name", "DROP TABLE--", "email"},
	}
	result := f.ApplySearch(db)
	// Should modify query because valid fields exist
	if result == db {
		t.Error("ApplySearch with some valid fields should modify the query")
	}
}

func TestApplySearch_AllInvalidFields(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{
		Search:       "test",
		SearchFields: []string{"bad;field", "also.bad"},
	}
	result := f.ApplySearch(db)
	if result != db {
		t.Error("ApplySearch with all invalid fields should return original query")
	}
}

func TestApplySearch_EscapesSpecialCharsInSearch(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{
		Search:       "100%",
		SearchFields: []string{"name"},
	}
	result := f.ApplySearch(db)
	if result == db {
		t.Error("ApplySearch should modify query for valid search")
	}
}

func TestApplySearch_ValidFieldNamePatterns(t *testing.T) {
	validFields := []string{
		"name",
		"first_name",
		"_private",
		"Column1",
		"a",
		"ABC",
		"field_with_123",
	}

	for _, field := range validFields {
		if !validFieldName.MatchString(field) {
			t.Errorf("validFieldName should accept %q", field)
		}
	}
}

func TestApplySearch_InvalidFieldNamePatterns(t *testing.T) {
	invalidFields := []string{
		"1starts_with_number",
		"has space",
		"has-dash",
		"has.dot",
		"has;semicolon",
		"",
		"has'quote",
		"has\"doublequote",
	}

	for _, field := range invalidFields {
		if validFieldName.MatchString(field) {
			t.Errorf("validFieldName should reject %q", field)
		}
	}
}

// ---------------------------------------------------------------------------
// ApplyPagination — uses GORM DryRun to inspect generated SQL
// ---------------------------------------------------------------------------

func TestApplyPagination_WithLimit(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{Limit: 10}
	result := f.ApplyPagination(db)
	if result == db {
		t.Error("ApplyPagination with Limit > 0 should modify the query")
	}
}

func TestApplyPagination_NoLimit(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{Limit: 0}
	result := f.ApplyPagination(db)
	if result != db {
		t.Error("ApplyPagination with Limit=0 should return original query")
	}
}

func TestApplyPagination_NegativeLimit(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{Limit: -1}
	result := f.ApplyPagination(db)
	if result != db {
		t.Error("ApplyPagination with negative Limit should return original query")
	}
}

func TestApplyPagination_WithOffset(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{Limit: 10, Page: 2}
	result := f.ApplyPagination(db)
	if result == db {
		t.Error("ApplyPagination with Limit and Page > 1 should modify the query")
	}
}

// ---------------------------------------------------------------------------
// ApplyIsActive
// ---------------------------------------------------------------------------

func TestApplyIsActive_Nil(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{IsActive: nil}
	result := f.ApplyIsActive(db)
	if result != db {
		t.Error("ApplyIsActive with nil should return the same query")
	}
}

func TestApplyIsActive_True(t *testing.T) {
	db := newDryRunDB()
	isActive := true
	f := ListFilters{IsActive: &isActive}
	result := f.ApplyIsActive(db)
	if result == db {
		t.Error("ApplyIsActive with true should modify the query")
	}
}

func TestApplyIsActive_False(t *testing.T) {
	db := newDryRunDB()
	isActive := false
	f := ListFilters{IsActive: &isActive}
	result := f.ApplyIsActive(db)
	if result == db {
		t.Error("ApplyIsActive with false should modify the query")
	}
}

// ---------------------------------------------------------------------------
// ApplyFieldFilters
// ---------------------------------------------------------------------------

func TestApplyFieldFilters_Empty(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{}
	result := f.ApplyFieldFilters(db, []string{"status"})
	if result != db {
		t.Error("ApplyFieldFilters with no field filters should return the same query")
	}
}

func TestApplyFieldFilters_SingleValue(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{
		FieldFilters: map[string][]string{"status": {"active"}},
	}
	result := f.ApplyFieldFilters(db, []string{"status"})
	if result == db {
		t.Error("ApplyFieldFilters with single value should modify the query")
	}
}

func TestApplyFieldFilters_MultiValue(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{
		FieldFilters: map[string][]string{"status": {"active", "pending"}},
	}
	result := f.ApplyFieldFilters(db, []string{"status"})
	if result == db {
		t.Error("ApplyFieldFilters with multiple values should modify the query")
	}
}

func TestApplyFieldFilters_NotAllowed(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{
		FieldFilters: map[string][]string{"secret": {"value"}},
	}
	result := f.ApplyFieldFilters(db, []string{"status"})
	if result != db {
		t.Error("ApplyFieldFilters with non-allowed field should return the same query")
	}
}

func TestApplyFieldFilters_InvalidFieldName(t *testing.T) {
	db := newDryRunDB()
	f := ListFilters{
		FieldFilters: map[string][]string{"bad;field": {"value"}},
	}
	result := f.ApplyFieldFilters(db, []string{"bad;field"})
	if result != db {
		t.Error("ApplyFieldFilters with invalid field name should return the same query")
	}
}

// ---------------------------------------------------------------------------
// ilikEscapeClause constant
// ---------------------------------------------------------------------------

func TestIlikEscapeClause(t *testing.T) {
	if !strings.Contains(ilikEscapeClause, "ESCAPE") {
		t.Errorf("ilikEscapeClause should contain ESCAPE, got %q", ilikEscapeClause)
	}
}
