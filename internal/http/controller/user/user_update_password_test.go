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
	"github.com/stretchr/testify/suite"
)

type UserUpdatePasswordTestSuite struct {
	suite.Suite
}

func (suite *UserUpdatePasswordTestSuite) SetupSuite() {
	test.RdbmsMigration.Run()
	test.UserService.Register()
}

func (suite *UserUpdatePasswordTestSuite) Test() {
	accessToken := test.UserService.Login()
	reqBody := user.UserUpdatePasswordRequest{
		CurrentPassword: "abcABC123",
		Password:        "abcABC000",
		ConfirmPassword: "abcABC000",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", "/api/user/password", bytes.NewReader(reqBodyBytes))
	test.AddCsrfToken(req)
	test.AddBearerToken(req, accessToken)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	suite.Equal(http.StatusNoContent, resp.Code)
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

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Equal(response.Unauthorized, respBody.Message)
	suite.Equal(response.MessageToCode[response.Unauthorized], respBody.Code)
}

func (suite *UserUpdatePasswordTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("PUT", "/api/user/password", nil)
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

func (suite *UserUpdatePasswordTestSuite) TestRequestValidationFailed() {
	accessToken := test.UserService.Login()
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

	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Equal(response.RequestValidationFailed, respBody.Message)
	suite.Equal(response.MessageToCode[response.RequestValidationFailed], respBody.Code)
}

func (suite *UserUpdatePasswordTestSuite) TearDownSuite() {
	test.RdbmsMigration.Reset()
}

func TestUserUpdatePasswordTestSuite(t *testing.T) {
	suite.Run(t, new(UserUpdatePasswordTestSuite))
}
