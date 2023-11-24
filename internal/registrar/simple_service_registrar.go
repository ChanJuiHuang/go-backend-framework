package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

var SimpleServiceRegistrar simpleServiceRegistrar

type simpleServiceRegistrar struct{}

func (*simpleServiceRegistrar) Register() {
	service.Registry.SetMany(map[string]any{
		"formDecoder": form.NewDecoder(),
		"modifier":    modifiers.New(),
	})

	registrars := []booter.ServiceRegistrar{
		&LoggerRegistrar{},
		&AuthenticationRegistrar{},
		&CasbinRegistrar{},
	}

	for _, registrar := range registrars {
		registrar.Register()
	}
}
