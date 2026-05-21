package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// Code is a stable machine-readable error identifier for clients.
type Code string

const (
	CodeValidation     Code = "VALIDATION_ERROR"
	CodeNotFound       Code = "NOT_FOUND"
	CodeConflict       Code = "CONFLICT"
	CodeBadRequest     Code = "BAD_REQUEST"
	CodeUnauthorized   Code = "UNAUTHORIZED"
	CodeForbidden      Code = "FORBIDDEN"
	CodeInternal       Code = "INTERNAL_ERROR"
	CodeUnprocessable  Code = "UNPROCESSABLE_ENTITY"
)

// AppError is the application error type mapped to HTTP responses.
type AppError struct {
	Code       Code           `json:"code"`
	Message    string         `json:"message"`
	HTTPStatus int            `json:"-"`
	Details    map[string]any `json:"details,omitempty"`
	err        error          `json:"-"`
}

func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.err
}

// Is supports errors.Is for AppError comparison by code.
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

func newAppError(code Code, message string, httpStatus int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		err:        err,
	}
}

func Validation(message string) *AppError {
	return newAppError(CodeValidation, message, http.StatusBadRequest, nil)
}

func ValidationWithDetails(message string, details map[string]any) *AppError {
	e := Validation(message)
	e.Details = details
	return e
}

func NotFound(message string) *AppError {
	return newAppError(CodeNotFound, message, http.StatusNotFound, nil)
}

func BadRequest(message string) *AppError {
	return newAppError(CodeBadRequest, message, http.StatusBadRequest, nil)
}

func Conflict(message string) *AppError {
	return newAppError(CodeConflict, message, http.StatusConflict, nil)
}

func Unauthorized(message string) *AppError {
	return newAppError(CodeUnauthorized, message, http.StatusUnauthorized, nil)
}

func Forbidden(message string) *AppError {
	return newAppError(CodeForbidden, message, http.StatusForbidden, nil)
}

func Unprocessable(message string) *AppError {
	return newAppError(CodeUnprocessable, message, http.StatusUnprocessableEntity, nil)
}

func Internal(message string, err error) *AppError {
	return newAppError(CodeInternal, message, http.StatusInternalServerError, err)
}

// Wrap attaches an underlying error while keeping the app error metadata.
func Wrap(err error, ae *AppError) *AppError {
	if err == nil {
		return ae
	}
	return &AppError{
		Code:       ae.Code,
		Message:    ae.Message,
		HTTPStatus: ae.HTTPStatus,
		Details:    ae.Details,
		err:        err,
	}
}

// AsAppError returns an AppError if err is or wraps one.
func AsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// FromError maps known errors to AppError; unknown errors become internal.
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}
	if ae, ok := AsAppError(err); ok {
		return ae
	}

	switch {
	case errors.Is(err, ErrSampleIDRequired):
		return ErrSampleIDRequired
	case errors.Is(err, ErrSampleNotFound):
		return ErrSampleNotFound
	case errors.Is(err, ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(err, ErrUserIDRequired):
		return ErrUserIDRequired
	case errors.Is(err, ErrUserNameRequired):
		return ErrUserNameRequired
	case errors.Is(err, ErrUserEmailRequired):
		return ErrUserEmailRequired
	}

	// Legacy string-based errors from repository / usecase.
	msg := err.Error()
	switch msg {
	case "sample not found", "user not found":
		return NotFound(msg)
	case "sample id is required", "id is required", "name is required", "email is required":
		return Validation(msg)
	case "invalid request body":
		return BadRequest(msg)
	}

	return Internal("an unexpected error occurred", err)
}
