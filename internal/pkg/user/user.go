package user

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user/model"
	"github.com/pkg/errors"
)

func Create(user any) error {
	if err := provider.Registry.DB().Create(user).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Get(query any, args ...any) (*model.User, error) {
	user := &model.User{}
	err := provider.Registry.DB().
		Where(query, args...).
		First(user).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}

func Update(id uint, values map[string]any) (int, error) {
	db := provider.Registry.DB().
		Table("users").
		Where("id = ?", id).
		Updates(values)
	if err := db.Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return int(db.RowsAffected), nil
}
