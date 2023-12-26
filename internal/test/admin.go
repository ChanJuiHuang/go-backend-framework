package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strconv"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/argon2"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

func AdminRegister() {
	database := service.Registry.Get("database").(*gorm.DB)
	db := database.Create(&model.User{
		Name:     "admin",
		Email:    "admin@test.com",
		Password: argon2.MakeArgon2IdHash("abcABC123"),
	})
	if err := db.Error; err != nil {
		panic(err)
	}
}

func AdminLogin() string {
	userLoginRequest := user.UserLoginRequest{
		Email:    "admin@test.com",
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

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}
	data := &user.UserLoginData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	return data.AccessToken

}

func AdminAddPolicies() {
	policies := [][]string{
		{"admin", "/api/admin/policy", "POST"},
		{"admin", "/api/admin/policy", "DELETE"},
		{"admin", "/api/admin/policy/subject", "GET"},
		{"admin", "/api/admin/policy/subject/:subject", "GET"},
		{"admin", "/api/admin/policy/subject", "DELETE"},
		{"admin", "/api/admin/policy/reload", "POST"},
		{"admin", "/api/admin/grouping-policy", "POST"},
		{"admin", "/api/admin/user/:userId/grouping-policy", "GET"},
		{"admin", "/api/admin/grouping-policy", "DELETE"},
	}
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	result, err := enforcer.AddPolicies(policies)
	if err != nil {
		panic(err)
	}
	if !result {
		panic("add casbin testing policies failed")
	}
}

func AdminAddRole() {
	user := &model.User{}
	database := service.Registry.Get("database").(*gorm.DB)
	db := database.Where("email = ?", "admin@test.com").
		First(user)
	if err := db.Error; err != nil {
		panic(err)
	}

	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	result, err := enforcer.AddRoleForUser(strconv.Itoa(int(user.Id)), "admin")
	if err != nil {
		panic(err)
	}
	if !result {
		panic("add casbin testing policies failed")
	}
}
