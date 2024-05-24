package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/test"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"
)

type UserLoginTestSuite struct {
	suite.Suite
}

func (suite *UserLoginTestSuite) SetupSuite() {
	test.RdbmsMigration.Run()
	test.UserService.Register()
}

func (suite *UserLoginTestSuite) Test() {
	reqBody := user.UserLoginRequest{
		Email:    test.UserService.User.Email,
		Password: test.UserService.UserPassword,
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &user.TokenData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusOK, resp.Code)
	suite.NotEmpty(data.AccessToken)
}

func (suite *UserLoginTestSuite) TestEmailIsWrong() {
	reqBody := user.UserLoginRequest{
		Email:    "john123@test.com",
		Password: "123456",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal(response.EmailIsWrong, respBody.Message)
	suite.Equal(response.MessageToCode[response.EmailIsWrong], respBody.Code)
}

func (suite *UserLoginTestSuite) TestPasswordIsWrong() {
	reqBody := user.UserLoginRequest{
		Email:    test.UserService.User.Email,
		Password: "abc123",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal(response.PasswordIsWrong, respBody.Message)
	suite.Equal(response.MessageToCode[response.PasswordIsWrong], respBody.Code)
}

func (suite *UserLoginTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader([]byte{}))
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusForbidden, resp.Code)
	suite.Equal(response.Forbidden, respBody.Message)
	suite.Equal(response.MessageToCode[response.Forbidden], respBody.Code)
}

func (suite *UserLoginTestSuite) TestRequestValidationFailed() {
	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader([]byte{}))
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal(response.RequestValidationFailed, respBody.Message)
	suite.Equal(response.MessageToCode[response.RequestValidationFailed], respBody.Code)
}

func (suite *UserLoginTestSuite) TearDownSuite() {
	test.RdbmsMigration.Reset()
}

func TestUserLoginTestSuite(t *testing.T) {
	suite.Run(t, new(UserLoginTestSuite))
}
