package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/ChanJuiHuang/go-backend-framework/internal/migration/seeder"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
)

func init() {
	globalConfig := newGlobalConfig()
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	registerConfig(*globalConfig)

	registerProvider()
}

func main() {
	seeder.Run(provider.Registry.DB())
}
