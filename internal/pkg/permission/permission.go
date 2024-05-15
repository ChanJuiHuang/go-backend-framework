package permission

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func Create(tx *gorm.DB, value any) error {
	if err := tx.Table("permissions").Create(value).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Get(tx *gorm.DB, query any, args ...any) (*model.Permission, error) {
	permission := &model.Permission{}
	err := tx.Table("permissions").
		Where(query, args...).
		First(permission).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return permission, nil
}
