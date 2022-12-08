package util

import (
	"errors"

	"gorm.io/gorm"
)

func IsDatabaseError(db *gorm.DB) bool {
	return db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound)
}
