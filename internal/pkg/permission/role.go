package permission

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateRole(tx *gorm.DB, value any) error {
	if err := tx.Table("roles").Create(value).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func CreateRolePermission(tx *gorm.DB, value any) error {
	if err := tx.Table("role_permissions").Create(value).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetRole(tx *gorm.DB, query any, args ...any) (*model.Role, error) {
	role := &model.Role{}
	err := tx.Table("roles").
		Where(query, args...).
		First(role).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return role, nil
}
