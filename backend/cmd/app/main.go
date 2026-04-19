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
	logger.Info("app starting",
		"version", os.Getenv("APP_VERSION"),
		"commit", os.Getenv("COMMIT_SHA"),
		"env", os.Getenv("APP_ENV"),
	)
	app.Run(context.Background(), cfg, logger)
}
