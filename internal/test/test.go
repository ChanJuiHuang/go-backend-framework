package test

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime caller cannot get file information")
	}
	wd := path.Join(path.Dir(file), "../..")
	env := "dev"
	if e := os.Getenv("ENV"); e != "" {
		env = e
	}
	envFile := fmt.Sprintf(".env.%s", env)
	configFile := fmt.Sprintf("config.%s.yml", env)

	booter.Boot(
		loadEnv(wd, envFile),
		NewConfig(wd, configFile),
		&registrar.RegisterExecutor,
	)

	HttpHandler = NewHttpHandler()
	RdbmsMigration = NewRdbmsMigration()
	ClickhouseMigration = NewClickhouseMigration()
	PermissionService = NewPermissionService()
	UserService = NewUserService()
	AdminService = NewAdminService()
}

func loadEnv(wd string, envFile string) func() {
	return func() {
		err := godotenv.Load(path.Join(wd, envFile))
		if err != nil {
			panic(err)
		}

		gin.SetMode(gin.ReleaseMode)
	}
}

func NewConfig(wd string, envFile string) func() *booter.Config {
	return func() *booter.Config {
		return booter.NewConfig(wd, envFile, false, true)
	}
}
