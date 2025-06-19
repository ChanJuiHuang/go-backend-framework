package registrar

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/service"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

var SimpleRegisterExecutor = simpleRegisterExecutor{
	registrarCenter: booter.NewRegistrarCenter([]booter.Registrar{
		&LoggerRegistrar{},
		&AuthenticationRegistrar{},
		&MapstructureDecoderRegistrar{},
	}),
}

type simpleRegisterExecutor struct {
	registrarCenter *booter.RegistrarCenter
}

func (*simpleRegisterExecutor) BeforeExecute() {
	config.Registry.RegisterMany(map[string]any{
		"httpServer":           &http.ServerConfig{},
		"middleware.csrf":      &middleware.CsrfConfig{},
		"middleware.rateLimit": &middleware.RateLimitConfig{},
	})
}

func (sre *simpleRegisterExecutor) Execute() {
	sre.registrarCenter.Execute()
}

func (*simpleRegisterExecutor) AfterExecute() {
	service.Registry.SetMany(map[string]any{
		"formDecoder": form.NewDecoder(),
		"modifier":    modifiers.New(),
	})
}
