package test

import (
	"os"
	"path"
	"runtime"

	internalConfig "github.com/ChanJuiHuang/go-backend-framework/internal/config"
	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	internalProvider "github.com/ChanJuiHuang/go-backend-framework/internal/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime caller cannot get file information")
	}
	wd := path.Join(path.Dir(file), "../..")

	globalConfig := newGlobalConfig(wd)
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	internalConfig.RegisterConfigWithFile(*globalConfig, "config.testing.yml")
	internalProvider.RegisterService()

	HttpHandler = NewHttpHandler()
	Migration = NewMigration()
}

func newGlobalConfig(rootDir string) *global.Config {
	return &global.Config{
		RootDir:  rootDir,
		Timezone: "UTC",
		Debug:    false,
		Testing:  true,
	}
}

func registerGlobalConfig(globalConfig *global.Config) {
	config.Registry.Set("global", globalConfig)
}

func setEnv(globalConfig global.Config) {
	err := godotenv.Load(path.Join(globalConfig.RootDir, ".env.testing"))
	if err != nil {
		panic(err)
	}

	err = os.Setenv("TZ", globalConfig.Timezone)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
}
