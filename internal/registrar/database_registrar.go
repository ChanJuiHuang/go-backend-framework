package registrar

import (
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/chan-jui-huang/go-backend-package/pkg/database"
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
