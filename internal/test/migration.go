package test

import (
	"path"

	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
	"github.com/casbin/casbin/v2"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

type migration struct {
	dir string
}

var Migration *migration

func NewMigration() *migration {
	booterConfig := config.Registry.Get("booter").(booter.Config)

	return &migration{
		dir: path.Join(booterConfig.RootDir, "internal/migration/mysql/test"),
	}
}

func (dt *migration) Run(callbacks ...func()) {
	databaseConfig := config.Registry.Get("database").(database.Config)
	database := service.Registry.Get("database").(*gorm.DB)
	db, err := database.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect(string(databaseConfig.Driver)); err != nil {
		panic(err)
	}
	if err := goose.Up(db, dt.dir); err != nil {
		panic(err)
	}

	for _, callback := range callbacks {
		callback()
	}
}

func (dt *migration) Reset() {
	database := service.Registry.Get("database").(*gorm.DB)
	db, err := database.DB()
	if err != nil {
		panic(err)
	}
	if err := goose.Reset(db, dt.dir); err != nil {
		panic(err)
	}

	err = database.Exec("DELETE FROM casbin_rules").Error
	if err != nil {
		panic(err)
	}

	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	if err := enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}
