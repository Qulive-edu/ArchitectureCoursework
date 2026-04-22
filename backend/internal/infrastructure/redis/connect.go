package redis

import (
	"context"
	"fmt"

	"backend/config" // 👈 Твой пакет с конфигом

	redispkg "github.com/redis/go-redis/v9"
)

// New принимает config.Redis и возвращает подключенный клиент
func New(cfg config.Redis) (*redispkg.Client, error) {
	rdb := redispkg.NewClient(&redispkg.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Проверяем подключение
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connect: %w", err)
	}

	return rdb, nil
}
