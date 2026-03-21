package gin

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	commonerrors "github.com/EduGoGroup/edugo-shared/common/errors"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
)

const (
	defaultLimit = 50
	maxLimit     = 200
)

// ParseListFilters extracts standard list filter parameters from the Gin
// request query string into a ListFilters value. It handles is_active, page,
// limit, search, and search_fields. Any extra field names supplied in
// extraFields are also parsed (comma-split) into FieldFilters.
//
// Defensive defaults: when no ?limit= is provided, Limit defaults to 50.
// Values above 200 are capped to 200. When no ?is_active= is provided,
// IsActive remains nil, meaning "show all" (repositories handle nil via
// ApplyIsActive). Clients that need the previous active-only behavior
// should explicitly send ?is_active=true.
//
// Returns a *commonerrors.AppError (HTTP 400) on validation failure, which is
// compatible with the ErrorHandler middleware.
func ParseListFilters(c *gin.Context, extraFields ...string) (sharedrepo.ListFilters, error) {
	var filters sharedrepo.ListFilters

	if err := parseIsActive(c, &filters); err != nil {
		return filters, err
	}
	if err := parseLimit(c, &filters); err != nil {
		return filters, err
	}
	if err := parsePage(c, &filters); err != nil {
		return filters, err
	}
	parseSearch(c, &filters)
	parseExtraFields(c, &filters, extraFields)
	applyDefaults(&filters)

	return filters, nil
}

func parseIsActive(c *gin.Context, f *sharedrepo.ListFilters) error {
	if s := c.Query("is_active"); s != "" {
		v, err := strconv.ParseBool(s)
		if err != nil {
			return commonerrors.NewValidationError("invalid is_active parameter")
		}
		f.IsActive = &v
	}
	return nil
}

func parseLimit(c *gin.Context, f *sharedrepo.ListFilters) error {
	if s := c.Query("limit"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v <= 0 {
			return commonerrors.NewValidationError("limit must be a positive integer")
		}
		f.Limit = v
	}
	return nil
}

func parsePage(c *gin.Context, f *sharedrepo.ListFilters) error {
	if s := c.Query("page"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v <= 0 {
			return commonerrors.NewValidationError("page must be a positive integer")
		}
		f.Page = v
	}
	return nil
}

func parseSearch(c *gin.Context, f *sharedrepo.ListFilters) {
	if search := c.Query("search"); search != "" {
		f.Search = search
		f.SearchFields = splitClean(c.Query("search_fields"))
	}
}

func parseExtraFields(c *gin.Context, f *sharedrepo.ListFilters, extra []string) {
	if len(extra) == 0 {
		return
	}
	fieldFilters := make(map[string][]string)
	for _, field := range extra {
		if vals := splitClean(c.Query(field)); len(vals) > 0 {
			fieldFilters[field] = vals
		}
	}
	if len(fieldFilters) > 0 {
		f.FieldFilters = fieldFilters
	}
}

func applyDefaults(f *sharedrepo.ListFilters) {
	// IsActive nil = "todos" — repositorios ya manejan nil con ApplyIsActive()
	if f.Limit == 0 {
		f.Limit = defaultLimit
	}
	if f.Limit > maxLimit {
		f.Limit = maxLimit
	}
}

// splitClean splits s by comma and returns non-empty trimmed values.
func splitClean(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	clean := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			clean = append(clean, p)
		}
	}
	return clean
}
