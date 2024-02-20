package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
	"github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserGetPolicyTestSuite struct {
	suite.Suite
}

func (suite *UserGetPolicyTestSuite) SetupTest() {
	test.RdbmsMigration.Run()
	test.UserRegister()
}

func (suite *UserGetPolicyTestSuite) TestGetPolicy() {
	database := service.Registry.Get("database").(*gorm.DB)
	u := &model.User{}
	db := database.Where("email = ?", "john@test.com").
		First(u)
	if err := db.Error; err != nil {
		panic(err)
	}

	test.AdminAddPolicies()
	test.AddRoleToUser(u.Id, "admin")

	accessToken := test.UserLogin()
	req := httptest.NewRequest("GET", "/api/user/policy", nil)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &user.UserGetPolicyData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.NotEmpty(suite.T(), data.Rules)
}

func (suite *UserGetPolicyTestSuite) TestWrongAccessToken() {
	req := httptest.NewRequest("GET", "/api/user/policy", nil)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(suite.T(), response.Unauthorized, respBody.Message)
	assert.Equal(suite.T(), response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *UserGetPolicyTestSuite) TearDownTest() {
	test.RdbmsMigration.Reset()
}

func TestUserGetPolicyTestSuite(t *testing.T) {
	suite.Run(t, new(UserGetPolicyTestSuite))
}
