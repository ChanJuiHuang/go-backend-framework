package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/redis"
)

type RedisRegistrar struct{}

func (*RedisRegistrar) Register() {
	service.Registry.Set("redis", redis.New(config.Registry.Get("redis").(redis.Config)))
}
