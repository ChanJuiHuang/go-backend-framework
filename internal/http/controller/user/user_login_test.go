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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserLoginTestSuite struct {
	suite.Suite
}

func (suite *UserLoginTestSuite) SetupSuite() {
	test.Migration.Run()
	test.UserRegister()
}

func (suite *UserLoginTestSuite) TestLogin() {
	userLoginRequest := user.UserLoginRequest{
		Email:    "john@test.com",
		Password: "abcABC123",
	}
	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &user.UserLoginData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.NotEmpty(suite.T(), data.AccessToken)
	assert.NotEmpty(suite.T(), data.RefreshToken)
}

func (suite *UserLoginTestSuite) TestEmailIsWrong() {
	userLoginRequest := user.UserLoginRequest{
		Email:    "john123@test.com",
		Password: "123456",
	}
	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(suite.T(), response.EmailIsWrong, respBody.Message)
	assert.Equal(suite.T(), response.MessageToCode[response.EmailIsWrong], respBody.Code)
}

func (suite *UserLoginTestSuite) TestPasswordIsWrong() {
	userLoginRequest := user.UserLoginRequest{
		Email:    "john@test.com",
		Password: "abc123",
	}
	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(suite.T(), response.PasswordIsWrong, respBody.Message)
	assert.Equal(suite.T(), response.MessageToCode[response.PasswordIsWrong], respBody.Code)
}

func (suite *UserLoginTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader([]byte{}))
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.ErrorResponse{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusForbidden, resp.Code)
	assert.Equal(suite.T(), response.Forbidden, respBody.Message)
	assert.Equal(suite.T(), response.MessageToCode[response.Forbidden], respBody.Code)
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

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(suite.T(), response.RequestValidationFailed, respBody.Message)
	assert.Equal(suite.T(), response.MessageToCode[response.RequestValidationFailed], respBody.Code)
}

func (suite *UserLoginTestSuite) TearDownSuite() {
	test.Migration.Reset()
}

func TestUserLoginTestSuite(t *testing.T) {
	suite.Run(t, new(UserLoginTestSuite))
}
