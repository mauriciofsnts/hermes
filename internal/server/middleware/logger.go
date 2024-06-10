package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		id := r.Context().Value(middleware.RequestIDKey)
		slog.Info(
			"Http request",
			"id", id,
			"remote_addr", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		next.ServeHTTP(ww, r)

		status := ww.Status()

		slog.Info(
			"Http response",
			"id", id,
			"status", status,
			"took", time.Since(start),
		)
	})
}
