package seeder

import (
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"gorm.io/gorm"
)

func Run() {
	err := provider.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := runUserSeeder(tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
