package cache

import (
	"context"
	"ratoneando/config"
	"ratoneando/utils/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client  *redis.Client
	enabled bool
)

func Init() {
	logger.Log(config.REDIS_URL)
	opts, err := redis.ParseURL(config.REDIS_URL)
	if err != nil {
		logger.LogWarn("Redis URL inválida — cache deshabilitado")
		return
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.LogWarn("Redis no disponible — cache deshabilitado")
		return
	}

	Client = client
	enabled = true
}

func Set(key string, value string, expiration int) error {
	if !enabled {
		return nil
	}
	ctx := context.Background()
	err := Client.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		logger.LogWarn("Error setting key: " + key)
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	if !enabled {
		return "", nil
	}
	ctx := context.Background()
	value, err := Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
