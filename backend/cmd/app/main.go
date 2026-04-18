package main

import (
    "backend/config"
    "backend/internal/app"
    "context"
    "flag"
    "log/slog"
    "os"
)

func main() {
    flag.Parse()
    cfg := config.NewConfig()
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
    app.Run(context.Background(), cfg, logger)
}
