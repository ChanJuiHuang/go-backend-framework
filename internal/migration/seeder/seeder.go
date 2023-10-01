package seeder

import (
	"gorm.io/gorm"
)

type runSeederFunc func(tx *gorm.DB) error

func Run(db *gorm.DB) {
	runSeederFuncs := []runSeederFunc{
		runUserSeeder,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, runSeeder := range runSeederFuncs {
			if err := runSeeder(tx); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
