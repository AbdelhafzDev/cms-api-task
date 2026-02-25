package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

const MaxBodySize = 10 * 1024 * 1024

func DecodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		var maxBytesErr *http.MaxBytesError
		switch {
		case errors.As(err, &maxBytesErr):
			return fmt.Errorf("request body too large")
		case errors.Is(err, io.EOF):
			return fmt.Errorf("request body is empty")
		default:
			return fmt.Errorf("invalid JSON: %w", err)
		}
	}

	if decoder.More() {
		return fmt.Errorf("request body must contain only one JSON object")
	}

	return nil
}

func GetQueryParam(r *http.Request, key string, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

func GetBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}

// Priority order: X-Forwarded-For (first IP), then single-IP proxy headers, then RemoteAddr
func GetClientIP(r *http.Request) string {
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		return strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
	}

	for _, header := range []string{"True-Client-IP", "CF-Connecting-IP", "X-Real-IP"} {
		if ip := strings.TrimSpace(r.Header.Get(header)); ip != "" {
			return ip
		}
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
