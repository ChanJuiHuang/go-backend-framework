package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

type appConfig struct {
	ProjectRoot string `ignored:"true"`
	Timezone    string `required:"true"`
	Locale      string `required:"true"`
	Debug       bool   `required:"true"`
}

var appCfg *appConfig

func App() *appConfig {
	if appCfg == nil {
		appCfg = new(appConfig)
		err := envconfig.Process("app", appCfg)
		if err != nil {
			panic(err)
		}

		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		appCfg.ProjectRoot = wd
	}

	return appCfg
}
