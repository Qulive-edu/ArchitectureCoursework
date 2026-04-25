package redis

import (
	"context"
	"fmt"

	"backend/config"

	redispkg "github.com/redis/go-redis/v9"
)

func New(cfg config.Redis) (*redispkg.Client, error) {
	rdb := redispkg.NewClient(&redispkg.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connect: %w", err)
	}

	return rdb, nil
}
