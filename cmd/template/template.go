package main

import (
	"github.com/ChanJuiHuang/go-backend-framework/v2/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/v2/pkg/booter"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	booter.Boot(
		func() {},
		booter.NewDefaultConfig,
		&registrar.RegisterExecutor,
	)
}

func main() {
}
