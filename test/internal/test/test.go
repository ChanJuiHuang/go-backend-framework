package test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"

	_ "github.com/ChanJuiHuang/go-backend-framework/test/internal/env"
	"github.com/pressly/goose/v3"

	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	frameworkHttp "github.com/ChanJuiHuang/go-backend-framework/app/http"
	"github.com/ChanJuiHuang/go-backend-framework/app/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/cmd/kit/permission"
	"github.com/ChanJuiHuang/go-backend-framework/database/seeder"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	Engine              *gin.Engine
	db                  *sql.DB
	TestingMigrationDir string
	AccessToken         string
	RefreshToken        string
	RootAccessToken     string
}

func (suite *TestSuite) SetupSuite() {
	suite.TestingMigrationDir = fmt.Sprintf("%s/%s", config.App().ProjectRoot, "database/testing-migration")
	db, err := provider.App.DB.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect(string(config.Database().Driver)); err != nil {
		panic(err)
	}
	if err := goose.Up(db, suite.TestingMigrationDir); err != nil {
		panic(err)
	}
	seeder.Run()

	suite.Engine = frameworkHttp.GetEngine()
	suite.db = db
	suite.RootAccessToken = permission.GenerateRootAccessToken()
	suite.registerUser()
	provider.App.ImportCasbinPolicies()
}

func (suite *TestSuite) TearDownSuite() {
	if err := goose.Reset(suite.db, suite.TestingMigrationDir); err != nil {
		panic(err)
	}

	_, err := provider.App.Redis.FlushDB(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}

func (suite *TestSuite) AddCsrfToken(req *http.Request) {
	cookie := &http.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    "1234567890",
		Path:     "/",
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	req.AddCookie(cookie)
	req.Header.Add("X-XSRF-TOKEN", "1234567890")
}

func (suite *TestSuite) AddAuthorizationHeader(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+suite.AccessToken)
}

func (suite *TestSuite) AddRootAccessToken(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+suite.RootAccessToken)
}

func (suite *TestSuite) registerUser() {
	requestBody, err := provider.App.Json.Marshal(user.UserCreateRequest{
		Name:     "go",
		Email:    "go@gogogo.com",
		Password: "abcABC123",
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/api/user", bytes.NewReader(requestBody))
	suite.AddCsrfToken(req)
	res := httptest.NewRecorder()
	suite.Engine.ServeHTTP(res, req)
	var resBody map[string]any
	provider.App.Json.Unmarshal(res.Body.Bytes(), &resBody)
	suite.AccessToken = resBody["access_token"].(string)
	suite.RefreshToken = resBody["refresh_token"].(string)
}
