package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChanJuiHuang/go-backend-framework/v2/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/v2/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/v2/internal/test"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"
)

type UserRegisterTestSuite struct {
	suite.Suite
}

func (suite *UserRegisterTestSuite) SetupSuite() {
	test.RdbmsMigration.Run()
}

func (suite *UserRegisterTestSuite) Test() {
	reqBody := user.UserRegisterRequest{
		Name:     "bob",
		Email:    "bob@test.com",
		Password: "abcABC123",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(reqBodyBytes))
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

func (suite *UserRegisterTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader([]byte{}))
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

func (suite *UserRegisterTestSuite) TestRequestValidationFailed() {
	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader([]byte{}))
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

func (suite *UserRegisterTestSuite) TearDownSuite() {
	test.RdbmsMigration.Reset()
}

func TestUserRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(UserRegisterTestSuite))
}
