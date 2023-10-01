package main

import (
	"fmt"

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
	fmt.Println(123)
}
