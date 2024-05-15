package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strconv"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	casbinrule "github.com/ChanJuiHuang/go-backend-framework/internal/pkg/casbin_rule"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/argon2"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
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
	data := &user.TokenData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	return data.AccessToken

}

func AdminAddPolicies() {
	policies := []gormadapter.CasbinRule{
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy", V2: "POST"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy", V2: "DELETE"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/subject", V2: "GET"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/subject/:subject", V2: "GET"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/subject/:subject/user", V2: "GET"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/subject", V2: "DELETE"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/policy/reload", V2: "POST"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/grouping-policy", V2: "POST"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/user/:userId/grouping-policy", V2: "GET"},
		{Ptype: "p", V0: "admin", V1: "/api/admin/grouping-policy", V2: "DELETE"},
	}
	if err := casbinrule.Create(database.NewTx(), policies); err != nil {
		panic(err)
	}

	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	if err := enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}

func AddRoleToUser(userId uint, role string) {
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	result, err := enforcer.AddRoleForUser(strconv.Itoa(int(userId)), role)
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

	AddRoleToUser(user.Id, "admin")
}
