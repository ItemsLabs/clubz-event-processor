package config

import "fmt"

func RMQConnectionURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.RMQUser,
		cfg.RMQPassword,
		cfg.RMQHost,
		cfg.RMQPort,
		cfg.RMQVHost,
	)
}

func RMQMatchEventExchange() string {
	return cfg.RMQMatchEventExchange
}

func RMQProcessorQueue() string {
	return cfg.RMQProcessorQueue
}

func RMQFCMExchange() string {
	return cfg.RMQFCMExchange
}

func RMQGamesExchange() string {
	return cfg.RMQGamesExchange
}

func RMQGamesListenerQueue() string {
	return cfg.RMQGamesListenerQueue
}

func RMQSystemExchange() string {
	return cfg.RMQSystemExchange
}

func RMQSystemListenerQueue() string {
	return cfg.RMQSystemListenerQueue
}

func RMQGameUpdatesExchange() string {
	return cfg.RMQGameUpdatesExchange
}
