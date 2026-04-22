package main

import (
	"backend/config"
	"backend/internal/app"
	"backend/internal/infrastructure/redis" // 👈 твой пакет подключения
	"context"
	"flag"
	"log/slog"
	"os"
)

func main() {
	flag.Parse()
	cfg := config.NewConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	// 👇 1. Инициализируем Redis
	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		logger.Error("redis connect: " + err.Error())
		panic(err) // или os.Exit(1), как у тебя с БД
	}
	defer rdb.Close()

	// 👇 2. Передаём rdb в app.Run
	app.Run(context.Background(), cfg, logger, rdb)
}
