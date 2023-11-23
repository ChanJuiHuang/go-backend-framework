package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strconv"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

func AdminRegister() {
	userRegisterRequest := user.UserRegisterRequest{
		Name:     "admin",
		Email:    "admin@test.com",
		Password: "abcABC123",
	}
	reqBody, err := json.Marshal(userRegisterRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(reqBody))
	AddCsrfToken(req)
	resp := httptest.NewRecorder()
	HttpHandler.ServeHTTP(resp, req)

	respBody := &user.UserRegisterResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}
}

func AdminLogin() (string, string) {
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

	respBody := &user.UserLoginResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	return respBody.AccessToken, respBody.RefreshToken
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
		{"admin", "/api/admin/grouping-policy/:userId", "GET"},
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
