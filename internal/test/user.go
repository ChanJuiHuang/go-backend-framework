package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/argon2"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"gorm.io/gorm"
)

func UserRegister() {
	database := service.Registry.Get("database").(*gorm.DB)
	db := database.Create(&model.User{
		Name:     "john",
		Email:    "john@test.com",
		Password: argon2.MakeArgon2IdHash("abcABC123"),
	})
	if err := db.Error; err != nil {
		panic(err)
	}
}

func UserLogin() (string, string) {
	userLoginRequest := user.UserLoginRequest{
		Email:    "john@test.com",
		Password: "abcABC123",
	}
	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBody))
	AddCsrfToken(req)
	resp := httptest.NewRecorder()
	HttpHandler.ServeHTTP(resp, req)

	respBody := &user.UserLoginResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	return respBody.AccessToken, respBody.RefreshToken
}
