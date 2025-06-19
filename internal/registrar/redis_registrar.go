package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/service"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/redis"
)

type RedisRegistrar struct {
	config redis.Config
}

func (rr *RedisRegistrar) Boot() {
	config.Registry.Register("redis", &rr.config)
}

func (rr *RedisRegistrar) Register() {
	service.Registry.Set("redis", redis.New(rr.config))
}
