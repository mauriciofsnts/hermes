package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mauriciofsnts/hermes/internal/metrics"
)

// MetricsMiddleware tracks HTTP request metrics.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap response writer to capture status code
		wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			metrics.RequestDuration.WithLabelValues(
				r.Method,
				r.URL.Path,
				strconv.Itoa(wrapped.statusCode),
			).Observe(duration)
		}()

		next.ServeHTTP(wrapped, r)
	})
}

// statusResponseWriter wraps http.ResponseWriter to capture the status code.
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
