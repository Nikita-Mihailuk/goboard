package redis

import (
	"context"
	"fmt"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/config"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	db *redis.Client
}

func NewStorage(cfg *config.Config) *Storage {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.DB.Host, cfg.DB.Port),
		Password: cfg.DB.Password,
		DB:       cfg.DB.DbNumber,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return &Storage{
		db: client,
	}
}
