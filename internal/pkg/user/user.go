package user

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/pkg/errors"
)

func Create(user any) error {
	if err := provider.Registry.DB().Create(user).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
