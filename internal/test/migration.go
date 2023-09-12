package test

import (
	"path"

	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
	"github.com/pressly/goose/v3"
)

type migration struct {
	dir string
}

var Migration *migration

func NewMigration() *migration {
	globalConfig := config.Registry.Get("global").(global.Config)

	return &migration{
		dir: path.Join(globalConfig.RootDir, "internal/migration/test"),
	}
}

func (dt *migration) Run(callbacks ...func()) {
	databaseConfig := config.Registry.Get("database").(database.Config)
	db, err := provider.Registry.DB().DB()
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
	db, err := provider.Registry.DB().DB()
	if err != nil {
		panic(err)
	}
	if err := goose.Reset(db, dt.dir); err != nil {
		panic(err)
	}
}
