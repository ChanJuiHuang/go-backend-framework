package registrar

import (
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/chan-jui-huang/go-backend-package/pkg/clickhouse"
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
