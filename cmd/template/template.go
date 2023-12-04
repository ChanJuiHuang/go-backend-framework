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
		&registrar.RegisterExecutor,
	)
}

func main() {
}
