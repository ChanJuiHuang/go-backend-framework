package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/redis"
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
