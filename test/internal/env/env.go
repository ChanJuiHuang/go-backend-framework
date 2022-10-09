package env

import (
	"os"
	"path"
	"runtime"

	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	"github.com/joho/godotenv"
)

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime caller cannot get file information")
	}
	projectRoot := path.Join(path.Dir(file), "../../..")
	envPath := path.Join(projectRoot, ".env")
	godotenv.Load(envPath)

	config.App().ProjectRoot = projectRoot
	config.Database().Driver = config.Sqlite
	config.Database().Database = ":memory:"
	config.Log().Level = config.Error
	config.Redis().DB = 1
	os.Setenv("GIN_MODE", "release")
}
