package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserMeTestSuite struct {
	suite.Suite
}

func (suite *UserMeTestSuite) SetupSuite() {
	test.Migration.Run()
	test.UserRegister()
}

func (suite *UserMeTestSuite) TestMe() {
	accessToken := test.UserLogin()
	req := httptest.NewRequest("GET", "/api/user/me", nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &user.UserMeData{}
	decoder := service.Registry.Get("mapstructureDecoder").(func(any, any) error)
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.NotEmpty(suite.T(), data.Id)
	assert.NotEmpty(suite.T(), data.Name)
	assert.NotEmpty(suite.T(), data.Email)
	assert.NotEmpty(suite.T(), data.CreatedAt)
	assert.NotEmpty(suite.T(), data.UpdatedAt)
}

func (suite *UserMeTestSuite) TestWrongAccessToken() {
	req := httptest.NewRequest("GET", "/api/user/me", nil)
	test.AddCsrfToken(req)
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

func (suite *UserMeTestSuite) TearDownSuite() {
	test.Migration.Reset()
}

func TestUserMeTestSuite(t *testing.T) {
	suite.Run(t, new(UserMeTestSuite))
}
