package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/clickhouse"
)

type ClickhouseRegistrar struct {
	config clickhouse.Config
}

func (cr *ClickhouseRegistrar) Boot() {
	config.Registry.Register("clickhouse", &cr.config)
}

func (cr *ClickhouseRegistrar) Register() {
	conn, err := clickhouse.New(cr.config)
	if err != nil {
		panic(err)
	}

	service.Registry.Set("clickhouse", conn)
}
