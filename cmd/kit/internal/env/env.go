package env

import "github.com/ChanJuiHuang/go-backend-framework/app/config"

func init() {
	config.Log().Level = config.Info
}
