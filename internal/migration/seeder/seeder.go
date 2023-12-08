package seeder

import (
	"fmt"

	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"gorm.io/gorm"
)

type runSeederFunc func(tx *gorm.DB) error

type SeederExecutor struct {
	runSeederFuncs map[string]runSeederFunc
}

func NewSeederExecutor() *SeederExecutor {
	runSeederFuncs := map[string]runSeederFunc{
		"user": runUserSeeder,
	}

	return &SeederExecutor{
		runSeederFuncs: runSeederFuncs,
	}
}

func (se *SeederExecutor) ShowSeeders() {
	for seeder := range se.runSeederFuncs {
		fmt.Println(seeder)
	}
}

func (se *SeederExecutor) Run(seeders []string) {
	if len(seeders) == 1 && seeders[0] == "" {
		seeders = make([]string, 0, len(se.runSeederFuncs))
		for seeder := range se.runSeederFuncs {
			seeders = append(seeders, seeder)
		}
	}

	database := service.Registry.Get("database").(*gorm.DB)
	err := database.Transaction(func(tx *gorm.DB) error {
		for _, seeder := range seeders {
			if fn, ok := se.runSeederFuncs[seeder]; ok {
				if err := fn(tx); err != nil {
					fmt.Printf("[%s] execute failed\n", seeder)
					return err
				}
			} else {
				fmt.Printf("[%s] does not exist\n", seeder)
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
