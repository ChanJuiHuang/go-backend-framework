package permission

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func Create(tx *gorm.DB, value any) error {
	if err := tx.Table("permissions").Create(value).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
