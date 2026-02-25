package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	"cms-api/internal/infra/telemetry"
	"cms-api/internal/pkg/contextutil"
)

type TracingMiddleware struct {
	tracer *telemetry.Tracer
}

func NewTracingMiddleware(tracer *telemetry.Tracer) *TracingMiddleware {
	return &TracingMiddleware{
		tracer: tracer,
	}
}

func (m *TracingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		spanName := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

		ctx, span := m.tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethod(r.Method),
				semconv.HTTPTarget(r.URL.Path),
				semconv.HTTPScheme(getScheme(r)),
				semconv.NetHostName(r.Host),
				semconv.UserAgentOriginal(r.UserAgent()),
				semconv.HTTPRequestContentLength(int(r.ContentLength)),
			),
		)
		defer span.End()

		if userID := contextutil.GetUserID(ctx); userID != "" {
			span.SetAttributes(telemetry.UserID(userID))
		}

		if requestID := contextutil.GetRequestID(ctx); requestID != "" {
			span.SetAttributes(attribute.String("request.id", requestID))
		}

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r.WithContext(ctx))

		span.SetAttributes(
			semconv.HTTPStatusCode(wrapped.statusCode),
			attribute.Int("http.response_content_length", wrapped.bytesWritten),
		)

		if wrapped.statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(wrapped.statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += n
	return n, err
}

func (w *responseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}
