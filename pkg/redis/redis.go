package redis

import "github.com/redis/go-redis/v9"

func New(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:            config.Address,
		Password:        config.Password,
		DB:              config.DB,
		MinIdleConns:    config.MinIdleConns,
		ConnMaxLifetime: config.ConnMaxLifetime,
	})
}
