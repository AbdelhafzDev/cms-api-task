package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"cms-api/internal/pkg/httputil"
)

func Logger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				duration := time.Since(start)

				log.Info("HTTP request",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("query", r.URL.RawQuery),
					zap.Int("status", ww.Status()),
					zap.Int("bytes", ww.BytesWritten()),
					zap.Duration("duration", duration),
					zap.String("ip", httputil.GetClientIP(r)),
					zap.String("x_forwarded_for", r.Header.Get("X-Forwarded-For")),
					zap.String("x_real_ip", r.Header.Get("X-Real-IP")),
					zap.String("cf_connecting_ip", r.Header.Get("CF-Connecting-IP")),
					zap.String("true_client_ip", r.Header.Get("True-Client-IP")),
					zap.String("user_agent", r.UserAgent()),
					zap.String("request_id", middleware.GetReqID(r.Context())),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
