package httputil

import (
	"encoding/json"
	"net/http"

	"cms-api/internal/pkg/apperror"
	"cms-api/internal/pkg/i18nutil"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Error(w http.ResponseWriter, statusCode int, code string, message string, details map[string]interface{}) {
	JSON(w, statusCode, Response{
		Success: false,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, "FORBIDDEN", message, nil)
}

func NotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(w, http.StatusNotFound, "NOT_FOUND", message, nil)
}

func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, "CONFLICT", message, nil)
}

func ValidationError(w http.ResponseWriter, err error) {
	Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", map[string]interface{}{
		"validation": err.Error(),
	})
}

func InternalServerError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "An unexpected error occurred"
	}
	Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	lang := i18nutil.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	statusCode := apperror.HTTPStatusCode(err)
	message := apperror.ErrorMessage(err, lang)

	var code string
	switch statusCode {
	case http.StatusBadRequest:
		code = "BAD_REQUEST"
	case http.StatusUnauthorized:
		code = "UNAUTHORIZED"
	case http.StatusForbidden:
		code = "FORBIDDEN"
	case http.StatusNotFound:
		code = "NOT_FOUND"
	case http.StatusConflict:
		code = "CONFLICT"
	case http.StatusGone:
		code = "GONE"
	default:
		code = "INTERNAL_ERROR"
	}

	Error(w, statusCode, code, message, nil)
}

type SSEWriter func(event, data string) error

func SSERaw(w http.ResponseWriter, _ *http.Request) (SSEWriter, error) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, http.ErrNotSupported
	}

	return func(event, data string) error {
		if event != "" {
			if _, err := w.Write([]byte("event: " + event + "\n")); err != nil {
				return err
			}
		}
		if _, err := w.Write([]byte("data: " + data + "\n\n")); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}, nil
}
