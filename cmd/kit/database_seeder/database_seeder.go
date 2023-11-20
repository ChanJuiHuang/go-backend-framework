package main

import (
	"flag"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"

	"github.com/ChanJuiHuang/go-backend-framework/internal/migration/seeder"
	internalProvider "github.com/ChanJuiHuang/go-backend-framework/internal/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/provider"
)

func init() {
	globalConfig := newGlobalConfig()
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	registerConfig(*globalConfig)
	internalProvider.RegisterService()
}

func main() {
	var seeders string
	flag.StringVar(&seeders, "seeders", "", "Type the seeders. EX: seeder1,seeder2")
	flag.Parse()

	db := provider.Registry.Get("database").(*gorm.DB)
	seeder.Run(db, strings.Split(seeders, ","))
}
