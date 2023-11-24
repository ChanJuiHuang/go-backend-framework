package main

import (
	"flag"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"

	"github.com/ChanJuiHuang/go-backend-framework/internal/migration/seeder"
	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
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
	var seeders string
	flag.StringVar(&seeders, "seeders", "", "Type the seeders. EX: seeder1,seeder2")
	flag.Parse()

	db := service.Registry.Get("database").(*gorm.DB)
	seeder.Run(db, strings.Split(seeders, ","))
}
