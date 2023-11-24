package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
)

type AuthenticationRegistrar struct{}

func (*AuthenticationRegistrar) Register() {
	authenticator, err := authentication.NewAuthenticator(config.Registry.Get("authentication.authenticator").(authentication.Config))
	if err != nil {
		panic(err)
	}

	service.Registry.Set("authentication.authenticator", authenticator)
}
