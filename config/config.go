package config

import (
	"github.com/caarlos0/env"
)

// Represents a structure with all env variables needed by the backend.
var cfg struct {
	DatabaseUser           string `env:"DATABASE_USER,required"`
	DatabasePassword       string `env:"DATABASE_PASSWORD,required"`
	DatabaseHost           string `env:"DATABASE_HOST,required"`
	DatabasePort           int    `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseName           string `env:"DATABASE_NAME,required"`
	DatabaseSSLMode        string `env:"DATABASE_SSLMODE" envDefault:"disable"`
	RMQHost                string `env:"RMQ_HOST,required"`
	RMQPort                int    `env:"RMQ_PORT,required"`
	RMQVHost               string `env:"RMQ_VHOST,required"`
	RMQUser                string `env:"RMQ_USER,required"`
	RMQPassword            string `env:"RMQ_PASSWORD,required"`
	RMQMatchEventExchange  string `env:"RMQ_MATCH_EVENT_EXCHANGE,required"`
	RMQProcessorQueue      string `env:"RMQ_PROCESSOR_QUEUE,required"`
	RMQFCMExchange         string `env:"RMQ_FCM_EXCHANGE,required"`
	RMQGamesExchange       string `env:"RMQ_GAMES_EXCHANGE,required"`
	RMQGamesListenerQueue  string `env:"RMQ_GAMES_LISTENER_QUEUE,required"`
	RMQSystemExchange      string `env:"RMQ_SYSTEM_EXCHANGE,required"`
	RMQSystemListenerQueue string `env:"RMQ_SYSTEM_LISTENER_QUEUE,required"`
	RMQGameUpdatesExchange string `env:"RMQ_GAME_UPDATES_EXCHANGE,required"`
	RedisAddress           string `env:"REDIS_ADDRESS,required"`
	RedisPassword          string `env:"REDIS_PASSWORD,required"`
	RedisDb                int    `env:"REDIS_DB,required"`
}

func init() {
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
}
