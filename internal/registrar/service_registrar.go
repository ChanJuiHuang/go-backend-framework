package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

var ServiceRegistrar serviceRegistrar

type serviceRegistrar struct{}

func (*serviceRegistrar) Register() {
	service.Registry.SetMany(map[string]any{
		"formDecoder": form.NewDecoder(),
		"modifier":    modifiers.New(),
	})

	registrars := []booter.ServiceRegistrar{
		&LoggerRegistrar{},
		&DatabaseRegistrar{},
		&RedisRegistrar{},
		&AuthenticationRegistrar{},
		&CasbinRegistrar{},
	}

	for _, registrar := range registrars {
		registrar.Register()
	}
}
