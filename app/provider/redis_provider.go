package provider

import (
	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	"github.com/go-redis/redis/v9"
)

func provideRedis() *redis.Client {
	redisConfig := config.Redis()
	if !redisConfig.Enabled {
		return new(redis.Client)
	}

	return redis.NewClient(&redis.Options{
		Addr:            redisConfig.Address,
		Password:        redisConfig.Password,
		DB:              redisConfig.DB,
		MinIdleConns:    redisConfig.MinIdleConns,
		ConnMaxLifetime: redisConfig.ConnMaxLifetime,
	})
}
