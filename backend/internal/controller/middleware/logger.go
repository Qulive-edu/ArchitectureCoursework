package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware" // 👈 Алиас
)

func StructuredLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor) // 👈 Алиас
			start := time.Now()

			reqID := chimiddleware.GetReqID(r.Context()) // 👈 Алиас

			defer func() {
				logger.Info("http_request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"bytes_written", ww.BytesWritten(),
					"duration_ms", time.Since(start).Milliseconds(),
					"remote_addr", r.RemoteAddr,
					"request_id", reqID,
					"user_agent", r.UserAgent(),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
