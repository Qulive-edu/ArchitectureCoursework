package app

import (
	"backend/config"
	"backend/internal/controller"
	"backend/internal/infrastructure/postgres"
	"backend/internal/usecase"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	redispkg "github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, cfg config.Config, logger *slog.Logger, rdb *redispkg.Client) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pg, err := postgres.New(ctx, cfg.Database, logger)
	if err != nil {
		logger.Error("db connect: " + err.Error())
		panic(err)
	}
	//defer pg.Close()

	placeRepo := postgres.NewPlaceRepo(pg)
	slotRepo := postgres.NewSlotRepo(pg)
	bookingRepo := postgres.NewBookingRepo(pg)
	userRepo := postgres.NewUserRepo(pg)

	bookingService := usecase.NewBookingService(bookingRepo, slotRepo)
	placeService := usecase.NewPlaceService(placeRepo, slotRepo)
	userService := usecase.NewUserService(userRepo)

	router := controller.NewRouter(ctx, cfg.Server, logger, placeService, bookingService, userService, rdb)

	srv := controller.NewServer(cfg.Server, router)
	logger.Info("starting server on port " + cfg.Server.Port)

	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	var shutdownErr error
	select {
	case s := <-interrupt:
		logger.Info(s.String() + " signal received")
	case err = <-srv.Notify():
		logger.Error("shutdown via http server error: " + err.Error())
		shutdownErr = err
	case <-ctx.Done():
		logger.Info("context canceled")
	}

	logger.Info("shutting down...")

	if err := srv.Shutdown(); err != nil {
		logger.Error("http server shutdown error", "error", err.Error())
		if shutdownErr == nil {
			shutdownErr = err
		}
	}
	logger.Info("http server stopped")

	pg.Close()
	logger.Info("postgres connection closed")

	if err := rdb.Close(); err != nil {
		logger.Error("redis close error", "error", err.Error())
	}
	logger.Info("redis connection closed")

	cancel()
	logger.Info("application shutdown complete")
	if shutdownErr != nil {
		os.Exit(1)
	}
}
