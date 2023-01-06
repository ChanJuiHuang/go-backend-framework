package seeder

import (
	"fmt"

	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"gorm.io/gorm"
)

func runUserSeeder(tx *gorm.DB) error {
	password := util.RandomString(16)
	user := &model.User{Name: "root"}
	if err := tx.Create(user).Error; err != nil {
		return err
	}

	err := tx.Create(&model.EmailUser{
		Email:    "root@root.com",
		Password: util.MakeArgon2IdHash(password),
		UserId:   user.Id,
	}).Error
	if err != nil {
		return err
	}

	fmt.Printf("root user password is [%s]\n", password)

	return nil
}
