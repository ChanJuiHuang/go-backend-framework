package init

import (
	"os"

	"github.com/ChanJuiHuang/go-backend-framework/app/config"
)

func init() {
	err := os.Setenv("TZ", config.App().Timezone)
	if err != nil {
		panic(err)
	}
}
