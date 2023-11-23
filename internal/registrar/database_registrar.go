package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
)

type DatabaseRegistrar struct{}

func (*DatabaseRegistrar) Register() {
	service.Registry.Set("database", database.New(config.Registry.Get("database").(database.Config)))
}
