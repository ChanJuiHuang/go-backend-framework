package util

import (
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/time/rate"
)

type rateLimiterConfig struct {
	PutTokenRate uint `required:"true" split_words:"true"`
	BurstNumber  uint `required:"true" split_words:"true"`
}

var rateLimiterCfg *rateLimiterConfig = &rateLimiterConfig{}

func init() {
	err := envconfig.Process("rate_limiter", rateLimiterCfg)
	if err != nil {
		panic(err)
	}
}

func NewRateLimiter() *rate.Limiter {
	return rate.NewLimiter(
		rate.Limit(rateLimiterCfg.PutTokenRate),
		int(rateLimiterCfg.BurstNumber),
	)
}
