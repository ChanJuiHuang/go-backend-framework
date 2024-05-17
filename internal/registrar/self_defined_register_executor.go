package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/http"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

type SelfDefinedRegisterExecutor struct {
	registrarCenter *booter.RegistrarCenter
}

func NewSelfDefinedRegisterExecutor(registrarCenter *booter.RegistrarCenter) *SelfDefinedRegisterExecutor {
	return &SelfDefinedRegisterExecutor{
		registrarCenter: registrarCenter,
	}
}

func (*SelfDefinedRegisterExecutor) BeforeExecute() {
	config.Registry.RegisterMany(map[string]any{
		"httpServer":           &http.ServerConfig{},
		"middleware.csrf":      &middleware.CsrfConfig{},
		"middleware.rateLimit": &middleware.RateLimitConfig{},
	})
}

func (executor *SelfDefinedRegisterExecutor) Execute() {
	executor.registrarCenter.Execute()
}

func (*SelfDefinedRegisterExecutor) AfterExecute() {
	service.Registry.SetMany(map[string]any{
		"formDecoder": form.NewDecoder(),
		"modifier":    modifiers.New(),
	})
}
