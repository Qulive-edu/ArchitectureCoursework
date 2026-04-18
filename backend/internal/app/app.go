package app

import (
    "context"
    "backend/config"
    "backend/internal/controller"
    "backend/internal/infrastructure/postgres"
    "backend/internal/usecase"
    "log/slog"
    "os"
    "os/signal"
    "syscall"
)

func Run(ctx context.Context, cfg config.Config, logger *slog.Logger) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    pg, err := postgres.New(ctx, cfg.Database, logger)
    if err != nil {
        logger.Error("db connect: " + err.Error())
        panic(err)
    }
    defer pg.Close()

    placeRepo := postgres.NewPlaceRepo(pg)
    slotRepo := postgres.NewSlotRepo(pg)
    bookingRepo := postgres.NewBookingRepo(pg)
    userRepo := postgres.NewUserRepo(pg)

    bookingService := usecase.NewBookingService(bookingRepo, slotRepo)
    placeService := usecase.NewPlaceService(placeRepo, slotRepo)
    userService := usecase.NewUserService(userRepo)

    router := controller.NewRouter(cfg.Server, logger, placeService, bookingService, userService)

    srv := controller.NewServer(cfg.Server, router)
    logger.Info("starting server on port " + cfg.Server.Port)

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

    select {
    case s := <-interrupt:
        logger.Info(s.String() + " signal received")
    case err = <-srv.Notify():
        logger.Error("shutdown via http server error: " + err.Error())
    case <-ctx.Done():
        logger.Info("context canceled")
    }

    logger.Info("shutting down...")
    cancel()
    _ = srv.Shutdown()
    logger.Info("server shut down")
}
