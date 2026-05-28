package gin

// PaginatedResponse es la respuesta estandar para endpoints paginados.
type PaginatedResponse struct {
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta contiene metadata de paginacion.
type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// NewPaginatedResponse creates a PaginatedResponse with automatic TotalPages calculation.
// The limit parameter maps to the "per_page" field in the JSON response.
// Negative values for total and limit are clamped to 0.
func NewPaginatedResponse(data any, total, page, limit int) PaginatedResponse {
	if total < 0 {
		total = 0
	}
	if limit < 0 {
		limit = 0
	}
	if page < 1 {
		page = 1
	}
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}
	return PaginatedResponse{
		Data: data,
		Pagination: PaginationMeta{
			Page:       page,
			PerPage:    limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
