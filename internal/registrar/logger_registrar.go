package registrar

import (
	"fmt"
	"path"

	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/logger"
)

type LoggerRegistrar struct{}

func (*LoggerRegistrar) Register() {
	consoleLogger, err := logger.NewLogger(
		config.Registry.Get("logger.console").(logger.Config),
		logger.ConsoleEncoder,
		logger.DefaultZapOptions...,
	)
	if err != nil {
		panic(err)
	}

	booterConfig := config.Registry.Get("booter").(booter.Config)
	fileConfig := config.Registry.Get("logger.file").(logger.Config)
	fileConfig.LogPath = path.Join(booterConfig.RootDir, fileConfig.LogPath)
	fileLogger, err := logger.NewLogger(
		fileConfig,
		logger.JsonEncoder,
		logger.DefaultZapOptions...,
	)
	if err != nil {
		panic(err)
	}

	accessConfig := config.Registry.Get("logger.access").(logger.Config)
	accessConfig.LogPath = path.Join(booterConfig.RootDir, accessConfig.LogPath)
	accessLogger, err := logger.NewLogger(
		accessConfig,
		logger.JsonEncoder,
	)
	if err != nil {
		panic(err)
	}
	service.Registry.SetMany(map[string]any{
		"consoleLogger": consoleLogger,
		"fileLogger":    fileLogger,
		"accessLogger":  accessLogger,
	})

	v := config.Registry.GetViper()
	settings := v.Sub("logger").AllSettings()
	defaultSetting := v.GetString("logger.default")

	for setting := range settings {
		if defaultSetting == setting {
			service.Registry.Set(
				"logger",
				service.Registry.Get(fmt.Sprintf("%sLogger", defaultSetting)),
			)
		}
	}
}
