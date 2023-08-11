package authentication

import (
	"time"
)

type Config struct {
	PrivateKey           string
	PublicKey            string
	AccessTokenLifeTime  time.Duration
	RefreshTokenLifeTime time.Duration
}
