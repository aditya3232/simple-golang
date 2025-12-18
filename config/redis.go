package config

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (cfg Config) RedisConfig() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		DB:       cfg.Redis.DB,
		Password: cfg.Redis.Password,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionRedis-1] Failed to get redis connection")
	}

	return redisClient
}
