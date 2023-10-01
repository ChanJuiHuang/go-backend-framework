package main

import (
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	globalConfig := newGlobalConfig()
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	registerConfig(*globalConfig)

	registerProvider()
}

func main() {
}
