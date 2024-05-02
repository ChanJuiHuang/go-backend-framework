package casbinrule

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func Create(tx *gorm.DB, value any) error {
	if err := tx.Table("casbin_rules").Create(value).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Delete(tx *gorm.DB, query any, args ...any) error {
	err := tx.Table("casbin_rules").
		Where(query, args...).
		Delete(&struct{}{}).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
