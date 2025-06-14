package config

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress, // Redis server address with port
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDb, // Default database
	})
}
