package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
)

type DatabaseRegistrar struct {
	config database.Config
}

func (dr *DatabaseRegistrar) Boot() {
	config.Registry.Register("database", &dr.config)
}

func (dr *DatabaseRegistrar) Register() {
	service.Registry.Set("database", database.New(dr.config))
}
