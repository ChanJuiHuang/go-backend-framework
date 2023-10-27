package seeder

import (
	"gorm.io/gorm"
)

type runSeederFunc func(tx *gorm.DB) error

func Run(db *gorm.DB, seeders []string) {
	runSeederFuncs := map[string]runSeederFunc{
		"user": runUserSeeder,
	}
	if len(seeders) == 1 && seeders[0] == "" {
		seeders = make([]string, 0, len(runSeederFuncs))
		for seeder := range runSeederFuncs {
			seeders = append(seeders, seeder)
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, seeder := range seeders {
			if err := runSeederFuncs[seeder](tx); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
