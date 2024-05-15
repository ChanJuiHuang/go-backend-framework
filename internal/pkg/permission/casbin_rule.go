package permission

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateCasbinRule(tx *gorm.DB, value any) error {
	if err := tx.Table("casbin_rules").Create(value).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetCasbinRules(tx *gorm.DB, query any, args ...any) ([]gormadapter.CasbinRule, error) {
	casbinRules := []gormadapter.CasbinRule{}
	err := tx.Table("casbin_rules").
		Where(query, args...).
		Find(&casbinRules).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return casbinRules, nil
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
