package provider

import (
	"path"

	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/logger"
	pkgRedis "github.com/ChanJuiHuang/go-backend-framework/pkg/redis"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ProvideLogger() (*zap.Logger, *zap.Logger, *zap.Logger) {
	consoleLogger := logger.NewConsoleLogger(
		config.Registry.Get("logger.console").(logger.ConsoleConfig),
		logger.ConsoleEncoder,
		logger.DefaultZapOptions...,
	)

	globalConfig := config.Registry.Get("global").(global.Config)
	fileConfig := config.Registry.Get("logger.file").(logger.FileConfig)
	fileConfig.LogPath = path.Join(globalConfig.RootDir, fileConfig.LogPath)
	fileLogger, err := logger.NewFileLogger(
		fileConfig,
		logger.JsonEncoder,
		logger.DefaultZapOptions...,
	)
	if err != nil {
		panic(err)
	}

	v := config.Registry.GetViper()
	var defaultLogger *zap.Logger
	switch logger.Type(v.GetString("logger.type")) {
	case logger.Console:
		defaultLogger = consoleLogger
	case "file":
		defaultLogger = fileLogger
	default:
		defaultLogger = consoleLogger
	}

	return defaultLogger, consoleLogger, fileLogger
}

func ProvideDB() *gorm.DB {
	return database.New(config.Registry.Get("database").(database.Config))
}

func ProvideRedis() *redis.Client {
	return pkgRedis.New(config.Registry.Get("redis").(pkgRedis.Config))
}

func ProvideAuthenticator() *authentication.Authenticator {
	authenticator, err := authentication.NewAuthenticator(config.Registry.Get("authentication.authenticator").(authentication.Config))
	if err != nil {
		panic(err)
	}

	return authenticator
}

func ProvideCasbinEnforcer(db *gorm.DB) *casbin.SyncedCachedEnforcer {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rules")
	if err != nil {
		panic(err)
	}

	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
`)
	if err != nil {
		panic(err)
	}

	enforcer, err := casbin.NewSyncedCachedEnforcer(m, adapter)
	if err != nil {
		panic(err)
	}

	return enforcer
}
