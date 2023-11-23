package main

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	booter.Boot(
		func() {},
		booter.NewDefaultConfig,
		&registrar.ConfigRegistrar,
		&registrar.ServiceRegistrar,
	)
}

func main() {
}
