package response

import (
	"net/http"

	"go-boilerplate-clean/internal/shared/apperrors"
	"go-boilerplate-clean/internal/shared/pagination"

	"github.com/labstack/echo/v4"
)

// ErrorBody is the standard error object in JSON responses.
type ErrorBody struct {
	Code    apperrors.Code `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// envelope is the uniform API response shape.
type envelope struct {
	Success bool        `json:"success"`
	Data    any         `json:"data,omitempty"`
	Meta    any         `json:"meta,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

// JSON sends a success response with data.
func JSON(c echo.Context, status int, data any) error {
	return c.JSON(status, envelope{
		Success: true,
		Data:    data,
	})
}

// JSONWithMeta sends a success response with data and meta (e.g. pagination).
func JSONWithMeta(c echo.Context, status int, data any, meta any) error {
	return c.JSON(status, envelope{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Paginated sends a success response for paginated lists.
func Paginated[T any](c echo.Context, list pagination.List[T]) error {
	return c.JSON(http.StatusOK, envelope{
		Success: true,
		Data:    list.Items,
		Meta:    list.Meta,
	})
}

// Created sends HTTP 201 with data.
func Created(c echo.Context, data any) error {
	return JSON(c, http.StatusCreated, data)
}

// OK sends HTTP 200 with data.
func OK(c echo.Context, data any) error {
	return JSON(c, http.StatusOK, data)
}

// NoContent sends HTTP 204 without a body.
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// Error maps err to AppError and writes a standard error JSON response.
func Error(c echo.Context, err error) error {
	ae := apperrors.FromError(err)
	return c.JSON(ae.HTTPStatus, envelope{
		Success: false,
		Error: &ErrorBody{
			Code:    ae.Code,
			Message: ae.Message,
			Details: ae.Details,
		},
	})
}

// BindError handles Echo bind/validation failures.
func BindError(c echo.Context, err error) error {
	return Error(c, apperrors.BadRequest("invalid request body"))
}
