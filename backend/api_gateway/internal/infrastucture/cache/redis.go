package cache

import (
	"context"
	"fmt"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/config"
	"github.com/redis/go-redis/v9"
	"reflect"
	"time"
)

type Cache struct {
	redis    *redis.Client
	cacheTTL time.Duration
}

func NewCache(cfg *config.Config) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DbNumber,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return &Cache{
		redis:    client,
		cacheTTL: cfg.Redis.CacheTTL,
	}
}

func structToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Tag.Get("redis")
		fieldValue := val.Field(i).Interface()

		if t, ok := fieldValue.(time.Time); ok {
			fieldValue = t.Format(time.RFC3339)
		}

		result[fieldName] = fieldValue
	}

	return result
}
