package main

import (
	"os"
	"path"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ChanJuiHuang/go-backend-framework/internal/migration/seeder"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
	"github.com/spf13/viper"
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	byteYaml, err := os.ReadFile(path.Join(wd, "config.yml"))
	if err != nil {
		panic(err)
	}
	stringYaml := os.ExpandEnv(string(byteYaml))

	v := viper.New()
	v.SetConfigType("yaml")
	v.ReadConfig(strings.NewReader(stringYaml))

	config.Registry.SetViper(v)
	config.Registry.Register(map[string]any{
		"database": &database.Config{},
	})
}

func main() {
	seeder.Run(provider.ProvideDB())
}
