package controller

import (
	"backend/internal/usecase"
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"

	"backend/config"
	"backend/internal/controller/handlers"
	"backend/internal/controller/middleware"

	redispkg "github.com/redis/go-redis/v9"
)

func NewRouter(ctx context.Context, cfg config.Server, logger *slog.Logger, placeSvc usecase.PlaceService, bookingSvc usecase.BookingService, userSvc usecase.UserService, rdb *redispkg.Client) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)             // Генерирует request_id
	r.Use(middleware.StructuredLogger(logger)) // 👈 Наш новый логгер
	r.Use(middleware.GracefulShutdownMiddleware(ctx, logger))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
	}))

	// Swagger документация
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger-ui.html")
	})

	r.Get("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("swagger.yaml")
		if err != nil {
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write(data)
	})

	jwt := jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)

	// Регистрация обработчиков
	handlers.NewPlaceHandler(r, placeSvc, logger)
	handlers.NewSlotHandler(r, placeSvc, logger)
	handlers.NewAuthHandler(r, userSvc, logger, jwt, rdb)
	handlers.NewBookingHandler(r, bookingSvc, logger, jwt)

	return r
}

func NewServer(cfg config.Server, handler http.Handler) *handlers.HttpServer {
	return handlers.NewServer(cfg, handler)
}
