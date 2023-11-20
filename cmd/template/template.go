package main

import (
	_ "github.com/joho/godotenv/autoload"

	internalProvider "github.com/ChanJuiHuang/go-backend-framework/internal/provider"
)

func init() {
	globalConfig := newGlobalConfig()
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	registerConfig(*globalConfig)
	internalProvider.RegisterService()
}

func main() {
}
