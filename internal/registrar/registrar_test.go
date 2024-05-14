package registrar_test

import (
	"testing"

	_ "github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/spf13/viper"

	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/stretchr/testify/suite"
)

type RegisterExecutorTestSuite struct {
	suite.Suite
	booterConfig booter.Config
	viper        viper.Viper
}

func (suite *RegisterExecutorTestSuite) SetupSuite() {
	suite.booterConfig = config.Registry.Get("booter").(booter.Config)
	suite.viper = config.Registry.GetViper()
}

func (suite *RegisterExecutorTestSuite) TestRegisterExecutor() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	service.Registry = service.NewRegistry()

	registrar.RegisterExecutor.AfterExecute()
	registrar.RegisterExecutor.Execute()
	registrar.RegisterExecutor.BeforeExecute()

	suite.NotEmpty(config.Registry.Get("httpServer"))
	suite.NotEmpty(config.Registry.Get("middleware.csrf"))
	suite.NotEmpty(config.Registry.Get("middleware.rateLimit"))
	suite.NotEmpty(config.Registry.Get("authentication.authenticator"))
	suite.NotEmpty(config.Registry.Get("database"))
	suite.NotEmpty(config.Registry.Get("logger.console"))
	suite.NotEmpty(config.Registry.Get("logger.file"))
	suite.NotEmpty(config.Registry.Get("logger.access"))
	suite.NotEmpty(config.Registry.Get("redis"))
	suite.NotEmpty(config.Registry.Get("clickhouse"))

	suite.NotEmpty(service.Registry.Get("authentication.authenticator"))
	suite.NotEmpty(service.Registry.Get("database"))
	suite.NotEmpty(service.Registry.Get("casbinEnforcer"))
	suite.NotEmpty(service.Registry.Get("logger"))
	suite.NotEmpty(service.Registry.Get("logger.console"))
	suite.NotEmpty(service.Registry.Get("logger.file"))
	suite.NotEmpty(service.Registry.Get("logger.access"))
	suite.NotEmpty(service.Registry.Get("redis"))
	suite.NotEmpty(service.Registry.Get("formDecoder"))
	suite.NotEmpty(service.Registry.Get("modifier"))
	suite.NotEmpty(service.Registry.Get("mapstructureDecoder"))
	suite.NotEmpty(service.Registry.Get("clickhouse"))
}

func (suite *RegisterExecutorTestSuite) TestSimpleRegisterExecutor() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	service.Registry = service.NewRegistry()

	registrar.SimpleRegisterExecutor.AfterExecute()
	registrar.SimpleRegisterExecutor.Execute()
	registrar.SimpleRegisterExecutor.BeforeExecute()

	suite.NotEmpty(config.Registry.Get("httpServer"))
	suite.NotEmpty(config.Registry.Get("middleware.csrf"))
	suite.NotEmpty(config.Registry.Get("middleware.rateLimit"))
	suite.NotEmpty(config.Registry.Get("authentication.authenticator"))
	suite.NotEmpty(config.Registry.Get("logger.console"))
	suite.NotEmpty(config.Registry.Get("logger.file"))
	suite.NotEmpty(config.Registry.Get("logger.access"))

	suite.NotEmpty(service.Registry.Get("authentication.authenticator"))
	suite.NotEmpty(service.Registry.Get("logger.console"))
	suite.NotEmpty(service.Registry.Get("logger.file"))
	suite.NotEmpty(service.Registry.Get("logger.access"))
	suite.NotEmpty(service.Registry.Get("formDecoder"))
	suite.NotEmpty(service.Registry.Get("modifier"))
	suite.NotEmpty(service.Registry.Get("mapstructureDecoder"))
}

func (suite *RegisterExecutorTestSuite) TearDownSuite() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	service.Registry = service.NewRegistry()

	registrar.RegisterExecutor.AfterExecute()
	registrar.RegisterExecutor.Execute()
	registrar.RegisterExecutor.BeforeExecute()
}

func TestRegistrarTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterExecutorTestSuite))
}
