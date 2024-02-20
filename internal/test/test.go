package test

import (
	"path"
	"runtime"

	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	booter.Boot(
		loadEnv,
		booter.NewTestingConfig,
		&registrar.RegisterExecutor,
	)

	HttpHandler = NewHttpHandler()
	RdbmsMigration = NewRdbmsMigration()
}

func loadEnv() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime caller cannot get file information")
	}
	wd := path.Join(path.Dir(file), "../..")
	err := godotenv.Load(path.Join(wd, ".env.testing"))
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
}
