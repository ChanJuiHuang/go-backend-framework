package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
)

type AuthenticationRegistrar struct {
	config authentication.Config
}

func (ar *AuthenticationRegistrar) Boot() {
	config.Registry.Register("authentication.authenticator", &ar.config)
}

func (ar *AuthenticationRegistrar) Register() {
	authenticator, err := authentication.NewAuthenticator(ar.config)
	if err != nil {
		panic(err)
	}

	service.Registry.Set("authentication.authenticator", authenticator)
}
