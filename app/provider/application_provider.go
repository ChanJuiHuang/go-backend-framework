package provider

import (
	"strconv"

	_ "github.com/ChanJuiHuang/go-backend-framework/app/provider/internal/init"

	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	Logger      *zap.Logger
	DB          *gorm.DB
	Json        jsoniter.API
	Redis       *redis.Client
	Validator   *validator.Validate
	FormDecoder *form.Decoder
	Modifier    *mold.Transformer
	Casbin      *casbin.SyncedEnforcer
}

func provideApplication(
	logger *zap.Logger,
	db *gorm.DB,
	json jsoniter.API,
	redis *redis.Client,
	validator *validator.Validate,
	formDecoder *form.Decoder,
	modifier *mold.Transformer,
	casbin *casbin.SyncedEnforcer,
) *Application {
	return &Application{
		Logger:      logger,
		DB:          db,
		Json:        json,
		Redis:       redis,
		Validator:   validator,
		FormDecoder: formDecoder,
		Modifier:    modifier,
		Casbin:      casbin,
	}
}

func (app *Application) ImportCasbinPolicies() {
	casbin := app.Casbin

	_, err := casbin.AddPolicies([][]string{
		{"root", "/scheduler/welcome", "GET"},
		{"root", "/scheduler/refresh-token-record", "DELETE"},
	})
	if err != nil {
		panic(err)
	}

	emailUser := new(model.EmailUser)
	db := app.DB.Where("email = ?", "root@root.com").
		First(emailUser)
	if db.Error != nil {
		panic(err)
	}

	_, err = casbin.AddGroupingPolicy(strconv.FormatUint(uint64(emailUser.UserId), 10), "root")
	if err != nil {
		panic(err)
	}
}
