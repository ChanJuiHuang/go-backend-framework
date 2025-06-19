package registrar

import (
	"fmt"
	"path"

	"github.com/ChanJuiHuang/go-backend-framework/v2/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/v2/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/v2/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/v2/pkg/logger"
)

type LoggerRegistrar struct {
	consoleConfig logger.Config
	fileConfig    logger.Config
	accessConfig  logger.Config
}

func (lr *LoggerRegistrar) Boot() {
	config.Registry.RegisterMany(map[string]any{
		"logger.console": &logger.Config{},
		"logger.file":    &logger.Config{},
		"logger.access":  &logger.Config{},
	})

	lr.consoleConfig = config.Registry.Get("logger.console").(logger.Config)
	lr.fileConfig = config.Registry.Get("logger.file").(logger.Config)
	lr.accessConfig = config.Registry.Get("logger.access").(logger.Config)
}

func (lr *LoggerRegistrar) Register() {
	consoleLogger, err := logger.NewLogger(
		lr.consoleConfig,
		logger.ConsoleEncoder,
		logger.DefaultZapOptions...,
	)
	if err != nil {
		panic(err)
	}

	booterConfig := config.Registry.Get("booter").(booter.Config)
	lr.fileConfig.LogPath = path.Join(booterConfig.RootDir, lr.fileConfig.LogPath)
	fileLogger, err := logger.NewLogger(
		lr.fileConfig,
		logger.JsonEncoder,
		logger.DefaultZapOptions...,
	)
	if err != nil {
		panic(err)
	}

	lr.accessConfig.LogPath = path.Join(booterConfig.RootDir, lr.accessConfig.LogPath)
	accessLogger, err := logger.NewLogger(
		lr.accessConfig,
		logger.JsonEncoder,
	)
	if err != nil {
		panic(err)
	}
	service.Registry.SetMany(map[string]any{
		"logger.console": consoleLogger,
		"logger.file":    fileLogger,
		"logger.access":  accessLogger,
	})

	v := config.Registry.GetViper()
	settings := v.Sub("logger").AllSettings()
	defaultSetting := v.GetString("logger.default")

	for setting := range settings {
		if defaultSetting == setting {
			service.Registry.Set(
				"logger",
				service.Registry.Get(fmt.Sprintf("logger.%s", defaultSetting)),
			)
		}
	}
}
