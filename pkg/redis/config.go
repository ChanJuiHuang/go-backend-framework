package redis

import (
	"time"
)

type Config struct {
	Address         string
	Password        string
	DB              int
	MinIdleConns    int
	ConnMaxLifetime time.Duration
}
