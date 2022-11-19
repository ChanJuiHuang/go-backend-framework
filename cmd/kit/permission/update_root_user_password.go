package permission

import (
	"fmt"

	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
)

func UpdateRootUserPassword() {
	password := util.RandomString(16)
	db := provider.App.DB.Model(&model.EmailUser{}).
		Where("email = ?", "root@root.com").
		Update("password", util.MakeArgon2IdHash(password))
	if db.Error != nil {
		panic(db.Error)
	}
	if db.RowsAffected != 1 {
		fmt.Println("update password failed")
		return
	}
	fmt.Printf("root user password is [%s]\n", password)
}
