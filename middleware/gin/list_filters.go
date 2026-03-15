package gin

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	commonerrors "github.com/EduGoGroup/edugo-shared/common/errors"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
)

// ParseListFilters extracts standard list filter parameters from the Gin
// request query string into a ListFilters value. It handles is_active, page,
// limit, search, and search_fields. Any extra field names supplied in
// extraFields are also parsed (comma-split) into FieldFilters.
//
// Returns a *commonerrors.AppError (HTTP 400) on validation failure, which is
// compatible with the ErrorHandler middleware.
func ParseListFilters(c *gin.Context, extraFields ...string) (sharedrepo.ListFilters, error) {
	var filters sharedrepo.ListFilters

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			return filters, commonerrors.NewValidationError("invalid is_active parameter")
		}
		filters.IsActive = &isActive
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return filters, commonerrors.NewValidationError("limit must be a positive integer")
		}
		filters.Limit = limit
	}

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			return filters, commonerrors.NewValidationError("page must be a positive integer")
		}
		filters.Page = page
	}

	if search := c.Query("search"); search != "" {
		filters.Search = search
		if fields := c.Query("search_fields"); fields != "" {
			rawFields := strings.Split(fields, ",")
			cleanFields := make([]string, 0, len(rawFields))
			for _, f := range rawFields {
				if f = strings.TrimSpace(f); f != "" {
					cleanFields = append(cleanFields, f)
				}
			}
			if len(cleanFields) > 0 {
				filters.SearchFields = cleanFields
			}
		}
	}

	if len(extraFields) > 0 {
		fieldFilters := make(map[string][]string)
		for _, field := range extraFields {
			if v := c.Query(field); v != "" {
				parts := strings.Split(v, ",")
				clean := make([]string, 0, len(parts))
				for _, p := range parts {
					if p = strings.TrimSpace(p); p != "" {
						clean = append(clean, p)
					}
				}
				if len(clean) > 0 {
					fieldFilters[field] = clean
				}
			}
		}
		if len(fieldFilters) > 0 {
			filters.FieldFilters = fieldFilters
		}
	}

	// Default and cap limit to protect against unbounded queries
	if filters.Limit == 0 {
		filters.Limit = 50
	}
	if filters.Limit > 200 {
		filters.Limit = 200
	}

	return filters, nil
}
