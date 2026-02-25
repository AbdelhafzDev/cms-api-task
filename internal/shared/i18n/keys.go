package i18n

type Key string

const (
	ErrNotFound           Key = "error.not_found"
	ErrBadRequest         Key = "error.bad_request"
	ErrUnauthorized       Key = "error.unauthorized"
	ErrForbidden          Key = "error.forbidden"
	ErrConflict           Key = "error.conflict"
	ErrInternalServer     Key = "error.internal_server"
	ErrServiceUnavailable Key = "error.service_unavailable"
	ErrInvalidCredentials Key = "error.invalid_credentials"
	ErrInvalidToken       Key = "error.invalid_token"
	ErrEmailAlreadyExists Key = "error.email_already_exists"
	ErrUserInactive       Key = "error.user_inactive"
	ErrTokenExpired       Key = "error.token_expired"
	ErrTokenRevoked       Key = "error.token_revoked"
	ErrValidationFailed   Key = "error.validation_failed"
)
