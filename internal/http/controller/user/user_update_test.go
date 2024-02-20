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
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserUpdateTestSuite struct {
	suite.Suite
}

func (suite *UserUpdateTestSuite) SetupSuite() {
	test.RdbmsMigration.Run()
	test.UserRegister()
}

func (suite *UserUpdateTestSuite) TestUpdate() {
	accessToken := test.UserLogin()
	userUpdateRequest := user.UserUpdateRequest{
		Name:  "bob",
		Email: "bob@test.com",
	}
	reqBody, err := json.Marshal(userUpdateRequest)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", "/api/user", bytes.NewReader(reqBody))
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

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.NotEmpty(suite.T(), data.Id)
	assert.Equal(suite.T(), userUpdateRequest.Name, data.Name)
	assert.Equal(suite.T(), userUpdateRequest.Email, data.Email)
	assert.NotEmpty(suite.T(), data.CreatedAt)
	assert.NotEmpty(suite.T(), data.UpdatedAt)
}

func (suite *UserUpdateTestSuite) TestWrongAccessToken() {
	req := httptest.NewRequest("PUT", "/api/user", nil)
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

func (suite *UserUpdateTestSuite) TestCsrfMismatch() {
	req := httptest.NewRequest("PUT", "/api/user", nil)
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

func (suite *UserUpdateTestSuite) TestRequestValidationFailed() {
	accessToken := test.UserLogin()
	req := httptest.NewRequest("PUT", "/api/user", bytes.NewReader([]byte{}))
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

func (suite *UserUpdateTestSuite) TearDownSuite() {
	test.RdbmsMigration.Reset()
}

func TestUserUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UserUpdateTestSuite))
}
