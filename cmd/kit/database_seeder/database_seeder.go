package main

import (
	"flag"
	"strings"

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
	var seeders string
	flag.StringVar(&seeders, "seeders", "", "Type the seeders. EX: seeder1,seeder2")
	flag.Parse()

	seeder.Run(provider.Registry.DB(), strings.Split(seeders, ","))
}
