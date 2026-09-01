package pagination

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v5"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// Request holds list query parameters (query string or JSON).
type Request struct {
	Page    int `query:"page" json:"page"`
	PerPage int `query:"per_page" json:"per_page"`
}

// Normalize applies defaults and caps page size.
func (r *Request) Normalize() {
	if r.Page < 1 {
		r.Page = DefaultPage
	}
	if r.PerPage < 1 {
		r.PerPage = DefaultPerPage
	}
	if r.PerPage > MaxPerPage {
		r.PerPage = MaxPerPage
	}
}

// Offset returns the SQL/GORM offset for the current page.
func (r *Request) Offset() int {
	r.Normalize()
	return (r.Page - 1) * r.PerPage
}

// Limit returns the page size after normalization.
func (r *Request) Limit() int {
	r.Normalize()
	return r.PerPage
}

// Meta describes pagination metadata in API responses.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewMeta builds pagination metadata from request and total row count.
func NewMeta(req Request, total int64) Meta {
	req.Normalize()
	totalPages := 0
	if req.PerPage > 0 && total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(req.PerPage)))
	}
	return Meta{
		Page:       req.Page,
		PerPage:    req.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

// List is a generic paginated payload.
type List[T any] struct {
	Items []T  `json:"items"`
	Meta  Meta `json:"meta"`
}

// NewList wraps items with pagination metadata.
func NewList[T any](items []T, req Request, total int64) List[T] {
	if items == nil {
		items = []T{}
	}
	return List[T]{
		Items: items,
		Meta:  NewMeta(req, total),
	}
}

// ParseQuery reads page and per_page from Echo query parameters.
func ParseQuery(c *echo.Context) Request {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	req := Request{Page: page, PerPage: perPage}
	req.Normalize()
	return req
}
