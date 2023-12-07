package registrar_test

import (
	"testing"

	_ "github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/spf13/viper"

	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/stretchr/testify/assert"
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

	assert.NotEmpty(suite.T(), config.Registry.Get("httpServer"))
	assert.NotEmpty(suite.T(), config.Registry.Get("middleware.csrf"))
	assert.NotEmpty(suite.T(), config.Registry.Get("middleware.rateLimit"))
	assert.NotEmpty(suite.T(), config.Registry.Get("authentication.authenticator"))
	assert.NotEmpty(suite.T(), config.Registry.Get("database"))
	assert.NotEmpty(suite.T(), config.Registry.Get("logger.console"))
	assert.NotEmpty(suite.T(), config.Registry.Get("logger.file"))
	assert.NotEmpty(suite.T(), config.Registry.Get("logger.access"))
	assert.NotEmpty(suite.T(), config.Registry.Get("redis"))

	assert.NotEmpty(suite.T(), service.Registry.Get("authentication.authenticator"))
	assert.NotEmpty(suite.T(), service.Registry.Get("database"))
	assert.NotEmpty(suite.T(), service.Registry.Get("casbinEnforcer"))
	assert.NotEmpty(suite.T(), service.Registry.Get("logger"))
	assert.NotEmpty(suite.T(), service.Registry.Get("logger.console"))
	assert.NotEmpty(suite.T(), service.Registry.Get("logger.file"))
	assert.NotEmpty(suite.T(), service.Registry.Get("logger.access"))
	assert.NotEmpty(suite.T(), service.Registry.Get("redis"))
	assert.NotEmpty(suite.T(), service.Registry.Get("formDecoder"))
	assert.NotEmpty(suite.T(), service.Registry.Get("modifier"))
}

func (suite *RegisterExecutorTestSuite) TestSimpleRegisterExecutor() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	service.Registry = service.NewRegistry()

	registrar.SimpleRegisterExecutor.AfterExecute()
	registrar.SimpleRegisterExecutor.Execute()
	registrar.SimpleRegisterExecutor.BeforeExecute()

	assert.NotEmpty(suite.T(), config.Registry.Get("httpServer"))
	assert.NotEmpty(suite.T(), config.Registry.Get("middleware.csrf"))
	assert.NotEmpty(suite.T(), config.Registry.Get("middleware.rateLimit"))
	assert.NotEmpty(suite.T(), config.Registry.Get("authentication.authenticator"))
	assert.NotEmpty(suite.T(), config.Registry.Get("logger.console"))
	assert.NotEmpty(suite.T(), config.Registry.Get("logger.file"))
	assert.NotEmpty(suite.T(), config.Registry.Get("logger.access"))

	assert.NotEmpty(suite.T(), service.Registry.Get("authentication.authenticator"))
	assert.NotEmpty(suite.T(), service.Registry.Get("logger.console"))
	assert.NotEmpty(suite.T(), service.Registry.Get("logger.file"))
	assert.NotEmpty(suite.T(), service.Registry.Get("logger.access"))
	assert.NotEmpty(suite.T(), service.Registry.Get("formDecoder"))
	assert.NotEmpty(suite.T(), service.Registry.Get("modifier"))
}

func (suite *RegisterExecutorTestSuite) TearDownSuite() {
	config.Registry = config.NewRegistry(&suite.viper)
	config.Registry.Set("booter", &suite.booterConfig)
	service.Registry = service.NewRegistry()

	registrar.RegisterExecutor.AfterExecute()
	registrar.RegisterExecutor.Execute()
	registrar.RegisterExecutor.BeforeExecute()
}

func TestUserUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterExecutorTestSuite))
}
