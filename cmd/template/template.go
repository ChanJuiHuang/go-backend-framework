package main

import (
	_ "github.com/joho/godotenv/autoload"

	internalConfig "github.com/ChanJuiHuang/go-backend-framework/internal/config"
	internalProvider "github.com/ChanJuiHuang/go-backend-framework/internal/provider"
)

func init() {
	globalConfig := newGlobalConfig()
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	internalConfig.RegisterConfig(*globalConfig)
	internalProvider.RegisterService()
}

func main() {
}
