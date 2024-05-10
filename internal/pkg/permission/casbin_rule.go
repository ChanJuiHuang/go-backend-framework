package permission

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateCasbinRule(tx *gorm.DB, value any) error {
	if err := tx.Table("casbin_rules").Create(value).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func DeleteCasbinRule(tx *gorm.DB, query any, args ...any) error {
	err := tx.Table("casbin_rules").
		Where(query, args...).
		Delete(&struct{}{}).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
