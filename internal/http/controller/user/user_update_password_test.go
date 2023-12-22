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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserUpdatePasswordTestSuite struct {
	suite.Suite
}

func (suite *UserUpdatePasswordTestSuite) SetupSuite() {
	test.Migration.Run()
	test.UserRegister()
}

func (suite *UserUpdatePasswordTestSuite) TestUpdatePassword() {
	accessToken := test.UserLogin()
	userUpdateRequest := user.UserUpdatePasswordRequest{
		Password:        "abcABC123",
		ConfirmPassword: "abcABC123",
	}
	reqBody, err := json.Marshal(userUpdateRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", "/api/user/password", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusNoContent, resp.Code)
}

func (suite *UserUpdatePasswordTestSuite) TestWrongAccessToken() {
	req := httptest.NewRequest("PUT", "/api/user/password", nil)
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

func (suite *UserUpdatePasswordTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("PUT", "/api/user/password", nil)
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

func (suite *UserUpdatePasswordTestSuite) TestRequestValidationFailed() {
	accessToken := test.UserLogin()
	userUpdateRequest := user.UserUpdatePasswordRequest{
		Password:        "abcABC123",
		ConfirmPassword: "abcABC",
	}
	reqBody, err := json.Marshal(userUpdateRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", "/api/user/password", bytes.NewReader(reqBody))
	test.AddBearerToken(req, accessToken)
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

func (suite *UserUpdatePasswordTestSuite) TearDownSuite() {
	test.Migration.Reset()
}

func TestUserUpdatePasswordTestSuite(t *testing.T) {
	suite.Run(t, new(UserUpdatePasswordTestSuite))
}
