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
	"github.com/stretchr/testify/suite"
)

type UserMeTestSuite struct {
	suite.Suite
}

func (suite *UserMeTestSuite) SetupSuite() {
	test.RdbmsMigration.Run()
	test.UserService.Register()
}

func (suite *UserMeTestSuite) TestMe() {
	accessToken := test.UserService.Login()
	req := httptest.NewRequest("GET", "/api/user/me", nil)
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &user.UserData{}
	decoder := service.Registry.Get("mapstructureDecoder").(func(any, any) error)
	if err := decoder(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.Id)
	suite.NotEmpty(data.Name)
	suite.NotEmpty(data.Email)
	suite.NotEmpty(data.CreatedAt)
	suite.NotEmpty(data.UpdatedAt)
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

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Equal(response.Unauthorized, respBody.Message)
	suite.Equal(response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *UserMeTestSuite) TearDownSuite() {
	test.RdbmsMigration.Reset()
}

func TestUserMeTestSuite(t *testing.T) {
	suite.Run(t, new(UserMeTestSuite))
}
