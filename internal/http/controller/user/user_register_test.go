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

type UserRegisterTestSuite struct {
	suite.Suite
}

func (suite *UserRegisterTestSuite) SetupSuite() {
	test.Migration.Run()
}

func (suite *UserRegisterTestSuite) TestRegister() {
	userRegisterRequest := user.UserRegisterRequest{
		Name:     "bob",
		Email:    "bob@test.com",
		Password: "abcABC123",
	}
	reqBody, err := json.Marshal(userRegisterRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(reqBody))
	test.AddCsrfToken(req)
	resp := httptest.NewRecorder()
	test.HttpHandler.ServeHTTP(resp, req)

	respBody := &response.Response{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		panic(err)
	}

	data := &user.UserRegisterData{}
	if err := mapstructure.Decode(respBody.Data, data); err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.NotEmpty(suite.T(), data.AccessToken)
	assert.NotEmpty(suite.T(), data.RefreshToken)
}

func (suite *UserRegisterTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader([]byte{}))
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

func (suite *UserRegisterTestSuite) TestRequestValidationFailed() {
	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader([]byte{}))
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

func (suite *UserRegisterTestSuite) TearDownSuite() {
	test.Migration.Reset()
}

func TestUserRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(UserRegisterTestSuite))
}
