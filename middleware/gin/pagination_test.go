package gin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaginatedResponse_Basic(t *testing.T) {
	resp := NewPaginatedResponse([]string{"a", "b"}, 10, 1, 5)

	assert.Equal(t, []string{"a", "b"}, resp.Data)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 5, resp.Pagination.PerPage)
	assert.Equal(t, 10, resp.Pagination.Total)
	assert.Equal(t, 2, resp.Pagination.TotalPages)
}

func TestNewPaginatedResponse_TotalPagesRoundsUp(t *testing.T) {
	resp := NewPaginatedResponse(nil, 11, 1, 5)
	assert.Equal(t, 3, resp.Pagination.TotalPages)
}

func TestNewPaginatedResponse_PageLessThanOne(t *testing.T) {
	resp := NewPaginatedResponse(nil, 10, 0, 5)
	assert.Equal(t, 1, resp.Pagination.Page)

	resp = NewPaginatedResponse(nil, 10, -1, 5)
	assert.Equal(t, 1, resp.Pagination.Page)
}

func TestNewPaginatedResponse_LimitZero(t *testing.T) {
	resp := NewPaginatedResponse(nil, 10, 1, 0)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
}

func TestNewPaginatedResponse_EmptyData(t *testing.T) {
	resp := NewPaginatedResponse([]string{}, 0, 1, 20)
	assert.Equal(t, 0, resp.Pagination.Total)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
}

func TestNewPaginatedResponse_SinglePage(t *testing.T) {
	resp := NewPaginatedResponse(nil, 5, 1, 20)
	assert.Equal(t, 1, resp.Pagination.TotalPages)
}

func TestNewPaginatedResponse_ExactFit(t *testing.T) {
	resp := NewPaginatedResponse(nil, 20, 1, 10)
	assert.Equal(t, 2, resp.Pagination.TotalPages)
}

func TestNewPaginatedResponse_NegativeTotal(t *testing.T) {
	resp := NewPaginatedResponse(nil, -5, 1, 10)
	assert.Equal(t, 0, resp.Pagination.Total)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
}

func TestNewPaginatedResponse_NegativeLimit(t *testing.T) {
	resp := NewPaginatedResponse(nil, 10, 1, -3)
	assert.Equal(t, 0, resp.Pagination.PerPage)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
}

func TestNewPaginatedResponse_BothNegative(t *testing.T) {
	resp := NewPaginatedResponse(nil, -5, 1, -3)
	assert.Equal(t, 0, resp.Pagination.Total)
	assert.Equal(t, 0, resp.Pagination.PerPage)
	assert.Equal(t, 0, resp.Pagination.TotalPages)
}
