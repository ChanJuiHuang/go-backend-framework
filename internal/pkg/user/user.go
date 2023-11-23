package user

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func Create(user any) error {
	database := service.Registry.Get("database").(*gorm.DB)
	if err := database.Create(user).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Get(query any, args ...any) (*model.User, error) {
	user := &model.User{}
	database := service.Registry.Get("database").(*gorm.DB)
	err := database.Where(query, args...).
		First(user).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}

func Update(id uint, values map[string]any) (int, error) {
	database := service.Registry.Get("database").(*gorm.DB)
	db := database.Table("users").
		Where("id = ?", id).
		Updates(values)
	if err := db.Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return int(db.RowsAffected), nil
}
