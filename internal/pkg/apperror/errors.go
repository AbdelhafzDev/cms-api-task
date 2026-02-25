package apperror

import (
	"errors"
	"net/http"

	"cms-api/internal/shared/i18n"
)

var (
	ErrNotFound           = errors.New("resource not found")
	ErrBadRequest         = errors.New("bad request")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrConflict           = errors.New("conflict")
	ErrInternalServer     = errors.New("internal server error")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenRevoked       = errors.New("token has been revoked")
	ErrValidationFailed   = errors.New("validation failed")
)

type AppError struct {
	Err        error
	Message    string
	StatusCode int
	Details    map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(err error, message string, statusCode int) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

func HTTPStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}

	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrBadRequest), errors.Is(err, ErrValidationFailed):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized), errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrInvalidToken), errors.Is(err, ErrTokenExpired), errors.Is(err, ErrTokenRevoked):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrUserInactive):
		return http.StatusForbidden
	case errors.Is(err, ErrConflict), errors.Is(err, ErrEmailAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, ErrServiceUnavailable):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func I18nKey(err error) i18n.Key {
	switch {
	case errors.Is(err, ErrNotFound):
		return i18n.ErrNotFound
	case errors.Is(err, ErrBadRequest):
		return i18n.ErrBadRequest
	case errors.Is(err, ErrUnauthorized):
		return i18n.ErrUnauthorized
	case errors.Is(err, ErrForbidden):
		return i18n.ErrForbidden
	case errors.Is(err, ErrInvalidCredentials):
		return i18n.ErrInvalidCredentials
	case errors.Is(err, ErrInvalidToken):
		return i18n.ErrInvalidToken
	case errors.Is(err, ErrEmailAlreadyExists):
		return i18n.ErrEmailAlreadyExists
	case errors.Is(err, ErrUserInactive):
		return i18n.ErrUserInactive
	case errors.Is(err, ErrTokenExpired):
		return i18n.ErrTokenExpired
	case errors.Is(err, ErrTokenRevoked):
		return i18n.ErrTokenRevoked
	case errors.Is(err, ErrValidationFailed):
		return i18n.ErrValidationFailed
	case errors.Is(err, ErrConflict):
		return i18n.ErrConflict
	case errors.Is(err, ErrServiceUnavailable):
		return i18n.ErrServiceUnavailable
	default:
		return i18n.ErrInternalServer
	}
}

func ErrorMessage(err error, lang string) string {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Message != "" {
		return appErr.Message
	}

	return i18n.GetMessage(I18nKey(err), lang)
}
