package seeder

import (
	"fmt"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/pkg/argon2"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/random"
	"gorm.io/gorm"
)

func runUserSeeder(tx *gorm.DB) error {
	password := random.RandomString(16)
	err := tx.Table("users").Create(&struct {
		Name      string
		Password  string
		Email     string
		CreatedAt time.Time
		UpdatedAt time.Time
	}{
		Name:     "admin",
		Password: argon2.MakeArgon2IdHash(password),
		Email:    "admin@admin.com",
	}).Error
	if err != nil {
		return err
	}

	fmt.Printf("admin user password is [%s]\n", password)

	return nil
}
