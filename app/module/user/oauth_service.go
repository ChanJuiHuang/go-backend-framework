package user

import (
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"gorm.io/gorm"
)

var GetUserId func(oauthProvider util.OauthInterface) (uint, error) = func(oauthProvider util.OauthInterface) (uint, error) {
	var userId uint
	var err error
	switch oauthProvider.(type) {
	case *util.GoogleOauth:
		userId, err = getGoogleUserId(oauthProvider.(*util.GoogleOauth))
	default:
		panic("oauth provider does not exist")
	}

	return userId, err
}

var CreateOauthUser = func(oauthProvider util.OauthInterface) (uint, error) {
	var userId uint
	var err error
	switch oauthProvider.(type) {
	case *util.GoogleOauth:
		userId, err = createGoogleUser(oauthProvider.(*util.GoogleOauth))
	default:
		panic("oauth provider does not exist")
	}

	return userId, err
}

func getGoogleUserId(googleOauth *util.GoogleOauth) (uint, error) {
	googleUser := new(model.GoogleUser)
	db := provider.App.DB.Where("google_id = ?", googleOauth.Id).First(googleUser)

	if util.IsDatabaseError(db) {
		return 0, util.WrapError(db.Error)
	}

	return googleUser.UserId, nil
}

func createGoogleUser(googleOauth *util.GoogleOauth) (uint, error) {
	user := new(model.User)
	err := provider.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return util.WrapError(err)
		}
		db := tx.Create(&model.EmailUser{
			Email:    googleOauth.Email,
			Password: util.RandomString(16),
			UserId:   user.Id,
		})
		if db.Error != nil {
			return util.WrapError(db.Error)
		}

		db = tx.Create(&model.GoogleUser{
			GoogleId: googleOauth.Id,
			UserId:   user.Id,
		})
		if db.Error != nil {
			return util.WrapError(db.Error)
		}

		return nil
	})

	return user.Id, err
}
