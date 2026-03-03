package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// validFieldName matches only safe column names (alphanumeric + underscore).
var validFieldName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// ilikEscapeClause is appended to every ILIKE condition so that backslash is
// recognised as the escape character by PostgreSQL.
const ilikEscapeClause = "ESCAPE '\\'"

// ListFilters represents common filters for listing entities.
//
// The Search and SearchFields fields enable flexible text search across multiple
// database columns using case-insensitive ILIKE matching.
//
// SearchFields must contain database column names (not Go struct field names).
// Each name must match the pattern ^[a-zA-Z_][a-zA-Z0-9_]*$ (letters, digits and
// underscores only). Invalid names are silently skipped to prevent SQL injection.
//
// Example - search users by name or email:
//
//	filters := repository.ListFilters{
//	    Search:       "john",
//	    SearchFields: []string{"name", "email"},
//	}
//	users, err := userRepo.List(ctx, filters)
//	// Executes: WHERE name ILIKE '%john%' OR email ILIKE '%john%'
//
// Example - paginated listing without search:
//
//	filters := repository.ListFilters{
//	    Limit:  20,
//	    Offset: 0,
//	}
type ListFilters struct {
	IsActive *bool
	Limit    int
	Offset   int
	// Page is a 1-based page number. When Page > 1 and Limit > 0, the offset is
	// calculated as (Page-1)*Limit, taking precedence over the Offset field.
	Page int
	// Search is the text to look for. It is applied with ILIKE '%value%' against
	// every column listed in SearchFields. An empty Search skips the search clause.
	Search string
	// SearchFields lists the database column names to search in.
	// Use snake_case column names as they appear in the database schema,
	// e.g. []string{"first_name", "last_name", "email"}.
	SearchFields []string
}

// escapeLikePattern escapes PostgreSQL ILIKE special characters in s so they
// are treated as literals instead of wildcards. The substitutions must happen
// in this order:
//  1. Backslash (\) – the escape character itself, must be escaped first to
//     avoid double-escaping the sequences added in the next steps.
//  2. Percent sign (%) – matches any sequence of characters in ILIKE.
//  3. Underscore (_) – matches any single character in ILIKE.
func escapeLikePattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// ApplySearch adds ILIKE search conditions to the given GORM query.
// It validates field names, escapes ILIKE special characters in the search term,
// and builds OR conditions like:
//
//	field1 ILIKE '%search%' ESCAPE '\' OR field2 ILIKE '%search%' ESCAPE '\'
func (f ListFilters) ApplySearch(query *gorm.DB) *gorm.DB {
	if f.Search == "" || len(f.SearchFields) == 0 {
		return query
	}
	var conditions []string
	var args []interface{}
	escaped := escapeLikePattern(f.Search)
	for _, field := range f.SearchFields {
		if !validFieldName.MatchString(field) {
			continue
		}
		conditions = append(conditions, fmt.Sprintf("%s ILIKE ? %s", field, ilikEscapeClause))
		args = append(args, "%"+escaped+"%")
	}
	if len(conditions) == 0 {
		return query
	}
	return query.Where(strings.Join(conditions, " OR "), args...)
}

// ApplyPagination applies LIMIT and OFFSET to the query when Limit > 0.
func (f *ListFilters) ApplyPagination(db *gorm.DB) *gorm.DB {
	if f.Limit > 0 {
		db = db.Limit(f.Limit)
		db = db.Offset(f.GetOffset())
	}
	return db
}

// GetOffset returns the calculated offset. If Page > 1 and Limit > 0, the
// offset is derived from (Page-1)*Limit; otherwise it falls back to Offset.
func (f *ListFilters) GetOffset() int {
	if f.Page > 1 && f.Limit > 0 {
		return (f.Page - 1) * f.Limit
	}
	return f.Offset
}

// UserRepository defines persistence operations for User
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters ListFilters) ([]*entities.User, int, error)
}

// MembershipRepository defines persistence operations for Membership
type MembershipRepository interface {
	Create(ctx context.Context, membership *entities.Membership) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error)
	FindByUser(ctx context.Context, userID uuid.UUID, filters ListFilters) ([]*entities.Membership, int, error)
	FindByUnit(ctx context.Context, unitID uuid.UUID, filters ListFilters) ([]*entities.Membership, int, error)
	FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool, filters ListFilters) ([]*entities.Membership, int, error)
	FindByUserAndSchool(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error)
	Update(ctx context.Context, membership *entities.Membership) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SchoolRepository defines persistence operations for School
type SchoolRepository interface {
	Create(ctx context.Context, school *entities.School) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error)
	FindByCode(ctx context.Context, code string) (*entities.School, error)
	Update(ctx context.Context, school *entities.School) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters ListFilters) ([]*entities.School, int, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
}
