package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

func GracefulShutdownMiddleware(ctx context.Context, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Если контекст отменён — сервер в процессе shutdown
			select {
			case <-ctx.Done():
				logger.Warn("rejecting request during shutdown", "path", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"error": "service shutting down", "retry_after": 5}`))
				return
			default:
				// Контекст активен — обрабатываем запрос
				next.ServeHTTP(w, r)
			}
		})
	}
}
