package provider

import (
	"path"

	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func provideCasbin(db *gorm.DB) *casbin.SyncedEnforcer {
	casbinCfg := config.Casbin()
	if !casbinCfg.Enabled {
		return new(casbin.SyncedEnforcer)
	}

	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rules")
	if err != nil {
		panic(err)
	}

	modelPath := path.Join(config.App().ProjectRoot, casbinCfg.ModelPath)
	enforcer, err := casbin.NewSyncedEnforcer(modelPath, adapter)
	if err != nil {
		panic(err)
	}

	return enforcer
}
